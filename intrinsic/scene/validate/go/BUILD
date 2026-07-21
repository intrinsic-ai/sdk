# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "scene_object_validation_c",
    srcs = ["scene_object_validation_c.cc"],
    hdrs = ["scene_object_validation_c.h"],
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/scene/validate:scene_object_validation",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/status",
    ],
)

go_library(
    name = "sceneobjectvalidation",
    srcs = ["scene_object_validation.go"],
    cdeps = [
        ":scene_object_validation_c",
    ],
    cgo = 1,
    importpath = "intrinsic/scene/validate/go/sceneobjectvalidation",
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "@org_golang_google_protobuf//proto",
    ],
)

go_test(
    name = "sceneobjectvalidation_test",
    srcs = ["scene_object_validation_test.go"],
    embed = [":sceneobjectvalidation"],
    deps = [
        "//intrinsic/scene/proto/v1:entity_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
    ],
)
