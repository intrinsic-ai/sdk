# Copyright 2023 Intrinsic Innovation LLC

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

package(default_visibility = ["//visibility:public"])

go_binary(
    name = "eni_datagen",
    embed = [":build_defs_lib"],
)

filegroup(
    name = "bzl",
    srcs = glob(["*.bzl"]),
)

go_library(
    name = "build_defs_lib",
    srcs = ["eni_datagen.go"],
    deps = [
        "//intrinsic/assets/data/proto/v1:data_manifest_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/icon/fieldbus/ethercat/device_service/v1:eni_go_proto",
        "@org_golang_google_protobuf//encoding/prototext",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
