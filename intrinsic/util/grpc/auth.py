# Copyright 2023 Intrinsic Innovation LLC

"""Provides methods to handle API keys for authentication.

This file implements a subset of the
`//intrinsic/tools/inctl/auth/auth.go` authorization library.
The implementation only contains methods to read the keys in Python.
Login etc. is handled by the inctl CLI.
"""

import dataclasses
import datetime
import json
import logging
import os.path
import threading
from typing import Dict
from typing import Optional
from typing import Tuple
import urllib.error
import urllib.request

import grpc
import retrying

from intrinsic.kubernetes.acl.py import jwt
from intrinsic.util.grpc import userconfig
ALIAS_DEFAULT_TOKEN = "default"
AUTH_CONFIG_EXTENSION = ".user-token"
STORE_DIRECTORY = "intrinsic/projects"
ORG_STORE_DIRECTORY = "intrinsic/organizations"

_MIN_TOKEN_LIFETIME = datetime.timedelta(minutes=1)


class APIKeyTokenSource(grpc.AuthMetadataPlugin):
  """Provides a JWT token retrieved using an API key.

  Can be used as a gRPC AuthMetadataPlugin.
  """

  def __init__(
      self,
      api_key: str,
      portal_domain: str,
  ):
    super().__init__()
    self._api_key = api_key
    self._portal_domain = portal_domain

    self._lock = threading.Lock()
    self._cached_token: str | None = None
    self._cached_expiry: datetime.datetime | None = None

  # Retry even "non-transient" errors like 401 as these can occur transiently.
  # http://b/504966163
  @retrying.retry(
      stop_max_attempt_number=4,
      wait_exponential_multiplier=1000,
      wait_jitter_max=1000,
      retry_on_exception=lambda e: isinstance(e, urllib.error.URLError),
  )
  def _exchange_token(self) -> str:
    url = f"https://{self._portal_domain}/api/v1/accountstokens:idtoken"
    data = json.dumps({"api_key": self._api_key, "do_fan_out": True}).encode(
        "utf-8"
    )

    headers = {
        "Content-Type": "application/json",
        "User-Agent": "sbl/1.0 (py, urllib)",
    }

    req = urllib.request.Request(url, data=data, headers=headers, method="POST")

    try:
      with urllib.request.urlopen(req, timeout=10) as response:
        res_data = json.loads(response.read().decode("utf-8"))
        return res_data["idToken"]
    except Exception as e:
      logging.warning("Failed to call accounts service: %s", e)
      raise

  def _get_now_utc(self) -> datetime.datetime:
    # Allows mocking the current time in unit tests, avoiding:
    #  TypeError: cannot set 'now' attribute of immutable type 'datetime.datetime'
    return datetime.datetime.now(datetime.timezone.utc)

  def token(self) -> str:
    """Returns the cached JWT, fetching it from the portal if expired or missing."""
    now = self._get_now_utc()
    with self._lock:
      if (
          self._cached_token is None
          or self._cached_expiry is None
          or self._cached_expiry - _MIN_TOKEN_LIFETIME < now
      ):
        token = self._exchange_token()
        self._cached_token = token
        self._cached_expiry = jwt.ExpiresAt(token)
      return self._cached_token

  def __call__(self, context, callback):
    try:
      token = self.token()
      callback((("cookie", f"auth-proxy={token}"),), None)
    except Exception as e:
      # gRPC wants the exception in the callback:
      # https://grpc.github.io/grpc/python/grpc.html#grpc.AuthMetadataPluginCallback
      callback(None, e)


@dataclasses.dataclass()
class ProjectToken:
  """Contains an API key and corresponding helpers."""

  api_key: str
  valid_until: Optional[datetime.datetime] = None

  def validate(self) -> None:
    if self.valid_until is not None:
      if datetime.datetime.now() > self.valid_until:
        raise AttributeError(f"project token expired: {self.valid_until}")

  def get_request_metadata(self) -> Tuple[Tuple[str, str], ...]:
    self.validate()
    return (("authorization", "Bearer " + self.api_key),)

  def as_id_token_credentials(
      self,
      portal_domain: str = "flowstate.intrinsic.ai",
  ) -> APIKeyTokenSource:
    return APIKeyTokenSource(
        api_key=self.api_key,
        portal_domain=portal_domain,
    )


@dataclasses.dataclass()
class ProjectConfiguration:
  """Contains a list of API keys for a given project."""

  name: str
  tokens: Dict[str, ProjectToken]

  def has_credentials(self, alias: str) -> bool:
    return alias in self.tokens

  def get_credentials(self, alias: str) -> ProjectToken:
    if not self.has_credentials(alias):
      raise KeyError(f"token with alias '{alias}' not found")

    return self.tokens[alias]

  def get_default_credentials(self) -> ProjectToken:
    return self.get_credentials(ALIAS_DEFAULT_TOKEN)


class CredentialsNotFoundError(ValueError):
  """Thrown in case the lookup for a given credential name failed.

  Attributes:
    message: the error message
    project_name: GCP project name for which the credentials were not found.
  """

  def __init__(self, message: str, project_name: str) -> None:
    super().__init__(message)
    self.project_name = project_name

  def __str__(self) -> str:
    return f"Credentials for project '{self.project_name}' could not be found!"


def get_configuration(name: str) -> ProjectConfiguration:
  """Reads the local project configuration for the provided project name.

  Args:
    name: name of the GCP project

  Raises:
    CredentialsNotFoundError: if configuration for the project could not be
    found.

  Returns:
    configuration for the project.
  """
  file_name = os.path.join(
      userconfig.get_user_config_dir(),
      STORE_DIRECTORY,
      name + AUTH_CONFIG_EXTENSION,
  )

  try:
    with open(file_name, "r") as f:
      config = json.load(f)
      tokens = {}
      for alias, token in config["tokens"].items():
        tokens[alias] = ProjectToken(
            api_key=token["apiKey"],
            valid_until=(
                datetime.datetime.fromisoformat(token["validUntil"])
                if "validUntil" in token
                else None
            ),
        )
      return ProjectConfiguration(name=name, tokens=tokens)
  except FileNotFoundError as e:
    raise CredentialsNotFoundError(message=e.strerror, project_name=name) from e


@dataclasses.dataclass()
class OrgInfo:
  """Encapsulates the information needed to access an organization."""

  organization: str
  project: str


class OrgNotFoundError(ValueError):
  """Thrown in case the lookup for a given organization name failed.

  Attributes:
    organization: Organization name for which no information was found.
  """

  def __init__(self, message: str, organization: str) -> None:
    super().__init__(message)
    self.organization = organization

  def __str__(self) -> str:
    return (
        f"Information for organization '{self.organization}' could not be"
        " found!"
    )


def parse_info_from_string(info_string: str) -> OrgInfo:
  """Parses the org and project information from a string.

  Args:
    info_string: String in the format of <org>@<project>

  Raises:
    ValueError: if either org or project can not be extracted from the input.

  Returns:
    The information for the organization.
  """
  org_and_project = info_string.split("@")
  if len(org_and_project) != 2:
    raise ValueError(f"Invalid org or project information: {info_string}")
  org, project = org_and_project
  if not org or not project:
    raise ValueError(f"Invalid org or project information: {info_string}")
  return OrgInfo(organization=org, project=project)


def read_org_info(organization: str) -> OrgInfo:
  """Reads the local org information for the provided organization name.

  Args:
    organization: name of the organization

  Raises:
    OrgNotFoundError: if information for the organization could not be found.

  Returns:
    The information for the organization.
  """
  file_name = os.path.join(
      userconfig.get_user_config_dir(),
      ORG_STORE_DIRECTORY,
      organization + ".json",
  )

  try:
    with open(file_name, "r") as f:
      org_info = json.load(f)
      return OrgInfo(organization=org_info["org"], project=org_info["project"])
  except FileNotFoundError as e:
    raise OrgNotFoundError(message=e.strerror, organization=organization) from e
