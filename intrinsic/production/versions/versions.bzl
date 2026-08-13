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

"""Constants for Intrinsic OS versions."""

# NOTE: careful! the OS versions here are *without* the prefix `xfa.` which
# they do use in versions.go. Make sure not to add `xfa.` here, as it will
# break the fleet-manager.
PREVIOUS_OS_VERSION = "20260617.RC02"

# Version that is currently running with the intrinsic stack
STABLE_OS_VERSION = "20260710.RC01"

CANARY_OS_VERSION = "20260710.RC01"
