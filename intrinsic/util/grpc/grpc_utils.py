# Copyright 2026 Intrinsic Innovation LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Utility functions for gRPC connections."""

from typing import Optional

import grpc

from intrinsic.kubernetes.acl.ipcidentity import ipcidentity
from intrinsic.kubernetes.acl.ipcidentity.internal import metadata
from intrinsic.kubernetes.acl.py import identity
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


def _create_auth_interceptor(
    ipc_identity: ipcidentity.IpcIdentity,
) -> interceptor.ClientCallDetailsInterceptor:
  """Creates an interceptor that adds an auth header to the gRPC channel."""
  return interceptor.HeaderAdderInterceptor(ipc_identity.metadata)


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


def create_assets_channel(
    assets_domain: str = "assets.intrinsic.ai",
) -> grpc.Channel:
  """Creates a gRPC channel for the assets service with an IPC ID auth header.

  Args:
    assets_domain: Hostname of the assets service (e.g. "assets.intrinsic.ai",
      "assets-qa.intrinsic.ai", or "assets-dev.intrinsic.ai").

  Returns:
    The gRPC channel for the assets service.
  """
  if not assets_domain:
    raise ValueError("assets_domain must not be empty.")
  assets_channel = _create_channel(f"{assets_domain}:443")
  return _add_auth_header(assets_channel, _shared_ipc_identity)


def create_cloud_channel() -> grpc.Channel:
  """Creates a gRPC channel for the cloud service with an IPC ID auth header.

  The organization ID is automatically derived from the IPC identity token.

  The returned channel has max_receive_message_length set to unlimited.

  Returns:
    The gRPC channel for the cloud service.
  """
  compute_project = _get_compute_project()
  address = f"www.endpoints.{compute_project}.cloud.goog:443"
  cloud_channel = _create_channel(address)
  return _add_auth_header(cloud_channel, _shared_ipc_identity)


def create_ingress_channel() -> grpc.Channel:
  """Creates a gRPC channel for in-cluster calls.

  This method is preferred for in-cluster communication when authentication
  is not required, as it avoids the overhead of retrieving an IPC identity
  (which requires an internet connection).

  Returns:
    The gRPC channel for in-cluster calls.
  """
  return grpc.insecure_channel(
      "istio-ingressgateway.app-ingress.svc.cluster.local:80",
      options=[("grpc.max_receive_message_length", -1)],
  )


def create_ingress_channel_with_ipc_identity() -> grpc.Channel:
  """Creates a gRPC channel for in-cluster calls with an IPC identity auth header.

  Use of this method is discouraged. Typically, services should communicate
  directly: use `create_ingress_channel` for in-cluster communication, and
  `create_cloud_channel` to communicate directly with the cloud.

  This method is only necessary when a service must communicate with another
  in-cluster service that proxies the request to the cloud, requiring the
  original caller's IPC identity to be forwarded to authenticate with the cloud.

  Returns:
    The gRPC channel for in-cluster calls.
  """
  return _add_auth_header(create_ingress_channel(), _shared_ipc_identity)
