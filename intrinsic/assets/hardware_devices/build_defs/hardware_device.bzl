# Copyright 2023 Intrinsic Innovation LLC

"""Build rules for HardwareDevice assets."""

load("//intrinsic/assets/build_defs:asset.bzl", "AssetCatalogRefInfo", "AssetInfo", "AssetLocalInfo")

HardwareDeviceAssetInfo = provider(
    doc = "Provided by intrinsic_hardware_device() rule.",
    fields = ["bundle_tar"],
)

def _intrinsic_hardware_device_impl(ctx):
    asset_bundles = [
        a[AssetLocalInfo].bundle_path
        for a in ctx.attr.assets
        if AssetLocalInfo in a
    ]
    local_assets = [
        a[AssetLocalInfo].local_info
        for a in ctx.attr.assets
        if AssetLocalInfo in a
    ]
    catalog_assets = [
        a[AssetCatalogRefInfo].catalog_info
        for a in ctx.attr.assets
        if AssetCatalogRefInfo in a
    ]

    args = ctx.actions.args().add(
        "--manifest",
        ctx.file.manifest,
    ).add_all(
        local_assets,
        format_each = "--local_asset=%s",
    ).add_all(
        catalog_assets,
        format_each = "--catalog_asset=%s",
    ).add(
        "--output_bundle",
        ctx.outputs.bundle_out,
    )

    ctx.actions.run(
        inputs = asset_bundles + local_assets + catalog_assets + [ctx.file.manifest],
        outputs = [ctx.outputs.bundle_out],
        executable = ctx.executable._hardwaredevicegen,
        arguments = [args],
        mnemonic = "HardwareDeviceBundle",
        progress_message = "HardwareDevice bundle %s" % ctx.outputs.bundle_out.short_path,
    )

    transitive_descriptor_sets = depset(transitive = [
        f[AssetInfo].transitive_descriptor_sets
        for f in ctx.attr.assets
    ])

    transitive_inputs = [transitive_descriptor_sets]
    asset_info_output = ctx.actions.declare_file(ctx.label.name + ".asset_info.binpb")
    asset_local_info_output = ctx.actions.declare_file(ctx.label.name + ".asset_local_info.binpb")
    local_info_args = ctx.actions.args().add(
        "--manifest",
        ctx.file.manifest,
    ).add(
        "--asset_type",
        "ASSET_TYPE_HARDWARE_DEVICE",
    ).add(
        "--bundle_path",
        ctx.outputs.bundle_out,
    ).add(
        "--bundle_short_path",
        ctx.outputs.bundle_out.short_path,
    ).add_all(
        transitive_descriptor_sets,
        before_each = "--file_descriptor_set",
        # Since we're aggregating multiple different assets here, we may see the
        # same file more than once.  Assume these are consistent.
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
        mnemonic = "Assetlocalinfo",
        arguments = [local_info_args],
    )

    return [
        DefaultInfo(
            executable = ctx.outputs.bundle_out,
        ),
        HardwareDeviceAssetInfo(
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

intrinsic_hardware_device = rule(
    implementation = _intrinsic_hardware_device_impl,
    attrs = {
        "manifest": attr.label(
            allow_single_file = [".textproto"],
            mandatory = True,
            doc = "A manifest that provides the HardwareDevice definition.",
        ),
        "assets": attr.label_list(
            providers = [
                [AssetInfo, AssetLocalInfo],
                [AssetInfo, AssetCatalogRefInfo],
            ],
            doc = """Assets to add to the HardwareDeviceManifest saved in the bundle. These assets
                  must not already be listed in the manifest.""",
        ),
        "_hardwaredevicegen": attr.label(
            default = Label("//intrinsic/assets/hardware_devices/build_defs:hardwaredevicegen_main"),
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
    provides = [HardwareDeviceAssetInfo, AssetInfo, AssetLocalInfo],
    doc = "Bundles a HardwareDevice asset into a tar file.",
)
