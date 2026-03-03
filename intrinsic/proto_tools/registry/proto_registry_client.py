# Copyright 2023 Intrinsic Innovation LLC

"""Access client for Proto Registry."""

from __future__ import annotations

import grpc

from intrinsic.proto_tools.proto import proto_registry_pb2
from intrinsic.proto_tools.proto import proto_registry_pb2_grpc
INTRINSIC_TYPE_URL_PREFIX = 'type.intrinsic.ai'


class ProtoRegistryClient:
  """Client for the proto registry gRPC service."""

  _stub: proto_registry_pb2_grpc.ProtoRegistryStub

  def __init__(self, stub: proto_registry_pb2_grpc.ProtoRegistryStub):
    """Constructs a new ProtoRegistryClient object.

    Args:
      stub: The gRPC stub to be used for communication with the service.
    """
    self._stub = stub

  @classmethod
  def connect(cls, grpc_channel: grpc.Channel) -> ProtoRegistryClient:
    """Connects to a proto registry for an existing channel.

    Args:
      grpc_channel: Channel to the gRPC service.

    Returns:
      A newly created instance of the wrapper class.

    Raises:
      grpc.RpcError: When gRPC call to service fails.
    """
    stub = proto_registry_pb2_grpc.ProtoRegistryStub(grpc_channel)
    return cls(stub)

  def get_descriptor_set_by_typeurl(
      self,
      type_url: str,
  ) -> descriptor_pb2.FileDescriptorSet:
    request = proto_registry_pb2.GetNamedFileDescriptorSetRequest(
        type_url=type_url,
    )

    response = self._stub.GetNamedFileDescriptorSet(request)
    return response.file_descriptor_set
