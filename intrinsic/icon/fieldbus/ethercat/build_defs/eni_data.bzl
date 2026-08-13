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
        asset_name = asset_name,
        asset_package = asset_package,
        display_name = display_name,
        eni_file = eni_file,
        tool = eni_datagen_label,
        vendor_display_name = vendor_display_name,
    )

    # Genrule to create the DataManifest textproto.
    native.genrule(
        name = manifest_name,
        srcs = [eni_file],
        outs = [manifest_file],
        cmd = eni_datagen_cmd,
        tools = [eni_datagen_label],
        visibility = ["//visibility:private"],
    )

    data_deps = [
        Label("//intrinsic/assets/data/proto/v1:data_manifest_proto"),
        Label("//intrinsic/icon/fieldbus/ethercat/device_service/v1:eni_proto"),
    ]

    intrinsic_data(
        name = name,
        manifest = ":" + manifest_file,
        visibility = visibility,
        deps = data_deps,
    )
