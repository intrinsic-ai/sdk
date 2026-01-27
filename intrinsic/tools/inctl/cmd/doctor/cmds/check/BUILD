# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "check",
    srcs = [
        "check.go",
        "topological_sort.go",
    ],
    deps = [
        "//intrinsic/tools/inctl/cmd/doctor/api",
        "//intrinsic/tools/inctl/cmd/doctor/checks",
        "//intrinsic/tools/inctl/cmd/doctor/proto/v1:report_go_proto",
        "//intrinsic/tools/inctl/cmd/doctor/reporters",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_spf13_cobra//:go_default_library",
        "@org_golang_google_protobuf//encoding/protojson:go_default_library",
    ],
)

go_test(
    name = "check_test",
    srcs = ["topological_sort_test.go"],
    embed = [":check"],
    deps = ["@com_github_google_go_cmp//cmp:go_default_library"],
)
