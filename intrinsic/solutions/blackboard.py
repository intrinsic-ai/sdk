# Copyright 2023 Intrinsic Innovation LLC

"""Blackboard access within the solution building library."""

import enum
import typing

from google.protobuf import any_pb2
from google.protobuf import message
from google.protobuf import wrappers_pb2

from intrinsic.executive.proto import blackboard_service_pb2
from intrinsic.executive.proto import blackboard_service_pb2_grpc
from intrinsic.solutions import blackboard_value
from intrinsic.solutions import ipython
from intrinsic.solutions.internal import skill_utils
from intrinsic.util.grpc import error_handling


class ScopedBlackboardKey(typing.NamedTuple):
  """An entry on the blackboard.

  Attributes:
    key: The key of the entry.
    scope: The scope of the entry.
    type_url: The type URL of the entry.
  """

  key: str
  scope: str
  type_url: str


_WRAPPER_CLASSES = [
    wrappers_pb2.DoubleValue,
    wrappers_pb2.FloatValue,
    wrappers_pb2.Int64Value,
    wrappers_pb2.UInt64Value,
    wrappers_pb2.Int32Value,
    wrappers_pb2.UInt32Value,
    wrappers_pb2.BoolValue,
    wrappers_pb2.StringValue,
    wrappers_pb2.BytesValue,
]
_WRAPPER_TYPES = {cls.DESCRIPTOR.full_name: cls for cls in _WRAPPER_CLASSES}

_PYTHON_TYPE_TO_WRAPPERS = {
    int: {
        wrappers_pb2.Int64Value,
        wrappers_pb2.UInt64Value,
        wrappers_pb2.Int32Value,
        wrappers_pb2.UInt32Value,
    },
    float: {
        wrappers_pb2.DoubleValue,
        wrappers_pb2.FloatValue,
    },
    bool: {wrappers_pb2.BoolValue},
    str: {wrappers_pb2.StringValue},
    bytes: {wrappers_pb2.BytesValue},
}


class Blackboard:
  """Convenience wrapper for blackboard access."""

  _stub: blackboard_service_pb2_grpc.ExecutiveBlackboardStub
  _operation_name: str

  def __init__(
      self,
      stub: blackboard_service_pb2_grpc.ExecutiveBlackboardStub,
      operation_name: str,
  ):
    """Initializes the blackboard.

    Args:
      stub: The gRPC stub to be used for blackboard related calls.
      operation_name: The name of the operation this blackboard belongs to.
    """
    self._stub = stub
    self._operation_name = operation_name

  def _resolve_key_and_scope(
      self, key: str | blackboard_value.BlackboardValue, scope: str | None
  ) -> tuple[str, str | None]:
    """Resolves the key and scope from the input.

    If key is a BlackboardValue, the key and scope are extracted from it.
    """
    if isinstance(key, blackboard_value.BlackboardValue):
      if scope is not None:
        raise ValueError(
            f"Cannot provide explicit scope '{scope}' when using a"
            " BlackboardValue."
        )
      if not key.is_toplevel_value:
        raise ValueError(
            f"BlackboardValue with path {key.value_access_path()} is not a"
            " toplevel value."
        )
      return key.value_access_path(), key.scope()
    if isinstance(key, str):
      return key, scope
    raise TypeError(
        f"Expected str or BlackboardValue for 'key', got {type(key)}"
    )

  @error_handling.retry_on_grpc_unavailable
  @error_handling.log_extended_status(
      ipython.display_extended_status_proto_if_ipython
  )
  def delete_value(
      self,
      key: str | blackboard_value.BlackboardValue,
      scope: str | None = None,
  ) -> None:
    """Deletes a specific value from the blackboard.

    Args:
      key: The key or BlackboardValue to delete.
      scope: Optional scope. If not specified, the value is deleted from the
        main process scope.
    """
    key, scope = self._resolve_key_and_scope(key, scope)
    request = blackboard_service_pb2.DeleteBlackboardValueRequest(
        key=key,
        scope=scope or "",
        operation_name=self._operation_name,
    )
    self._stub.DeleteBlackboardValue(request)

  @error_handling.retry_on_grpc_unavailable
  @error_handling.log_extended_status(
      ipython.display_extended_status_proto_if_ipython
  )
  def list_keys(self, scope: str | None = None) -> list[ScopedBlackboardKey]:
    """Lists keys on the blackboard.

    Args:
      scope: Optional scope to filter by.

    Returns:
      A list of ScopedBlackboardKey objects containing key and scope.
    """
    request = blackboard_service_pb2.ListBlackboardValuesRequest(
        operation_name=self._operation_name,
        scope=scope,
        view=blackboard_service_pb2.ListBlackboardValuesRequest.ANY_TYPEURL_ONLY,
    )
    response = self._stub.ListBlackboardValues(request)
    return [
        ScopedBlackboardKey(key=v.key, scope=v.scope, type_url=v.value.type_url)
        for v in response.values
    ]

  @error_handling.retry_on_grpc_unavailable
  @error_handling.log_extended_status(
      ipython.display_extended_status_proto_if_ipython
  )
  def get_value_any(
      self,
      key: str | blackboard_value.BlackboardValue,
      scope: str | None = None,
  ) -> any_pb2.Any:
    """Gets a value from the blackboard as an Any proto.

    Args:
      key: The key or BlackboardValue to retrieve.
      scope: Optional scope.

    Returns:
      The value from the blackboard as an Any proto.
    """
    key, scope = self._resolve_key_and_scope(key, scope)
    request = blackboard_service_pb2.GetBlackboardValueRequest(
        key=key,
        scope=scope,
        operation_name=self._operation_name,
    )
    response = self._stub.GetBlackboardValue(request)
    return response.value

  def get_value(
      self,
      key: str | blackboard_value.BlackboardValue,
      scope: str | None = None,
  ) -> int | float | bool | str | bytes | any_pb2.Any:
    """Gets a value from the blackboard.

    Args:
      key: The key or BlackboardValue to retrieve.
      scope: Optional scope.

    Returns:
      The value from the blackboard. If it is a known wrapper type, the native
      Python value is returned. Otherwise, the Any proto is returned.
    """
    any_value = self.get_value_any(key, scope)

    type_name = any_value.type_url.rpartition("/")[-1]
    if type_name in _WRAPPER_TYPES:
      wrapper = _WRAPPER_TYPES[type_name]()
      any_value.Unpack(wrapper)
      return wrapper.value

    return any_value

  @error_handling.retry_on_grpc_unavailable
  @error_handling.log_extended_status(
      ipython.display_extended_status_proto_if_ipython
  )
  def update_value(
      self,
      key: str | blackboard_value.BlackboardValue,
      value: (
          int
          | float
          | bool
          | str
          | bytes
          | any_pb2.Any
          | message.Message
          | skill_utils.MessageWrapper
      ),
      scope: str | None = None,
  ) -> None:
    """Updates a value on the blackboard.

    Args:
      key: The key or BlackboardValue to update.
      value: The value to set. Can be an Any proto, a generic protobuf message,
        a MessageWrapper or a native Python type (int, float, bool, str, bytes).
      scope: Optional scope.

    Raises:
      TypeError: If the value type does not match the existing blackboard value.
    """
    key, scope = self._resolve_key_and_scope(key, scope)
    existing_any = self.get_value_any(key, scope)
    existing_type = existing_any.type_url.rpartition("/")[-1]

    if isinstance(value, any_pb2.Any):
      any_val = value
    elif isinstance(value, skill_utils.MessageWrapper):
      any_val = value.to_any()
    elif isinstance(value, message.Message):
      any_val = any_pb2.Any()
      any_val.Pack(value)
    elif (py_type := type(value)) in _PYTHON_TYPE_TO_WRAPPERS:
      wrapper_cls = _WRAPPER_TYPES.get(existing_type)
      if wrapper_cls in _PYTHON_TYPE_TO_WRAPPERS[py_type]:
        any_val = any_pb2.Any()
        any_val.Pack(wrapper_cls(value=value))
      else:
        article = "an" if py_type is int else "a"
        raise TypeError(
            f"Type mismatch for key '{key}': existing type {existing_type} is"
            f" not {article} {py_type.__name__} wrapper"
        )
    else:
      raise TypeError(
          "Expected Any, Message, MessageWrapper or native type for 'value',"
          f" got {type(value)}"
      )

    # Check for type mismatch
    if any_val.type_url != existing_any.type_url:
      raise TypeError(
          f"Type mismatch for key '{key}': existing type"
          f" {existing_any.type_url}, new type {any_val.type_url}"
      )

    request = blackboard_service_pb2.UpdateBlackboardValueRequest(
        value=blackboard_service_pb2.BlackboardValue(
            key=key,
            scope=scope or "",
            operation_name=self._operation_name,
            value=any_val,
        )
    )
    self._stub.UpdateBlackboardValue(request)
