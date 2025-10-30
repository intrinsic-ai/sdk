# Copyright 2023 Intrinsic Innovation LLC

"""Provides utility functions for Asset dependencies."""

import contextlib

from google.protobuf import any_pb2
import grpc
from intrinsic.assets.data.proto.v1 import data_assets_pb2
from intrinsic.assets.data.proto.v1 import data_assets_pb2_grpc
from intrinsic.assets.proto.v1 import resolved_dependency_pb2
from intrinsic.util.grpc import interceptor

_INGRESS_ADDRESS = "istio-ingressgateway.app-ingress.svc.cluster.local:80"


class MissingInterfaceError(ValueError):
  """Raised when an interface is not found in the resolved dependency."""


class NotDataError(ValueError):
  """Raised when an interface is not data."""


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
  iface_proto = _find_interface(dep, iface)
  if not iface_proto.HasField("grpc_connection"):
    raise NotGRPCError(
        "Interface is not gRPC or no connection information is available:"
        f" {iface}"
    )

  metadata = iface_proto.grpc_connection.metadata
  channel = grpc.intercept_channel(
      grpc.insecure_channel(iface_proto.grpc_connection.address),
      interceptor.HeaderAdderInterceptor(
          lambda: [(m.key, m.value) for m in metadata]
      ),
  )

  return channel


def get_data_payload(
    dep: resolved_dependency_pb2.ResolvedDependency,
    iface: str,
    data_assets_client: data_assets_pb2_grpc.DataAssetsStub | None = None,
) -> any_pb2.Any:
  """Returns the DataAsset payload for the specified interface.

  If no DataAssets client is provided, an insecure connection to the DataAssets
  service via the ingress gateway will be created. This connection is valid for
  services running in the same cluster as the DataAssets service.

  Args:
    dep: The resolved dependency.
    iface: The interface to get the data payload for.
    data_assets_client: The DataAssets client to use.

  Returns:
    The DataAsset payload for the specified interface.

  Raises:
    MissingInterfaceError: If the specified interface is not found in the
      resolved dependency.
    NotDataError: If the specified interface is not data or no data dependency
      information is available.
  """
  iface_proto = _find_interface(dep, iface)
  if not iface_proto.HasField("data"):
    raise NotDataError(
        "Interface is not data or no data dependency information is available:"
        f" {iface}"
    )

  channel = contextlib.nullcontext()
  if data_assets_client is None:
    data_assets_client, channel = _make_default_data_assets_client()

  with channel:
    da = data_assets_client.GetDataAsset(
        data_assets_pb2.GetDataAssetRequest(id=iface_proto.data.id)
    )

  return da.data


def _find_interface(
    dep: resolved_dependency_pb2.ResolvedDependency, iface: str
) -> resolved_dependency_pb2.ResolvedDependency.Interface:
  """Returns the interface for the specified interface."""
  if iface not in dep.interfaces:
    if not dep.interfaces:
      explanation = "no interfaces provided"
    else:
      keys = ", ".join(dep.interfaces.keys())
      explanation = f"got interfaces: {keys}"
    raise MissingInterfaceError(
        f"Interface not found in resolved dependency (want {iface},"
        f" {explanation})"
    )
  return dep.interfaces[iface]


def _make_default_data_assets_client() -> (
    tuple[data_assets_pb2_grpc.DataAssetsStub, grpc.Channel]
):
  """Creates an insecure channel to the DataAssets service via the ingress gateway."""
  channel = grpc.insecure_channel(_INGRESS_ADDRESS)
  return data_assets_pb2_grpc.DataAssetsStub(channel), channel
