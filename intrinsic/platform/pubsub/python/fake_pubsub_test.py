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

"""Tests for the SDK in-memory PubSub key-value store fakes."""

import contextlib
import threading
import types

from absl.testing import absltest
from google.protobuf import any_pb2
from pybind11_abseil import status

from intrinsic.platform.common.proto import test_pb2
from intrinsic.platform.pubsub.python import fake_pubsub


class _BytesSubclass(bytes):
  """Bytes subclass used to verify exact-`bytes` normalization."""


@contextlib.contextmanager
def _assert_raises_status(
    test_case: absltest.TestCase,
    expected_code: status.StatusCode,
    expected_regex: str | None = None,
):
  cm = (
      test_case.assertRaisesRegex(status.StatusNotOk, expected_regex)
      if expected_regex is not None
      else test_case.assertRaises(status.StatusNotOk)
  )
  with cm as ctx:
    yield ctx
  test_case.assertEqual(ctx.exception.status.code(), expected_code)


class FakeKeyValueStoreTest(absltest.TestCase):

  def test_set_packs_message_and_get_unpacks_it(self):
    store = fake_pubsub.FakeKeyValueStore()
    msg = test_pb2.TestMessageString(data="cam0_payload")

    store.set("captures/cam0", msg)

    stored = store.get("captures/cam0")
    self.assertIsInstance(stored, any_pb2.Any)
    self.assertTrue(stored.Is(test_pb2.TestMessageString.DESCRIPTOR))
    unpacked = test_pb2.TestMessageString()
    self.assertTrue(stored.Unpack(unpacked))
    self.assertEqual(unpacked, msg)

  def test_set_with_any_value_copies_input_verbatim(self):
    store = fake_pubsub.FakeKeyValueStore()
    packed = any_pb2.Any()
    packed.Pack(test_pb2.TestMessageString(data="abc"))
    expected = any_pb2.Any()
    expected.CopyFrom(packed)

    store.set("buffer", packed)
    store.set_with_verification("verified_buffer", packed)
    packed.Clear()

    self.assertEqual(store.get("buffer"), expected)
    self.assertEqual(store.get("verified_buffer"), expected)

  def test_get_returns_a_copy(self):
    store = fake_pubsub.FakeKeyValueStore()
    store.set("buffer", test_pb2.TestMessageString(data="abc"))

    first = store.get("buffer")
    first.Clear()

    self.assertTrue(
        store.get("buffer").Is(test_pb2.TestMessageString.DESCRIPTOR)
    )

  def test_get_missing_key_raises_not_found_status(self):
    store = fake_pubsub.FakeKeyValueStore()
    with _assert_raises_status(self, status.StatusCode.NOT_FOUND, "missing"):
      store.get("missing")
    with _assert_raises_status(self, status.StatusCode.NOT_FOUND, "missing"):
      store.get_raw("missing")

  def test_set_raw_and_get_raw_with_cross_encoding_checks(self):
    store = fake_pubsub.FakeKeyValueStore()
    payload = b"\x00\x01\x02raw\xff"
    mutable_buf = bytearray(payload)
    store.set_raw("captures/cam0/0", mutable_buf)
    mutable_buf[0] = 0xAA

    raw_result = store.get_raw("captures/cam0/0", timeout=1.5)
    self.assertIsInstance(raw_result, bytes)
    self.assertEqual(raw_result, payload)

    # Strided memoryviews, set_raw_with_verification, and zero-byte payloads.
    strided_view = memoryview(b"a0b1c2")[::2]
    store.set_raw_with_verification(
        "captures/cam0/view",
        strided_view,
        options=types.SimpleNamespace(timeout=5.0),
    )
    self.assertEqual(store.get_raw("captures/cam0/view"), b"abc")
    store.set_raw("captures/cam0/empty", b"")
    self.assertEqual(store.get_raw("captures/cam0/empty"), b"")

    # Bytes subclasses are normalized to exact `bytes`, matching C++ `GetRaw`.
    store.set_raw("captures/cam0/sub", _BytesSubclass(b"sub_raw"))
    sub_result = store.get_raw("captures/cam0/sub")
    self.assertIs(type(sub_result), bytes)
    self.assertEqual(sub_result, b"sub_raw")

    store.set_raw_with_verification(
        "captures/cam0/sub_v",
        _BytesSubclass(b"sub_verified"),
        options=types.SimpleNamespace(timeout=5.0),
    )
    sub_v_result = store.get_raw("captures/cam0/sub_v")
    self.assertIs(type(sub_v_result), bytes)
    self.assertEqual(sub_v_result, b"sub_verified")

    with _assert_raises_status(
        self,
        status.StatusCode.INVALID_ARGUMENT,
        "kv_store/captures/cam0/0' stores raw bytes, not a Protobuf message",
    ):
      store.get("captures/cam0/0")

    store.set("captures/cam0", test_pb2.TestMessageString(data="meta"))
    with _assert_raises_status(
        self,
        status.StatusCode.INVALID_ARGUMENT,
        "kv_store/captures/cam0' stores a Protobuf message, not raw bytes",
    ):
      store.get_raw("captures/cam0")

    # Overwriting between raw and Protobuf encodings in the shared keyspace.
    store.set_raw("captures/cam0", b"overwritten")
    self.assertEqual(store.get_raw("captures/cam0"), b"overwritten")
    store.set("captures/cam0/0", test_pb2.TestMessageString(data="now_proto"))
    unpacked = test_pb2.TestMessageString()
    self.assertTrue(store.get("captures/cam0/0").Unpack(unpacked))
    self.assertEqual(unpacked.data, "now_proto")

    self.assertCountEqual(
        store.list_all_keys(),
        [
            "kv_store/captures/cam0",
            "kv_store/captures/cam0/0",
            "kv_store/captures/cam0/view",
            "kv_store/captures/cam0/empty",
            "kv_store/captures/cam0/sub",
            "kv_store/captures/cam0/sub_v",
        ],
    )
    store.delete("captures/**")
    self.assertEmpty(store.list_all_keys())

    for bad_value in ("text", 123, test_pb2.TestMessageString(data="x"), None):
      with self.subTest(bad_value=bad_value):
        with self.assertRaises(TypeError):
          store.set_raw("captures/bad", bad_value)  # pytype: disable=wrong-arg-types
        with self.assertRaises(TypeError):
          store.set_raw_with_verification("captures/bad", bad_value)  # pytype: disable=wrong-arg-types

  def test_signature_compatibility_keyword_arguments(self):
    pubsub = fake_pubsub.FakePubSub(
        participant_name="test_participant", config="{}"
    )
    store = pubsub.key_value_store(prefix_override="custom_store")
    msg = test_pb2.TestMessageString(data="val")
    verified_msg = test_pb2.TestMessageString(data="verified")
    store.set(key="key", value=msg, high_consistency=True)
    store.set_raw(key="raw_key", value=b"raw_val")
    store.set_with_verification(
        key="verified_key",
        value=verified_msg,
        options=types.SimpleNamespace(timeout=5.0),
    )
    store.set_raw_with_verification(
        key="verified_raw_key",
        value=b"verified_raw_val",
        options=types.SimpleNamespace(timeout=5.0),
    )

    self.assertCountEqual(
        store.list_all_keys(timeout=0),
        [
            "custom_store/key",
            "custom_store/raw_key",
            "custom_store/verified_key",
            "custom_store/verified_raw_key",
        ],
    )
    unpacked = test_pb2.TestMessageString()
    self.assertTrue(store.get(key="key", timeout=0).Unpack(unpacked))
    self.assertEqual(unpacked, msg)
    self.assertEqual(store.get_raw(key="raw_key", timeout=0.0), b"raw_val")
    self.assertTrue(store.get(key="verified_key", timeout=5).Unpack(unpacked))
    self.assertEqual(unpacked, verified_msg)
    self.assertEqual(
        store.get_raw(key="verified_raw_key", timeout=5.0), b"verified_raw_val"
    )

    store.delete(key="key")
    store.delete(key="raw_key")
    store.delete(key="verified_raw_key")
    self.assertCountEqual(
        store.list_all_keys(timeout=0), ["custom_store/verified_key"]
    )
    fake_pubsub.FakePubSub.destroy_session_when_unused()

  def test_invalid_keys_and_timeouts_raise_invalid_argument(self):
    store = fake_pubsub.FakeKeyValueStore()
    msg = test_pb2.TestMessageString(data="val")
    invalid_arg = status.StatusCode.INVALID_ARGUMENT
    for bad_key in ("", "/bad", "bad/", "bad//key"):
      with self.subTest(bad_key=bad_key):
        with _assert_raises_status(self, invalid_arg):
          store.set(bad_key, msg)
        with _assert_raises_status(self, invalid_arg):
          store.set_raw(bad_key, b"raw")
        with _assert_raises_status(self, invalid_arg):
          store.set_with_verification(bad_key, msg)
        with _assert_raises_status(self, invalid_arg):
          store.set_raw_with_verification(bad_key, b"raw")
        with _assert_raises_status(self, invalid_arg):
          store.get(bad_key)
        with _assert_raises_status(self, invalid_arg):
          store.get_raw(bad_key)
        with _assert_raises_status(self, invalid_arg):
          store.delete(bad_key)

    with _assert_raises_status(self, invalid_arg):
      store.get("valid/key", timeout=-1)
    with _assert_raises_status(self, invalid_arg):
      store.list_all_keys(timeout=-1)

    for non_int_timeout in (1.5, 0.0, float("nan"), float("inf"), True):
      with self.subTest(non_int_timeout=non_int_timeout):
        with self.assertRaises(TypeError):
          store.get("valid/key", timeout=non_int_timeout)  # pytype: disable=wrong-arg-types
        with self.assertRaises(TypeError):
          store.list_all_keys(timeout=non_int_timeout)  # pytype: disable=wrong-arg-types

    for bad_raw_timeout in (
        -1,
        -0.5,
        float("nan"),
        float("inf"),
        float("-inf"),
    ):
      with self.subTest(bad_raw_timeout=bad_raw_timeout):
        with _assert_raises_status(self, invalid_arg):
          store.get_raw("valid/key", timeout=bad_raw_timeout)
    with self.assertRaises(TypeError):
      store.get_raw("valid/key", timeout=True)

    for bad_timeout in (0.0, -1.0, float("nan"), float("inf"), float("-inf")):
      with self.subTest(bad_timeout=bad_timeout):
        with _assert_raises_status(self, invalid_arg):
          store.set_with_verification(
              "valid/key",
              msg,
              options=types.SimpleNamespace(timeout=bad_timeout),
          )
        with _assert_raises_status(self, invalid_arg):
          store.set_raw_with_verification(
              "valid/key",
              b"raw",
              options=types.SimpleNamespace(timeout=bad_timeout),
          )

    for non_numeric_opt_timeout in (True, "1.0"):
      with self.subTest(non_numeric_opt_timeout=non_numeric_opt_timeout):
        with self.assertRaises(TypeError):
          store.set_with_verification(
              "valid/key",
              msg,
              options=types.SimpleNamespace(timeout=non_numeric_opt_timeout),
          )
        with self.assertRaises(TypeError):
          store.set_raw_with_verification(
              "valid/key",
              b"raw",
              options=types.SimpleNamespace(timeout=non_numeric_opt_timeout),
          )

    # Buffer protocol validation happens before option timeout checks, matching
    # pybind11 argument conversion order.
    with self.assertRaises(TypeError):
      store.set_raw_with_verification(
          "valid/key",
          "not_a_buffer",  # pytype: disable=wrong-arg-types
          options=types.SimpleNamespace(timeout=0.0),
      )

    self.assertEmpty(store.list_all_keys())

  def test_delete_and_list_all_keys(self):
    store = fake_pubsub.FakeKeyValueStore()
    store.set("a", test_pb2.TestMessageString(data="1"))
    store.set("b", test_pb2.TestMessageString(data="2"))
    self.assertCountEqual(store.list_all_keys(), ["kv_store/a", "kv_store/b"])

    store.delete("a")
    store.delete("not_there")

    self.assertCountEqual(store.list_all_keys(), ["kv_store/b"])

  def test_delete_wildcard_removes_matching_subtree_only(self):
    store = fake_pubsub.FakeKeyValueStore()
    store.set("cam_a/latest", test_pb2.TestMessageString(data="root"))
    store.set("cam_a/latest/0", test_pb2.TestMessageString(data="s0"))
    store.set("cam_a/latest/0/rgb", test_pb2.TestMessageString(data="s0_rgb"))
    store.set("cam_a/latest/1", test_pb2.TestMessageString(data="s1"))
    store.set("cam_a/latest_other", test_pb2.TestMessageString(data="other"))

    store.delete("cam_a/latest/**")
    self.assertCountEqual(
        store.list_all_keys(), ["kv_store/cam_a/latest_other"]
    )

  def test_delete_mid_path_and_multiple_double_star_wildcards(self):
    store = fake_pubsub.FakeKeyValueStore()
    for key in (
        "cam_a/rgb",
        "cam_a/0/rgb",
        "cam_a/0/left/rgb",
        "cam_a/0/depth",
        "a/x/y/z/leaf",
        "a/leaf",
        "a/x/y/other",
    ):
      store.set(key, test_pb2.TestMessageString(data=key))

    store.delete("cam_a/**/rgb")
    store.delete("a/**/**/missing")
    store.delete("a/**/**/leaf")

    self.assertCountEqual(
        store.list_all_keys(),
        ["kv_store/cam_a/0/depth", "kv_store/a/x/y/other"],
    )

  def test_delete_single_segment_and_subsegment_wildcards(self):
    # Each case gets a fresh store so delete patterns cannot mask over-deletion
    # across shared key prefixes.
    cases = [
        (
            "leading_subsegment_wildcard",
            "cam_a/$*_suffix",
            ("cam_a/_suffix", "cam_a/frame_suffix"),
            ("cam_a/frame_other",),
        ),
        (
            "trailing_subsegment_wildcard",
            "cam_a/prefix_$*",
            ("cam_a/prefix_", "cam_a/prefix_001"),
            ("cam_a/other_001",),
        ),
        (
            "consecutive_subsegment_wildcards",
            "grp/a$*$*c",
            ("grp/ac", "grp/axbyc"),
            ("grp/axd",),
        ),
        (
            "mid_subsegment_and_single_segment_wildcard",
            "cam_a/frame_$*2/*",
            (
                "cam_a/frame_002/rgb",
                "cam_a/frame_002/depth",
                "cam_a/frame_2/rgb",
            ),
            ("cam_a/frame_001/rgb", "cam_a/meta_002/rgb"),
        ),
        (
            "trailing_digit_subsegment_wildcard",
            "cam_a/item_$*9",
            ("cam_a/item_9", "cam_a/item_abc9"),
            ("cam_a/item_10",),
        ),
        (
            "standalone_dollar_star_segment",
            "cam_a/$*/meta",
            ("cam_a/any_seg/meta",),
            ("cam_a/a/b/meta",),
        ),
        (
            "multiple_literals_between_wildcards",
            "grp/a$*b$*c",
            ("grp/abc", "grp/axbyc"),
            ("grp/axc",),
        ),
        (
            "overlapping_prefix_and_suffix",
            "ovl/ab$*bc",
            ("ovl/abbc", "ovl/abXYZbc"),
            ("ovl/abc",),
        ),
        (
            "overlapping_middle_and_suffix",
            "ovl2/a$*b$*bc",
            ("ovl2/abbc", "ovl2/axbybc"),
            ("ovl2/abc",),
        ),
        (
            "repeated_middle_literals",
            "sbr/a$*b$*b$*c",
            ("sbr/abbc", "sbr/axbbyc"),
            ("sbr/abc",),
        ),
        (
            "single_segment_star_does_not_cross_slash",
            "seg/*",
            ("seg/x",),
            ("seg/x/y",),
        ),
        (
            # Guards linear-time greedy matching against exponential regex
            # backtracking on repeated `$*` wildcards.
            "repeated_subsegment_wildcards_no_backtracking_explosion",
            "ovl/" + ("a$*" * 25) + "b",
            ("ovl/" + ("a" * 25) + "b",),
            ("ovl/" + ("a" * 60),),
        ),
    ]
    for case_name, pattern, matching, non_matching in cases:
      with self.subTest(case=case_name, pattern=pattern):
        store = fake_pubsub.FakeKeyValueStore()
        for key in (*matching, *non_matching):
          store.set(key, test_pb2.TestMessageString(data=key))

        store.delete(pattern)

        self.assertCountEqual(
            store.list_all_keys(),
            [f"kv_store/{key}" for key in non_matching],
        )

  def test_concurrent_set_delete_and_list_all_keys(self):
    store = fake_pubsub.FakeKeyValueStore()
    errors: list[Exception] = []
    barrier = threading.Barrier(3, timeout=10.0)

    def writer() -> None:
      try:
        barrier.wait()
        for i in range(100):
          store.set(f"tmp/{i}", test_pb2.TestMessageString(data=str(i)))
          store.set(f"keep/{i}", test_pb2.TestMessageString(data=str(i)))
      except Exception as exc:  # pylint: disable=broad-exception-caught
        errors.append(exc)

    def deleter() -> None:
      try:
        barrier.wait()
        for _ in range(50):
          store.delete("tmp/*")
      except Exception as exc:  # pylint: disable=broad-exception-caught
        errors.append(exc)

    def lister() -> None:
      try:
        barrier.wait()
        for _ in range(100):
          _ = store.list_all_keys()
      except Exception as exc:  # pylint: disable=broad-exception-caught
        errors.append(exc)

    threads = [
        threading.Thread(target=writer, daemon=True),
        threading.Thread(target=deleter, daemon=True),
        threading.Thread(target=lister, daemon=True),
    ]
    for thread in threads:
      thread.start()
    for thread in threads:
      thread.join(timeout=10.0)
      self.assertFalse(thread.is_alive())

    self.assertEmpty(errors)
    store.delete("tmp/*")
    self.assertCountEqual(
        store.list_all_keys(), [f"kv_store/keep/{i}" for i in range(100)]
    )


class FakePubSubTest(absltest.TestCase):

  def test_constructor_overloads_and_destroy_session_when_unused(self):
    ps_default = fake_pubsub.FakePubSub()
    ps_named = fake_pubsub.FakePubSub("participant")
    ps_configured = fake_pubsub.FakePubSub("participant", "{}")

    for ps in (ps_default, ps_named, ps_configured):
      ps.key_value_store().set("k", test_pb2.TestMessageString(data="v"))
      self.assertCountEqual(
          ps.key_value_store().list_all_keys(), ["kv_store/k"]
      )

    self.assertIsNone(fake_pubsub.FakePubSub.destroy_session_when_unused())
    self.assertIsNone(ps_default.destroy_session_when_unused())

  def test_default_store_uses_kv_store_prefix_and_rejects_empty_prefix(self):
    pubsub = fake_pubsub.FakePubSub()
    pubsub.key_value_store().set("key", test_pb2.TestMessageString(data="x"))

    self.assertIs(pubsub.key_value_store(), pubsub.key_value_store(None))
    self.assertIs(pubsub.key_value_store(), pubsub.key_value_store("kv_store"))
    self.assertCountEqual(
        pubsub.key_value_store().list_all_keys(), ["kv_store/key"]
    )

    empty_prefix_store = pubsub.key_value_store("")
    self.assertIsNot(pubsub.key_value_store(), empty_prefix_store)
    with _assert_raises_status(self, status.StatusCode.INVALID_ARGUMENT):
      empty_prefix_store.set("key", test_pb2.TestMessageString(data="x"))
    with _assert_raises_status(self, status.StatusCode.INVALID_ARGUMENT):
      empty_prefix_store.set_raw("key", b"x")
    with _assert_raises_status(self, status.StatusCode.INVALID_ARGUMENT):
      empty_prefix_store.list_all_keys()

  def test_replication_store_list_all_keys_raises_unimplemented(self):
    pubsub = fake_pubsub.FakePubSub()
    repl_store = pubsub.replication_key_value_store()
    self.assertIs(repl_store, pubsub.key_value_store("kv_store_repl"))
    repl_store.set("global/item", test_pb2.TestMessageString(data="v"))

    with _assert_raises_status(self, status.StatusCode.UNIMPLEMENTED):
      repl_store.list_all_keys()

  def test_named_stores_are_isolated_and_prefix_listed_keys(self):
    pubsub = fake_pubsub.FakePubSub()
    pubsub.key_value_store("alpha").set(
        "key", test_pb2.TestMessageString(data="a")
    )

    self.assertCountEqual(
        pubsub.key_value_store("alpha").list_all_keys(), ["alpha/key"]
    )
    self.assertEmpty(pubsub.key_value_store("beta").list_all_keys())
    self.assertEmpty(pubsub.key_value_store().list_all_keys())

  def test_same_name_returns_same_store(self):
    pubsub = fake_pubsub.FakePubSub()
    self.assertIs(
        pubsub.key_value_store("alpha"), pubsub.key_value_store("alpha")
    )


if __name__ == "__main__":
  absltest.main()
