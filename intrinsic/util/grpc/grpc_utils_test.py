# Copyright 2023 Intrinsic Innovation LLC

"""Tests for the grpc_utils module."""

from unittest import mock

from absl.testing import absltest
import grpc
from intrinsic.kubernetes.acl.ipcidentity import ipcidentity
from intrinsic.util.grpc import grpc_utils


class GrpcUtilsTest(absltest.TestCase):

  @mock.patch.object(grpc_utils, '_get_environment', autospec=True)
  @mock.patch.object(grpc_utils, '_create_channel', autospec=True)
  @mock.patch.object(grpc_utils, '_add_auth_header', autospec=True)
  def test_create_assets_channel(
      self,
      mock_add_auth_header,
      mock_create_channel,
      mock_get_environment,
  ):
    """Tests that create_assets_channel creates a channel with an IPC identity."""
    mock_get_environment.return_value = 'prod'
    mock_create_channel.return_value = mock.MagicMock()
    mock_add_auth_header.return_value = mock.MagicMock()

    channel = grpc_utils.create_assets_channel()

    mock_get_environment.assert_called_once()
    mock_create_channel.assert_called_once_with('assets.intrinsic.ai:443')
    mock_add_auth_header.assert_called_once_with(
        mock_create_channel.return_value, mock.ANY
    )
    self.assertEqual(channel, mock_add_auth_header.return_value)

  @mock.patch.object(grpc_utils, '_get_compute_project', autospec=True)
  @mock.patch.object(grpc_utils, '_create_channel', autospec=True)
  @mock.patch.object(grpc_utils, '_add_auth_header', autospec=True)
  def test_create_cloud_channel(
      self,
      mock_add_auth_header,
      mock_create_channel,
      mock_get_compute_project,
  ):
    """Tests that create_cloud_channel creates a channel with an IPC identity."""
    mock_get_compute_project.return_value = 'test-project'
    mock_create_channel.return_value = mock.MagicMock()
    mock_add_auth_header.return_value = mock.MagicMock()

    channel = grpc_utils.create_cloud_channel()

    mock_get_compute_project.assert_called_once()
    mock_create_channel.assert_called_once_with(
        'www.endpoints.test-project.cloud.goog:443'
    )
    mock_add_auth_header.assert_called_once_with(
        mock_create_channel.return_value, mock.ANY
    )
    self.assertEqual(channel, mock_add_auth_header.return_value)

  @mock.patch.object(grpc, 'ssl_channel_credentials', autospec=True)
  @mock.patch.object(grpc, 'secure_channel', autospec=True)
  def test_create_channel(
      self, mock_secure_channel, mock_ssl_channel_credentials
  ):
    """Tests that _create_channel creates a secure channel with the correct options."""
    mock_secure_channel.return_value = mock.MagicMock()
    mock_ssl_channel_credentials.return_value = mock.MagicMock()

    channel = grpc_utils._create_channel('test-address')

    mock_secure_channel.assert_called_once_with(
        'test-address',
        mock_ssl_channel_credentials.return_value,
        options=[('grpc.max_receive_message_length', -1)],
    )
    self.assertEqual(channel, mock_secure_channel.return_value)

  @mock.patch.object(grpc_utils, '_get_environment', autospec=True)
  @mock.patch.object(grpc_utils, '_create_channel', autospec=True)
  @mock.patch.object(grpc_utils, '_add_auth_header', autospec=True)
  def test_create_assets_channel_calls_add_auth_header(
      self,
      mock_add_auth_header,
      mock_create_channel,
      mock_get_environment,
  ):
    """Tests that create_assets_channel calls _add_auth_header."""
    mock_get_environment.return_value = 'prod'
    mock_create_channel.return_value = mock.MagicMock()
    mock_add_auth_header.return_value = mock.MagicMock()

    grpc_utils.create_assets_channel()

    mock_add_auth_header.assert_called_once()

  @mock.patch.object(grpc_utils, '_get_compute_project', autospec=True)
  @mock.patch.object(grpc_utils, '_create_channel', autospec=True)
  @mock.patch.object(grpc_utils, '_add_auth_header', autospec=True)
  def test_create_cloud_channel_calls_add_auth_header(
      self,
      mock_add_auth_header,
      mock_create_channel,
      mock_get_compute_project,
  ):
    """Tests that create_cloud_channel calls _add_auth_header."""
    mock_get_compute_project.return_value = 'test-project'
    mock_create_channel.return_value = mock.MagicMock()
    mock_add_auth_header.return_value = mock.MagicMock()

    grpc_utils.create_cloud_channel()

    mock_add_auth_header.assert_called_once()

  @mock.patch.object(grpc, 'intercept_channel', autospec=True)
  @mock.patch.object(grpc, 'Channel', autospec=True)
  @mock.patch.object(ipcidentity, 'IpcIdentity', autospec=True)
  def test_add_auth_header(
      self, mock_ipc_identity, mock_channel, mock_intercept_channel
  ):
    """Tests that _add_auth_header adds the auth header to the channel."""
    mock_ipc_identity.token.return_value = 'test_token'
    mock_intercept_channel.return_value = mock_channel.return_value

    # Call _add_auth_header
    intercepted_channel = grpc_utils._add_auth_header(
        mock_channel.return_value, mock_ipc_identity
    )

    # Assert that the channel was intercepted
    self.assertEqual(intercepted_channel, mock_channel.return_value)
    mock_intercept_channel.assert_called_once()

    # Get the interceptor from the call
    interceptor = mock_intercept_channel.call_args[0][1]

    # Mock the continuation function to capture the ClientCallDetails
    mock_continuation = mock.Mock()
    mock_continuation.return_value = 'test_response'

    # Create a dummy ClientCallDetails
    client_call_details = grpc_utils.interceptor.ClientCallDetails(
        method='test_method',
        timeout=None,
        metadata=None,
        credentials=None,
        wait_for_ready=None,
    )

    # Call the captured interceptor
    response = interceptor.intercept_call(
        mock_continuation, client_call_details, 'test_request'
    )

    # Assert that the continuation was called
    mock_continuation.assert_called_once()
    self.assertEqual(response, 'test_response')

    # Assert that the ClientCallDetails was modified with the correct metadata
    modified_call_details = mock_continuation.call_args[0][0]
    self.assertIsNotNone(modified_call_details.metadata)
    self.assertIn(
        ('cookie', 'auth-proxy=test_token'), modified_call_details.metadata
    )


if __name__ == '__main__':
  absltest.main()
