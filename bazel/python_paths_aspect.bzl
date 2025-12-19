# Copyright 2023 Intrinsic Innovation LLC

"""Defines a Bazel aspect for collecting Python paths for import resolution in IDEs.
See https://bazel.build/extending/aspects and
https://blog.bazel.build/2016/06/10/ide-support.html .
"""

load("@rules_python//python:py_info.bzl", "PyInfo")

def _python_paths_aspect_impl(target, ctx):
    """Aspect for collecting Python paths for import resolution in IDEs.

    This produces information for all rule targets that return a PyInfo provider
    (py_library, py_binary, py_proto_library, ...).
    For each target //path/to:foo to which this aspect is applied and which has a
    PyInfo, an output file "path/to/foo.python_paths.json" will be generated. This
    file contains the import paths for external dependencies in all transitive
    sources of the target. The paths are relative to the folders
    'bazel-bin/external' and 'bazel-workspace_name/external'.
    This aspect has a single output group called 'python_paths'.
    Example output file:
    ```json
    {
        "python_paths": [
            "ai_intrinsic_sdks",
            "ai_intrinsic_sdks_pip_deps_numpy",
            "ai_intrinsic_sdks_pip_deps_numpy/site-packages",
            "com_google_protobuf",
            "com_google_protobuf/python",
        ]
    }
    ```
    Note that for external repositories in which the "Python root" is nested inside
    of the repository the output will contain (at least) two entries. E.g.,
    "com_google_protobuf" and "com_google_protobuf/python" in the example above.
    Even though the repository root might not strictly be required as an import path,
    this is consistent with how Bazel constructs PYTHONPATH at runtime (the repository
    root gets included in PYTHONPATH) and an IDE should thus try to resolve imports
    relative to "com_google_protobuf" *and* "com_google_protobuf/python". See the
    construction of PYTHONPATH in Bazel's bootstrap template for Python which serves
    as the main script for every py_binary:
    https://github.com/bazelbuild/bazel/blob/0696ba32a789bbf3100f62fa2c1547fc74e36006/tools/python/python_bootstrap_template.txt#L444-L446
    """
    if PyInfo not in target:
        return []

    py_info = target[PyInfo]

    # Use a dict with mock values since sets are not generally supported in
    # Starlark.
    paths = {}

    for importPath in py_info.imports.to_list():
        paths[importPath] = ""

    for sourcePath in py_info.transitive_sources.to_list():
        if sourcePath.owner.repo_name:
            paths[sourcePath.owner.repo_name] = ""

    paths_struct = struct(python_paths = sorted(paths.keys()))
    paths_json = json.encode(paths_struct)

    file = ctx.actions.declare_file("%s.python_paths.json" % ctx.label.name)
    ctx.actions.write(output = file, content = paths_json + "\n")

    return [
        DefaultInfo(files = depset([file])),
        OutputGroupInfo(python_paths = depset([file])),
    ]

python_paths_aspect = aspect(
    implementation = _python_paths_aspect_impl,
    # Do not propagate automatically. PyInfo.imports and
    # PyInfo.transitive_sources are already accumulated across all transitive
    # dependencies so no additional traversal is necessary.
    attr_aspects = [],
)
