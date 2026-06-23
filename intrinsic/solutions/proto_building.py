# Copyright 2023 Intrinsic Innovation LLC

"""Lightweight Python wrapper around the proto builder service."""

from __future__ import annotations

import dataclasses

from google.protobuf import descriptor_pb2
from google.protobuf import message as protobuf_message
from google.protobuf import message_factory
import grpc

from intrinsic.executive.proto import proto_builder_pb2
from intrinsic.executive.proto import proto_builder_pb2_grpc
from intrinsic.solutions import errors as solutions_errors
from intrinsic.util.grpc import error_handling
from intrinsic.util.proto import descriptors

_DEFAULT_PARAM_MSG_NAME = "Params"
_DEFAULT_RETURN_MSG_NAME = "ReturnValue"
_DEFAULT_SIGNATURE_PROTO_FILENAME = "gen/sbl/signature.proto"


@dataclasses.dataclass(kw_only=True, frozen=True)
class FieldSpec:
  """Specifies a singular or repeated field in a proto message.

  Attributes:
    type: Proto type name as it would appear in a .proto file. Can be a scalar
      value type (see https://protobuf.dev/programming-guides/proto3/#scalar) or
      the full name of a well-known message type (see
      ProtoBuilder.get_well_known_types()).
    name: Name of the field in the message.
    repeated: True if the field should be repeated (cannot be combined with
      'optional').
    optional: True if the field should be optional (cannot be combined with
      'repeated').
    number: Number of the field in the message. To ensure backwards
      compatibility and follow Protobuf best practices always use new field
      numbers and never reuse field numbers, e.g., of removed fields.
    description: Description of the field. Can be a multi-line string.
  """

  type: str
  name: str
  number: int
  repeated: bool = False
  optional: bool = False
  description: str | None = None

  def __post_init__(self):
    if self.repeated and self.optional:
      raise ValueError(
          "A field cannot be repeated and optional at the same time"
      )
    if self.number is not None and self.number <= 0:
      raise ValueError("Field numbers must be positive")
    if any(c.isspace() for c in self.type):
      raise ValueError(f"Field type '{self.type}' cannot contain whitespace")
    if self.type.startswith("map<"):
      raise ValueError(
          f"FieldSpec cannot be used for map type '{self.type}'. "
          "Use MapFieldSpec instead."
      )


@dataclasses.dataclass(kw_only=True, frozen=True)
class MapFieldSpec:
  """Specifies a map field in a proto message.

  Attributes:
    key_type: Proto type name for the key as it would appear in a .proto file.
      Must be a scalar value type except float, double, and bytes (see
      https://protobuf.dev/programming-guides/proto3/#scalar).
    value_type: Proto type name for the value as it would appear in a .proto
      file. Can be a scalar value type (see
      https://protobuf.dev/programming-guides/proto3/#scalar) or the full name
        of a well-known message type (see ProtoBuilder.get_well_known_types()).
    name: Name of the field in the message.
    number: Number of the field in the message. To ensure backwards
      compatibility and follow Protobuf best practices always use new field
      numbers and never reuse field numbers, e.g., of removed fields.
    description: Description of the field. Can be a multi-line string.
  """

  key_type: str
  value_type: str
  name: str
  number: int
  description: str | None = None

  def __post_init__(self):
    if self.number is not None and self.number <= 0:
      raise ValueError("Field numbers must be positive")
    if any(c.isspace() for c in self.key_type):
      raise ValueError(f"Key type '{self.key_type}' cannot contain whitespace")
    if any(c.isspace() for c in self.value_type):
      raise ValueError(
          f"Value type '{self.value_type}' cannot contain whitespace"
      )


@dataclasses.dataclass(kw_only=True, frozen=True)
class MessageSpec:
  """Specifies a proto message.

  Attributes:
    fields: The fields in the proto message.
  """

  fields: list[FieldSpec | MapFieldSpec]


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
    desc_pool = descriptors.create_descriptor_pool(file_descriptor_set)
    message_type = desc_pool.FindMessageTypeByName(package + "." + name)
    assert message_type is not None

    return message_factory.GetMessageClass(message_type)()

  def _generate_import_lines(
      self,
      *specs: MessageSpec | None,
  ) -> list[str]:
    used_types = set()
    for spec in specs:
      if spec is not None:
        for field in spec.fields:
          if isinstance(field, FieldSpec):
            used_types.add(field.type)
          elif isinstance(field, MapFieldSpec):
            used_types.add(field.value_type)

    if not used_types:
      return []

    self._load_well_known_types()
    imports_to_add = set()
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
        field_line += f"{field.type} {field.name} = {field.number};"
      elif isinstance(field, MapFieldSpec):
        field_line += (
            f"map<{field.key_type}, {field.value_type}> {field.name} ="
            f" {field.number};"
        )
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
    """
    if parameters is None and return_value is None:
      return Signature()

    # Use a fixed file and package name. This can be changed to something unique
    # later when using the signature in a behavior tree - we don't know here
    # where and how often the same signature object will be used.
    proto_filename = _DEFAULT_SIGNATURE_PROTO_FILENAME
    package_name = proto_filename.rpartition("/")[0].replace("/", ".")
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
