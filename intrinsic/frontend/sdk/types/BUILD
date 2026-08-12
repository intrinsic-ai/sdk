# Copyright 2023 Intrinsic Innovation LLC

load("@aspect_rules_js//npm:defs.bzl", "npm_package")
load("@bazel_skylib//rules:copy_file.bzl", "copy_file")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")

package(default_visibility = ["//visibility:public"])

copy_file(
    name = "license",
    src = "//intrinsic/frontend/sdk:LICENSE",
    out = "LICENSE",
)

npm_package(
    name = "package",
    srcs = [
        "package.json",
        ":license",
        "//intrinsic/assets/proto/v1:asset_configuration_jspb_proto",
        "//intrinsic/executive/proto:blackboard_service_jspb_proto",
        "//intrinsic/executive/proto:executive_service_jspb_proto",
        "//intrinsic/executive/proto:proto_builder_jspb_proto",
        "//intrinsic/frontend/solution_service/proto:solution_service_jspb_proto",
        "//intrinsic/geometry/proto:geometry_service_jspb_proto",
        "//intrinsic/httpjson/proto/v1alpha:controlplane_jspb_proto",
        "//intrinsic/httpjson/proto/v1alpha:dataplane_jspb_proto",
        "//intrinsic/httpjson/proto/v1alpha:signalling_jspb_proto",
        "//intrinsic/httpjson/proto/v1alpha:webpubsub_jspb_proto",
        "//intrinsic/perception/proto/v1:train_service_jspb_proto",
        "//intrinsic/platform/pubsub/adapters:pubsub_jspb_proto",
        "//intrinsic/proto_tools/proto:proto_registry_jspb_proto",
        "//intrinsic/scene/proto/v1:scene_object_import_jspb_proto",
        "//intrinsic/world/proto:object_world_service_jspb_proto",
    ],
    include_external_repositories = ["*"],
    include_srcs_packages = ["**"],
    package = "@intrinsic-ai/types",
    publishable = True,
)

pkg_tar(
    name = "intrinsic-ai-types",
    srcs = [":package"],
    extension = "tgz",
)

alias(
    name = "pack",
    actual = ":intrinsic-ai-types",
)
