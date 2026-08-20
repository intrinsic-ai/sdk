# Copyright 2026 Intrinsic Innovation LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@intrinsic_js_rules//:defs.bzl", "js_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "dio_config_proto",
    srcs = ["dio_config.proto"],
)

cc_proto_library(
    name = "dio_config_cc_proto",
    deps = [":dio_config_proto"],
)

go_proto_library(
    name = "dio_config_go_proto",
    importpath = "intrinsic/icon/server/config/dio_config_go_proto",
    protos = [":dio_config_proto"],
)

proto_library(
    name = "services_config_proto",
    srcs = ["services_config.proto"],
)

cc_proto_library(
    name = "services_config_cc_proto",
    deps = [":services_config_proto"],
)

go_proto_library(
    name = "services_config_go_proto",
    importpath = "intrinsic/icon/server/config/services_config_go_proto",
    protos = [":services_config_proto"],
)

proto_library(
    name = "icon_main_config_proto",
    srcs = ["icon_main_config.proto"],
    deps = [
        ":realtime_control_config_proto",
        ":services_config_proto",
        "//intrinsic/assets/proto:field_metadata_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_proto",
        "//intrinsic/icon/proto:logging_mode_proto",
        "@com_google_protobuf//:duration_proto",
    ],
)

cc_proto_library(
    name = "icon_main_config_cc_proto",
    deps = [":icon_main_config_proto"],
)

py_proto_library(
    name = "icon_main_config_py_pb2",
    deps = [":icon_main_config_proto"],
)

js_library(
    name = "icon_main_config_jspb_proto",
    deps = [":icon_main_config_proto"],
)

go_proto_library(
    name = "icon_main_config_go_proto",
    importpath = "intrinsic/icon/server/config/icon_main_config_go_proto",
    protos = [":icon_main_config_proto"],
    deps = [
        ":realtime_control_config_go_proto",
        ":services_config_go_proto",
        "//intrinsic/assets/proto:field_metadata_go_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_go_proto",
        "//intrinsic/icon/proto:logging_mode_go_proto",
    ],
)

proto_library(
    name = "realtime_control_config_proto",
    srcs = ["realtime_control_config.proto"],
    # Always depend on all default HAL part proto types to mainloop targets simpler.
    deps = [
        "//intrinsic/icon/control/parts/hal/adio_part:hal_adio_part_config_proto",  # buildcleaner: keep
        "//intrinsic/icon/control/parts/hal/arm_part:hal_arm_part_config_proto",  # buildcleaner: keep
        "//intrinsic/icon/control/parts/hal/force_torque_sensor_part:hal_force_torque_sensor_part_config_proto",  # buildcleaner: keep
        "//intrinsic/icon/control/parts/hal/imu_part:hal_imu_part_config_proto",  # buildcleaner: keep
        "//intrinsic/icon/control/parts/hal/laser_tracker_part:hal_laser_tracker_part_config_proto",  # buildcleaner: keep
        "//intrinsic/icon/control/parts/hal/linear_gripper_part:hal_linear_gripper_part_config_proto",  # buildcleaner: keep
        "//intrinsic/icon/control/parts/hal/rangefinder_part:hal_rangefinder_part_config_proto",  # buildcleaner: keep
        "//intrinsic/icon/control/parts/hal/v1:hal_part_config_proto",
        "//intrinsic/icon/control/parts/proto/v1:realtime_part_config_proto",
        "//intrinsic/icon/hal/lib/fieldbus/v1:value_parsing_proto",
    ],
)

cc_proto_library(
    name = "realtime_control_config_cc_proto",
    deps = [":realtime_control_config_proto"],
)

js_library(
    name = "realtime_control_config_jspb_proto",
    deps = [":realtime_control_config_proto"],
)

go_proto_library(
    name = "realtime_control_config_go_proto",
    importpath = "intrinsic/icon/server/config/realtime_control_config_go_proto",
    protos = [":realtime_control_config_proto"],
    deps = [
        "//intrinsic/icon/control/parts/hal/adio_part:hal_adio_part_config_go_proto",  # buildcleaner: keep
        "//intrinsic/icon/control/parts/hal/arm_part:hal_arm_part_config_go_proto",  # buildcleaner: keep
        "//intrinsic/icon/control/parts/hal/force_torque_sensor_part:hal_force_torque_sensor_part_config_go_proto",  # buildcleaner: keep
        "//intrinsic/icon/control/parts/hal/imu_part:hal_imu_part_config_go_proto",  # buildcleaner: keep
        "//intrinsic/icon/control/parts/hal/laser_tracker_part:hal_laser_tracker_part_config_go_proto",  # buildcleaner: keep
        "//intrinsic/icon/control/parts/hal/linear_gripper_part:hal_linear_gripper_part_config_go_proto",  # buildcleaner: keep
        "//intrinsic/icon/control/parts/hal/rangefinder_part:hal_rangefinder_part_config_go_proto",  # buildcleaner: keep
        "//intrinsic/icon/control/parts/hal/v1:hal_part_config_go_proto",
        "//intrinsic/icon/control/parts/proto/v1:realtime_part_config_go_proto",
        "//intrinsic/icon/hal/lib/fieldbus/v1:value_parsing_go_proto",
    ],
)
