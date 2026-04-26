# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "recordings",
    srcs = [
        "create.go",
        "generate.go",
        "get.go",
        "json_util.go",
        "list.go",
        "recordings.go",
        "server.go",
        "visualize.go",
    ],
    importpath = "intrinsic/tools/inctl/cmd/recordings/recordings",
    deps = [
        "//intrinsic/assets:cmdutils",
        "//intrinsic/kubernetes/vmpool/manager/api/v1:lease_api_go_proto",
        "//intrinsic/logging/proto:bag_metadata_go_proto",
        "//intrinsic/logging/proto:bag_packager_service_go_proto",
        "//intrinsic/logging/proto:logger_service_go_proto",
        "//intrinsic/logging/proto:recording_data_categories_go",
        "//intrinsic/logging/proto:recording_data_categories_go_proto",
        "//intrinsic/logging/proto:replay_service_go_proto",
        "//intrinsic/logging/proto:replay_service_options_go_proto",
        "//intrinsic/tools/inctl/auth",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util",
        "//intrinsic/tools/inctl/util:cobrautil",
        "//intrinsic/tools/inctl/util:color",
        "//intrinsic/tools/inctl/util:orgutil",
        "//intrinsic/tools/inctl/util:printer",
        "//intrinsic/tools/inctl/util/promptutil",
        "@com_github_cenkalti_backoff_v4//:go_default_library",
        "@com_github_pborman_uuid//:go_default_library",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//encoding/protojson:go_default_library",
        "@org_golang_google_protobuf//encoding/prototext:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

go_test(
    name = "recordings_test",
    srcs = [
        "create_test.go",
        "generate_test.go",
        "get_test.go",
        "list_test.go",
        "visualize_test.go",
    ],
    embed = [":recordings"],
    importpath = "intrinsic/tools/inctl/cmd/recordings/recordings_test",
    deps = [
        "//intrinsic/kubernetes/vmpool/manager/api/v1:lease_api_go_proto",
        "//intrinsic/logging/proto:bag_metadata_go_proto",
        "//intrinsic/logging/proto:bag_packager_service_go_proto",
        "//intrinsic/logging/proto:logger_service_go_proto",
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
