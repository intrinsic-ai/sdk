# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "skillbundle",
    srcs = ["skillbundle.go"],
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        "//intrinsic/assets:imageutils",
        "//intrinsic/assets:ioutils",
        "//intrinsic/skills/internal:skillmanifest",
        "//intrinsic/skills/proto:processed_skill_manifest_go_proto",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util/archive:tartooling",
        "//intrinsic/util/proto:registryutil",
        "@com_github_google_safearchive//tar",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_protobuf//proto",
    ],
)
