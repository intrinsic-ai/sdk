# Copyright 2023 Intrinsic Innovation LLC

"""Provides a robot JWT from the metadata server."""

import http
import http.client
import io
from typing import Optional


METADATA_HOST = "169.254.169.254"
METADATA_URL = f"http://{METADATA_HOST}/computeMetadata/v1/"
IDENTITY_URL = METADATA_URL + "instance/service-accounts/default/identity"
PROJECT_URL = METADATA_URL + "project/project-id"


class MetadataClient:
  """Provides a robot JWT from the metadata server."""

  def __init__(self, http_client: Optional[http.client.HTTPConnection] = None):
    if http_client is None:
      self.http_client = http.client.HTTPConnection(METADATA_HOST)
    else:
      self.http_client = http_client

  def token(self) -> str:
    """Returns the robot JWT from the metadata server."""
    return self._get(IDENTITY_URL)

  def compute_project(self) -> str:
    """Returns the current compute project from the metadata server."""
    return self._get(PROJECT_URL)

  def _get(self, url: str) -> str:
    """Gets a value from the metadata server."""
    self.http_client.request("GET", url)
    response = self.http_client.getresponse()
    if response.status != http.HTTPStatus.OK:
      raise RuntimeError(f"status code {response.status}")
    ret = io.BytesIO(response.read())
    return ret.read().decode("utf-8")
