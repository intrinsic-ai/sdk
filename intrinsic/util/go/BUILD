# Copyright 2023 Intrinsic Innovation LLC

# Generic language utilities for golang.

load("//bazel:go_macros.bzl", "go_library", "go_test")

go_library(
    name = "pointer",
    srcs = ["pointer.go"],
    importpath = "intrinsic/util/go/pointer",
    visibility = ["//intrinsic:internal_api_users"],
)

go_library(
    name = "validate",
    srcs = ["validate.go"],
    importpath = "intrinsic/util/go/validate",
    visibility = ["//intrinsic:internal_api_users"],
    deps = ["//intrinsic/kubernetes/workcell_spec/proto:image_go_proto"],
)

go_test(
    name = "validate_test",
    srcs = ["validate_test.go"],
    embed = [":validate"],
    importpath = "intrinsic/util/go/validate_test",
    deps = ["//intrinsic/kubernetes/workcell_spec/proto:image_go_proto"],
)
