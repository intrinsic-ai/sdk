# Copyright 2023 Intrinsic Innovation LLC

"""Helpers for dealing with C++ docker images."""

load("@bazel_skylib//lib:paths.bzl", "paths")
load(
    "//bazel:container.bzl",
    "container_image",
    "container_layer",
)

def cc_oci_image(
        name,
        binary,
        base = None,
        extra_tars = None,
        symlinks = None,
        **kwargs):
    """Wrapper for creating a oci_image from a cc_binary target.

    Will create both an oci_image ($name) and a container_tarball ($name.tar) target.

    Args:
      name: name of the image.
      base: base image to use.
      binary: the cc_binary target.
      extra_tars: additional layers to add to the image with e.g. supporting files.
      symlinks: if specified, symlinks to add to the final image (analogous to rules_docker container_image#sylinks).
      **kwargs: extra arguments to pass on to the oci_image target.
    """

    if base == None:
        base = Label("@distroless_cc")

    layer_kwargs = {key: value for key, value in kwargs.items() if key in ["compatible_with", "data_path", "directory", "testonly"]}
    container_layer(
        name = name + "_binary_layer",
        files = [binary],
        include_runfiles = True,  # Include dynamic libraries
        visibility = ["//visibility:private"],
        **layer_kwargs
    )
    layers = [name + "_binary_layer"]

    binary_label = native.package_relative_label(binary)
    binary_path = paths.join("/", kwargs.get("directory", ""), binary_label.package, binary_label.name)

    if kwargs.get("cmd") == None:
        kwargs["cmd"] = [binary_path]

    # By default, the runfiles path is prefixed by the repo name while the binary path is not. Adding a symlink to fix the runfiles detection logic in case the repo is not empty.
    if native.repo_name():
        binary_runfiles_path = paths.join("/", kwargs.get("directory", ""), native.repo_name(), binary_label.package, binary_label.name + ".runfiles")
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
