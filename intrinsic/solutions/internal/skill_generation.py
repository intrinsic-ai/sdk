# Copyright 2023 Intrinsic Innovation LLC

"""Contains the implementation of generated skill classes."""

from __future__ import annotations

import collections
import enum
import inspect
import re
import textwrap
from typing import Any, Callable, Optional, Set, Type
import uuid

from google.protobuf import descriptor
from google.protobuf import descriptor_pb2
from google.protobuf import descriptor_pool
from google.protobuf import message
from intrinsic.assets import id_utils
from intrinsic.executive.proto import behavior_call_pb2
from intrinsic.skills.proto import skills_pb2
from intrinsic.solutions import blackboard_value
from intrinsic.solutions import cel
from intrinsic.solutions import provided
from intrinsic.solutions import utils
from intrinsic.solutions.internal import skill_utils

# Typing aliases
# Maps slot name to resource names.
ResourceMap = dict[str, provided.ResourceHandle]


def print_failed_descriptor_infra(
    descriptor_set_message: descriptor_pb2.FileDescriptorSet,
) -> bool:
  """Prints proto message that failed to generate infra for.

  Tries to generate proto infra for the descriptor_set_message and if that
  fails prints the proto message that failed. Otherwise nothing is printed.

  Args:
    descriptor_set_message: The file descriptor set message to try to generate
      infra for.

  Returns:
    True if no error was encountered during generation, False otherwise.
  """
  failed_proto = (
      skill_utils.determine_failed_generate_proto_infra_from_filedescriptorset(
          descriptor_set_message
      )
  )
  if failed_proto:
    print(f"Failed creating message types for {failed_proto}.")
    return False
  return True


class SkillInfoImpl(provided.SkillInfo):
  """Implementation of the SkillInfo interface."""

  _skill_proto: skills_pb2.Skill

  def __init__(self, skill_proto: skills_pb2.Skill):
    """Creates a SkillInfoImpl object from the skill_proto.

    Args:
      skill_proto: the protobuf description of this skill.

    Raises:
      TypeError if the skill_proto does not contain all transitive dependencies
      in skill_proto.parameter_description.parameter_descriptor_fileset.
    """

    self._skill_proto = skill_proto
    # Each SkillInfoImpl class uses its own descriptor pool so that the
    # creation of each SkillBase class is hermetic. Ie., Skill A and Skill B
    # do not incidentally clash over the definition of a proto.
    parameter_description = self._skill_proto.parameter_description
    file_descriptor_set = descriptor_pb2.FileDescriptorSet()
    file_descriptor_set.MergeFrom(
        parameter_description.parameter_descriptor_fileset
    )
    file_descriptor_set.MergeFrom(
        self._skill_proto.return_value_description.descriptor_fileset
    )

    desc_pool, message_classes = None, None
    try:
      desc_pool, message_classes = (
          skill_utils.generate_proto_infra_from_filedescriptorset(
              file_descriptor_set
          )
      )
    except NotImplementedError as e:
      print(
          f"Failed to load proto file descriptors for skill: {skill_proto.id}"
      )
      # Try to "dummy" generate pools, etc. individually to determine, which
      # part failed for a more informative error message.
      if not print_failed_descriptor_infra(
          self._skill_proto.parameter_description.parameter_descriptor_fileset
      ):
        print(
            "Could not generate file descriptor infra for the parameter"
            " description."
        )
      if not print_failed_descriptor_infra(
          self._skill_proto.return_value_description.descriptor_fileset
      ):
        print(
            "Could not generate file descriptor infra for the return value"
            " description."
        )
      print(
          "Were this skill's protos build against a different code base than"
          " the workcell API? An example case is a skill build in the internal"
          " code base, but executed in an externally supplied Jupyter."
          " This can be a direct dependency of the parameter or return value"
          " proto or also an indirect dependency via the contained protos."
      )
      raise e

    self._message_pool: descriptor_pool.DescriptorPool = desc_pool
    self._message_classes: dict[str, Type[message.Message]] = message_classes

    self._field_names: Set[str] = set()
    if self._skill_proto.HasField("parameter_description"):
      self._field_names = set(
          [field.name for field in self.parameter_descriptor().fields]
      )

  @property
  def id(self) -> str:
    return self._skill_proto.id

  @property
  def id_version(self) -> str:
    return self._skill_proto.id_version

  @property
  def version(self) -> str:
    return id_utils.version_from(self._skill_proto.id_version)

  @property
  def skill_name(self) -> str:
    # Use skill ID as ground truth. Don't use the 'skill_name' in the proto
    # which is a display name that might contain, e.g., spaces.
    return _skill_name_from_id(self._skill_proto.id)

  @property
  def package_name(self) -> str:
    if self._skill_proto.package_name:
      return self._skill_proto.package_name

    # Extract from the skill ID if the package name is not explicitly set.
    return _skill_package_name_from_id(self._skill_proto.id)

  @property
  def skill_proto(self) -> skills_pb2.Skill:
    return self._skill_proto

  @property
  def parameter_message_full_name(self) -> str:
    return self._skill_proto.parameter_description.parameter_message_full_name

  @property
  def return_value_message_full_name(self) -> str:
    return (
        self._skill_proto.return_value_description.return_value_message_full_name
    )

  def create_param_message(self) -> message.Message:
    return self._message_classes[
        self._skill_proto.parameter_description.parameter_message_full_name
    ]()

  def create_result_message(self) -> message.Message:
    return self._message_classes[
        self._skill_proto.return_value_description.return_value_message_full_name
    ]()

  def get_param_message_type(self) -> Type[message.Message]:
    return self._message_classes[
        self._skill_proto.parameter_description.parameter_message_full_name
    ]

  def get_result_message_type(self) -> Type[message.Message]:
    return self._message_classes[
        self._skill_proto.return_value_description.return_value_message_full_name
    ]

  def parameter_descriptor(self) -> descriptor.Descriptor:
    return self._message_pool.FindMessageTypeByName(
        self._skill_proto.parameter_description.parameter_message_full_name
    )

  def return_value_descriptor(self) -> descriptor.Descriptor:
    return self._message_pool.FindMessageTypeByName(
        self._skill_proto.return_value_description.return_value_message_full_name
    )

  @property
  def field_names(self) -> Set[str]:
    return self._field_names

  @property
  def message_classes(self) -> dict[str, Type[message.Message]]:
    return self._message_classes

  def get_message_class(
      self, msg_descriptor: descriptor.Descriptor
  ) -> Type[message.Message]:
    return self._message_classes[msg_descriptor.full_name]

  def get_parameter_field_comments(self, full_field_name: str) -> str:
    return textwrap.dedent(
        self._skill_proto.parameter_description.parameter_field_comments[
            full_field_name
        ]
    )

  def get_result_field_comments(self, full_field_name: str) -> str:
    return textwrap.dedent(
        self._skill_proto.return_value_description.return_value_field_comments[
            full_field_name
        ]
    )


def _gen_class_docstring(info: provided.SkillInfo) -> str:
  """Generates the class docstring for a skill class.

  Args:
    info: The skill's SkillInfo.

  Returns:
    The generated Python docstring.
  """
  docstring: list[str] = [f"Skill class for {info.skill_proto.id}.\n"]

  # Expect 80 chars width
  is_first_line = True
  for description_line in textwrap.dedent(
      info.skill_proto.description
  ).splitlines():
    wrapped_lines = textwrap.wrap(description_line, 80)
    # Make sure that an empty line is wrapped to an empty line
    # and not removed. We assume that the skill author intended
    # the extra line break there unless it is the first line.
    if not wrapped_lines and not is_first_line:
      wrapped_lines = [""]
    docstring += wrapped_lines
    is_first_line = False

  return "\n".join(docstring)


def _gen_init_docstring(
    info: provided.SkillInfo,
    compatible_resources: dict[str, provided.ResourceList],
) -> str:
  """Generates the __init__ docstring for a skill class.

  Args:
    info: The skill's SkillInfo.
    compatible_resources: Map from resource slot names to resources suitable for
      that slot. It is used to determine whether a default value can be assigned
      for resource parameters.

  Returns:
    The generated Python docstring.

  Raises:
    NameError: if a parameter name and resource name are the same and even
      disambiguation adding a "_resource" suffix failed.
  """
  docstring: list[str] = [
      f"Initializes an instance of the skill {info.skill_proto.id}."
  ]

  params: list[skill_utils.ParameterInformation] = []
  param_names: list[str] = []

  if info.skill_proto.HasField("parameter_description"):
    param_defaults = info.create_param_message()

    if info.skill_proto.parameter_description.HasField("default_value"):
      info.skill_proto.parameter_description.default_value.Unpack(
          param_defaults
      )

    params = skill_utils.extract_docstring_from_message(
        param_defaults, info.skill_proto.parameter_description
    )
    param_names = [p.name for p in params]

  return_values: list[tuple[str, str]] = []
  if info.skill_proto.HasField("return_value_description"):
    result_defaults = info.create_result_message()

    for field in result_defaults.DESCRIPTOR.fields:
      doc_string = info.get_result_field_comments(field.full_name)
      return_values.append((field.name, doc_string))

    params.append(
        skill_utils.ParameterInformation(
            has_default=False,
            name="return_value_key",
            default=None,
            doc_string=["Blackboard key where to store the return value"],
            message_full_name=None,
            enum_full_name=None,
        )
    )
    param_names.append("return_value_key")

  for slot, selector in sorted(info.skill_proto.resource_selectors.items()):
    param_name = skill_utils.deconflict_param_and_resources(slot, param_names)

    slot_docstring = []
    if len(selector.capability_names) == 1:
      slot_docstring.append(
          f"Resource with capability {selector.capability_names[0]}"
      )
    else:
      slot_docstring.append(
          "Resource having all of the following capabilities:"
      )
      for t in selector.capability_names:
        slot_docstring.append(f"  - {t}")

    slot_compatible_resources = compatible_resources[param_name]
    if len(slot_compatible_resources) == 1:
      default_resource = slot_compatible_resources[
          dir(slot_compatible_resources)[0]
      ]
      slot_docstring.append(f"Default resource: {default_resource.name}")

    params.append(
        skill_utils.ParameterInformation(
            has_default=False,
            name=param_name,
            default=None,
            doc_string=slot_docstring,
            message_full_name=None,
            enum_full_name=None,
        )
    )

  skill_utils.append_used_proto_full_names(info.skill_name, params, docstring)

  if not info.skill_proto.resource_selectors:
    docstring.append("\nThis skill requires no resources.")

  # Generate actual docstring for arguments
  if params:
    docstring.append("\nArgs:")
    params.sort(key=lambda p: (p.has_default, p.name, p.default, p.doc_string))
    for p in params:
      docstring.append(f"    {p.name}:")
      for param_doc_string in p.doc_string:
        # Expect 80 chars width, subtract 8 for leading spaces in args string.
        for line in textwrap.wrap(param_doc_string.strip(), 72):
          docstring.append(f"        {line}")
      if p.has_default:
        docstring.append(f"        Default value: {p.default}")

  if return_values:
    docstring.append("\nReturns:")
    for name, result_doc_string in sorted(return_values):
      docstring.append(f"    {name}:")
      # Expect 80 chars width, subtract 8 for leading spaces in args string.
      for line in textwrap.wrap(result_doc_string, 72):
        docstring.append(f"        {line}")

  return "\n".join(docstring)


def _gen_init_params(
    info: provided.SkillInfo,
    compatible_resources: dict[str, provided.ResourceList],
    wrapper_classes: dict[str, Type[skill_utils.MessageWrapper]],
    enum_classes: dict[str, Type[enum.IntEnum]],
) -> list[inspect.Parameter]:
  """Create argument typing information for a given skill info.

  This iterates over the parameters in the skill info and suggests the most
  pythonic type available, for example, STRING is mapped to Str, STRING_VECTOR
  is mapped to list[str].

  Args:
    info: Skill info to create argument typing info for.
    compatible_resources: Map from resource slot names to resources suitable for
      that slot. It is used to determine whether a default value can be assigned
      for resource parameters.
    wrapper_classes: Map from proto message names to corresponding message
      wrapper classes.
    enum_classes: Map from full enum names to corresponding enum classes.

  Returns:
    A dict mapping from field name to pythonic type.

  Raises:
    NameError: a parameter name and resource name are the same and even
      disambiguation adding a "_resource" suffix failed.
  """
  params: list[inspect.Parameter] = []
  param_names: list[str] = []

  if info.skill_proto.HasField("parameter_description"):
    param_defaults = info.create_param_message()

    if info.skill_proto.parameter_description.HasField("default_value"):
      info.skill_proto.parameter_description.default_value.Unpack(
          param_defaults
      )

    # Extract those fields from the default parameter proto which may contain
    # actual default values.
    param_info = skill_utils.extract_parameter_information_from_message(
        param_defaults,
        info.skill_proto.parameter_description,
        wrapper_classes,
        enum_classes,
    )
    if param_info:
      params, param_names = map(list, zip(*param_info))

  # Add resource slot names, if a resource slot has the same name as
  # a skill parameter (actually different namespaces) disambiguate by suffixing
  # the slot name with "_resource".
  for slot in sorted(info.skill_proto.resource_selectors.keys()):
    default = inspect.Parameter.empty
    if slot in compatible_resources:
      slot_compatible_resource = compatible_resources[slot]
      if len(slot_compatible_resource) == 1:
        # A resource might be contained twice for the slot in
        # compatible_resources if the name had to be mangled to be compatible as
        # a Python attribute name. Using dir() will filter this for us, so just
        # take the first and only of these entries to assign.
        default = slot_compatible_resource[dir(slot_compatible_resource)[0]]
    params.append(
        inspect.Parameter(
            skill_utils.deconflict_param_and_resources(slot, param_names),
            inspect.Parameter.KEYWORD_ONLY,
            annotation=provided.ResourceHandle,
            default=default,
        )
    )

  if info.skill_proto.HasField("return_value_description"):
    params.append(
        inspect.Parameter(
            "return_value_key",
            inspect.Parameter.KEYWORD_ONLY,
            default=None,
            annotation=Optional[str],
        )
    )

  # Sort items without default arguments before the ones with defaults.
  # This is required to generate valid function signatures.
  params.sort(key=lambda f: f.default == inspect.Parameter.empty, reverse=True)
  return params


def _gen_init_fun(
    info: provided.SkillInfo,
    compatible_resources: dict[str, provided.ResourceList],
    wrapper_classes: dict[str, Type[skill_utils.MessageWrapper]],
    enum_classes: dict[str, Type[enum.IntEnum]],
) -> Callable[[Any, Any], "GeneratedSkill"]:
  """Generate custom __init__ class method with proper auto-completion info.

  Args:
    info: Skill information.
    compatible_resources: Map from resource slot names to resources suitable for
      that slot. It is used to determine whether a default value can be assigned
      for resource parameters.
    wrapper_classes: Map from proto message names to corresponding message
      wrapper classes.
    enum_classes: Map from full enum names to corresponding enum classes.

  Returns:
    A function suitable to be used as __init__ function for a GeneratedSkill
    derivative.
  """

  def new_init_fun(self, **kwargs):
    # We disable the warning because we are generating a sub-class function that
    # can actually access the protected method.
    GeneratedSkill.__init__(self)
    self._resources: ResourceMap = {}
    params_set = self._set_params(**kwargs)  # pylint: disable=protected-access
    resource_set = self._set_resources(**kwargs)  # pylint: disable=protected-access
    return_value_key_set = self._set_return_value_key(**kwargs)  # pylint: disable=protected-access
    # Arguments which are neither skill param, resources nor return_value_key
    extra_args_set = set(kwargs.keys()) - set(
        params_set + resource_set + return_value_key_set
    )
    if extra_args_set:
      raise NameError(f"Unknown argument(s): {', '.join(extra_args_set)}")

  params = [
      inspect.Parameter(
          "self",
          inspect.Parameter.POSITIONAL_OR_KEYWORD,
          annotation="Skill_" + info.skill_name,
      )
  ] + _gen_init_params(
      info, compatible_resources, wrapper_classes, enum_classes
  )
  new_init_fun.__signature__ = inspect.Signature(params)
  new_init_fun.__annotations__ = collections.OrderedDict(
      [(p.name, p.annotation) for p in params]
  )
  new_init_fun.__doc__ = _gen_init_docstring(info, compatible_resources)

  return new_init_fun


def gen_skill_class(
    info: provided.SkillInfo,
    compatible_resources: dict[str, provided.ResourceList],
) -> Type[Any]:
  """This generates a new skill wrapper class type.

  We need to do this because we already need the constructor to pass instance
  information we do not want to spill to user space (skill info), and therefore
  for a nice notation need to overload __init__. In order to be able to augment
  it with meta info for auto-completion, we need to dynamically generate
  it. Since __init__ is a class and not an instance method, we cannot simply
  assign the function, but need to generate an entire type for it.

  Args:
    info: Skill information.
    compatible_resources: Map with compatible resources.

  Returns:
    A new type for a GeneratedSkill sub-class.
  """
  message_classes_to_wrap = {}
  enum_descriptors_to_wrap = {}
  if info.skill_proto.HasField("parameter_description"):
    skill_utils.collect_message_classes_to_wrap(
        info.parameter_descriptor(),
        info.message_classes,
        message_classes_to_wrap,
        enum_descriptors_to_wrap,
    )
  if info.skill_proto.HasField("return_value_description"):
    skill_utils.collect_message_classes_to_wrap(
        info.return_value_descriptor(),
        info.message_classes,
        message_classes_to_wrap,
        enum_descriptors_to_wrap,
    )

  type_class = type(
      # E.g.: 'move_robot'
      info.skill_name,
      (GeneratedSkill,),
      {
          "_info": info,
          "_compatible_resources": provided.SkillCompatibleResourcesMap(
              compatible_resources
          ),
          # We use the __init__ documentation because that is shown in the
          # auto-completion tooltip, not __init__.__doc__.
          "__doc__": _gen_class_docstring(info),
          # E.g.: 'intrinsic.solutions.skills.ai.intrinsic'.
          "__module__": skill_utils.module_for_generated_skill(
              info.package_name
          ),
      },
  )

  wrapper_classes, enum_classes = skill_utils.update_message_class_modules(
      type_class,
      info,
      message_classes_to_wrap,
      enum_descriptors_to_wrap,
  )

  init_fun = _gen_init_fun(
      info, compatible_resources, wrapper_classes, enum_classes
  )
  type_class.__init__ = init_fun

  return type_class


def _field_to_repr(field: descriptor.FieldDescriptor, field_value: Any) -> str:
  """Generates a representation for a field value.

  Args:
    field: proto field descriptor
    field_value: the value to represent

  Returns:
    string representation of value with respect the given descriptor.
  """
  if field.message_type is None:
    return repr(field_value)

  if field.label == descriptor.FieldDescriptor.LABEL_REPEATED:
    if (
        field.type == descriptor.FieldDescriptor.TYPE_MESSAGE
        and field.message_type.GetOptions().map_entry
    ):

      def quote_str(value):
        if isinstance(value, str):
          return f'"{value}"'
        else:
          return value

      item_strs = [
          f"{quote_str(k)}: {repr(v)}" for (k, v) in field_value.items()
      ]
      return f'{{{", ".join(item_strs)}}}'
    else:
      return f'[{", ".join(repr(value) for value in field_value)}]'

  return repr(field_value)


_SKILL_ID_VERSION_REGEX = (
    r"^(?P<id>(?:(?P<package>(?:\D\w*\.)*\D\w*)\.)?(?P<name>\D\w*))$"
)


def _skill_name_from_id(skill_id: str) -> str:
  """Extracts the name from the given skill ID.

  Args:
    skill_id: The skill ID.

  Returns:
    The name extracted from the ID.
  """

  m = re.search(_SKILL_ID_VERSION_REGEX, skill_id)

  if m is not None:
    return m.group("name")

  return skill_id


def _skill_package_name_from_id(skill_id: str) -> str:
  """Extracts the name from the given skill ID.

  Args:
    skill_id: The skill ID.

  Returns:
    The name extracted from the ID.
  """
  m = re.search(_SKILL_ID_VERSION_REGEX, skill_id)

  if m is not None:
    package = m.group("package")
    if package is not None:
      return package

  return ""


class GeneratedSkill(provided.SkillBase):
  """Base class for skill classes dynamically generated at runtime."""

  _info: provided.SkillInfo
  _compatible_resources: provided.SkillCompatibleResourcesMap

  def __init__(self, **kwargs):
    """This constructor normally will not be called.

    It needs to accepts arbitrary args, or we might get errors during type-
    checking for scripts that invoke skills (The skills are not online
    during type checking).

    Args:
     **kwargs: Keyword arguments for skill construction.
    """
    super().__init__()

    self._resources: ResourceMap = {}
    self._param_message: message.Message = None
    self._blackboard_params: dict[str, Any] = {}
    self._return_value_key = (
        self._info.skill_name + "_" + str(uuid.uuid4()).replace("-", "_")
    )

  @property
  def result_key(self) -> str:
    return self._return_value_key

  @property
  def result(self) -> blackboard_value.BlackboardValue:
    msg = self._info.create_result_message()
    return blackboard_value.BlackboardValue(
        msg.DESCRIPTOR.fields_by_name,
        self._return_value_key,
        self._info.get_result_message_type(),
        None,
    )

  def _set_params(self, **kwargs) -> list[str]:
    """Set parameters of skill.

    Args:
      **kwargs: Map from parameter name to value as specified by the skill.
        Unknown arguments are silently ignored.

    Returns:
      List of keys in arguments consumed as parameters.

    Raises:
      TypeError: If passing a value that does not match a field's type
      ValueError: If passing a negative value for a UINT64 field
      KeyError: If failing to provide a value for any skill argument
    """
    params_set = []

    if self._info.skill_proto.HasField("parameter_description"):
      default_message = None
      if self._info.skill_proto.parameter_description.HasField("default_value"):
        default_message = self._info.create_param_message()
        self._info.skill_proto.parameter_description.default_value.Unpack(
            default_message
        )

      self._param_message = self._info.create_param_message()
      if default_message:
        self._param_message.CopyFrom(default_message)

      fields = {}
      for param_name, value in kwargs.items():
        if param_name in self._param_message.DESCRIPTOR.fields_by_name:
          if isinstance(value, blackboard_value.BlackboardValue):
            self._blackboard_params[param_name] = value.value_access_path()
            continue

          if isinstance(value, cel.CelExpression):
            self._blackboard_params[param_name] = str(value)
            continue

          fields[param_name] = value
          if isinstance(value, skill_utils.MessageWrapper):
            for k, v in value.blackboard_params.items():
              self._blackboard_params[param_name + "." + k] = v

          if isinstance(value, list):
            if value and (
                isinstance(value[0], skill_utils.MessageWrapper)
                or isinstance(value[0], blackboard_value.BlackboardValue)
                or isinstance(value[0], cel.CelExpression)
            ):
              # Guard this check for list non-scalar list values to allow
              # assigning a VectorNd from a list[float].
              if (
                  self._param_message.DESCRIPTOR.fields_by_name[
                      param_name
                  ].label
                  != descriptor.FieldDescriptor.LABEL_REPEATED
              ):
                raise TypeError(
                    f"Cannot set field {param_name} to list, not a repeated"
                    " field"
                )
            for i, listelem in enumerate(value):
              if isinstance(listelem, skill_utils.MessageWrapper):
                for k, v in listelem.blackboard_params.items():
                  self._blackboard_params[param_name + f"[{i}]." + k] = v
              elif isinstance(listelem, blackboard_value.BlackboardValue):
                self._blackboard_params[param_name + f"[{i}]"] = (
                    listelem.value_access_path()
                )
              elif isinstance(listelem, cel.CelExpression):
                self._blackboard_params[param_name + f"[{i}]"] = str(listelem)

          if isinstance(value, dict):
            field = self._param_message.DESCRIPTOR.fields_by_name[param_name]
            if (
                field.label != descriptor.FieldDescriptor.LABEL_REPEATED
                or field.type != descriptor.FieldDescriptor.TYPE_MESSAGE
                or not field.message_type.GetOptions().map_entry
            ):
              raise TypeError(
                  f"Cannot set field {param_name} to dict, not a map field"
              )

      params_set = skill_utils.set_fields_in_msg(self._param_message, fields)
      params_set.extend(self._blackboard_params.keys())

    return params_set

  def _set_return_value_key(self, **kwargs) -> list[str]:
    """Set return value key, if provided.

    Args:
      **kwargs: Map from parameter name to value as specified by the skill.
        Unknown arguments are silently ignored.

    Returns:
      List of keys in arguments consumed as parameters.

    Raises:
      TypeError: if passing a value that satisfy the slot's type requirements
    """
    if "return_value_key" in kwargs:
      if kwargs["return_value_key"] is not None:
        self._return_value_key = kwargs["return_value_key"]
      return ["return_value_key"]
    return []

  def _set_resources(self, **kwargs) -> list[str]:
    """Set resource requirements of skill.

    Args:
      **kwargs: Map from parameter name to value as specified by the skill.
        Unknown arguments are silently ignored.

    Returns:
      List of keys in arguments consumed as parameters.

    Raises:
      TypeError: if passing a value that satisfy the slot's type requirements
      KeyError: if failing to provide a value for any skill argument
    """
    resource_set = []

    for slot, selector in self._info.skill_proto.resource_selectors.items():
      slot_param_name = slot
      if slot in self._info.field_names:
        slot_param_name = slot + skill_utils.RESOURCE_SLOT_DECONFLICT_SUFFIX
      compatible_resources = self.compatible_resources[slot_param_name]
      if slot_param_name not in kwargs and len(compatible_resources) != 1:
        if not compatible_resources:
          raise KeyError(
              f"Resource argument '{slot_param_name}' is missing and "
              "no compatible resource has been configured for this "
              "solution."
          )
        raise KeyError(
            f"Resource argument '{slot_param_name}' is missing. "
            "There is more than one compatible resource ("
            f"{', '.join([e.name for e in compatible_resources])})."
        )

      resource = (
          kwargs[slot_param_name]
          if slot_param_name in kwargs and kwargs[slot_param_name] is not None
          else compatible_resources[dir(compatible_resources)[0]]
      )

      if not isinstance(resource, provided.ResourceHandle):
        raise TypeError(
            f"Given value for resource slot '{slot_param_name}' "
            "is not a ResourceHandle."
        )

      given_capabilities = set(resource.types)
      expected_capabilities = set(selector.capability_names)

      if not expected_capabilities.issubset(given_capabilities):
        raise TypeError(
            f"Expected capabilities ({', '.join(expected_capabilities)}) "
            f"for slot {slot}, but given resource handle has "
            f"({', '.join(given_capabilities)})"
        )

      self._resources[slot] = resource
      resource_set.append(slot_param_name)

    return resource_set

  @utils.classproperty
  def info(cls) -> skills_pb2.Skill:  # pylint:disable=no-self-argument
    return cls._info.skill_proto

  @property
  def proto(self) -> behavior_call_pb2.BehaviorCall:
    proto = behavior_call_pb2.BehaviorCall(
        skill_id=self._info.skill_proto.id,
        return_value_name=self._return_value_key,
    )

    if self._execute_timeout:
      proto.skill_execution_options.execute_timeout.FromTimedelta(
          self._execute_timeout
      )
    if self._project_timeout:
      proto.skill_execution_options.project_timeout.FromTimedelta(
          self._project_timeout
      )

    if self._param_message is not None:
      proto.parameters.Pack(
          self._param_message,
          type_url_prefix=skill_utils.type_url_prefix_for_skill(self._info),
      )

    for slot, handle in self._resources.items():
      proto.resources[slot].handle = handle.name

    for name, cel_expression in self._blackboard_params.items():
      proto.assignments.append(
          behavior_call_pb2.BehaviorCall.ParameterAssignment(
              parameter_path=name, cel_expression=cel_expression
          )
      )

    return proto

  def __repr__(self) -> str:
    params = []
    if self._param_message is not None:
      params.extend([
          f"{field.name}={_field_to_repr(field, value)}"
          for field, value in self._param_message.ListFields()
      ])
    resource_params = []
    if self.proto.resources:
      for key, value in sorted(self.proto.resources.items()):
        slot_param_name = key
        if key in self._info.field_names:
          slot_param_name = key + skill_utils.RESOURCE_SLOT_DECONFLICT_SUFFIX
        resource_params.append(
            "{}={{{}}}".format(slot_param_name, repr(value).strip())
        )
    return (
        f"skills.{_skill_name_from_id(self.proto.skill_id)}({', '.join(params+resource_params)})"
    )

  @utils.classproperty
  def compatible_resources(cls) -> provided.SkillCompatibleResourcesMap:  # pylint:disable=no-self-argument
    return cls._compatible_resources

  @utils.classproperty
  def skill_info(cls) -> provided.SkillInfo:  # pylint:disable=no-self-argument
    return cls._info

  @utils.classproperty
  def message_classes(cls) -> dict[str, Type[message.Message]]:  # pylint:disable=no-self-argument
    return cls._info.message_classes
