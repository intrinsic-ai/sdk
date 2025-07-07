# Copyright 2023 Intrinsic Innovation LLC

load("@pybind11_bazel//:build_defs.bzl", "pybind_extension")
load("@rules_cc//cc:defs.bzl", "cc_library", "cc_test")
load("@rules_python//python:defs.bzl", "py_library")
load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "any",
    hdrs = ["any.h"],
    deps = [
        ":merge",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:str_format",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "any_test",
    size = "small",
    srcs = ["any_test.cc"],
    deps = [
        ":any",
        "//intrinsic/util/proto/testing:param_message_cc_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@com_google_protobuf//:protobuf",
    ],
)

go_library(
    name = "descriptor",
    srcs = ["descriptor.go"],
    deps = [
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
    ],
)

cc_library(
    name = "descriptors",
    srcs = ["descriptors.cc"],
    hdrs = ["descriptors.h"],
    deps = [
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings:str_format",
        "@com_google_protobuf//:protobuf",
    ],
)

py_library(
    name = "descriptors_py",
    srcs = ["descriptors.py"],
    deps = ["@com_google_protobuf//:protobuf_python"],
)

cc_library(
    name = "get_text_proto",
    srcs = ["get_text_proto.cc"],
    hdrs = ["get_text_proto.h"],
    deps = [
        ":any",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "merge",
    srcs = ["merge.cc"],
    hdrs = ["merge.h"],
    deps = [
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/types:span",
        "@com_google_protobuf//:protobuf",
    ],
)

go_library(
    name = "protoio",
    srcs = ["protoio.go"],
    deps = [
        "@com_github_protocolbuffers_txtpbfmt//config:go_default_library",
        "@com_github_protocolbuffers_txtpbfmt//parser:go_default_library",
        "@org_golang_google_protobuf//encoding/prototext:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)

go_library(
    name = "resolvercache",
    srcs = ["resolver_cache.go"],
    deps = [
        ":registryutil",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)

go_library(
    name = "registryutil",
    srcs = ["registryutil.go"],
    deps = [
        ":protoio",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
        "@org_golang_google_protobuf//types/dynamicpb:go_default_library",
    ],
)

cc_library(
    name = "source_code_info_view",
    srcs = ["source_code_info_view.cc"],
    hdrs = ["source_code_info_view.h"],
    deps = [
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@com_google_protobuf//:protobuf",
        "@com_google_protobuf//:protobuf_lite",
    ],
)

pybind_extension(
    name = "source_code_info_view_py",
    srcs = ["source_code_info_view_py.cc"],
    deps = [
        ":source_code_info_view",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:string_view",
        "@com_google_protobuf//:protobuf",
        "@com_google_protobuf//:protobuf_lite",
        "@pybind11_abseil//pybind11_abseil:absl_casters",
        "@pybind11_abseil//pybind11_abseil:status_casters",
        "@pybind11_protobuf//pybind11_protobuf:wrapped_proto_caster",
    ],
)

cc_library(
    name = "type_url",
    hdrs = ["type_url.h"],
    deps = [
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "type_url_test",
    srcs = ["type_url_test.cc"],
    deps = [
        ":type_url",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/strings:string_view",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "error_collector",
    hdrs = ["error_collector.h"],
    deps = [
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/strings:string_view",
        "@com_google_protobuf//src/google/protobuf/io:tokenizer",
    ],
)

cc_library(
    name = "status_specs",
    hdrs = ["status_specs.h"],
    deps = [
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "parsed_type_url",
    srcs = ["parsed_type_url.cc"],
    hdrs = [
        "parsed_type_url.h",
    ],
    deps = [
        ":status_specs",
        ":type_url",
        "//intrinsic/util/status:status_builder",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/strings:string_view",
    ],
)

cc_test(
    name = "parsed_type_url_test",
    srcs = ["parsed_type_url_test.cc"],
    deps = [
        ":parsed_type_url",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "parse_text_proto",
    srcs = ["parse_text_proto.cc"],
    hdrs = ["parse_text_proto.h"],
    deps = [
        ":error_collector",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
        "@com_google_protobuf//:protobuf",
        "@com_googlesource_code_re2//:re2",
    ],
)

cc_test(
    name = "parse_text_proto_test",
    srcs = ["parse_text_proto_test.cc"],
    deps = [
        ":parse_text_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@com_google_protobuf//:protobuf",
    ],
)

go_library(
    name = "sourcecodeinfoview",
    srcs = ["source_code_info_view.go"],
    deps = [
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
    ],
)

go_library(
    name = "names",
    srcs = ["names.go"],
    deps = ["@org_golang_google_protobuf//types/known/anypb"],
)

go_test(
    name = "names_test",
    srcs = ["names_test.go"],
    library = ":names",
    deps = [
        "//intrinsic/util/proto/testing:diamond_a_go_proto",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/known/anypb",
        "@org_golang_google_protobuf//types/known/emptypb",
    ],
)
