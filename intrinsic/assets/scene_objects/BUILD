# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "sceneobjectmanifest",
    srcs = ["sceneobjectmanifest.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:metadatautils",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
