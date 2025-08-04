# Copyright 2023 Intrinsic Innovation LLC

load("@rules_python//python:defs.bzl", "py_library", "py_test")

package(default_visibility = ["//intrinsic:public_api_users"])

py_library(
    name = "behavior_tree_visitor",
    srcs = ["behavior_tree_visitor.py"],
    deps = ["//intrinsic/executive/proto:behavior_tree_py_pb2"],
)

py_test(
    name = "behavior_tree_visitor_test",
    srcs = ["behavior_tree_visitor_test.py"],
    deps = [
        ":behavior_tree_visitor",
        "//intrinsic/executive/proto:behavior_tree_py_pb2",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
        "@com_google_protobuf//:protobuf_python",
    ],
)
