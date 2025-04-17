# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

go_library(
    name = "sceneobjectmanifest",
    srcs = ["sceneobjectmanifest.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets:metadatautils",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
    ],
)

go_test(
    name = "sceneobjectmanifest_test",
    srcs = ["sceneobjectmanifest_test.go"],
    library = ":sceneobjectmanifest",
    deps = [
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "@org_golang_google_protobuf//proto",
    ],
)
