# Copyright 2023 Intrinsic Innovation LLC

# Common code for assets.

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_python//python:defs.bzl", "py_library")
load("//bazel:go_macros.bzl", "go_library", "go_test")

go_library(
    name = "baseclientutils",
    srcs = ["baseclientutils.go"],
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        "@com_github_pkg_errors//:go_default_library",
        "@io_opencensus_go//plugin/ocgrpc:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials:go_default_library",
        "@org_golang_google_grpc//credentials/insecure:go_default_library",
    ],
)

go_library(
    name = "clientutils",
    srcs = ["clientutils.go"],
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        ":baseclientutils",
        ":cmdutils",
        "//intrinsic/frontend/cloud/api/v1:clusterdiscovery_api_go_grpc_proto",
        "//intrinsic/frontend/cloud/api/v1:solutiondiscovery_api_go_grpc_proto",
        "//intrinsic/kubernetes/acl:identity",
        "//intrinsic/tools/inctl/auth",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_google_go_containerregistry//pkg/authn:go_default_library",
        "@com_github_google_go_containerregistry//pkg/v1/google:go_default_library",
        "@com_github_google_go_containerregistry//pkg/v1/remote:go_default_library",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials:go_default_library",
        "@org_golang_google_grpc//credentials/insecure:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)

go_library(
    name = "cmdutils",
    srcs = ["cmdutils.go"],
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        ":imagetransfer",
        ":imageutils",
        ":typeutils",
        ":viewutils",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:installed_assets_go_grpc_proto",
        "//intrinsic/assets/proto:view_go_proto",
        "//intrinsic/skills/tools/resource/cmd:bundleimages",
        "//intrinsic/tools/inctl/util:orgutil",
        "@com_github_google_go_containerregistry//pkg/authn:go_default_library",
        "@com_github_google_go_containerregistry//pkg/v1/google:go_default_library",
        "@com_github_google_go_containerregistry//pkg/v1/remote:go_default_library",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
        "@org_golang_x_exp//maps",
    ],
)

cc_library(
    name = "id_utils",
    srcs = ["id_utils.cc"],
    hdrs = ["id_utils.h"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto:id_cc_proto",
        "//intrinsic/util/status:annotate",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@boringssl//:crypto",
        "@com_googlesource_code_re2//:re2",
    ],
)

py_library(
    name = "id_utils_py",
    srcs = ["id_utils.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/assets/proto:id_py_pb2"],
)

go_library(
    name = "idutils",
    srcs = ["idutils.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto:id_go_proto",
        "@org_golang_x_exp//slices",
    ],
)

go_library(
    name = "imagetransfer",
    srcs = ["imagetransfer.go"],
    visibility = [
        "//intrinsic:internal_api_users",
        "//intrinsic:public_api_users",
    ],
    deps = [
        "@com_github_cenkalti_backoff_v4//:go_default_library",
        "@com_github_google_go_containerregistry//pkg/name:go_default_library",
        "@com_github_google_go_containerregistry//pkg/v1:go_default_library",
        "@com_github_google_go_containerregistry//pkg/v1/remote:go_default_library",
        "@com_github_google_go_containerregistry//pkg/v1/remote/transport:go_default_library",
        "@com_github_pkg_errors//:go_default_library",
    ],
)

go_library(
    name = "imageutils",
    srcs = ["imageutils.go"],
    visibility = ["//intrinsic:public_api_users"],
)

go_library(
    name = "tagutils",
    srcs = ["tagutils.go"],
    embedsrcs = [
        "//intrinsic/assets/proto:asset_tags",  # keep
    ],
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        "//intrinsic/assets/proto:asset_tag_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "@org_golang_google_protobuf//encoding/prototext:go_default_library",
    ],
)

go_test(
    name = "tagutils_test",
    srcs = ["tagutils_test.go"],
    library = ":tagutils",
    deps = [
        "//intrinsic/assets/proto:asset_tag_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
    ],
)

go_library(
    name = "typeutils",
    srcs = ["typeutils.go"],
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        "//intrinsic/assets/proto:asset_type_go_proto",
        "@org_golang_x_exp//slices",
    ],
)

go_library(
    name = "viewutils",
    srcs = ["viewutils.go"],
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/assets/proto:view_go_proto",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
    ],
)

go_library(
    name = "bundleio",
    srcs = [
        "bundle_io.go",
        "bundle_io_data.go",
        "bundle_io_hardware_devices.go",
        "bundle_io_services.go",
        "bundle_io_skills.go",
    ],
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        ":idutils",
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_grpc_proto",
        "//intrinsic/assets/data:utils",
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/data/proto/v1:referenced_data_go_proto",
        "//intrinsic/assets/hardware_devices:hardwaredevicemanifest",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/scene_objects:sceneobjectmanifest",
        "//intrinsic/assets/services:servicemanifest",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/kubernetes/workcell_spec/proto:image_go_proto",
        "//intrinsic/skills/internal:skillmanifest",
        "//intrinsic/skills/proto:processed_skill_manifest_go_proto",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util/archive:tartooling",
        "//intrinsic/util/go:pointer",
        "//intrinsic/util/proto:registryutil",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_google_safearchive//tar",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)

go_library(
    name = "version",
    srcs = ["version.go"],
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        ":idutils",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:installed_assets_go_grpc_proto",
        "@org_golang_google_protobuf//proto",
    ],
)

go_library(
    name = "listutils",
    srcs = ["listutils.go"],
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_grpc_proto",
        "//intrinsic/assets/proto:view_go_proto",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_library(
    name = "metadatautils",
    srcs = ["metadatautils.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":idutils",
        ":tagutils",
        "//intrinsic/assets/data/proto/v1:data_manifest_go_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/processes/proto:process_manifest_go_proto",
        "//intrinsic/assets/proto:asset_tag_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//proto",
    ],
)

go_test(
    name = "metadatautils_test",
    srcs = ["metadatautils_test.go"],
    library = ":metadatautils",
    deps = [
        "//intrinsic/assets/data/proto/v1:data_manifest_go_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/processes/proto:process_manifest_go_proto",
        "//intrinsic/assets/proto:asset_tag_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)
