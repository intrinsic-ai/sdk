# Copyright 2023 Intrinsic Innovation LLC

from absl.testing import absltest
from intrinsic.kubernetes.acl.py import jwt

# JWT with {"email": "doe@example.com"}
TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZW1haWwiOiJkb2VAZXhhbXBsZS5jb20iLCJpYXQiOjE1MTYyMzkwMjJ9.qRdA3amFU5P4jl4LvErW8876QAfRXryMfI9LSiLVlS8"


class JwtTest(absltest.TestCase):

  def test_payload_unsafe(self):
    got = jwt.PayloadUnsafe(TOKEN)
    self.assertNotEmpty(got)

  def test_email(self):
    got = jwt.Email(TOKEN)
    self.assertEqual(got, "doe@example.com")


if __name__ == "__main__":
  absltest.main()
