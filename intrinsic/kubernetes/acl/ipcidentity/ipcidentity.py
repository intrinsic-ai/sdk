# Copyright 2023 Intrinsic Innovation LLC

"""Provides an IPC identity token by exchanging a robot JWT via AccountsTokensService."""

import datetime
import logging
from typing import Optional

import grpc  # For gRPC communication
from intrinsic.config import environments  # For service discovery
from intrinsic.kubernetes.accounts.service.api.tokens.v2 import tokens_pb2
from intrinsic.kubernetes.accounts.service.api.tokens.v2 import tokens_pb2_grpc
from intrinsic.kubernetes.acl.ipcidentity.internal import metadata
from intrinsic.kubernetes.acl.py import jwt

TOKEN_EXPIRY_MARGIN = datetime.timedelta(seconds=30)
GRPC_TIMEOUT_SECONDS = 10
GRPC_USER_AGENT = "ipcidentity/1.0 (py, grpc)"


class IpcIdentity:
  """Provides an IPC identity token.

  Retrieves a robot JWT from the MetaDataClient and then calls the
  AccountsTokensService via gRPC to exchange it for an IPC identity token.
  The class maintains that token and its expiry data internally.
  """

  def __init__(
      self,
      metadata_client: Optional[metadata.MetadataClient] = None,
  ):
    self._metadata_client = metadata_client or metadata.MetadataClient()
    self._token = None
    self._expires = None
    self._grpc_channel: Optional[grpc.Channel] = None
    self._tokens_stub: Optional[tokens_pb2_grpc.AccountsTokensServiceStub] = (
        None
    )

  def _ensure_grpc_setup(self) -> None:
    """Ensures gRPC channel and stub are initialized."""
    if self._tokens_stub and self._grpc_channel:
      return

    logging.info("gRPC channel and stub not initialized, setting up.")

    compute_project = self._metadata_client.compute_project()
    env = environments.from_compute_project(compute_project)
    accounts_service_domain = environments.accounts_domain(env)

    if not accounts_service_domain:
      raise RuntimeError(
          f"Could not determine accounts domain for env '{env}' from project"
          f" '{compute_project}'"
      )

    target_address = f"{accounts_service_domain}:443"  # Standard gRPC TLS port
    logging.info("Connecting to AccountsTokensService at: %s", target_address)

    self._grpc_channel = grpc.secure_channel(
        target_address,
        grpc.ssl_channel_credentials(),
        # https://github.com/grpc/grpc/issues/23644#issuecomment-669344588
        options=[("grpc.primary_user_agent", GRPC_USER_AGENT)],
    )
    self._tokens_stub = tokens_pb2_grpc.AccountsTokensServiceStub(
        self._grpc_channel
    )

  def _fetch_ipc_identity_token(self) -> None:
    """Fetches the IPC identity token via AccountsTokensService."""
    self._ensure_grpc_setup()

    if not self._tokens_stub:  # Should be set by _ensure_grpc_setup
      raise RuntimeError("gRPC stub not initialized.")

    robot_jwt_str = self._metadata_client.token()
    if not robot_jwt_str:
      raise RuntimeError("Failed to get robot JWT from metadata client.")

    logging.debug("Obtained robot JWT, preparing GetIPCTokenRequest.")

    request = tokens_pb2.GetIPCTokenRequest(
        credential=tokens_pb2.IPCCredential(
            robot_jwt=tokens_pb2.RobotJWT(jwt=robot_jwt_str)
        )
    )

    try:
      logging.debug("Sending GetIPCTokenRequest to AccountsTokensService.")
      response = self._tokens_stub.GetIPCToken(
          request, timeout=GRPC_TIMEOUT_SECONDS
      )
    except grpc.RpcError as e:
      raise RuntimeError(
          f"Failed to get IPC token from AccountsTokensService: {e}"
      ) from e

    ipc_token_str = response.ipc_token
    if not ipc_token_str:
      raise RuntimeError("AccountsTokensService returned an empty IPC token.")

    self._token = ipc_token_str
    logging.info("Obtained IPC token: %s", self._token)
    jwt_values = jwt.PayloadUnsafe(self._token)
    if "exp" not in jwt_values:
      raise RuntimeError("IPC token does not contain an expiration time.")
    self._expires = datetime.datetime.fromtimestamp(
        float(jwt_values["exp"]), tz=datetime.timezone.utc
    )
    logging.info(
        "IPC token received for '%s', expires at: %s",
        jwt_values["user_id"],
        self._expires,
    )

  def token(self) -> str:
    """Returns the IPC identity token.

    Returns:
      The IPC identity token.

    Raises:
      RuntimeError: If the token could not be fetched or is invalid.
    """
    # Ensure datetime.now() is timezone-aware if self._expires is.
    # datetime.datetime.fromtimestamp(ts, tz=datetime.timezone.utc) makes it
    # offset-aware.
    now_utc = datetime.datetime.now(datetime.timezone.utc)
    if (
        self._token is not None
        and self._expires is not None
        and self._expires > now_utc + TOKEN_EXPIRY_MARGIN
    ):
      logging.debug("Returning cached IPC token.")
      return self._token

    logging.info("Cached IPC token is invalid or expired, fetching a new one.")
    try:
      self._fetch_ipc_identity_token()
    except Exception as e:
      raise RuntimeError(f"Failed to fetch IPC identity token: {e}") from e

    if (
        not self._token
    ):  # Should be set by _fetch_ipc_identity_token if successful
      raise RuntimeError(
          "Failed to obtain a valid IPC token after refresh attempt."
      )
    return self._token

  def close(self) -> None:
    """Closes any open resources, like the gRPC channel."""
    if self._grpc_channel:
      logging.info("Closing gRPC channel.")
      self._grpc_channel.close()
      self._grpc_channel = None
      self._tokens_stub = None

  def __enter__(self):
    return self

  def __exit__(self, exc_type, exc_val, exc_tb):
    self.close()
