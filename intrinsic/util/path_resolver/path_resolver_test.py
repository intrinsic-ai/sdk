# Copyright 2023 Intrinsic Innovation LLC

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
