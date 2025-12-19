# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//intrinsic/tools/inbuild:__subpackages__"])

go_library(
    name = "bundle",
    srcs = ["bundle.go"],
    deps = [
        "//intrinsic/assets/services/build_defs:servicegen",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/util/proto:protoio",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_library(
    name = "service",
    srcs = ["service.go"],
    deps = [
        ":bundle",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_test(
    name = "bundle_test",
    srcs = ["bundle_test.go"],
    data = [
        "//intrinsic/tools/inbuild/cmd/service/test_data:example_service_files",
        "//intrinsic/tools/inbuild/cmd/service/test_data:example_service_proto",
        "//intrinsic/tools/inbuild/cmd/service/test_data:example_service_py_image.tar",
    ],
    library = ":bundle",
    deps = [
        "//intrinsic/assets/services:servicebundle",
        "//intrinsic/tools/inbuild/cmd/service/test_data:example_service_go_proto",
        "//intrinsic/util/path_resolver:pathresolver",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//encoding/prototext:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
