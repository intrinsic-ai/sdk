# Copyright 2023 Intrinsic Innovation LLC

"""This file contains starlark rules for building shell targets."""

load("@rules_shell//shell:sh_binary.bzl", _sh_binary = "sh_binary")
load("@rules_shell//shell:sh_library.bzl", _sh_library = "sh_library")
load("@rules_shell//shell:sh_test.bzl", _sh_test = "sh_test")

def _append_gbash_path(kwargs):
    env = kwargs.get("env", {})

    # NOTE Default PATH=.:/bin:/usr/bin:/usr/local/bin, "." is always prepended
    PATH = env.get("PATH", "/bin:/usr/bin:/usr/local/bin").split(":")

    # NOTE this needs to come BEFORE /usr/bin such that it is preferred over /usr/bin/gbash.sh on Cloudtops
    PATH.insert(0, "../ai_intrinsic_sdks+/third_party/imported/google/util/shell/gbash")

    env["PATH"] = ":".join(PATH)
    kwargs["env"] = env

def gbash_binary(name, **kwargs):
    _append_gbash_path(kwargs)
    _sh_binary(
        name = name,
        **kwargs
    )

def gbash_test(name, **kwargs):
    _append_gbash_path(kwargs)
    _sh_test(
        name = name,
        **kwargs
    )

sh_binary = gbash_binary
sh_library = _sh_library
sh_test = gbash_test
