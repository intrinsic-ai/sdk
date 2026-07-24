# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "gzfprocessor",
    srcs = ["gzfprocessor.go"],
    importpath = "intrinsic/assets/scene_objects/gzfprocessor",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":gzffile",
        "//intrinsic/assets:referenceddata",
        "//intrinsic/assets/data/proto/v1:referenced_data_go_proto",
        "//intrinsic/geometry/proto:geometry_storage_refs_go_proto",
        "//intrinsic/geometry/proto/v1:geometry_go_proto",
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_go_proto",
        "//intrinsic/geometry/proto/v1:transformed_geometry_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "//intrinsic/world/proto:geometry_component_go_proto",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_x_sync//errgroup",
    ],
)

go_library(
    name = "gzffile",
    srcs = ["gzffile.go"],
    importpath = "intrinsic/assets/scene_objects/gzffile",
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "@org_golang_google_protobuf//proto",
    ],
)

go_library(
    name = "sceneobjectbundle",
    srcs = ["sceneobjectbundle.go"],
    importpath = "intrinsic/assets/scene_objects/sceneobjectbundle",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        ":gzfprocessor",
        ":sceneobjectvalidate",
        "//intrinsic/assets:ioutils",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "//intrinsic/util/archive:tartooling",
        "//intrinsic/util/go:pointer",
        "//intrinsic/util/proto:registryutil",
        "@com_github_google_safearchive//tar",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
        "@org_golang_google_protobuf//types/descriptorpb:go_default_library",
    ],
)

go_library(
    name = "sceneobjectfix",
    srcs = ["sceneobjectfix.go"],
    importpath = "intrinsic/assets/scene_objects/sceneobjectfix",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "@org_golang_google_protobuf//types/descriptorpb:go_default_library",
    ],
)

go_library(
    name = "sceneobjectvalidate",
    srcs = ["sceneobjectvalidate.go"],
    importpath = "intrinsic/assets/scene_objects/sceneobjectvalidate",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":gzffile",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:metadatautils",
        "//intrinsic/assets/errors:report",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "//intrinsic/scene/validate/go:sceneobjectvalidation",
        "//intrinsic/util/proto:registryutil",
        "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)
