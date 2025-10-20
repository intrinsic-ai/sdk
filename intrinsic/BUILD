# Copyright 2023 Intrinsic Innovation LLC

package(default_visibility = ["//visibility:public"])

# Do not use insrc_filegroup to avoid the dependency from various //third_party libs, see cl/756637985.

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
    packages = ["//intrinsic/..."],
)

# public_asset_users contains the set of packages allowed to depend on assets built by Intrinsic.
package_group(
    name = "public_asset_users",
    packages = [
        "public",
    ],
)

###############################################################################
