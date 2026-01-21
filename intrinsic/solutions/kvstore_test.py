# Copyright 2023 Intrinsic Innovation LLC

"""Tests for kvstore.py."""

from unittest import mock

from absl.testing import absltest
from google.protobuf import any_pb2
from google.protobuf import wrappers_pb2
import grpc

from intrinsic.platform.pubsub.kvstore_grpc import kvstore_pb2
from intrinsic.platform.pubsub.kvstore_grpc import kvstore_pb2_grpc
from intrinsic.solutions import kvstore


class KVStoreTest(absltest.TestCase):

  def setUp(self):
    super().setUp()
    self.stub = mock.create_autospec(
        kvstore_pb2_grpc.KVStoreStub,
        instance=True,
    )
    self.stub.Get = mock.MagicMock()
    self.stub.Set = mock.MagicMock()
    self.stub.Delete = mock.MagicMock()
    self.stub.List = mock.MagicMock()
    self.kv_store = kvstore.KVStore(self.stub)

  def test_get_success(self):
    expected_value = wrappers_pb2.StringValue(value="test_value")
    any_value = any_pb2.Any()
    any_value.Pack(expected_value)

    self.stub.Get.return_value = kvstore_pb2.GetResponse(value=any_value)

    result = self.kv_store.get("test_key")

    self.stub.Get.assert_called_once()
    self.assertEqual(result, any_value)

    unpacked = wrappers_pb2.StringValue()
    result.Unpack(unpacked)
    self.assertEqual(unpacked, expected_value)

  def test_set_success(self):
    value = wrappers_pb2.StringValue(value="test_value")
    any_value = any_pb2.Any()
    any_value.Pack(value)
    self.stub.Set.return_value = kvstore_pb2.SetResponse()

    self.kv_store.set("test_key", any_value)

    self.stub.Set.assert_called_once()
    args = self.stub.Set.call_args[0][0]
    self.assertEqual(args.key, "test_key")
    self.assertEqual(args.value, any_value)

  def test_delete_success(self):
    self.stub.Delete.return_value = kvstore_pb2.DeleteResponse()

    self.kv_store.delete("test_key")

    self.stub.Delete.assert_called_once()
    args = self.stub.Delete.call_args[0][0]
    self.assertEqual(args.key, "test_key")

  def test_keys_success(self):
    self.stub.List.return_value = kvstore_pb2.ListResponse(keys=["k1", "k2"])

    keys = self.kv_store.keys()

    self.stub.List.assert_called_once()
    self.assertEqual(keys, ["k1", "k2"])

  def test_get_not_found(self):
    state = grpc.RpcError()
    state.code = lambda: grpc.StatusCode.NOT_FOUND
    self.stub.Get.side_effect = state

    result = self.kv_store.get("test_key")

    self.stub.Get.assert_called_once()
    self.assertIsNone(result)

  def test_get_failure(self):
    state = grpc.RpcError()
    state.code = lambda: grpc.StatusCode.INTERNAL
    self.stub.Get.side_effect = state

    with self.assertRaises(grpc.RpcError):
      self.kv_store.get("test_key")

    self.stub.Get.assert_called_once()

  def test_set_failure(self):
    state = grpc.RpcError()
    state.code = lambda: grpc.StatusCode.INTERNAL
    self.stub.Set.side_effect = state

    value = wrappers_pb2.StringValue(value="test_value")
    any_value = any_pb2.Any()
    any_value.Pack(value)

    with self.assertRaises(grpc.RpcError):
      self.kv_store.set("test_key", any_value)

    self.stub.Set.assert_called_once()

  def test_delete_failure(self):
    state = grpc.RpcError()
    state.code = lambda: grpc.StatusCode.INTERNAL
    self.stub.Delete.side_effect = state

    with self.assertRaises(grpc.RpcError):
      self.kv_store.delete("test_key")

    self.stub.Delete.assert_called_once()

  def test_keys_failure(self):
    state = grpc.RpcError()
    state.code = lambda: grpc.StatusCode.INTERNAL
    self.stub.List.side_effect = state

    with self.assertRaises(grpc.RpcError):
      self.kv_store.keys()

    self.stub.List.assert_called_once()


if __name__ == "__main__":
  absltest.main()
