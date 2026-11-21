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


@contextlib.contextmanager
def _assert_raises_status(
    test_case: absltest.TestCase, expected_code: status.StatusCode
):
  with test_case.assertRaises(status.StatusNotOk) as ctx:
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
    with self.assertRaisesRegex(status.StatusNotOk, "missing") as ctx:
      store.get("missing")
    self.assertEqual(ctx.exception.status.code(), status.StatusCode.NOT_FOUND)

  def test_signature_compatibility_keyword_arguments(self):
    pubsub = fake_pubsub.FakePubSub(
        participant_name="test_participant", config="{}"
    )
    store = pubsub.key_value_store(prefix_override="custom_store")
    msg = test_pb2.TestMessageString(data="val")
    verified_msg = test_pb2.TestMessageString(data="verified")
    store.set(key="key", value=msg, high_consistency=True)
    store.set_with_verification(
        key="verified_key",
        value=verified_msg,
        options=types.SimpleNamespace(timeout=5.0),
    )

    self.assertCountEqual(
        store.list_all_keys(timeout=0),
        ["custom_store/key", "custom_store/verified_key"],
    )
    unpacked = test_pb2.TestMessageString()
    self.assertTrue(store.get(key="key", timeout=0).Unpack(unpacked))
    self.assertEqual(unpacked, msg)
    self.assertTrue(store.get(key="verified_key", timeout=5).Unpack(unpacked))
    self.assertEqual(unpacked, verified_msg)

    store.delete(key="key")
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
          store.set_with_verification(bad_key, msg)
        with _assert_raises_status(self, invalid_arg):
          store.get(bad_key)
        with _assert_raises_status(self, invalid_arg):
          store.delete(bad_key)

    with _assert_raises_status(self, invalid_arg):
      store.get("valid/key", timeout=-1)

    with _assert_raises_status(self, invalid_arg):
      store.list_all_keys(timeout=-1)

    for bad_timeout in (0.0, -1.0):
      with self.subTest(bad_timeout=bad_timeout):
        with _assert_raises_status(self, invalid_arg):
          store.set_with_verification(
              "valid/key",
              msg,
              options=types.SimpleNamespace(timeout=bad_timeout),
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
