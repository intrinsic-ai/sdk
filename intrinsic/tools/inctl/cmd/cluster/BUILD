# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//intrinsic/tools/inctl:__subpackages__"])

go_library(
    name = "cluster",
    srcs = [
        "cluster.go",
        "cluster_delete.go",
        "cluster_list.go",
        "cluster_mode.go",
        "cluster_poweroff.go",
        "cluster_upgrade.go",
        "cluster_ups.go",
    ],
    visibility = [
        "//intrinsic/tools/inctl:__subpackages__",
    ],
    deps = [
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/frontend/cloud/api/v1:clusterdiscovery_api_go_proto",
        "//intrinsic/frontend/cloud/api/v1:clustermanager_api_go_proto",
        "//intrinsic/frontend/cloud/api/v1alpha1:clustermanager_api_go_proto",
        "//intrinsic/frontend/cloud/devicemanager:version",
        "//intrinsic/kubernetes/inversion/v1:inversion_go_proto",
        "//intrinsic/kubernetes/workcellmode/proto:workcellmode_service_go_proto",
        "//intrinsic/skills/tools/skill/cmd:dialerutil",
        "//intrinsic/tools/inctl/auth",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:cobrautil",
        "//intrinsic/tools/inctl/util:orgutil",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//credentials/insecure:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//types/known/fieldmaskpb",
    ],
)
