# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//intrinsic/tools/inbuild:__subpackages__"])

go_library(
    name = "skillmanifest",
    srcs = ["skillmanifest.go"],
    deps = [
        "//intrinsic/skills/internal:skillmanifest",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/proto:registryutil",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
    ],
)
