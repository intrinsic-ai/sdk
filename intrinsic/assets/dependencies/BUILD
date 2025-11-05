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
        "//intrinsic/assets/data/proto/v1:data_assets_go_grpc_proto",
        "//intrinsic/assets/proto:field_metadata_go_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_go_proto",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/insecure:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)

go_test(
    name = "utils_test",
    srcs = ["utils_test.go"],
    library = ":utils",
    deps = [
        "//intrinsic/assets/data:fakedataassets",
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/dependencies/testing:test_configs_go_proto",
        "//intrinsic/assets/dependencies/testing:test_service_go_grpc_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_go_proto",
        "//intrinsic/testing:grpctest",
        "//intrinsic/util/proto:descriptor",
        "//intrinsic/util/proto/testing:prototestutil",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
        "@org_golang_google_protobuf//types/known/emptypb",
    ],
)

cc_library(
    name = "utils_cc",
    srcs = ["utils.cc"],
    hdrs = ["utils.h"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/data/proto/v1:data_assets_cc_grpc_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_cc_proto",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "utils_test_cc",
    srcs = ["utils_test.cc"],
    deps = [
        ":utils_cc",
        "//intrinsic/assets/data:fake_data_assets_cc",
        "//intrinsic/assets/data/proto/v1:data_asset_cc_proto",
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
        "@com_google_protobuf//:protobuf",
    ],
)

py_library(
    name = "utils_py",
    srcs = ["utils.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/data/proto/v1:data_assets_py_pb2",
        "//intrinsic/assets/data/proto/v1:data_assets_py_pb2_grpc",
        "//intrinsic/assets/proto/v1:resolved_dependency_py_pb2",
        "//intrinsic/util/grpc:interceptor",
        requirement("grpcio"),
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "utils_py_test",
    srcs = ["utils_test.py"],
    main = "utils_test.py",
    deps = [
        ":utils_py",
        "//intrinsic/assets/data:fake_data_assets_py",
        "//intrinsic/assets/data/proto/v1:data_asset_py_pb2",
        "//intrinsic/assets/dependencies/testing:test_service_py_grpc_proto",
        "//intrinsic/assets/dependencies/testing:test_service_py_pb2",
        "//intrinsic/assets/proto:asset_type_py_pb2",
        "//intrinsic/assets/proto:id_py_pb2",
        "//intrinsic/assets/proto:metadata_py_pb2",
        "//intrinsic/assets/proto/v1:resolved_dependency_py_pb2",
        requirement("grpcio"),
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
        "@com_google_protobuf//:protobuf_python",
    ],
)
