# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")
load("@rules_python//python:defs.bzl", "py_library", "py_test")
load("//bazel:go_macros.bzl", "go_library", "go_test")

go_library(
    name = "utils",
    srcs = ["utils.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto/v1:resolved_dependency_go_proto",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/insecure:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)

go_test(
    name = "utils_test",
    srcs = ["utils_test.go"],
    library = ":utils",
    deps = [
        "//intrinsic/assets/dependencies/testing:test_service_go_grpc_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_go_proto",
        "//intrinsic/testing:grpctest",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)

cc_library(
    name = "utils_cc",
    srcs = ["utils.cc"],
    hdrs = ["utils.h"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto/v1:resolved_dependency_cc_proto",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
        "@com_github_grpc_grpc//:grpc++",
    ],
)

cc_test(
    name = "utils_test_cc",
    srcs = ["utils_test.cc"],
    deps = [
        ":utils_cc",
        "//intrinsic/assets/dependencies/testing:test_service_cc_grpc_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_cc_proto",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpc_security_base",
    ],
)

py_library(
    name = "utils_py",
    srcs = ["utils.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto/v1:resolved_dependency_py_pb2",
        "//intrinsic/util/grpc:interceptor",
        requirement("grpcio"),
    ],
)

py_test(
    name = "utils_py_test",
    srcs = ["utils_test.py"],
    main = "utils_test.py",
    deps = [
        ":utils_py",
        "//intrinsic/assets/dependencies/testing:test_service_py_grpc_proto",
        "//intrinsic/assets/dependencies/testing:test_service_py_pb2",
        "//intrinsic/assets/proto:id_py_pb2",
        "//intrinsic/assets/proto/v1:resolved_dependency_py_pb2",
        requirement("grpcio"),
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
    ],
)
