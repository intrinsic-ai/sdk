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

"""Sanity test for the FakeKeyValueStore Python bindings."""

from absl.testing import absltest

from intrinsic.platform.common.proto import test_pb2
from intrinsic.platform.pubsub.python import fake_kvstore


class FakeKvStoreTest(absltest.TestCase):

  def test_fake_kvstore_basic_ops(self):
    store = fake_kvstore.FakeKeyValueStore()
    msg = test_pb2.TestMessageString(data="hello")
    store.set("my_key", msg)

    any_msg = store.get("my_key")
    got = test_pb2.TestMessageString()
    self.assertTrue(any_msg.Unpack(got))
    self.assertEqual(got.data, "hello")

    store.delete("my_key")
    with self.assertRaises(RuntimeError):
      store.get("my_key")

    prefixed = fake_kvstore.FakeKeyValueStore(prefix_override="custom")
    prefixed.set("k", msg)
    self.assertEqual(prefixed.list_all_keys(), ["custom/k"])
    got_prefixed = test_pb2.TestMessageString()
    self.assertTrue(prefixed.get("k").Unpack(got_prefixed))
    self.assertEqual(got_prefixed.data, "hello")


if __name__ == "__main__":
  absltest.main()
