# Copyright 2023 Intrinsic Innovation LLC

"""Build rule for creating Data assets from ENI files."""

load("//intrinsic/assets/data/build_defs:data.bzl", "intrinsic_data")

def intrinsic_eni_data(
        name,
        eni_file,
        asset_package,
        asset_name,
        display_name,
        vendor_display_name,
        visibility = None):
    """Creates an intrinsic_data asset from an ENI file.

    The content of the eni_file is embedded in an intrinsic_proto.fieldbus.ethercat.device_service.v1.Eni
    message, which is then packed into the DataManifest's 'data' Any field. See
    @@ai_intrinsic_sdks+/intrinsic/assets/data/proto/v1/data_manifest.proto for details.

    Args:
      name: The name of the intrinsic_data target to generate.
      eni_file: The label of the ENI file (e.g., ":my_config.eni").
      asset_package: The package of the Data asset ID (e.g., "ai.intrinsic.ethercat").
      asset_name: The name of the Data asset ID (e.g., "my_config_eni").
      display_name: The display name of the Data asset.
      vendor_display_name: The display name of the vendor.
      visibility: Visibility of the generated intrinsic_data target.
    """

    manifest_name = name + "_manifest"
    manifest_file = name + "_manifest.textproto"

    eni_datagen_label = Label("//intrinsic/icon/fieldbus/ethercat/build_defs:eni_datagen")

    eni_datagen_cmd = """
        $(location {tool}) \
            --eni_file=$(location {eni_file}) \
            --output_file=$@ \
            --asset_package=\"{asset_package}\" \
            --asset_name=\"{asset_name}\" \
            --display_name=\"{display_name}\" \
            --vendor_display_name=\"{vendor_display_name}\" \
    """.format(
        tool = eni_datagen_label,
        eni_file = eni_file,
        asset_package = asset_package,
        asset_name = asset_name,
        display_name = display_name,
        vendor_display_name = vendor_display_name,
    )

    # Genrule to create the DataManifest textproto.
    native.genrule(
        name = manifest_name,
        srcs = [eni_file],
        outs = [manifest_file],
        tools = [eni_datagen_label],
        cmd = eni_datagen_cmd,
        visibility = ["//visibility:private"],
    )

    data_deps = [
        Label("//intrinsic/assets/data/proto/v1:data_manifest_proto"),
        Label("//intrinsic/icon/fieldbus/ethercat/device_service/v1:eni_proto"),
    ]

    intrinsic_data(
        name = name,
        manifest = ":" + manifest_file,
        deps = data_deps,
        visibility = visibility,
    )
