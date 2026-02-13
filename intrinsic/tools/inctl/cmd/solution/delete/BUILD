# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "delete",
    srcs = ["delete.go"],
    visibility = [
        "//intrinsic/tools/inctl/cmd/solution:__pkg__",
    ],
    deps = [
        "//intrinsic/solution_versions/proto/v1:solution_version_service_go_proto",
        "//intrinsic/tools/inctl/auth",
        "//intrinsic/tools/inctl/util:color",
        "//intrinsic/tools/inctl/util:orgutil",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
