# Copyright 2026 Intrinsic Innovation LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Helper for extracting claims from jwts."""

import base64
import datetime
import json


def PayloadUnsafe(j: str) -> dict[str, str]:
  """Decodes the jwt payload into a dict.

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
  """Returns the email claim from a jwt payload.

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


def ExpiresAt(j: str) -> datetime.datetime:
  """Returns the expiry claim from a jwt payload.

  Args:
    j (str): A json web token.

  Returns:
    datetime.datetime: The expiry time.

  Raises:
    KeyError: If the expiry value is missing.
    ValueError: If the jwt cannot be parsed.
  """
  p = PayloadUnsafe(j)
  if 'exp' not in p:
    raise KeyError('failed to extract expiry from JWT')
  try:
    return datetime.datetime.fromtimestamp(
        int(p['exp']), tz=datetime.timezone.utc
    )
  except (ValueError, TypeError) as e:
    raise ValueError('Error parsing expiry') from e
