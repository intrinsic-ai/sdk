# Copyright 2023 Intrinsic Innovation LLC

"""Implements go_grpc_http_proto_library Bazel macro."""

load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("//bazel:go_macros.bzl", "calculate_importpath")

def go_grpc_http_library(name, protos, **kwargs):
    """Generate Golang code for a gRPC service supporting gRPC and gRPC Gateway.

    Args:
      name: The name of the generated target
      protos: The proto_library targets this macro should generate Golang code for
      **kwargs: Everything else you would normally pass to go_proto_library.
    """

    go_proto_library(
        name = name,
        compilers = [
            # Existing default compilers
            "@io_bazel_rules_go//proto:go_proto",
            "@io_bazel_rules_go//proto:go_grpc_v2",
            # Add compiler that generates code for grpc-gateway
            Label("//bazel:go_gen_grpc_gateway"),
        ],
        protos = protos,
        importpath = calculate_importpath(name, kwargs.pop("importpath", None)),
        **kwargs
    )
