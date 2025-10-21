# Copyright 2023 Intrinsic Innovation LLC

"""Tests for the Asset dependencies utility functions."""

import dataclasses
from typing import cast

from absl.testing import absltest
from absl.testing import parameterized
import grpc
from grpc.framework.foundation import logging_pool
from intrinsic.assets.dependencies import utils
from intrinsic.assets.dependencies.testing import test_service_pb2
from intrinsic.assets.dependencies.testing import test_service_pb2_grpc
from intrinsic.assets.proto import id_pb2
from intrinsic.assets.proto.v1 import resolved_dependency_pb2


@dataclasses.dataclass(frozen=True)
class ConnectTestCase:
  desc: str
  dep: resolved_dependency_pb2.ResolvedDependency
  iface: str
  want_metadata: dict[str, list[str]] | None = None
  want_error: type[BaseException] | None = None


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
        desc="missing interface",
        dep=resolved_dependency_pb2.ResolvedDependency(),
        iface="grpc://intrinsic_proto.assets.dependencies.testing.TestService",
        want_error=utils.MissingInterfaceError,
    ),
    ConnectTestCase(
        desc="not a gRPC connection",
        dep=resolved_dependency_pb2.ResolvedDependency(
            interfaces={
                "data://foo": resolved_dependency_pb2.ResolvedDependency.Interface(
                    data=resolved_dependency_pb2.ResolvedDependency.Interface.Data(
                        id=id_pb2.Id(package="foo", name="bar"),
                    ),
                ),
            },
        ),
        iface="data://foo",
        want_error=utils.NotGRPCError,
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
      with self.assertRaises(tc.want_error):
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


if __name__ == "__main__":
  absltest.main()
