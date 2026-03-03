# Copyright 2023 Intrinsic Innovation LLC

"""Access client for Proto Registry."""

from __future__ import annotations

# Thus MUST be from google.protobuf in order for monkey patching to work
from google.protobuf import descriptor_pb2
from google.protobuf import descriptor_pool
from google.protobuf import json_format
from google.protobuf import message as protobuf_message
from google.protobuf import message_factory
import grpc

from intrinsic.proto_tools.proto import proto_registry_pb2
from intrinsic.proto_tools.proto import proto_registry_pb2_grpc
from intrinsic.util.proto import descriptors as descriptor_utils

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

  def message_to_json(self, message: protobuf_message.Message) -> str:
    """Converts a message to JSON resolving Any protos from proto registry.

    This is similar to a normal json_format conversion, but Any protos that it
    encounters are resolved using the proto registry (which requires to have
    Intrinsic-style type URLs for them).

    Args:
      message: Message to convert to JSON.

    Returns:
      JSON string representing message.
    """
    # Monkey patching json format so we can use proto registry for Any encoding
    old_create = json_format._CreateMessageFromTypeUrl  # pylint: disable=protected-access

    def create_message_from_type_url_with_proto_registry(
        type_url: str, pool: descriptor_pool.DescriptorPool
    ) -> protobuf_message.Message:
      del pool  # unused
      fds = self.get_descriptor_set_by_typeurl(type_url)
      pool = descriptor_utils.create_descriptor_pool(fds)

      type_name = type_url.split('/')[-1]
      try:
        message_descriptor = pool.FindMessageTypeByName(type_name)
      except KeyError as e:
        raise TypeError(
            f'Can not find message descriptor by type_url: {type_url}'
        ) from e
      message_class = message_factory.GetMessageClass(message_descriptor)
      return message_class()

    setattr(
        json_format,
        '_CreateMessageFromTypeUrl',
        create_message_from_type_url_with_proto_registry,
    )
    json = json_format.MessageToJson(message)
    setattr(json_format, '_CreateMessageFromTypeUrl', old_create)
    return json
