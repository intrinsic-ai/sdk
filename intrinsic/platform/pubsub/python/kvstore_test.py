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

import gc
import threading
import time
import weakref

from absl.testing import absltest
from google.protobuf import any_pb2
from google.protobuf import wrappers_pb2

from intrinsic.platform.pubsub.python import pubsub


class _Sentinel:
  """Weak-referenceable object captured by Python query callbacks."""


class KVStoreTest(absltest.TestCase):

  def test_make_key(self):
    self.assertEqual(
        pubsub.KeyValueStore.make_key("/foo", "bar", "baz/"), "foo/bar/baz"
    )
    self.assertEqual(
        pubsub.KeyValueStore.make_key("foo", "bar", "baz"), "foo/bar/baz"
    )
    self.assertEqual(
        pubsub.KeyValueStore.make_key("///foo", "bar///", "///baz///"),
        "foo/bar/baz",
    )
    self.assertEqual(
        pubsub.KeyValueStore.make_key("/foo/", "/bar/", "/baz/"), "foo/bar/baz"
    )
    self.assertEqual(pubsub.KeyValueStore.make_key("foo", "", "bar"), "foo/bar")
    self.assertEqual(
        pubsub.KeyValueStore.make_key("foo", "///", "bar"), "foo/bar"
    )
    self.assertEqual(pubsub.KeyValueStore.make_key("///", "///", "///"), "")
    self.assertEqual(pubsub.KeyValueStore.make_key(), "")
    self.assertEqual(
        pubsub.KeyValueStore.make_key("foo/bar", "baz"), "foo/bar/baz"
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

  def test_get_all_destroys_python_callbacks_with_gil_and_releases_on_error(
      self,
  ):
    ps = pubsub.PubSub()
    store = ps.key_value_store()
    try:
      # Repeated async queries followed immediately by list_all_keys exercise
      # background-thread callback destruction under the GIL without deadlock.
      for i in range(10):
        done = threading.Event()
        sentinel = _Sentinel()
        weak_sentinel = weakref.ref(sentinel)

        def on_kv(key, value, pin=sentinel):
          del key, value, pin

        def on_done(ev=done, pin=sentinel):
          del pin
          ev.set()

        query = store.get_all(f"py_gil_test/{i}/**", on_kv, on_done)
        del sentinel, on_kv, on_done
        self.assertTrue(done.wait(5))
        del query
        self.assertIsInstance(store.list_all_keys(1), list)
        # Wait briefly if the Zenoh background thread is still running the
        # deleter after on_done returned.
        for _ in range(50):
          gc.collect()
          if weak_sentinel() is None:
            break
          threading.Event().wait(0.01)
        self.assertIsNone(weak_sentinel())

      # Synchronous startup failure must immediately release Python closures.
      err_sentinel = _Sentinel()
      weak_err = weakref.ref(err_sentinel)

      def err_cb(key, value, pin=err_sentinel):
        del key, value, pin

      def err_done(pin=err_sentinel):
        del pin

      with self.assertRaisesRegex(Exception, "Error getting a key"):
        store.get_all("invalid?key", err_cb, err_done)
      del err_sentinel, err_cb, err_done
      gc.collect()
      self.assertIsNone(weak_err())
    finally:
      del store
      del ps
      pubsub.PubSub.destroy_session_when_unused()

  def test_set_with_verification_and_delete_release_gil(self):
    ps = pubsub.PubSub()
    store = ps.key_value_store()
    try:
      # Plain set and delete succeed with call_guard<gil_scoped_release>().
      store.set("py_gil_test/plain", any_pb2.Any())
      store.delete("py_gil_test/plain")

      # A blocking set_with_verification must release the GIL so another Python
      # thread (or Zenoh callback teardown) can make progress during the poll.
      write_started = threading.Event()
      write_done = threading.Event()
      progressed_during_write = threading.Event()

      def check_gil_released():
        write_started.wait(5)
        start = time.monotonic()
        while not write_done.is_set():
          if time.monotonic() - start >= 0.05:
            progressed_during_write.set()
            return
          time.sleep(0.01)

      worker = threading.Thread(target=check_gil_released)
      worker.start()
      opts = pubsub.SetWithVerificationOptions()
      opts.timeout = 0.5
      try:
        write_started.set()
        with self.assertRaisesRegex(
            Exception, "Timeout waiting for verification"
        ):
          store.set_with_verification(
              "py_gil_test/blocking", any_pb2.Any(), opts
          )
      finally:
        write_done.set()
        worker.join(timeout=5)
      self.assertTrue(progressed_during_write.is_set())
    finally:
      del store
      del ps
      pubsub.PubSub.destroy_session_when_unused()

  def test_set_and_set_with_verification_do_not_double_pack_any(self):
    ps = pubsub.PubSub()
    store = ps.key_value_store()
    received: list[tuple[str, any_pb2.Any]] = []
    got_all = threading.Event()

    def on_kv(key: str, value: any_pb2.Any) -> None:
      received.append((key, value))
      if len(received) >= 4:
        got_all.set()

    sub = None
    try:
      sub = store.create_subscription(
          "py_any_test/*", pubsub.TopicConfig(), on_kv, None
      )
      payload = wrappers_pb2.StringValue(value="already_packed")
      packed = any_pb2.Any()
      packed.Pack(payload)

      store.set("py_any_test/direct", packed)
      store.set("py_any_test/direct_hc", packed, high_consistency=False)
      store.set("py_any_test/unpacked", payload)

      # No KV storage queryable runs in this unit test, so read-back
      # verification times out after publishing the Any payload to subscribers.
      opts = pubsub.SetWithVerificationOptions()
      opts.timeout = 0.05
      with self.assertRaisesRegex(
          Exception, "Timeout waiting for verification"
      ):
        store.set_with_verification("py_any_test/verified", packed, opts)

      self.assertTrue(got_all.wait(5))
      self.assertCountEqual(
          [key for key, _ in received],
          [
              "kv_store/py_any_test/direct",
              "kv_store/py_any_test/direct_hc",
              "kv_store/py_any_test/unpacked",
              "kv_store/py_any_test/verified",
          ],
      )
      for _, wire_any in received:
        self.assertTrue(wire_any.Is(wrappers_pb2.StringValue.DESCRIPTOR))
        unpacked = wrappers_pb2.StringValue()
        self.assertTrue(wire_any.Unpack(unpacked))
        self.assertEqual(unpacked.value, "already_packed")
    finally:
      if sub is not None:
        sub.unsubscribe()
        del sub
      del store
      del ps
      pubsub.PubSub.destroy_session_when_unused()


if __name__ == "__main__":
  absltest.main()
