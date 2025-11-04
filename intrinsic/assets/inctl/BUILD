# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "assetcmd",
    srcs = ["assetcmd.go"],
    deps = [
        ":getreleased",
        ":install",
        ":list",
        ":listreleased",
        ":listreleasedversions",
        ":release",
        ":uninstall",
        ":updatereleasemetadata",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:cobrautil",
    ],
)

go_library(
    name = "assetviews",
    srcs = ["assetviews.go"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/proto:metadata_go_proto",
        "@org_golang_google_protobuf//proto",
    ],
)

go_test(
    name = "assetviews_test",
    srcs = ["assetviews_test.go"],
    library = ":assetviews",
    deps = [
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_grpc_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:installed_assets_go_grpc_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//proto",
    ],
)

go_library(
    name = "getreleased",
    srcs = ["getreleased.go"],
    deps = [
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_grpc_proto",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_library(
    name = "install",
    srcs = ["install.go"],
    deps = [
        "//intrinsic/assets:bundleio",
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:imagetransfer",
        "//intrinsic/assets:imageutils",
        "//intrinsic/assets/proto:installed_assets_go_grpc_proto",
        "//intrinsic/assets/services:bundleimages",
        "//intrinsic/kubernetes/acl:identity",
        "//intrinsic/skills/tools/skill/cmd/directupload",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_library(
    name = "list",
    srcs = ["list.go"],
    deps = [
        ":assetviews",
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets/proto:installed_assets_go_grpc_proto",
        "//intrinsic/assets/proto/v1:search_go_proto",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_library(
    name = "listreleased",
    srcs = ["listreleased.go"],
    deps = [
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:listutils",
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_grpc_proto",
        "//intrinsic/assets/proto:view_go_proto",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_spf13_cobra//:go_default_library",
        "@org_golang_google_protobuf//proto",
    ],
)

go_library(
    name = "listreleasedversions",
    srcs = ["listreleasedversions.go"],
    deps = [
        ":assetviews",
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets:listutils",
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_grpc_proto",
        "//intrinsic/assets/proto:view_go_proto",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@org_golang_google_protobuf//proto",
    ],
)

go_library(
    name = "release",
    srcs = ["release.go"],
    deps = [
        "//intrinsic/assets:bundleio",
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:imagetransfer",
        "//intrinsic/assets:imageutils",
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_grpc_proto",
        "//intrinsic/assets/catalog/proto/v1:release_metadata_go_proto",
        "//intrinsic/assets/services:bundleimages",
        "//intrinsic/skills/tools/skill/cmd/directupload",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_google_go_containerregistry//pkg/v1/remote:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_library(
    name = "uninstall",
    srcs = ["uninstall.go"],
    deps = [
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/proto:installed_assets_go_grpc_proto",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_library(
    name = "updatereleasemetadata",
    srcs = ["updatereleasemetadata.go"],
    deps = [
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_grpc_proto",
        "//intrinsic/assets/catalog/proto/v1:release_metadata_go_proto",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_spf13_cobra//:go_default_library",
        "@org_golang_google_protobuf//types/known/fieldmaskpb",
    ],
)
