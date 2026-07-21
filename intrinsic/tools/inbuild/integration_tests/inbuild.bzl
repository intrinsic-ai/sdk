# Copyright 2023 Intrinsic Innovation LLC

"""inbuild.bzl contains rules that invoke inbuild to build assets for Flowstate."""

load("@bazel_skylib//lib:paths.bzl", "paths")
load("@com_google_protobuf//bazel/common:proto_info.bzl", "ProtoInfo")

InbuildSkillManifestInfo = provider(
    doc = "Encapsulates binary manifest and file descriptor set artifacts.",
    fields = [
        "manifest",
        "file_descriptor_set",
    ],
)

def _inbuild_skill_manifest_impl(ctx):
    manifest_out = ctx.actions.declare_file(ctx.label.name + ".pbbin")
    fds_out = ctx.actions.declare_file(ctx.label.name + "_filedescriptor.pbbin")
    args = ctx.actions.args()
    args.add("skill").add("manifest")
    args.add("--manifest", ctx.file.manifest)
    args.add_joined(
        "--file_descriptor_sets",
        ctx.attr.proto[ProtoInfo].transitive_descriptor_sets,
        join_with = ",",
    )
    args.add("--output", manifest_out.path)
    args.add("--file_descriptor_set_out", fds_out.path)

    ctx.actions.run(
        arguments = [args],
        executable = ctx.executable._inbuild,
        inputs = depset([ctx.file.manifest], transitive = [ctx.attr.proto[ProtoInfo].transitive_descriptor_sets]),
        outputs = [manifest_out, fds_out],
    )

    return [
        DefaultInfo(files = depset([manifest_out, fds_out])),
        InbuildSkillManifestInfo(
            file_descriptor_set = fds_out,
            manifest = manifest_out,
        ),
    ]

inbuild_skill_manifest = rule(
    attrs = {
        "manifest": attr.label(
            allow_single_file = True,
            mandatory = True,
        ),
        "proto": attr.label(
            mandatory = True,
            providers = [ProtoInfo],
        ),
        "_inbuild": attr.label(
            cfg = "exec",
            default = Label("//intrinsic/tools/inbuild:inbuild"),
            doc = "The inbuild executable.",
            executable = True,
        ),
    },
    doc = "Generates binary manifest and consolidated file descriptor set",
    implementation = _inbuild_skill_manifest_impl,
)

def _inbuild_skill_bundle_impl(ctx):
    output_file = ctx.actions.declare_file(ctx.label.name + ".bundle.tar")
    manifest = ctx.attr.manifest[InbuildSkillManifestInfo].manifest
    fds = ctx.attr.manifest[InbuildSkillManifestInfo].file_descriptor_set
    args = ctx.actions.args()
    args.add("skill").add("bundle")
    args.add("--augmented_manifest", manifest.path)
    args.add("--augmented_file_descriptor_set", fds.path)
    args.add("--oci_image", ctx.file.oci_image.path)
    args.add("--output", output_file.path)
    ctx.actions.run(
        arguments = [args],
        executable = ctx.executable._inbuild,
        inputs = [manifest, fds, ctx.file.oci_image],
        outputs = [output_file],
    )

    return [
        DefaultInfo(files = depset([output_file])),
    ]

inbuild_skill_bundle = rule(
    attrs = {
        "manifest": attr.label(
            mandatory = True,
            providers = [InbuildSkillManifestInfo],
        ),
        "oci_image": attr.label(
            allow_single_file = True,
        ),
        "_inbuild": attr.label(
            cfg = "exec",
            default = Label("//intrinsic/tools/inbuild:inbuild"),
            doc = "The inbuild executable.",
            executable = True,
        ),
    },
    doc = "Generates the final skill bundle",
    implementation = _inbuild_skill_bundle_impl,
)

def _inbuild_skill_generate_entrypoint_py_impl(ctx):
    output_file = ctx.actions.declare_file(ctx.label.name + ".py")
    augmented_manifest_out = ctx.actions.declare_file(ctx.label.name + "_augmented_manifest.pbbin")
    augmented_fds_out = ctx.actions.declare_file(ctx.label.name + "_augmented_filedescriptor.pbbin")
    manifest = ctx.attr.manifest[InbuildSkillManifestInfo].manifest
    fds = ctx.attr.manifest[InbuildSkillManifestInfo].file_descriptor_set
    args = ctx.actions.args()
    args.add("skill").add("generate").add("entrypoint")
    args.add("--manifest", manifest.path)
    args.add("--language", "python")
    args.add("--file_descriptor_set", fds.path)
    args.add("--augmented_manifest_out", augmented_manifest_out.path)
    args.add("--augmented_file_descriptor_set_out", augmented_fds_out.path)
    args.add("--output", output_file.path)
    ctx.actions.run(
        arguments = [args],
        executable = ctx.executable._inbuild,
        inputs = [manifest, fds],
        outputs = [output_file, augmented_manifest_out, augmented_fds_out],
    )

    return [
        DefaultInfo(files = depset([output_file, augmented_manifest_out, augmented_fds_out])),
        InbuildSkillManifestInfo(
            file_descriptor_set = augmented_fds_out,
            manifest = augmented_manifest_out,
        ),
    ]

inbuild_skill_generate_entrypoint_py = rule(
    attrs = {
        "manifest": attr.label(
            mandatory = True,
            providers = [InbuildSkillManifestInfo],
        ),
        "_inbuild": attr.label(
            cfg = "exec",
            default = Label("//intrinsic/tools/inbuild:inbuild"),
            doc = "The inbuild executable.",
            executable = True,
        ),
    },
    doc = "Generates the main entry point for a Python skill",
    implementation = _inbuild_skill_generate_entrypoint_py_impl,
)

def _inbuild_skill_generate_entrypoint_cc_impl(ctx):
    output_file = ctx.actions.declare_file(ctx.label.name + ".cc")
    augmented_manifest_out = ctx.actions.declare_file(ctx.label.name + "_augmented_manifest.pbbin")
    augmented_fds_out = ctx.actions.declare_file(ctx.label.name + "_augmented_filedescriptor.pbbin")
    manifest = ctx.attr.manifest[InbuildSkillManifestInfo].manifest
    fds = ctx.attr.manifest[InbuildSkillManifestInfo].file_descriptor_set
    args = ctx.actions.args()
    args.add("skill").add("generate").add("entrypoint")
    args.add("--manifest", manifest.path)
    args.add("--language", "cpp")
    args.add("--cc_header", ctx.attr.cc_header)
    args.add("--file_descriptor_set", fds.path)
    args.add("--augmented_manifest_out", augmented_manifest_out.path)
    args.add("--augmented_file_descriptor_set_out", augmented_fds_out.path)
    args.add("--output", output_file.path)
    ctx.actions.run(
        arguments = [args],
        executable = ctx.executable._inbuild,
        inputs = [manifest, fds],
        outputs = [output_file, augmented_manifest_out, augmented_fds_out],
    )

    return [
        DefaultInfo(files = depset([output_file, augmented_manifest_out, augmented_fds_out])),
        InbuildSkillManifestInfo(
            file_descriptor_set = augmented_fds_out,
            manifest = augmented_manifest_out,
        ),
    ]

inbuild_skill_generate_entrypoint_cc = rule(
    attrs = {
        "cc_header": attr.string(
            mandatory = True,
        ),
        "manifest": attr.label(
            mandatory = True,
            providers = [InbuildSkillManifestInfo],
        ),
        "_inbuild": attr.label(
            cfg = "exec",
            default = Label("//intrinsic/tools/inbuild:inbuild"),
            doc = "The inbuild executable.",
            executable = True,
        ),
    },
    doc = "Generates the main entry point for a C++ skill",
    implementation = _inbuild_skill_generate_entrypoint_cc_impl,
)

def _inbuild_skill_generate_config_impl(ctx):
    output_file = ctx.actions.declare_file(ctx.label.name + ".pbbin")
    manifest = ctx.attr.manifest[InbuildSkillManifestInfo].manifest
    fds = ctx.attr.manifest[InbuildSkillManifestInfo].file_descriptor_set
    args = ctx.actions.args()
    args.add("skill").add("generate").add("config")
    args.add("--augmented_manifest", manifest.path)
    args.add("--augmented_file_descriptor_set", fds.path)
    args.add("--output", output_file.path)
    ctx.actions.run(
        arguments = [args],
        executable = ctx.executable._inbuild,
        inputs = [manifest, fds],
        outputs = [output_file],
    )

    return [
        DefaultInfo(files = depset([output_file])),
    ]

inbuild_skill_generate_config = rule(
    attrs = {
        "manifest": attr.label(
            mandatory = True,
            providers = [InbuildSkillManifestInfo],
        ),
        "_inbuild": attr.label(
            cfg = "exec",
            default = Label("//intrinsic/tools/inbuild:inbuild"),
            doc = "The inbuild executable.",
            executable = True,
        ),
    },
    doc = "Generates the SkillServiceConfig for a skill",
    implementation = _inbuild_skill_generate_config_impl,
)

def _inbuild_service_bundle_impl(ctx):
    output_file = ctx.actions.declare_file(ctx.label.name + ".bundle.tar")
    args = ctx.actions.args()
    args.add("service").add("bundle")
    args.add("--manifest", ctx.file.manifest)
    args.add_all(ctx.attr.proto[ProtoInfo].transitive_descriptor_sets, before_each = "--file_descriptor_set")
    args.add("--oci_image", ctx.file.oci_image.path)
    args.add("--default_config", ctx.file.default_config)
    args.add("--output", output_file.path)
    ctx.actions.run(
        arguments = [args],
        executable = ctx.executable._inbuild,
        inputs = depset([ctx.file.manifest, ctx.file.oci_image, ctx.file.default_config], transitive = [ctx.attr.proto[ProtoInfo].transitive_descriptor_sets]),
        outputs = [output_file],
    )

    return [
        DefaultInfo(files = depset([output_file])),
    ]

inbuild_service_bundle = rule(
    attrs = {
        "default_config": attr.label(
            allow_single_file = True,
        ),
        "manifest": attr.label(
            allow_single_file = True,
            mandatory = True,
        ),
        "oci_image": attr.label(
            allow_single_file = True,
            mandatory = True,
        ),
        "proto": attr.label(
            providers = [ProtoInfo],
        ),
        "_inbuild": attr.label(
            cfg = "exec",
            default = Label("//intrinsic/tools/inbuild:inbuild"),
            doc = "The inbuild executable.",
            executable = True,
        ),
    },
    doc = "Generates the final service bundle",
    implementation = _inbuild_service_bundle_impl,
)

def _inbuild_data_bundle_impl(ctx):
    output_file = ctx.actions.declare_file(ctx.label.name + ".bundle.tar")
    args = ctx.actions.args()
    args.add("data").add("bundle")
    args.add("--manifest", ctx.file.manifest)

    input_files = [ctx.file.manifest]
    transitive_inputs = []

    if ctx.attr.proto:
        args.add_all(ctx.attr.proto[ProtoInfo].transitive_descriptor_sets, before_each = "--file_descriptor_set")
        transitive_inputs.append(ctx.attr.proto[ProtoInfo].transitive_descriptor_sets)

    if ctx.files.data:
        for f in ctx.files.data:
            rel_path = paths.relativize(f.short_path, paths.dirname(ctx.file.manifest.short_path))
            args.add("--reference_to_path", rel_path + "=" + f.path)
        input_files.extend(ctx.files.data)

    args.add("--output", output_file.path)
    ctx.actions.run(
        arguments = [args],
        executable = ctx.executable._inbuild,
        inputs = depset(input_files, transitive = transitive_inputs),
        outputs = [output_file],
    )

    return [
        DefaultInfo(files = depset([output_file])),
    ]

inbuild_data_bundle = rule(
    attrs = {
        "data": attr.label_list(
            allow_files = True,
            doc = "Data files to map via reference_to_path.",
        ),
        "manifest": attr.label(
            allow_single_file = True,
            mandatory = True,
        ),
        "proto": attr.label(
            providers = [ProtoInfo],
        ),
        "_inbuild": attr.label(
            cfg = "exec",
            default = Label("//intrinsic/tools/inbuild:inbuild"),
            doc = "The inbuild executable.",
            executable = True,
        ),
    },
    doc = "Generates the Data Asset bundle using inbuild data bundle",
    implementation = _inbuild_data_bundle_impl,
)
