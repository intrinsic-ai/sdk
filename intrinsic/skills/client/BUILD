# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_python//python:defs.bzl", "py_library", "py_test")

package(default_visibility = ["//visibility:public"])

py_library(
    name = "skill_registry_client",
    srcs = ["skill_registry_client.py"],
    deps = [
        "//intrinsic/skills/proto:skill_registry_py_pb2",
        "//intrinsic/skills/proto:skill_registry_py_pb2_grpc",
        "//intrinsic/skills/proto:skills_py_pb2",
        "//intrinsic/util/grpc:error_handling",
        requirement("grpcio"),
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "skill_registry_client_test",
    srcs = ["skill_registry_client_test.py"],
    deps = [
        ":skill_registry_client",
        "//intrinsic/skills/proto:skill_registry_py_pb2",
        "//intrinsic/skills/proto:skills_py_pb2",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)
