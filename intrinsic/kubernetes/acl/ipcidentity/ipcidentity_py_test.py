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

"""Tests for ipcidentity."""

import http
import http.client
import json
import os
import tempfile
import time
from typing import Any
from unittest import mock

from absl.testing import absltest

from intrinsic.kubernetes.acl.ipcidentity import ipcidentity


def create_test_jwt(payload: dict[str, Any]) -> str:
  """Creates a minimal JWT string parseable by jwt.PayloadUnsafe."""
  header_b64 = (  # {"alg":"RS256","typ":"JWT"}
      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9"
  )
  payload_json = json.dumps(payload)
  payload_b64 = (
      ipcidentity.jwt.base64.urlsafe_b64encode(payload_json.encode("utf-8"))
      .decode("utf-8")
      .rstrip("=")
  )
  signature_b64 = "dummySignature"
  return f"{header_b64}.{payload_b64}.{signature_b64}"


def _make_mock_http_response(status: int, body: str) -> mock.MagicMock:
  response = mock.MagicMock(spec=http.client.HTTPResponse)
  response.status = status
  response.read.return_value = body.encode("utf-8")
  return response


def _setup_mock_provider(
    responses: list[mock.MagicMock],
    bearer_token: str = "test-k8s-sa-token",
) -> tuple[ipcidentity.IpcIdentity, mock.MagicMock, mock.MagicMock]:
  """Creates an IpcIdentity backed by a mocked HTTPConnection factory."""
  mock_conn = mock.MagicMock(spec=http.client.HTTPSConnection)
  mock_conn.getresponse.side_effect = responses
  conn_factory = mock.MagicMock(return_value=mock_conn)
  bearer_source = mock.MagicMock(return_value=bearer_token)
  provider = ipcidentity.IpcIdentity(
      http_connection_factory=conn_factory,
      bearer_token_source=bearer_source,
  )
  return provider, mock_conn, bearer_source


class IpcidentityTest(absltest.TestCase):
  """Tests for ipcidentity."""

  def test_token_retrieval_success_and_caching(self):
    expiration_time = int(time.time() + 3600)
    expected_token = create_test_jwt({
        "exp": expiration_time,
        "user_id": "test_user@ipc.intrinsic.ai",
        "orgid": "test_org",
    })
    response = _make_mock_http_response(
        http.HTTPStatus.OK,
        json.dumps({"token": expected_token, "expires": expiration_time}),
    )
    provider, mock_conn, bearer_source = _setup_mock_provider([response])

    actual_token = provider.token()

    self.assertEqual(actual_token, expected_token)
    bearer_source.assert_called_once()
    mock_conn.request.assert_called_once_with(
        "POST",
        ipcidentity.XFA_CONFIG_PATH,
        body="{}",
        headers={
            "Authorization": "Bearer test-k8s-sa-token",
            "Content-Type": "application/json",
        },
    )
    mock_conn.close.assert_called_once()

    # Subsequent call should return the cached token without issuing HTTP calls.
    cached_token = provider.token()
    self.assertEqual(cached_token, expected_token)
    mock_conn.request.assert_called_once()
    provider.close()

  def test_token_retrieval_expired_refreshes(self):
    expired_time = int(time.time() - 3600)
    expired_token = create_test_jwt({
        "exp": expired_time,
        "user_id": "test_user@ipc.intrinsic.ai",
    })
    valid_time = int(time.time() + 3600)
    valid_token = create_test_jwt({
        "exp": valid_time,
        "user_id": "test_user@ipc.intrinsic.ai",
    })
    responses = [
        _make_mock_http_response(
            http.HTTPStatus.OK,
            json.dumps({"token": expired_token, "expires": expired_time}),
        ),
        _make_mock_http_response(
            http.HTTPStatus.OK,
            json.dumps({"token": valid_token, "expires": valid_time}),
        ),
    ]
    provider, mock_conn, _ = _setup_mock_provider(responses)

    self.assertEqual(provider.token(), expired_token)
    self.assertEqual(mock_conn.request.call_count, 1)

    self.assertEqual(provider.token(), valid_token)
    self.assertEqual(mock_conn.request.call_count, 2)
    self.assertEqual(mock_conn.close.call_count, 2)
    provider.close()

  def test_token_reads_k8s_sa_token_file(self):
    expiration_time = int(time.time() + 3600)
    expected_token = create_test_jwt({"exp": expiration_time})
    response = _make_mock_http_response(
        http.HTTPStatus.OK,
        json.dumps({"token": expected_token, "expires": expiration_time}),
    )
    mock_conn = mock.MagicMock(spec=http.client.HTTPSConnection)
    mock_conn.getresponse.return_value = response

    with tempfile.TemporaryDirectory() as tmpdir:
      token_file = os.path.join(tmpdir, "token")
      with open(token_file, "w", encoding="utf-8") as f:
        f.write("file-sa-token\n")

      with ipcidentity.IpcIdentity(
          http_connection_factory=lambda: mock_conn,
          token_path=token_file,
      ) as provider:
        self.assertEqual(provider.token(), expected_token)

    mock_conn.request.assert_called_once_with(
        "POST",
        ipcidentity.XFA_CONFIG_PATH,
        body="{}",
        headers={
            "Authorization": "Bearer file-sa-token",
            "Content-Type": "application/json",
        },
    )

  def test_token_retrieval_missing_sa_token_file(self):
    provider = ipcidentity.IpcIdentity(
        http_connection_factory=mock.MagicMock(),
        token_path="/nonexistent/serviceaccount/token",
    )
    with self.assertRaisesRegex(RuntimeError, "Failed to read bearer token"):
      provider.token()

  def test_token_retrieval_bearer_source_error(self):
    provider = ipcidentity.IpcIdentity(
        http_connection_factory=mock.MagicMock(),
        bearer_token_source=mock.MagicMock(
            side_effect=ValueError("SA token error")
        ),
    )
    with self.assertRaisesRegex(RuntimeError, "SA token error"):
      provider.token()

  def test_token_retrieval_http_status_error(self):
    response = _make_mock_http_response(
        http.HTTPStatus.FORBIDDEN, "rbac denied"
    )
    provider, mock_conn, _ = _setup_mock_provider([response])

    with self.assertRaisesRegex(RuntimeError, "status code 403"):
      provider.token()
    mock_conn.close.assert_called_once()

  def test_token_retrieval_invalid_json(self):
    response = _make_mock_http_response(http.HTTPStatus.OK, "not-json")
    provider, _, _ = _setup_mock_provider([response])

    with self.assertRaisesRegex(RuntimeError, "invalid JSON"):
      provider.token()

  def test_token_retrieval_missing_token_field(self):
    response = _make_mock_http_response(
        http.HTTPStatus.OK, json.dumps({"expires": int(time.time() + 3600)})
    )
    provider, _, _ = _setup_mock_provider([response])

    with self.assertRaisesRegex(RuntimeError, "valid 'token' field"):
      provider.token()

  def test_token_retrieval_missing_expires_field(self):
    response = _make_mock_http_response(
        http.HTTPStatus.OK, json.dumps({"token": "abc"})
    )
    provider, _, _ = _setup_mock_provider([response])

    with self.assertRaisesRegex(RuntimeError, "valid 'expires' field"):
      provider.token()

  def test_metadata(self):
    expiration_time = int(time.time() + 3600)
    expected_token = create_test_jwt({
        "exp": expiration_time,
        "user_id": "test_user@ipc.intrinsic.ai",
        "orgid": "test_org",
    })
    response = _make_mock_http_response(
        http.HTTPStatus.OK,
        json.dumps({"token": expected_token, "expires": expiration_time}),
    )
    provider, _, _ = _setup_mock_provider([response])

    metadata_list = provider.metadata()

    metadata_dict = dict(metadata_list)
    self.assertIn("cookie", metadata_dict)
    cookies = [v for k, v in metadata_list if k == "cookie"]
    self.assertTrue(any("org-id=test_org" in c for c in cookies))
    self.assertTrue(any("auth-proxy=" + expected_token in c for c in cookies))
    self.assertEqual(metadata_dict["x-intrinsic-org"], "test_org")
    provider.close()

  def test_metadata_missing_orgid(self):
    expiration_time = int(time.time() + 3600)
    expected_token = create_test_jwt({
        "exp": expiration_time,
        "user_id": "test_user@ipc.intrinsic.ai",
    })
    response = _make_mock_http_response(
        http.HTTPStatus.OK,
        json.dumps({"token": expected_token, "expires": expiration_time}),
    )
    provider, _, _ = _setup_mock_provider([response])

    with self.assertRaisesRegex(
        RuntimeError, "does not contain an 'orgid' claim"
    ):
      provider.metadata()
    provider.close()


if __name__ == "__main__":
  absltest.main()
