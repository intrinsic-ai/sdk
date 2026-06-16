# Copyright 2023 Intrinsic Innovation LLC

"""Implements gen_http_bridge Bazel rule."""

load("@io_bazel_rules_go//go:def.bzl", "GoInfo")

def _gen_http_bridge_impl(ctx):
    output_path = ctx.actions.declare_file(ctx.label.name + "/main.go")

    openapi_file = ctx.file.openapi_path

    args = ctx.actions.args()
    args.add("httpjson", "generatemain")
    args.add("--openapi_path", openapi_file.path)
    args.add("--output", output_path)

    for go_proto_target, grpc_service in ctx.attr.services.items():
        go_info = go_proto_target[GoInfo]
        service_go_importpath = go_info.importpath
        args.add("--http_service", "%s:%s" % (grpc_service, service_go_importpath))

    ctx.actions.run(
        arguments = [args],
        executable = ctx.executable._inbuild,
        inputs = [openapi_file],
        mnemonic = "GenHttpBridge",
        outputs = [output_path],
        progress_message = "Generating HTTP bridge files for %s" % ctx.label.name,
    )

    return [
        DefaultInfo(files = depset([output_path])),
    ]

# gen_http_bridge calls `inbuild httpjson generatemain` for the intrinsic_http_image Bazel macro.
gen_http_bridge = rule(
    attrs = {
        "services": attr.label_keyed_string_dict(
            doc = "A map from go_proto_library targets to their gRPC service FQNs.",
            mandatory = True,
            providers = [GoInfo],
        ),
        "openapi_path": attr.label(
            allow_single_file = True,
            doc = "The path to an OpenAPI file whose content will be inserted into main.go, and which will be returned from the /openapi.yaml endpoint.",
            mandatory = True,
        ),
        "_inbuild": attr.label(
            cfg = "exec",
            default = Label("//intrinsic/tools/inbuild:inbuild"),
            doc = "The inbuild executable.",
            executable = True,
        ),
    },
    doc = "Generate files name/main.go needed to build an HTTP bridge service asset",
    implementation = _gen_http_bridge_impl,
)
