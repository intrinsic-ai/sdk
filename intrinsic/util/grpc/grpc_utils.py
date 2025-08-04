# Copyright 2023 Intrinsic Innovation LLC

"""Utility functions for gRPC connections."""

from typing import Optional

import grpc
from intrinsic.config import environments  # For service discovery
from intrinsic.kubernetes.acl.ipcidentity import ipcidentity
from intrinsic.kubernetes.acl.ipcidentity.internal import metadata
from intrinsic.util.grpc import interceptor

# Create a module-level instance of IpcIdentity. This ensures that the same
# instance is used throughout the application, which avoids the overhead of
# creating new instances and ensures that the same token is used for all
# requests.
_shared_ipc_identity = ipcidentity.IpcIdentity()


def _get_compute_project() -> str:
  """Gets the compute project from the metadata server.

  Returns:
    The compute project.
  """
  return metadata.MetadataClient().compute_project()


def _get_environment() -> str:
  """Gets the environment based on the compute project.

  The environment is either prod, dev or staging. It is important to
  consistently use the correct environment because each environment has
  associated its dedicated account, assets, compute projects and virtual
  machines.

  This function is supposed to be used within a cluster and first determines the
  compute project from the metadata server. Then it uses the compute project to
  determine the environment.

  Returns:
    The environment based on the compute project.
  """
  compute_project = _get_compute_project()
  return environments.from_compute_project(compute_project)


def _create_auth_interceptor(
    ipc_identity: ipcidentity.IpcIdentity,
) -> interceptor.ClientCallDetailsInterceptor:
  """Creates an interceptor that adds an auth header to the gRPC channel."""
  auth_header = lambda: [("cookie", "auth-proxy=" + ipc_identity.token())]
  return interceptor.HeaderAdderInterceptor(auth_header)


def _add_auth_header(
    channel: grpc.Channel,
    ipc_identity: Optional[ipcidentity.IpcIdentity],
) -> grpc.Channel:
  """Adds an auth header to the gRPC channel."""
  if ipc_identity:
    return grpc.intercept_channel(
        channel, _create_auth_interceptor(ipc_identity)
    )
  return channel


def _create_channel(address: str) -> grpc.Channel:
  """Creates a gRPC channel with the given address."""
  return grpc.secure_channel(
      address,
      grpc.ssl_channel_credentials(),
      options=[("grpc.max_receive_message_length", -1)],
  )


def create_assets_channel() -> grpc.Channel:
  """Creates a gRPC channel for the assets service with an IPC ID auth header.

  Returns:
    The gRPC channel for the assets service.
  """
  env = _get_environment()
  assets_domain = environments.assets_domain(env)
  assets_channel = _create_channel(f"{assets_domain}:443")
  return _add_auth_header(assets_channel, _shared_ipc_identity)


def create_cloud_channel() -> grpc.Channel:
  """Creates a gRPC channel for the cloud service with an IPC ID auth header.

  Returns:
    The gRPC channel for the cloud service.
  """
  compute_project = _get_compute_project()
  cloud_channel = _create_channel(
      f"www.endpoints.{compute_project}.cloud.goog:443"
  )
  return _add_auth_header(cloud_channel, _shared_ipc_identity)


def create_ingress_channel() -> grpc.Channel:
  """Creates a gRPC channel for in-cluster calls with an IPC ID auth header.

  Returns:
    The gRPC channel for in-cluster calls.
  """
  ingress_channel = grpc.insecure_channel(
      "istio-ingressgateway.app-ingress.svc.cluster.local:80",
      options=[("grpc.max_receive_message_length", -1)],
  )
  return _add_auth_header(ingress_channel, _shared_ipc_identity)
