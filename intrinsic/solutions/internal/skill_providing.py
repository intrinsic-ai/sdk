# Copyright 2023 Intrinsic Innovation LLC

"""Contains the implementation of SkillProvider and related interfaces."""

from __future__ import annotations

from typing import Any, Iterable, Type, Union

from intrinsic.assets.proto import asset_tag_pb2
from intrinsic.assets.proto import asset_type_pb2
from intrinsic.assets.proto import view_pb2
from intrinsic.resources.client import resource_registry_client
from intrinsic.skills.client import skill_registry_client
from intrinsic.skills.proto import skills_pb2
from intrinsic.solutions import provided
from intrinsic.solutions import providers
from intrinsic.solutions.internal import installed_assets_client
from intrinsic.solutions.internal import resources
from intrinsic.solutions.internal import skill_generation
from intrinsic.solutions.internal import skill_utils


_SKILL_ID_SEP = "."


def _create_skill_packages(
    current_package: str,
    relative_child_package_names: set[str],
    skill_infos: list[skill_generation.SkillInfoImpl],
    compatible_resources_by_id: dict[str, dict[str, provided.ResourceList]],
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
  _installed_assets: installed_assets_client.InstalledAssetsClient | None

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
      installed_assets: (
          installed_assets_client.InstalledAssetsClient | None
      ) = None,
  ):
    self._skill_registry = skill_registry
    self._resource_registry = resource_registry
    self._installed_assets = installed_assets

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
    skill_infos = [
        skill_generation.SkillInfoImpl(skill_proto)
        for skill_proto in self._get_skill_protos()
    ]

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
      for slot_selector in skill_info.skill_proto.resource_selectors.values():
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
      for resource_slot in skill_info.skill_proto.resource_selectors:
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
        "", package_names, skill_infos, self._compatible_resources_by_id
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
            )
        )

  def _get_skill_protos(self) -> list[skills_pb2.Skill]:
    # Get all Process assets with the SUBPROCESS tag, i.e., the Process assets
    # that are marked to be listed alongside regular "skills" in UIs.
    process_asset_skills: list[skills_pb2.Skill] = []
    if self._installed_assets is not None:
      for asset in self._installed_assets.list_all_installed_assets(
          asset_types=[asset_type_pb2.AssetType.ASSET_TYPE_PROCESS],
          asset_tag=asset_tag_pb2.AssetTag.ASSET_TAG_SUBPROCESS,
          view=view_pb2.AssetViewType.ASSET_VIEW_TYPE_FULL,
      ):
        bt = asset.deployment_data.process.process.behavior_tree
        # Skip behavior trees without a Skill proto (this is unexpected)
        if bt.HasField("description"):
          process_asset_skills.append(bt.description)

    # Merge skill protos from process assets with skill protos from the skill
    # registry (legacy processes and regular skills), assets taking precedence.
    result: list[skills_pb2.Skill] = process_asset_skills
    collected_ids = set(skill.id for skill in process_asset_skills)
    for skill in self._skill_registry.get_skills():
      if skill.id not in collected_ids:
        result.append(skill)
        collected_ids.add(skill.id)

    result.sort(key=lambda skill: skill.id)
    return result
