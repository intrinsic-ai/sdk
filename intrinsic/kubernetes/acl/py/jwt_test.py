# Copyright 2023 Intrinsic Innovation LLC

import base64
import datetime
import json

from absl.testing import absltest
from intrinsic.kubernetes.acl.py import jwt


def _get_test_token() -> str:
  header = {"alg": "HS256", "typ": "JWT"}
  encoded_header = (
      base64.urlsafe_b64encode(json.dumps(header).encode())
      .rstrip(b"=")
      .decode()
  )
  payload = {
      "email": "doe@example.com",
  }
  encoded_payload = (
      base64.urlsafe_b64encode(json.dumps(payload).encode())
      .rstrip(b"=")
      .decode()
  )
  signature = "dummy_signature"
  encoded_signature = (
      base64.urlsafe_b64encode(signature.encode()).rstrip(b"=").decode()
  )
  return f"{encoded_header}.{encoded_payload}.{encoded_signature}"


def _get_test_token_with_expiry() -> str:
  header = {"alg": "HS256", "typ": "JWT"}
  encoded_header = (
      base64.urlsafe_b64encode(json.dumps(header).encode())
      .rstrip(b"=")
      .decode()
  )
  date = datetime.datetime(2025, 4, 1, 0, 0, 0, tzinfo=datetime.timezone.utc)
  payload = {
      "email": "doe@example.com",
      "exp": int(date.timestamp()),
  }
  encoded_payload = (
      base64.urlsafe_b64encode(json.dumps(payload).encode())
      .rstrip(b"=")
      .decode()
  )
  signature = "dummy_signature"
  encoded_signature = (
      base64.urlsafe_b64encode(signature.encode()).rstrip(b"=").decode()
  )
  return f"{encoded_header}.{encoded_payload}.{encoded_signature}"


class JwtTest(absltest.TestCase):

  def test_payload_unsafe(self):
    got = jwt.PayloadUnsafe(_get_test_token())
    self.assertNotEmpty(got)

  def test_email(self):
    got = jwt.Email(_get_test_token())
    self.assertEqual(got, "doe@example.com")

  def test_expires_at(self):
    got = jwt.ExpiresAt(_get_test_token_with_expiry())
    self.assertEqual(
        got,
        datetime.datetime(2025, 4, 1, 0, 0, 0, tzinfo=datetime.timezone.utc),
    )


if __name__ == "__main__":
  absltest.main()
