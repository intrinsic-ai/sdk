# Copyright 2023 Intrinsic Innovation LLC

"""Deprecated. Import from intrinsic.util.grpc.auth instead."""


import warnings

from intrinsic.util.grpc import auth

warnings.warn(
    "The intrinsic.solutions.auth module has been moved, import"
    " intrinsic.util.grpc.auth instead",
    DeprecationWarning,
    stacklevel=2,
)

ALIAS_DEFAULT_TOKEN = auth.ALIAS_DEFAULT_TOKEN
AUTH_CONFIG_EXTENSION = auth.AUTH_CONFIG_EXTENSION
ORG_STORE_DIRECTORY = auth.ORG_STORE_DIRECTORY
STORE_DIRECTORY = auth.STORE_DIRECTORY

ProjectToken = auth.ProjectToken
ProjectConfiguration = auth.ProjectConfiguration
CredentialsNotFoundError = auth.CredentialsNotFoundError
get_configuration = auth.get_configuration
OrgInfo = auth.OrgInfo
OrgNotFoundError = auth.OrgNotFoundError
parse_info_from_string = auth.parse_info_from_string
read_org_info = auth.read_org_info
