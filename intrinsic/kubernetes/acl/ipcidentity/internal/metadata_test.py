# Copyright 2023 Intrinsic Innovation LLC

"""Tests for the metadata client."""

import http
import http.client
from unittest import mock

from absl.testing import absltest
from intrinsic.kubernetes.acl.ipcidentity.internal import metadata


class MetadataTest(absltest.TestCase):

  def test_get_token(self):
    mock_http_client = mock.create_autospec(http.client.HTTPConnection)
    mock_response = mock.create_autospec(http.client.HTTPResponse)
    mock_response.status = http.HTTPStatus.OK
    mock_response.read.return_value = b"test-token"
    mock_http_client.getresponse.return_value = mock_response

    client = metadata.MetadataClient(mock_http_client)
    token = client.token()
    self.assertEqual(token, "test-token")

    mock_http_client.request.assert_called_once_with(
        "GET",
        metadata.IDENTITY_URL,
    )

  def test_get_project(self):
    mock_http_client = mock.create_autospec(http.client.HTTPConnection)
    mock_response = mock.create_autospec(http.client.HTTPResponse)
    mock_response.status = 200
    mock_response.read.return_value = b"test-project"
    mock_http_client.getresponse.return_value = mock_response

    client = metadata.MetadataClient(mock_http_client)
    project = client.compute_project()
    self.assertEqual(project, "test-project")

    mock_http_client.request.assert_called_once_with(
        "GET",
        metadata.PROJECT_URL,
    )

  def test_get_token_error(self):
    mock_http_client = mock.create_autospec(http.client.HTTPConnection)
    mock_response = mock.create_autospec(http.client.HTTPResponse)
    mock_response.status = 404
    mock_http_client.getresponse.return_value = mock_response

    client = metadata.MetadataClient(mock_http_client)
    with self.assertRaisesRegex(RuntimeError, "status code 404"):
      client.token()


if __name__ == "__main__":
  absltest.main()
