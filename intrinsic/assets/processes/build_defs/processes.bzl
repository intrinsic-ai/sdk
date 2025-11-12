# Copyright 2023 Intrinsic Innovation LLC

"""Bazel rules for Process assets."""

load("//intrinsic/assets/build_defs:asset.bzl", "AssetInfo", "AssetLocalInfo")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "ProtoSourceCodeInfo", "gen_source_code_info_descriptor_set")

ProcessAssetInfo = provider(
    "Provided by the intrinsic_process() rule.",
    fields = ["bundle_tar"],
)

def _intrinsic_process_impl(ctx):
    # The file descriptor set of the Process asset as a binary proto. In contrast to other asset
    # types, this cannot be specified via the rule's dependencies/parameters but must be extracted
    # from the behavior tree proto, hence we generate it in 'processgen'. This file descriptor set
    # is optional and at loading/analysis time we can't tell whether it is present. We thus always
    # generate a descriptor set file in 'processgen' and pass it to 'assetlocalinfogen' below. The
    # file will be empty if the Process asset does not have a parameter file descriptor set.
    file_descriptor_set_file = ctx.actions.declare_file(ctx.label.name + ".file_descriptor_set.binpb")

    # The full manifest as a binary proto. We generate it as a *binary* proto in 'processgen' so
    # that it can be parsed later by 'assetlocalinfogen' without all the additional file descriptor
    # sets required to parse the input manifest textproto (see 'textproto_deps').
    manifest_binary_file = ctx.actions.declare_file(ctx.label.name + ".manifest.binpb")

    transitive_textproto_descriptor_sets = depset(transitive = [
        f[ProtoSourceCodeInfo].transitive_descriptor_sets
        for f in ctx.attr.textproto_deps
    ])

    processgen_args = ctx.actions.args().add(
        "--manifest",
        ctx.file.manifest,
    ).add_all(
        transitive_textproto_descriptor_sets,
        uniquify = True,
        before_each = "--textproto_file_descriptor_set",
    ).add(
        "--output_bundle",
        ctx.outputs.bundle_out,
    ).add(
        "--output_file_descriptor_set",
        file_descriptor_set_file,
    ).add(
        "--output_manifest_binary",
        manifest_binary_file,
    )

    processgen_inputs = [ctx.file.manifest]

    if ctx.file.behavior_tree:
        processgen_args.add("--behavior_tree", ctx.file.behavior_tree)
        processgen_inputs.append(ctx.file.behavior_tree)

    ctx.actions.run(
        inputs = depset(processgen_inputs, transitive = [transitive_textproto_descriptor_sets]),
        outputs = [ctx.outputs.bundle_out, file_descriptor_set_file, manifest_binary_file],
        executable = ctx.executable._processgen,
        arguments = [processgen_args],
        mnemonic = "Processbundle",
        progress_message = "Process bundle %s" % ctx.outputs.bundle_out.short_path,
    )

    asset_info_output = ctx.actions.declare_file(ctx.label.name + ".asset_info.binpb")
    asset_local_info_output = ctx.actions.declare_file(ctx.label.name + ".asset_local_info.binpb")
    assetlocalinfogen_args = ctx.actions.args().add(
        "--manifest",
        manifest_binary_file,
    ).add(
        "--asset_type",
        "ASSET_TYPE_PROCESS",
    ).add(
        "--bundle_path",
        ctx.outputs.bundle_out,
    ).add(
        "--bundle_short_path",
        ctx.outputs.bundle_out.short_path,
    ).add(
        # This file is empty if the Process asset does not have a parameter file descriptor set
        # (also see comment above).
        "--file_descriptor_set",
        file_descriptor_set_file,
    ).add(
        "--output_asset_info",
        asset_info_output,
    ).add(
        "--output_asset_local_info",
        asset_local_info_output,
    )
    ctx.actions.run(
        inputs = depset([ctx.file.manifest, file_descriptor_set_file, manifest_binary_file]),
        outputs = [asset_info_output, asset_local_info_output],
        executable = ctx.executable._assetlocalinfogen,
        arguments = [assetlocalinfogen_args],
        mnemonic = "AssetLocalInfo",
        progress_message = "Writing asset local info %{output} for %{label}",
    )

    return [
        DefaultInfo(
            executable = ctx.outputs.bundle_out,
        ),
        ProcessAssetInfo(
            bundle_tar = ctx.outputs.bundle_out,
        ),
        AssetInfo(
            asset_info = asset_info_output,
            transitive_descriptor_sets = depset([file_descriptor_set_file]),
        ),
        AssetLocalInfo(
            bundle_path = ctx.outputs.bundle_out,
            local_info = asset_local_info_output,
        ),
    ]

intrinsic_process = rule(
    implementation = _intrinsic_process_impl,
    doc = "Bundles a Process asset into a tar file.",
    attrs = {
        "manifest": attr.label(
            allow_single_file = [".txtpb"],
            mandatory = True,
            doc = "A ProcessManifest message that provides the metadata and process definition.",
        ),
        "behavior_tree": attr.label(
            allow_single_file = [".txtpb"],
            doc = "An optional BehaviorTree representing the process " +
                  "(if the manifest does not contain a behavior tree).",
        ),
        "textproto_deps": attr.label_list(
            providers = [ProtoInfo],
            aspects = [gen_source_code_info_descriptor_set],
            doc = "Optional proto dependencies for parsing expanded Any protos in the input " +
                  "textprotos. This attribute is for convenience only to allow for more readable " +
                  "input textprotos. These deps will NOT be included in the file descriptor set " +
                  "of the created Process asset. If there are no expanded Any protos in the " +
                  "input, you can safely omit this attribute.",
        ),
        "_processgen": attr.label(
            default = Label("//intrinsic/assets/processes/build_defs:processgen_main"),
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
    provides = [ProcessAssetInfo, AssetInfo, AssetLocalInfo],
)
