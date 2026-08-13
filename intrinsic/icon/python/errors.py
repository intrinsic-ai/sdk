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

"""A common location for ICON exceptions.

Centralises the exceptions for the ICON module through a single import. A common
base Error is provided for users to easily handle all ICON related exceptions.

  Usage example:

  try:
    session.do_something(...)
  except errors.Session.ActionError:
    # Handle specific error.
"""


class Error(Exception):
  pass


class Client:
  """Errors raised by icon.py."""

  class ServerError(Error):
    """Errors related to server connection."""

  class InvalidArgumentError(Error):
    """Errors related to bad arguments."""


class Session:
  """Errors raised by _session.py."""

  class ActionError(Error):
    """Errors related to OpenSessionResponse."""

  class StreamError(Error):
    """Errors related to OpenWriteStreamResponse."""

  class SessionEndedError(Error):
    """Errors related to the (expected or unexpected) end of an ICON sesson."""
