# Copyright 2023 Intrinsic Innovation LLC

import base64
import datetime
import json
import pathlib
import time
from unittest import mock
import urllib.error
import urllib.request

from absl.testing import absltest

from intrinsic.util.grpc import auth
from intrinsic.util.grpc import userconfig


class AuthTest(absltest.TestCase):

  @mock.patch.object(userconfig, "get_user_config_dir", autospec=True)
  def test_get_configuration_configuration_not_found(
      self, mock_user_config_dir
  ):
    testdir = self.create_tempdir()
    mock_user_config_dir.return_value = testdir.full_path

    with self.assertRaises(auth.CredentialsNotFoundError):
      auth.get_configuration("bad_configuration")

  @mock.patch.object(userconfig, "get_user_config_dir", autospec=True)
  def test_get_configuration_returns_valid_configuration_if_exists(
      self, mock_user_config_dir
  ):
    testdir = self.create_tempdir()
    testfile = pathlib.Path(
        testdir.full_path,
        auth.STORE_DIRECTORY,
        f"test-project{auth.AUTH_CONFIG_EXTENSION}",
    )
    testfile.parent.mkdir(exist_ok=True, parents=True)
    testfile.write_text("""{
  "name": "test-project",
  "tokens": {
    "default": {
      "apiKey": "ink_0000000000000000000000000000000000000000000000000000000000000000",
      "validUntil": "2023-06-01T00:00:00.000000+00:00"
    },
    "test": {
      "apiKey": "ink_1111111111111111111111111111111111111111111111111111111111111111",
      "validUntil": "2023-07-01T00:00:00.000000+00:00"
    }
  },
  "lastUpdated": "2023-05-24T09:22:26Z"
}""")

    gold = auth.ProjectConfiguration(
        "test-project",
        tokens={
            "default": auth.ProjectToken(
                "ink_0000000000000000000000000000000000000000000000000000000000000000",
                datetime.datetime.fromisoformat(
                    "2023-06-01T00:00:00.000000+00:00"
                ),
            ),
            "test": auth.ProjectToken(
                "ink_1111111111111111111111111111111111111111111111111111111111111111",
                datetime.datetime.fromisoformat(
                    "2023-07-01T00:00:00.000000+00:00"
                ),
            ),
        },
    )
    mock_user_config_dir.return_value = testdir.full_path
    result = auth.get_configuration("test-project")
    self.assertEqual(gold, result)

  @mock.patch.object(userconfig, "get_user_config_dir", autospec=True)
  @mock.patch("intrinsic.config.environments.from_any_project")
  def test_get_configuration_returns_env_configuration_if_exists(
      self, mock_from_any_project, mock_user_config_dir
  ):
    mock_from_any_project.return_value = "test-env"
    testdir = self.create_tempdir()
    testfile = pathlib.Path(
        testdir.full_path,
        auth.ENV_STORE_DIRECTORY,
        f"test-env{auth.AUTH_CONFIG_EXTENSION}",
    )
    testfile.parent.mkdir(exist_ok=True, parents=True)
    testfile.write_text("""{
  "name": "test-env",
  "tokens": {
    "default": {
      "apiKey": "ink_env_000000000000000000000000000000000000000000000000000000000000",
      "validUntil": "2023-06-01T00:00:00.000000+00:00"
    }
  },
  "lastUpdated": "2023-05-24T09:22:26Z"
}""")

    gold = auth.ProjectConfiguration(
        "test-env",
        tokens={
            "default": auth.ProjectToken(
                "ink_env_000000000000000000000000000000000000000000000000000000000000",
                datetime.datetime.fromisoformat(
                    "2023-06-01T00:00:00.000000+00:00"
                ),
            ),
        },
    )
    mock_user_config_dir.return_value = testdir.full_path
    result = auth.get_configuration("test-project")
    self.assertEqual(gold, result)

  @mock.patch.object(auth, "get_configuration")
  def test_get_api_key_calls_get_configuration(self, mock_get_configuration):
    mock_config = mock.Mock()
    mock_creds = mock.Mock()
    mock_creds.api_key = "fake_api_key"
    mock_config.get_default_credentials.return_value = mock_creds
    mock_get_configuration.return_value = mock_config

    api_key = auth.get_api_key("my-project")

    self.assertEqual(api_key, "fake_api_key")
    mock_get_configuration.assert_called_once_with("my-project")

  def test_project_token_validate(self):
    token = auth.ProjectToken("ink_00000")

    # should not raise if API key is set (valid_until is optional)
    token.validate()

    token.valid_until = datetime.datetime.now() + datetime.timedelta(days=1)
    token.validate()

    # expired token should raise
    token.valid_until = datetime.datetime.now() - datetime.timedelta(days=1)
    with self.assertRaisesRegex(AttributeError, "project token expired: .*"):
      token.validate()

  def test_project_token_get_request_metadata(self):
    token = auth.ProjectToken("ink_00000")
    self.assertEqual(
        token.get_request_metadata(), (("authorization", "Bearer ink_00000"),)
    )

  def test_parse_info_from_string(self):
    org_info = auth.parse_info_from_string("test-org@test-project")
    self.assertEqual(
        org_info,
        auth.OrgInfo(organization="test-org", project="test-project"),
    )

  def test_parse_info_from_string_raises_on_invalid_input(self):
    with self.assertRaisesRegex(ValueError, "Invalid org or project .*"):
      auth.parse_info_from_string("test-org")

  @mock.patch.object(userconfig, "get_user_config_dir", autospec=True)
  def test_read_org_info_raises_if_file_not_found(self, mock_user_config_dir):
    testdir = self.create_tempdir()
    mock_user_config_dir.return_value = testdir.full_path

    with self.assertRaises(auth.OrgNotFoundError):
      auth.read_org_info("org_for_which_no_local_info_is_stored")

  @mock.patch.object(userconfig, "get_user_config_dir", autospec=True)
  def test_read_org_info_returns_stored_org_info(self, mock_user_config_dir):
    testdir = self.create_tempdir()
    testfile = pathlib.Path(
        testdir.full_path, auth.ORG_STORE_DIRECTORY, "my_org.json"
    )
    testfile.parent.mkdir(exist_ok=True, parents=True)
    testfile.write_text("""{
  "org": "my_org",
  "project": "my_project"
}""")
    mock_user_config_dir.return_value = testdir.full_path

    result = auth.read_org_info("my_org")

    self.assertEqual(
        result, auth.OrgInfo(organization="my_org", project="my_project")
    )


class APIKeyTokenSourceTest(absltest.TestCase):

  def setUp(self):
    super().setUp()
    self.api_key = "ink_00000"
    self.portal_domain = "flowstate.intrinsic.ai"
    self.token_source = auth.APIKeyTokenSource(
        api_key=self.api_key,
        portal_domain=self.portal_domain,
    )

  def _mint_dummy_jwt(self, exp_timestamp: int) -> str:
    header = {"alg": "none", "typ": "JWT"}
    payload = {
        "user_id": "test-user",
        "exp": exp_timestamp,
        "email_verified": True,
        "authorized": True,
    }

    header_b64 = (
        base64.urlsafe_b64encode(json.dumps(header).encode("utf-8"))
        .decode("utf-8")
        .rstrip("=")
    )
    payload_b64 = (
        base64.urlsafe_b64encode(json.dumps(payload).encode("utf-8"))
        .decode("utf-8")
        .rstrip("=")
    )
    signature_b64 = "dummy_signature"

    return f"{header_b64}.{payload_b64}.{signature_b64}"

  def _mock_jwt_response(self, jwt: str, status: int = 200) -> mock.MagicMock:
    """Helper to create a mock response for urllib.request.urlopen."""
    mock_context = mock.MagicMock()
    mock_response = mock_context.__enter__.return_value
    mock_response.status = status
    mock_response.read.return_value = json.dumps({"idToken": jwt}).encode(
        "utf-8"
    )
    return mock_context

  @mock.patch.object(urllib.request, "urlopen")
  def test_initial_exchange_and_cache(self, mock_urlopen):
    # Setup mock response
    exp_time = int(time.time()) + 3600  # expires in 1 hour
    dummy_jwt = self._mint_dummy_jwt(exp_time)

    mock_urlopen.return_value = self._mock_jwt_response(dummy_jwt)

    # First call: triggers HTTP post exchange
    t1 = self.token_source.token()
    self.assertEqual(t1, dummy_jwt)
    self.assertEqual(mock_urlopen.call_count, 1)

    # Verify HTTP post payload details
    call_arg = mock_urlopen.call_args[0][0]
    self.assertIsInstance(call_arg, urllib.request.Request)
    self.assertEqual(
        call_arg.full_url,
        f"https://{self.portal_domain}/api/v1/accountstokens:idtoken",
    )
    self.assertEqual(call_arg.get_header("Content-type"), "application/json")

    body = json.loads(call_arg.data.decode("utf-8"))
    self.assertEqual(body["api_key"], self.api_key)
    self.assertTrue(body["do_fan_out"])

    # Second call: fetches directly from cache, call count unchanged
    t2 = self.token_source.token()
    self.assertEqual(t2, dummy_jwt)
    self.assertEqual(mock_urlopen.call_count, 1)

  @mock.patch.object(urllib.request, "urlopen")
  def test_token_source_refresh_margin(self, mock_urlopen):
    # Setup mock responses
    t_start = 1700000000
    self.token_source._get_now_utc = mock.Mock(
        return_value=datetime.datetime.fromtimestamp(
            t_start, tz=datetime.timezone.utc
        )
    )
    jwt_below_min_lifetime = self._mint_dummy_jwt(t_start + 50)
    jwt_with_normal_expiry = self._mint_dummy_jwt(t_start + 3600)

    mock_resp_1 = self._mock_jwt_response(jwt_below_min_lifetime)
    mock_resp_2 = self._mock_jwt_response(jwt_with_normal_expiry)

    mock_urlopen.side_effect = [mock_resp_1, mock_resp_2]

    # 1. Fetch first token.
    t1 = self.token_source.token()
    self.assertEqual(t1, jwt_below_min_lifetime)
    self.assertEqual(mock_urlopen.call_count, 1)

    # 2. Fetch again: Should be refreshed as lifetime is so short.
    t2 = self.token_source.token()
    self.assertEqual(t2, jwt_with_normal_expiry)
    self.assertEqual(mock_urlopen.call_count, 2)

    # 3. Fetch again: Should not refresh due to ample lifetime.
    t3 = self.token_source.token()
    self.assertEqual(t3, jwt_with_normal_expiry)
    self.assertEqual(mock_urlopen.call_count, 2)

  @mock.patch.object(urllib.request, "urlopen")
  # Mock sleep for faster retries
  @mock.patch.object(time, "sleep")
  def test_exchange_exponential_backoff(self, mock_sleep, mock_urlopen):
    # Setup mock failures then success
    mock_response_err = urllib.error.HTTPError(
        url="http://test",
        code=500,
        msg="Internal Server Error",
        hdrs=None,
        fp=None,
    )

    exp_time = int(time.time()) + 3600
    dummy_jwt = self._mint_dummy_jwt(exp_time)

    mock_resp_ok = self._mock_jwt_response(dummy_jwt)

    # Set up fake behavior: Two errors followed by success
    mock_urlopen.side_effect = [
        mock_response_err,
        mock_response_err,
        mock_resp_ok,
    ]

    t = self.token_source.token()

    self.assertEqual(t, dummy_jwt)
    self.assertEqual(mock_urlopen.call_count, 3)
    self.assertEqual(mock_sleep.call_count, 2)

  @mock.patch.object(urllib.request, "urlopen")
  # Mock sleep for faster retries
  @mock.patch.object(time, "sleep")
  def test_exchange_exhausts_retries_and_raises(self, mock_sleep, mock_urlopen):
    # Set up fake behavior: Always fails
    mock_response_err = urllib.error.HTTPError(
        url="http://test",
        code=503,
        msg="Service Unavailable",
        hdrs=None,
        fp=None,
    )
    mock_urlopen.side_effect = mock_response_err

    # It should raise HTTPError after exhausting 3 retries (4 total attempts)
    with self.assertRaisesRegex(
        urllib.error.HTTPError, "HTTP Error 503: Service Unavailable"
    ):
      self.token_source.token()

    self.assertEqual(mock_urlopen.call_count, 4)
    self.assertEqual(mock_sleep.call_count, 3)

  @mock.patch.object(urllib.request, "urlopen")
  def test_token_source_grpc_plugin_metadata(self, mock_urlopen):
    exp_time = int(time.time()) + 3600
    dummy_jwt = self._mint_dummy_jwt(exp_time)

    mock_urlopen.return_value = self._mock_jwt_response(dummy_jwt)

    # Verify __call__ invokes gRPC callback with correct structure
    mock_callback = mock.Mock()
    self.token_source(context=None, callback=mock_callback)

    mock_callback.assert_called_once_with(
        (("cookie", f"auth-proxy={dummy_jwt}"),), None
    )


if __name__ == "__main__":
  absltest.main()
