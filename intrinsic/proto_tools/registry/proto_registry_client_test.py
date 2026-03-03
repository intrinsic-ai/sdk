# Copyright 2023 Intrinsic Innovation LLC

import datetime
from unittest import mock

from absl.testing import absltest
from google.protobuf import any_pb2
from google.protobuf import descriptor_pb2
from google.protobuf import message_factory

from intrinsic.proto_tools.proto import proto_registry_pb2
from intrinsic.proto_tools.proto import proto_registry_pb2_grpc
from intrinsic.proto_tools.registry import proto_registry_client
from intrinsic.util.path_resolver import path_resolver
from intrinsic.util.proto import descriptors as descriptor_utils

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

  @mock.patch.object(proto_registry_pb2_grpc, "ProtoRegistryStub")
  def test_message_to_json(
      self, stub: proto_registry_pb2_grpc.ProtoRegistryStub
  ):
    with open(
        path_resolver.resolve_runfiles_path(DESCRIPTOR_TEST_MESSAGE_FILENAME),
        "rb",
    ) as fp:
      fds = descriptor_pb2.FileDescriptorSet.FromString(fp.read())

    stub.GetNamedFileDescriptorSet.return_value = (
        proto_registry_pb2.NamedFileDescriptorSet(file_descriptor_set=fds)
    )
    proto_registry = proto_registry_client.ProtoRegistryClient(stub)

    known_time = datetime.datetime(
        2024, 3, 26, 11, 51, 13, tzinfo=datetime.timezone.utc
    )
    pool = descriptor_utils.create_descriptor_pool(fds)
    message_descriptor = pool.FindMessageTypeByName(
        "intrinsic.testing.proto.DescriptorTestMessage"
    )
    message_class = message_factory.GetMessageClass(message_descriptor)
    msg = message_class()
    msg.timestamp.FromDatetime(known_time)

    any_message = any_pb2.Any()
    any_message.Pack(msg)
    any_message.type_url = f"{proto_registry_client.INTRINSIC_TYPE_URL_PREFIX}/something/test/1.2.3/{msg.DESCRIPTOR.full_name}"
    json = proto_registry.message_to_json(any_message)
    self.assertEqual(
        json,
        (
            '{\n  "@type":'
            ' "type.intrinsic.ai/something/test/1.2.3/intrinsic.testing.proto.DescriptorTestMessage",\n'
            '  "timestamp": "2024-03-26T11:51:13Z"\n}'
        ),
    )

    stub.GetNamedFileDescriptorSet.assert_called_once_with(
        proto_registry_pb2.GetNamedFileDescriptorSetRequest(
            type_url=any_message.type_url
        )
    )


if __name__ == "__main__":
  absltest.main()
