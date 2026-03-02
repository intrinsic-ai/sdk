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
from intrinsic.util.status import extended_status_pb2


class BlackboardTest(absltest.TestCase):

  def setUp(self):
    super().setUp()
    self._stub = mock.MagicMock()
    self._blackboard = blackboard.Blackboard(self._stub, "test_operation")

    # Default success response for load_snapshot
    success_diagnostics = extended_status_pb2.ExtendedStatus()
    success_diagnostics.status_code.code = 80114
    success_diagnostics.severity = (
        extended_status_pb2.ExtendedStatus.Severity.INFO
    )
    self._stub.LoadBlackboardSnapshot.return_value = (
        blackboard_service_pb2.LoadBlackboardSnapshotResponse(
            diagnostics=success_diagnostics
        )
    )

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


  def test_create_snapshot(self):
    expected_snapshot = blackboard_service_pb2.BlackboardSnapshot(
        handle="new_handle", display_name="new_name"
    )
    self._stub.CreateBlackboardSnapshot.return_value = (
        blackboard_service_pb2.CreateBlackboardSnapshotResponse(
            snapshot=expected_snapshot
        )
    )

    result = self._blackboard.create_snapshot("new_name")

    self.assertEqual(result, expected_snapshot)
    self._stub.CreateBlackboardSnapshot.assert_called_once_with(
        blackboard_service_pb2.CreateBlackboardSnapshotRequest(
            operation_name="test_operation",
            display_name="new_name",
            snapshot_source=blackboard.SnapshotSource.USER.value,
        )
    )

  def test_create_snapshot_default_name(self):
    self._stub.CreateBlackboardSnapshot.return_value = (
        blackboard_service_pb2.CreateBlackboardSnapshotResponse(
            snapshot=blackboard_service_pb2.BlackboardSnapshot(handle="h")
        )
    )

    self._blackboard.create_snapshot()

    self._stub.CreateBlackboardSnapshot.assert_called_once_with(
        blackboard_service_pb2.CreateBlackboardSnapshotRequest(
            operation_name="test_operation",
            display_name="",
            snapshot_source=blackboard.SnapshotSource.USER.value,
        )
    )

  def test_load_snapshot(self):
    self._blackboard.load_snapshot("test_handle")

    self._stub.LoadBlackboardSnapshot.assert_called_once_with(
        blackboard_service_pb2.LoadBlackboardSnapshotRequest(
            operation_name="test_operation",
            handle="test_handle",
            integration_mode=blackboard.IntegrationMode.MERGE.value,
        )
    )

  def test_load_snapshot_with_proto(self):
    snapshot = blackboard_service_pb2.BlackboardSnapshot(handle="proto_handle")
    self._blackboard.load_snapshot(snapshot)

    self._stub.LoadBlackboardSnapshot.assert_called_once_with(
        blackboard_service_pb2.LoadBlackboardSnapshotRequest(
            operation_name="test_operation",
            handle="proto_handle",
            integration_mode=blackboard.IntegrationMode.MERGE.value,
        )
    )

  @mock.patch(
      "intrinsic.solutions.ipython.display_extended_status_proto_if_ipython"
  )
  def test_load_snapshot_with_diagnostics(self, mock_display):
    diagnostics = extended_status_pb2.ExtendedStatus()
    diagnostics.status_code.code = 12345
    diagnostics.severity = extended_status_pb2.ExtendedStatus.Severity.ERROR
    self._stub.LoadBlackboardSnapshot.return_value = (
        blackboard_service_pb2.LoadBlackboardSnapshotResponse(
            diagnostics=diagnostics
        )
    )

    self._blackboard.load_snapshot("test_handle")

    mock_display.assert_called_once_with(diagnostics)


class BlackboardSnapshotsTest(absltest.TestCase):

  def setUp(self):
    super().setUp()
    self._stub = mock.MagicMock()
    self._snapshots = blackboard.BlackboardSnapshots(self._stub)

  def test_list(self):
    self._stub.ListBlackboardSnapshots.return_value = (
        blackboard_service_pb2.ListBlackboardSnapshotsResponse(
            snapshots=[
                blackboard_service_pb2.BlackboardSnapshot(
                    handle="h1", display_name="s1"
                ),
                blackboard_service_pb2.BlackboardSnapshot(
                    handle="h2", display_name="s2"
                ),
            ]
        )
    )

    result = self._snapshots.list()

    self.assertEqual(len(result), 2)
    self.assertEqual(result[0].handle, "h1")
    self.assertEqual(result[1].handle, "h2")
    self._stub.ListBlackboardSnapshots.assert_called_once()

  def test_list_paginated(self):
    self._stub.ListBlackboardSnapshots.side_effect = [
        blackboard_service_pb2.ListBlackboardSnapshotsResponse(
            snapshots=[
                blackboard_service_pb2.BlackboardSnapshot(
                    handle="h1", display_name="s1"
                )
            ],
            next_page_token="token2",
        ),
        blackboard_service_pb2.ListBlackboardSnapshotsResponse(
            snapshots=[
                blackboard_service_pb2.BlackboardSnapshot(
                    handle="h2", display_name="s2"
                )
            ],
            next_page_token="",
        ),
    ]

    result = self._snapshots.list()

    self.assertEqual(len(result), 2)
    self.assertEqual(result[0].handle, "h1")
    self.assertEqual(result[1].handle, "h2")
    self.assertEqual(self._stub.ListBlackboardSnapshots.call_count, 2)
    self._stub.ListBlackboardSnapshots.assert_has_calls([
        mock.call(
            blackboard_service_pb2.ListBlackboardSnapshotsRequest(
                page_size=100, page_token=""
            )
        ),
        mock.call(
            blackboard_service_pb2.ListBlackboardSnapshotsRequest(
                page_size=100, page_token="token2"
            )
        ),
    ])

  def test_delete(self):
    self._snapshots.delete("test_handle")

    self._stub.DeleteBlackboardSnapshot.assert_called_once_with(
        blackboard_service_pb2.DeleteBlackboardSnapshotRequest(
            handle="test_handle"
        )
    )

  def test_delete_with_proto(self):
    snapshot = blackboard_service_pb2.BlackboardSnapshot(handle="proto_handle")
    self._snapshots.delete(snapshot)

    self._stub.DeleteBlackboardSnapshot.assert_called_once_with(
        blackboard_service_pb2.DeleteBlackboardSnapshotRequest(
            handle="proto_handle"
        )
    )




if __name__ == "__main__":
  absltest.main()
