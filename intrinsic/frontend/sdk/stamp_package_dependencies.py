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

"""Stamps package.json versions and internal dependencies for release builds.

In development, internal packages use "0.0.0". When building release tarballs
with `--stamp`, this script extracts the release semver from Bazel's workspace
status files and rewrites the package version and specified dependency versions
so that published npm artifacts correctly reference each other.
"""

import argparse
from collections.abc import Sequence
import json
import sys


def main(argv: Sequence[str] | None = None) -> None:
  parser = argparse.ArgumentParser(description=__doc__)
  parser.add_argument(
      "--src", required=True, help="Path to input package.json file."
  )
  parser.add_argument(
      "--out", required=True, help="Path to output package.json file."
  )
  parser.add_argument(
      "--stamp_var",
      default="STABLE_SEMVER",
      help="Workspace status variable to read (default: STABLE_SEMVER).",
  )
  parser.add_argument(
      "--stamp_dependency",
      dest="stamp_dependencies",
      action="append",
      default=[],
      help="Dependency package name in package.json to stamp (repeatable).",
  )
  parser.add_argument(
      "--status_file",
      dest="status_files",
      action="append",
      default=[],
      help="Path to Bazel workspace status file to inspect (repeatable).",
  )
  args = parser.parse_args(argv)

  version: str | None = None
  for sf in args.status_files:
    try:
      with open(sf, "r", encoding="utf-8") as f:
        for line in f:
          key, _, val = line.partition(" ")
          if key == args.stamp_var:
            version = val.strip()
            break
    except OSError as e:
      sys.stderr.write(
          "[stamp_package_dependencies] ERROR: failed to read status file"
          f" {sf}: {e}\n"
      )
      sys.exit(1)
    if version is not None:
      break

  try:
    with open(args.src, "r", encoding="utf-8") as f:
      pkg = json.load(f)
  except (OSError, json.JSONDecodeError) as e:
    sys.stderr.write(
        f"[stamp_package_dependencies] ERROR: failed to parse {args.src}: {e}\n"
    )
    sys.exit(1)

  if not isinstance(pkg, dict):
    sys.stderr.write(
        f"[stamp_package_dependencies] ERROR: {args.src} must contain a valid"
        f" JSON object (got {type(pkg).__name__})\n"
    )
    sys.exit(1)

  if version is None:
    version = pkg.get("version", "0.0.0")

  pkg["version"] = version

  if args.stamp_dependencies:
    for section in ("dependencies", "peerDependencies", "optionalDependencies"):
      deps = pkg.get(section)
      if isinstance(deps, dict):
        for dep in args.stamp_dependencies:
          if dep in deps:
            deps[dep] = version

  try:
    with open(args.out, "w", encoding="utf-8") as f:
      json.dump(pkg, f, indent=2)
      f.write("\n")
  except OSError as e:
    sys.stderr.write(
        "[stamp_package_dependencies] ERROR: failed to write output"
        f" {args.out}: {e}\n"
    )
    sys.exit(1)


if __name__ == "__main__":
  main()
