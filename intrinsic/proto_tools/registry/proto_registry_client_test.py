# Copyright 2023 Intrinsic Innovation LLC

import datetime
from unittest import mock

from absl.testing import absltest

# isort: off
# isort: on
from google.protobuf import descriptor_pb2

# isort: off
# isort: on

from intrinsic.proto_tools.proto import proto_registry_pb2
from intrinsic.proto_tools.proto import proto_registry_pb2_grpc
from intrinsic.proto_tools.registry import proto_registry_client
from intrinsic.util.path_resolver import path_resolver

DESCRIPTOR_TEST_MESSAGE_FILENAME = "intrinsic/proto_tools/registry/test_data/descriptor_test_message_proto_descriptor_set_transitive_set_sci.proto.bin"


class ProtoRegistryTest(absltest.TestCase):

  @mock.patch.object(proto_registry_pb2_grpc, "ProtoRegistryStub")
  def test_get_descriptor_set_by_typeurl(
      self, stub: proto_registry_pb2_grpc.ProtoRegistryStub
  ):
    stub.GetNamedFileDescriptorSet.return_value = (
        proto_registry_pb2.NamedFileDescriptorSet(
            file_descriptor_set=descriptor_pb2.FileDescriptorSet(
                file=[descriptor_pb2.FileDescriptorProto(name="foo")]
            )
        )
    )
    proto_registry = proto_registry_client.ProtoRegistryClient(stub)

    fds = proto_registry.get_descriptor_set_by_typeurl("foo/url")
    self.assertEqual(fds.file[0].name, "foo")

    stub.GetNamedFileDescriptorSet.assert_called_once_with(
        proto_registry_pb2.GetNamedFileDescriptorSetRequest(type_url="foo/url")
    )

if __name__ == "__main__":
  absltest.main()
