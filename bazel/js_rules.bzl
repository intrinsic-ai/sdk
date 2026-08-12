# Copyright 2023 Intrinsic Innovation LLC

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
