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

from absl.testing import absltest

from intrinsic.assets.proto import status_spec_pb2
from intrinsic.solutions.testing import compare

# isort: off
# isort: on
from intrinsic.util.status import extended_status_pb2
from intrinsic.util.status import status_specs

_COMPONENT = "ai.intrinsic.test"


class StatusSpecsTest(absltest.TestCase):

  def test_from_map(self):
    status_specs.init_once(
        _COMPONENT,
        status_specs=[
            status_spec_pb2.StatusSpec(
                code=10001,
                title="Error 1",
                recovery_instructions="Test instructions 1",
            ),
            status_spec_pb2.StatusSpec(
                code=10002,
                title="Error 2",
                recovery_instructions="Test instructions 2",
            ),
        ],
    )

    with self.subTest(name="Create status for declared code"):
      compare.assertProto2Equal(
          self,
          status_specs.create(10001, "user message").proto,
          extended_status_pb2.ExtendedStatus(
              status_code=extended_status_pb2.StatusCode(
                  component=_COMPONENT, code=10001
              ),
              title="Error 1",
              user_report=extended_status_pb2.ExtendedStatus.UserReport(
                  message="user message", instructions="Test instructions 1"
              ),
          ),
          ignored_fields=["timestamp"],
      )

      compare.assertProto2Equal(
          self,
          status_specs.create(10002, "user message").proto,
          extended_status_pb2.ExtendedStatus(
              status_code=extended_status_pb2.StatusCode(
                  component=_COMPONENT, code=10002
              ),
              title="Error 2",
              user_report=extended_status_pb2.ExtendedStatus.UserReport(
                  message="user message", instructions="Test instructions 2"
              ),
          ),
          ignored_fields=["timestamp"],
      )

    with self.subTest(name="Create status for undeclared code"):
      compare.assertProto2Equal(
          self,
          status_specs.create(99999, "user message").proto,
          extended_status_pb2.ExtendedStatus(
              status_code=extended_status_pb2.StatusCode(
                  component=_COMPONENT, code=99999
              ),
              title=f"Undeclared error {_COMPONENT}:99999",
              user_report=extended_status_pb2.ExtendedStatus.UserReport(
                  message="user message"
              ),
              context=[
                  extended_status_pb2.ExtendedStatus(
                      status_code=extended_status_pb2.StatusCode(
                          component="ai.intrinsic.errors", code=604
                      ),
                      title="Error code not declared",
                      severity=extended_status_pb2.ExtendedStatus.Severity.WARNING,
                      user_report=extended_status_pb2.ExtendedStatus.UserReport(
                          message=(
                              f"The code {_COMPONENT}:99999 has not been"
                              " declared by the component."
                          ),
                          instructions=(
                              "Inform the component owner to add the missing"
                              " error code to the status specs file."
                          ),
                      ),
                  )
              ],
          ),
          ignored_fields=["timestamp"],
      )


if __name__ == "__main__":
  absltest.main()
