# Copyright 2023 Intrinsic Innovation LLC

load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("//bazel:go_macros.bzl", "go_binary", "go_library")

package(
    default_visibility = ["//visibility:public"],
)

bzl_library(
    name = "services_bzl",
    srcs = ["services.bzl"],
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/util/proto/build_defs:descriptor_set",
    ],
)

go_library(
    name = "servicegen",
    srcs = ["servicegen.go"],
    deps = [
        "//intrinsic/assets:bundleio",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/proto:registryutil",
        "//intrinsic/util/proto:sourcecodeinfoview",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)

go_binary(
    name = "servicegen_main",
    srcs = ["servicegen_main.go"],
    deps = [
        ":servicegen",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/production:intrinsic",
        "//intrinsic/util/proto:protoio",
        "@com_github_golang_glog//:go_default_library",
    ],
)
