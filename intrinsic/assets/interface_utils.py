# Copyright 2023 Intrinsic Innovation LLC

"""Provides utilities for working with Asset URI interfaces."""

import re

# The prefix used for gRPC service dependencies.
GRPC_URI_PREFIX = "grpc://"
# The prefix used for proto-based data dependencies.
DATA_URI_PREFIX = "data://"

_URI_REGEX = re.compile(
    r"^(grpc://|data://)([A-Za-z_][A-Za-z0-9_]*\.)+[A-Za-z_][A-Za-z0-9_]*$"
)


def ValidateInterfaceName(uri: str) -> None:
  """Validates an interface name with a protocol prefix.

  Args:
    uri: The URI to validate.

  Raises:
    ValueError: If the URI is not formatted as
    '<protocol>://<package>.<message>'.
  """
  if not _URI_REGEX.fullmatch(uri):
    raise ValueError(
        "Invalid interface name: expected URI to be formatted as"
        f" '<protocol>://<package>.<message>', got '{uri}'"
    )
