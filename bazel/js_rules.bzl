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

"""Conditional wrapper for js_library in Bzlmod workspaces.

In root workspaces (such as internal builds or SDK NPM building), it forwards to rules_js.
When consumed as an external downstream module, it acts as a no-op filegroup to avoid dev dependencies.
"""

def _js_rules_repo_impl(ctx):
    if ctx.attr.is_root:
        defs_content = """load("@aspect_rules_js//js:defs.bzl", _real_js_library = "js_library")

def js_library(**kwargs):
    _real_js_library(**kwargs)
"""
    else:
        defs_content = """def js_library(**kwargs):
    native.filegroup(
        name = kwargs.get("name"),
        visibility = kwargs.get("visibility"),
    )
"""
    ctx.file("BUILD.bazel", 'package(default_visibility = ["//visibility:public"])')
    ctx.file("defs.bzl", defs_content)

_js_rules_repo = repository_rule(
    attrs = {"is_root": attr.bool()},
    implementation = _js_rules_repo_impl,
)

def _js_rules_impl(ctx):
    is_root = False
    for mod in ctx.modules:
        if mod.is_root:
            is_root = True
    _js_rules_repo(name = "intrinsic_js_rules", is_root = is_root)

js_rules = module_extension(implementation = _js_rules_impl)
