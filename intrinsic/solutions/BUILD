# Copyright 2023 Intrinsic Innovation LLC

# Solution building Python libraries.

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_python//python:defs.bzl", "py_library", "py_test")
load("@rules_python//python:packaging.bzl", "py_package", "py_wheel")

package(default_visibility = [
    "//visibility:public",
])

py_library(
    name = "behavior_tree",
    srcs = ["behavior_tree.py"],
    deps = [
        ":blackboard_value",
        ":cel",
        ":errors",
        ":ipython",
        ":proto_building",
        ":utils",
        "//intrinsic/executive/proto:any_list_py_pb2",
        "//intrinsic/executive/proto:any_with_assignments_py_pb2",
        "//intrinsic/executive/proto:behavior_call_py_pb2",
        "//intrinsic/executive/proto:behavior_tree_py_pb2",
        "//intrinsic/executive/proto:code_execution_py_pb2",
        "//intrinsic/executive/proto:world_query_py_pb2",
        "//intrinsic/skills/proto:skills_py_pb2",
        "//intrinsic/solutions/internal:actions",
        "//intrinsic/solutions/internal:skill_generation",
        "//intrinsic/solutions/internal:skill_utils",
        "//intrinsic/util/status:extended_status_py_pb2",
        "//intrinsic/world/proto:object_world_refs_py_pb2",
        "//intrinsic/world/python:object_world_resources",
        "@ai_intrinsic_sdks_pip_deps//graphviz",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "behavior_tree_test",
    srcs = ["behavior_tree_test.py"],
    data = [
        "//intrinsic/solutions/testing:test_skill_params_proto_descriptors",
    ],
    deps = [
        ":behavior_tree",
        ":blackboard_value",
        ":cel",
        ":errors",
        ":proto_building",
        "//intrinsic/executive/proto:any_with_assignments_py_pb2",
        "//intrinsic/executive/proto:behavior_tree_py_pb2",
        "//intrinsic/executive/proto:proto_builder_py_pb2",
        "//intrinsic/executive/proto:test_message_py_pb2",
        "//intrinsic/executive/proto:world_query_py_pb2",
        "//intrinsic/solutions/internal:behavior_call",
        "//intrinsic/solutions/internal:skill_providing",
        "//intrinsic/solutions/testing:compare",
        "//intrinsic/solutions/testing:skill_test_utils",
        "//intrinsic/solutions/testing:test_skill_params_py_pb2",
        "//intrinsic/util/path_resolver:path_resolver_py",
        "//intrinsic/util/status:extended_status_py_pb2",
        "//intrinsic/world/proto:object_world_refs_py_pb2",
        "//intrinsic/world/proto:object_world_service_py_pb2",
        "//intrinsic/world/python:object_world_resources",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_library(
    name = "blackboard_value",
    srcs = ["blackboard_value.py"],
    deps = [
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_library(
    name = "cel",
    srcs = ["cel.py"],
)

py_library(
    name = "deployments",
    srcs = ["deployments.py"],
    deps = [
        ":auth",
        ":dialerutil",
        ":error_processing",
        ":errors",
        ":execution",
        ":ipython",
        ":pbt_registration",
        ":pose_estimation",
        ":proto_building",
        ":providers",
        ":simulation",
        ":userconfig",
        ":worlds",
        "//intrinsic/assets/proto:installed_assets_pb2_grpc",
        "//intrinsic/frontend/cloud/api/v1:solutiondiscovery_api_py_pb2",
        "//intrinsic/frontend/cloud/api/v1:solutiondiscovery_api_py_pb2_grpc",
        "//intrinsic/frontend/solution_service/proto:solution_service_py_pb2",
        "//intrinsic/frontend/solution_service/proto:solution_service_py_pb2_grpc",
        "//intrinsic/frontend/solution_service/proto:status_py_pb2",
        "//intrinsic/kubernetes/workcell_spec/proto:installer_py_pb2_grpc",
        "//intrinsic/resources/client:resource_registry_client",
        "//intrinsic/scene/product/client:product_client",
        "//intrinsic/skills/client:skill_registry_client",
        "//intrinsic/solutions/internal:behavior_tree_providing",
        "//intrinsic/solutions/internal:products",
        "//intrinsic/solutions/internal:resources",
        "//intrinsic/solutions/internal:skill_providing",
        "//intrinsic/solutions/internal:stubs",
        requirement("grpcio"),
    ],
)

py_test(
    name = "deployments_test",
    srcs = ["deployments_test.py"],
    deps = [
        ":auth",
        ":deployments",
        ":dialerutil",
        ":error_processing",
        ":errors",
        ":execution",
        ":simulation",
        ":userconfig",
        ":worlds",
        "//intrinsic/frontend/solution_service/proto:solution_service_py_pb2",
        "//intrinsic/frontend/solution_service/proto:status_py_pb2",
        "//intrinsic/skills/client:skill_registry_client",
        "//intrinsic/skills/proto:skill_registry_py_pb2",
        requirement("grpcio"),
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_library(
    name = "userconfig",
    srcs = ["userconfig.py"],
)

py_test(
    name = "userconfig_test",
    srcs = ["userconfig_test.py"],
    deps = [
        ":userconfig",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

py_library(
    name = "errors",
    srcs = ["errors.py"],
    deps = [
        "@ai_intrinsic_sdks_pip_deps//retrying",
    ],
)

py_library(
    name = "error_processing",
    srcs = ["error_processing.py"],
    deps = [
        ":errors",
        ":ipython",
        "//intrinsic/logging/errors/proto:error_report_py_pb2",
        requirement("grpcio"),  # buildcleaner: keep
        "@com_google_googleapis//google/longrunning:operations_py_proto",
        "@com_google_googleapis//google/rpc:code_py_proto",
    ],
)

py_test(
    name = "error_processing_test",
    srcs = ["error_processing_test.py"],
    deps = [
        ":error_processing",
        "//intrinsic/logging/errors/proto:error_report_py_pb2",
        "//intrinsic/solutions/testing:compare",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_googleapis//google/longrunning:operations_py_proto",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_library(
    name = "execution",
    srcs = ["execution.py"],
    deps = [
        ":behavior_tree",
        ":blackboard_value",
        ":error_processing",
        ":errors",
        ":ipython",
        ":provided",
        ":simulation",
        ":utils",
        "//intrinsic/executive/proto:behavior_tree_py_pb2",
        "//intrinsic/executive/proto:blackboard_service_py_pb2",
        "//intrinsic/executive/proto:blackboard_service_py_pb2_grpc",
        "//intrinsic/executive/proto:executive_execution_mode_py_pb2",
        "//intrinsic/executive/proto:executive_service_py_pb2",
        "//intrinsic/executive/proto:executive_service_py_pb2_grpc",
        "//intrinsic/executive/proto:run_metadata_py_pb2",
        "//intrinsic/executive/proto:run_response_py_pb2",
        "//intrinsic/solutions/internal:actions",
        "//intrinsic/util/grpc:error_handling",
        "//intrinsic/util/proto:descriptors_py",
        "//intrinsic/util/status:extended_status_py_pb2",
        "//intrinsic/util/status:status_exception",
        requirement("grpcio"),
        "@com_google_googleapis//google/longrunning:operations_py_proto",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "execution_test",
    srcs = ["execution_test.py"],
    deps = [
        ":behavior_tree",
        ":blackboard_value",
        ":error_processing",
        ":errors",
        ":execution",
        ":provided",
        ":simulation",
        "//intrinsic/executive/proto:behavior_tree_py_pb2",
        "//intrinsic/executive/proto:blackboard_service_py_pb2",
        "//intrinsic/executive/proto:blackboard_service_py_pb2_grpc",
        "//intrinsic/executive/proto:executive_execution_mode_py_pb2",
        "//intrinsic/executive/proto:executive_service_py_pb2",
        "//intrinsic/executive/proto:executive_service_py_pb2_grpc",
        "//intrinsic/executive/proto:run_metadata_py_pb2",
        "//intrinsic/executive/proto:run_response_py_pb2",
        "//intrinsic/logging/errors/proto:error_report_py_pb2",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/solutions/internal:behavior_call",
        "//intrinsic/solutions/testing:compare",
        "//intrinsic/solutions/testing:test_skill_params_py_pb2",
        requirement("grpcio"),
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
        "@com_google_googleapis//google/longrunning:operations_py_proto",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_library(
    name = "ipython",
    srcs = ["ipython.py"],
    deps = [
        "//intrinsic/util/status:status_exception",
        "@ai_intrinsic_sdks_pip_deps//ipywidgets",
    ],
)

py_test(
    name = "ipython_test",
    srcs = ["ipython_test.py"],
    deps = [
        ":ipython",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

# Collection of all "public" modules for packaging.
py_library(
    name = "solutions_lib",
    deps = [
        ":blackboard_value",
        ":cel",
        ":deployments",
        ":errors",
        ":execution",
        ":installed_assets",
        ":motion_planning",
        ":perception",
        ":pose_estimation",
        ":providers",
        ":scene_object_import",
        ":simulation",
        ":worlds",
    ],
)

py_library(
    name = "perception",
    srcs = ["perception.py"],
    deps = [
        ":deployments",
        ":execution",
        ":utils",
        "//intrinsic/perception/client/v1/python/camera:_camera_utils",
        "//intrinsic/perception/client/v1/python/camera:camera_client",
        "//intrinsic/perception/client/v1/python/camera:data_classes",
        "//intrinsic/perception/proto/v1:image_buffer_py_pb2",
        "//intrinsic/resources/client:resource_registry_client",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/util/grpc:connection",
        "@ai_intrinsic_sdks_pip_deps//matplotlib",
        requirement("grpcio"),
    ],
)

py_test(
    name = "perception_external_test",
    srcs = ["perception_external_test.py"],
    deps = [
        ":perception",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

py_library(
    name = "pose_estimation",
    srcs = ["pose_estimation.py"],
    deps = [
        ":ipython",
        "//intrinsic/assets/proto:asset_type_py_pb2",
        "//intrinsic/assets/proto:installed_assets_pb2_grpc",
        "//intrinsic/assets/proto:installed_assets_py_pb2",
        "//intrinsic/assets/proto:view_py_pb2",
        "//intrinsic/perception/proto/v1:perception_model_py_pb2",
        "//intrinsic/resources/client:resource_registry_client",
        "//intrinsic/util/grpc:error_handling",
        requirement("grpcio"),
    ],
)

py_test(
    name = "pose_estimation_test",
    srcs = ["pose_estimation_test.py"],
    deps = [
        ":pose_estimation",
        "//intrinsic/assets/data/proto/v1:data_asset_py_pb2",
        "//intrinsic/assets/proto:asset_type_py_pb2",
        "//intrinsic/assets/proto:id_py_pb2",
        "//intrinsic/assets/proto:installed_assets_py_pb2",
        "//intrinsic/assets/proto:metadata_py_pb2",
        "//intrinsic/perception/proto/v1:perception_model_py_pb2",
        "//intrinsic/resources/proto:resource_registry_py_pb2",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

py_library(
    name = "proto_building",
    srcs = ["proto_building.py"],
    deps = [
        ":errors",
        "//intrinsic/executive/proto:proto_builder_py_pb2",
        "//intrinsic/executive/proto:proto_builder_py_pb2_grpc",
        "//intrinsic/solutions/internal:skill_utils",
        "//intrinsic/util/grpc:error_handling",
        "//intrinsic/util/proto:descriptors",
        "//intrinsic/util/proto:descriptors_py",
        requirement("grpcio"),
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "proto_building_test",
    srcs = ["proto_building_test.py"],
    deps = [
        ":behavior_tree",
        ":errors",
        ":proto_building",
        "//intrinsic/executive/proto:proto_builder_py_pb2",
        "//intrinsic/solutions/testing:compare",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_library(
    name = "pbt_registration",
    srcs = ["pbt_registration.py"],
    deps = [
        ":behavior_tree",
        "//intrinsic/executive/proto:behavior_tree_py_pb2",
        "//intrinsic/skills/proto:behavior_tree_registry_py_pb2",
        "//intrinsic/skills/proto:behavior_tree_registry_py_pb2_grpc",
        "//intrinsic/skills/proto:skill_registry_config_py_pb2",
        "//intrinsic/util/grpc:error_handling",
        requirement("grpcio"),
    ],
)

py_library(
    name = "simulation",
    srcs = ["simulation.py"],
    deps = [
        ":errors",
        "//intrinsic/math/python:proto_conversion",
        "//intrinsic/simulation/service/proto/v1:simulation_service_py_pb2",
        "//intrinsic/simulation/service/proto/v1:simulation_service_py_pb2_grpc",
        "//intrinsic/util/grpc:error_handling",
        "//intrinsic/world/proto:object_world_service_py_pb2",
        "//intrinsic/world/proto:object_world_service_py_pb2_grpc",
        requirement("grpcio"),
    ],
)

py_test(
    name = "simulation_external_test",
    srcs = ["simulation_external_test.py"],
    deps = [
        ":simulation",
        "//intrinsic/simulation/service/proto/v1:simulation_service_py_pb2",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_library(
    name = "provided",
    srcs = ["provided.py"],
    deps = [
        ":blackboard_value",
        ":cel",
        ":utils",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/scene/proto/v1:scene_object_py_pb2",
        "//intrinsic/skills/proto:skills_py_pb2",
        "//intrinsic/solutions/internal:actions",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_library(
    name = "providers",
    srcs = ["providers.py"],
    deps = [
        ":behavior_tree",
        ":provided",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
    ],
)

py_library(
    name = "auth",
    srcs = ["auth.py"],
    deps = [
        ":userconfig",
    ],
)

py_library(
    name = "dialerutil",
    srcs = ["dialerutil.py"],
    deps = [
        ":auth",
        "//intrinsic/kubernetes/acl/py:identity",
        requirement("grpcio"),
    ],
)

py_test(
    name = "auth_test",
    srcs = ["auth_test.py"],
    deps = [
        ":auth",
        ":userconfig",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

py_test(
    name = "dialerutil_test",
    srcs = ["dialerutil_test.py"],
    deps = [
        ":auth",
        ":dialerutil",
        requirement("grpcio"),
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

py_library(
    name = "utils",
    srcs = ["utils.py"],
    deps = [":errors"],
)

py_test(
    name = "utils_test",
    srcs = ["utils_test.py"],
    deps = [
        ":utils",
        "//intrinsic/executive/proto:test_message_py_pb2",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

py_library(
    name = "scene_object_import",
    deps = [
        "//intrinsic/geometry/proto/v1:material_py_pb2",
        "//intrinsic/scene/proto/v1:scene_object_import_py_pb2",
        "//intrinsic/scene/proto/v1:scene_object_import_py_pb2_grpc",
        "//intrinsic/scene/proto/v1:scene_object_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//googleapis_common_protos",
        "@com_google_googleapis//google/longrunning:operations_py_proto",
    ],
)

py_library(
    name = "installed_assets",
    deps = [
        "//intrinsic/assets/proto:asset_type_py_pb2",
        "//intrinsic/assets/proto:id_py_pb2",
        "//intrinsic/assets/proto:installed_assets_pb2_grpc",
        "//intrinsic/assets/proto:installed_assets_py_pb2",
        "//intrinsic/assets/proto:metadata_py_pb2",
        "//intrinsic/assets/proto:view_py_pb2",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_py_proto",
        "@ai_intrinsic_sdks_pip_deps//googleapis_common_protos",
        "@com_google_googleapis//google/longrunning:operations_py_proto",
    ],
)

py_library(
    name = "worlds",
    srcs = ["worlds.py"],
    deps = [
        ":utils",
        "//intrinsic/geometry/proto:geometry_service_py_pb2_grpc",
        "//intrinsic/math/python:data_types",
        "//intrinsic/math/python:proto_conversion",
        "//intrinsic/motion_planning/proto:motion_target_py_pb2",
        "//intrinsic/world/proto:collision_settings_py_pb2",
        "//intrinsic/world/proto:object_world_refs_py_pb2",
        "//intrinsic/world/proto:object_world_service_py_pb2_grpc",
        "//intrinsic/world/python:object_world_client",
        "//intrinsic/world/python:object_world_resources",
        requirement("grpcio"),
    ],
)

py_library(
    name = "motion_planning",
    srcs = ["motion_planning.py"],
    deps = [
        ":deployments",
        "//intrinsic/motion_planning:motion_planner_client_py",
        "//intrinsic/motion_planning/proto/v1:motion_planner_service_py_pb2_grpc",
    ],
)

py_test(
    name = "motion_planning_external_test",
    srcs = ["motion_planning_external_test.py"],
    deps = [
        ":motion_planning",
        "//intrinsic/motion_planning/proto/v1:motion_planner_service_py_pb2_grpc",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "worlds_external_test",
    srcs = ["worlds_external_test.py"],
    deps = [
        ":worlds",
        "//intrinsic/world/proto:object_world_service_py_pb2",
        "//intrinsic/world/proto:object_world_service_py_pb2_grpc",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

py_package(
    name = "solutions_py_package",
    packages = [
        "intrinsic",
        "third_party",
    ],
    deps = [
        ":solutions_lib",
    ],
)

py_wheel(
    name = "solutions_wheel",
    distribution = "solutions",
    requires = [
        "absl-py",
        "black",
        "googleapis-common-protos",
        "graphviz",
        "grpcio",
        "ipywidgets",
        "matplotlib",
        "numpy",
        "retrying",
        "typing-extensions",
    ],
    version = "0.0.1",
    deps = [":solutions_py_package"],
)
