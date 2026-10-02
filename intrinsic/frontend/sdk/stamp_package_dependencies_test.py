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

import json
import os

from absl.testing import absltest

from intrinsic.frontend.sdk import stamp_package_dependencies


class StampPackageDependenciesTest(absltest.TestCase):

  def setUp(self):
    super().setUp()
    tmpdir = self.create_tempdir().full_path
    self.src_file = os.path.join(tmpdir, "input_package.json")
    self.out_file = os.path.join(tmpdir, "output_package.json")
    self.status_file = os.path.join(tmpdir, "stable-status.txt")

  def test_stamp_version_and_dependencies(self):
    initial_pkg = {
        "name": "@intrinsic-ai/core",
        "version": "0.0.0",
        "dependencies": {
            "other-dep": "3.2.1",
            "@intrinsic-ai/types": "0.0.0",
        },
    }
    with open(self.src_file, "w", encoding="utf-8") as f:
      json.dump(initial_pkg, f)

    with open(self.status_file, "w", encoding="utf-8") as f:
      f.write("BUILD_USER test\nSTABLE_SEMVER 1.2.3-test\n")

    stamp_package_dependencies.main([
        "--src",
        self.src_file,
        "--out",
        self.out_file,
        "--stamp_var",
        "STABLE_SEMVER",
        "--stamp_dependency",
        "@intrinsic-ai/types",
        "--status_file",
        self.status_file,
    ])

    with open(self.out_file, "r", encoding="utf-8") as f:
      stamped = json.load(f)

    self.assertEqual(stamped["version"], "1.2.3-test")
    self.assertEqual(
        stamped["dependencies"]["@intrinsic-ai/types"], "1.2.3-test"
    )
    self.assertEqual(stamped["dependencies"]["other-dep"], "3.2.1")

  def test_unstamped_fallback_version(self):
    initial_pkg = {
        "name": "@intrinsic-ai/types",
        "version": "0.1.2",
    }
    with open(self.src_file, "w", encoding="utf-8") as f:
      json.dump(initial_pkg, f)

    # Status file does not contain STABLE_SEMVER
    with open(self.status_file, "w", encoding="utf-8") as f:
      f.write("BUILD_USER test\n")

    stamp_package_dependencies.main([
        "--src",
        self.src_file,
        "--out",
        self.out_file,
        "--stamp_var",
        "STABLE_SEMVER",
        "--status_file",
        self.status_file,
    ])

    with open(self.out_file, "r", encoding="utf-8") as f:
      stamped = json.load(f)

    self.assertEqual(stamped["version"], "0.1.2")

  def test_empty_status_files(self):
    initial_pkg = {
        "name": "@intrinsic-ai/types",
        "version": "0.0.0",
    }
    with open(self.src_file, "w", encoding="utf-8") as f:
      json.dump(initial_pkg, f)

    stamp_package_dependencies.main([
        "--src",
        self.src_file,
        "--out",
        self.out_file,
        "--stamp_var",
        "STABLE_SEMVER",
    ])

    with open(self.out_file, "r", encoding="utf-8") as f:
      stamped = json.load(f)

    self.assertEqual(stamped["version"], "0.0.0")

  def test_multiple_stamp_dependencies(self):
    initial_pkg = {
        "name": "@intrinsic-ai/bundle",
        "version": "0.0.0",
        "dependencies": {
            "@intrinsic-ai/types": "0.0.0",
            "@intrinsic-ai/core": "0.0.0",
            "other-dep": "3.2.1",
        },
    }
    with open(self.src_file, "w", encoding="utf-8") as f:
      json.dump(initial_pkg, f)

    with open(self.status_file, "w", encoding="utf-8") as f:
      f.write("STABLE_SEMVER 1.2.3-test\n")

    stamp_package_dependencies.main([
        "--src",
        self.src_file,
        "--out",
        self.out_file,
        "--stamp_var",
        "STABLE_SEMVER",
        "--stamp_dependency",
        "@intrinsic-ai/types",
        "--stamp_dependency",
        "@intrinsic-ai/core",
        "--status_file",
        self.status_file,
    ])

    with open(self.out_file, "r", encoding="utf-8") as f:
      stamped = json.load(f)

    self.assertEqual(stamped["version"], "1.2.3-test")
    self.assertEqual(
        stamped["dependencies"]["@intrinsic-ai/types"], "1.2.3-test"
    )
    self.assertEqual(
        stamped["dependencies"]["@intrinsic-ai/core"], "1.2.3-test"
    )
    self.assertEqual(stamped["dependencies"]["other-dep"], "3.2.1")

  def test_malformed_json_exits(self):
    with open(self.src_file, "w", encoding="utf-8") as f:
      f.write("{ invalid json")

    with self.assertRaises(SystemExit):
      stamp_package_dependencies.main([
          "--src",
          self.src_file,
          "--out",
          self.out_file,
          "--stamp_var",
          "STABLE_SEMVER",
      ])

  def test_non_dict_json_exits(self):
    with open(self.src_file, "w", encoding="utf-8") as f:
      json.dump([], f)

    with self.assertRaises(SystemExit):
      stamp_package_dependencies.main([
          "--src",
          self.src_file,
          "--out",
          self.out_file,
          "--stamp_var",
          "STABLE_SEMVER",
      ])

  def test_status_file_read_error_exits(self):
    initial_pkg = {
        "name": "@intrinsic-ai/types",
        "version": "0.0.0",
    }
    with open(self.src_file, "w", encoding="utf-8") as f:
      json.dump(initial_pkg, f)

    nonexistent_file = os.path.join(
        os.path.dirname(self.status_file), "nonexistent.txt"
    )

    with self.assertRaises(SystemExit):
      stamp_package_dependencies.main([
          "--src",
          self.src_file,
          "--out",
          self.out_file,
          "--status_file",
          nonexistent_file,
      ])


if __name__ == "__main__":
  absltest.main()
