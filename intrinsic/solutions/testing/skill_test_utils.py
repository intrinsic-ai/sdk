# Copyright 2023 Intrinsic Innovation LLC

"""Utilities for testing skills in the solution building library."""

from typing import cast
from typing import Optional
from unittest import mock

from google.protobuf import any_pb2
from google.protobuf import descriptor
from google.protobuf import descriptor_pb2
from google.protobuf import message
from google.protobuf import text_format

from intrinsic.assets import id_utils
from intrinsic.assets.configuration import asset_configuration_client
from intrinsic.assets.install import installed_assets_client
from intrinsic.assets.processes.proto import process_asset_pb2
from intrinsic.assets.proto import asset_tag_pb2
from intrinsic.assets.proto import asset_type_pb2
from intrinsic.assets.proto import documentation_pb2
from intrinsic.assets.proto import id_pb2
from intrinsic.assets.proto import installed_assets_pb2
from intrinsic.assets.proto import metadata_pb2
from intrinsic.assets.proto import view_pb2
from intrinsic.assets.proto.v1 import asset_configuration_pb2
from intrinsic.executive.proto import behavior_tree_pb2
from intrinsic.resources.client import resource_registry_client
from intrinsic.resources.proto import resource_handle_pb2
from intrinsic.resources.proto import resource_registry_pb2
from intrinsic.skills.client import skill_registry_client
from intrinsic.skills.proto import equipment_pb2
from intrinsic.skills.proto import processed_skill_manifest_pb2
from intrinsic.skills.proto import skill_registry_pb2
from intrinsic.skills.proto import skills_pb2
from intrinsic.util.proto import descriptors


def read_file_descriptor_set(
    file_descriptor_set_pbbin_filename: str,
) -> descriptor_pb2.FileDescriptorSet:
  """Returns the file descriptor set loaded from the given file.

  Args:
    file_descriptor_set_pbbin_filename: The filename of the file descriptor set
      binary proto file.

  Returns:
    The file descriptor set.
  """
  with open(file_descriptor_set_pbbin_filename, 'rb') as fileobj:
    return descriptor_pb2.FileDescriptorSet.FromString(fileobj.read())


class SkillTestUtils:
  """Utility class for testing skills in the solution building library.

  Provides helpers for creating SkillProvider instances.
  """

  def __init__(self, file_descriptor_set_pbbin_filename: str | None = None):
    """Initializes a new instance.

    Args:
      file_descriptor_set_pbbin_filename: The filename of a file descriptor set
        binary proto file. This file descriptor set will be used for all skills
        that are created with this instance and which have parameters or return
        values.
    """
    if file_descriptor_set_pbbin_filename is not None:
      self._file_descriptor_set = read_file_descriptor_set(
          file_descriptor_set_pbbin_filename
      )
    else:
      self._file_descriptor_set = None

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
    id_version = (
        f'{skill_id}.{"0.0.1" if skill_version is None else skill_version}'
    )
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
    id_version = (
        f'{skill_id}.{"0.0.1" if skill_version is None else skill_version}'
    )
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

  def create_empty_skill_registry(
      self,
  ) -> skill_registry_client.SkillRegistryClient:
    """Creates a mock client for an empty skill registry.

    Returns:
      The mock skill registry client.
    """
    return self.create_skill_registry_for_skill_infos([])

  def create_legacy_process(
      self,
      skill_id: str,
      version: str = '0.0.1',
      description: str | None = None,
      parameter_message: type[message.Message] | None = None,
      return_value_message: type[message.Message] | None = None,
      default_params: message.Message | None = None,
  ) -> skills_pb2.Skill:
    """Creates a Skill proto for a legacy process in the skill registry.

    Args:
      skill_id: The skill ID of the process.

    Returns:
      The Skill proto.
    """
    package, _, name = skill_id.rpartition('.')
    display_name = f'Default display name ({name})'
    if description is None:
      description = f'Default description ({name})'
    # If no parameter/result message is specified this will create an empty file
    # descriptor set.
    file_descriptor_set = descriptors.gen_file_descriptor_set([
        cast(descriptor.Descriptor, msg.DESCRIPTOR)
        for msg in (parameter_message, return_value_message)
        if msg is not None
    ])

    skill = skills_pb2.Skill(
        skill_name=name,
        package_name=package,
        id=skill_id,
        id_version=f'{skill_id}.{version}',
        display_name=display_name,
        description=description,
        behavior_tree_description=skills_pb2.BehaviorTreeDescription(),
    )
    if parameter_message is not None:
      skill.parameter_description.parameter_message_full_name = cast(
          descriptor.Descriptor, parameter_message.DESCRIPTOR
      ).full_name
      skill.parameter_description.parameter_descriptor_fileset.CopyFrom(
          file_descriptor_set
      )
      if default_params is not None:
        skill.parameter_description.default_value.Pack(default_params)
    if return_value_message is not None:
      skill.return_value_description.return_value_message_full_name = cast(
          descriptor.Descriptor, return_value_message.DESCRIPTOR
      ).full_name
      skill.return_value_description.descriptor_fileset.CopyFrom(
          file_descriptor_set
      )

    return skill

  def _generate_asset_file_descriptor_set(
      self,
      parameter_message: type[message.Message] | None,
      return_value_message: type[message.Message] | None,
  ) -> tuple[descriptor_pb2.FileDescriptorSet, str | None, str | None]:
    file_descriptor_set = descriptor_pb2.FileDescriptorSet()
    parameter_message_full_name = None
    return_value_message_full_name = None

    if parameter_message is not None or return_value_message is not None:
      file_descriptor_set = descriptors.gen_file_descriptor_set([
          cast(descriptor.Descriptor, msg.DESCRIPTOR)
          for msg in (parameter_message, return_value_message)
          if msg is not None
      ])
      if parameter_message is not None:
        parameter_message_full_name = cast(
            descriptor.Descriptor, parameter_message.DESCRIPTOR
        ).full_name
      if return_value_message is not None:
        return_value_message_full_name = cast(
            descriptor.Descriptor, return_value_message.DESCRIPTOR
        ).full_name

    return (
        file_descriptor_set,
        parameter_message_full_name,
        return_value_message_full_name,
    )

  def create_skill_asset_with_file_descriptor_set(
      self,
      id: str,
      version: str = '0.0.1',
      description: str | None = None,
      file_descriptor_set: descriptor_pb2.FileDescriptorSet | None = None,
      parameter_message_full_name: str | None = None,
      return_value_message_full_name: str | None = None,
      default_params: message.Message | None = None,
      resource_selectors: dict[str, list[str]] | None = None,
  ) -> installed_assets_pb2.InstalledAsset:
    """Creates a configurable Skill asset with a custom file descriptor set.

    Use this method if the tests needs a file descriptor set with source code
    info (i.e., proto comments) which cannot be generated from the compiled-in
    descriptors. Otherwise, prefer create_skill_asset() below which is easier to
    use and leads to more readable tests.
    """
    package = id_utils.package_from(id)
    name = id_utils.name_from(id)
    if description is None:
      description = f'Default description ({name})'
    if file_descriptor_set is None:
      file_descriptor_set = descriptor_pb2.FileDescriptorSet()

    metadata = metadata_pb2.Metadata(
        id_version=id_pb2.IdVersion(
            id=id_pb2.Id(package=package, name=name),
            version=version,
        ),
        display_name=f'Default display name ({name})',
        documentation=documentation_pb2.Documentation(description=description),
        asset_type=asset_type_pb2.ASSET_TYPE_SKILL,
        file_descriptor_set=file_descriptor_set,
    )

    details = processed_skill_manifest_pb2.SkillDetails()
    if parameter_message_full_name is not None:
      details.parameter.message_full_name = parameter_message_full_name
      if default_params is not None:
        details.parameter.default_value.Pack(default_params)
    if return_value_message_full_name is not None:
      details.execute_result.message_full_name = return_value_message_full_name
    if resource_selectors is not None:
      for slot, capability_names in resource_selectors.items():
        details.dependencies.required_equipment[slot].CopyFrom(
            equipment_pb2.ResourceSelector(capability_names=capability_names)
        )

    return installed_assets_pb2.InstalledAsset(
        metadata=metadata,
        skill_specific_metadata=(
            installed_assets_pb2.InstalledAsset.SkillMetadata(details=details)
        ),
    )

  def create_skill_asset(
      self,
      id: str,
      version: str = '0.0.1',
      description: str | None = None,
      parameter_message: type[message.Message] | None = None,
      return_value_message: type[message.Message] | None = None,
      default_params: message.Message | None = None,
      resource_selectors: dict[str, list[str]] | None = None,
  ) -> installed_assets_pb2.InstalledAsset:
    """Creates a configurable Skill asset.

    The file descriptor set gets auto-generated from the compiled-in
    descriptors.
    """
    (
        file_descriptor_set,
        parameter_message_full_name,
        return_value_message_full_name,
    ) = self._generate_asset_file_descriptor_set(
        parameter_message, return_value_message
    )
    return self.create_skill_asset_with_file_descriptor_set(
        id,
        version,
        description,
        file_descriptor_set,
        parameter_message_full_name,
        return_value_message_full_name,
        default_params,
        resource_selectors,
    )

  def create_process_asset_with_file_descriptor_set(
      self,
      id: str,
      version: str = '0.0.1',
      description: str | None = None,
      file_descriptor_set: descriptor_pb2.FileDescriptorSet | None = None,
      parameter_message_full_name: str | None = None,
      return_value_message_full_name: str | None = None,
      default_params: message.Message | None = None,
      resource_selectors: dict[str, list[str]] | None = None,
  ) -> installed_assets_pb2.InstalledAsset:
    """Creates a configurable Process asset with a custom file descriptor set.

    Use this method if the tests needs a file descriptor set with source code
    info (i.e., proto comments) which cannot be generated from the compiled-in
    descriptors. Otherwise, prefer create_process_asset() below which is easier
    to use and leads to more readable tests.
    """
    package = id_utils.package_from(id)
    name = id_utils.name_from(id)
    display_name = f'Default display name ({name})'
    if description is None:
      description = f'Default description ({name})'
    if file_descriptor_set is None:
      file_descriptor_set = descriptor_pb2.FileDescriptorSet()

    metadata = metadata_pb2.Metadata(
        id_version=id_pb2.IdVersion(
            id=id_pb2.Id(package=package, name=name),
            version=version,
        ),
        display_name=display_name,
        documentation=documentation_pb2.Documentation(description=description),
        asset_type=asset_type_pb2.ASSET_TYPE_PROCESS,
        # This file descriptor set can not yet be relied upon for Process
        # assets. Thus in tests we set this to an empty set so that the prod
        # code is forced to use the file descriptor set from the Skill proto
        # below (see b/510731384).
        file_descriptor_set=descriptor_pb2.FileDescriptorSet(),
    )
    skill = skills_pb2.Skill(
        skill_name=name,
        package_name=package,
        id=id,
        id_version=f'{id}.{version}',
        display_name=display_name,
        description=description,
        behavior_tree_description=skills_pb2.BehaviorTreeDescription(),
    )
    if parameter_message_full_name is not None:
      skill.parameter_description.parameter_message_full_name = (
          parameter_message_full_name
      )
      skill.parameter_description.parameter_descriptor_fileset.CopyFrom(
          file_descriptor_set
      )
      if default_params is not None:
        skill.parameter_description.default_value.Pack(default_params)
    if return_value_message_full_name is not None:
      skill.return_value_description.return_value_message_full_name = (
          return_value_message_full_name
      )
      skill.return_value_description.descriptor_fileset.CopyFrom(
          file_descriptor_set
      )
    if resource_selectors is not None:
      for slot, capability_names in resource_selectors.items():
        skill.resource_selectors[slot].CopyFrom(
            equipment_pb2.ResourceSelector(capability_names=capability_names)
        )

    behavior_tree = behavior_tree_pb2.BehaviorTree(
        name=display_name,
        root=behavior_tree_pb2.BehaviorTree.Node(
            sequence=behavior_tree_pb2.BehaviorTree.SequenceNode()
        ),
        description=skill,
    )
    return installed_assets_pb2.InstalledAsset(
        metadata=metadata,
        deployment_data=installed_assets_pb2.InstalledAsset.DeploymentData(
            process=installed_assets_pb2.InstalledAsset.ProcessDeploymentData(
                process=process_asset_pb2.ProcessAsset(
                    behavior_tree=behavior_tree
                )
            )
        ),
    )

  def create_process_asset(
      self,
      id: str,
      version: str = '0.0.1',
      description: str | None = None,
      parameter_message: type[message.Message] | None = None,
      return_value_message: type[message.Message] | None = None,
      default_params: message.Message | None = None,
      resource_selectors: dict[str, list[str]] | None = None,
  ) -> installed_assets_pb2.InstalledAsset:
    """Creates a configurable Process asset.

    The file descriptor set gets auto-generated from the compiled-in
    descriptors.
    """
    (
        file_descriptor_set,
        parameter_message_full_name,
        return_value_message_full_name,
    ) = self._generate_asset_file_descriptor_set(
        parameter_message, return_value_message
    )
    return self.create_process_asset_with_file_descriptor_set(
        id,
        version,
        description,
        file_descriptor_set,
        parameter_message_full_name,
        return_value_message_full_name,
        default_params,
        resource_selectors,
    )

  def create_installed_assets(
      self, assets: list[installed_assets_pb2.InstalledAsset]
  ) -> installed_assets_client.InstalledAssetsClient:
    process_assets = [
        asset
        for asset in assets
        if asset.metadata.asset_type
        == asset_type_pb2.AssetType.ASSET_TYPE_PROCESS
    ]
    skill_assets = [
        asset
        for asset in assets
        if asset.metadata.asset_type
        == asset_type_pb2.AssetType.ASSET_TYPE_SKILL
    ]

    # Only support exactly the calls that are made by the production code.
    def mock_list_all_installed_assets(
        asset_types: list[asset_type_pb2.AssetType],
        asset_tag: asset_tag_pb2.AssetTag | None = None,
        view: view_pb2.AssetViewType | None = None,
    ):
      if (
          asset_types == [asset_type_pb2.AssetType.ASSET_TYPE_PROCESS]
          and asset_tag == asset_tag_pb2.AssetTag.ASSET_TAG_SUBPROCESS
          and view == view_pb2.AssetViewType.ASSET_VIEW_TYPE_FULL
      ):
        return process_assets
      elif (
          asset_types == [asset_type_pb2.AssetType.ASSET_TYPE_SKILL]
          and view == view_pb2.AssetViewType.ASSET_VIEW_TYPE_ALL_METADATA
      ):
        return skill_assets
      else:
        raise ValueError(
            'Unsupported call to mock_list_all_installed_assets'
            f' (asset_types={asset_types}, asset_tag={asset_tag}, view={view})'
        )

    client = mock.MagicMock()
    client.list_all_installed_assets.side_effect = (
        mock_list_all_installed_assets
    )
    return client

  def create_empty_installed_assets(
      self,
  ) -> installed_assets_client.InstalledAssetsClient:
    """Creates a mock client for the installed assets service with no assets."""
    client = mock.MagicMock()
    client.list_all_installed_assets.return_value = []
    return client

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

  def create_asset_configuration_client(
      self,
      recommendations: (
          dict[str, list[tuple[message.Message | None, message.Message]]] | None
      ) = None,
  ) -> asset_configuration_client.AssetConfigurationClient:
    """Creates a client for the asset configuration service.

    If no recommendations are given, a default service will be created which
    will always return an empty response, i.e., it will never provide any
    recommended configuration. Otherwise, if recommendations are given, it is
    expected that they cover all requests made by the test and the service will
    raise an error for any request which does not match any preconfigured
    recommendation.

    Args:
      recommendations: Mapping from asset ID to a list of preconfigured
        recommendations in the form of (input_config, output_config) tuples.
    """

    def mock_recommend_asset_configuration(
        name: str, input_configuration: any_pb2.Any | None
    ) -> asset_configuration_pb2.RecommendAssetConfigurationResponse:
      if recommendations is None:
        return asset_configuration_pb2.RecommendAssetConfigurationResponse()

      if name not in recommendations:
        raise LookupError(f'No recommendations defined for {name}')

      for input_config_canditate, output_config_candidate in recommendations[
          name
      ]:
        if (input_config_canditate is None and input_configuration is None) or (
            input_config_canditate.SerializeToString()
            == input_configuration.value
        ):
          response = (
              asset_configuration_pb2.RecommendAssetConfigurationResponse()
          )
          response.config.Pack(output_config_candidate)
          return response

      raise LookupError(
          f'No defined recommendation for {name} matches the input config'
      )

    client = mock.MagicMock()
    client.recommend_asset_configuration.side_effect = (
        mock_recommend_asset_configuration
    )
    return client
