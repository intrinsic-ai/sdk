# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:container.bzl", "container_image")
load("//bazel:go_macros.bzl", "go_binary", "go_library", "go_test")
load("//intrinsic/assets/services/build_defs:services.bzl", "intrinsic_service")

package(default_visibility = [
    "//visibility:public",
])

go_library(
    name = "deviceservice",
    srcs = [
        "device_service.go",
        "esi_handler.go",
    ],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/data/proto/v1:data_assets_go_proto",
        "//intrinsic/assets/dependencies:utils",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/icon/fieldbus/ethercat/device_service/v1:device_service_config_go_proto",
        "//intrinsic/icon/fieldbus/ethercat/device_service/v1:device_service_go_proto",
        "//intrinsic/icon/fieldbus/ethercat/device_service/v1:esi_go_proto",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_protobuf//proto",
    ],
)

go_binary(
    name = "device_service_main",
    srcs = ["device_service_main.go"],
    deps = [
        ":deviceservice",
        "//intrinsic/assets/data/proto/v1:data_assets_go_proto",
        "//intrinsic/assets/services/proto/v1:service_state_go_proto",
        "//intrinsic/icon/fieldbus/ethercat/device_service/v1:device_service_config_go_proto",
        "//intrinsic/icon/fieldbus/ethercat/device_service/v1:device_service_go_proto",
        "//intrinsic/production:intrinsic",
        "//intrinsic/resources/proto:runtime_context_go_proto",
        "//intrinsic/util/status:extended_status_go_proto",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//credentials/insecure:go_default_library",
        "@org_golang_google_grpc//reflection:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//proto",
    ],
)

container_image(
    name = "device_service_image",
    base = "@distroless_base",
    cmd = ["intrinsic/icon/fieldbus/ethercat/device_service/device_service_main"],
    data_path = "/",
    files = [":device_service_main"],
)

intrinsic_service(
    name = "device_service",
    default_config = ":device_service_default_config.textproto",
    images = ["//intrinsic/icon/fieldbus/ethercat/device_service:device_service_image.tar"],
    manifest = "device_service_manifest.textproto",
    deps = [
        "//intrinsic/assets/services/proto/v1:service_state_proto",
        "//intrinsic/icon/fieldbus/ethercat/device_service/v1:device_service_config_proto",
        "//intrinsic/icon/fieldbus/ethercat/device_service/v1:device_service_proto",
    ],
)

go_test(
    name = "deviceservice_test",
    srcs = ["device_service_test.go"],
    embed = [":deviceservice"],
    deps = [
        "//intrinsic/assets/data:fakedataassets",
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/data/proto/v1:data_assets_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_go_proto",
        "//intrinsic/icon/fieldbus/ethercat/device_service/v1:device_service_config_go_proto",
        "//intrinsic/icon/fieldbus/ethercat/device_service/v1:device_service_go_proto",
        "//intrinsic/icon/fieldbus/ethercat/device_service/v1:esi_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
