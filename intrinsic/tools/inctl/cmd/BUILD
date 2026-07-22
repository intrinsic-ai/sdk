# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "icon",
    srcs = [
        "icon.go",
        "icon_clear_faults.go",
        "icon_config.go",
        "icon_disable.go",
        "icon_enable.go",
        "icon_list_actions.go",
        "icon_list_parts.go",
        "icon_status.go",
    ],
    importpath = "intrinsic/tools/inctl/cmd/icon",
    deps = [
        ":root",
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/icon/go:icon",
        "//intrinsic/icon/proto/v1:types_go_proto",
        "//intrinsic/tools/inctl/util:cobrautil",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
        "@org_golang_google_protobuf//encoding/prototext:go_default_library",
    ],
)

go_library(
    name = "root",
    srcs = [
        "full_help.go",
        "root.go",
    ],
    importpath = "intrinsic/tools/inctl/cmd/root",
    deps = [
        "//intrinsic/production:intrinsic",
        "//intrinsic/skills/tools/skill/cmd:dialerutil",
        "//intrinsic/tools/inctl/util:orgutil",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@io_opencensus_go//trace:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_x_exp//slices",
    ],
)

go_library(
    name = "markdown",
    srcs = ["markdown.go"],
    importpath = "intrinsic/tools/inctl/cmd/markdown",
    deps = [
        ":root",
        "//intrinsic/util/archive:tartooling",
        "@com_github_google_safearchive//tar",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_cobra//doc:go_default_library",
    ],
)
