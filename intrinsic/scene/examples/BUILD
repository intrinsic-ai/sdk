# Copyright 2023 Intrinsic Innovation LLC

load("//intrinsic/scene/build_defs:sdf_scene_object.bzl", "sdf_scene_object")

package(
    default_visibility = ["//visibility:public"],
)

sdf_scene_object(
    name = "box_scene_object",
    src = "box.sdf",
)

sdf_scene_object(
    name = "box_with_updates_scene_object",
    src = "box.sdf",
    updates_pbtxts = ["box_updates.pbtxt"],
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_updates_proto",
    ],
)
