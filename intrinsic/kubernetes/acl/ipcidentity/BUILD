# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_python//python:defs.bzl", "py_library", "py_test")

py_library(
    name = "ipcidentity_py",
    srcs = ["ipcidentity.py"],
    visibility = ["//intrinsic/util/grpc:__pkg__"],
    deps = [
        "//intrinsic/config:environments_py",
        "//intrinsic/kubernetes/accounts/service/api/tokens/v2:tokens_py_pb2",
        "//intrinsic/kubernetes/accounts/service/api/tokens/v2:tokens_py_pb2_grpc",
        "//intrinsic/kubernetes/acl/ipcidentity/internal:metadata",
        "//intrinsic/kubernetes/acl/py:jwt",
        requirement("grpcio"),
    ],
)

py_test(
    name = "ipcidentity_py_test",
    srcs = ["ipcidentity_py_test.py"],
    deps = [
        ":ipcidentity_py",
        "//intrinsic/kubernetes/accounts/service/api/tokens/v2:tokens_py_pb2",
        "//intrinsic/kubernetes/acl/ipcidentity/internal:metadata",
        requirement("grpcio"),
        "@com_google_absl_py//absl/testing:absltest",
    ],
)
