# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "customer",
    srcs = [
        "customer.go",
        "organizations.go",
        "rolebindings.go",
        "roles.go",
        "users.go",
    ],
    visibility = ["//intrinsic/tools/inctl:__pkg__"],
    deps = [
        "//intrinsic/config:environments",
        "//intrinsic/kubernetes/accounts/service/api/accesscontrol/v1:accesscontrol_go_grpc_proto",
        "//intrinsic/kubernetes/accounts/service/api/resourcemanager/v1:resourcemanager_go_grpc_proto",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:cobrautil",
        "//intrinsic/tools/inctl/util:orgutil",
        "//intrinsic/tools/inctl/util:printer",
        "//intrinsic/tools/inctl/util/accounts",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
        "@org_golang_google_protobuf//encoding/protojson:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)
