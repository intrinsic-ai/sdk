# Copyright 2023 Intrinsic Innovation LLC

"""Provides utility functions for Asset dependencies."""

import grpc
from intrinsic.assets.proto.v1 import resolved_dependency_pb2
from intrinsic.util.grpc import interceptor


class MissingInterfaceError(ValueError):
  """Raised when an interface is not found in the resolved dependency."""


class NotGRPCError(ValueError):
  """Raised when an interface is not gRPC."""


def connect(
    dep: resolved_dependency_pb2.ResolvedDependency,
    iface: str,
) -> grpc.Channel:
  """Creates a gRPC channel to the provider of the specified interface.

  The returned channel will be intercepted to include any needed metadata for
  communicating with the provider.

  Args:
    dep: The resolved dependency.
    iface: The interface to connect to.

  Returns:
    A gRPC channel to the provider of the specified interface.

  Raises:
    MissingInterfaceError: If the specified interface is not found in the
      resolved dependency.
    NotGRPCError: If the specified interface is not gRPC.
  """
  if iface not in dep.interfaces:
    raise MissingInterfaceError(
        f"Interface {iface} not found in resolved dependency"
    )
  iface_proto = dep.interfaces[iface]

  if not iface_proto.HasField("grpc_connection"):
    raise NotGRPCError(
        f"Interface {iface} is not gRPC or no connection information is"
        " available."
    )

  metadata = iface_proto.grpc_connection.metadata
  channel = grpc.intercept_channel(
      grpc.insecure_channel(iface_proto.grpc_connection.address),
      interceptor.HeaderAdderInterceptor(
          lambda: [(m.key, m.value) for m in metadata]
      ),
  )

  return channel
