# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@rules_oci//oci:defs.bzl", "oci_image", "oci_load")
load("@rules_pkg//:mappings.bzl", "pkg_attributes", "pkg_files", "strip_prefix")
load("@rules_pkg//:pkg.bzl", "pkg_tar")
load("//bazel:go_macros.bzl", "go_binary", "go_library", "go_proto_library", "go_test")
load("//intrinsic/assets/services/build_defs:services.bzl", "intrinsic_service")
load("//intrinsic/httpjson:intrinsic_http_image.bzl", "intrinsic_http_image")
load(
    "//intrinsic/httpjson/openapi:protoc_gen_openapi.bzl",
    "protoc_gen_openapi",
)

proto_library(
    name = "inventory_service_proto",
    srcs = ["inventory_service.proto"],
    deps = [
        "@com_google_googleapis//google/api:annotations_proto",
    ],
)

go_proto_library(
    name = "inventory_service_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
        "//bazel:go_gen_grpc_gateway",
    ],
    importpath = "intrinsic/httpjson/test/inventory_service_go_proto",
    protos = [
        ":inventory_service_proto",
    ],
    deps = [
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)

protoc_gen_openapi(
    name = "inventory_service_openapi",
    protos = [
        ":inventory_service_proto",
    ],
)

go_test(
    name = "openapi_test",
    srcs = ["openapi_test.go"],
    data = [":inventory_service_openapi"],
    importpath = "intrinsic/httpjson/test/openapi",
    deps = [
        "@io_bazel_rules_go//go/runfiles",
    ],
)

go_library(
    name = "inventoryserver",
    srcs = ["inventoryserver.go"],
    importpath = "intrinsic/httpjson/test/inventoryserver",
    deps = [
        ":inventory_service_go_proto",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "inventoryserver_test",
    srcs = ["inventoryserver_test.go"],
    data = [":inventory_service_openapi"],
    importpath = "intrinsic/httpjson/test/inventoryserver_test",
    deps = [
        ":inventory_service_go_proto",
        ":inventoryserver",
        "//intrinsic/httpjson/openapi:handlers",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//runtime",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//test/bufconn:go_default_library",
        "@org_golang_google_protobuf//encoding/protojson:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_binary(
    name = "inventoryserver_main",
    srcs = ["inventoryserver_main.go"],
    importpath = "intrinsic/httpjson/test/inventoryserver_main",
    deps = [
        ":inventory_service_go_proto",
        ":inventoryserver",
        "//intrinsic/resources/proto:runtime_context_go_proto",
        "//intrinsic/util/proto:protoio",
        "@org_golang_google_grpc//:grpc",
    ],
)

# https://blog.engflow.com/2024/08/20/migrating-to-bazel-modules-aka-bzlmod---repo-names-and-rules_pkg/
pkg_files(
    name = "inventoryserver_files",
    srcs = [":inventoryserver_main"],
    attributes = pkg_attributes(mode = "0555"),  # all: Read + Execute
    prefix = "/opt/intrinsic",
    strip_prefix = strip_prefix.from_pkg(),
)

pkg_tar(
    name = "inventoryserver_layer",
    srcs = [":inventoryserver_files"],
    extension = "tar.gz",
)

oci_image(
    name = "inventoryserver_image",
    base = "@distroless_base",
    cmd = ["/opt/intrinsic/inventoryserver_main"],
    tars = [":inventoryserver_layer"],
)

oci_load(
    name = "inventoryserver_image_tarball",
    image = ":inventoryserver_image",
    repo_tags = ["inventoryserver_image_tarball:latest"],
)

filegroup(
    name = "inventoryserver_image.tar",
    srcs = [":inventoryserver_image_tarball"],
    output_group = "tarball",
)

intrinsic_http_image(
    name = "inventoryserver_http_image",
    go_proto = ":inventory_service_go_proto",
    grpc_service = "intrinsic_proto.httpjson.v1.InventoryService",
    proto = ":inventory_service_proto",
)

intrinsic_service(
    name = "inventory_service",
    images = [
        ":inventoryserver_image.tar",
        ":inventoryserver_http_image",
    ],
    manifest = "inventoryservice_manifest.textproto",
    deps = [":inventory_service_proto"],
)
