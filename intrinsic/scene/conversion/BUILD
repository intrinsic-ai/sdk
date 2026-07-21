# Copyright 2023 Intrinsic Innovation LLC

# Libraries for conversion to Scene Object data structures.

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")
load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "object_properties_conversion",
    srcs = ["object_properties_conversion.cc"],
    hdrs = ["object_properties_conversion.h"],
    visibility = [
        "//intrinsic/scene/sdf:__subpackages__",
        "//intrinsic/scene/util:__subpackages__",
        "//intrinsic/scene/validate:__subpackages__",
        "//intrinsic/world:__subpackages__",
    ],
    deps = [
        "//intrinsic/icon/proto:eigen_conversion",
        "//intrinsic/kinematics/types:cartesian_limits",
        "//intrinsic/scene/proto/v1:object_properties_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_test(
    name = "object_properties_conversion_test",
    srcs = ["object_properties_conversion_test.cc"],
    deps = [
        ":object_properties_conversion",
        "//intrinsic/icon/utils:cart_state_test_utils",
        "//intrinsic/kinematics/types:cartesian_limits",
        "//intrinsic/scene/proto/v1:object_properties_cc_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
    ],
)

cc_library(
    name = "scene_object_model_utils",
    srcs = ["scene_object_model_utils.cc"],
    hdrs = ["scene_object_model_utils.h"],
    deps = [
        "//intrinsic/scene/proto/v1:entity_cc_proto",
        "//intrinsic/world/proto:kinematics_component_cc_proto",
        "@abseil-cpp//absl/container:flat_hash_set",
    ],
)

cc_library(
    name = "user_data",
    srcs = ["user_data.cc"],
    hdrs = ["user_data.h"],
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/util/proto:descriptors",
        "//intrinsic/util/proto:dynamic_message_parser",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "user_data_test",
    srcs = ["user_data_test.cc"],
    deps = [
        ":user_data",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_updates_cc_proto",
        "//intrinsic/util/proto:descriptors",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "world_to_scene_object_updates",
    srcs = ["world_to_scene_object_updates.cc"],
    hdrs = ["world_to_scene_object_updates.h"],
    deps = [
        "//intrinsic/math:pose3",
        "//intrinsic/math:proto_conversion",
        "//intrinsic/scene/proto/v1:scene_object_updates_cc_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world:geometry_types",
        "//intrinsic/world/proto:object_world_refs_cc_proto",
        "//intrinsic/world/proto:object_world_service_cc_proto",
        "//intrinsic/world/proto:object_world_updates_cc_proto",
        "@abseil-cpp//absl/algorithm:container",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "world_to_scene_object_updates_test",
    srcs = ["world_to_scene_object_updates_test.cc"],
    deps = [
        ":world_to_scene_object_updates",
        "//intrinsic/eigenmath:matchers",
        "//intrinsic/scene/proto/v1:scene_object_updates_cc_proto",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//intrinsic/world/proto:object_world_service_cc_proto",
        "//intrinsic/world/proto:object_world_updates_cc_proto",
        "@abseil-cpp//absl/status",
    ],
)

cc_test(
    name = "world_to_scene_object_updates_fuzz_test",
    timeout = "long",
    srcs = ["world_to_scene_object_updates_fuzz_test.cc"],
    deps = [
        ":world_to_scene_object_updates",
        "//intrinsic/config/proto:application_cc_proto",
        "//intrinsic/geometry/proto:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/storage:dummy_storage",
        "//intrinsic/math/proto:matrix_cc_proto",
        "//intrinsic/resources/proto:geometric_resource_data_cc_proto",
        "//intrinsic/resources/proto:resource_type_cc_proto",
        "//intrinsic/scene/config:scene_object_config_updater",
        "//intrinsic/scene/proto/v1:entity_cc_proto",
        "//intrinsic/scene/proto/v1:object_properties_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/scene/testing:scene_object_fuzz_domains",
        "//intrinsic/scene/testing:scene_object_pretty_printer",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/testing:gtest_wrapper",
        "//intrinsic/world:cc_proto",
        "//intrinsic/world/proto:attachment_component_cc_proto",
        "//intrinsic/world/proto:collections_component_cc_proto",
        "//intrinsic/world/proto:collision_component_cc_proto",
        "//intrinsic/world/proto:entity_search_cc_proto",
        "//intrinsic/world/proto:object_world_refs_cc_proto",
        "//intrinsic/world/proto:object_world_service_cc_grpc_proto",
        "//intrinsic/world/proto:object_world_service_cc_proto",
        "//intrinsic/world/proto:object_world_updates_cc_proto",
        "//intrinsic/world/proto:world_entity_cc_proto",
        "//intrinsic/world/proto:world_fragment_cc_proto",
        "//intrinsic/world/service:object_world_service_test_server",
        "//intrinsic/world/test:world_update_fuzz_domains",
        "@abseil-cpp//absl/base:log_severity",
        "@abseil-cpp//absl/log:globals",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_fuzztest//fuzztest",
        "@com_google_fuzztest//fuzztest:googletest_fixture_adapter",
        "@com_google_fuzztest//fuzztest:init_fuzztest",
        "@com_google_googleapis//google/rpc:code_cc_proto",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "world_to_scene_object_updates_c",
    srcs = ["world_to_scene_object_updates_c.cc"],
    hdrs = ["world_to_scene_object_updates_c.h"],
    deps = [
        ":world_to_scene_object_updates",
        "//intrinsic/scene/proto/v1:scene_object_updates_cc_proto",
        "//intrinsic/util/cgo:c_types",
        "//intrinsic/util/cgo:call_cpp",
        "//intrinsic/world/proto:object_world_service_cc_proto",
        "//intrinsic/world/proto:object_world_updates_cc_proto",
    ],
)

go_library(
    name = "world_to_scene_object_updates_go",
    srcs = ["go/world_to_scene_object_updates.go"],
    cdeps = [":world_to_scene_object_updates_c"],
    cgo = 1,
    importpath = "intrinsic/scene/conversion/worldtosceneobjectupdates",
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_updates_go_proto",
        "//intrinsic/util/cgo:goc",
        "//intrinsic/util/cgo:handle",
        "//intrinsic/world/proto:object_world_service_go_proto",
        "//intrinsic/world/proto:object_world_updates_go_proto",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//proto",
    ],
)

go_test(
    name = "world_to_scene_object_updates_go_test",
    srcs = ["go/world_to_scene_object_updates_test.go"],
    embed = [":world_to_scene_object_updates_go"],
    deps = [
        "//intrinsic/icon/proto:cart_space_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_updates_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
        "//intrinsic/world/proto:object_world_service_go_proto",
        "//intrinsic/world/proto:object_world_updates_go_proto",
        "//intrinsic/world/proto:simulation_component_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
    ],
)
