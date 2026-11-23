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

"""Provides an IPC identity token via xfa-config-service."""

from collections.abc import Callable
import datetime
import http
import http.client
import json
import logging
import pathlib
import ssl
import threading
from typing import Optional

from intrinsic.kubernetes.acl.py import identity
from intrinsic.kubernetes.acl.py import jwt

XFA_CONFIG_HOST = "xfa-config-control-plane.default.svc.cluster.local"
XFA_CONFIG_PORT = 443
XFA_CONFIG_PATH = "/v1alpha1/ipcidentity"
K8S_SA_TOKEN_PATH = "/var/run/secrets/kubernetes.io/serviceaccount/token"

TOKEN_EXPIRY_MARGIN = datetime.timedelta(seconds=30)
HTTP_TIMEOUT_SECONDS = 10


def _read_k8s_sa_token(token_path: str) -> str:
  """Reads the Kubernetes ServiceAccount bearer token from disk."""
  try:
    token = pathlib.Path(token_path).read_text(encoding="utf-8").strip()
  except OSError as e:
    raise RuntimeError(
        f"Failed to read bearer token from {token_path}: {e}"
    ) from e
  if not token:
    raise RuntimeError(f"Bearer token file {token_path} is empty.")
  return token


def _default_http_connection() -> http.client.HTTPSConnection:
  """Creates an HTTPS connection to xfa-config-service."""
  # xfa-config-control-plane uses a self-signed certificate; match the C++ and
  # Go ipcidentity clients by skipping server certificate verification.
  return http.client.HTTPSConnection(
      XFA_CONFIG_HOST,
      XFA_CONFIG_PORT,
      timeout=HTTP_TIMEOUT_SECONDS,
      context=ssl._create_unverified_context(),  # pylint: disable=protected-access
  )


class IpcIdentity:
  """Provides an IPC identity token from xfa-config-service.

  Retrieves the pod's Kubernetes ServiceAccount bearer token and requests an
  IPC identity token from the cluster-local xfa-config-service endpoint.
  The token and its expiry are cached in-process and refreshed automatically.
  """

  def __init__(
      self,
      http_connection_factory: Optional[
          Callable[[], http.client.HTTPConnection]
      ] = None,
      bearer_token_source: Optional[Callable[[], str]] = None,
      token_path: str = K8S_SA_TOKEN_PATH,
  ):
    self._http_connection_factory = (
        http_connection_factory or _default_http_connection
    )
    self._bearer_token_source = bearer_token_source or (
        lambda: _read_k8s_sa_token(token_path)
    )
    self._lock = threading.Lock()
    self._token: Optional[str] = None
    self._expires: Optional[datetime.datetime] = None

  def _fetch_ipc_identity_token(self) -> None:
    """Fetches the IPC identity token from xfa-config-service."""
    bearer_token = self._bearer_token_source()
    if not bearer_token:
      raise RuntimeError("Failed to get Kubernetes ServiceAccount token.")

    conn = self._http_connection_factory()
    try:
      conn.request(
          "POST",
          XFA_CONFIG_PATH,
          body="{}",
          headers={
              "Authorization": f"Bearer {bearer_token}",
              "Content-Type": "application/json",
          },
      )
      response = conn.getresponse()
      body = response.read().decode("utf-8", errors="replace")
    finally:
      conn.close()

    if response.status != http.HTTPStatus.OK:
      raise RuntimeError(
          f"xfa-config-service returned status code {response.status}:"
          f" {body[:100]}"
      )

    try:
      payload = json.loads(body)
    except json.JSONDecodeError as e:
      raise RuntimeError(
          f"xfa-config-service returned invalid JSON: {e}"
      ) from e

    if not isinstance(payload, dict):
      raise RuntimeError("xfa-config-service response must be a JSON object.")

    ipc_token_str = payload.get("token")
    if not isinstance(ipc_token_str, str) or not ipc_token_str:
      raise RuntimeError(
          "xfa-config-service response did not contain a valid 'token' field."
      )

    expires = payload.get("expires")
    if isinstance(expires, bool) or not isinstance(expires, (int, float)):
      raise RuntimeError(
          "xfa-config-service response did not contain a valid 'expires' field."
      )

    self._token = ipc_token_str
    self._expires = datetime.datetime.fromtimestamp(
        float(expires), tz=datetime.timezone.utc
    )
    logging.info("IPC token received, expires at: %s", self._expires)

  def token(self) -> str:
    """Returns the IPC identity token.

    Returns:
      The IPC identity token.

    Raises:
      RuntimeError: If the token could not be fetched or is invalid.
    """
    with self._lock:
      now_utc = datetime.datetime.now(datetime.timezone.utc)
      if (
          self._token is not None
          and self._expires is not None
          and self._expires > now_utc + TOKEN_EXPIRY_MARGIN
      ):
        logging.debug("Returning cached IPC token.")
        return self._token

      logging.info(
          "Cached IPC token is invalid or expired, fetching a new one."
      )
      try:
        self._fetch_ipc_identity_token()
      except Exception as e:
        raise RuntimeError(f"Failed to fetch IPC identity token: {e}") from e

      if not self._token:
        raise RuntimeError(
            "Failed to obtain a valid IPC token after refresh attempt."
        )
      return self._token

  def metadata(self) -> list[tuple[str, str]]:
    """Returns the IPC identity token and organization ID as gRPC metadata.

    The metadata includes the IPC identity token as an 'auth-proxy' cookie,
    and the organization identifier as an 'org-id' cookie and
    'x-intrinsic-org' header.

    Returns:
      A list of (key, value) pairs for gRPC metadata.

    Raises:
      RuntimeError: If the token or orgid could not be retrieved.
    """
    token = self.token()
    payload = jwt.PayloadUnsafe(token)
    org_id = payload.get("orgid")
    if not org_id:
      raise RuntimeError("IPC token does not contain an 'orgid' claim.")

    res = identity.OrgIDToGRPCMetadata(org_id)
    res.append(
        (identity.COOKIE_KEY, f"{identity.AUTH_PROXY_COOKIE_NAME}={token}")
    )
    return res

  def close(self) -> None:
    """Closes any open resources (no-op; HTTP connections close per fetch)."""

  def __enter__(self):
    return self

  def __exit__(self, exc_type, exc_val, exc_tb):
    self.close()
