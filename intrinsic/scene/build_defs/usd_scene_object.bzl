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

"""Custom rules for creating SceneObject from USD."""

load("@com_google_protobuf//bazel/common:proto_info.bzl", "ProtoInfo")
load("//intrinsic/scene/build_defs:scene_object.bzl", "SceneObjectInfo")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "ProtoSourceCodeInfo", "gen_source_code_info_descriptor_set")

def _usd_scene_object_impl(ctx):
    # Determine output files up front.
    if ctx.attr.updates_pbtxts:
        # If updates are provided, then output the gzf that has updates applied to it.
        raw_gzf = ctx.actions.declare_file(ctx.label.name + "__raw.gzf")
        raw_pbtxt = ctx.actions.declare_file(ctx.label.name + "__raw.pbtxt")
        updated_gzf = ctx.actions.declare_file(ctx.label.name + ".gzf")
        updated_pbtxt = ctx.actions.declare_file(ctx.label.name + ".pbtxt")
        output_gzf = updated_gzf
        output_pbtxt = updated_pbtxt
    else:
        # If updates are not provided, then output the raw gzf.
        raw_gzf = ctx.actions.declare_file(ctx.label.name + ".gzf")
        raw_pbtxt = ctx.actions.declare_file(ctx.label.name + ".pbtxt")
        updated_gzf = None
        updated_pbtxt = None
        output_gzf = raw_gzf
        output_pbtxt = raw_pbtxt

    transitive_descriptor_sets = depset(transitive = [
        f[ProtoSourceCodeInfo].transitive_descriptor_sets
        for f in ctx.attr.deps
    ])

    ## Convert the USD to a SceneObject pbtxt.
    raw_gzf_args = ctx.actions.args()
    raw_gzf_args.add("--scene_object_name", ctx.label.name)
    raw_gzf_args.add("--input_usd_file", ctx.file.src.path)
    raw_gzf_args.add("--output_scene_object_gzf_file", raw_gzf.path)
    raw_gzf_args.add("--output_scene_object_pbtxt_file", raw_pbtxt.path)

    ctx.actions.run(
        arguments = [raw_gzf_args],
        executable = ctx.executable._usd_file_converter,
        inputs = [ctx.file.src],
        mnemonic = "UsdToSceneObject",
        outputs = [raw_gzf, raw_pbtxt],
        progress_message = "Converting %s to SceneObject" % ctx.file.src.short_path,
    )

    ## Apply scene object updates to the proto.
    if ctx.attr.updates_pbtxts:
        update_args = ctx.actions.args()
        update_args.add("--input_scene_object_gzf_file", raw_gzf.path)
        update_args.add_joined("--input_updates_proto_filenames", ctx.files.updates_pbtxts, join_with = ",")
        update_args.add("--output_scene_object_gzf_file", updated_gzf.path)
        update_args.add("--output_scene_object_pbtxt_file", updated_pbtxt.path)
        update_args.add_joined("--file_descriptor_sets", transitive_descriptor_sets, join_with = ",")
        ctx.actions.run(
            arguments = [update_args],
            executable = ctx.executable._update_scene_object,
            inputs = [raw_gzf] + ctx.files.updates_pbtxts + transitive_descriptor_sets.to_list(),
            mnemonic = "UpdateSceneObject",
            outputs = [updated_gzf, updated_pbtxt],
            progress_message = "Updating SceneObject %s " % ctx.label.name,
        )

    return [
        DefaultInfo(
            files = depset([output_gzf, output_pbtxt]),
            runfiles = ctx.runfiles(files = [output_gzf, output_pbtxt]),
        ),
        SceneObjectInfo(
            fdsets = transitive_descriptor_sets,
            gzf = output_gzf,
            pbtxt = output_pbtxt,
        ),
    ]

_usd_scene_object = rule(
    attrs = {
        "deps": attr.label_list(
            aspects = [gen_source_code_info_descriptor_set],
            providers = [ProtoInfo],
        ),
        "src": attr.label(allow_single_file = True),
        "updates_pbtxts": attr.label_list(allow_files = True),
        "_update_scene_object": attr.label(
            allow_single_file = True,
            cfg = "exec",
            default = Label("//intrinsic/scene/tools:update_scene_object"),
            executable = True,
        ),
        "_usd_file_converter": attr.label(
            allow_single_file = True,
            cfg = "exec",
            default = Label("//intrinsic/scene/usd:usd_file_converter"),
            executable = True,
        ),
    },
    doc = (
        "Converts a USD file to a scene object proto, to be used as a " +
        "product or resource in a PPR-based application (or anything which " +
        "accepts a SceneObjectInfo provider, which is returned by this rule)."
    ),
    implementation = _usd_scene_object_impl,
)

def usd_scene_object(
        name,
        src,
        updates_pbtxts = None,
        deps = [],
        compatible_with = [],
        testonly = False,
        visibility = None):
    """Creates SceneObject from USD file and optionally SceneObjectUpdates proto file.

    Args:
      name: The name of the SceneObject to be created.
      src: Source USD file. May be a .usd, .usda, .usdc, or .usdz file
      updates_pbtxts: list of labels, Optional. Text proto files
        (intrinsic_proto::scene_object::v1::SceneObjectUpdate) that specify
        post-processings of the converted SceneObject.
      deps: list of labels, Optional. Proto dependencies needed to parse the updates pbtxts.
      compatible_with: The list of environments this target can be built for,
        in addition to default-supported environments.
      testonly: If true, only testonly targets can depend on this target.
      visibility: Visibility of the build rule.
    """

    if updates_pbtxts and type(updates_pbtxts) != "list":
        fail("Expected updates_pbtxts to be of type list")

    # Actually generate the scene object gzf/textproto.
    _usd_scene_object(
        name = name,
        testonly = testonly,
        src = src,
        compatible_with = compatible_with,
        updates_pbtxts = updates_pbtxts,
        visibility = visibility,
        deps = deps,
    )
