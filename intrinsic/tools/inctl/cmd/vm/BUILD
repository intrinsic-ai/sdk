# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = [
    "//intrinsic/production:__subpackages__",
])

go_library(
    name = "vm",
    srcs = [
        "expire.go",
        "lease.go",
        "return.go",
        "vm.go",
    ],
    visibility = [
        "//intrinsic/assets:__pkg__",
        "//intrinsic/tools/inctl:__subpackages__",
        "//visibility:public",
    ],
    deps = [
        "//intrinsic/config:environments",
        "//intrinsic/kubernetes/acl:identity",
        "//intrinsic/kubernetes/vmpool/manager/api/v1:lease_api_go_grpc_proto",
        "//intrinsic/kubernetes/vmpool/service/api/v1:vmpool_api_go_grpc_proto",
        "//intrinsic/kubernetes/vmpool/service/pkg/defaults",
        "//intrinsic/tools/inctl/auth",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/cmd/vm/pool",
        "//intrinsic/tools/inctl/util:cobrautil",
        "//intrinsic/tools/inctl/util:color",
        "//intrinsic/tools/inctl/util:orgutil",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_pborman_uuid//:go_default_library",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_rs_xid//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
        "@io_opencensus_go//plugin/ocgrpc:go_default_library",
        "@io_opencensus_go//trace:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//credentials:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)
