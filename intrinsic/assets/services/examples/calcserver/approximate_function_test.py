# Copyright 2023 Intrinsic Innovation LLC

"""Tests for the ApproximateFunction CustomCalculation."""

import dataclasses

from absl.testing import absltest
from absl.testing import parameterized
from google.protobuf import text_format
import grpc
from grpc.framework.foundation import logging_pool
from intrinsic.assets.data.proto.v1 import data_manifest_pb2
from intrinsic.assets.services.examples.calcserver import approximate_function
from intrinsic.assets.services.examples.calcserver import calc_server_pb2
from intrinsic.assets.services.examples.calcserver import calc_server_pb2_grpc
from intrinsic.assets.services.examples.calcserver import two_dimensional_function_data_pb2
from intrinsic.util.path_resolver import path_resolver
import numpy as np

_INVERTED_PYRAMID_DATA_MANIFEST_PATH = "intrinsic/assets/services/examples/calcserver/inverted_pyramid_function_data.textproto"


def _make_linear_data(
    n, x1, x2, x3
) -> two_dimensional_function_data_pb2.TwoDimensionalFunctionData:
  """Returns a TwoDimensionalFunctionData for a linear function."""
  x = 2 * np.arange(n)
  y = 2 * np.arange(n) + 10
  f = x1 + x2 * x + x3 * y
  return two_dimensional_function_data_pb2.TwoDimensionalFunctionData(
      x=x.astype(int).tolist(),
      y=y.astype(int).tolist(),
      f=f.astype(int).tolist(),
  )


def _make_quadratic_data(
    n, x1, x2, x3, x4
) -> two_dimensional_function_data_pb2.TwoDimensionalFunctionData:
  """Returns a TwoDimensionalFunctionData for a quadratic function."""
  x = 2 * np.arange(n)
  y = 2 * np.arange(n) + 10
  f = x1 + x2 * x + x3 * y + x4 * y**2
  return two_dimensional_function_data_pb2.TwoDimensionalFunctionData(
      x=x.astype(int).tolist(),
      y=y.astype(int).tolist(),
      f=f.astype(int).tolist(),
  )


def _load_inverted_pyramid_data() -> (
    two_dimensional_function_data_pb2.TwoDimensionalFunctionData
):
  """Loads the inverted pyramid function data."""
  manifest_path = path_resolver.resolve_runfiles_path(
      _INVERTED_PYRAMID_DATA_MANIFEST_PATH
  )
  with open(manifest_path, "r") as f:
    manifest = text_format.Parse(f.read(), data_manifest_pb2.DataManifest())

  payload = two_dimensional_function_data_pb2.TwoDimensionalFunctionData()
  if not manifest.data.Unpack(payload):
    raise RuntimeError("Failed to unpack data")

  return payload


@dataclasses.dataclass(frozen=True)
class ApproximateFunctionTestCase:
  desc: str
  data: two_dimensional_function_data_pb2.TwoDimensionalFunctionData | None
  degree: int
  x: int
  y: int
  want: int | None = None
  tolerance: int | None = None
  want_code: grpc.StatusCode = grpc.StatusCode.OK


_APPROXIMATE_FUNCTION_TEST_CASES = [
    ApproximateFunctionTestCase(
        desc="no data",
        data=None,
        degree=3,
        x=0,
        y=0,
        want_code=grpc.StatusCode.FAILED_PRECONDITION,
    ),
    ApproximateFunctionTestCase(
        desc="empty data",
        data=two_dimensional_function_data_pb2.TwoDimensionalFunctionData(),
        degree=3,
        x=0,
        y=0,
        want=0,
    ),
    ApproximateFunctionTestCase(
        desc="single data point",
        data=two_dimensional_function_data_pb2.TwoDimensionalFunctionData(
            x=[1], y=[2], f=[3]
        ),
        degree=1,
        x=1,
        y=2,
        want=3,
    ),
    ApproximateFunctionTestCase(
        desc="linear function",
        data=_make_linear_data(n=5, x1=0, x2=1, x3=1),  # x + y
        degree=1,
        x=3,
        y=13,
        want=16,
        tolerance=1,
    ),
    ApproximateFunctionTestCase(
        desc="polynomial function",
        data=_make_quadratic_data(n=10, x1=3, x2=2, x3=0, x4=1),  # 3 + 2x + y^2
        degree=2,
        x=3,
        y=13,
        want=178,
        tolerance=1,
    ),
    ApproximateFunctionTestCase(
        desc="inverted pyramid function",
        data=_load_inverted_pyramid_data(),
        degree=2,
        x=8,
        y=-8,
        want=8,
        tolerance=1,
    ),
]


class ApproximateFunctionTest(parameterized.TestCase):

  @parameterized.parameters(*_APPROXIMATE_FUNCTION_TEST_CASES)
  def test_approximate_function(self, tc: ApproximateFunctionTestCase):
    service = approximate_function.ApproximateFunction(tc.data, tc.degree)
    server = grpc.server(logging_pool.pool(max_workers=1))
    calc_server_pb2_grpc.add_CustomCalculationServicer_to_server(
        service, server
    )
    port = server.add_secure_port("[::]:0", grpc.local_server_credentials())
    server.start()
    channel = grpc.secure_channel(
        f"localhost:{port}", grpc.local_channel_credentials()
    )
    stub = calc_server_pb2_grpc.CustomCalculationStub(channel)

    try:
      request = calc_server_pb2.CustomCalculateRequest(x=tc.x, y=tc.y)
      if tc.want_code != grpc.StatusCode.OK:
        with self.assertRaises(grpc.RpcError) as e:
          stub.Calculate(request)
        self.assertEqual(_get_status_code(e.exception), tc.want_code)
      else:
        response = stub.Calculate(request)
        min_valid = tc.want - tc.tolerance if tc.tolerance else tc.want
        max_valid = tc.want + tc.tolerance if tc.tolerance else tc.want
        self.assertBetween(response.result, min_valid, max_valid)
    finally:
      server.stop(None)
      channel.close()


def _get_status_code(e: grpc.RpcError) -> grpc.StatusCode:
  return e.code()  # pytype: disable=attribute-error


if __name__ == "__main__":
  absltest.main()
