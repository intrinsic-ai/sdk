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

"""Build definitions for content addressable storage service protos."""

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

proto_library(
    name = "cas_service_proto",
    srcs = ["cas_service.proto"],
)

go_proto_library(
    name = "cas_service_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/storage/content_addressable_storage/proto/cas_service_go_proto",
    protos = [":cas_service_proto"],
    visibility = [
        "//intrinsic/assets/artifacts:__pkg__",
        "//intrinsic/assets/catalog:__pkg__",
        "//intrinsic/assets/catalog/data:__pkg__",
        "//intrinsic/assets/catalog/testing/test_assets:__pkg__",
        "//intrinsic/assets/data:__pkg__",
        "//intrinsic/assets/hardware_devices/inctl:__pkg__",
        "//intrinsic/assets/inctl:__pkg__",
        "//intrinsic/assets/install:__pkg__",
        "//intrinsic/assets/scene_objects:__pkg__",
        "//intrinsic/assets/scene_objects/inctl:__pkg__",
        "//intrinsic/config:__pkg__",
        "//intrinsic/datasilo/access/examples/journalservice:__pkg__",
        "//intrinsic/frontend/cloud:__pkg__",
        "//intrinsic/frontend/go/api:__pkg__",
        "//intrinsic/frontend/go/cas:__pkg__",
        "//intrinsic/kubernetes/data_store:__pkg__",
        "//intrinsic/kubernetes/data_store/firestore_shim:__pkg__",
        "//intrinsic/kubernetes/workcell_spec:__pkg__",
        "//intrinsic/modified_solution:__pkg__",
        "//intrinsic/modified_solution/copier:__subpackages__",
        "//intrinsic/resources/catalog:__pkg__",
        "//intrinsic/skills/tools/resource/cmd:__pkg__",
        "//intrinsic/solution_copy/handlers:__pkg__",
        "//intrinsic/solution_versions:__pkg__",
        "//intrinsic/storage/content_addressable_storage:__pkg__",
        "//intrinsic/storage/content_addressable_storage/cmd:__pkg__",
        "//intrinsic/storage/content_addressable_storage/internal:__pkg__",
        "//intrinsic/storage/content_addressable_storage/internal/service:__pkg__",
        "//intrinsic/storage/content_addressable_storage/pkg/client:__pkg__",
        "//intrinsic/storage/content_addressable_storage/pkg/clienthelpers:__pkg__",
        "//intrinsic/storage/content_addressable_storage/pkg/crossproject:__pkg__",
        "//intrinsic/storage/content_addressable_storage/pkg/filetocas:__pkg__",
        "//intrinsic/storage/content_addressable_storage/pkg/lamprey:__pkg__",
        "//intrinsic/storage/content_addressable_storage/pkg/onpremproxy:__pkg__",
        "//intrinsic/storage/content_addressable_storage/testing/bin:__pkg__",
        "//intrinsic/storage/content_addressable_storage/testing/go:__pkg__",
        "//intrinsic/storage/download_data_files:__pkg__",
        "//intrinsic/storage/download_data_files/downloader:__pkg__",
        "//intrinsic/storage/hot_shared_state:__pkg__",
        "//intrinsic/storage/hot_shared_state/testing/go:__pkg__",
        "//intrinsic/storage/onprem_cas:__subpackages__",
        "//intrinsic/storage/onprem_cas_preload:__pkg__",
        "//intrinsic/storage/onprem_cas_preload/cmd:__pkg__",
        "//intrinsic/storage/onprem_cas_preload/internal/service:__pkg__",
        "//intrinsic/storage/onprem_cas_preload/pkg/testing:__pkg__",
        "//intrinsic/storage/polymorphic_data_frontend:__pkg__",
        "//intrinsic/tools/inctl/cmd:__pkg__",
        "//intrinsic/tools/inctl/cmd/cas:__pkg__",
        "//intrinsic/tools/inctl/cmd/modified_solution:__pkg__",
        "//intrinsic/tools/inctl/cmd/solution_version:__pkg__",
        "//intrinsic/tools/inctl/util:__pkg__",
    ],
)

cc_proto_library(
    name = "cas_service_cc_proto",
    visibility = [
        # Please don't use __subpackages__, specify the exact __pkg__!
        "//intrinsic/logging/bag_packager:__pkg__",
        "//intrinsic/perception/service:__pkg__",
        "//intrinsic/storage/content_addressable_storage:__pkg__",
        "//intrinsic/storage/content_addressable_storage/cpp:__pkg__",
        "//intrinsic/storage/content_addressable_storage/testing:__pkg__",
    ],
    deps = [":cas_service_proto"],
)

cc_grpc_library(
    name = "cas_service_cc_grpc_proto",
    srcs = [":cas_service_proto"],
    generate_mocks = False,  # Please use FakeCAS instead.
    grpc_only = True,
    visibility = [
        # Please don't use __subpackages__, specify the exact __pkg__!
        "//intrinsic/geometry/service:__pkg__",
        "//intrinsic/geometry/storage:__pkg__",
        "//intrinsic/logging/bag_packager:__pkg__",
        "//intrinsic/perception/public:__pkg__",
        "//intrinsic/perception/service:__pkg__",
        "//intrinsic/storage/content_addressable_storage:__pkg__",
        "//intrinsic/storage/content_addressable_storage/cpp:__pkg__",
        "//intrinsic/storage/content_addressable_storage/testing:__pkg__",
        "//intrinsic/world/binary:__pkg__",
        "//intrinsic/intrinsic_motion_planning/service/debug:__pkg__",
        "//intrinsic/intrinsic_motion_planning/service/nonvolatile_cache:__pkg__",
    ],
    deps = [":cas_service_cc_proto"],
)

py_proto_library(
    name = "cas_service_py_pb2",
    visibility = [
        # Please don't use __subpackages__, specify the exact __pkg__!
        "//intrinsic/storage/content_addressable_storage/python:__pkg__",
        "//intrinsic/storage/content_addressable_storage/testing/python:__pkg__",
    ],
    deps = [":cas_service_proto"],
)

py_grpc_library(
    name = "cas_service_py_pb2_grpc",
    srcs = [":cas_service_proto"],
    grpc_library = requirement("grpcio"),
    visibility = [
        # Please don't use __subpackages__, specify the exact __pkg__!
        "//intrinsic:public_api_users",
        "//intrinsic/perception/tools/utils:__pkg__",
        "//intrinsic/storage/content_addressable_storage/python:__pkg__",
        "//intrinsic/storage/content_addressable_storage/testing/python:__pkg__",
    ],
    deps = [":cas_service_py_pb2"],
)

proto_library(
    name = "cas_testing_service_proto",
    testonly = True,
    srcs = ["cas_testing_service.proto"],
)

go_proto_library(
    name = "cas_testing_service_go_proto",
    testonly = True,
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/storage/content_addressable_storage/proto/cas_testing_service_go_proto",
    protos = [":cas_testing_service_proto"],
    visibility = [
        "//intrinsic/solution_copy/handlers:__pkg__",
        "//intrinsic/storage/content_addressable_storage/testing/bin:__pkg__",
        "//intrinsic/storage/content_addressable_storage/testing/go:__pkg__",
    ],
)

cc_proto_library(
    name = "cas_testing_service_cc_proto",
    testonly = True,
    visibility = [
        # Please don't use __subpackages__, specify the exact __pkg__!
        "//intrinsic/storage/content_addressable_storage/testing:__pkg__",
    ],
    deps = [":cas_testing_service_proto"],
)

cc_grpc_library(
    name = "cas_testing_service_cc_grpc_proto",
    testonly = True,
    srcs = [":cas_testing_service_proto"],
    generate_mocks = False,
    grpc_only = True,
    visibility = [
        # Please don't use __subpackages__, specify the exact __pkg__!
        "//intrinsic/storage/content_addressable_storage/testing:__pkg__",
    ],
    deps = [":cas_testing_service_cc_proto"],
)

py_proto_library(
    name = "cas_testing_service_py_pb2",
    testonly = True,
    visibility = [
        # Please don't use __subpackages__, specify the exact __pkg__!
        "//intrinsic/storage/content_addressable_storage/testing/python:__pkg__",
    ],
    deps = [":cas_testing_service_proto"],
)

py_grpc_library(
    name = "cas_testing_service_py_pb2_grpc",
    testonly = True,
    srcs = [":cas_testing_service_proto"],
    grpc_library = requirement("grpcio"),
    visibility = [
        # Please don't use __subpackages__, specify the exact __pkg__!
        "//intrinsic/storage/content_addressable_storage/testing/python:__pkg__",
    ],
    deps = [":cas_testing_service_py_pb2"],
)
