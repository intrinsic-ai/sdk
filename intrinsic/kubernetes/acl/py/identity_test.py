# Copyright 2023 Intrinsic Innovation LLC

import collections
import unittest.mock

from absl.testing import absltest
from absl.testing import parameterized
import grpc
from intrinsic.kubernetes.acl.py import identity

# JWT with {"email": "doe@example.com"}
TOKEN = b'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZW1haWwiOiJkb2VAZXhhbXBsZS5jb20iLCJpYXQiOjE1MTYyMzkwMjJ9.qRdA3amFU5P4jl4LvErW8876QAfRXryMfI9LSiLVlS8'


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
    ctx = TestContext(((ckey, cvalue),))
    u = identity.UserFromContext(ctx)
    self.assertIsNotNone(u)


class OrgTest(absltest.TestCase):
  @unittest.mock.patch.multiple(TestContext, __abstractmethods__=set())
  def test_from_context(self):
    organization_name = 'my-organization'
    ctx = TestContext(((identity.ORG_ID_COOKIE, organization_name),))
    organization = identity.OrgFromContext(ctx)
    self.assertEqual(organization.org_id, organization_name)

  @unittest.mock.patch.multiple(TestContext, __abstractmethods__=set())
  def test_from_context_cookie_field(self):
    organization_name = 'my-organization'
    ctx = TestContext((
        (identity.COOKIE_KEY, f'{identity.ORG_ID_COOKIE}={organization_name}'),
    ))
    organization = identity.OrgFromContext(ctx)
    self.assertEqual(organization.org_id, organization_name)


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
