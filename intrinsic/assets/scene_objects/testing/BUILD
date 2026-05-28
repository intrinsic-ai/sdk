# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "utils",
    testonly = True,
    srcs = ["utils.go"],
    importpath = "intrinsic/assets/scene_objects/testing/utils",
    visibility = [
        "//intrinsic:intrinsic_internal",
    ],
    deps = [
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/scene/proto/v1:entity_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_config_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/descriptorpb:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
