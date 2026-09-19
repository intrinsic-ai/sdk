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
