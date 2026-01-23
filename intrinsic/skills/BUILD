# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "skillbundle",
    srcs = ["skillbundle.go"],
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        ":skillmanifest",
        ":skillvalidate",
        "//intrinsic/assets:imageutils",
        "//intrinsic/assets:ioutils",
        "//intrinsic/skills/proto:processed_skill_manifest_go_proto",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util/archive:tartooling",
        "@com_github_google_safearchive//tar",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
    ],
)

go_library(
    name = "skillmanifest",
    srcs = ["skillmanifest.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util/proto:sourcecodeinfoview",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
    ],
)

go_library(
    name = "skillvalidate",
    srcs = ["skillvalidate.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:metadatautils",
        "//intrinsic/assets/dependencies:utils",
        "//intrinsic/skills/proto:processed_skill_manifest_go_proto",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)
