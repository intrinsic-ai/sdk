# Copyright 2023 Intrinsic Innovation LLC

# Provides generators for creating skill source code.

load("@rules_python//python:defs.bzl", "py_library")
load("//bazel:go_macros.bzl", "go_binary", "go_library")

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

go_binary(
    name = "skill_service_generator",
    srcs = ["skill_service_generator.go"],
    importpath = "intrinsic/skills/generator/skill_service_generator",
    visibility = ["//visibility:public"],
    deps = [
        ":gen",
        "//intrinsic/production:intrinsic",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_protobuf//proto",
    ],
)

py_library(
    name = "app",
    srcs = ["app.py"],
    visibility = ["//visibility:public"],
    deps = ["@com_google_absl_py//absl:app"],
)
