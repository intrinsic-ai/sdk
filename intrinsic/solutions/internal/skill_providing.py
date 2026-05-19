# Copyright 2023 Intrinsic Innovation LLC

"""Contains the implementation of SkillProvider and related interfaces."""

from __future__ import annotations

from typing import Any
from typing import Iterable
from typing import Type
from typing import Union

from intrinsic.assets.configuration import asset_configuration_client
from intrinsic.assets.install import installed_assets_client
from intrinsic.assets.proto import asset_tag_pb2
from intrinsic.assets.proto import asset_type_pb2
from intrinsic.assets.proto import id_pb2
from intrinsic.assets.proto import installed_assets_pb2
from intrinsic.assets.proto import view_pb2
from intrinsic.resources.client import resource_registry_client
from intrinsic.skills.client import skill_registry_client
from intrinsic.solutions import provided
from intrinsic.solutions import providers
from intrinsic.solutions.internal import resources
from intrinsic.solutions.internal import skill_generation
from intrinsic.solutions.internal import skill_utils
from intrinsic.util.proto import descriptors

_SKILL_ID_SEP = "."


def _create_skill_packages(
    current_package: str,
    relative_child_package_names: set[str],
    skill_infos: list[skill_generation.SkillInfoImpl],
    compatible_resources_by_id: dict[str, dict[str, provided.ResourceList]],
    asset_config_client: asset_configuration_client.AssetConfigurationClient,
) -> dict[str, _SkillPackageImpl]:
  """Recursively creates a hierarchy of skill packages.

  Args:
    current_package: The current package name, e.g., "" or "ai".
    relative_child_package_names: The package names to create below
      'current_package' and  relative to 'current_package'. Passed as a set to
      avoid duplicates. E.g., {'intrinsic', 'intrinsic.special_skills'} if the
      current package is 'ai'.
    skill_infos: All skill infos from the skill registry.
    compatible_resources_by_id: All compatible resource definitions by skill id.

  Returns:
    A mapping from skill package prefix to child skill packages. E.g., if the
    current package is 'ai' the returned mapping might be:
      { 'intrinsic': <SkillPackageImpl for 'ai.intrinsic'> }
  """

  prefixes = {package.split(".")[0] for package in relative_child_package_names}

  skill_packages = {}
  for prefix in prefixes:
    prefix_dotted = prefix + "."
    relative_child_package_names_of_child = {
        package.removeprefix(prefix_dotted)
        for package in relative_child_package_names
        if package.startswith(prefix_dotted)
    }
    skill_packages[prefix] = _SkillPackageImpl(
        current_package + "." + prefix if current_package else prefix,
        relative_child_package_names_of_child,
        skill_infos,
        compatible_resources_by_id,
        asset_config_client,
    )

  return skill_packages


class _SkillPackageImpl(provided.SkillPackage):
  """Implementation of the SkillPackage interface."""

  # Full package name (e.g. 'ai.intrinsic').
  _package_name: str

  # All skills in the package by name, not including skills from child packages
  _skills_by_name: dict[str, provided.SkillInfo]
  _compatible_resources_by_name: dict[str, dict[str, provided.ResourceList]]
  _skill_type_classes_by_name: dict[str, Type[provided.SkillBase]]

  _child_packages_by_prefix: dict[str, _SkillPackageImpl]

  def __init__(
      self,
      package_name: str,
      relative_child_package_names: set[str],
      skill_infos: list[skill_generation.SkillInfoImpl],
      compatible_resources_by_id: dict[str, dict[str, provided.ResourceList]],
      asset_config_client: asset_configuration_client.AssetConfigurationClient,
  ):
    """Creates a new skills package including all child packages.

    Args:
      package_name: The full package name (e.g. 'ai.intrinsic').
      relative_child_package_names: The relative names of any child packages.
        E.g. {'intrinsic', 'intrinsic.special_skills'} if the package name is
        'ai'.
      skill_infos: All skill infos from the skill registry.
      compatible_resources_by_id: All compatible resource definitions by skill
        id.
    """

    self._package_name = package_name
    self._asset_config_client = asset_config_client
    self._skills_by_name = {
        info.skill_name: info
        for info in skill_infos
        if info.package_name == package_name
    }
    self._compatible_resources_by_name = {
        name: compatible_resources_by_id[info.id]
        for name, info in self._skills_by_name.items()
    }
    self._skill_type_classes_by_name = {}
    self._child_packages_by_prefix = _create_skill_packages(
        package_name,
        relative_child_package_names,
        skill_infos,
        compatible_resources_by_id,
        asset_config_client,
    )

  @property
  def package_name(self) -> str:
    return self._package_name

  @property
  def relative_package_name(self) -> str:
    return self._package_name.split(_SKILL_ID_SEP)[-1]

  def __getattr__(self, name: str) -> Union[Type[Any], provided.SkillPackage]:
    if name in self._child_packages_by_prefix:
      return self._child_packages_by_prefix[name]
    elif name not in self._skills_by_name:
      raise AttributeError(
          f"Could not resolve the attribute '{name}'."
          f" '{self._package_name}.{name}' is not an available skill id or"
          " skill package. Use update() or reconnect to the deployed solution"
          " to update the available skills."
      )
    if name not in self._skill_type_classes_by_name:
      self._skill_type_classes_by_name[name] = skill_generation.gen_skill_class(
          self._skills_by_name[name],
          self._compatible_resources_by_name[name],
          self._asset_config_client,
      )
    return self._skill_type_classes_by_name[name]

  def __dir__(self) -> list[str]:
    return sorted(
        self._child_packages_by_prefix.keys() | self._skills_by_name.keys()
    )


class Skills(providers.SkillProvider):
  """Implementation of the SkillProvider interface."""

  _skill_registry: skill_registry_client.SkillRegistryClient
  _resource_registry: resource_registry_client.ResourceRegistryClient
  _installed_assets: installed_assets_client.InstalledAssetsClient

  # All skills in the solution by id
  _skills_by_id: dict[str, provided.SkillInfo]
  _skill_type_classes_by_id: dict[str, Type[provided.SkillBase]]
  _compatible_resources_by_id: dict[str, dict[str, provided.ResourceList]]

  # Only "global" skills with an empty package name
  _global_skills_by_name: dict[str, provided.SkillInfo]
  _global_skill_type_classes_by_name: dict[str, Type[provided.SkillBase]]
  _global_compatible_resources_by_name: dict[
      str, dict[str, provided.ResourceList]
  ]

  _skill_packages: dict[str, _SkillPackageImpl]

  def __init__(
      self,
      skill_registry: skill_registry_client.SkillRegistryClient,
      resource_registry: resource_registry_client.ResourceRegistryClient,
      installed_assets: installed_assets_client.InstalledAssetsClient,
      asset_config_client: asset_configuration_client.AssetConfigurationClient,
  ):
    self._skill_registry = skill_registry
    self._resource_registry = resource_registry
    self._installed_assets = installed_assets
    self._asset_config_client = asset_config_client

    self._skills_by_id = {}
    self._skill_type_classes_by_id = {}
    self._compatible_resources_by_id = {}

    self._global_skills_by_name = {}
    self._global_skill_type_classes_by_name = {}
    self._global_compatible_resources_by_name = {}

    self._skill_packages = {}
    self.update()

  def update(self) -> None:
    """Retrieve most recent list of available skills."""
    skill_infos = self._collect_skill_infos()

    self._global_skills_by_name = {
        info.skill_name: info for info in skill_infos if not info.package_name
    }
    self._global_skill_type_classes_by_name = {}
    self._global_compatible_resources_by_name = {}

    self._skills_by_id = {info.id: info for info in skill_infos}
    self._skill_type_classes_by_id = {}
    self._compatible_resources_by_id = {}

    # Get handles for all selectors that we need in the loop below in a single
    # batch request for performance reasons.
    capability_names_batch = []
    for skill_info in skill_infos:
      for slot_selector in skill_info.resource_selectors.values():
        capability_names_batch.append(slot_selector.capability_names)

    handles_by_selector = self._resource_registry.batch_list_all_resource_handles(  # pytype: disable=wrong-arg-types  # always-use-property-annotation
        capability_names_batch=capability_names_batch
    )

    # Update compatible resources for skills
    selector_index = 0
    for skill_info in skill_infos:
      skill_name = skill_info.skill_name
      is_global = not skill_info.package_name

      if is_global:
        self._global_compatible_resources_by_name[skill_name] = {}
      self._compatible_resources_by_id[skill_info.id] = {}
      for resource_slot in skill_info.resource_selectors:
        handles = handles_by_selector[selector_index]
        slot = resource_slot
        if slot in skill_info.field_names:
          slot = slot + skill_utils.RESOURCE_SLOT_DECONFLICT_SUFFIX
        resource_list = resources.ResourceListImpl(
            [provided.ResourceHandle(h) for h in handles]
        )

        if is_global:
          self._global_compatible_resources_by_name[skill_name][
              slot
          ] = resource_list
        self._compatible_resources_by_id[skill_info.id][slot] = resource_list

        selector_index += 1

    package_names = {
        package_name
        for info in skill_infos
        if (package_name := info.package_name)
    }
    self._skill_packages = _create_skill_packages(
        "",
        package_names,
        skill_infos,
        self._compatible_resources_by_id,
        self._asset_config_client,
    )

  def __getattr__(self, name: str) -> Union[Type[Any], provided.SkillPackage]:
    if name in self._skill_packages:
      return self._skill_packages[name]
    elif name not in self._global_skills_by_name:
      raise AttributeError(
          f"Could not resolve the attribute '{name}'. '{name}' is not an"
          " available skill id, skill package or name of a package-less skill."
          " Use update() or reconnect to the deployed solution to update the"
          " available skills."
      )
    # Provide global skills with an empty package name (e.g.
    # `skills.global_skill`).
    if name not in self._global_skill_type_classes_by_name:
      self._global_skill_type_classes_by_name[name] = (
          skill_generation.gen_skill_class(
              self._global_skills_by_name[name],
              self._global_compatible_resources_by_name[name],
              self._asset_config_client,
          )
      )
    return self._global_skill_type_classes_by_name[name]

  def __len__(self) -> int:
    return len(self._skills_by_id)

  def __dir__(self) -> list[str]:
    return sorted(
        self._skill_packages.keys() | self._global_skills_by_name.keys()
    )

  def __getitem__(self, skill_id: str) -> Type[Any]:
    if skill_id in self._skills_by_id:
      if skill_id not in self._skill_type_classes_by_id:
        self._skill_type_classes_by_id[skill_id] = (
            skill_generation.gen_skill_class(
                self._skills_by_id[skill_id],
                self._compatible_resources_by_id[skill_id],
                self._asset_config_client,
            )
        )
      return self._skill_type_classes_by_id[skill_id]
    else:
      raise KeyError(
          f"Could not resolve the attribute '{skill_id}'. '{skill_id}' is not"
          " an available skill skill id. Use update() or reconnect to the"
          " deployed solution to update the available skills."
      )

  def __contains__(self, skill_id: str) -> bool:
    return skill_id in self._skills_by_id

  def __iter__(self):
    for skill_id in self._skills_by_id.keys():
      yield self[skill_id]

  def get_skill_ids(self) -> Iterable[str]:
    return self._skills_by_id.keys()

  def get_skill_classes(self) -> Iterable[Type[Any]]:
    self._generate_all_skills()
    return self._skill_type_classes_by_id.values()

  def get_skill_ids_and_classes(self) -> Iterable[tuple[str, Type[Any]]]:
    self._generate_all_skills()
    return self._skill_type_classes_by_id.items()

  def _generate_all_skills(self) -> None:
    for skill_id in self._skills_by_id:
      if skill_id not in self._skill_type_classes_by_id:
        self._skill_type_classes_by_id[skill_id] = (
            skill_generation.gen_skill_class(
                self._skills_by_id[skill_id],
                self._compatible_resources_by_id[skill_id],
                self._asset_config_client,
            )
        )

  def _process_asset_to_skill_info(
      self,
      asset: installed_assets_pb2.InstalledAsset,
  ) -> provided.SkillInfo:
    skill = asset.deployment_data.process.process.behavior_tree.description

    file_descriptor_set = descriptors.merge_file_descriptor_sets([
        skill.parameter_description.parameter_descriptor_fileset,
        skill.return_value_description.descriptor_fileset,
    ])
    default_value = (
        skill.parameter_description.default_value
        if skill.parameter_description.HasField("default_value")
        else None
    )
    # Comments are currently only available from the skill proto and not from
    # the file descriptor set in the asset metadata (b/510731384).
    proto_comments = dict(
        skill.parameter_description.parameter_field_comments
    ) | dict(skill.return_value_description.return_value_field_comments)
    return skill_generation.SkillInfoImpl(
        skill,
        asset.metadata.id_version,
        asset.metadata.documentation.description,
        skill.parameter_description.parameter_message_full_name,
        skill.return_value_description.return_value_message_full_name,
        file_descriptor_set,
        default_value,
        dict(skill.resource_selectors),
        proto_comments,
        provided.SkillType.PROCESS,
        skill_utils.INTRINSIC_TYPE_URL_AREA_ASSETS,
    )

  def _skill_proto_from_skill_registry_to_skill_info(
      self,
      skill: skills_pb2.Skill,
  ) -> provided.SkillInfo:
    # Manually decompose the 'id_version' here. We cannot use 'idutils'
    # since PBTs from the skill registry can still have invalid package
    # names (empty or with only a single period). See b/512066439.
    package, _, name = skill.id.rpartition(".")
    version = skill.id_version.removeprefix(skill.id + ".")
    id_version = id_pb2.IdVersion(
        id=id_pb2.Id(package=package, name=name), version=version
    )

    file_descriptor_set = descriptors.merge_file_descriptor_sets([
        skill.parameter_description.parameter_descriptor_fileset,
        skill.return_value_description.descriptor_fileset,
    ])
    default_value = (
        skill.parameter_description.default_value
        if skill.parameter_description.HasField("default_value")
        else None
    )
    proto_comments = dict(
        skill.parameter_description.parameter_field_comments
    ) | dict(skill.return_value_description.return_value_field_comments)

    # A skill proto from the skill registry can represent either an actual skill
    # or a legacy PBT (but not a Process asset). This can be told based on the
    # presence of 'behavior_tree_description'. For actual skills we should use
    # the 'assets' area since those can also be looked up via the installed
    # assets service. Legacy PBT's are only available in the skill registry and
    # thus need the 'skills' area.
    if skill.HasField("behavior_tree_description"):
      skill_type = provided.SkillType.PROCESS
      area = skill_utils.INTRINSIC_TYPE_URL_AREA_SKILLS
    else:
      skill_type = provided.SkillType.REGULAR_SKILL
      area = skill_utils.INTRINSIC_TYPE_URL_AREA_ASSETS

    return skill_generation.SkillInfoImpl(
        skill,
        id_version,
        skill.description,
        skill.parameter_description.parameter_message_full_name,
        skill.return_value_description.return_value_message_full_name,
        file_descriptor_set,
        default_value,
        dict(skill.resource_selectors),
        proto_comments,
        skill_type,
        area,
    )

  def _collect_skill_infos(self) -> list[provided.SkillInfo]:
    # Get all Process assets with the SUBPROCESS tag, i.e., the Process assets
    # that are marked to be listed alongside regular "skills" in UIs.
    process_asset_skills: list[provided.SkillInfo] = []
    for asset in self._installed_assets.list_all_installed_assets(
        asset_types=[asset_type_pb2.AssetType.ASSET_TYPE_PROCESS],
        asset_tag=asset_tag_pb2.AssetTag.ASSET_TAG_SUBPROCESS,
        view=view_pb2.AssetViewType.ASSET_VIEW_TYPE_FULL,
    ):
      bt = asset.deployment_data.process.process.behavior_tree
      # Skip behavior trees without a Skill proto (this is unexpected)
      if not bt.HasField("description"):
        continue

      process_asset_skills.append(self._process_asset_to_skill_info(asset))

    # Merge skill protos from process assets with skill protos from the skill
    # registry (legacy processes and regular skills), assets taking precedence.
    result: list[provided.SkillInfo] = process_asset_skills
    collected_ids = set(skill.id for skill in process_asset_skills)
    for skill in self._skill_registry.get_skills():
      if skill.id not in collected_ids:
        result.append(
            self._skill_proto_from_skill_registry_to_skill_info(skill)
        )
        collected_ids.add(skill.id)

    result.sort(key=lambda skill: skill.id)

    return result
