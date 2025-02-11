# Copyright 2023 Intrinsic Innovation LLC

"""inbuild.bzl contains rules that invoke inbuild to build assets for Flowstate."""

def _inbuild_skill_bundle_impl(ctx):
    output_file = ctx.actions.declare_file(ctx.label.name + ".bundle.tar")
    args = ctx.actions.args()
    args.add("skill").add("bundle")
    args.add("--manifest", ctx.file.manifest)
    for fds in ctx.attr.proto[ProtoInfo].transitive_descriptor_sets.to_list():
        args.add("--file_descriptor_set", fds.path)
    args.add("--oci_image", ctx.file.oci_image.path)
    args.add("--output", output_file.path)
    ctx.actions.run(
        outputs = [output_file],
        executable = ctx.executable._inbuild,
        inputs = depset([ctx.file.manifest, ctx.file.oci_image], transitive = [ctx.attr.proto[ProtoInfo].transitive_descriptor_sets]),
        arguments = [args],
    )

    return [
        DefaultInfo(files = depset([output_file])),
    ]

inbuild_skill_bundle = rule(
    implementation = _inbuild_skill_bundle_impl,
    doc = "Generates the final skill bundle",
    attrs = {
        "manifest": attr.label(
            mandatory = True,
            allow_single_file = True,
        ),
        "proto": attr.label(
            providers = [ProtoInfo],
        ),
        "oci_image": attr.label(
            allow_single_file = True,
        ),
        "_inbuild": attr.label(
            default = Label("//intrinsic/tools/inbuild:inbuild"),
            doc = "The inbuild executable.",
            executable = True,
            cfg = "exec",
        ),
    },
)

def _inbuild_skill_generate_entrypoint_py_impl(ctx):
    output_file = ctx.actions.declare_file(ctx.label.name + ".py")
    args = ctx.actions.args()
    args.add("skill").add("generate").add("entrypoint")
    args.add("--manifest", ctx.file.manifest)
    args.add("--language", "python")
    args.add("--output", output_file.path)
    ctx.actions.run(
        outputs = [output_file],
        executable = ctx.executable._inbuild,
        inputs = [ctx.file.manifest],
        arguments = [args],
    )

    return [
        DefaultInfo(files = depset([output_file])),
    ]

inbuild_skill_generate_entrypoint_py = rule(
    implementation = _inbuild_skill_generate_entrypoint_py_impl,
    doc = "Generates the main entry point for a Python skill",
    attrs = {
        "manifest": attr.label(
            mandatory = True,
            allow_single_file = True,
        ),
        "_inbuild": attr.label(
            default = Label("//intrinsic/tools/inbuild:inbuild"),
            doc = "The inbuild executable.",
            executable = True,
            cfg = "exec",
        ),
    },
)

def _inbuild_skill_generate_entrypoint_cc_impl(ctx):
    output_file = ctx.actions.declare_file(ctx.label.name + ".cc")
    args = ctx.actions.args()
    args.add("skill").add("generate").add("entrypoint")
    args.add("--manifest", ctx.file.manifest)
    args.add("--language", "cpp")
    args.add("--cc_header", ctx.attr.cc_header)
    args.add("--output", output_file.path)
    ctx.actions.run(
        outputs = [output_file],
        executable = ctx.executable._inbuild,
        inputs = [ctx.file.manifest],
        arguments = [args],
    )

    return [
        DefaultInfo(files = depset([output_file])),
    ]

inbuild_skill_generate_entrypoint_cc = rule(
    implementation = _inbuild_skill_generate_entrypoint_cc_impl,
    doc = "Generates the main entry point for a C++ skill",
    attrs = {
        "manifest": attr.label(
            mandatory = True,
            allow_single_file = True,
        ),
        "cc_header": attr.string(
            mandatory = True,
        ),
        "_inbuild": attr.label(
            default = Label("//intrinsic/tools/inbuild:inbuild"),
            doc = "The inbuild executable.",
            executable = True,
            cfg = "exec",
        ),
    },
)

def _inbuild_skill_generate_config_impl(ctx):
    output_file = ctx.actions.declare_file(ctx.label.name + ".pbbin")
    args = ctx.actions.args()
    args.add("skill").add("generate").add("config")
    args.add("--manifest", ctx.file.manifest)
    for fds in ctx.attr.proto[ProtoInfo].transitive_descriptor_sets.to_list():
        args.add("--file_descriptor_set", fds.path)
    args.add("--output", output_file.path)
    ctx.actions.run(
        outputs = [output_file],
        executable = ctx.executable._inbuild,
        inputs = depset([ctx.file.manifest], transitive = [ctx.attr.proto[ProtoInfo].transitive_descriptor_sets]),
        arguments = [args],
    )

    return [
        DefaultInfo(files = depset([output_file])),
    ]

inbuild_skill_generate_config = rule(
    implementation = _inbuild_skill_generate_config_impl,
    doc = "Generates the SkillServiceConfig for a skill",
    attrs = {
        "manifest": attr.label(
            mandatory = True,
            allow_single_file = True,
        ),
        "proto": attr.label(
            providers = [ProtoInfo],
        ),
        "_inbuild": attr.label(
            default = Label("//intrinsic/tools/inbuild:inbuild"),
            doc = "The inbuild executable.",
            executable = True,
            cfg = "exec",
        ),
    },
)

def _inbuild_service_bundle_impl(ctx):
    output_file = ctx.actions.declare_file(ctx.label.name + ".bundle.tar")
    args = ctx.actions.args()
    args.add("service").add("bundle")
    args.add("--manifest", ctx.file.manifest)
    for fds in ctx.attr.proto[ProtoInfo].transitive_descriptor_sets.to_list():
        args.add("--file_descriptor_set", fds.path)
    args.add("--oci_image", ctx.file.oci_image.path)
    args.add("--default_config", ctx.file.default_config)
    args.add("--output", output_file.path)
    ctx.actions.run(
        outputs = [output_file],
        executable = ctx.executable._inbuild,
        inputs = depset([ctx.file.manifest, ctx.file.oci_image, ctx.file.default_config], transitive = [ctx.attr.proto[ProtoInfo].transitive_descriptor_sets]),
        arguments = [args],
    )

    return [
        DefaultInfo(files = depset([output_file])),
    ]

inbuild_service_bundle = rule(
    implementation = _inbuild_service_bundle_impl,
    doc = "Generates the final service bundle",
    attrs = {
        "manifest": attr.label(
            mandatory = True,
            allow_single_file = True,
        ),
        "proto": attr.label(
            providers = [ProtoInfo],
        ),
        "oci_image": attr.label(
            mandatory = True,
            allow_single_file = True,
        ),
        "default_config": attr.label(
            allow_single_file = True,
        ),
        "_inbuild": attr.label(
            default = Label("//intrinsic/tools/inbuild:inbuild"),
            doc = "The inbuild executable.",
            executable = True,
            cfg = "exec",
        ),
    },
)
