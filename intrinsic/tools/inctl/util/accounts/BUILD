# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "accounts",
    srcs = [
        "accounts.go",
        "clients.go",
        "operations.go",
    ],
    importpath = "intrinsic/tools/inctl/util/accounts/accounts",
    visibility = ["//intrinsic/tools/inctl:__subpackages__"],
    deps = [
        "//intrinsic/config:environments",
        "//intrinsic/kubernetes/accounts/service/api/accesscontrol/v1:accesscontrol_go_proto",
        "//intrinsic/kubernetes/accounts/service/api/resourcemanager/v1:resourcemanager_go_proto",
        "//intrinsic/tools/inctl/auth",
        "//intrinsic/tools/inctl/util:orgutil",
        "@com_github_spf13_viper//:go_default_library",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
