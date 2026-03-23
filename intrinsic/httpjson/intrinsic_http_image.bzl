# Copyright 2023 Intrinsic Innovation LLC

"""Implements intrinsic_http_service Bazel macro."""

load("@bazel_lib//lib:copy_file.bzl", "copy_file")
load("@rules_oci//oci:defs.bzl", "oci_image", "oci_load")
load("@rules_pkg//:mappings.bzl", "pkg_attributes", "pkg_files", "strip_prefix")
load("@rules_pkg//:pkg.bzl", "pkg_tar")
load("//bazel:go_macros.bzl", "go_binary")
load("//intrinsic/httpjson/openapi:protoc_gen_openapi.bzl", "protoc_gen_openapi")
load("//intrinsic/httpjson/private:gen_http_bridge.bzl", _gen_http_bridge = "gen_http_bridge")

def intrinsic_http_image(
        name,
        grpc_service,
        proto,
        go_proto):
    """Generate a Service Asset that offers HTTP/JSON endpoints for another Service Asset.

    Args:
        name: A name for the HTTP Bridge to be generated
        grpc_service: The fully qualified name of an annotated gRPC service
        proto: The proto_library target of the gRPC service
        go_proto: The go_grpc_http_library of the gRPC service
    """

    gen_name = "_" + name + "_generate"
    openapi_name = "_" + name + "_openapi"
    gobin_name = "_" + name + "_gobin"
    binfiles_name = gobin_name + "_files"
    tarbin_name = "_" + name + "_tarbin"
    ociimage_name = "_" + name + "_ociimage"
    ocitarball_name = "_" + name + "_tarball"
    ocitar_name = ocitarball_name + ".tar"

    protoc_gen_openapi(
        name = openapi_name,
        protos = [proto],
    )

    # Generate main.go using `inbuild httpservice generate`
    _gen_http_bridge(
        name = gen_name,
        grpc_service = grpc_service,
        openapi_path = ":" + openapi_name,
        service_go_proto = go_proto,
    )

    go_binary(
        name = gobin_name,
        srcs = [":" + gen_name],
        embedsrcs = [":" + openapi_name],
        deps = [
            go_proto,
            Label("//intrinsic/httpjson/openapi:handlers"),
            Label("//intrinsic/httpjson/any:anyresolver"),
            Label("//intrinsic/resources/proto:runtime_context_go_proto"),
            Label("@org_golang_google_grpc//credentials/insecure"),
            Label("//intrinsic/util/proto:protoio"),
            Label("@org_golang_google_grpc//:grpc"),
            Label("@com_github_grpc_ecosystem_grpc_gateway_v2//runtime"),
            Label("@org_golang_google_protobuf//encoding/protojson:go_default_library"),
        ],
    )

    pkg_files(
        name = binfiles_name,
        srcs = [":" + gobin_name],
        attributes = pkg_attributes(mode = "0555"),  # all: Read + Execute
        include_runfiles = True,
        prefix = "/opt/intrinsic",
        strip_prefix = strip_prefix.from_pkg(),
    )

    pkg_tar(
        name = tarbin_name,
        srcs = [":" + binfiles_name],
        extension = "tar.gz",
    )

    oci_image(
        name = ociimage_name,
        base = Label("@distroless_base"),
        entrypoint = ["/opt/intrinsic/" + gobin_name],
        tars = [":" + tarbin_name],
    )

    oci_load(
        name = ocitarball_name,
        image = ":" + ociimage_name,
        repo_tags = [ocitarball_name + ":latest"],
    )

    native.filegroup(
        name = ocitar_name,
        srcs = [":" + ocitarball_name],
        output_group = "tarball",
    )

    # Must rename file because intrinsic_service() only looks at an image's basename.
    copy_file(
        name = name,
        src = ":" + ocitar_name,
        out = name + ".tar",
        allow_symlink = True,
    )
