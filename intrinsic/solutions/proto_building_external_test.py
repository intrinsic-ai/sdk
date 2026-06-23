# Copyright 2023 Intrinsic Innovation LLC

"""Tests of the proto_building module which also work externally.

These tests are simple smoke tests which use mocks and don't have dependencies
on internal fake implementations.
"""

from unittest import mock

from absl.testing import absltest
from google.protobuf import descriptor_pb2
from google.protobuf import text_format

from intrinsic.executive.proto import proto_builder_pb2
from intrinsic.solutions import proto_building
from intrinsic.solutions.testing import compare


class ProtoBuildingExternalTest(absltest.TestCase):

  def setUp(self):
    super().setUp()
    self.stub = mock.MagicMock()
    self.proto_builder = proto_building.ProtoBuilder(self.stub)

  def test_get_well_known_types(self):
    self.stub.GetWellKnownTypes.return_value = (
        proto_builder_pb2.GetWellKnownTypesResponse(type_names=['foo', 'bar'])
    )

    types = self.proto_builder.get_well_known_types()

    self.assertEqual(types, ['foo', 'bar'])
    self.stub.GetWellKnownTypes.assert_called_once_with(
        proto_builder_pb2.GetWellKnownTypesRequest()
    )

  def test_compile(self):
    expected_fds = descriptor_pb2.FileDescriptorSet(
        file=[
            descriptor_pb2.FileDescriptorProto(
                name='foo.proto', syntax='proto3'
            )
        ]
    )
    self.stub.Compile.return_value = proto_builder_pb2.ProtoCompileResponse(
        file_descriptor_set=expected_fds
    )

    fds = self.proto_builder.compile('foo.proto', 'syntax = "proto3";')

    self.assertEqual(fds, expected_fds)
    self.stub.Compile.assert_called_once_with(
        proto_builder_pb2.ProtoCompileRequest(
            proto_filename='foo.proto', proto_schema='syntax = "proto3";'
        )
    )

  def test_compose(self):
    expected_fds = descriptor_pb2.FileDescriptorSet(
        file=[
            descriptor_pb2.FileDescriptorProto(
                name='foo.proto',
                syntax='proto3',
                message_type=[descriptor_pb2.DescriptorProto(name='Foo')],
            )
        ]
    )
    self.stub.Compose.return_value = proto_builder_pb2.ProtoComposeResponse(
        file_descriptor_set=expected_fds
    )
    input_descriptors = [descriptor_pb2.DescriptorProto(name='Foo')]

    fds = self.proto_builder.compose(
        'foo.proto', 'my_package', input_descriptors
    )

    self.assertEqual(fds, expected_fds)
    self.stub.Compose.assert_called_once_with(
        proto_builder_pb2.ProtoComposeRequest(
            proto_filename='foo.proto',
            proto_package='my_package',
            input_descriptor=input_descriptors,
        )
    )

  def test_create_message(self):
    self.stub.GetWellKnownTypes.return_value = (
        proto_builder_pb2.GetWellKnownTypesResponse()
    )
    expected_fds = descriptor_pb2.FileDescriptorSet(
        file=[
            descriptor_pb2.FileDescriptorProto(
                name='my_pkg_MyMessage.proto',
                package='my_pkg',
                syntax='proto3',
                message_type=[descriptor_pb2.DescriptorProto(name='MyMessage')],
            )
        ]
    )
    self.stub.Compose.return_value = proto_builder_pb2.ProtoComposeResponse(
        file_descriptor_set=expected_fds
    )

    msg = self.proto_builder.create_message('my_pkg', 'MyMessage', {})

    self.assertIsNotNone(msg)
    self.assertEqual(msg.DESCRIPTOR.full_name, 'my_pkg.MyMessage')

    expected_input_descriptor = descriptor_pb2.DescriptorProto(name='MyMessage')
    self.stub.Compose.assert_called_once_with(
        proto_builder_pb2.ProtoComposeRequest(
            proto_filename='my_pkg_MyMessage.proto',
            proto_package='my_pkg',
            input_descriptor=[expected_input_descriptor],
        )
    )

  def test_create_signature(self):
    self.stub.GetWellKnownTypes.return_value = (
        proto_builder_pb2.GetWellKnownTypesResponse()
    )
    expected_fds = text_format.Parse(
        """
        file {
          name: "gen/sbl/signature.proto"
          package: "gen.sbl"
          message_type {
            name: "Params"
            field {
              name: "x"
              number: 1
              type: TYPE_INT64
            }
          }
          message_type {
            name: "ReturnValue"
            field {
              name: "y"
              number: 1
              type: TYPE_STRING
            }
          }
          syntax: "proto3"
        }""",
        descriptor_pb2.FileDescriptorSet(),
    )
    self.stub.Compile.return_value = proto_builder_pb2.ProtoCompileResponse(
        file_descriptor_set=expected_fds
    )

    sig = self.proto_builder.create_signature(
        parameters=proto_building.MessageSpec(
            fields=[proto_building.FieldSpec(name='x', type='int64', number=1)]
        ),
        return_value=proto_building.MessageSpec(
            fields=[
                proto_building.FieldSpec(name='y', type='string', number=42)
            ]
        ),
    )

    self.assertEqual(
        sig.parameter_message_full_name,
        f'gen.sbl.{proto_building._DEFAULT_PARAM_MSG_NAME}',
    )
    self.assertEqual(
        sig.return_value_message_full_name,
        f'gen.sbl.{proto_building._DEFAULT_RETURN_MSG_NAME}',
    )
    compare.assertProto2Equal(self, sig.file_descriptor_set, expected_fds)

    expected_schema = (
        'syntax = "proto3";\n'
        'package gen.sbl;\n'
        '\n'
        '\n'
        f'message {proto_building._DEFAULT_PARAM_MSG_NAME} {{\n'
        '  int64 x = 1;\n'
        '}\n'
        '\n'
        f'message {proto_building._DEFAULT_RETURN_MSG_NAME} {{\n'
        '  string y = 42;\n'
        '}\n'
    )
    self.stub.Compile.assert_called_once_with(
        proto_builder_pb2.ProtoCompileRequest(
            proto_filename=proto_building._DEFAULT_SIGNATURE_PROTO_FILENAME,
            proto_schema=expected_schema,
        )
    )
    self.stub.GetWellKnownTypes.assert_called_once()


if __name__ == '__main__':
  absltest.main()
