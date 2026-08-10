# Copyright 2023 Intrinsic Innovation LLC

"""Lightweight Python wrapper around the proto builder service."""

from __future__ import annotations

import dataclasses
import re
import typing
from typing import Any
import uuid

from google.protobuf import descriptor_pb2
from google.protobuf import message as protobuf_message
import grpc

from intrinsic.executive.proto import proto_builder_pb2
from intrinsic.executive.proto import proto_builder_pb2_grpc
from intrinsic.solutions import blackboard_value
from intrinsic.solutions import cel
from intrinsic.solutions import errors as solutions_errors
from intrinsic.solutions import provided
from intrinsic.solutions.internal import skill_utils
from intrinsic.util.grpc import error_handling
from intrinsic.util.proto import descriptors

_DEFAULT_PARAM_MSG_NAME = "Params"
_DEFAULT_RETURN_MSG_NAME = "ReturnValue"

# Expected to be filled with a UUID (32 hex digits without '-' separators).
_UNIQUE_MAIN_PROTO_FILENAME_PATTERN = "gen/sbl_%s/node.proto"
# Matches file names generated with the above pattern.
_UNIQUE_MAIN_PROTO_FILENAME_REGEX = re.compile(
    r"gen/sbl_[0-9A-Fa-f]{32}/node.proto"
)

_RESOLVED_DEPENDENCY_FILE = (
    "intrinsic/assets/proto/v1/resolved_dependency.proto"
)
_RESOLVED_DEPENDENCY_TYPE = "intrinsic_proto.assets.v1.ResolvedDependency"

_FIELD_METADATA_FILE = "intrinsic/assets/proto/field_metadata.proto"
_FIELD_METADATA_OPTION = "intrinsic_proto.assets.field_metadata"


def _generate_proto_filename_and_package() -> tuple[str, str]:
  filename = _UNIQUE_MAIN_PROTO_FILENAME_PATTERN % uuid.uuid4().hex
  package = filename.rpartition("/")[0].replace("/", ".")
  return filename, package


def _is_generated_proto_filename(filename: str):
  return re.match(_UNIQUE_MAIN_PROTO_FILENAME_REGEX, filename)


def _is_generated_signature(signature: Signature) -> bool:
  """Indicates whether the signature is a generated signature.

  Returns True if the given signature can be assumed to have been generated with
  the ProtoBuilder class (see ProtoBuilder.create_signature() and
  ProtoBuilder.create_signature_with_args()).
  """

  for file in signature.file_descriptor_set.file:
    if _is_generated_proto_filename(file.name):
      return True
  return False


def _partition_message_full_name(full_name: str) -> tuple[str, str]:
  package, _, name = full_name.rpartition(".")
  if not package:
    raise ValueError(f"{full_name} has an empty package")
  if not name:
    raise ValueError(f"{full_name} has an empty name")
  return package, name


@dataclasses.dataclass(kw_only=True, frozen=True)
class FieldSpecBase:
  """Base class for all field specs.

  Attributes:
    name: Name of the field in the message.
    number: Number of the field in the message. To ensure backwards
      compatibility and follow Protobuf best practices always use new field
      numbers and never reuse field numbers, e.g., of removed fields.
    description: Description of the field. Can be a multi-line string.
  """

  name: str
  number: int
  description: str = ""

  def __new__(cls, *args, **kwargs):
    del args, kwargs  # unused
    if cls == FieldSpecBase:
      raise TypeError("Cannot instantiate abstract class.")
    return super().__new__(cls)

  def __post_init__(self):
    if self.number is not None and self.number <= 0:
      raise ValueError("Field numbers must be positive")


@dataclasses.dataclass(kw_only=True, frozen=True)
class FieldSpec(FieldSpecBase):
  """Specifies a singular or repeated field in a proto message.

  Attributes:
    type: Proto type name as it would appear in a .proto file. Can be a scalar
      value type (see https://protobuf.dev/programming-guides/proto3/#scalar) or
      the full name of a well-known message type (see
      ProtoBuilder.get_well_known_types()).
    repeated: True if the field should be repeated (cannot be combined with
      'optional').
    optional: True if the field should be optional (cannot be combined with
      'repeated').
    unit: Unit of the field (e.g., "m" for meters or "rad" for radian). If
      empty, the field is treated as unitless.
    is_icon2_position_part: True if the field holds an Icon2PositionPart
      reference. This can be used to designate a particular field in the UI that
      should be used to select a part name.
    is_installed_scene_object_asset: True if a field of type
      "intrinsic_proto.assets.Id" holds an installed scene object asset. This
      can be used to designate a particular field in the UI that should be used
      to select a scene object asset.
    arg: Value or parameter assignment to be applied to this field (see
      Signature.with_args() for a of list accepted values). Passing arguments is
      only applicable in some contexts, in other contexts setting this will
      cause an error.
  """

  type: str
  repeated: bool = False
  optional: bool = False

  # Attributes for field metadata
  unit: str = ""
  is_icon2_position_part: bool = False
  is_installed_scene_object_asset: bool = False

  arg: Any | provided.ParamAssignment | None = None

  def __post_init__(self):
    super().__post_init__()
    if self.repeated and self.optional:
      raise ValueError(
          "A field cannot be repeated and optional at the same time"
      )
    if any(c.isspace() for c in self.type):
      raise ValueError(f"Field type '{self.type}' cannot contain whitespace")
    if self.type.startswith("map<"):
      raise ValueError(
          f"FieldSpec cannot be used for map type '{self.type}'. "
          "Use MapFieldSpec instead."
      )
    if self.type == _RESOLVED_DEPENDENCY_TYPE:
      raise ValueError(
          f"FieldSpec cannot be used for {_RESOLVED_DEPENDENCY_TYPE}. "
          "Use DependencySpec instead."
      )
    if (
        self.is_installed_scene_object_asset
        and self.type != "intrinsic_proto.assets.Id"
    ):
      raise ValueError(
          "'is_installed_scene_object_asset' can only be set for fields of type"
          " 'intrinsic_proto.assets.Id' (type is '{self.type}')"
      )

  @property
  def has_field_metadata(self) -> bool:
    return (
        bool(self.unit)
        or self.is_icon2_position_part
        or self.is_installed_scene_object_asset
    )


@dataclasses.dataclass(kw_only=True, frozen=True)
class MapFieldSpec(FieldSpecBase):
  """Specifies a map field in a proto message.

  Attributes:
    key_type: Proto type name for the key as it would appear in a .proto file.
      Must be a scalar value type except float, double, and bytes (see
      https://protobuf.dev/programming-guides/proto3/#scalar).
    value_type: Proto type name for the value as it would appear in a .proto
      file. Can be a scalar value type (see
      https://protobuf.dev/programming-guides/proto3/#scalar) or the full name
        of a well-known message type (see ProtoBuilder.get_well_known_types()).
    arg: Value or parameter assignment to be applied to this field (see
      Signature.with_args() for a of list accepted values). Passing arguments is
      only applicable in some contexts, in other contexts setting this will
      cause an error.
  """

  key_type: str
  value_type: str
  arg: Any | provided.ParamAssignment | None = None

  def __post_init__(self):
    super().__post_init__()
    if any(c.isspace() for c in self.key_type):
      raise ValueError(f"Key type '{self.key_type}' cannot contain whitespace")
    if any(c.isspace() for c in self.value_type):
      raise ValueError(
          f"Value type '{self.value_type}' cannot contain whitespace"
      )


@dataclasses.dataclass(kw_only=True, frozen=True)
class DependencySpec(FieldSpecBase):
  """Specifies a ResolvedDependency field in a proto message.

  Attributes:
    repeated: True if the field should be repeated (cannot be combined with
      'optional').
    optional: True if the field should be optional (cannot be combined with
      'repeated').
    requires: List of interface URIs which the dependency (instance or Asset)
      must provide. For example
      "grpc://intrinsic_proto.motion_planning.MotionPlannerService" or
      "data://intrinsic_proto.perception.v1.PoseEstimationConfig".
    requires_object: True if an an object from the world is required to satisfy
      this dependency.
  """

  repeated: bool = False
  optional: bool = False

  # Attributes for dependency requirements
  requires: list[str] = dataclasses.field(default_factory=lambda: [])
  requires_object: bool = False

  def __post_init__(self):
    super().__post_init__()
    if self.repeated and self.optional:
      raise ValueError(
          "A dependency field cannot be repeated and optional at the same time"
      )
    if not self.requires and not self.requires_object:
      raise ValueError(
          "A dependency field needs to have at least one requirement (via"
          " 'requires' or 'requires_object')"
      )


@dataclasses.dataclass(kw_only=True, frozen=True)
class MessageSpec:
  """Specifies a proto message.

  Attributes:
    fields: The fields in the proto message.
  """

  fields: list[FieldSpec | MapFieldSpec | DependencySpec]


class Signature:
  """Represents the signature of a skill, process or Python script node.

  Attributes:
    parameter_message_full_name: Full type name of the parameter message. Empty
      means "no parameters".
    return_value_message_full_name: Full type name of the return value message.
      Empty means "no return value".
    file_descriptor_set: File descriptor set containing a main file which
      defines the parameter and return value messages and all required
      dependencies. Always set but can be empty (=list of file descriptors is
      empty) if there is no parameter and no return value message.
  """

  parameter_message_full_name: str
  return_value_message_full_name: str
  file_descriptor_set: descriptor_pb2.FileDescriptorSet

  def __init__(
      self,
      *,
      parameter_message_full_name: str = "",
      return_value_message_full_name: str = "",
      file_descriptor_set: descriptor_pb2.FileDescriptorSet | None = None,
  ):
    if parameter_message_full_name or return_value_message_full_name:
      if file_descriptor_set is None:
        raise ValueError(
            "'file_descriptor_set' is required if 'parameter_message_full_name'"
            " or 'return_value_message_full_name' is set"
        )
      self.file_descriptor_set = file_descriptor_set
    else:
      if file_descriptor_set is not None and file_descriptor_set.file:
        raise ValueError(
            "'file_descriptor_set' non-empty but neither"
            " 'parameter_message_full_name' nor"
            " 'return_value_message_full_name' is set"
        )
      self.file_descriptor_set = descriptor_pb2.FileDescriptorSet()
    self.parameter_message_full_name = parameter_message_full_name
    self.return_value_message_full_name = return_value_message_full_name

  def with_args(self, **kwargs) -> SignatureWithArgs:
    """Adds a parameter message and/or parameter assignments to the signature.

    The parameter message and/or parameter assignments are created from the
    given key-value pairs (similar to skill constructors). The keys are the
    field names of the parameter message and each value can be one of:

     - A Python value matching the scalar type of the corresponding field
     - A proto message matching the message type of the corresponding field
     - A provided.ParamAssignment (= BlackboardValue or CelExpression) for
       dynamically assigning a value from the blackboard during process
       execution

    Args:
      **kwargs: Arguments matching the fields of the parameter message.

    Returns:
      A SignatureWithArgs instance with the applied arguments.
    """
    if not self.parameter_message_full_name:
      if kwargs:
        raise ValueError(
            "Cannot apply arguments to a signature without parameters"
        )
      return SignatureWithArgs(
          signature=self,
          params_message=None,
          blackboard_params={},
      )

    params_message = skill_utils.create_message_from_file_descriptor_set(
        self.file_descriptor_set, self.parameter_message_full_name
    )
    blackboard_params = {}
    consumed = skill_utils.set_params(params_message, blackboard_params, kwargs)

    unconsumed = set(kwargs) - set(consumed)
    if unconsumed:
      # A TypeError on unsupported kwargs is idiomatic Python.
      raise TypeError(
          f"Unexpected keyword arguments: {', '.join(sorted(unconsumed))}"
      )

    return SignatureWithArgs(
        signature=self,
        params_message=params_message,
        blackboard_params=blackboard_params,
    )

  def create_params_message(self, **kwargs) -> protobuf_message.Message:
    """Creates a parameter message for the signature.

    The parameter message is created from the given key-value pairs (similar to
    skill constructors). The keys are the field names of the parameter message
    and each value can be one of:

     - A Python value matching the scalar type of the corresponding field
     - A proto message matching the message type of the corresponding field

    The values must be fixed values, so provided.ParamAssignment values (=
    BlackboardValue or CelExpression) for dynamically assigning a value from the
    blackboard are not supported.

    This method can be used, e.g., to create a fixed parameter message for
    testing a behavior tree with parameters:

    ```
    params = behavior_tree.get_signature().create_params_message(x=42)
    solution.executive.run(behavior_tree, parameters=params)
    ```

    Args:
      **kwargs: Arguments matching the fields of the parameter message.

    Returns:
      A protobuf message populated with the given key-value pairs.

    Raises:
      ValueError: If the signature has no parameter message, or if kwargs
        contains parameter assignments.
      TypeError: If kwargs contains keys that do not match fields of the
        parameter message.
    """
    if not self.parameter_message_full_name:
      raise ValueError(
          "Cannot create parameter message for a signature without parameters"
      )

    params_message = skill_utils.create_message_from_file_descriptor_set(
        self.file_descriptor_set, self.parameter_message_full_name
    )

    # We technically don't need to create blackboard_params but it allows us to
    # raise a specific error below if the user tries to pass a parameter
    # assignment.
    blackboard_params: dict[str, str] = {}
    consumed = skill_utils.set_params(params_message, blackboard_params, kwargs)

    unconsumed = set(kwargs) - set(consumed)
    if unconsumed:
      # A TypeError on unsupported kwargs is idiomatic Python.
      raise TypeError(
          f"Unexpected keyword arguments: {', '.join(sorted(unconsumed))}"
      )

    if blackboard_params:
      param_assignment_types = " or ".join(
          cls.__name__ for cls in typing.get_args(provided.ParamAssignment)
      )
      raise ValueError(
          f"Dynamic parameter assignments ({param_assignment_types}) cannot be"
          " used to create a fixed params message. Got parameter assignments"
          f" for paths: {', '.join(sorted(blackboard_params))}"
      )

    return params_message

  def create_return_value_expression(self, **kwargs) -> cel.CelExpression:
    """Creates a CEL expression for the return value message.

    Creates a CEL expression which generates an instance of the return value
    message from values on the blackboard (e.g. skill results). This is useful
    for setting 'BehaviorTree.return_value_expression' for a behavior tree with
    a return value. Example:

    ```
    sig = solution.proto_builder.create_signature(
        return_value=pb.MessageSpec(fields=[
            pb.FieldSpec(type="string", name="str_out", number=1),
            pb.FieldSpec(
                type="intrinsic_proto.Pose", name="poses_out", number=2,
                repeated=True
            ),
    ]))

    skill = solution.skills.ai.intrinsic.my_skill(...)
    tree = bt.BehaviorTree(..., root=skill)
    tree.set_asset_metadata(...)
    tree.set_signature(sig)

    tree.return_value_expression = sig.create_return_value_expression(
      str_out = my_skill.result.string_result,
      poses_out = [my_skill.result.poses[3], my_skill.result.another_pose],
    )
    ```

    Args:
      **kwargs: Mapping from field name to CEL expression (or equivalent). The
        keys must be field names of the return value message specified by the
        signature. Each value can be a CelExpression, BlackboardValue, a plain
        string representing a CEL expression or a list thereof.

    Returns:
      A CEL expression which generates an instance of the return value message
      from the given CEL expressions.

    Raises:
      ValueError: If the signature does not specify a return value message.
      TypeError: If kwargs contains keys that do not match fields of the return
        value message, or if any value is not a CelExpression (or equivalent).
    """
    if not self.return_value_message_full_name:
      raise ValueError("Signature has no return value message")

    pool = descriptors.create_descriptor_pool(self.file_descriptor_set)
    descriptor = pool.FindMessageTypeByName(self.return_value_message_full_name)
    unexpected = set(kwargs) - set(descriptor.fields_by_name)
    if unexpected:
      raise TypeError(
          "Unexpected keyword arguments (not matching field names in the"
          f" return value message): {', '.join(sorted(unexpected))}"
      )

    # Converts 'value' to a CEL expression, recursing into lists.
    def to_expr(value: Any, key: str) -> str:
      if isinstance(
          value, (blackboard_value.BlackboardValue, cel.CelExpression, str)
      ):
        return str(value)
      elif isinstance(value, list):
        return f"[{', '.join(to_expr(item, key) for item in value)}]"
      else:
        raise TypeError(
            "Expected BlackboardValue, CelExpression, str, or list thereof,"
            f" got {type(value).__name__} for '{key}'"
        )

    field_exprs = [
        f"{key}: {to_expr(value, key)}" for key, value in kwargs.items()
    ]

    return cel.CelExpression(
        f"{self.return_value_message_full_name}{{{', '.join(field_exprs)}}}"
    )


@dataclasses.dataclass(frozen=True)
class SignatureWithArgs:
  """Represents a signature with concrete arguments or parameter assignments.

  Attributes:
    signature: The signature itself.
    params_message: The parameter message containing all concretely assigned
      values. None if the signature does not define a parameter message.
    blackboard_params: Parameter assignments to be applied dynamically during
      process execution, specified as a mapping from paths (into the parameter
      message) to CEL expressions.
  """

  signature: Signature
  params_message: protobuf_message.Message | None
  blackboard_params: dict[str, str]

  @property
  def parameter_message_full_name(self) -> str:
    """Returns signature.parameter_message_full_name."""
    return self.signature.parameter_message_full_name

  @property
  def return_value_message_full_name(self) -> str:
    """Returns signature.return_value_message_full_name."""
    return self.signature.return_value_message_full_name

  @property
  def file_descriptor_set(self) -> descriptor_pb2.FileDescriptorSet:
    """Returns signature.file_descriptor_set."""
    return self.signature.file_descriptor_set


class ProtoBuilder:
  """Wrapper for the proto builder gRPC service."""

  _stub: proto_builder_pb2_grpc.ProtoBuilderStub
  _well_known_types_loaded: bool
  _well_known_types: proto_builder_pb2.GetWellKnownTypesResponse
  # Mapping from full type name to file name
  _well_known_types_imports: dict[str, str]

  def __init__(self, stub: proto_builder_pb2_grpc.ProtoBuilderStub):
    """Constructs a new ProtoBuilder object.

    Args:
      stub: The gRPC stub to be used for communication with the service.
    """
    self._stub = stub
    self._well_known_types_loaded = False
    self._well_known_types = proto_builder_pb2.GetWellKnownTypesResponse()
    self._well_known_types_imports = {}

  @classmethod
  def connect(cls, grpc_channel: grpc.Channel) -> ProtoBuilder:
    """Connects to a proto builder for an existing channel.

    Args:
      grpc_channel: Channel to the gRPC service.

    Returns:
      A newly created instance of the wrapper class.

    Raises:
      grpc.RpcError: When gRPC call to service fails.
    """
    stub = proto_builder_pb2_grpc.ProtoBuilderStub(grpc_channel)
    return cls(stub)

  @error_handling.retry_on_grpc_unavailable
  def compile(
      self, proto_filename: str, proto_schema: str
  ) -> descriptor_pb2.FileDescriptorSet:
    """Compiles a proto schema into a FileDescriptorSet proto.

    Args:
      proto_filename: file name to assume for the generated FileDescriptor.
      proto_schema: The schema, e.g., the contents of a .proto file.

    Returns:
      A FileDescriptorSet for the proto_schema.

    Raises:
      grpc.RpcError: When gRPC call fails.
    """
    request = proto_builder_pb2.ProtoCompileRequest(
        proto_filename=proto_filename, proto_schema=proto_schema
    )

    response = self._stub.Compile(request)
    return response.file_descriptor_set

  @error_handling.retry_on_grpc_unavailable
  def compose(
      self,
      proto_filename: str,
      proto_package: str,
      input_descriptors: list[descriptor_pb2.DescriptorProto],
  ) -> descriptor_pb2.FileDescriptorSet:
    """Composes a list of DescriptorProtos into a FileDescriptorSet proto.

    The fields in the input descriptors must point to either native proto types
    or a message contained in the well known types. Use get_well_known_types()
    for a list of the available message types.

    Args:
      proto_filename: file name to assume for the generated FileDescriptor.
      proto_package: the proto package the input_descriptors are in.
      input_descriptors: list of DescriptorProto describing the messages.

    Returns:
      A FileDescriptorSet for the input_descriptors.

    Raises:
      grpc.RpcError: When gRPC call fails.
    """
    request = proto_builder_pb2.ProtoComposeRequest(
        proto_filename=proto_filename,
        proto_package=proto_package,
        input_descriptor=input_descriptors,
    )

    response = self._stub.Compose(request)
    return response.file_descriptor_set

  @error_handling.retry_on_grpc_unavailable
  def _load_well_known_types(self):
    if self._well_known_types_loaded:
      return

    response = self._stub.GetWellKnownTypes(
        proto_builder_pb2.GetWellKnownTypesRequest()
    )
    self._well_known_types = response
    self._well_known_types_imports = {}
    for type_with_version in response.types_with_versions:
      for version in type_with_version.versions:
        self._well_known_types_imports[version.message_full_name] = version.file
    self._well_known_types_loaded = True

  def get_well_known_types(self) -> list[str]:
    """Retrieves a list of well known types.

    Returns:
      A list of full names for the well known types.

    Raises:
      grpc.RpcError: When gRPC call fails.
    """
    self._load_well_known_types()
    return list(self._well_known_types.type_names)

  def create_message(
      self,
      package: str,
      name: str,
      fields: dict[
          str,
          type[int] | type[float] | type[str] | type[bool] | type[bytes] | str,
      ],
  ) -> protobuf_message.Message:
    """Creates a new custom message.

    Example usage:
      create_message('my_pkg', 'MyMessage', {
          'x': float,
          'a': str,
          'object' : 'intrinsic_proto.world.ObjectReference'})

    Args:
      package: The proto package for the message.
      name: The name of the message.
      fields: dict from field name to field type. The type can either be a
        python type, i.e., `int`, `float`, `str`, `bool`, `bytes` or the name of
        a well known type (see get_well_known_types).

    Returns:
      An instance of the new message.
    """
    # 1. Check that all fields refer to built-in types or well known types
    well_known_types = self.get_well_known_types()
    for field_name, field_type in fields.items():
      if isinstance(field_type, str):
        if field_type not in well_known_types:
          raise solutions_errors.InvalidArgumentError(
              f"Field {field_name} with type {field_type} is not a well known"
              " type."
          )
    # 2. Compose a file descriptor set
    proto_filename = package + "_" + name
    proto_filename = proto_filename.replace(".", "_")
    proto_filename = proto_filename + ".proto"
    msg_descriptor = descriptor_pb2.DescriptorProto(name=name)
    for field_number, field in enumerate(fields.items()):
      field_name, field_type = field
      if isinstance(field_type, str):
        msg_descriptor.field.append(
            descriptor_pb2.FieldDescriptorProto(
                name=field_name,
                number=field_number + 1,
                type_name=field_type,
                type=descriptor_pb2.FieldDescriptorProto.TYPE_MESSAGE,
                label=descriptor_pb2.FieldDescriptorProto.LABEL_OPTIONAL,
            )
        )
      elif isinstance(field_type, type):
        # There are multiple proto types mapping to the same python type. Thus
        # when inverting the mapping here we must pick a type, e.g., a python
        # int is TYPE_INT64.
        if field_type == int:
          proto_field_type = descriptor_pb2.FieldDescriptorProto.TYPE_INT64
        elif field_type == float:
          proto_field_type = descriptor_pb2.FieldDescriptorProto.TYPE_DOUBLE
        elif field_type == str:
          proto_field_type = descriptor_pb2.FieldDescriptorProto.TYPE_STRING
        elif field_type == bool:
          proto_field_type = descriptor_pb2.FieldDescriptorProto.TYPE_BOOL
        elif field_type == bytes:
          proto_field_type = descriptor_pb2.FieldDescriptorProto.TYPE_BYTES
        else:
          raise solutions_errors.InvalidArgumentError(
              f"Field {field_name} does not have a supported type:"
              f" {field_type}."
          )
        msg_descriptor.field.append(
            descriptor_pb2.FieldDescriptorProto(
                name=field_name,
                number=field_number + 1,
                type=proto_field_type,
                label=descriptor_pb2.FieldDescriptorProto.LABEL_OPTIONAL,
            )
        )
      else:
        raise solutions_errors.InvalidArgumentError(
            f"For field {field_name}, type {field_type} is not supported"
        )
    file_descriptor_set = self.compose(
        proto_filename, package, [msg_descriptor]
    )

    # 3. Construct the message out of the file descriptor set
    return skill_utils.create_message_from_file_descriptor_set(
        file_descriptor_set, package + "." + name
    )

  def _generate_import_lines(
      self,
      *specs: MessageSpec | None,
  ) -> list[str]:
    used_types: set[str] = set()
    imports_to_add: set[str] = set()
    for spec in specs:
      if spec is not None:
        for field in spec.fields:
          if isinstance(field, FieldSpec):
            if field.has_field_metadata:
              imports_to_add.add(_FIELD_METADATA_FILE)
            used_types.add(field.type)
          elif isinstance(field, MapFieldSpec):
            used_types.add(field.value_type)
          elif isinstance(field, DependencySpec):
            imports_to_add.add(_RESOLVED_DEPENDENCY_FILE)
            imports_to_add.add(_FIELD_METADATA_FILE)

    if used_types:
      self._load_well_known_types()

      for used_type in used_types:
        if used_type in self._well_known_types_imports:
          imports_to_add.add(self._well_known_types_imports[used_type])

    return [f'import "{imp}";' for imp in sorted(imports_to_add)]

  def _message_spec_to_proto_lines(
      self, name: str, spec: MessageSpec
  ) -> list[str]:
    lines = [f"message {name} {{"]
    for field in spec.fields:
      if field.description is not None and field.description.strip():
        for desc_line in field.description.strip().splitlines():
          lines.append(f"  // {desc_line}")

      field_line = "  "
      if isinstance(field, FieldSpec):
        if field.repeated:
          field_line += "repeated "
        elif field.optional:
          field_line += "optional "
        field_line += f"{field.type} {field.name} = {field.number}"
        if field.has_field_metadata:
          field_line += f" [({_FIELD_METADATA_OPTION}) = {{"
          if field.unit:
            field_line += f'unit: "{field.unit}"'
          if field.is_icon2_position_part:
            field_line += " is_icon2_position_part: true"
          if field.is_installed_scene_object_asset:
            field_line += " is_installed_scene_object_asset: true"
          field_line += "}]"
        field_line += ";"
      elif isinstance(field, MapFieldSpec):
        field_line += (
            f"map<{field.key_type}, {field.value_type}> {field.name} ="
            f" {field.number};"
        )
      elif isinstance(field, DependencySpec):
        if field.repeated:
          field_line += "repeated "
        elif field.optional:
          field_line += "optional "
        field_line += (
            f"{_RESOLVED_DEPENDENCY_TYPE} {field.name} = {field.number}"
            f" [({_FIELD_METADATA_OPTION}).dependency = {{"
        )
        for requirement in field.requires:
          field_line += f' requires: "{requirement}"'
        if field.requires_object:
          field_line += " requires_object: {}"
        field_line += "}];"

      lines.append(field_line)
    lines.append("}")
    return lines

  def _proto_schema_for_signature(
      self,
      package_name: str,
      parameters: MessageSpec | None = None,
      return_value: MessageSpec | None = None,
  ) -> str:
    proto_lines = [
        'syntax = "proto3";',
        f"package {package_name};",
        "",
    ]

    proto_lines.extend(self._generate_import_lines(parameters, return_value))
    proto_lines.append("")

    if parameters is not None:
      proto_lines.extend(
          self._message_spec_to_proto_lines(_DEFAULT_PARAM_MSG_NAME, parameters)
      )
      proto_lines.append("")

    if return_value is not None:
      proto_lines.extend(
          self._message_spec_to_proto_lines(
              _DEFAULT_RETURN_MSG_NAME, return_value
          )
      )
      proto_lines.append("")

    return "\n".join(proto_lines)

  def create_signature(
      self,
      *,
      parameters: MessageSpec | None = None,
      return_value: MessageSpec | None = None,
  ) -> Signature:
    """Creates a Signature from parameters and return value specs.

    The FileDescriptorSet of the returned Signature consists of a file
    descriptor containing definitions for the messages specified by 'parameters'
    and 'return_value'. If the given MessageSpecs use well-known types (see
    get_well_known_types()), the corresponding file descriptors are included
    automatically.

    If a given MessageSpec is None, no message definition is generated for it.
    If both MessageSpecs are None, a Signature with an empty file descriptor set
    is returned.

    Args:
      parameters: Specification of the parameters message. If None, no parameter
        message is created.
      return_value: Specification of the return value message. If None, no
        return value message is created.

    Returns:
      A Signature object containing the compiled descriptors and message names.

    Raises:
      ValueError: If the field specs contained in the given message specs
        do specify any args.
    """

    # Check that the passed specs do not define any args
    if parameters is not None:
      for field in parameters.fields:
        if (
            isinstance(field, (FieldSpec, MapFieldSpec))
            and field.arg is not None
        ):
          raise ValueError(
              f"Parameter field '{field.name}' has arg set. Use"
              " create_signature_with_args() to directly create a"
              " SignatureWithArgs."
          )
    if return_value is not None:
      for field in return_value.fields:
        if (
            isinstance(field, (FieldSpec, MapFieldSpec))
            and field.arg is not None
        ):
          raise ValueError(f"Return value field '{field.name}' has arg set")

    if parameters is None and return_value is None:
      return Signature()

    # Generate a unique file and package name.
    proto_filename, package_name = _generate_proto_filename_and_package()
    proto_schema = self._proto_schema_for_signature(
        package_name, parameters, return_value
    )

    file_descriptor_set = self.compile(proto_filename, proto_schema)

    param_full_name = (
        f"{package_name}.{_DEFAULT_PARAM_MSG_NAME}"
        if parameters is not None
        else ""
    )
    return_full_name = (
        f"{package_name}.{_DEFAULT_RETURN_MSG_NAME}"
        if return_value is not None
        else ""
    )

    return Signature(
        parameter_message_full_name=param_full_name,
        return_value_message_full_name=return_full_name,
        file_descriptor_set=file_descriptor_set,
    )

  def create_signature_with_args(
      self,
      *,
      parameters: MessageSpec | None = None,
      return_value: MessageSpec | None = None,
  ) -> SignatureWithArgs:
    """Creates a SignatureWithArgs from parameters and return value specs.

    This is a convenience method for creating a SignatureWithArgs in one step
    instead of using create_signature() followed by Signature.with_args(). The
    desired args can directly be specified in the passed message specs using
    FieldSpec.arg / MapFieldSpec.arg.

    Args:
      parameters: Specification of the parameters message which includes the
        args to be passed. If None, no parameter message is created.
      return_value: Specification of the return value message. Cannot include
        any args (return values cannot have args). If None, no return value
        message is created.

    Returns:
      A SignatureWithArgs object containing the compiled Signature together with
      corresponding parameter assignments.

    Raises:
      ValueError: If the field specs contained in the return_value message spec
        do specify any args.
    """
    # Create a copy of 'parameters' without args, create_signature() does not
    # accept any field specs with args.
    argless_parameters = None
    kwargs = {}
    if parameters is not None:
      argless_fields = []
      for field in parameters.fields:
        if field.arg is not None:
          kwargs[field.name] = field.arg
        argless_fields.append(dataclasses.replace(field, arg=None))
      argless_parameters = dataclasses.replace(
          parameters, fields=argless_fields
      )

    signature = self.create_signature(
        parameters=argless_parameters,
        return_value=return_value,
    )

    return signature.with_args(**kwargs)
