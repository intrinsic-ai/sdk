# Copyright 2023 Intrinsic Innovation LLC

"""Bazel rules for Data assets."""

load("@bazel_skylib//rules:common_settings.bzl", "BuildSettingInfo")
load("//intrinsic/assets/build_defs:asset.bzl", "AssetInfo", "AssetLocalInfo")

DataAssetInfo = provider(
    "Provided by intrinsic_data() rule.",
    fields = ["bundle_tar"],
)

def _intrinsic_data_impl(ctx):
    inputs = []
    transitive_inputs = []

    inputs.append(ctx.file.manifest)
    inputs.extend(ctx.files.data)

    transitive_descriptor_sets = depset(transitive = [
        f[ProtoInfo].transitive_descriptor_sets
        for f in ctx.attr.deps
    ])
    transitive_inputs.append(transitive_descriptor_sets)

    args = ctx.actions.args().add(
        "--manifest",
        ctx.file.manifest,
    ).add(
        "--output_bundle",
        ctx.outputs.bundle_out,
    ).add_all(
        transitive_descriptor_sets,
        before_each = "--file_descriptor_set",
        uniquify = True,
    ).add_all(
        ctx.files.data,
        before_each = "--expected_referenced_file",
    )

    ctx.actions.run(
        inputs = depset(inputs, transitive = transitive_inputs),
        outputs = [ctx.outputs.bundle_out],
        executable = ctx.executable._datagen,
        arguments = [args],
        mnemonic = "Databundle",
        progress_message = "Data bundle %s" % ctx.outputs.bundle_out.short_path,
    )

    asset_info_output = ctx.actions.declare_file(ctx.label.name + ".asset_info.binpb")
    asset_local_info_output = ctx.actions.declare_file(ctx.label.name + ".asset_local_info.binpb")
    local_info_args = ctx.actions.args().add(
        "--manifest",
        ctx.file.manifest,
    ).add(
        "--asset_type",
        "ASSET_TYPE_DATA",
    ).add(
        "--bundle_path",
        ctx.outputs.bundle_out,
    ).add(
        "--bundle_short_path",
        ctx.outputs.bundle_out.short_path,
    ).add_all(
        transitive_descriptor_sets,
        before_each = "--file_descriptor_set",
        uniquify = True,
    ).add(
        "--output_asset_info",
        asset_info_output,
    ).add(
        "--output_asset_local_info",
        asset_local_info_output,
    )
    ctx.actions.run(
        inputs = depset([ctx.file.manifest], transitive = transitive_inputs),
        outputs = [asset_info_output, asset_local_info_output],
        executable = ctx.executable._assetlocalinfogen,
        arguments = [local_info_args],
        mnemonic = "AssetLocalInfo",
        progress_message = "Writing asset local info %{output} for %{label}",
    )

    return [
        DefaultInfo(
            executable = ctx.outputs.bundle_out,
        ),
        DataAssetInfo(
            bundle_tar = ctx.outputs.bundle_out,
        ),
        AssetInfo(
            asset_info = asset_info_output,
            transitive_descriptor_sets = transitive_descriptor_sets,
        ),
        AssetLocalInfo(
            bundle_path = ctx.outputs.bundle_out,
            local_info = asset_local_info_output,
        ),
    ]

intrinsic_data = rule(
    implementation = _intrinsic_data_impl,
    attrs = {
        "manifest": attr.label(
            allow_single_file = [".textproto"],
            mandatory = True,
            doc = "A manifest that provides the data payload and metadata.",
        ),
        "data": attr.label_list(
            allow_empty = True,
            allow_files = True,
            mandatory = False,
            doc = "Data files that are referenced via ReferencedData in the data payload.",
        ),
        "deps": attr.label_list(
            mandatory = True,
            providers = [ProtoInfo],
            doc = "Proto dependencies needed to construct the data payload's FileDescriptorSet.",
        ),
        "_datagen": attr.label(
            default = Label("//intrinsic/assets/data/build_defs:datagen_main"),
            cfg = "exec",
            executable = True,
        ),
        "_assetlocalinfogen": attr.label(
            default = Label("//intrinsic/assets/build_defs:assetlocalinfogen"),
            cfg = "exec",
            executable = True,
        ),
    },
    outputs = {
        "bundle_out": "%{name}.bundle.tar",
    },
    provides = [DataAssetInfo, AssetInfo, AssetLocalInfo],
    doc = "Bundles a Data asset into a tar file.",
)
