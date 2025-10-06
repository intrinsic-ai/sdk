# Copyright 2023 Intrinsic Innovation LLC

from unittest import mock
from absl.testing import absltest
import grpc
from intrinsic.solutions import auth
from intrinsic.solutions import dialerutil


class DialerutilTest(absltest.TestCase):

  @mock.patch.object(dialerutil, "_create_channel")
  @mock.patch.object(dialerutil, "_get_cluster_from_solution")
  def test_create_channel_from_org_and_solution(
      self,
      mock_get_cluster_from_solution: mock.MagicMock,
      mock_create_channel: mock.MagicMock,
  ):
    mock_get_cluster_from_solution.return_value = "test-cluster"
    mock_create_channel.return_value = grpc.insecure_channel("localhost:1234")

    dialerutil.create_channel_from_solution(
        org_info=auth.OrgInfo(organization="test-org", project="test-project"),
        solution="test-solution",
    )

    mock_get_cluster_from_solution.assert_called_with(
        auth.OrgInfo(organization="test-org", project="test-project"),
        "test-solution",
        None,
    )
    mock_create_channel.assert_called_with(
        org_info=auth.OrgInfo(organization="test-org", project="test-project"),
        cluster="test-cluster",
        grpc_options=None,
    )

  @mock.patch.object(auth, "get_configuration", autospec=True)
  def test_dial_channel_opens_grpc_connection(self, mock_get_configuration):
    mock_get_configuration.return_value = auth.ProjectConfiguration(
        name="test-project",
        tokens={"default": auth.ProjectToken("test-token", None)},
    )
    channel = dialerutil._create_channel(
        org_info=auth.OrgInfo(organization="test-org", project="test-project"),
        cluster="test-cluster",
    )
    self.assertIsInstance(channel, grpc.Channel)

  @mock.patch.object(grpc, "secure_channel", autospec=True)
  @mock.patch.object(grpc, "metadata_call_credentials", autospec=True)
  @mock.patch.object(auth, "parse_info_from_string", autospec=True)
  def test_create_channel_from_token(
      self,
      mock_parse_info_from_string: mock.MagicMock,
      mock_metadata_call_credentials: mock.MagicMock,
      mock_secure_channel: mock.MagicMock,
  ):
    mock_parse_info_from_string.return_value = auth.OrgInfo(
        organization="test-org", project="test-project"
    )
    mock_metadata_call_credentials.return_value = mock.MagicMock()

    dialerutil.create_channel_from_token(
        auth_token="test-auth-token",
        org="test-org",
        cluster="test-cluster",
    )

    mock_parse_info_from_string.assert_called_with("test-org")
    self.assertTrue(
        any(
            isinstance(c.args[0], dialerutil._AuthProxy)
            for c in mock_metadata_call_credentials.call_args_list
        ),
        "grpc.metadata_call_credentials was not called with _AuthProxy",
    )
    self.assertTrue(mock_secure_channel.called)


if __name__ == "__main__":
  absltest.main()
