# Copyright 2023 Intrinsic Innovation LLC

# For libraries and services covering the installation of assets.  Initial
# design was detailed in go/intrinsic-asset-install.

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_python//python:defs.bzl", "py_library", "py_test")

py_library(
    name = "installed_assets_client",
    srcs = ["installed_assets_client.py"],
    visibility = [
        "//intrinsic:internal_api_users",
    ],
    deps = [
        "//intrinsic/assets:id_utils_py",
        "//intrinsic/assets/proto:asset_tag_py_pb2",
        "//intrinsic/assets/proto:asset_type_py_pb2",
        "//intrinsic/assets/proto:id_py_pb2",
        "//intrinsic/assets/proto:installed_assets_pb2_grpc",
        "//intrinsic/assets/proto:installed_assets_py_pb2",
        "//intrinsic/assets/proto:view_py_pb2",
        "//intrinsic/util/grpc:error_handling",
        "@ai_intrinsic_sdks_pip_deps//grpcio_status",
        requirement("grpcio"),
        "@com_google_googleapis//google/longrunning:operations_py_proto",
        "@com_google_googleapis//google/rpc:status_py_proto",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "installed_assets_client_test",
    srcs = ["installed_assets_client_test.py"],
    deps = [
        ":installed_assets_client",
        "//intrinsic/assets/processes/proto:process_asset_py_pb2",
        "//intrinsic/assets/proto:asset_tag_py_pb2",
        "//intrinsic/assets/proto:asset_type_py_pb2",
        "//intrinsic/assets/proto:id_py_pb2",
        "//intrinsic/assets/proto:installed_assets_py_pb2",
        "//intrinsic/assets/proto:metadata_py_pb2",
        "//intrinsic/assets/proto:view_py_pb2",
        "//intrinsic/executive/proto:behavior_tree_py_pb2",
        "//intrinsic/solutions/testing:compare",
        requirement("grpcio"),
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
        "@com_google_googleapis//google/longrunning:operations_py_proto",
        "@com_google_googleapis//google/rpc:code_py_proto",
        "@com_google_googleapis//google/rpc:status_py_proto",
        "@com_google_protobuf//:protobuf_python",
    ],
)
