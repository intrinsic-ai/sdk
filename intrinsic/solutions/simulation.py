# Copyright 2023 Intrinsic Innovation LLC

"""Provides functionality to interact with a running simulation.

Typical usage example:
  from intrinsic.executive.jupyter.workcell import intrinsic

  workcell = intrinsic.connect()
  simulation = workcell.simulation

  simulation.reset()
"""

import grpc
from intrinsic.math.python import proto_conversion
from intrinsic.simulation.service.proto import simulation_service_pb2
from intrinsic.simulation.service.proto import simulation_service_pb2_grpc
from intrinsic.solutions import errors
from intrinsic.util.grpc import error_handling
from intrinsic.world.proto import object_world_service_pb2
from intrinsic.world.proto import object_world_service_pb2_grpc


SimulationServiceStub = simulation_service_pb2_grpc.SimulationServiceStub
ObjectWorldServiceStub = object_world_service_pb2_grpc.ObjectWorldServiceStub



class Simulation:
  """Provides commands to interact with a running simulation."""

  def __init__(
      self,
      simulation_service: SimulationServiceStub,
      object_world_service: ObjectWorldServiceStub,
  ):
    """Constructs a new Simulation object.

    Args:
      simulation_service: The gRPC stub to be used for communication with the
        simulation service.
      object_world_service: The gRPC stub to be used for communication with the
        object world service.
    """
    self._simulation_service: SimulationServiceStub = simulation_service
    self._object_world_service: ObjectWorldServiceStub = object_world_service

  @classmethod
  def connect(cls, grpc_channel: grpc.Channel) -> 'Simulation':
    """Create a Simulation instance using the given gRPC address.

    Args:
      grpc_channel: Address of the simulation gRPC service to use.

    Returns:
      A newly created Simulation instance.
    """
    simulation_service = SimulationServiceStub(grpc_channel)
    object_world_service = ObjectWorldServiceStub(grpc_channel)
    return cls(simulation_service, object_world_service)

  @error_handling.retry_on_grpc_unavailable
  def _call_simulation_service_reset(
      self, request: simulation_service_pb2.ResetSimulationRequest
  ):
    return self._simulation_service.ResetSimulation(request)

  @error_handling.retry_on_grpc_unavailable
  def reset(self) -> None:
    """Resets the simulation world to its initial state.

    Also makes sure that all affected components such as ICON are in a working
    state.
    """
    request = simulation_service_pb2.ResetSimulationRequest()
    self._call_simulation_service_reset(request)
