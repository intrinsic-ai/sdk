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

"""Conversions to/from numpy ndarray from/to matrix proto types."""

import numpy as np

from intrinsic.icon.proto import matrix_pb2


def from_ndarray(matrix: np.ndarray) -> matrix_pb2.Matrix6d:
  """Converts from a numpy ndarray to a Matrix6d proto.

  Args:
    matrix: the matrix as a ndarray.

  Returns:
    A Matrix6d proto.

  Raises:
    ValueError: If the input ndarray is not 6x6.
  """
  if matrix.shape[0] != matrix.shape[1] != 6:
    raise ValueError('Matrix is not 6x6, received size [%s,%s]' % matrix.shape)

  proto_matrix = matrix_pb2.Matrix6d()
  for val in matrix.flatten():
    proto_matrix.data.append(val)

  return proto_matrix


def to_ndarray(proto_matrix: matrix_pb2.Matrix6d) -> np.ndarray:
  """Converts from a Matrix6d proto to a numpy ndarray.

  Args:
    proto_matrix: the matrix as a Matrix6d.

  Returns:
    A numpy ndarray.

  Raises:
    ValueError: if the input matrix does not have 36 elements.
  """
  if len(proto_matrix.data) != 36:
    raise ValueError(
        'Matrix is not 6x6, received size %s' % len(proto_matrix.data)
    )

  return np.array([
      proto_matrix.data[0:6],
      proto_matrix.data[6:12],
      proto_matrix.data[12:18],
      proto_matrix.data[18:24],
      proto_matrix.data[24:30],
      proto_matrix.data[30:36],
  ])
