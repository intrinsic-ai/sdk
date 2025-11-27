# Copyright 2023 Intrinsic Innovation LLC

import json
import time
from typing import Any
from unittest import mock

from absl.testing import absltest

from intrinsic.kubernetes.accounts.service.api.tokens.v2 import tokens_pb2
from intrinsic.kubernetes.acl.ipcidentity import ipcidentity
from intrinsic.kubernetes.acl.ipcidentity.internal import metadata


# Helper to create a simple JWT string for testing
def create_test_jwt(payload: dict[str, Any]) -> str:
  """Creates a simple JWT string for testing.

  Args:
    payload: The payload dictionary to include in the JWT.

  In a real test, you might not need a fully valid JWT structure,
  just something that jwt.PayloadUnsafe can parse.
  For simplicity, we'll just use the payload part for the mock response.
  The actual token returned by AccountsTokensService would be a full JWT.
  Here, we're focusing on what IpcIdentity expects *after* decoding.
  The important part is that jwt.PayloadUnsafe(token_str) works.
  A minimal JWT is "header.payload.signature".
  We'll make a dummy header and signature for the structure.

  Returns:
    A JWT string with a dummy header and signature.
  """
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


def setup_mock_identity_provider(
    robot_token: str = "test_robot_jwt_string",
) -> tuple[ipcidentity.IpcIdentity, mock.MagicMock, mock.MagicMock]:
  """Sets up a mock IpcIdentity provider for testing.

  Args:
    robot_token: The robot token to return from the metadata client.

  Returns:
    A tuple containing the IpcIdentity instance, the mock tokens stub,
    and the mock metadata client.
  """
  mock_metadata_client = mock.create_autospec(
      metadata.MetadataClient, instance=True
  )
  mock_metadata_client.token.return_value = robot_token
  # Mocking compute_project for _ensure_grpc_setup if it were to run fully
  mock_metadata_client.compute_project.return_value = "test-project"

  identity_provider = ipcidentity.IpcIdentity(
      metadata_client=mock_metadata_client
  )

  # Create a mock for the AccountsTokensServiceStub
  mock_tokens_stub = mock.MagicMock()

  # Directly assign the mock stub and a mock channel to the instance
  identity_provider._tokens_stub = mock_tokens_stub  # pylint: disable=protected-access
  identity_provider._grpc_channel = mock.MagicMock()  # pylint: disable=protected-access

  return identity_provider, mock_tokens_stub, mock_metadata_client


class IpcidentityTest(absltest.TestCase):
  """Tests for ipcidentity."""

  def test_token_retrieval_success(self):
    """Tests successful token retrieval and caching."""
    identity_provider, mock_tokens_stub, mock_metadata_client = (
        setup_mock_identity_provider()
    )

    # Prepare the mock response from GetIPCToken
    expiration_time = time.time() + 3600  # Expires in 1 hour
    expected_ipc_token_payload = {
        "exp": int(expiration_time),
        "user_id": "test_user@example.com",
        # Add other claims as needed by your tests or code
    }
    # The actual token string that GetIPCToken would return
    expected_ipc_token_str = create_test_jwt(expected_ipc_token_payload)

    mock_response = tokens_pb2.GetIPCTokenResponse(
        ipc_token=expected_ipc_token_str
    )
    mock_tokens_stub.GetIPCToken.return_value = mock_response

    # Call the method under test
    actual_token = identity_provider.token()

    # Assertions
    self.assertEqual(actual_token, expected_ipc_token_str)
    mock_metadata_client.token.assert_called_once()
    mock_tokens_stub.GetIPCToken.assert_called_once()

    # Check the request passed to GetIPCToken
    expected_request = tokens_pb2.GetIPCTokenRequest(
        credential=tokens_pb2.IPCCredential(
            robot_jwt=tokens_pb2.RobotJWT(jwt="test_robot_jwt_string")
        )
    )
    mock_tokens_stub.GetIPCToken.assert_called_with(
        expected_request, timeout=10
    )

    # Test caching: calling token() again should not call GetIPCToken
    mock_tokens_stub.GetIPCToken.reset_mock()  # Reset call count
    cached_token = identity_provider.token()
    self.assertEqual(cached_token, expected_ipc_token_str)
    mock_tokens_stub.GetIPCToken.assert_not_called()

    # Close to avoid resource warnings when the test runner cares
    identity_provider.close()

  def test_token_retrieval_expired(self):
    """Tests token refresh when the cached token is expired."""
    identity_provider, mock_tokens_stub, _ = setup_mock_identity_provider()

    # Prepare the mock response from GetIPCToken with an expired token
    expired_time = time.time() - 3600  # Expired 1 hour ago
    expired_ipc_token_payload = {
        "exp": int(expired_time),
        "user_id": "test_user@example.com",
    }
    expired_ipc_token_str = create_test_jwt(expired_ipc_token_payload)

    # Prepare a second response with a valid token
    valid_expiration_time = time.time() + 3600  # Expires in 1 hour
    valid_ipc_token_payload = {
        "exp": int(valid_expiration_time),
        "user_id": "test_user@example.com",
    }
    valid_ipc_token_str = create_test_jwt(valid_ipc_token_payload)

    mock_response_expired = tokens_pb2.GetIPCTokenResponse(
        ipc_token=expired_ipc_token_str
    )
    mock_response_valid = tokens_pb2.GetIPCTokenResponse(
        ipc_token=valid_ipc_token_str
    )
    mock_tokens_stub.GetIPCToken.side_effect = [
        mock_response_expired,
        mock_response_valid,
    ]

    # Call the method under test for the first time (expired token)
    first_token = identity_provider.token()
    self.assertEqual(first_token, expired_ipc_token_str)
    mock_tokens_stub.GetIPCToken.assert_called_once()

    # Call the method again, it should have refreshed the token
    second_token = identity_provider.token()
    self.assertEqual(second_token, valid_ipc_token_str)
    self.assertEqual(mock_tokens_stub.GetIPCToken.call_count, 2)

    # Remember to close to avoid resource warnings if your test runner cares
    identity_provider.close()

  def test_token_retrieval_robot_token_error(self):
    """Tests handling of errors when retrieving the robot token."""
    mock_metadata_client = mock.create_autospec(
        metadata.MetadataClient, instance=True
    )
    mock_metadata_client.token.side_effect = ValueError("Robot token error")

    identity_provider = ipcidentity.IpcIdentity(
        metadata_client=mock_metadata_client
    )

    with self.assertRaisesRegex(RuntimeError, "Robot token error"):
      identity_provider.token()

    # Remember to close to avoid resource warnings if your test runner cares
    identity_provider.close()

  def test_token_retrieval_grpc_error(self):
    """Tests handling of errors when calling the gRPC service."""
    identity_provider, mock_tokens_stub, _ = setup_mock_identity_provider()
    mock_tokens_stub.GetIPCToken.side_effect = Exception("gRPC error")

    with self.assertRaisesRegex(RuntimeError, "gRPC error"):
      identity_provider.token()

    # Remember to close to avoid resource warnings if your test runner cares
    identity_provider.close()


if __name__ == "__main__":
  absltest.main()
