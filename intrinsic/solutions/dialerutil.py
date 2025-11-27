# Copyright 2023 Intrinsic Innovation LLC

"""Provides methods to API key autorized gRPC calls.

This file implements a subset of the
`//intrinsic/skills/tools/skill/cmd/dialerutil.go` library.
"""

from typing import Any

import grpc
from intrinsic.frontend.cloud.api.v1 import solutiondiscovery_api_pb2
from intrinsic.frontend.cloud.api.v1 import solutiondiscovery_api_pb2_grpc
from intrinsic.kubernetes.acl.py import identity
from intrinsic.solutions import auth


class _TokenAuth(grpc.AuthMetadataPlugin):
  """gRPC Metadata Plugin that adds an API key to the header."""

  _token: auth.ProjectToken

  def __init__(self, token: auth.ProjectToken):
    self._token = token

  def __call__(self, context, callback):
    callback(self._token.get_request_metadata(), None)


class _AuthProxy(grpc.AuthMetadataPlugin):
  """gRPC Metadata Plugin that adds an auth-proxy cookie."""

  _token: str

  def __init__(self, token: str):
    self._token = token

  def __call__(self, context, callback):
    callback((("cookie", f"auth-proxy={self._token}"),), None)


class _ServerName(grpc.AuthMetadataPlugin):
  """gRPC Metadata Plugin that adds the cluster name to the header."""

  _server_name: str

  def __init__(self, server_name: str):
    self._server_name = server_name

  def __call__(self, context, callback):
    callback((("x-server-name", self._server_name),), None)


def create_channel_from_address(
    address: str,
    grpc_options: list[tuple[str, Any]] | None = None,
) -> grpc.Channel:
  """Creates a gRPC channel based on the provided address.

  Args:
    address: The address of the server to connect to.
    grpc_options: List of gRPC channel options.

  Returns:
    A gRPC channel to the provided address.
  """
  return grpc.insecure_channel(address, options=grpc_options)


def create_channel_from_cluster(
    org_info: auth.OrgInfo,
    cluster: str,
    grpc_options: list[tuple[str, Any]] | None = None,
) -> grpc.Channel:
  """Creates a gRPC channel based on the provided cluster.

  Args:
    org_info: The organization and project information.
    cluster: The name of the cluster.
    grpc_options: List of gRPC channel options.

  Returns:
    A gRPC channel to the cluster.
  """
  return _create_channel(
      org_info=org_info,
      cluster=cluster,
      grpc_options=grpc_options,
  )


def _get_cluster_from_solution(
    org_info: auth.OrgInfo,
    solution: str,
    grpc_options: list[tuple[str, Any]] | None = None,
) -> str:
  """Returns the name of the cluster in which the given solution is running.

  Args:
    org_info: The organization and project information.
    solution: The name of the solution.
    grpc_options: List of gRPC channel options.

  Returns:
    The name of the cluster in which the given solution is running.
  """
  # Open a temporary gRPC channel to the cloud cluster to resolve the cluster
  # on which the solution is running.
  channel = _create_channel(
      org_info=org_info,
      grpc_options=grpc_options,
  )
  stub = solutiondiscovery_api_pb2_grpc.SolutionDiscoveryServiceStub(channel)
  response = stub.GetSolutionDescription(
      solutiondiscovery_api_pb2.GetSolutionDescriptionRequest(name=solution)
  )
  channel.close()

  return response.solution.cluster_name


def create_channel_from_solution(
    org_info: auth.OrgInfo,
    solution: str,
    grpc_options: list[tuple[str, Any]] | None = None,
) -> grpc.Channel:
  """Creates a gRPC channel based on the provided solution.

  Args:
    org_info: The organization and project information.
    solution: The name of the solution.
    grpc_options: List of gRPC channel options.

  Returns:
    A gRPC channel to the cluster running the solution.
  """
  return _create_channel(
      org_info=org_info,
      cluster=_get_cluster_from_solution(org_info, solution, grpc_options),
      grpc_options=grpc_options,
  )


def create_channel_from_token(
    auth_token: str,
    org: str,
    cluster: str,
    grpc_options: list[tuple[str, Any]] | None = None,
) -> grpc.Channel:
  """Creates a gRPC channel based on the provided token.

  Args:
    auth_token: The auth-proxy token to use for authentication.
    org: The organization and project in the format <org>@<project>.
    cluster: The name of the cluster.
    grpc_options: List of gRPC channel options.

  Returns:
    A gRPC channel to the cluster.
  """
  return _create_channel(
      org_info=auth.parse_info_from_string(org),
      cluster=cluster,
      grpc_options=grpc_options,
      auth_token=auth_token,
  )


def _create_channel(
    org_info: auth.OrgInfo,
    cluster: str | None = None,
    grpc_options: list[tuple[str, Any]] | None = None,
    auth_token: str | None = None,
) -> grpc.Channel:
  """Creates a gRPC channel based on the provided connection parameters.

  Args:
    org_info: The organization and project information.
    cluster: The name of the cluster.
    grpc_options: List of gRPC channel options.
    auth_token: The auth-proxy token to use for authentication.

  Returns:
    A gRPC channel to the cluster.
  """
  channel_credentials = grpc.ssl_channel_credentials()
  call_credentials = []

  if auth_token is None:
    token = auth.get_configuration(org_info.project).get_default_credentials()
    call_credentials.append(
        grpc.metadata_call_credentials(_TokenAuth(token), name="TokenAuth")
    )
  else:
    call_credentials.append(
        grpc.metadata_call_credentials(_AuthProxy(auth_token), name="AuthProxy")
    )

  call_credentials.append(
      identity.OrgNameCallCredentials(org_info.organization)
  )

  if cluster is not None:
    call_credentials.append(
        grpc.metadata_call_credentials(_ServerName(cluster), name="ServerName")
    )

  return grpc.secure_channel(
      f"dns:///www.endpoints.{org_info.project}.cloud.goog:443",
      grpc.composite_channel_credentials(
          channel_credentials, *call_credentials
      ),
      options=grpc_options,
  )
