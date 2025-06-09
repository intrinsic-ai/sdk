# Copyright 2023 Intrinsic Innovation LLC

"""
Bazel rules for service types.
"""

load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "ProtoSourceCodeInfo", "gen_source_code_info_descriptor_set")

ServiceTypeInfo = provider(
    "provided by intrinsic_service() rule",
    fields = ["bundle_tar"],
)

def _intrinsic_service_impl(ctx):
    bundle_output = ctx.outputs.bundle_out

    basenames = {}
    for file in ctx.files.images:
        if file.basename in basenames:
            # This is a requirement based on how we place the files into the tar
            # archive.  The files are placed into the root of the tar file
            # currently, so having ones with the same base name would cause them
            # to conflict or potentially silently overwrite.
            fail("Basenames of images must be unique; got multiple {}".format(file.basename))
        basenames[file.basename] = None

    transitive_descriptor_sets = depset(transitive = [
        f[ProtoSourceCodeInfo].transitive_descriptor_sets
        for f in ctx.attr.deps
    ])

    inputs = [ctx.file.manifest] + ctx.files.images
    transitive_inputs = [transitive_descriptor_sets]
    args = ctx.actions.args().add(
        "--manifest",
        ctx.file.manifest,
    ).add(
        "--output_bundle",
        bundle_output,
    ).add_joined(
        "--image_tars",
        ctx.files.images,
        join_with = ",",
    ).add_joined(
        "--file_descriptor_sets",
        transitive_descriptor_sets,
        join_with = ",",
    )
    if ctx.file.default_config:
        inputs.append(ctx.file.default_config)
        args.add("--default_config", ctx.file.default_config.path)

    ctx.actions.run(
        inputs = depset(inputs, transitive = transitive_inputs),
        outputs = [bundle_output],
        executable = ctx.executable._servicegen,
        arguments = [args],
        mnemonic = "Servicebundle",
        progress_message = "Creating service bundle %{output} for %{label}",
    )

    return [
        DefaultInfo(
            executable = bundle_output,
        ),
        ServiceTypeInfo(
            bundle_tar = bundle_output,
        ),
    ]

intrinsic_service = rule(
    implementation = _intrinsic_service_impl,
    attrs = {
        "default_config": attr.label(
            allow_single_file = [".pbtxt", ".textproto"],
            doc = """The path to the default configuration text proto for the service. If
            unspecified, the default configuration will be an empty message of the type specified in
            the manifest's ServiceDef.config_message_full_name.""",
        ),
        "images": attr.label_list(
            allow_empty = True,
            allow_files = [".tar"],
            doc = "Image tarballs referenced by the service type.",
        ),
        "manifest": attr.label(
            allow_single_file = [".textproto"],
            mandatory = True,
            doc = (
                "A manifest that can be used to provide the service definition and metadata."
            ),
        ),
        "deps": attr.label_list(
            providers = [ProtoInfo],
            aspects = [gen_source_code_info_descriptor_set],
        ),
        "_servicegen": attr.label(
            default = Label("//intrinsic/assets/services/build_defs:servicegen_main"),
            cfg = "exec",
            executable = True,
        ),
    },
    outputs = {
        "bundle_out": "%{name}.bundle.tar",
    },
)
