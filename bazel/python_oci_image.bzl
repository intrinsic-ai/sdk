# Copyright 2023 Intrinsic Innovation LLC

"""Helpers for dealing with Python docker images."""

load("@bazel_skylib//lib:paths.bzl", "paths")
load("@tar.bzl", "mtree_mutate", "mtree_spec", "tar")
load("//bazel:container.bzl", "container_image")

_INTERPRETER_SYMLINKS = {
    "/bin/2to3": "2to3-3.11",
    "/bin/idle3": "idle3.11",
    "/bin/pip": "pip3.11",
    "/bin/pip3": "pip3.11",
    "/bin/pydoc3": "pydoc3.11",
    "/bin/python": "python3.11",
    "/bin/python3": "python3.11",
    "/bin/python3-config": "python3.11-config",
    "/libpython3.11.so": "libpython3.11.so.1.0",
}

# By default, py_binary runfiles manifests list aliases like /bin/python, /bin/python3,
# and /libpython3.11.so as duplicate regular files (`type=file`) rather than true symlinks.
# If archived without modification, the image tarball contains multiple 55+ MB copies of
# identical binaries and shared libraries. Furthermore, using the standard `symlinks`
# attribute on `container_image` only layers symlink entries on top of lower layers
# without removing the physical duplicate files archived inside `_interpreter_layer`.
# To prevent this bloating, we intercept and rewrite the mtree manifest lines with `sed`
# before the tarball is built, converting duplicate file entries into lightweight `type=link`
# symlink headers directly inside `_interpreter_layer`.
def _build_interpreter_manifest_cmd(interpreter_regex):
    """Returns the shell command to build the interpreter tar manifest with symlinks restored."""
    sed_exprs = [
        "-e 's|{path} uid=.*|{path} uid=0 gid=0 time=1672560000 mode=0777 type=link link={target}|'".format(
            path = path,
            target = target,
        )
        for path, target in _INTERPRETER_SYMLINKS.items()
    ]
    return "grep '{}' $< | sed {} >$@".format(interpreter_regex, " ".join(sed_exprs))

def python_layers(name, binary, **kwargs):
    """Create list of layers for a py_binary target.

    We use three layers for the interpreter, third-party dependencies, and application code.

    The setup is adapted from https://github.com/aspect-build/bazel-examples/blob/main/oci_python_image/py_layer.bzl.

    Args:
        name: prefix for generated targets, to ensure they are unique within the package
        binary: the py_binary target.
        **kwargs: extra arguments to pass on to the layers.
    Returns:
        a list of labels for the layers, which are tar files
    """

    layers = []

    # Produce the manifest for a tar file of our py_binary, but don't tar it up yet, so we can split
    # into fine-grained layers for better docker performance.
    mtree_spec(
        name = name + "_tar_manifest_raw",
        testonly = kwargs.get("testonly"),
        srcs = [binary],
        compatible_with = kwargs.get("compatible_with"),
        include_runfiles = True,
        tags = kwargs.get("tags"),
    )

    # ADDITION: Handle local_repository sub repos by removing '../' and ' external/' from paths.
    # Without this the resulting image manifest is malformed and tools like dive cannot open the image.
    native.genrule(
        name = name + "_tar_manifest_filtered",
        testonly = kwargs.get("testonly"),
        srcs = [":" + name + "_tar_manifest_raw"],
        outs = [name + "_tar_manifest_filtered.spec"],
        cmd = "sed -e 's#^\\.\\./##' $< | sed -e 's# external/##g' >$@",
        compatible_with = kwargs.get("compatible_with"),
        tags = kwargs.get("tags"),
    )

    # Apply mutations for path prefixes
    mtree_mutate(
        name = name + "_tar_manifest_prefix",
        testonly = kwargs.get("testonly"),
        compatible_with = kwargs.get("compatible_with"),
        mtree = ":" + name + "_tar_manifest_filtered",
        # BUG: https://github.com/bazel-contrib/bazel-lib/issues/946
        # strip_prefix = kwargs.pop("data_path", None),
        package_dir = kwargs.pop("directory", None),
        tags = kwargs.get("tags"),
    )

    # Workaround unsupported "strip_prefix"
    native.genrule(
        name = name + "_tar_manifest",
        testonly = kwargs.get("testonly"),
        srcs = [":" + name + "_tar_manifest_prefix"],
        outs = [name + "_tar_manifest.spec"],
        cmd = "sed -e 's,^/,,' $< >$@",
        compatible_with = kwargs.get("compatible_with"),
        tags = kwargs.get("tags"),
    )

    # One layer with only the python interpreter.
    # Bzlmod: "runfiles/rules_python~0.27.1~python~python_3_11_x86_64-unknown-linux-gnu/"
    PY_INTERPRETER_REGEX = "\\S*\\.runfiles/\\S*\\(rules_python\\S*_x86_64-unknown-linux-gnu/\\|.*rules_Upython++python+python.*libpython.*\\)"

    native.genrule(
        name = name + "_interpreter_tar_manifest",
        testonly = kwargs.get("testonly"),
        srcs = [":" + name + "_tar_manifest"],
        outs = [name + "_interpreter_tar_manifest.spec"],
        cmd = _build_interpreter_manifest_cmd(PY_INTERPRETER_REGEX),
        compatible_with = kwargs.get("compatible_with"),
        tags = kwargs.get("tags"),
    )

    tar(
        name = name + "_interpreter_layer",
        testonly = kwargs.get("testonly"),
        srcs = [binary],
        compatible_with = kwargs.get("compatible_with"),
        compress = "gzip",
        compute_unused_inputs = 1,
        mtree = ":" + name + "_interpreter_tar_manifest",
        tags = kwargs.get("tags"),
    )
    layers.append(":" + name + "_interpreter_layer")

    # Attempt to match all external (3P) dependencies. Since these can come in as either
    # `requirement` or native Bazel deps, do our best to guess the runfiles path.
    PACKAGES_REGEX = "\\S*\\.runfiles/\\S*\\(site-packages\\|com_\\|pip_deps_\\)"

    # One layer with the third-party pip packages.
    # To make sure some dependencies with surprising paths are not included twice, exclude the interpreter from the site-packages layer.
    native.genrule(
        name = name + "_packages_tar_manifest",
        testonly = kwargs.get("testonly"),
        srcs = [":" + name + "_tar_manifest"],
        outs = [name + "_packages_tar_manifest.spec"],
        cmd = "if ! grep -v '{}' $< | grep '{}' >$@; then touch $@; fi".format(PY_INTERPRETER_REGEX, PACKAGES_REGEX),
        compatible_with = kwargs.get("compatible_with"),
        tags = kwargs.get("tags"),
    )

    tar(
        name = name + "_packages_layer",
        testonly = kwargs.get("testonly"),
        srcs = [binary],
        compatible_with = kwargs.get("compatible_with"),
        compress = "gzip",
        compute_unused_inputs = 1,
        mtree = ":" + name + "_packages_tar_manifest",
        tags = kwargs.get("tags"),
    )
    layers.append(":" + name + "_packages_layer")

    # Any lines that didn't match one of the two grep above...
    native.genrule(
        name = name + "_app_tar_manifest",
        testonly = kwargs.get("testonly"),
        srcs = [":" + name + "_tar_manifest"],
        outs = [name + "_app_tar_manifest.spec"],
        cmd = "grep -v '{}' $< | grep -v '{}' >$@".format(PACKAGES_REGEX, PY_INTERPRETER_REGEX),
        compatible_with = kwargs.get("compatible_with"),
        tags = kwargs.get("tags"),
    )

    # ... go into the third layer which is the application. We assume it changes the most frequently.
    tar(
        name = name + "_app_layer",
        testonly = kwargs.get("testonly"),
        srcs = [binary],
        compatible_with = kwargs.get("compatible_with"),
        compress = "gzip",
        compute_unused_inputs = 1,
        mtree = ":" + name + "_app_tar_manifest",
        tags = kwargs.get("tags"),
    )
    layers.append(":" + name + "_app_layer")

    return layers

def python_oci_image(
        name,
        binary,
        base = None,
        extra_tars = None,
        symlinks = None,
        **kwargs):
    """Wrapper for creating a oci_image from a py_binary target.

    Will create both an oci_image ($name) and a container_tarball ($name.tar) target.

    The setup is inspired by https://github.com/aspect-build/bazel-examples/blob/main/oci_python_image/hello_world/BUILD.bazel.

    Args:
      name: name of the image.
      base: base image to use.
      binary: the py_binary target.
      extra_tars: additional layers to add to the image with e.g. supporting files.
      symlinks: if specified, symlinks to add to the final image (analogous to rules_docker container_image#sylinks).
      **kwargs: extra arguments to pass on to the oci_image target.
    """

    if base == None:
        base = Label("@distroless_python3")

    layer_kwargs = {key: value for key, value in kwargs.items() if key in ["compatible_with", "directory", "tags", "testonly"]}
    layers = python_layers(
        name = name,
        binary = binary,
        visibility = ["//visibility:private"],
        **layer_kwargs
    )

    binary_label = native.package_relative_label(binary)
    package_str = binary_label.package

    binary_path = paths.join("/", kwargs.get("directory", ""), package_str, binary_label.name)

    if kwargs.get("cmd") == None:
        kwargs["cmd"] = [binary_path]

    # By default, the runfiles path is prefixed by the repo name while the binary path is not. Adding a symlink to fix the runfiles detection logic in case the repo is not empty.
    if binary_label.repo_name:
        binary_runfiles_path = paths.join("/", kwargs.get("directory", ""), binary_label.repo_name, binary_label.package, binary_label.name + ".runfiles")
        symlinks = (symlinks or {}) | {
            binary_path + ".runfiles": binary_runfiles_path,
        }

    if extra_tars:
        layers.extend(extra_tars)

    container_image(
        name = name,
        base = base,
        layers = layers,
        symlinks = symlinks,
        **kwargs
    )
