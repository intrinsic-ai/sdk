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
