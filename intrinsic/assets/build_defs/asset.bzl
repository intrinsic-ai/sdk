# Copyright 2026 Intrinsic Innovation LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
        "asset": "The asset ID of the instance",
        "config": "InstanceConfig textproto file",
        "name": "The name of the asset instance",
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
        "deps": attr.label_list(
            doc = "Proto dependencies that are compatible with the catalog " +
                  "asset. These are optional but required to parse config " +
                  "files. Note that version skew or other errors may happen " +
                  "if the wrong protos are used.",
            providers = [ProtoInfo],
        ),
        "id": attr.string(
            mandatory = True,
        ),
        "type": attr.string(),
        "version": attr.string(
            mandatory = True,
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
    files = []
    config = None
    if ctx.attr.required_node_hostname or ctx.file.service_config:
        config = ctx.actions.declare_file(ctx.label.name + ".txtpb")
        files.append(config)
        args = ctx.actions.args().add(
            "--output",
            config,
        ).add(
            "--required_node_hostname",
            ctx.attr.required_node_hostname,
        )
        inputs = []
        if ctx.file.service_config:
            args.add(
                "--service_config",
                ctx.file.service_config,
            )
            inputs.append(ctx.file.service_config)

        ctx.actions.run(
            arguments = [args],
            executable = ctx.executable._assetinstancegen,
            inputs = inputs,
            mnemonic = "AssetInstance",
            outputs = [config],
            progress_message = "Writing %{output} for %{label}",
        )

    return [
        DefaultInfo(
            files = depset(files),
        ),
        AssetInstanceInfo(
            name = name,
            asset = ctx.attr.asset,
            config = config,
        ),
    ]

intrinsic_asset_instance = rule(
    attrs = {
        "asset": attr.string(
            mandatory = True,
        ),
        "instance_name": attr.string(
            doc = "Name of the instance, if it should be different than 'name'",
        ),
        "required_node_hostname": attr.string(
            mandatory = False,
        ),
        "service_config": attr.label(
            allow_single_file = [
                ".pbtxt",
                ".txtpb",
                ".textproto",
            ],
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
