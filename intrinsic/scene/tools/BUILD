# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_binary.bzl", "cc_binary")

package(default_visibility = ["//visibility:public"])

cc_binary(
    name = "update_scene_object",
    srcs = ["update_scene_object.cc"],
    deps = [
        "//intrinsic/geometry/storage:dummy_storage",
        "//intrinsic/geometry/storage:gzf_storage",
        "//intrinsic/icon/release/portable:init_intrinsic",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_updates_cc_proto",
        "//intrinsic/scene/util:object_user_data",
        "//intrinsic/scene/util:scene_object_gzf",
        "//intrinsic/scene/util:scene_object_updates",
        "//intrinsic/scene/validate:scene_object_validate_geo",
        "//intrinsic/scene/validate:scene_object_validation",
        "//intrinsic/util/proto:descriptors",
        "//intrinsic/util/proto:dynamic_message_parser",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/gzfile",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/log:flags",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
        "@or_tools//ortools/base:file",
    ],
)
