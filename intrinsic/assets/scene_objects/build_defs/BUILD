# Copyright 2023 Intrinsic Innovation LLC

load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("//bazel:go_macros.bzl", "go_binary", "go_library")

# Starlark rules and macros for building scene object related targets.

package(
    default_visibility = ["//visibility:public"],
)

bzl_library(
    name = "scene_object_bzl",
    srcs = ["scene_object.bzl"],
    deps = [
        "//intrinsic/assets/build_defs:asset_bzl",
        "//intrinsic/scene/build_defs:scene_object_bzl",
        "//intrinsic/util/proto/build_defs:descriptor_set",
        "@bazel_lib//lib:paths.bzl",
    ],
)

go_library(
    name = "sceneobjectgen",
    srcs = ["sceneobjectgen.go"],
    importpath = "intrinsic/assets/scene_objects/build_defs/sceneobjectgen",
    deps = [
        "//intrinsic/assets/scene_objects:sceneobjectbundle",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/proto:registryutil",
        "//intrinsic/util/proto:sourcecodeinfoview",
    ],
)

go_binary(
    name = "sceneobjectgen_main",
    srcs = ["sceneobjectgen_main.go"],
    importpath = "intrinsic/assets/scene_objects/build_defs/sceneobjectgen_main",
    deps = [
        ":sceneobjectgen",
        "//intrinsic/production:intrinsic",
        "//intrinsic/util:flag",
        "@com_github_golang_glog//:go_default_library",
    ],
)
