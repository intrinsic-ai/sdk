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

"""Tests for helper functions to resolve runfiles paths."""

import os

from absl.testing import absltest

from intrinsic.util.path_resolver import path_resolver

TEST_FILE = 'intrinsic/util/path_resolver/path_resolver_test.py'
NONEXISTENT_FILE = 'intrinsic/util/path_resolver/nonexistent_file.txt'


class PathResolverTest(absltest.TestCase):

  def test_resolve_runfiles_path(self):
    path = path_resolver.resolve_runfiles_path(TEST_FILE)
    self.assertTrue(os.path.exists(path))
    self.assertTrue(path.endswith(TEST_FILE))

  def test_resolve_runfiles_path_absolute(self):
    absolute_path = '/usr/local/some/absolute/path.txt'
    path = path_resolver.resolve_runfiles_path(absolute_path)
    self.assertEqual(path, absolute_path)

  def test_rlocation_valid_path(self):
    path = path_resolver.rlocation(
        os.path.join(path_resolver._repo_name, TEST_FILE)
    )
    self.assertTrue(os.path.exists(path))
    self.assertTrue(path.endswith(TEST_FILE))

  def test_rlocation_nonexistent_path(self):
    path = path_resolver.rlocation(
        os.path.join(path_resolver._repo_name, NONEXISTENT_FILE)
    )
    self.assertFalse(os.path.exists(path))
    # It should still return a normalized path ending with the file,
    # or handle the non-existent gracefully.
    self.assertTrue(path.endswith(NONEXISTENT_FILE))

  def test_resolve_runfiles_path_nonexistent_path(self):
    path = path_resolver.resolve_runfiles_path(NONEXISTENT_FILE)
    self.assertFalse(os.path.exists(path))
    self.assertTrue(path.endswith(NONEXISTENT_FILE))


if __name__ == '__main__':
  absltest.main()
