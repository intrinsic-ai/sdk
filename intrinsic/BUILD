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

package(default_visibility = ["//visibility:public"])

###############################################################################
# Common package groups for visibility
# public_api_users contains the set of packages allowed to depend on
# Intrinsic's public APIs.
package_group(
    name = "public_api_users",
    packages = [
        "public",
    ],
)

# internal_api_users contains the set of packages allowed to depend on
# Intrinsic's internal APIs.
package_group(
    name = "internal_api_users",
    packages = [
        "//intrinsic/...",
    ],
)

# intrinsic_internal contains the set of packages allowed to depend on internal code.
package_group(
    name = "intrinsic_internal",
    packages = [
        "//intrinsic/...",
    ],
)

# public_asset_users contains the set of packages allowed to depend on assets built by Intrinsic.
package_group(
    name = "public_asset_users",
    packages = [
        "public",
    ],
)

###############################################################################
