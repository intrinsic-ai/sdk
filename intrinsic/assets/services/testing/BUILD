# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "utils",
    testonly = True,
    srcs = ["utils.go"],
    importpath = "intrinsic/assets/services/testing/utils",
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = [
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/assets/services/proto/v1:service_inspection_go_proto",
        "//intrinsic/kubernetes/workcell_spec/proto:image_go_proto",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
