# Copyright 2023 Intrinsic Innovation LLC

"""Bazel rule to generate an OpenAPI spec from proto_library targets."""

load("@com_google_protobuf//bazel/common:proto_info.bzl", "ProtoInfo")

def _make_proto_path_arg(proto_path):
    # Create a --proto_path arg for use with with args.add_all(...map_each=)
    # Don't use before_each because it adds leading whitespace that causes an error like.
    # Could not map to virtual file: bazel-out/haswell-fastbuild/bin/external/protobuf+/src/google/protobuf/_virtual_imports/any_proto: Input file is a directory.
    return "--proto_path=" + proto_path

def _protoc_gen_openapi_impl(ctx):
    """Implementation for the protoc_gen_openapi rule."""

    # Create a directory so multiple targets in one BUILD file don't conflict
    output_file = ctx.actions.declare_file("_%s/openapi.yaml" % ctx.attr.name)

    # Collect ProtoInfo from the user-provided targets
    all_proto_infos = [p[ProtoInfo] for p in ctx.attr.protos]

    # Gather all transitive .proto files. These are the inputs to the action.
    transitive_sources = depset(
        # direct = ctx.files._well_known_protos,
        transitive = [
            info.transitive_sources
            for info in all_proto_infos
        ],
    )

    # Gather all proto source roots for the -I/--proto_path flags.
    transitive_proto_paths = depset(transitive = [
        info.transitive_proto_path
        for info in all_proto_infos
    ])

    # The direct sources from the user's targets are the files passed directly
    # to the protoc command line. Transitive dependencies are found via the import paths.
    direct_sources = []
    for proto in ctx.attr.protos:
        direct_sources.extend(proto[ProtoInfo].direct_sources)

    # Use an args object to build the command line arguments for protoc.
    args = ctx.actions.args()

    # Add the plugin command, specifying the plugin executable's path.
    args.add("--plugin=protoc-gen-openapi=" + ctx.executable._plugin.path)

    # Add the output flag for the openapi plugin. This specifies the output directory
    args.add("--openapi_out=" + output_file.dirname)

    args.add("--openapi_opt=naming=json")
    args.add("--openapi_opt=fq_schema_naming=True")
    args.add("--openapi_opt=enum_type=string")

    # Add all the necessary import paths.
    args.add_all(transitive_proto_paths, map_each = _make_proto_path_arg)

    # Add the .proto files to be processed.
    args.add_all(direct_sources)

    # Define the build action that runs protoc.
    ctx.actions.run(
        arguments = [args],
        executable = ctx.executable._protoc,
        # All transitive .proto files are inputs to the action.
        inputs = transitive_sources,
        mnemonic = "ProtocGenOpenAPI",
        outputs = [output_file],
        progress_message = "Generating OpenAPI spec from %d protos" % len(direct_sources),
        # The protoc executable and the openapi plugin are the tools.
        tools = [ctx.executable._protoc, ctx.executable._plugin],
    )

    return [
        DefaultInfo(
            files = depset([output_file]),
        ),
    ]

protoc_gen_openapi = rule(
    attrs = {
        "protos": attr.label_list(
            doc = "A list of proto_library targets to generate the OpenAPI spec from.",
            mandatory = True,
            providers = [ProtoInfo],
        ),
        "_plugin": attr.label(
            cfg = "exec",
            # The label for the protoc-gen-openapi executable from gnostic.
            default = Label("@com_github_google_gnostic//cmd/protoc-gen-openapi:protoc-gen-openapi"),
            executable = True,
        ),
        "_protoc": attr.label(
            cfg = "exec",
            default = Label("@com_google_protobuf//:protoc"),
            executable = True,
        ),
    },
    doc = "Generates an OpenAPI v3 specification from a set of proto_library targets.",
    implementation = _protoc_gen_openapi_impl,
)
