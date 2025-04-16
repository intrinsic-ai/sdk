# Copyright 2023 Intrinsic Innovation LLC

"""Helpers for dealing with the rules_docker->rules_oci transition.
."""

load("@container_structure_test//:defs.bzl", "container_structure_test")
load(
    "@rules_oci//oci:defs.bzl",
    "oci_image",
    "oci_load",
)
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")

def _container_import_impl(ctx):
    output = ctx.actions.declare_directory(ctx.label.name)
    ctx.actions.run(
        outputs = [output],
        inputs = [ctx.file.tarball],
        executable = ctx.toolchains["@rules_oci//oci:regctl_toolchain_type"].regctl_info.binary,
        arguments = [
            "image",
            "import",
            "ocidir://%s" % output.path,
            ctx.file.tarball.path,
        ],
        mnemonic = "ExtractContainerTarball",
    )
    return DefaultInfo(
        files = depset([output]),
        runfiles = ctx.runfiles(files = [output]),
    )

container_import = rule(
    implementation = _container_import_impl,
    doc = "Imports an image tarball into an oci-layout directory",
    attrs = {
        "tarball": attr.label(
            allow_single_file = [".tar"],
        ),
    },
    toolchains = [
        "@rules_oci//oci:regctl_toolchain_type",
    ],
)

def _symlink_tarball_impl(ctx):
    ctx.actions.symlink(output = ctx.outputs.output, target_file = ctx.attr.src[OutputGroupInfo].tarball.to_list()[0])

_symlink_tarball = rule(
    implementation = _symlink_tarball_impl,
    doc = "Creates a symlink to tarball.tar in src's DefaultInfo at output",
    attrs = {
        "src": attr.label(
            providers = [OutputGroupInfo],
            mandatory = True,
        ),
        "output": attr.output(),
    },
)

def _container_tarball(name, image, **kwargs):
    oci_load(
        name = name,
        image = image,
        **kwargs
    )
    _symlink_tarball(
        name = "%s_symlink" % name,
        src = name,
        output = "%s.tar" % image,
        compatible_with = kwargs.get("compatible_with"),
        visibility = kwargs.get("visibility"),
        testonly = kwargs.get("testonly"),
    )

def container_layer(name, **kwargs):
    pkg_tar(
        name = name,
        extension = "tar.gz",
        compressor_args = "--fast",
        package_dir = kwargs.pop("directory", None),
        strip_prefix = kwargs.pop("data_path", None),
        srcs = kwargs.pop("files", None),
        deps = kwargs.pop("tars", None),
        **kwargs
    )

# buildozer: disable=function-docstring-args
def container_image(
        name,
        base = None,
        cmd = None,
        data_path = None,
        directory = None,
        entrypoint = None,
        layers = None,
        tars = None,
        files = None,
        symlinks = None,
        labels = None,
        **kwargs):
    """Wrapper for creating an oci_image from a rules_docker container_image target.

    Will create both an oci_image ($name) and a _container_tarball ($name.tar) target.

    Note that it does not support the experimental_tarball_format attribute:
    - All tarballs created by this macro will be in .tar.gz format.
    - Existing tarballs won't be compressed if they are not already compressed.

    See https://docs.aspect.build/guides/rules_oci_migration/#container_image for the official conversion documentation.
    """
    if not layers:
        layers = []

    if tars:
        container_layer(
            name = name + "_tar_layer",
            tars = tars,
            data_path = data_path,
            directory = directory,
            compatible_with = kwargs.get("compatible_with"),
            visibility = kwargs.get("visibility"),
            testonly = kwargs.get("testonly"),
        )
        layers.append(name + "_tar_layer")

    if files:
        container_layer(
            name = name + "_files_layer",
            files = files,
            data_path = data_path,
            directory = directory,
            compatible_with = kwargs.get("compatible_with"),
            visibility = kwargs.get("visibility"),
            testonly = kwargs.get("testonly"),
        )
        layers.append(name + "_files_layer")

    if symlinks:
        container_layer(
            name = name + "_symlink_layer",
            symlinks = symlinks,
            data_path = "/",
            compatible_with = kwargs.get("compatible_with"),
            visibility = kwargs.get("visibility"),
            testonly = kwargs.get("testonly"),
        )
        layers.append(name + "_symlink_layer")

    oci_image(
        name = name,
        base = base,
        tars = layers,
        entrypoint = entrypoint,
        cmd = cmd,
        labels = labels,
        **kwargs
    )

    tag = "%s:latest" % name
    package = native.package_name()
    if package:
        tag = "%s/%s" % (package, tag)

    tarball_name = "_%s_tarball" % name
    _container_tarball(
        name = tarball_name,
        image = name,
        compatible_with = kwargs.get("compatible_with"),
        repo_tags = [tag],
        visibility = kwargs.get("visibility"),
        testonly = kwargs.get("testonly"),
    )

    path_under_test = None
    if entrypoint != None and len(entrypoint) > 0 and (entrypoint[0].startswith("/") or entrypoint[0].startswith("intrinsic")):
        path_under_test = entrypoint[0]
    if cmd != None and len(cmd) > 0 and (cmd[0].startswith("/") or cmd[0].startswith("intrinsic")):
        path_under_test = cmd[0]
    if path_under_test != None:
        test_config_name = "_%s_test_config" % name
        test_config_template = """schemaVersion: "2.0.0"

fileExistenceTests:
- name: "Entrypoint existence test"
  path: "{}"
  shouldExist: true
  isExecutableBy: "owner"
        """.format(path_under_test)
        native.genrule(
            name = test_config_name,
            srcs = [],
            outs = ["_%s_test_config.yaml" % name],
            cmd = "echo '{}' > $@".format(test_config_template),
            tools = [],
        )
        container_structure_test(
            name = "_%s_test" % name,
            timeout = "eternal",
            configs = [test_config_name],
            driver = "tar",
            image = "%s.tar" % name,
        )
