# Copyright 2023 Intrinsic Innovation LLC

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


if __name__ == "__main__":
  absltest.main()
