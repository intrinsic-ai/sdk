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

"""Tests for the connect utility functions."""

from absl.testing import absltest
from absl.testing import parameterized
import grpc
from grpc.framework.foundation import logging_pool

from intrinsic.assets.instances.connect import connect
from intrinsic.assets.instances.connect.testing import test_service_pb2
from intrinsic.assets.instances.connect.testing import test_service_pb2_grpc
from intrinsic.assets.proto.v1 import grpc_connection_pb2


class TestService(test_service_pb2_grpc.TestServiceServicer):

  def Test(
      self,
      request: test_service_pb2.TestRequest,
      context: grpc.ServicerContext,
  ) -> test_service_pb2.TestResponse:
    response = test_service_pb2.TestResponse(context_metadata={})
    for k, vs in context.invocation_metadata():
      response.context_metadata[k].values.append(vs)
    return response


class ConnectTest(parameterized.TestCase):

  def test_connect_with_metadata(self):
    server = grpc.server(
        logging_pool.pool(max_workers=1),
        options=[("grpc.max_receive_message_length", -1)],
    )
    test_service_pb2_grpc.add_TestServiceServicer_to_server(
        TestService(), server
    )
    port = server.add_insecure_port("localhost:0")
    server.start()
    address = f"localhost:{port}"

    connection = grpc_connection_pb2.GrpcConnection(
        address=address,
        metadata=[
            grpc_connection_pb2.GrpcConnection.Metadata(
                key="test_key", value="test_value1"
            ),
            grpc_connection_pb2.GrpcConnection.Metadata(
                key="test_key", value="test_value2"
            ),
        ],
    )

    channel = connect.connect(connection)
    stub = test_service_pb2_grpc.TestServiceStub(channel)
    response = stub.Test(test_service_pb2.TestRequest())
    channel.close()
    server.stop(None)

    self.assertIn("test_key", response.context_metadata)
    self.assertCountEqual(
        ["test_value1", "test_value2"],
        response.context_metadata["test_key"].values,
    )


if __name__ == "__main__":
  absltest.main()
