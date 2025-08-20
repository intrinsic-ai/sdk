# Copyright 2023 Intrinsic Innovation LLC

"""Build rules for creating Hardware Module Assets."""

def _hal_manifest_impl(ctx):
    out = ctx.actions.declare_file("%s.textproto" % ctx.label.name)

    args = ctx.actions.args().add(
        ctx.file.manifest,
        format = "--manifest=%s",
    ).add(
        ctx.attr.manifest_type,
        format = "--manifest_type=%s",
    ).add(
        ctx.attr.provides_service_inspection,
        format = "--provides_service_inspection=%s",
    ).add(
        ctx.attr.requires_rtpc_node,
        format = "--requires_rtpc_node=%s",
    ).add(
        ctx.attr.requires_atemsys,
        format = "--requires_atemsys=%s",
    ).add(
        ctx.attr.running_ethercat_oss,
        format = "--running_ethercat_oss=%s",
    ).add_all(
        ctx.attr.service_proto_prefixes,
        format_each = "--service_proto_prefix=%s",
    ).add(
        out,
        format = "--output=%s",
    )

    if not ctx.file.image and not ctx.file.image_sim:
        fail("One of image or image_sim must be provided")
    if ctx.file.image:
        args = args.add(
            ctx.file.image,
            format = "--image=%s",
        )
    if ctx.file.image_sim:
        args = args.add(
            ctx.file.image_sim,
            format = "--image_sim=%s",
        )

    ctx.actions.run(
        inputs = depset([ctx.file.manifest]),
        outputs = [out],
        executable = ctx.executable._hal_manifest,
        arguments = [args],
        mnemonic = "HalManifest",
        progress_message = "Generating complete hal manifest %s" % out.short_path,
    )
    return [
        DefaultInfo(
            files = depset([out]),
            runfiles = ctx.runfiles(files = [out]),
        ),
    ]

hardware_module_manifest = rule(
    implementation = _hal_manifest_impl,
    doc = "Completes a partial resource manifest file for hardware modules",
    attrs = {
        "manifest": attr.label(
            allow_single_file = [".textproto"],
            doc = "The partially complete manifest containing metadata",
            mandatory = True,
        ),
        "manifest_type": attr.string(
            default = "service",
            doc = """The type of manifest that is being provided.  Allowed
            values are: service.""",
        ),
        "image": attr.label(
            allow_single_file = [".tar"],
            doc = "The image archive to be included in the bundle",
        ),
        "image_sim": attr.label(
            allow_single_file = [".tar"],
            mandatory = True,
            doc = """The image archive to be included in the bundle for
            simulation.  This can be the same as image if it supports both sim
            and real""",
        ),
        "provides_service_inspection": attr.bool(
            default = False,
            doc = "Flag to indicate that the module provides service inspection",
        ),
        "requires_rtpc_node": attr.bool(
            default = True,
            doc = "If the hardware module requires a real-time PC or not",
        ),
        "requires_atemsys": attr.bool(
            default = False,
            doc = "If the hardware module requires a atemsys ethercat",
        ),
        "running_ethercat_oss": attr.bool(
            default = False,
            doc = "If the hardware module is running ethercat oss",
        ),
        "service_proto_prefixes": attr.string_list(
            doc = """The services the module exposes, if any. The ServiceState service is always
            exposed.""",
        ),
        "_hal_manifest": attr.label(
            default = Label("//intrinsic/icon/hal/bzl:hal_manifest"),
            cfg = "exec",
            executable = True,
        ),
    },
)
