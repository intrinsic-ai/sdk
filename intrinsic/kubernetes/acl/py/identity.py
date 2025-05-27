# Copyright 2023 Intrinsic Innovation LLC

"""Helpers to work with user identities inside the Intrinsic stack."""

import http.cookies
import re

from absl import logging
import grpc
from intrinsic.kubernetes.acl.py import jwt

COOKIE_KEY = 'cookie'
AUTH_PROXY_COOKIE_NAME = 'auth-proxy'
PORTAL_TOKEN_COOKIE_NAME = 'portal-token'
ONPREM_TOKEN_COOKIE_NAME = 'onprem-token'
APIKEY_TOKEN_HEADER_NAME = 'apikey-token'
ORG_ID_COOKIE = 'org-id'


class Organization:
  """Represents an organization inside the Intrinsic stack."""

  def __init__(self, org_id: str):
    self.org_id = org_id


class User:
  """Represents a user inside the Intrinsic stack.

  Attributes:
    jwt (str): A json webtoken.
  """

  jwt: str

  def __init__(self, j: str) -> None:
    self.jwt = j

  def Email(self) -> str:
    """Retrieves the user or service email of an identity.

    Deprecated: Use EmailCanonicalized instead or better EmailRaw.

    Returns:
      str: The canonicalized email of the user.
    """
    return self.EmailCanonicalized()

  def EmailRaw(self) -> str:
    """Retrieves the user or service email of an identity as stored in the JWT."""
    return jwt.Email(self.jwt)

  def EmailCanonicalized(self) -> str:
    """Retrieves canonicalized user or service email of an identity.

    Use only for ACL lookups. For other use cases prefer EmailRaw.

    Returns:
      str: The canonicalized email of the user.
    """
    return CanonicalizeEmail(jwt.Email(self.jwt))


def UserFromContext(context: grpc.ServicerContext) -> User:
  """Get user identity from grpc context.

  Args:
    context: The grpc context.

  Returns:
     User: New user object.

  Raises:
    KeyError: If no jwt found.
  """
  metadata = {c.key: c.value for c in context.invocation_metadata()}
  for cn in (AUTH_PROXY_COOKIE_NAME, APIKEY_TOKEN_HEADER_NAME):
    if cn in metadata:
      if cn == AUTH_PROXY_COOKIE_NAME:
        logging.warning('Deprecated metadata key auth-proxy.')
      return User(j=metadata[cn])

  if COOKIE_KEY in metadata:
    cks = http.cookies.SimpleCookie()
    cks.load(str(metadata[COOKIE_KEY]))
    for cookie in cks.values():
      if cookie.key in {AUTH_PROXY_COOKIE_NAME, PORTAL_TOKEN_COOKIE_NAME}:
        return User(j=cookie.value)

  raise KeyError('no jwt found')


def OrgFromContext(context: grpc.ServicerContext) -> Organization:
  """Get organization from grpc context.

  Args:
    context: The grpc context.

  Returns:
     Organization: A new organization.

  Raises:
    KeyError: If no org-id found.
  """
  metadata = {c.key: c.value for c in context.invocation_metadata()}
  if COOKIE_KEY in metadata:
    cks = http.cookies.SimpleCookie()
    cks.load(str(metadata[COOKIE_KEY]))
    for cookie in cks.values():
      if cookie.key == ORG_ID_COOKIE:
        return Organization(cookie.value)
  if ORG_ID_COOKIE in metadata:
    return Organization(org_id=metadata[ORG_ID_COOKIE])

  logging.error('No organization information in context.')
  raise KeyError('no org-id found')


def CanonicalizeEmail(email: str) -> str:
  """Ensures that different valid forms of emails map to the same user account.

  Args:
    email: Any email.

  Returns:
    str: The canonicalized email.

  Raises:
    ValueError: If input is not a well formed email.
  """
  parts = email.lower().split('@', 2)
  if len(parts) != 2:
    raise ValueError('Missing "@" in email "%s"' % email)
  user, provider = parts
  if not user:
    raise ValueError('Missing user part in email "%s"' % email)
  if not provider:
    raise ValueError('Missing provider part in email "%s"' % email)

  # First canonicalize the provider part.
  if provider == 'googlemail.com':
    provider = 'gmail.com'

  # Next canonicalize the user part.
  # Cut everything starting with '+' on the part before the @ (including the
  # '+') (RFC 5233).
  user = user.split('+', 2)[0]

  # Finally canonicalize user based on provider.
  if provider == 'gmail.com':
    user = re.sub(r'[^a-zA-Z0-9]', '', user)
  else:
    user = re.sub(r'[^a-zA-Z0-9!#$%&\'*+\-/=?^_{|}~`]', '', user)

  return user + '@' + provider
