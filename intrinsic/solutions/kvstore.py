# Copyright 2023 Intrinsic Innovation LLC

"""API to interface with the Key-Value Store GRPC service."""

from typing import Iterable

from google.protobuf import any_pb2
from google.protobuf import message as proto_message
from google.protobuf import wrappers_pb2
import grpc

from intrinsic.platform.pubsub.kvstore_grpc import kvstore_pb2
from intrinsic.platform.pubsub.kvstore_grpc import kvstore_pb2_grpc
from intrinsic.solutions import deployments
from intrinsic.util.grpc import error_handling


class KVStore:
  """Wrapper for interacting with the Key-Value Store."""

  @classmethod
  def for_solution(cls, solution: deployments.Solution) -> "KVStore":
    """Connects to the KVStore gRPC service for a given solution."""
    return cls(
        stub=kvstore_pb2_grpc.KVStoreStub(solution.grpc_channel),
    )

  def __init__(self, stub: kvstore_pb2_grpc.KVStoreStub):
    """Initializes the KVStore client.

    Args:
      stub: The GRPC stub for the KVStore service.
    """
    self._stub = stub

  @error_handling.retry_on_grpc_transient_errors
  def get(self, key: str) -> any_pb2.Any:
    """Gets a value from the store.

    Args:
      key: The key to look up.

    Returns:
      The value associated with the key, unpacked as a protobuf message, or None
      if the key does not exist.

    Raises:
      grpc.RpcError: If the GRPC call fails.
    """
    try:
      response = self._stub.Get(kvstore_pb2.GetRequest(key=key))
      return response.value
    except grpc.RpcError as e:
      if e.code() == grpc.StatusCode.NOT_FOUND:
        return None
      raise

  @error_handling.retry_on_grpc_transient_errors
  def set(self, key: str, value: any_pb2.Any) -> None:
    """Sets a value in the store.

    Args:
      key: The key to set.
      value: The protobuf message to set as the value.
    """
    self._stub.Set(kvstore_pb2.SetRequest(key=key, value=value))

  @error_handling.retry_on_grpc_transient_errors
  def delete(self, key: str) -> None:
    """Deletes a value from the store.

    Args:
      key: The key to delete.
    """
    self._stub.Delete(kvstore_pb2.DeleteRequest(key=key))

  @error_handling.retry_on_grpc_transient_errors
  def keys(self) -> list[str]:
    """Lists all keys in the store.

    Returns:
      A list of all keys in the store.
    """
    response = self._stub.List(kvstore_pb2.ListRequest())
    return list(response.keys)
