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

"""Tests for the KeyValueStore Python bindings."""

from absl.testing import absltest

from intrinsic.platform.pubsub.python import pubsub


class KVStoreTest(absltest.TestCase):

  def test_make_key(self):
    self.assertEqual(
        pubsub.KeyValueStore.MakeKey("/foo", "bar", "baz/"), "foo/bar/baz"
    )
    self.assertEqual(
        pubsub.KeyValueStore.MakeKey("foo", "bar", "baz"), "foo/bar/baz"
    )
    self.assertEqual(
        pubsub.KeyValueStore.MakeKey("///foo", "bar///", "///baz///"),
        "foo/bar/baz",
    )
    self.assertEqual(
        pubsub.KeyValueStore.MakeKey("/foo/", "/bar/", "/baz/"), "foo/bar/baz"
    )
    self.assertEqual(pubsub.KeyValueStore.MakeKey("foo", "", "bar"), "foo/bar")
    self.assertEqual(
        pubsub.KeyValueStore.MakeKey("foo", "///", "bar"), "foo/bar"
    )
    self.assertEqual(pubsub.KeyValueStore.MakeKey("///", "///", "///"), "")
    self.assertEqual(pubsub.KeyValueStore.MakeKey(), "")
    self.assertEqual(
        pubsub.KeyValueStore.MakeKey("foo/bar", "baz"), "foo/bar/baz"
    )

  def test_set_with_verification_options(self):
    mode = pubsub.SetWithVerificationOptions.VerificationMode
    self.assertCountEqual(
        [m.name for m in mode], ["FIRST_REPLY", "HIGH_CONSISTENCY"]
    )

    options = pubsub.SetWithVerificationOptions()
    self.assertEqual(options.mode, mode.HIGH_CONSISTENCY)
    self.assertEqual(options.timeout, 30.0)

    options.mode = mode.FIRST_REPLY
    options.timeout = 5.0
    self.assertEqual(options.mode, mode.FIRST_REPLY)
    self.assertEqual(options.timeout, 5.0)


if __name__ == "__main__":
  absltest.main()
