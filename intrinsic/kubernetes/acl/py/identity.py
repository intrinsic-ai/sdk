# Copyright 2023 Intrinsic Innovation LLC

"""Helpers to work with user identities inside the Intrinsic stack."""

import http.cookies
import re
from typing import List
from typing import Optional

from absl import logging
import grpc

from intrinsic.kubernetes.acl.py import jwt

COOKIE_KEY = 'cookie'
AUTH_PROXY_COOKIE_NAME = 'auth-proxy'
PORTAL_TOKEN_COOKIE_NAME = 'portal-token'
ONPREM_TOKEN_COOKIE_NAME = 'onprem-token'
APIKEY_TOKEN_HEADER_NAME = 'apikey-token'
AUTH_HEADER_NAME = 'authorization'
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
      str: The canonicalized email of the user. e.g. "user@gmail.com".
    """
    return self.EmailCanonicalized()

  def EmailRaw(self) -> str:
    """Retrieves the user or service email of an identity as stored in the JWT."""
    return jwt.Email(self.jwt)

  def EmailCanonicalized(self) -> str:
    """Retrieves canonicalized user or service email of an identity.

    Use only for ACL lookups. For other use cases prefer EmailRaw.

    Returns:
      str: The canonicalized email of the user. e.g. "user@gmail.com".
    """
    return CanonicalizeEmail(jwt.Email(self.jwt))


def GetJWTFromContext(context: grpc.ServicerContext) -> Optional[str]:
  """Extracts the raw JWT token from the context.

  See GetJWTFromContext in intrinsic/kubernetes/acl/identity.go.

  Priority:
    1. auth-proxy cookie
    2. apikey-token header
    3. authorization header
  """
  metadata = dict(context.invocation_metadata())

  # 1. Check cookies (auth-proxy or portal-token)
  if COOKIE_KEY in metadata:
    cks = http.cookies.SimpleCookie()
    cks.load(str(metadata[COOKIE_KEY]))
    for key in (AUTH_PROXY_COOKIE_NAME, PORTAL_TOKEN_COOKIE_NAME):
      if key in cks:
        return cks[key].value

  # 2. Check apikey-token header
  if APIKEY_TOKEN_HEADER_NAME in metadata:
    return metadata[APIKEY_TOKEN_HEADER_NAME]

  # 3. Check authorization header
  if AUTH_HEADER_NAME in metadata:
    val = metadata[AUTH_HEADER_NAME]
    if val.lower().startswith('bearer '):
      return val[7:]
    return val

  return None


def UserFromContext(context: grpc.ServicerContext) -> User:
  """Get user identity from grpc context.

  Does not verify JWT signatures.

  Args:
    context: The grpc context.

  Returns:
     User: New user object. e.g. User(jwt='...')

  Raises:
    KeyError: If no jwt found.
  """
  token = GetJWTFromContext(context)
  if token:
    return User(j=token)
  raise KeyError('no jwt found')


def OrgFromContext(context: grpc.ServicerContext) -> Organization:
  """Get organization from grpc context.

  Args:
    context: The grpc context.

  Returns:
     Organization: A new organization. e.g. Organization('my-org')

  Raises:
    KeyError: If no org-id found.
  """
  metadata = dict(context.invocation_metadata())
  if COOKIE_KEY in metadata:
    cks = http.cookies.SimpleCookie()
    cks.load(str(metadata[COOKIE_KEY]))
    for cookie in cks.values():
      if cookie.key == ORG_ID_COOKIE:
        return Organization(cookie.value)
  if ORG_ID_COOKIE in metadata:
    logging.error("""Found org-id in metadata directly instead of a cookie.
                  Update your code to use cookies instead.""")
    raise KeyError("""Tried using org from context metadata instead of a
                    cookie. Update your code to use cookies instead.""")

  logging.error('No organization information in context.')
  raise KeyError('no org-id found')


def UserToGRPCMetadata(user: User) -> List[tuple[str, str]]:
  """Converts a user's identity to a gRPC metadata list.

  Example:
    Input: User(jwt='abc')
    Output: [('cookie', 'auth-proxy=abc')]

  Args:
    user: The user to add to the metadata.

  Returns:
    A list of (key, value) pairs containing the user's identity as a cookie.
  """
  return CookiesToGRPCMetadata(
      http.cookies.SimpleCookie({AUTH_PROXY_COOKIE_NAME: user.jwt})
  )


def ToGRPCMetadataFromIncoming(
    context: grpc.ServicerContext,
) -> List[tuple[str, str]]:
  """Copies auth-related incoming GRPC metadata to a metadata list.

  Selectively copies and normalizes auth-related incoming GRPC metadata.
  Collapses multiple JWT identity sources into a single auth-proxy cookie.
  Backfills the org-id cookie if it is missing but the org-id header is present.

  See: `RequestToContext` in intrinsic/kubernetes/acl/identity.go.
  We require the user JWT to be on metadata["Cookie"]["auth-proxy"] because
  that is where the auth proxy service reads it from.

  Example:
    Input context metadata:
      cookie: "org-id=my-org"
      apikey-token: "token-val"

    Output:
      [
        ('cookie', 'auth-proxy=token-val; org-id=my-org'),
      ]

  Args:
    context: The grpc context.

  Returns:
    A list of (key, value) pairs containing auth-related metadata.
  """
  metadata = dict(context.invocation_metadata())
  outgoing_cks = http.cookies.SimpleCookie()

  # Extract JWT token (identity) from context, to "normalize it" into the
  # auth-proxy cookie like in identity.go.
  token = GetJWTFromContext(context)
  if token:
    outgoing_cks[AUTH_PROXY_COOKIE_NAME] = token

  # Selectively copy org-id cookie
  if COOKIE_KEY in metadata:
    incoming_cks = http.cookies.SimpleCookie()
    incoming_cks.load(str(metadata[COOKIE_KEY]))
    if ORG_ID_COOKIE in incoming_cks:
      outgoing_cks[ORG_ID_COOKIE] = incoming_cks[ORG_ID_COOKIE].value

  # Backfill org-id from metadata if not present in cookies
  if ORG_ID_COOKIE not in outgoing_cks and ORG_ID_COOKIE in metadata:
    val = metadata[ORG_ID_COOKIE]
    if isinstance(val, (bytes, bytearray, memoryview)):
      val = bytes(val).decode('utf-8')
    outgoing_cks[ORG_ID_COOKIE] = str(val)

  result = []
  if outgoing_cks:
    result.extend(CookiesToGRPCMetadata(outgoing_cks))

  return result


def CookiesToGRPCMetadata(
    cookies: http.cookies.BaseCookie,
) -> List[tuple[str, str]]:
  """Converts cookies to a GRPC metadata entry.

  Example:
    Input: SimpleCookie({'key': 'val'})
    Output: [('cookie', 'key=val')]

  Args:
    cookies: The cookies to convert.

  Returns:
    A tuple of (key, value) pairs.
  """
  return [
      (COOKIE_KEY, '; '.join([f'{k}={v.value}' for k, v in cookies.items()]))
  ]


def OrgIDToGRPCMetadata(org_id: str) -> List[tuple[str, str]]:
  """Writes an org-id to a GRPC metadata entry.

  Example:
    Input: 'my-org'
    Output: [('cookie', 'org-id=my-org')]

  Args:
    org_id: The org-id to convert.

  Returns:
    A tuple of (key, value) pairs containing the org-id as a cookie.
  """
  return CookiesToGRPCMetadata(
      http.cookies.SimpleCookie({ORG_ID_COOKIE: org_id})
  )


class _OrgName(grpc.AuthMetadataPlugin):
  """gRPC Metadata Plugin that adds the org name to the header."""

  def __init__(self, org: str):
    self._organization_name = org.split('@')[0]

  def __call__(self, context, callback):
    callback(tuple(OrgIDToGRPCMetadata(self._organization_name)), None)


def OrgNameCallCredentials(org: str) -> grpc.CallCredentials:
  """Returns call credentials for the org name."""
  return grpc.metadata_call_credentials(_OrgName(org), name='OrgName')


def CanonicalizeEmail(email: str) -> str:
  """Ensures that different valid forms of emails map to the same user account.

  Example:
    Input: "User+tag@GoogleMail.com"
    Output: "user@gmail.com"

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
