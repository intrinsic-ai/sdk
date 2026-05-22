# Copyright 2023 Intrinsic Innovation LLC

# Provides generators for creating skill source code.

load("@rules_python//python:defs.bzl", "py_library")
load("//bazel:go_macros.bzl", "go_binary", "go_library")
load(
    "//intrinsic/util/proto/build_defs:descriptor_set.bzl",
    "proto_source_code_info_transitive_descriptor_set",
)

package(
    default_visibility = ["//visibility:private"],
)

go_library(
    name = "gen",
    srcs = ["gen.go"],
    embedsrcs = [
        "skill_service_main.py.tmpl",
        "skill_service_main_tmpl.cc",
    ],
    importpath = "intrinsic/skills/generator/gen",
    visibility = ["//intrinsic/tools/inbuild:__subpackages__"],
    deps = ["//intrinsic/skills/proto:skill_manifest_go_proto"],
)

proto_source_code_info_transitive_descriptor_set(
    name = "skill_services_provided_to_platform",
    visibility = ["//visibility:private"],
    deps = ["//intrinsic/skills/proto:skill_service_proto"],
)

go_binary(
    name = "skill_service_generator",
    srcs = ["skill_service_generator.go"],
    embedsrcs = [":skill_services_provided_to_platform"],
    importpath = "intrinsic/skills/generator/skill_service_generator",
    visibility = ["//visibility:public"],
    deps = [
        ":gen",
        "//intrinsic/production:intrinsic",
        "//intrinsic/skills:skillmanifest",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util/proto:descriptor",
        "//intrinsic/util/proto:protoio",
        "@com_github_golang_glog//:go_default_library",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_protobuf//proto",
    ],
)

py_library(
    name = "app",
    srcs = ["app.py"],
    visibility = ["//visibility:public"],
    deps = ["@com_google_absl_py//absl:app"],
)
