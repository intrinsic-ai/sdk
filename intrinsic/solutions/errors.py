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

"""Common errors in the solution building library."""

import retrying


class Error(Exception):
  """Top-level module error for the solution building library."""


class InvalidArgumentError(Error):
  """Thrown when invalid arguments are passed to the solution building library."""


class NotFoundError(Error):
  """Thrown when an element cannot be found in the solution building library."""


class UnavailableError(Error):
  """Thrown when a backend of the solution building library cannot be reached."""


class FailedPreconditionError(Error):
  """Thrown when a precondition about the state of the solution is broken."""


class BackendPendingError(Error):
  """Thrown if a backend is unhealthy but expected to become healthy again."""


class BackendHealthError(Error):
  """Thrown if a backend is unhealthy and not expected to recover."""


class BackendNoWorkcellError(Error):
  """Thrown if no workcell spec has been installed."""


def _is_backend_pending_error(e: Exception) -> bool:
  """Determines whether a backend is expected to become healthy.

  Args:
    e: The exception under evaluation.

  Returns:
    True if the exception indicates that the backend is expected to become
      healthy again.
  """
  return isinstance(e, BackendPendingError)


# Decorator that retries if a backend's health is expected to recover.
retry_on_pending_backend = retrying.retry(
    retry_on_exception=_is_backend_pending_error,
    stop_max_attempt_number=10,
    wait_exponential_multiplier=2,
    wait_exponential_max=4000,  # in milliseconds
    wait_incrementing_start=250,  # in milliseconds
)
