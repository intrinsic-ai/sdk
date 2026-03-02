# Copyright 2023 Intrinsic Innovation LLC

"""Tests for blackboard.py."""

from unittest import mock

from absl.testing import absltest
from google.protobuf import any_pb2
from google.protobuf import wrappers_pb2

from intrinsic.executive.proto import blackboard_service_pb2
from intrinsic.executive.proto import test_message_pb2
from intrinsic.solutions import blackboard
from intrinsic.solutions import blackboard_value
from intrinsic.solutions.internal import skill_utils


class BlackboardTest(absltest.TestCase):

  def setUp(self):
    super().setUp()
    self._stub = mock.MagicMock()
    self._blackboard = blackboard.Blackboard(self._stub, "test_operation")

  def test_delete_value_without_scope(self):
    self._blackboard.delete_value("test_key")

    self._stub.DeleteBlackboardValue.assert_called_once_with(
        blackboard_service_pb2.DeleteBlackboardValueRequest(
            key="test_key",
            scope="",
            operation_name="test_operation",
        )
    )

  def test_delete_value_with_scope(self):
    self._blackboard.delete_value("test_key", scope="test_scope")

    self._stub.DeleteBlackboardValue.assert_called_once_with(
        blackboard_service_pb2.DeleteBlackboardValueRequest(
            key="test_key",
            scope="test_scope",
            operation_name="test_operation",
        )
    )

  def test_list_keys_without_scope(self):
    self._stub.ListBlackboardValues.return_value = (
        blackboard_service_pb2.ListBlackboardValuesResponse(
            values=[
                blackboard_service_pb2.BlackboardValue(
                    key="key1",
                    scope="scope1",
                    value=any_pb2.Any(type_url="type.googleapis.com/test1"),
                ),
                blackboard_service_pb2.BlackboardValue(
                    key="key2",
                    scope="scope2",
                    value=any_pb2.Any(type_url="type.googleapis.com/test2"),
                ),
            ]
        )
    )

    result = self._blackboard.list_keys()

    self.assertEqual(
        result,
        [
            blackboard.ScopedBlackboardKey(
                key="key1",
                scope="scope1",
                type_url="type.googleapis.com/test1",
            ),
            blackboard.ScopedBlackboardKey(
                key="key2",
                scope="scope2",
                type_url="type.googleapis.com/test2",
            ),
        ],
    )
    self._stub.ListBlackboardValues.assert_called_once_with(
        blackboard_service_pb2.ListBlackboardValuesRequest(
            operation_name="test_operation",
            scope=None,
            view=blackboard_service_pb2.ListBlackboardValuesRequest.ANY_TYPEURL_ONLY,
        )
    )

  def test_list_keys_with_scope(self):
    self._stub.ListBlackboardValues.return_value = (
        blackboard_service_pb2.ListBlackboardValuesResponse(
            values=[
                blackboard_service_pb2.BlackboardValue(
                    key="key1",
                    scope="test_scope",
                    value=any_pb2.Any(type_url="type.googleapis.com/test1"),
                ),
            ]
        )
    )

    result = self._blackboard.list_keys(scope="test_scope")

    self.assertEqual(
        result,
        [
            blackboard.ScopedBlackboardKey(
                key="key1",
                scope="test_scope",
                type_url="type.googleapis.com/test1",
            )
        ],
    )
    self._stub.ListBlackboardValues.assert_called_once_with(
        blackboard_service_pb2.ListBlackboardValuesRequest(
            operation_name="test_operation",
            scope="test_scope",
            view=blackboard_service_pb2.ListBlackboardValuesRequest.ANY_TYPEURL_ONLY,
        )
    )

  def test_get_value_any(self):
    val = wrappers_pb2.UInt64Value(value=42)
    any_val = any_pb2.Any()
    any_val.Pack(val)
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=any_val)
    )

    result = self._blackboard.get_value_any("test_key")

    self.assertIsInstance(result, any_pb2.Any)
    self.assertEqual(result.type_url, any_val.type_url)
    self.assertEqual(result.value, any_val.value)
    self._stub.GetBlackboardValue.assert_called_once_with(
        blackboard_service_pb2.GetBlackboardValueRequest(
            key="test_key",
            scope=None,
            operation_name="test_operation",
        )
    )

  def test_get_value_uint64(self):
    val = wrappers_pb2.UInt64Value(value=42)
    any_val = any_pb2.Any()
    any_val.Pack(val)
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=any_val)
    )

    result = self._blackboard.get_value("test_key")

    self.assertEqual(result, 42)
    self._stub.GetBlackboardValue.assert_called_once_with(
        blackboard_service_pb2.GetBlackboardValueRequest(
            key="test_key",
            scope=None,
            operation_name="test_operation",
        )
    )

  def test_get_value_string(self):
    val = wrappers_pb2.StringValue(value="hello")
    any_val = any_pb2.Any()
    any_val.Pack(val)
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=any_val)
    )

    result = self._blackboard.get_value("test_key")

    self.assertEqual(result, "hello")

  def test_get_value_with_blackboard_value(self):
    mock_bv = mock.MagicMock(spec=blackboard_value.BlackboardValue)
    mock_bv.value_access_path.return_value = "bv_key"
    mock_bv.scope.return_value = "bv_scope"
    mock_bv.is_toplevel_value = True

    val = wrappers_pb2.StringValue(value="hello")
    any_val = any_pb2.Any()
    any_val.Pack(val)
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=any_val)
    )

    result = self._blackboard.get_value(mock_bv)

    self.assertEqual(result, "hello")
    self._stub.GetBlackboardValue.assert_called_once_with(
        blackboard_service_pb2.GetBlackboardValueRequest(
            key="bv_key",
            scope="bv_scope",
            operation_name="test_operation",
        )
    )

  def test_get_value_with_blackboard_value_and_explicit_scope(self):
    mock_bv = mock.MagicMock(spec=blackboard_value.BlackboardValue)
    mock_bv.value_access_path.return_value = "bv_key"
    mock_bv.scope.return_value = "bv_scope"
    mock_bv.is_toplevel_value = True

    with self.assertRaisesRegex(ValueError, "Cannot provide explicit scope"):
      self._blackboard.get_value(mock_bv, scope="explicit_scope")

  def test_resolve_key_fails_for_non_toplevel_bv(self):
    mock_bv = mock.MagicMock(spec=blackboard_value.BlackboardValue)
    mock_bv.value_access_path.return_value = "bv_key.sub_field"
    mock_bv.is_toplevel_value = False

    with self.assertRaisesRegex(ValueError, "not a toplevel value"):
      self._blackboard.get_value(mock_bv)

  def test_get_value_any_fallback(self):
    val = test_message_pb2.TestMessage(int32_value=123)
    any_val = any_pb2.Any()
    any_val.Pack(val)
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=any_val)
    )

    result = self._blackboard.get_value("test_key")

    self.assertIsInstance(result, any_pb2.Any)
    self.assertEqual(result.type_url, any_val.type_url)
    self.assertEqual(result.value, any_val.value)

  def test_update_value_any(self):
    existing_any = any_pb2.Any()
    existing_any.Pack(wrappers_pb2.StringValue(value="existing"))
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=existing_any)
    )

    any_val = any_pb2.Any()
    any_val.Pack(wrappers_pb2.StringValue(value="test"))

    self._blackboard.update_value("test_key", any_val)

    self._stub.UpdateBlackboardValue.assert_called_once_with(
        blackboard_service_pb2.UpdateBlackboardValueRequest(
            value=blackboard_service_pb2.BlackboardValue(
                key="test_key",
                scope="",
                operation_name="test_operation",
                value=any_val,
            )
        )
    )

  def test_update_value_message(self):
    existing_any = any_pb2.Any()
    existing_any.Pack(test_message_pb2.TestMessage())
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=existing_any)
    )

    msg = test_message_pb2.TestMessage(int32_value=123)
    any_val = any_pb2.Any()
    any_val.Pack(msg)

    self._blackboard.update_value("test_key", msg, scope="test_scope")

    self._stub.UpdateBlackboardValue.assert_called_once_with(
        blackboard_service_pb2.UpdateBlackboardValueRequest(
            value=blackboard_service_pb2.BlackboardValue(
                key="test_key",
                scope="test_scope",
                operation_name="test_operation",
                value=any_val,
            )
        )
    )

  def test_update_value_message_wrapper(self):
    existing_any = any_pb2.Any()
    existing_any.Pack(wrappers_pb2.StringValue(value="existing"))
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=existing_any)
    )

    mock_wrapper = mock.MagicMock(spec=skill_utils.MessageWrapper)
    any_val = any_pb2.Any()
    any_val.Pack(wrappers_pb2.StringValue(value="wrapped"))
    mock_wrapper.to_any.return_value = any_val

    self._blackboard.update_value("test_key", mock_wrapper)

    self._stub.UpdateBlackboardValue.assert_called_once_with(
        blackboard_service_pb2.UpdateBlackboardValueRequest(
            value=blackboard_service_pb2.BlackboardValue(
                key="test_key",
                scope="",
                operation_name="test_operation",
                value=any_val,
            )
        )
    )

  def test_update_value_int_matching_uint64(self):
    existing_any = any_pb2.Any()
    existing_any.Pack(wrappers_pb2.UInt64Value(value=1))
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=existing_any)
    )

    any_val = any_pb2.Any()
    any_val.Pack(wrappers_pb2.UInt64Value(value=123))

    self._blackboard.update_value("test_key", 123)

    self._stub.UpdateBlackboardValue.assert_called_once_with(
        blackboard_service_pb2.UpdateBlackboardValueRequest(
            value=blackboard_service_pb2.BlackboardValue(
                key="test_key",
                scope="",
                operation_name="test_operation",
                value=any_val,
            )
        )
    )

  def test_update_value_int_matching_int64(self):
    existing_any = any_pb2.Any()
    existing_any.Pack(wrappers_pb2.Int64Value(value=1))
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=existing_any)
    )

    any_val = any_pb2.Any()
    any_val.Pack(wrappers_pb2.Int64Value(value=123))

    # Should match Int64Value even if positive
    self._blackboard.update_value("test_key", 123)

    self._stub.UpdateBlackboardValue.assert_called_once_with(
        blackboard_service_pb2.UpdateBlackboardValueRequest(
            value=blackboard_service_pb2.BlackboardValue(
                key="test_key",
                scope="",
                operation_name="test_operation",
                value=any_val,
            )
        )
    )

  def test_update_value_type_mismatch(self):
    existing_any = any_pb2.Any()
    existing_any.Pack(wrappers_pb2.StringValue(value="existing"))
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=existing_any)
    )

    with self.assertRaisesRegex(TypeError, "Type mismatch"):
      self._blackboard.update_value("test_key", 123)

  def test_update_value_float_type_mismatch(self):
    existing_any = any_pb2.Any()
    existing_any.Pack(wrappers_pb2.StringValue(value="existing"))
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=existing_any)
    )

    with self.assertRaisesRegex(TypeError, "Type mismatch"):
      self._blackboard.update_value("test_key", 1.23)

  def test_update_value_negative_int(self):
    existing_any = any_pb2.Any()
    existing_any.Pack(wrappers_pb2.Int64Value(value=-1))
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=existing_any)
    )

    any_val = any_pb2.Any()
    any_val.Pack(wrappers_pb2.Int64Value(value=-123))

    self._blackboard.update_value("test_key", -123)

    self._stub.UpdateBlackboardValue.assert_called_once_with(
        blackboard_service_pb2.UpdateBlackboardValueRequest(
            value=blackboard_service_pb2.BlackboardValue(
                key="test_key",
                scope="",
                operation_name="test_operation",
                value=any_val,
            )
        )
    )

  def test_update_value_float(self):
    existing_any = any_pb2.Any()
    existing_any.Pack(wrappers_pb2.DoubleValue(value=1.0))
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=existing_any)
    )

    any_val = any_pb2.Any()
    any_val.Pack(wrappers_pb2.DoubleValue(value=1.23))

    self._blackboard.update_value("test_key", 1.23)

    self._stub.UpdateBlackboardValue.assert_called_once_with(
        blackboard_service_pb2.UpdateBlackboardValueRequest(
            value=blackboard_service_pb2.BlackboardValue(
                key="test_key",
                scope="",
                operation_name="test_operation",
                value=any_val,
            )
        )
    )

  def test_update_value_float_matching_float_value(self):
    existing_any = any_pb2.Any()
    existing_any.Pack(wrappers_pb2.FloatValue(value=1.0))
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=existing_any)
    )

    any_val = any_pb2.Any()
    any_val.Pack(wrappers_pb2.FloatValue(value=1.23))

    self._blackboard.update_value("test_key", 1.23)

    self._stub.UpdateBlackboardValue.assert_called_once_with(
        blackboard_service_pb2.UpdateBlackboardValueRequest(
            value=blackboard_service_pb2.BlackboardValue(
                key="test_key",
                scope="",
                operation_name="test_operation",
                value=any_val,
            )
        )
    )

  def test_update_value_bool(self):
    existing_any = any_pb2.Any()
    existing_any.Pack(wrappers_pb2.BoolValue(value=False))
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=existing_any)
    )

    any_val = any_pb2.Any()
    any_val.Pack(wrappers_pb2.BoolValue(value=True))

    self._blackboard.update_value("test_key", True)

    self._stub.UpdateBlackboardValue.assert_called_once_with(
        blackboard_service_pb2.UpdateBlackboardValueRequest(
            value=blackboard_service_pb2.BlackboardValue(
                key="test_key",
                scope="",
                operation_name="test_operation",
                value=any_val,
            )
        )
    )

  def test_update_value_str(self):
    existing_any = any_pb2.Any()
    existing_any.Pack(wrappers_pb2.StringValue(value="existing"))
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=existing_any)
    )

    any_val = any_pb2.Any()
    any_val.Pack(wrappers_pb2.StringValue(value="test"))

    self._blackboard.update_value("test_key", "test")

    self._stub.UpdateBlackboardValue.assert_called_once_with(
        blackboard_service_pb2.UpdateBlackboardValueRequest(
            value=blackboard_service_pb2.BlackboardValue(
                key="test_key",
                scope="",
                operation_name="test_operation",
                value=any_val,
            )
        )
    )

  def test_update_value_bytes(self):
    existing_any = any_pb2.Any()
    existing_any.Pack(wrappers_pb2.BytesValue(value=b"existing"))
    self._stub.GetBlackboardValue.return_value = (
        blackboard_service_pb2.BlackboardValue(value=existing_any)
    )

    any_val = any_pb2.Any()
    any_val.Pack(wrappers_pb2.BytesValue(value=b"test"))

    self._blackboard.update_value("test_key", b"test")

    self._stub.UpdateBlackboardValue.assert_called_once_with(
        blackboard_service_pb2.UpdateBlackboardValueRequest(
            value=blackboard_service_pb2.BlackboardValue(
                key="test_key",
                scope="",
                operation_name="test_operation",
                value=any_val,
            )
        )
    )

  def test_update_value_invalid_type(self):
    with self.assertRaises(TypeError):
      self._blackboard.update_value("test_key", [1, 2, 3])


if __name__ == "__main__":
  absltest.main()
