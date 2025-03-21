# Copyright 2023 Intrinsic Innovation LLC

import time
from unittest import mock

from absl.testing import absltest
import grpc
from intrinsic.util.grpc import error_handling


class _GrpcError(grpc.RpcError, grpc.Call):
  """Helper class that emulates a gRPC error."""

  def __init__(self, code: int):
    self._code = code

  def code(self) -> int:
    return self._code

  def details(self):
    return '_GrpcError'


@error_handling.retry_on_grpc_unavailable
def _call_with_unavailable_retry(stub) -> str:
  return stub.call()


@error_handling.retry_on_grpc_resource_exhausted
def _call_with_resource_exhausted_retry(stub) -> str:
  return stub.call()


@error_handling.retry_on_grpc_transient_errors
def _call_with_transient_errors_retry(stub) -> str:
  return stub.call()


class ErrorsTest(absltest.TestCase):

  @mock.patch.object(time, 'sleep')
  def test_retry_on_grpc_unavailable_retries_on_certain_errors(self, _):
    stub = mock.MagicMock()
    stub.call.side_effect = [
        _GrpcError(grpc.StatusCode.UNAVAILABLE),
        _GrpcError(grpc.StatusCode.UNIMPLEMENTED),
        'some result',
    ]
    result = _call_with_unavailable_retry(stub)
    stub.call.assert_called_with()
    self.assertEqual(result, 'some result')

  @mock.patch.object(time, 'sleep')
  def test_retry_on_grpc_unavailable_fails_after_max_retries(self, mock_sleep):
    stub = mock.MagicMock()
    # We receive an UNAVAILABLE error when we directly try to contact a grpc
    # server that is not (yet) running.
    stub.call.side_effect = _GrpcError(grpc.StatusCode.UNAVAILABLE)
    with self.assertRaises(_GrpcError) as context:
      _call_with_unavailable_retry(stub)
    self.assertEqual(context.exception.code(), grpc.StatusCode.UNAVAILABLE)
    # Stub gets called for the max number of attempts.
    self.assertEqual(stub.call.call_count, 15)
    mock_sleep.assert_has_calls([mock.call(mock.ANY)])

  def test_retry_on_grpc_unavailable_does_not_retry_on_other_grpc_error(self):
    stub = mock.MagicMock()
    stub.call.side_effect = _GrpcError(grpc.StatusCode.INVALID_ARGUMENT)
    with self.assertRaises(Exception) as context:
      _call_with_unavailable_retry(stub)
    self.assertEqual(context.exception.code(), grpc.StatusCode.INVALID_ARGUMENT)
    self.assertEqual(stub.call.call_count, 1)

  def test_retry_on_grpc_unavailable_does_not_retry_on_non_grpc_error(self):
    stub = mock.MagicMock()
    stub.call.side_effect = Exception('non-grpc error')
    with self.assertRaises(Exception) as context:
      _call_with_unavailable_retry(stub)
    self.assertEqual(str(context.exception), 'non-grpc error')
    self.assertEqual(stub.call.call_count, 1)

  @mock.patch.object(time, 'sleep')
  def test_retry_on_grpc_resource_exhausted_retries_on_certain_errors(self, _):
    stub = mock.MagicMock()
    stub.call.side_effect = [
        _GrpcError(grpc.StatusCode.RESOURCE_EXHAUSTED),
        'some result',
    ]
    result = _call_with_resource_exhausted_retry(stub)
    stub.call.assert_called_with()
    self.assertEqual(result, 'some result')

  @mock.patch.object(time, 'sleep')
  def test_retry_on_grpc_resource_exhausted_fails_after_max_retries(
      self, mock_sleep
  ):
    stub = mock.MagicMock()
    # We receive an RESOURCE_EXHAUSTED error when we directly try to contact a
    # grpc server that is not (yet) running.
    stub.call.side_effect = _GrpcError(grpc.StatusCode.RESOURCE_EXHAUSTED)
    with self.assertRaises(_GrpcError) as context:
      _call_with_resource_exhausted_retry(stub)
    self.assertEqual(
        context.exception.code(), grpc.StatusCode.RESOURCE_EXHAUSTED
    )
    # Stub gets called for the max number of attempts.
    self.assertEqual(stub.call.call_count, 15)
    mock_sleep.assert_has_calls([mock.call(mock.ANY)])

  def test_retry_on_grpc_resource_exhausted_does_not_retry_on_other_grpc_error(
      self,
  ):
    stub = mock.MagicMock()
    stub.call.side_effect = _GrpcError(grpc.StatusCode.INVALID_ARGUMENT)
    with self.assertRaises(Exception) as context:
      _call_with_resource_exhausted_retry(stub)
    self.assertEqual(context.exception.code(), grpc.StatusCode.INVALID_ARGUMENT)
    self.assertEqual(stub.call.call_count, 1)

  def test_retry_on_grpc_resource_exhausted_does_not_retry_on_non_grpc_error(
      self,
  ):
    stub = mock.MagicMock()
    stub.call.side_effect = Exception('non-grpc error')
    with self.assertRaises(Exception) as context:
      _call_with_resource_exhausted_retry(stub)
    self.assertEqual(str(context.exception), 'non-grpc error')
    self.assertEqual(stub.call.call_count, 1)

  @mock.patch.object(time, 'sleep')
  def test_retry_on_grpc_transient_errors_retries_on_certain_errors(self, _):
    stub = mock.MagicMock()
    stub.call.side_effect = [
        _GrpcError(grpc.StatusCode.UNAVAILABLE),
        _GrpcError(grpc.StatusCode.UNIMPLEMENTED),
        _GrpcError(grpc.StatusCode.RESOURCE_EXHAUSTED),
        'some result',
    ]
    result = _call_with_transient_errors_retry(stub)
    stub.call.assert_called_with()
    self.assertEqual(result, 'some result')

  @mock.patch.object(time, 'sleep')
  def test_retry_on_grpc_transient_errors_fails_after_max_retries(
      self, mock_sleep
  ):
    stub = mock.MagicMock()
    # We receive an UNAVAILABLE error when we directly try to contact a grpc
    # server that is not (yet) running.
    stub.call.side_effect = _GrpcError(grpc.StatusCode.UNAVAILABLE)
    with self.assertRaises(_GrpcError) as context:
      _call_with_transient_errors_retry(stub)
    self.assertEqual(context.exception.code(), grpc.StatusCode.UNAVAILABLE)
    # Stub gets called for the max number of attempts.
    self.assertEqual(stub.call.call_count, 15)
    mock_sleep.assert_has_calls([mock.call(mock.ANY)])

  def test_retry_on_grpc_transient_errors_does_not_retry_on_other_grpc_error(
      self,
  ):
    stub = mock.MagicMock()
    stub.call.side_effect = _GrpcError(grpc.StatusCode.INVALID_ARGUMENT)
    with self.assertRaises(Exception) as context:
      _call_with_transient_errors_retry(stub)
    self.assertEqual(context.exception.code(), grpc.StatusCode.INVALID_ARGUMENT)
    self.assertEqual(stub.call.call_count, 1)

  def test_retry_on_grpc_transient_errors_does_not_retry_on_non_grpc_error(
      self,
  ):
    stub = mock.MagicMock()
    stub.call.side_effect = Exception('non-grpc error')
    with self.assertRaises(Exception) as context:
      _call_with_transient_errors_retry(stub)
    self.assertEqual(str(context.exception), 'non-grpc error')
    self.assertEqual(stub.call.call_count, 1)


if __name__ == '__main__':
  absltest.main()
