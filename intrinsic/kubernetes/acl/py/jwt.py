# Copyright 2023 Intrinsic Innovation LLC

"""Helper for extracting claims from jwts."""

import base64
import json


def PayloadUnsafe(j: str) -> dict[str, str]:
  """Decode the jwt payload into a dict.

  Does not validate the signature.

  Args:
    j (str): A json web token.

  Returns:
    dict[str,str]: The payload.

  Raises:
    ValueError If the jwt cannot be parsed.
  """
  parts = j.split('.')
  if len(parts) < 3:
    raise ValueError('Invalid JWT, token must have 3 parts')
  payload_str = base64.urlsafe_b64decode(parts[1] + '==').decode('utf-8')
  try:
    return json.loads(payload_str)
  except json.JSONDecodeError as e:
    raise ValueError('Error parsing json') from e


def Email(j: str) -> str:
  """Return the email claim from a jwt payload.

  Args:
    j (str): A json web token.

  Returns:
    str: The email.

  Raises:
    KeyError: If the email value is missing.
    ValueError: If the jwt cannot be parsed.
  """
  p = PayloadUnsafe(j)
  for k in ('email', 'uid'):
    if k in p:
      return p[k]
  raise KeyError('failed to extract email from JWT')
