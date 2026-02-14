# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//intrinsic:internal_api_users"])

go_library(
    name = "behaviortree",
    srcs = ["behavior_tree_visitor.go"],
    importpath = "intrinsic/executive/go/behaviortree",
    deps = ["//intrinsic/executive/proto:behavior_tree_go_proto"],
)

go_test(
    name = "behaviortreevisitor_test",
    srcs = ["behavior_tree_visitor_test.go"],
    importpath = "intrinsic/executive/go/behaviortreevisitor_test",
    deps = [
        ":behaviortree",
        "//intrinsic/executive/proto:behavior_tree_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//proto",
    ],
)
