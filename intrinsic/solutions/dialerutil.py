# Copyright 2023 Intrinsic Innovation LLC

"""Deprecated. Import from intrinsic.util.grpc.dialerutil instead."""

import warnings

from intrinsic.util.grpc import dialerutil

warnings.warn(
    "The intrinsic.solutions.dialerutil module has been moved, import"
    " intrinsic.util.grpc.dialerutil instead",
    DeprecationWarning,
    stacklevel=2,
)

create_channel_from_address = dialerutil.create_channel_from_address
create_channel_from_org = dialerutil.create_channel_from_org
create_channel_from_cluster = dialerutil.create_channel_from_cluster
create_channel_from_solution = dialerutil.create_channel_from_solution
create_channel_from_token = dialerutil.create_channel_from_token
