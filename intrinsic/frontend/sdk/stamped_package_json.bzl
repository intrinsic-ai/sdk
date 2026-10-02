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

"""Rule to stamp package.json version and dependencies for release builds."""

load("@bazel_lib//lib:stamping.bzl", "STAMP_ATTRS", "maybe_stamp")

def _stamped_package_json_impl(ctx):
    out = ctx.actions.declare_file(ctx.attr.out)
    stamp = maybe_stamp(ctx)
    status_files = [stamp.volatile_status_file, stamp.stable_status_file] if stamp else []

    args = ctx.actions.args()
    args.add("--src", ctx.file.src)
    args.add("--out", out)
    args.add("--stamp_var", ctx.attr.stamp_var)
    for dep in ctx.attr.stamp_dependencies:
        args.add("--stamp_dependency", dep)
    for sf in status_files:
        args.add("--status_file", sf)

    ctx.actions.run(
        arguments = [args],
        executable = ctx.executable._stamper,
        inputs = [ctx.file.src] + status_files,
        mnemonic = "StampPackageDependencies",
        outputs = [out],
        progress_message = "Stamping %s" % out.short_path,
        tools = [ctx.attr._stamper[DefaultInfo].files_to_run],
    )
    return [DefaultInfo(files = depset([out]))]

_stamped_package_json = rule(
    attrs = {
        "out": attr.string(default = "package.json"),
        "src": attr.label(
            allow_single_file = [".json"],
            mandatory = True,
        ),
        "stamp_dependencies": attr.string_list(default = []),
        "stamp_var": attr.string(default = "STABLE_SEMVER"),
        "_stamper": attr.label(
            cfg = "exec",
            default = Label("//intrinsic/frontend/sdk:stamp_package_dependencies"),
            executable = True,
        ),
    } | STAMP_ATTRS,
    implementation = _stamped_package_json_impl,
)

def stamped_package_json(
        name = "package",
        src = "package.json",
        stamp_dependencies = [],
        stamp_var = "STABLE_SEMVER",
        **kwargs):
    """Stamps package.json version and specific internal dependencies.

    Args:
        name: Target name.
        src: Input package.json file.
        stamp_dependencies: List of dependency package names in package.json
            whose versions should also be set to the stamped version.
        stamp_var: Name of the workspace status variable to read (default: STABLE_SEMVER).
        **kwargs: Additional rule attributes.
    """
    _stamped_package_json(
        name = name,
        src = src,
        stamp_dependencies = stamp_dependencies,
        stamp_var = stamp_var,
        **kwargs
    )
