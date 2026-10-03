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

"""Implements a CustomCalculation for approximating a function."""

from typing import Callable

from absl import logging
import grpc
import numpy as np

from intrinsic.assets.services.examples.calcserver import calc_server_pb2
from intrinsic.assets.services.examples.calcserver import calc_server_pb2_grpc
from intrinsic.assets.services.examples.calcserver import two_dimensional_function_data_pb2


class ApproximateFunction(calc_server_pb2_grpc.CustomCalculationServicer):
  """Approximates a function, given data from a Data Asset.

  The function is approximated using a polynomial regression model.

  Can be used as a custom calculation dependency for the Calculator Service.
  """

  _estimate: Callable[[int, int], int] | None

  def __init__(
      self,
      data: two_dimensional_function_data_pb2.TwoDimensionalFunctionData | None,
      degree: int,
  ):
    """Initializes the object."""
    self._estimate = None if data is None else _fit_model(data, degree)

  def Calculate(
      self,
      request: calc_server_pb2.CustomCalculateRequest,
      context: grpc.ServicerContext,
  ) -> calc_server_pb2.CalculatorResponse:
    """Approximates a function, given data from a Data Asset."""
    if self._estimate is None:
      context.abort(
          grpc.StatusCode.FAILED_PRECONDITION,
          "No data provided to approximate function.",
      )

    result = self._estimate(request.x, request.y)

    logging.info(
        "Calculate: x: %d, y: %d, result: %d",
        request.x,
        request.y,
        result,
    )

    return calc_server_pb2.CalculatorResponse(result=result)


def _fit_model(
    data: two_dimensional_function_data_pb2.TwoDimensionalFunctionData,
    degree: int,
) -> Callable[[int, int], int]:
  """Fits a model to the specified data."""
  x = np.reshape(data.x, (-1, 1))
  y = np.reshape(data.y, (-1, 1))
  f = np.reshape(data.f, (-1, 1))
  n = len(f)
  if len(data.x) != n or len(data.y) != n:
    raise ValueError(
        "All data must have the same length (got x: %d, y: %d, f: %d)"
        % (len(x), len(y), len(f))
    )

  if n == 0:
    return _zero

  return _fit_polynomial_regressor(x=x, y=y, f=f, degree=degree)


def _zero(x: int, y: int) -> int:
  del x, y  # Unused.
  return 0


def _fit_polynomial_regressor(
    x: np.ndarray, y: np.ndarray, f: np.ndarray, degree: int
) -> Callable[[int, int], int]:
  """Predicts the function value at the specified coordinates."""

  # Fit the polynomial to the data.
  design = _create_design_matrix(x=x, y=y, degree=degree)
  theta, _, _, _ = np.linalg.lstsq(design, f)

  def _estimate(x: int, y: int) -> int:
    design = _create_design_matrix(
        x=np.array([x]).reshape(-1, 1),
        y=np.array([y]).reshape(-1, 1),
        degree=degree,
    )
    f_pred = design @ theta
    return int(np.round(f_pred[0]))

  return _estimate


def _create_design_matrix(
    x: np.ndarray, y: np.ndarray, degree: int
) -> np.ndarray:
  """Creates the design matrix for the specified polynomial degree."""
  features = []
  for d in range(1, degree + 1):
    for y_degree in range(d + 1):
      x_poly = np.power(x, degree - y_degree)
      y_poly = np.power(y, y_degree)
      features.append(x_poly * y_poly)

  return np.hstack((np.ones_like(x), np.hstack(features)))
