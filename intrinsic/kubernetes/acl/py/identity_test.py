# Copyright 2023 Intrinsic Innovation LLC

import collections
import http.cookies
import unittest.mock

from absl.testing import absltest
from absl.testing import parameterized
import grpc

from intrinsic.kubernetes.acl.py import identity

# JWT with {"email": "doe@example.com"}
TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZW1haWwiOiJkb2VAZXhhbXBsZS5jb20iLCJpYXQiOjE1MTYyMzkwMjJ9.qRdA3amFU5P4jl4LvErW8876QAfRXryMfI9LSiLVlS8'


_Metadatum = collections.namedtuple(
    '_Metadatum',
    (
        'key',
        'value',
    ),
)


class TestContext(grpc.ServicerContext):
  md: list[_Metadatum]

  def __init__(self, md: list[tuple[str, str]]) -> None:
    super(grpc.ServicerContext, self).__init__()
    self.md = (_Metadatum(key=k, value=v) for k, v in md)

  def invocation_metadata(self) -> list[_Metadatum]:
    return self.md


class IdentityTest(parameterized.TestCase):

  @parameterized.parameters(
      # Deprecated metadata["auth-proxy"] JWT.
      (identity.AUTH_PROXY_COOKIE_NAME, TOKEN),
      # metadata["apikey-token"] we still support as header for now
      (identity.APIKEY_TOKEN_HEADER_NAME, TOKEN),
      # JWT from metadata["cookie"]["auth-proxy"]
      (
          identity.COOKIE_KEY,
          f'{identity.AUTH_PROXY_COOKIE_NAME}={TOKEN};',
      ),
      # JWT from metadata["cookie"]["portal-token"]
      (
          identity.COOKIE_KEY,
          f'{identity.PORTAL_TOKEN_COOKIE_NAME}={TOKEN};',
      ),
  )
  @unittest.mock.patch.multiple(TestContext, __abstractmethods__=set())
  def test_from_context(self, ckey, cvalue):
    ctx = TestContext([(ckey, cvalue)])
    u = identity.UserFromContext(ctx)
    self.assertIsNotNone(u)

  @unittest.mock.patch.multiple(TestContext, __abstractmethods__=set())
  def test_from_context_with_cookies(self):
    ctx = TestContext(
        identity.CookiesToGRPCMetadata(
            http.cookies.SimpleCookie(
                {identity.PORTAL_TOKEN_COOKIE_NAME: TOKEN}
            )
        )
    )
    u = identity.UserFromContext(ctx)
    self.assertIsNotNone(u)

  def test_user_to_grpc_metadata(self):
    user = identity.User(j=TOKEN)
    metadata = identity.UserToGRPCMetadata(user)
    self.assertLen(metadata, 1)
    self.assertEqual(metadata[0][0], identity.COOKIE_KEY)
    self.assertIn(f'{identity.AUTH_PROXY_COOKIE_NAME}={TOKEN}', metadata[0][1])

  @unittest.mock.patch.multiple(TestContext, __abstractmethods__=set())
  def test_to_grpc_metadata_from_incoming(self):
    ctx = TestContext([
        (identity.COOKIE_KEY, f'{identity.AUTH_PROXY_COOKIE_NAME}={TOKEN};'),
        (identity.APIKEY_TOKEN_HEADER_NAME, 'apikey'),
        (identity.AUTH_HEADER_NAME, 'auth'),
        ('other', 'stuff'),
    ])
    metadata = identity.ToGRPCMetadataFromIncoming(ctx)
    self.assertCountEqual(
        metadata,
        [
            (
                identity.COOKIE_KEY,
                f'{identity.AUTH_PROXY_COOKIE_NAME}={TOKEN};',
            ),
            (identity.APIKEY_TOKEN_HEADER_NAME, 'apikey'),
            (identity.AUTH_HEADER_NAME, 'auth'),
        ],
    )


class OrgTest(absltest.TestCase):

  @unittest.mock.patch.multiple(TestContext, __abstractmethods__=set())
  def test_from_context_cookie_field(self):
    organization_name = 'my-organization'
    ctx = TestContext(identity.OrgIDToGRPCMetadata(organization_name))
    organization = identity.OrgFromContext(ctx)
    self.assertEqual(organization.org_id, organization_name)

  def test_org_name_call_credentials(self):
    self.assertEqual(
        identity._OrgName('my-organization')._organization_name,
        'my-organization',
    )
    self.assertEqual(
        identity._OrgName('my-organization@my-project')._organization_name,
        'my-organization',
    )


class CanonicalizationTest(parameterized.TestCase):

  @parameterized.parameters('', 'john', 'john@', '@gmail.com', '@')
  def test_invalid_email(self, email):
    with self.assertRaises(ValueError):
      identity.CanonicalizeEmail(email)

  @parameterized.parameters(
      ('doe@gmail.com', 'doe@gmail.com'),
      ('john.doe@gmail.com', 'johndoe@gmail.com'),
      ('.john..doe.@gmail.com', 'johndoe@gmail.com'),
      ('John.Doe@gmail.com', 'johndoe@gmail.com'),
      ('doe+foo@gmail.com', 'doe@gmail.com'),
      ('doe@googlemail.com', 'doe@gmail.com'),
      ('!john.doe#@gmail.com', 'johndoe@gmail.com'),
      ('!john.doe#@yahoo.com', '!johndoe#@yahoo.com'),
  )
  def test_email_cononicalization(self, email, want):
    got = identity.CanonicalizeEmail(email)
    self.assertEqual(got, want)


if __name__ == '__main__':
  absltest.main()
