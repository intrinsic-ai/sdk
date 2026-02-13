# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//intrinsic/tools/inctl/cmd/vm:__subpackages__"])

go_library(
    name = "pool",
    srcs = [
        "create.go",
        "describe.go",
        "leases.go",
        "list.go",
        "pool.go",
        "state.go",
    ],
    deps = [
        "//intrinsic/assets:cmdutils",
        "//intrinsic/kubernetes/acl:identity",
        "//intrinsic/kubernetes/vmpool/manager/api/v1:lease_api_go_proto",
        "//intrinsic/kubernetes/vmpool/service/api/v1:vmpool_api_go_proto",
        "//intrinsic/kubernetes/vmpool/service/pkg/defaults",
        "//intrinsic/tools/inctl/auth",
        "//intrinsic/tools/inctl/util:cobrautil",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
        "@io_opencensus_go//trace:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//encoding/protojson:go_default_library",
        "@org_golang_google_protobuf//encoding/prototext:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/known/emptypb",
        "@org_golang_google_protobuf//types/known/fieldmaskpb",
    ],
)
