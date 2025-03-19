# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

package(default_visibility = [
    "//intrinsic/assets:__subpackages__",
    "//intrinsic/frontend:__subpackages__",
    "//intrinsic/kubernetes:__subpackages__",
    "//intrinsic/skills/tools/skill:__subpackages__",
    "//intrinsic/solutions:__subpackages__",
    "//intrinsic/tools/inctl:__subpackages__",
])

proto_library(
    name = "clusterdeletion_api_proto",
    srcs = ["clusterdeletion_api.proto"],
    deps = [
        "@com_google_googleapis//google/api:field_behavior_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

go_proto_library(
    name = "clusterdeletion_api_go_proto",
    go_deps = ["@org_golang_google_genproto_googleapis_api//annotations"],
    deps = [":clusterdeletion_api_proto"],
)

go_grpc_library(
    name = "clusterdeletion_api_go_grpc_proto",
    srcs = [":clusterdeletion_api_proto"],
    deps = [
        ":clusterdeletion_api_go_proto",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)
