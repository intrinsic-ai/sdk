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
