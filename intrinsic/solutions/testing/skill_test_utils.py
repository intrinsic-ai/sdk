# Copyright 2023 Intrinsic Innovation LLC

"""Utilities for testing skills in the solution building library."""

from typing import Optional
from unittest import mock

from absl import flags
from google.protobuf import descriptor
from google.protobuf import descriptor_pb2
from google.protobuf import message
from google.protobuf import text_format
from intrinsic.resources.client import resource_registry_client
from intrinsic.resources.proto import resource_handle_pb2
from intrinsic.resources.proto import resource_registry_pb2
from intrinsic.skills.client import skill_registry_client
from intrinsic.skills.proto import skill_registry_pb2
from intrinsic.skills.proto import skills_pb2
from intrinsic.solutions.testing import test_skill_params_pb2


FLAGS = flags.FLAGS


def _read_message_from_pbbin_file(filename):
  with open(filename, 'rb') as fileobj:
    return descriptor_pb2.FileDescriptorSet.FromString(fileobj.read())


def _get_test_message_file_descriptor_set(
    file_descriptor_set_pbbin_filename: str,
) -> descriptor_pb2.FileDescriptorSet:
  """Returns the file descriptor set loaded from the given file.

  Requires FLAGS to be parsed prior to invocation.

  Args:
    file_descriptor_set_pbbin_filename: The filename of the file descriptor set
      binary proto file.

  Returns:
    The file descriptor set.
  """
  return _read_message_from_pbbin_file(file_descriptor_set_pbbin_filename)


class SkillTestUtils:
  """Utility class for testing skills in the solution building library.

  Provides helpers for creating SkillProvider instances.
  """

  def __init__(self, file_descriptor_set_pbbin_filename: str):
    """Initializes a new instance.

    Args:
      file_descriptor_set_pbbin_filename: The filename of a file descriptor set
        binary proto file. This file descriptor set will be used for all skills
        that are created with this instance and which have parameters or return
        values.
    """
    self._file_descriptor_set = _get_test_message_file_descriptor_set(
        file_descriptor_set_pbbin_filename
    )

  def create_parameterless_skill_info(self, skill_id: str) -> skills_pb2.Skill:
    """Creates a skill proto for a skill without parameters or return values.

    Args:
      skill_id: The ID of the skill.

    Returns:
      The skill proto.
    """
    id_parts = skill_id.split('.')
    skill_info = skills_pb2.Skill(
        id=skill_id,
        skill_name=id_parts[-1],
        package_name='.'.join(id_parts[:-1]),
    )

    return skill_info

  def create_test_skill_info(
      self,
      skill_id: str,
      parameter_defaults: message.Message,
      resource_selectors: Optional[dict[str, str]] = None,
      skill_version: str | None = None,
  ) -> skills_pb2.Skill:
    """Creates a skill proto for a skill with parameters.

    Args:
      skill_id: The ID of the skill.
      parameter_defaults: The default values for the skill's parameters. The
        type of this message will be used as the skill's parameter message type.
      resource_selectors: A mapping from resource selector names to capability
        names.
      skill_version: An optional version string to generate id_version.

    Returns:
      The skill proto.
    """
    id_parts = skill_id.split('.')
    id_version = None
    if skill_version is not None:
      id_version = f'{skill_id}.{skill_version}'
    skill_info = skills_pb2.Skill(
        id=skill_id,
        id_version=id_version,
        skill_name=id_parts[-1],
        package_name='.'.join(id_parts[:-1]),
    )

    skill_info.parameter_description.parameter_descriptor_fileset.CopyFrom(
        self._file_descriptor_set
    )

    skill_info.parameter_description.default_value.Pack(parameter_defaults)

    skill_info.parameter_description.parameter_message_full_name = (
        parameter_defaults.DESCRIPTOR.full_name
    )

    for field in parameter_defaults.DESCRIPTOR.fields:
      skill_info.parameter_description.parameter_field_comments[
          field.full_name
      ] = 'Mockup comment'

    if resource_selectors:
      for key, value in resource_selectors.items():
        skill_info.resource_selectors[key].capability_names.append(value)

    return skill_info

  def create_test_skill_info_with_return_value(
      self,
      skill_id: str,
      parameter_defaults: message.Message,
      resource_selectors: Optional[dict[str, str]] = None,
      skill_description: Optional[str] = None,
      return_value_defaults: message.Message | None = None,
      skill_version: str | None = None,
  ) -> skills_pb2.Skill:
    """Creates a skill proto for a skill with parameters and return values.

    Args:
      skill_id: The ID of the skill.
      parameter_defaults: The default values for the skill's parameters. The
        type of this message will be used as the skill's parameter message type
        and return value message type.
      resource_selectors: A mapping from resource selector names to capability
        names.
      skill_description: The description of the skill.
      return_value_defaults: Defaults values for the skill's return value. If
        not set the parameter_defaults are taken.
      skill_version: An optional version string to generate id_version.

    Returns:
      The skill proto.
    """
    id_parts = skill_id.split('.')
    id_version = None
    if skill_version is not None:
      id_version = f'{skill_id}.{skill_version}'
    skill_info = skills_pb2.Skill(
        id=skill_id,
        id_version=id_version,
        skill_name=id_parts[-1],
        package_name='.'.join(id_parts[:-1]),
    )

    if skill_description is not None:
      skill_info.description = skill_description

    skill_info.parameter_description.parameter_descriptor_fileset.CopyFrom(
        self._file_descriptor_set
    )

    skill_info.parameter_description.default_value.Pack(parameter_defaults)

    skill_info.parameter_description.parameter_message_full_name = (
        parameter_defaults.DESCRIPTOR.full_name
    )

    skill_info.return_value_description.descriptor_fileset.CopyFrom(
        self._file_descriptor_set
    )

    if return_value_defaults is not None:
      skill_info.return_value_description.return_value_message_full_name = (
          return_value_defaults.DESCRIPTOR.full_name
      )

    # Prevents infinite recursion due to recursive messages
    messages_done: set[str] = set()

    def add_field_comments(descr: descriptor.Descriptor):
      if descr.full_name in messages_done:
        return
      messages_done.add(descr.full_name)
      for field in descr.fields:
        skill_info.parameter_description.parameter_field_comments[
            field.full_name
        ] = 'Mockup comment'
        skill_info.return_value_description.return_value_field_comments[
            field.full_name
        ] = 'Mockup comment'
        if field.message_type is not None:
          add_field_comments(field.message_type)

    add_field_comments(parameter_defaults.DESCRIPTOR)

    if resource_selectors:
      for key, value in resource_selectors.items():
        skill_info.resource_selectors[key].capability_names.append(value)

    return skill_info

  def create_get_skills_response(
      self,
      skill_id: str,
      parameter_defaults: test_skill_params_pb2.TestMessage,
      resource_selectors: Optional[dict[str, str]] = None,
      skill_version: str | None = None,
  ) -> skill_registry_pb2.GetSkillsResponse:
    """Creates a GetSkillsResponse for a skill with parameters.

    Args:
      skill_id: The ID of the skill.
      parameter_defaults: The default values for the skill's parameters. The
        type of this message will be used as the skill's parameter message type.
      resource_selectors: A mapping from resource selector names to capability
        names.
      skill_version: An optional version string to generate id_version.

    Returns:
      The skill proto.
    """
    skill_info = self.create_test_skill_info(
        skill_id, parameter_defaults, resource_selectors, skill_version
    )

    skill_registry_response = skill_registry_pb2.GetSkillsResponse()
    skill_registry_response.skills.add().CopyFrom(skill_info)
    return skill_registry_response

  def create_skill_registry_for_skill_infos(
      self,
      skill_infos: list[skills_pb2.Skill],
  ) -> skill_registry_client.SkillRegistryClient:
    """Creates a mock client for a skill registry with the given skills.

    Args:
      skill_infos: The skills to add to the registry.

    Returns:
      The mock skill registry client.
    """

    skill_registry_stub = mock.MagicMock()
    skill_registry_response = skill_registry_pb2.GetSkillsResponse()
    for info in skill_infos:
      skill_registry_response.skills.add().CopyFrom(info)
    skill_registry_stub.GetSkills.return_value = skill_registry_response
    return skill_registry_client.SkillRegistryClient(skill_registry_stub)

  def create_skill_registry_for_skill_info(
      self,
      skill_info: skills_pb2.Skill,
  ) -> skill_registry_client.SkillRegistryClient:
    """Creates a mock client for a skill registry with the given skill.

    Args:
      skill_info: The skill to add to the registry.

    Returns:
      The mock skill registry client.
    """
    return self.create_skill_registry_for_skill_infos([skill_info])

  def create_resource_registry_with_handles(
      self, handles: list[resource_handle_pb2.ResourceHandle]
  ) -> resource_registry_client.ResourceRegistryClient:
    """Creates a client for a resource registry with the given handles.

    Args:
      handles: The handles to add to the registry.

    Returns:
      The mock resource registry client.
    """
    resource_registry_stub = mock.MagicMock()
    resource_registry_stub.ListResourceInstances.return_value = (
        resource_registry_pb2.ListResourceInstanceResponse(
            instances=[
                resource_registry_pb2.ResourceInstance(
                    name=handle.name, resource_handle=handle
                )
                for handle in handles
            ],
        )
    )
    return resource_registry_client.ResourceRegistryClient(
        resource_registry_stub
    )

  def create_resource_registry_with_single_type_handles(
      self, handles: dict[str, str]
  ) -> resource_registry_client.ResourceRegistryClient:
    """Creates a client for a resource registry with a single handle.

    Args:
      handles: Map from resource handle name to resource handle type.

    Returns:
      The mock resource registry client.
    """
    return self.create_resource_registry_with_handles([
        text_format.Parse(
            f"""name: '{name}'
                    resource_data {{
                      key: '{type_name}'
                    }}""",
            resource_handle_pb2.ResourceHandle(),
        )
        for name, type_name in handles.items()
    ])

  def create_resource_registry_with_single_handle(
      self, name: str, type_name: str
  ) -> resource_registry_client.ResourceRegistryClient:
    """Creates a client for a resource registry with a single handle.

    Args:
      name: The name of the resource.
      type_name: The type name of the resource.

    Returns:
      The mock resource registry client.
    """
    return self.create_resource_registry_with_handles([
        text_format.Parse(
            f"""name: '{name}'
                    resource_data {{
                      key: '{type_name}'
                    }}""",
            resource_handle_pb2.ResourceHandle(),
        )
    ])

  def create_empty_resource_registry(
      self,
  ) -> resource_registry_client.ResourceRegistryClient:
    """Creates a client for an empty resource registry.

    Returns:
      The mock resource registry client.
    """
    return self.create_resource_registry_with_handles([])
