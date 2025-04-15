# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_library")

package(default_visibility = [
    "//visibility:public",
])

cc_library(
    name = "icon_api_service",
    hdrs = ["icon_api_service.h"],
    deps = [
        "//intrinsic/icon/proto/v1:service_cc_grpc_proto",
    ],
)

cc_library(
    name = "operational_state_interface",
    hdrs = ["operational_state_interface.h"],
    deps = [
        "//intrinsic/icon/cc_client:operational_status",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
    ],
)
