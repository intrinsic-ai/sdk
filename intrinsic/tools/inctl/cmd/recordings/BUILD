# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "recordings",
    srcs = [
        "generate.go",
        "get.go",
        "list.go",
        "recordings.go",
        "server.go",
        "visualize.go",
    ],
    deps = [
        "//intrinsic/assets:cmdutils",
        "//intrinsic/kubernetes/vmpool/manager/api/v1:lease_api_go_grpc_proto",
        "//intrinsic/logging/proto:bag_metadata_go_proto",
        "//intrinsic/logging/proto:bag_packager_service_go_grpc_proto",
        "//intrinsic/logging/proto:replay_service_go_grpc_proto",
        "//intrinsic/tools/inctl/auth",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:cobrautil",
        "//intrinsic/tools/inctl/util:color",
        "//intrinsic/tools/inctl/util:orgutil",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_pborman_uuid//:go_default_library",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//encoding/prototext:go_default_library",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

go_test(
    name = "recordings_test",
    srcs = [
        "generate_test.go",
        "get_test.go",
        "list_test.go",
        "visualize_test.go",
    ],
    library = ":recordings",
    deps = [
        "//intrinsic/kubernetes/vmpool/manager/api/v1:lease_api_go_grpc_proto",
        "//intrinsic/kubernetes/vmpool/manager/api/v1:lease_api_go_proto",
        "//intrinsic/logging/proto:bag_metadata_go_proto",
        "//intrinsic/logging/proto:bag_packager_service_go_grpc_proto",
        "//intrinsic/logging/proto:replay_service_go_grpc_proto",
        "//intrinsic/logging/proto:replay_service_go_proto",
        "//intrinsic/tools/inctl/util:orgutil",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//mock",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)
