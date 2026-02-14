# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "list",
    srcs = ["list.go"],
    importpath = "intrinsic/tools/inctl/cmd/solution/list/list",
    visibility = ["//intrinsic/tools/inctl/cmd/solution:__pkg__"],
    deps = [
        "//intrinsic/frontend/cloud/api/v1:clusterdiscovery_api_go_proto",
        "//intrinsic/frontend/cloud/api/v1:solutiondiscovery_api_go_proto",
        "//intrinsic/tools/inctl/auth",
        "//intrinsic/tools/inctl/util:orgutil",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
