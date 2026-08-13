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

load("//bazel:container.bzl", "container_image")
load("//bazel:go_macros.bzl", "go_binary", "go_library", "go_test")
load("//intrinsic/assets/services/build_defs:services.bzl", "intrinsic_service")

go_library(
    name = "volumemount",
    srcs = ["volumemount.go"],
    importpath = "intrinsic/assets/services/examples/volume_mount/volumemount",
    deps = [
        "//intrinsic/assets/services/examples/volume_mount/proto/v1:volume_mount_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "volumemount_test",
    srcs = ["volumemount_test.go"],
    embed = [":volumemount"],
    importpath = "intrinsic/assets/services/examples/volume_mount/volumemount_test",
    deps = [
        "//intrinsic/assets/services/examples/volume_mount/proto/v1:volume_mount_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_binary(
    name = "volumemount_main",
    srcs = ["volumemount_main.go"],
    importpath = "intrinsic/assets/services/examples/volume_mount/volumemount_main",
    deps = [
        ":volumemount",
        "//intrinsic/assets/services:config",
        "//intrinsic/assets/services/examples/volume_mount/proto/v1:volume_mount_go_proto",
        "//intrinsic/production:intrinsic",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

container_image(
    name = "volume_mount_image",
    base = "@distroless_base",
    cmd = [
        "intrinsic/assets/services/examples/volume_mount/volumemount_main",
    ],
    data_path = "/",
    files = [
        ":volumemount_main",
    ],
)

intrinsic_service(
    name = "volume_mount_service",
    default_config = ":volume_mount_default_config.textproto",
    images = [
        ":volume_mount_image.tar",
    ],
    manifest = "volume_mount_manifest.textproto",
    deps = [
        "//intrinsic/assets/services/examples/volume_mount/proto/v1:volume_mount_proto",
    ],
)
