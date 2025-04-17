# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_library")
load("//bazel:go_macros.bzl", "go_library", "go_test")
load("//intrinsic/scene/build_defs:sdf_scene_object.bzl", "sdf_scene_object")

go_library(
    name = "bundlegeometry",
    srcs = ["bundlegeometry.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":gzffile",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "//intrinsic/tools/inctl/util:casgeometryrewriter",
        "//intrinsic/tools/inctl/util:casgeometryuploader",
        "@com_github_golang_glog//:go_default_library",
    ],
)

go_test(
    name = "bundlegeometry_test",
    srcs = ["bundlegeometry_test.go"],
    data = [
    ],
    library = ":bundlegeometry",
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "//intrinsic/storage/content_addressable_storage:filetocas",
        "//intrinsic/storage/content_addressable_storage/testing/go:contentaddressablestoragetesting",
        "//intrinsic/tools/inctl/util:casgeometryuploader",
        "//intrinsic/util/testing:testio",
    ],
)

go_library(
    name = "gzffile",
    srcs = ["gzffile.go"],
    cdeps = [":gzf_uploader_c"],
    cgo = 1,
    visibility = [
        "//intrinsic/resources/catalog:__subpackages__",
        "//intrinsic/skills/tools/resource:__subpackages__",
    ],
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "//intrinsic/util/status:status_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_protobuf//proto",
    ],
)

sdf_scene_object(
    name = "dice_scene_object",
    src = "//intrinsic/models/assets/bluebird:dice/dice.sdf",
)

go_test(
    name = "gzffile_test",
    srcs = ["gzffile_test.go"],
    data = [
        ":dice_scene_object",
    ],
    library = ":gzffile",
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "//intrinsic/util/testing:testio",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
    ],
)

cc_library(
    name = "gzf_uploader_c",
    srcs = ["gzf_uploader_c.cc"],
    hdrs = ["gzf_uploader_c.h"],
    deps = [
        "//intrinsic/geometry/api/storage:geometry_library",
        "//intrinsic/geometry/api/storage:gfile_storage",
        "//intrinsic/geometry/api/storage:gzf_storage",
        "//intrinsic/scene/proto/v1:entity_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/scene/util:scene_object_gzf",
        "//intrinsic/util/status:status_cc_proto",
        "//intrinsic/util/status:status_conversion_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/component:geometry_component",
        "//intrinsic/world/gzfile",
        "@com_google_absl//absl/base:core_headers",
        "@com_google_absl//absl/status:statusor",
        "@com_google_absl//absl/strings",
    ],
)
