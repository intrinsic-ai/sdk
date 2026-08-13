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

from unittest import mock

from absl.testing import absltest
from google.protobuf import empty_pb2

from intrinsic.simulation.service.proto.v1 import simulation_service_pb2
from intrinsic.solutions import simulation as simulation_mod


class SimulationTest(absltest.TestCase):
  def setUp(self):
    super().setUp()
    self.simulation_service_stub = mock.MagicMock()
    self.object_world_service_stub = mock.MagicMock()
    self.simulation = simulation_mod.Simulation(
        self.simulation_service_stub, self.object_world_service_stub
    )

  def test_reset(self):
    self.simulation_service_stub.ResetSimulation.return_value = (
        empty_pb2.Empty()
    )

    self.simulation.reset()

    self.simulation_service_stub.ResetSimulation.assert_called_once_with(
        simulation_service_pb2.ResetSimulationRequest()
    )


if __name__ == '__main__':
  absltest.main()
