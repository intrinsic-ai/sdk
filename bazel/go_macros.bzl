# Copyright 2023 Intrinsic Innovation LLC

"""This file contains starlark rules for building golang targets."""

load("@com_google_protobuf//bazel/common:proto_info.bzl", "ProtoInfo")
load("@io_bazel_rules_go//go:def.bzl", "GoInfo", _go_binary = "go_binary", _go_library = "go_library", _go_test = "go_test")
load("@io_bazel_rules_go//proto:def.bzl", _go_grpc_library = "go_grpc_library", _go_proto_library = "go_proto_library")

def calculate_importpath(name, importpath):
    label = native.package_relative_label(name)
    return label.package + "/" + label.name

def go_binary(name, **kwargs):
    if kwargs.get("linkmode", None) in (None, "normal", "pie"):  # executables
        # Create the output binary at "name" instead of "name_/name"
        kwargs["out"] = name

    _go_binary(
        name = name,
        importpath = calculate_importpath(name, kwargs.pop("importpath", None)),
        **kwargs
    )

def go_library(name, **kwargs):
    _go_library(
        name = name,
        importpath = calculate_importpath(name, kwargs.pop("importpath", None)),
        **kwargs
    )

def go_grpc_library(name, deps = None, **kwargs):
    _go_grpc_library(
        name = name,
        deps = deps,
        importpath = calculate_importpath(name, kwargs.pop("importpath", None)),
        **kwargs
    )

def go_proto_library(name, deps, go_deps = None, **kwargs):
    _go_proto_library(
        name = name,
        protos = deps,
        deps = go_deps,
        importpath = calculate_importpath(name, kwargs.pop("importpath", None)),
        **kwargs
    )

    importpath_matches_go_package_test(
        name = "importpath_matches_go_package_" + name + "_test",
        # Some proto_library files contain multiple proto files; however, protoc-gen-go requires
        # all proto files in a proto_library target have the same go_package. We only need to
        # check the go_package of one of the files.
        proto = deps[0],
        go_proto_library = name,
    )

def go_test(name, **kwargs):
    _go_test(
        name = name,
        importpath = calculate_importpath(name, kwargs.pop("importpath", None)),
        **kwargs
    )

def _importpath_matches_go_package_test_impl(ctx):
    """Implementation for the importpath_matches_go_package_test rule."""

    # The executable shell script that Bazel will run as the test.
    test_script = ctx.actions.declare_file(ctx.label.name + "_test.sh")

    # Call the checker binary with the import path and the proto file path.
    # The 'set -e' ensures that the script will exit with a non-zero status
    # if the checker fails, which then fails the test.
    ctx.actions.write(
        output = test_script,
        content = """#!/bin/bash
set -e
checker="{checker_path}"
import_path="{import_path}"
proto_file="{proto_file_path}"
"$checker" "$import_path" "$proto_file"
""".format(
            checker_path = ctx.executable._checker.short_path,
            import_path = ctx.attr.go_proto_library[GoInfo].importpath,
            proto_file_path = ctx.attr.proto[ProtoInfo].direct_sources[0].short_path,
        ),
        is_executable = True,
    )

    proto_runfiles = ctx.runfiles(ctx.attr.proto[ProtoInfo].direct_sources)

    return [DefaultInfo(
        executable = test_script,
        runfiles = ctx.attr._checker.default_runfiles.merge(proto_runfiles),
    )]

importpath_matches_go_package_test = rule(
    implementation = _importpath_matches_go_package_test_impl,
    test = True,
    attrs = {
        "go_proto_library": attr.label(
            mandatory = True,
            providers = [GoInfo],
            doc = "The go_proto_library target to check.",
        ),
        "proto": attr.label(
            mandatory = True,
            providers = [ProtoInfo],
            doc = "The proto_library target that has a go_package option.",
        ),
        "_checker": attr.label(
            default = Label("//bazel:checkimportpathmatchesgopackage"),
            cfg = "exec",
            executable = True,
            doc = "Internal: The binary that performs the check.",
        ),
    },
    doc = "A test rule that ensures a go_proto_library's importpath matches the go_package option in its .proto file.",
)
