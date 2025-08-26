# Copyright 2023 Intrinsic Innovation LLC

"""Entry point of the Intrinsic solution building libraries.

## Usage example

```
from intrinsic.solutions import deployments

solution = deployments.connect(...)
skills = solution.skills
executive = solution.executive

throw_ball = skills.throw_ball(...)
executive.run(throw_ball)
```
"""

import enum
import inspect
import sys
from typing import Optional

import grpc
from intrinsic.assets.proto import installed_assets_pb2_grpc
from intrinsic.frontend.solution_service.proto import solution_service_pb2
from intrinsic.frontend.solution_service.proto import solution_service_pb2_grpc
from intrinsic.frontend.solution_service.proto import status_pb2 as solution_status_pb2
from intrinsic.kubernetes.workcell_spec.proto import installer_pb2_grpc
from intrinsic.resources.client import resource_registry_client
from intrinsic.scene.product.client import product_client as product_client_mod
from intrinsic.skills.client import skill_registry_client
from intrinsic.solutions import auth
from intrinsic.solutions import dialerutil
from intrinsic.solutions import error_processing
from intrinsic.solutions import errors as solution_errors
from intrinsic.solutions import execution
from intrinsic.solutions import ipython
from intrinsic.solutions import pbt_registration
from intrinsic.solutions import pose_estimation
from intrinsic.solutions import proto_building
from intrinsic.solutions import providers
from intrinsic.solutions import simulation
from intrinsic.solutions import userconfig
from intrinsic.solutions import worlds
from intrinsic.solutions.internal import behavior_tree_providing
from intrinsic.solutions.internal import products as products_mod
from intrinsic.solutions.internal import resources as resources_mod
from intrinsic.solutions.internal import skill_providing
from intrinsic.solutions.internal import stubs

_GRPC_OPTIONS = [
    # Remove limit on message size for e.g. images.
    ("grpc.max_receive_message_length", -1),
    ("grpc.max_send_message_length", -1),
]

_CSS_FAILURE_STYLE = (
    "color: #ab0000; font-family: monospace; font-weight: bold; "
    "padding-left: var(--jp-code-padding);"
)

_WORLD_ID = "world"


class Solution:
  """Wraps a connection to a deployed solution.

  Attributes:
    grpc_channel: gRPC channel to the cluster which hosts the deployed solution.
    is_simulated: Whether the solution is deployed on a simulated workcell
      rather than on a physical workcell.
    executive: Executive instance to communicate with executive.
    behavior_trees: Behavior trees stored on the solution.
    skills: Wrapper to easily access skills.
    resources: Provides access to resources.
    products: Provides access to products.
    simulator: Simulator instance for controlling simulation.
    errors: Exposes error reports from executions.
    pose_estimators: Optional. Wrapper to access pose estimators.
    world: default world in world service.
    pbt_registry: gRPC wrapper to sideload PBTs
    proto_builder: service to build proto FileDescriptorSets from proto schemas
  """

  class HealthStatus(enum.Enum):
    """Health status of the solution's backend."""

    UNKNOWN = 0
    # Ready to receive requests.
    HEALTHY = 1
    # Not ready to receive requests, but should fix itself.
    PENDING = 2
    # Non-recoverable error.
    ERROR = 3

  grpc_channel: grpc.Channel
  is_simulated: bool
  executive: execution.Executive
  resources: providers.ResourceProvider
  products: providers.ProductProvider
  world: worlds.ObjectWorld
  simulator: Optional[simulation.Simulation]
  behavior_trees: providers.BehaviorTreeProvider
  skills: providers.SkillProvider
  errors: error_processing.ErrorsLoader
  pose_estimators: Optional[pose_estimation.PoseEstimators]
  _solution_service: solution_service_pb2_grpc.SolutionServiceStub
  _skill_registry: skill_registry_client.SkillRegistryClient
  _resource_registry: resource_registry_client.ResourceRegistryClient
  _installer_service_stub: installer_pb2_grpc.InstallerServiceStub
  pbt_registry: Optional[pbt_registration.BehaviorTreeRegistry]
  proto_builder: Optional[proto_building.ProtoBuilder]

  def __init__(
      self,
      grpc_channel: grpc.Channel,
      is_simulated: bool,
      executive: execution.Executive,
      solution_service: solution_service_pb2_grpc.SolutionServiceStub,
      skill_registry: skill_registry_client.SkillRegistryClient,
      resource_registry: resource_registry_client.ResourceRegistryClient,
      product_client: product_client_mod.ProductClient,
      object_world: worlds.ObjectWorld,
      simulator: Optional[simulation.Simulation],
      errors: error_processing.ErrorsLoader,
      pose_estimators: Optional[pose_estimation.PoseEstimators],
      pbt_registry: Optional[pbt_registration.BehaviorTreeRegistry] = None,
      proto_builder: Optional[proto_building.ProtoBuilder] = None,
  ):
    self.grpc_channel: grpc.Channel = grpc_channel
    self.is_simulated: bool = is_simulated

    self.executive = executive
    self._solution_service = solution_service
    self._skill_registry = skill_registry
    self._resource_registry = resource_registry
    self._product_client = product_client
    self.resources = resources_mod.Resources(self._resource_registry)
    self.products = products_mod.Products(self._product_client)

    self.world: worlds.ObjectWorld = object_world
    self.simulator: Optional[simulation.Simulation] = simulator

    self.behavior_trees = behavior_tree_providing.BehaviorTrees(
        self._solution_service
    )
    self.skills = skill_providing.Skills(
        self._skill_registry,
        self._resource_registry,
    )

    self.pose_estimators = pose_estimators
    self.errors = errors
    self.pbt_registry = pbt_registry
    self.proto_builder = proto_builder

  @classmethod
  def for_channel(
      cls,
      grpc_channel: grpc.Channel,
  ) -> "Solution":
    """Creates a Solution for the given channel and options.

    Args:
      grpc_channel: gRPC channel to the cluster which hosts the deployed
        solution.

    Returns:
      A fully initialized Workcell instance.
    """

    print("Connecting to deployed solution...")

    solution_service = solution_service_pb2_grpc.SolutionServiceStub(
        grpc_channel
    )

    try:
      solution_status = _get_solution_status_with_retry(solution_service)
    except solution_errors.BackendNoWorkcellError as e:
      ipython.display_html_or_print_msg(
          f'<span style="{_CSS_FAILURE_STYLE}">{str(e)}</span>', str(e)
      )
      raise

    # Optional backends.
    simulator = None
    if solution_status.simulated:
      simulator = simulation.Simulation.connect(grpc_channel)

    # Required backends.
    error_loader = error_processing.ErrorsLoader()
    executive = execution.Executive.connect(
        grpc_channel, error_loader, simulator
    )
    skill_registry = skill_registry_client.SkillRegistryClient.connect(
        grpc_channel
    )
    resource_registry = resource_registry_client.ResourceRegistryClient.connect(
        grpc_channel
    )

    # Remaining backends.
    product_client = product_client_mod.ProductClient.connect(grpc_channel)

    object_world = worlds.ObjectWorld.connect(_WORLD_ID, grpc_channel)
    installed_assets_stub = installed_assets_pb2_grpc.InstalledAssetsStub(
        grpc_channel
    )
    pose_estimators = pose_estimation.PoseEstimators(
        installed_assets_stub,
    )

    pbt_registry = pbt_registration.BehaviorTreeRegistry.connect(grpc_channel)
    proto_builder = proto_building.ProtoBuilder.connect(grpc_channel)

    print(
        "Connected successfully to"
        f' "{solution_status.display_name}({solution_status.platform_version})"'
        f' at "{solution_status.cluster_name}".'
    )
    return cls(
        grpc_channel,
        solution_status.simulated,
        executive,
        solution_service,
        skill_registry,
        resource_registry,
        product_client,
        object_world,
        simulator,
        error_loader,
        pose_estimators,
        pbt_registry,
        proto_builder,
    )

  def get_health_status(self) -> "HealthStatus":
    """Returns the health status of the solution backend.

    Can be called before or after connecting to the deployed solution.

    Returns:
      Health status of solution backend
    """
    status = self._solution_service.GetStatus(
        solution_service_pb2.GetStatusRequest()
    ).state
    if status == solution_status_pb2.Status.State.READY:
      return self.HealthStatus.HEALTHY
    if status == solution_status_pb2.Status.State.DEPLOYING:
      return self.HealthStatus.PENDING
    if status == solution_status_pb2.Status.State.ERROR:
      return self.HealthStatus.ERROR
    return self.HealthStatus.UNKNOWN

  def skills_overview(
      self,
      with_signatures: bool = False,
      with_type_annotations: bool = False,
      with_doc: bool = False,
  ) -> None:
    """Prints an overview of the registered skills.

    Args:
      with_signatures: Include signatures for skill construction.
      with_type_annotations: Include type annotations and not just the parameter
        name.
      with_doc: Also print out docstring for each skill.
    """

    def build_signature(skill, with_type_annotations: bool) -> str:
      """Build a signature for skill, optionally including type annotations.

      Args:
        skill: The skill to build the signature for.
        with_type_annotations: Include type annotations and not just the
          parameter name.

      Returns:
        The skill signature.
      """
      signature = inspect.signature(skill)
      parameters = []
      for k, v in signature.parameters.items():
        if not with_type_annotations:
          parameters.append(k)
          continue
        param = str(v)
        parameters.append(param)
      return ", ".join(parameters)

    for skill_id, skill in self.skills.get_skill_ids_and_classes():
      if with_signatures:
        signature_str = build_signature(skill, with_type_annotations)
        print(f"{skill_id}({signature_str})")
      else:
        print(skill_id)
      if with_doc:
        print(f"\n{inspect.getdoc(skill)}\n")

  def update_skills(self) -> None:
    self.skills.update()

  def generate_stubs(self, output_path: str) -> None:
    """Generates type stubs for all available skill classes in the solution.

    The generated stubs can be provided to an IDE or type checker to get proper
    language and type support when working with the auto-generated skill classes
    of the solution building library. Usage examples:
      - VS Code: The given 'output_path' should match the value of the
        'python.analysis.stubPath' setting. After generating the stubs, a
        restart of the Python language server is usually required.
      - mypy: The given 'output_path' should be included in the $MYPYPATH
        environment variable.

    The generated stubs are specific to a solution. They match the skills
    installed in the solution at their respective version. Hence the stubs need
    to be updated when the skills in the solution change or when connecting to a
    different solution.

    Args:
      output_path: The root directory into which the stubs shall be written.
        E.g., the file '<output_path>/intrinsic-stubs/solutions/providers.pyi'
        will be generated.
    """
    stubs.generate(output_path, self.skills, sys.stdout)


def connect(
    *,
    grpc_channel: Optional[grpc.Channel] = None,
    address: Optional[str] = None,
    org: Optional[str] = None,
    solution: Optional[str] = None,
    cluster: Optional[str] = None,
) -> "Solution":
  """Connects to a deployed solution.

  If arguments are not provided, connect to the solution specified in the user
  config.

  Args:
    grpc_channel: gRPC channel to use for connection.
    address: Connect directly to an address (e.g. localhost). Only one of
      solution and address is allowed.
    org: Organization of the solution to connect to.
    solution: Id (not display name!) of the solution to connect to.
    cluster: Name of cluster to connect to (instead of specifying 'solution').

  Raises:
    ValueError: if parameter combination is incorrect.

  Returns:
    A fully initialized Solution object that represents the deployed solution.
  """
  if (
      sum([
          (grpc_channel is not None),
          (solution is not None),
          (cluster is not None),
          (address is not None),
      ])
      > 1
  ):
    raise ValueError(
        "Only one of grpc_channel, solution, cluster or address is allowed!"
    )

  if (org is not None) and (
      (grpc_channel is not None) or (address is not None)
  ):
    raise ValueError(
        "Org is not supported when connecting via grpc_channel or address."
    )

  if not any([
      grpc_channel,
      solution,
      cluster,
      address,
  ]):
    return connect_to_selected_solution()

  if grpc_channel is None:
    grpc_channel = _create_grpc_channel(
        address=address,
        org=org,
        solution=solution,
        cluster=cluster,
    )

  return Solution.for_channel(grpc_channel)


_NO_SOLUTION_SELECTED_ERROR = (
    "No solution selection can be found in the current environment! E.g., in VS"
    " Code you can use the Intrinsic extension to select a deployed solution."
)

_INVALID_SOLUTION_SELECTION_ERROR = (
    "The solution selection found in the current environment is invalid. To"
    " correctly select a solution you can use, e.g., the Intrinsic extension in"
    " VS Code."
)


def connect_to_selected_solution() -> "Solution":
  """Connects to the solution specified in the user config.

  Connects to the deployed solution that is selected in the current environment.
  E.g., in VS Code you can use the Intrinsic extension to select a deployed
  solution from a list of available solutions and then use this method to
  connect to this solution.
  Raises:
    NotFoundError: If no valid solution is specified in the user config.

  Returns:
    A fully initialized Solution object that represents the deployed solution.
  """
  try:
    config = userconfig.read()
  except userconfig.NotFoundError as e:
    raise solution_errors.NotFoundError(_NO_SOLUTION_SELECTED_ERROR) from e

  selected_org = config.get(userconfig.SELECTED_ORGANIZATION, None)
  selected_solution = config.get(userconfig.SELECTED_SOLUTION, None)
  selected_cluster = config.get(userconfig.SELECTED_CLUSTER, None)
  selected_address = config.get(userconfig.SELECTED_ADDRESS, None)

  try:
    grpc_channel = _create_grpc_channel(
        address=selected_address,
        org=selected_org,
        solution=selected_solution,
        cluster=selected_cluster,
    )
  except ValueError as e:
    raise solution_errors.NotFoundError(
        _INVALID_SOLUTION_SELECTION_ERROR
    ) from e

  return Solution.for_channel(grpc_channel)


# Disable pytype error since the raise is incorrectly detected as returning None
# pytype: disable=bad-return-type
def _create_grpc_channel(
    *,
    address: Optional[str] = None,
    org: Optional[str] = None,
    solution: Optional[str] = None,
    cluster: Optional[str] = None,
) -> grpc.Channel:
  """Creates a gRPC channel to a deployed solution.

  Args:
    address: Connect directly to an address (e.g. localhost). Only one of
      solution and address is allowed.
    org: Organization of the solution to connect to.
    solution: Id (not display name!) of the solution to connect to.
    cluster: Name of cluster to connect to (instead of specifying 'solution').

  Returns:
    gRPC channel to the deployed solution.
  """
  org_info = None
  if org is not None:
    try:
      org_info = auth.read_org_info(org)
    except auth.OrgNotFoundError as error:
      raise solution_errors.NotFoundError(
          f"Credentials for organization '{error.organization}' not found."
          f" Run 'inctl auth login --org {error.organization}' on a terminal"
          " to login with this organization, or run 'inctl auth list' to see"
          " the organizations you are currently logged in with."
      ) from error

  if solution is not None:
    if org_info is None:
      raise ValueError(
          f"'org' is required when 'solution' ({solution}) is set!"
      )
    return dialerutil.create_channel_from_solution(
        org_info, solution, grpc_options=_GRPC_OPTIONS
    )
  elif cluster is not None:
    if org_info is None:
      raise ValueError(f"'org' is required when 'cluster' ({cluster}) is set!")
    return dialerutil.create_channel_from_cluster(
        org_info, cluster, grpc_options=_GRPC_OPTIONS
    )
  elif address is not None:
    return dialerutil.create_channel_from_address(
        address, grpc_options=_GRPC_OPTIONS
    )
  else:
    raise ValueError("No connection parameters specified!")


# pytype: enable=bad-return-type


@solution_errors.retry_on_pending_backend
def _get_solution_status_with_retry(
    stub: solution_service_pb2_grpc.SolutionServiceStub,
) -> solution_status_pb2.Status:
  """Loads a solution's status.

  Args:
    stub: Solution service to query health.

  Returns:
    Solution status

  Raises:
    solution_errors.BackendPendingError: Will lead to retry.
    solution_errors.BackendHealthError: Not recoverable.
  """
  try:
    response = stub.GetStatus(solution_service_pb2.GetStatusRequest())

    if response.state != solution_status_pb2.Status.State.READY:
      if response.state in [
          solution_status_pb2.Status.State.PLATFORM_DEPLOYING,
          solution_status_pb2.Status.State.PLATFORM_READY,
          solution_status_pb2.Status.State.DEPLOYING,
      ]:
        print("Solution not ready yet. Retrying...")
        print(f"Reason: {response.state_reason}")
        # Note this error leads to a retry given the retry_on_pending_backend
        # decorator.
        raise solution_errors.BackendPendingError(
            f"Solution not ready yet. {response.state_reason}"
        )
      if response.state == solution_status_pb2.Status.State.ERROR:
        raise solution_errors.BackendHealthError(
            "Solution backend is unhealthy and not expected to recover "
            "without intervention. Try restarting your solution. "
            f"{response.state_reason}"
        )
      else:
        raise solution_errors.BackendHealthError(
            "Unexpected solution status. Try restarting your "
            f"solution. {response.state_reason}"
        )
    return response
  except grpc.RpcError as e:
    if hasattr(e, "code"):
      if e.code() in [
          grpc.StatusCode.UNIMPLEMENTED,
          grpc.StatusCode.UNAVAILABLE,
      ]:
        raise solution_errors.BackendPendingError(
            "Transfer service is not available yet."
        )
      if e.code() == grpc.StatusCode.NOT_FOUND:
        raise solution_errors.BackendNoWorkcellError(
            "No solution has been started. Start a solution before connecting."
        )
    raise
