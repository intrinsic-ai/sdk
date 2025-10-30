# Copyright 2023 Intrinsic Innovation LLC

"""Tests for the Asset dependencies utility functions."""

import dataclasses
from typing import cast

from absl.testing import absltest
from absl.testing import parameterized
from google.protobuf import any_pb2
from google.protobuf import descriptor_pb2
from google.protobuf import empty_pb2
import grpc
from grpc.framework.foundation import logging_pool
from intrinsic.assets.data import fake_data_assets
from intrinsic.assets.data.proto.v1 import data_asset_pb2
from intrinsic.assets.dependencies import utils
from intrinsic.assets.dependencies.testing import test_service_pb2
from intrinsic.assets.dependencies.testing import test_service_pb2_grpc
from intrinsic.assets.proto import asset_type_pb2
from intrinsic.assets.proto import id_pb2
from intrinsic.assets.proto import metadata_pb2
from intrinsic.assets.proto.v1 import resolved_dependency_pb2


def _make_empty_data_asset(
    name: str = "data_asset",
) -> data_asset_pb2.DataAsset:
  payload = empty_pb2.Empty()
  payload_any = any_pb2.Any()
  payload_any.Pack(payload)

  return data_asset_pb2.DataAsset(
      data=payload_any,
      file_descriptor_set=descriptor_pb2.FileDescriptorSet(),
      metadata=metadata_pb2.Metadata(
          asset_type=asset_type_pb2.AssetType.ASSET_TYPE_DATA,
          id_version=id_pb2.IdVersion(
              id=id_pb2.Id(package="ai.intrinsic", name=name),
              version="0.0.1",
          ),
      ),
  )


@dataclasses.dataclass(frozen=True)
class ConnectTestCase:
  desc: str
  dep: resolved_dependency_pb2.ResolvedDependency
  iface: str
  want_metadata: dict[str, list[str]] | None = None
  want_error: type[BaseException] | None = None
  want_error_message_regex: str | None = None


@dataclasses.dataclass(frozen=True)
class GetDataPayloadTestCase:
  desc: str
  dep: resolved_dependency_pb2.ResolvedDependency
  iface: str
  want_payload: any_pb2.Any | None = None
  want_error: type[BaseException] | None = None
  want_error_message_regex: str | None = None


_CONNECT_TEST_CASES = [
    ConnectTestCase(
        desc="success",
        dep=resolved_dependency_pb2.ResolvedDependency(
            interfaces={
                "grpc://intrinsic_proto.assets.dependencies.testing.TestService": resolved_dependency_pb2.ResolvedDependency.Interface(
                    grpc_connection=resolved_dependency_pb2.ResolvedDependency.Interface.GrpcConnection(
                        address="localhost:12345",
                        metadata=[
                            resolved_dependency_pb2.ResolvedDependency.Interface.GrpcConnection.Metadata(
                                key="test_key",
                                value="test_value1",
                            ),
                            resolved_dependency_pb2.ResolvedDependency.Interface.GrpcConnection.Metadata(
                                key="test_key",
                                value="test_value2",
                            ),
                        ],
                    ),
                ),
            },
        ),
        iface="grpc://intrinsic_proto.assets.dependencies.testing.TestService",
        want_metadata={
            "test_key": ["test_value1", "test_value2"],
        },
    ),
    ConnectTestCase(
        desc="no interfaces",
        dep=resolved_dependency_pb2.ResolvedDependency(),
        iface="grpc://intrinsic_proto.assets.dependencies.testing.TestService",
        want_error=utils.MissingInterfaceError,
    ),
    ConnectTestCase(
        desc="wrong interface type",
        dep=resolved_dependency_pb2.ResolvedDependency(
            interfaces={
                "data://google.protobuf.Empty": resolved_dependency_pb2.ResolvedDependency.Interface(
                    data=resolved_dependency_pb2.ResolvedDependency.Interface.Data(
                        id=id_pb2.Id(package="ai.intrinsic", name="data_asset"),
                    ),
                ),
            },
        ),
        iface="grpc://intrinsic_proto.assets.dependencies.testing.TestService",
        want_error=utils.MissingInterfaceError,
        want_error_message_regex="got interfaces: data://google.protobuf.Empty",
    ),
    ConnectTestCase(
        desc="not a gRPC connection",
        dep=resolved_dependency_pb2.ResolvedDependency(
            interfaces={
                "data://google.protobuf.Empty": resolved_dependency_pb2.ResolvedDependency.Interface(
                    data=resolved_dependency_pb2.ResolvedDependency.Interface.Data(
                        id=id_pb2.Id(package="ai.intrinsic", name="data_asset"),
                    ),
                ),
            },
        ),
        iface="data://google.protobuf.Empty",
        want_error=utils.NotGRPCError,
    ),
]

_GET_DATA_PAYLOAD_TEST_CASES = [
    GetDataPayloadTestCase(
        desc="success",
        dep=resolved_dependency_pb2.ResolvedDependency(
            interfaces={
                "data://google.protobuf.Empty": resolved_dependency_pb2.ResolvedDependency.Interface(
                    data=resolved_dependency_pb2.ResolvedDependency.Interface.Data(
                        id=id_pb2.Id(package="ai.intrinsic", name="data_asset"),
                    ),
                ),
            },
        ),
        iface="data://google.protobuf.Empty",
        want_payload=_make_empty_data_asset().data,
    ),
    GetDataPayloadTestCase(
        desc="no interfaces",
        dep=resolved_dependency_pb2.ResolvedDependency(),
        iface="data://google.protobuf.Empty",
        want_error=utils.MissingInterfaceError,
        want_error_message_regex="no interfaces provided",
    ),
    GetDataPayloadTestCase(
        desc="wrong interface type",
        dep=resolved_dependency_pb2.ResolvedDependency(
            interfaces={
                "grpc://intrinsic_proto.assets.dependencies.testing.TestService": resolved_dependency_pb2.ResolvedDependency.Interface(
                    grpc_connection=resolved_dependency_pb2.ResolvedDependency.Interface.GrpcConnection(
                        address="localhost:12345",
                    ),
                ),
            },
        ),
        iface="data://google.protobuf.Empty",
        want_error=utils.MissingInterfaceError,
        want_error_message_regex=(
            "got interfaces:"
            " grpc://intrinsic_proto.assets.dependencies.testing.TestService"
        ),
    ),
    GetDataPayloadTestCase(
        desc="not data",
        dep=resolved_dependency_pb2.ResolvedDependency(
            interfaces={
                "grpc://intrinsic_proto.assets.dependencies.testing.TestService": resolved_dependency_pb2.ResolvedDependency.Interface(
                    grpc_connection=resolved_dependency_pb2.ResolvedDependency.Interface.GrpcConnection(
                        address="localhost:12345",
                    ),
                ),
            },
        ),
        iface="grpc://intrinsic_proto.assets.dependencies.testing.TestService",
        want_error=utils.NotDataError,
    ),
]


class TestService(test_service_pb2_grpc.TestServiceServicer):

  def Test(
      self,
      request: test_service_pb2.TestRequest,
      context: grpc.ServicerContext,
  ) -> test_service_pb2.TestResponse:
    response = test_service_pb2.TestResponse(
        context_metadata={},
    )
    for k, vs in context.invocation_metadata():
      response.context_metadata[k].values.append(vs)
    return response


class UtilsTest(parameterized.TestCase):

  @parameterized.parameters(*_CONNECT_TEST_CASES)
  def test_connect(self, tc: ConnectTestCase):
    if tc.want_error is not None:
      with self.assertRaisesRegex(tc.want_error, tc.want_error_message_regex):
        utils.connect(tc.dep, tc.iface)
    else:
      server = grpc.server(
          logging_pool.pool(max_workers=1),
          options=[("grpc.max_receive_message_length", -1)],
      )
      test_service_pb2_grpc.add_TestServiceServicer_to_server(
          TestService(), server
      )
      server.add_insecure_port(
          tc.dep.interfaces[tc.iface].grpc_connection.address
      )
      server.start()

      channel = utils.connect(tc.dep, tc.iface)
      stub = test_service_pb2_grpc.TestServiceStub(channel)
      response = stub.Test(test_service_pb2.TestRequest())
      channel.close()
      server.stop(None)
      for k, vs in cast(dict[str, list[str]], tc.want_metadata).items():
        self.assertIn(k, response.context_metadata)
        self.assertCountEqual(vs, response.context_metadata[k].values)

  @parameterized.parameters(*_GET_DATA_PAYLOAD_TEST_CASES)
  def test_get_data_payload(self, tc: GetDataPayloadTestCase):
    stub, cleanup = fake_data_assets.FakeDataAssetsService.start_server(
        [_make_empty_data_asset()]
    )

    try:
      if tc.want_error is not None:
        with self.assertRaisesRegex(tc.want_error, tc.want_error_message_regex):
          utils.get_data_payload(
              dep=tc.dep, iface=tc.iface, data_assets_client=stub
          )
      else:
        got_payload = utils.get_data_payload(
            dep=tc.dep, iface=tc.iface, data_assets_client=stub
        )
        self.assertEqual(got_payload, tc.want_payload, msg=tc.desc)
    finally:
      cleanup()


if __name__ == "__main__":
  absltest.main()
