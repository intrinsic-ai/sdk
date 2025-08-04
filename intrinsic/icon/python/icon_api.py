# Copyright 2023 Intrinsic Innovation LLC

"""Python client for the ICON Application Layer.

Provides a Python client API for application developers and skill authors who
wish to interact with ICON-compatible robots.
"""

from __future__ import annotations

from collections.abc import Iterable, Mapping
import enum
from typing import Optional, Union

import grpc
from intrinsic.icon.proto import logging_mode_pb2
from intrinsic.icon.proto.v1 import service_pb2
from intrinsic.icon.proto.v1 import service_pb2_grpc
from intrinsic.icon.proto.v1 import types_pb2
from intrinsic.icon.python import _session
from intrinsic.icon.python import actions
from intrinsic.icon.python import errors
from intrinsic.icon.python import reactions
from intrinsic.icon.python import state_variable_path
from intrinsic.logging.proto import context_pb2
from intrinsic.solutions import deployments
from intrinsic.util.grpc import connection
from intrinsic.util.grpc import interceptor
from intrinsic.world.robot_payload.python import robot_payload

# Type forwarding, to enable instantiating these without loading the respective
# modules in client code. We believe that wrapping all class definitions in a
# single ICON module will increase usability.
Action = actions.Action
Condition = reactions.Condition
Reaction = reactions.Reaction
StartActionInRealTime = reactions.StartActionInRealTime
StartParallelActionInRealTime = reactions.StartParallelActionInRealTime
TriggerCallback = reactions.TriggerCallback
Event = reactions.Event
EventFlag = reactions.EventFlag
OperationalState = types_pb2.OperationalState
StateVariablePath = state_variable_path.StateVariablePath
# For generating documentation, Session needs to be publicly visible, but we
# override the settings to at least hide the constructor since it's not meant
# to be directly created.
Session = _session.Session
Stream = _session.Stream
__pdoc__ = {}
__pdoc__["Session.__init__"] = None

_DEFAULT_INSECURE = True
_DEFAULT_RPC_TIMEOUT_INFINITE = None
_DEFAULT_CONNECT_TIMEOUT_SECONDS = 20


def _create_stub(
    connection_params: connection.ConnectionParams,
    insecure: bool = _DEFAULT_INSECURE,
    connect_timeout: int = _DEFAULT_CONNECT_TIMEOUT_SECONDS,
) -> service_pb2_grpc.IconApiStub:
  """Creates a stub for the ICON gRPC service.

  Args:
    connection_params: The required parameters to talk to the specific ICON
      instance.
    insecure: Whether to use insecure channel credentials.
    connect_timeout: Time in seconds to wait for the ICON gRPC server to be
      ready.

  Returns:
    The ICON Client API stub.
  """
  if insecure:
    channel = grpc.insecure_channel(connection_params.address)
  else:
    channel_creds = grpc.local_channel_credentials()
    channel = grpc.secure_channel(connection_params.address, channel_creds)

  try:
    grpc.channel_ready_future(channel).result(timeout=connect_timeout)
  except grpc.FutureTimeoutError as e:
    raise errors.Client.ServerError("Failed to connect to ICON server") from e

  channel = grpc.intercept_channel(
      channel, interceptor.HeaderAdderInterceptor(connection_params.headers)
  )

  return service_pb2_grpc.IconApiStub(channel)


class Client:
  """Wrapper for the ICON gRPC service.

  Attributes:
    ActionType: Dynamically generated enum of all available action type names.
  """

  class HardwareGroup(enum.Enum):
    ALL_HARDWARE = 1
    OPERATIONAL_HARDWARE_ONLY = 2

  # Explicitly avoid errors around dynamically-populated action enums.
  _HAS_DYNAMIC_ATTRIBUTES = True
  _rpc_timeout_seconds: Optional[int] = None

  def __init__(
      self,
      stub: service_pb2_grpc.IconApiStub,
      rpc_timeout: Optional[int] = _DEFAULT_RPC_TIMEOUT_INFINITE,
  ):
    # Ensure the timeout is set before calling any methods that do RPCs, like
    # self._generate_action_types. By setting it before self._stub we should be
    # safe.
    self._rpc_timeout_seconds = rpc_timeout
    self._stub = stub
    self._generate_action_types()

  # Disable lint warnings since this is a class, not a standard attribute.
  # pylint: disable=invalid-name
  @property
  def ActionType(self) -> enum.Enum:
    return self._ActionType

  @ActionType.setter
  def ActionType(self, value: enum.Enum):
    self._ActionType = value

  # pylint: enable=invalid-name

  def _generate_action_types(self) -> None:
    """Dynamically generates the ActionType enum from the available actions."""
    action_signatures = self.list_action_signatures()
    action_type_names = {}
    for action_signature in action_signatures:
      if not action_signature.action_type_name:
        continue
      # Strip out namespace prefixes and convert to upper case constant.
      const_name = action_signature.action_type_name.split(".")[-1].upper()
      action_type_names[const_name] = action_signature.action_type_name
    # Disable lint warnings since this is a class, not a standard attribute.
    # pylint: disable=invalid-name
    self.ActionType = enum.Enum("ActionType", action_type_names)
    # pylint: enable=invalid-name

  @classmethod
  def connect(
      cls,
      grpc_host: str = "localhost",
      grpc_port: int = 8128,
      insecure: bool = _DEFAULT_INSECURE,
      connect_timeout: int = _DEFAULT_CONNECT_TIMEOUT_SECONDS,
      rpc_timeout: Optional[int] = _DEFAULT_RPC_TIMEOUT_INFINITE,
  ) -> Client:
    """Connects to the ICON gRPC service.

    This is a convenience wrapper around creating the stub and instantiating the
    Client separately.

    Args:
      grpc_host: Host to connect to for the ICON gRPC server.
      grpc_port: Port to connect to for the ICON gRPC server.
      insecure: Whether to use insecure channel credentials.
      connect_timeout: Time in seconds to wait for the ICON gRPC server to be
        ready.
      rpc_timeout: Time in seconds to wait for RPCs to complete.

    Returns:
      An instance of the ICON Client.
    """
    return cls.connect_with_params(
        connection.ConnectionParams.no_ingress(f"{grpc_host}:{grpc_port}"),
        insecure=insecure,
        connect_timeout=connect_timeout,
        rpc_timeout=rpc_timeout,
    )

  @classmethod
  def connect_with_params(
      cls,
      connection_params: connection.ConnectionParams,
      insecure: bool = _DEFAULT_INSECURE,
      connect_timeout: int = _DEFAULT_CONNECT_TIMEOUT_SECONDS,
      rpc_timeout: Optional[int] = _DEFAULT_RPC_TIMEOUT_INFINITE,
  ) -> Client:
    """Connects to the ICON gRPC service.

    This is a convenience wrapper around creating the stub and instantiating the
    Client separately.

    Args:
      connection_params: The required parameters to talk to the specific ICON
        instance.
      insecure: Whether to use insecure channel credentials.
      connect_timeout: Time in seconds to wait for the ICON gRPC server to be
        ready.
      rpc_timeout: Time in seconds to wait for RPCs to complete.

    Returns:
      An instance of the ICON Client.
    """
    return cls(
        _create_stub(
            connection_params=connection_params,
            insecure=insecure,
            connect_timeout=connect_timeout,
        ),
        rpc_timeout,
    )

  @classmethod
  def for_solution(
      cls,
      solution: deployments.Solution,
      instance_name: Optional[str] = None,
      header: str = "x-resource-instance-name",
  ) -> Client:
    """Connects to the ICON gRPC service for a given solution."""
    channel = solution.grpc_channel
    if instance_name:
      channel = grpc.intercept_channel(
          channel,
          interceptor.HeaderAdderInterceptor(lambda: [(header, instance_name)]),
      )
    return cls(service_pb2_grpc.IconApiStub(channel))

  def get_action_signature_by_name(
      self, action_type_name: str
  ) -> Optional[types_pb2.ActionSignature]:
    """Gets details of an action type, by name.

    Args:
      action_type_name: The action type to lookup.

    Returns:
      ActionSignature, or None if the action type is not found.
      Propagates gRPC exceptions.
    """
    response = self._stub.GetActionSignatureByName(
        service_pb2.GetActionSignatureByNameRequest(name=action_type_name),
        timeout=self._rpc_timeout_seconds,
    )
    if not response.HasField("action_signature"):
      return None
    return response.action_signature

  def get_config(self) -> service_pb2.GetConfigResponse:
    """Gets part-specific config properties.

    These are fixed properties for the lifetime of the server (for
    example, the number of DOFs for a robot arm.)
    Returns:
      GetConfigResponse.
      Propagates gRPC exceptions.
    """
    return self._stub.GetConfig(
        service_pb2.GetConfigRequest(), timeout=self._rpc_timeout_seconds
    )

  def get_status(self) -> service_pb2.GetStatusResponse:
    """Gets a snapshot of the server-side status, including part-specific status.

    Returns:
      GetStatusResponse.
      Propagates gRPC exceptions.
    """
    return self._stub.GetStatus(
        service_pb2.GetStatusRequest(), timeout=self._rpc_timeout_seconds
    )

  def is_action_compatible(self, action_type_name: str, part: str) -> bool:
    """Reports whether actions of type `action_type_name` are compatible with `part`.

    Args:
      action_type_name: The action type to check.
      part: Name of the part to check.

    Returns:
      True iff actions of type `action_type_name` can be instantiated using
      `part` (in one of their slots).
      Propagates gRPC exceptions.
    """
    return self._stub.IsActionCompatible(
        service_pb2.IsActionCompatibleRequest(
            action_type_name=action_type_name, part_name=part
        ),
        timeout=self._rpc_timeout_seconds,
    ).is_compatible

  def list_action_signatures(self) -> Iterable[types_pb2.ActionSignature]:
    """Lists details of all available action types.

    Returns:
      Iterable of ActionSignatures.
    """
    return self._stub.ListActionSignatures(
        service_pb2.ListActionSignaturesRequest(),
        timeout=self._rpc_timeout_seconds,
    ).action_signatures

  def list_compatible_parts(
      self, action_type_names: Iterable[str]
  ) -> list[str]:
    """Lists the parts that are compatible with all of the listed action types.

    Args:
      action_type_names: The action types to check.

    Returns:
      List of individual parts that can be controlled by actions listed in
      `action_type_name`. If `action_type_names` is empty, returns all parts.
    """
    return self._stub.ListCompatibleParts(
        service_pb2.ListCompatiblePartsRequest(
            action_type_names=action_type_names
        ),
        timeout=self._rpc_timeout_seconds,
    ).parts

  def list_parts(self) -> list[str]:
    """Lists all available parts.

    Returns:
      List of available parts.
    """
    return self._stub.ListParts(
        service_pb2.ListPartsRequest(), timeout=self._rpc_timeout_seconds
    ).parts

  def start_session(
      self, parts: list[str], context: Optional[context_pb2.Context] = None
  ) -> _session.Session:
    """Starts a new `Session` for the given parts.

    Context management is supported, and it is recommended to obtain the Session
    using the `with` statement. For example:

      with icon_client.start_session(["robot_arm", "robot_gripper"]) as session:
        # ...

    Otherwise, do not forget to call `end()` once done with the Session. For
    example:

      session = icon_client.start_session(["robot_arm", "robot_gripper"])
      try:
        # ...
      finally:
        session.end()

    Attempts to recreate the same Session without calling `end()` will cause an
    exception since parts are exclusive to a Session. Otherwise once the Python
    process ends, the Session will be cleaned up via garbage collection. Note
    that this is not always guaranteed, see
    https://docs.python.org/3.3/reference/datamodel.html). In a notebook
    environment such as Jupyter, this can be triggered by restarting the kernel.

    Args:
      parts: List of parts to control.
      context: The log context passed to the session. Needed to sync ICON logs
        to the cloud. In skills use `context.logging_context`.

    Returns:
      A new Session.

    Raises:
      grpc.RpcError: An error occurred while starting the `Session`.
    """
    return _session.Session(self._stub, parts, context)

  def enable(self) -> None:
    """Enables all parts on the server.

    Performs all steps necessary to get the parts ready to receive commands.
    Since the server auto-enables at startup, this is only needed after a
    call to disable().

    NOTE: Enabling a server is something the user does directly. DO NOT call
    this from library code automatically to make things more convenient. Human
    users must be able to rely on the robot to stay still unless they enable
    it.

    Raises:
      grpc.RpcError: An error occurred while enabling.
    """
    self._stub.Enable(
        service_pb2.EnableRequest(), timeout=self._rpc_timeout_seconds
    )

  def disable(self, group: HardwareGroup = HardwareGroup.ALL_HARDWARE) -> None:
    """Disables all parts on the server.

    NOTE: Disabling a server is something the user does directly. DO NOT call
    this from library code automatically to make things more convenient. Human
    users must be able to rely on the robot to stay enabled unless they
    explicitly disable it (or the robot encounters a fault).

    Ends all currently-active sessions.

    Args:
      group: The group of hardware to disable. `ALL_HARDWARE`: All hardware
        modules and parts will be disabled. This is the default.
        `OPERATIONAL_HARDWARE_ONLY`: Only operational hardware modules and parts
        that use them will be disabled. So, hardware modules that are configured
        with `IconMainConfig.hardware_config.cell_control_hardware` will be kept
        enabled (if they are enabled). One use case is to integrate cell-level
        control where operational robot hardware can be paused such that
        automatic mode is not needed, while still reading/writing input/output
        on a fieldbus hardware module for cell-level control.

    Raises:
      grpc.RpcError: An error occurred while disabling.
    """
    request = service_pb2.DisableRequest()
    if group == self.HardwareGroup.OPERATIONAL_HARDWARE_ONLY:
      request.group = (
          service_pb2.DisableRequest.HardwareGroup.OPERATIONAL_HARDWARE_ONLY
      )
    else:
      request.group = service_pb2.DisableRequest.HardwareGroup.ALL_HARDWARE
    self._stub.Disable(request, timeout=self._rpc_timeout_seconds)

  def clear_faults(self) -> None:
    """Clears all faults and returns the server to an enabled state.

    NOTE: Clearing faults is something the user does directly. DO NOT call this
    from library code automatically to make things more convenient, ESPECIALLY
    not in connection with re-enabling the server afterwards! Human users must
    be able to rely on the robot to stay still unless they explicitly clear the
    fault(s) and enable it again.

    Some classes of faults (internal server errors or issues that have a
    physical root cause) may require additional server- or hardware-specific
    mitigation before clear_faults can successfully clear the fault.
    Returns RESOURCE_EXHAUSTED when a fatal fault is being cleared,
    which is not completed yet and involves a process restart. In this case,
    the client should retry until receiving OK.

    Raises:
      grpc.RpcError: An error occurred while clearing faults.
    """
    self._stub.ClearFaults(
        service_pb2.ClearFaultsRequest(), timeout=self._rpc_timeout_seconds
    )

  def get_operational_status(self) -> types_pb2.OperationalStatus:
    """Returns the summarized status of the server.

    This is the status of all hardware and the server.
    It can differ from `get_cell_control_hardware_status`, which is the state of
    a subset of hardware.

    This status may indicate that the server is ENABLED, DISABLED, or FAULTED.
    If FAULTED, OperationalStatus also includes a string explaining why the
    robot faulted.

    Returns:
      The operational status of the server.

    Raises:
      grpc.RpcError: An error occurred while getting the state.
    """
    resp = self._stub.GetOperationalStatus(
        service_pb2.GetOperationalStatusRequest(),
        timeout=self._rpc_timeout_seconds,
    )
    return resp.operational_status

  def get_cell_control_hardware_status(self) -> types_pb2.OperationalStatus:
    """Returns the status of the cell control hardware.

    Returns the status of cell control hardware, which is marked with
    `IconMainConfig.hardware_config.cell_control_hardware`.
    Cell control hardware is a group of hardware modules that does not inherit
    faults from operational hardware, and only gets disabled when any cell
    control hardware module faults (or when `disable` is called).

    This status may indicate that the server is ENABLED, DISABLED, or FAULTED.
    If FAULTED, OperationalStatus also includes a string explaining why the
    robot faulted.

    Returns:
      The status of the cell control hardware.

    Raises:
      grpc.RpcError: An error occurred while getting the state.
    """
    resp = self._stub.GetOperationalStatus(
        service_pb2.GetOperationalStatusRequest(),
        timeout=self._rpc_timeout_seconds,
    )
    return resp.cell_control_hardware_status

  def get_speed_override(self) -> float:
    """Returns the current speed override value.

    This is a value between 0 and 1, and acts as a multiplier to the speed of
    compatible actions.
    """
    resp = self._stub.GetSpeedOverride(
        service_pb2.GetSpeedOverrideRequest(), timeout=self._rpc_timeout_seconds
    )
    return resp.override_factor

  def set_speed_override(self, new_speed_override: float) -> None:
    """Sets the speed override value.

    Args:
      new_speed_override: A value between 0 and 1. Compatible actions will do
        their best to scale their speed.

    Raises:
      grpc.RpcError on errors, including invalid values
    """
    self._stub.SetSpeedOverride(
        service_pb2.SetSpeedOverrideRequest(override_factor=new_speed_override),
        timeout=self._rpc_timeout_seconds,
    )

  def get_logging_mode(self) -> logging_mode_pb2.LoggingMode:
    """Gets the logging mode."""
    return self._stub.GetLoggingMode(
        service_pb2.GetLoggingModeRequest(), timeout=self._rpc_timeout_seconds
    ).logging_mode

  def set_logging_mode(
      self, logging_mode: logging_mode_pb2.LoggingMode
  ) -> None:
    """Sets the logging mode.

    The logging mode defines which robot-status logs are logged to the cloud.
    ICON logs only to the cloud if a session is active. PuSub is not influenced
    by this setting.

    Args:
      logging_mode: The logging mode to set.
    """
    self._stub.SetLoggingMode(
        service_pb2.SetLoggingModeRequest(logging_mode=logging_mode),
        timeout=self._rpc_timeout_seconds,
    )

  def get_part_properties(self) -> service_pb2.GetPartPropertiesResponse:
    """Gets the values of all part properties.

    Returns:
      A GetPartPropertiesResponse proto that contains:
      * The control timestamp at the time the properties were reported
      * The wall time at the time the properties were reported
      * A map from part name to a map from property name to value
        For instance: {'robot': {'motor_0_current_amps': 2.0}}

    Raises:
      grpc.RpcError: The ICON server responds with an error. See message for
        details.
    """
    return self._stub.GetPartProperties(
        service_pb2.GetPartPropertiesRequest(),
        timeout=self._rpc_timeout_seconds,
    )

  def set_part_properties(
      self, part_properties: Mapping[str, Mapping[str, Union[bool, float]]]
  ) -> None:
    """Sets part properties.

    Check the output of get_part_properties to learn the available properties
    and their types.

    Args:
      part_properties: A map from part name to a map from property name to
        value. For instance: {'robot': {'internal_controller_p_value': 0.2}}

    Raises:
      grpc.RpcError: Server responded with an error. Common errors include
        unknown part or property names, or wrong property types.
    """
    request = service_pb2.SetPartPropertiesRequest()
    for part_name, properties in part_properties.items():
      properties_proto = service_pb2.PartPropertyValues()
      for property_name, property_value in properties.items():
        value_proto = service_pb2.PartPropertyValue()
        if isinstance(property_value, bool):
          value_proto.bool_value = property_value
        if isinstance(property_value, float):
          value_proto.double_value = property_value
        properties_proto.property_values_by_name[property_name].CopyFrom(
            value_proto
        )
      request.part_properties_by_part_name[part_name].CopyFrom(properties_proto)
    self._stub.SetPartProperties(request, timeout=self._rpc_timeout_seconds)

  def set_payload(
      self,
      payload: robot_payload.RobotPayload,
      part_name: str,
      payload_name: str,
  ) -> None:
    """Sets a payload of a part.

    After setting the payload, ICON disables and then re-enables all parts.

    Args:
      payload: The payload to set.
      part_name: The name of the part to set the payload for.
      payload_name: The name of the payload to set.
    """
    request = service_pb2.SetPayloadRequest(
        payload=robot_payload.payload_to_proto(payload),
        part_name=part_name,
        payload_name=payload_name,
    )
    self._stub.SetPayload(request, timeout=self._rpc_timeout_seconds)

  def get_payload(
      self, part_name: str, payload_name: str
  ) -> robot_payload.RobotPayload:
    """Gets a payload of a part.

    Args:
      part_name: The name of the part to get the payload for.
      payload_name: The name of the payload to get.

    Returns:
      The requested payload.
    """
    request = service_pb2.GetPayloadRequest(
        part_name=part_name,
        payload_name=payload_name,
    )
    response = self._stub.GetPayload(request, timeout=self._rpc_timeout_seconds)
    return robot_payload.payload_from_proto(response.payload)
