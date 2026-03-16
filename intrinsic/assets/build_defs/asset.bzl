# Copyright 2023 Intrinsic Innovation LLC

"""Bazel rules for Assets."""

load("@com_google_protobuf//bazel/common:proto_info.bzl", "ProtoInfo")

AssetInfo = provider(
    "Info about an asset.",
    fields = {
        "asset_info": "An AssetInfo proto",
        "transitive_descriptor_sets": "depset of file descriptor sets",
    },
)

AssetLocalInfo = provider(
    "Info about a built asset bundle file.",
    fields = {
        "bundle_path": "The full path to the asset's bundle file",
        "local_info": "An AssetLocalInfo proto",
    },
)

AssetCatalogRefInfo = provider(
    "Info about an asset catalog reference.",
    fields = {
        "catalog_info": "An AssetCatalogRefInfo proto",
    },
)

AssetInstanceInfo = provider(
    "An asset instance.",
    fields = {
        "instance_info": "An AssetInstanceInfo proto",
        "config": "Optional any proto text file of the asset's configuration",
    },
)

def _intrinsic_asset_reference_impl(ctx):
    asset_info_output = ctx.actions.declare_file(ctx.label.name + ".asset_info.binpb")
    asset_catalog_ref_info_output = ctx.actions.declare_file(ctx.label.name + ".asset_catalog_ref_info.binpb")

    transitive_descriptor_sets = depset(transitive = [
        f[ProtoInfo].transitive_descriptor_sets
        for f in ctx.attr.deps
    ])

    args = ctx.actions.args().add(
        "--asset_type",
        ctx.attr.type,
    ).add(
        "--id",
        ctx.attr.id,
    ).add_all(
        transitive_descriptor_sets,
        before_each = "--file_descriptor_set",
    ).add(
        "--version",
        ctx.attr.version,
    ).add(
        "--output_asset_info",
        asset_info_output,
    ).add(
        "--output_asset_catalog_ref_info",
        asset_catalog_ref_info_output,
    )

    ctx.actions.run(
        arguments = [args],
        executable = ctx.executable._assetcatalogrefinfogen,
        inputs = transitive_descriptor_sets,
        mnemonic = "AssetReference",
        outputs = [asset_info_output, asset_catalog_ref_info_output],
        progress_message = "Writing %{output} for %{label}",
    )
    return [
        DefaultInfo(
            executable = asset_catalog_ref_info_output,
            files = depset([asset_catalog_ref_info_output]),
        ),
        AssetInfo(
            asset_info = asset_info_output,
            transitive_descriptor_sets = transitive_descriptor_sets,
        ),
        AssetCatalogRefInfo(
            catalog_info = asset_catalog_ref_info_output,
        ),
    ]

intrinsic_asset_reference = rule(
    attrs = {
        "type": attr.string(),
        "id": attr.string(
            mandatory = True,
        ),
        "version": attr.string(
            mandatory = True,
        ),
        "deps": attr.label_list(
            doc = "Proto dependencies that are compatible with the catalog " +
                  "asset. These are optional but required to parse config " +
                  "files. Note that version skew or other errors may happen " +
                  "if the wrong protos are used.",
            providers = [ProtoInfo],
        ),
        "_assetcatalogrefinfogen": attr.label(
            cfg = "exec",
            default = Label("//intrinsic/assets/build_defs:assetcatalogrefinfogen"),
            executable = True,
        ),
    },
    provides = [
        AssetInfo,
        AssetCatalogRefInfo,
    ],
    implementation = _intrinsic_asset_reference_impl,
)

def _intrinsic_asset_instance_impl(ctx):
    name = ctx.attr.instance_name if ctx.attr.instance_name else ctx.label.name
    asset_instance_output = ctx.actions.declare_file(ctx.label.name + ".binpb")
    args = ctx.actions.args().add(
        "--id",
        ctx.attr.id,
    ).add(
        "--instance_name",
        name,
    ).add(
        "--required_node_hostname",
        ctx.attr.required_node_hostname,
    ).add(
        "--output_asset_instance",
        asset_instance_output,
    )
    inputs = []
    transitive_inputs = []
    transitive_runfiles = []
    if ctx.file.config:
        args.add(
            "--config_path",
            ctx.file.config,
        )
        inputs.append(ctx.file.config)

    ctx.actions.run(
        arguments = [args],
        executable = ctx.executable._assetinstancegen,
        inputs = depset(inputs, transitive = transitive_inputs),
        mnemonic = "AssetInstance",
        outputs = [asset_instance_output],
        progress_message = "Writing %{output} for %{label}",
    )
    return [
        DefaultInfo(
            executable = asset_instance_output,
            files = depset([asset_instance_output]),
            runfiles = ctx.runfiles(
                transitive_files = depset(transitive = transitive_runfiles),
            ),
        ),
        AssetInstanceInfo(
            config = ctx.file.config,
            instance_info = asset_instance_output,
        ),
    ]

intrinsic_asset_instance = rule(
    attrs = {
        "id": attr.string(
            mandatory = True,
        ),
        "instance_name": attr.string(
            doc = "Name of the instance, if it should be different than 'name'",
        ),
        "config": attr.label(
            allow_single_file = [
                ".pbtxt",
                ".txtpb",
                ".textproto",
            ],
        ),
        "required_node_hostname": attr.string(
            mandatory = False,
        ),
        "_assetinstancegen": attr.label(
            cfg = "exec",
            default = Label("//intrinsic/assets/build_defs:assetinstancegen"),
            executable = True,
        ),
    },
    provides = [AssetInstanceInfo],
    implementation = _intrinsic_asset_instance_impl,
)
