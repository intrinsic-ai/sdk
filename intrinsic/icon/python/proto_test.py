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

"""Empty test that imports ICON python protos.

Validates that Python protobuf codegen is working correctly in the released ICON
codebase built with Bazel.
"""

from absl.testing import absltest

# pylint: disable=unused-import
from intrinsic.icon.proto import cart_space_pb2
from intrinsic.icon.proto import ik_options_pb2
from intrinsic.icon.proto import joint_space_pb2
from intrinsic.icon.proto import matrix_pb2
from intrinsic.icon.proto import part_status_pb2
from intrinsic.icon.proto import streaming_output_pb2
from intrinsic.icon.proto.v1 import service_pb2
from intrinsic.icon.proto.v1 import types_pb2

# pylint: enable=unused-import

if __name__ == '__main__':
  absltest.main()
