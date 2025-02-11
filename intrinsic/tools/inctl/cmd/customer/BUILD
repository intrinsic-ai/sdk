# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "customer",
    srcs = [
        "adduser.go",
        "clients.go",
        "customer.go",
        "organizations.go",
        "rolebindings.go",
        "roles.go",
    ],
    visibility = ["//intrinsic/tools/inctl:__pkg__"],
    deps = [
        "//intrinsic/config:environments",
        "//intrinsic/kubernetes/accounts/service/api/accesscontrol/v1:accesscontrol_go_grpc_proto",
        "//intrinsic/kubernetes/accounts/service/api/resourcemanager/v1:resourcemanager_go_grpc_proto",
        "//intrinsic/kubernetes/acl:cookies",
        "//intrinsic/tools/inctl/auth",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:cobrautil",
        "//intrinsic/tools/inctl/util:orgutil",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@io_opencensus_go//plugin/ocgrpc:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
        "@org_golang_google_protobuf//encoding/protojson:go_default_library",
        "@org_golang_google_protobuf//proto",
    ],
)
