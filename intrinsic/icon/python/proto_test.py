# Copyright 2023 Intrinsic Innovation LLC

"""Empty test that imports ICON python protos.

Validates that Python protobuf codegen is working correctly in the released ICON
codebase built with Bazel.
"""

# pylint: disable=unused-import
from intrinsic.icon.proto import cart_space_pb2
from intrinsic.icon.proto import ik_options_pb2
from intrinsic.icon.proto import joint_space_pb2
from intrinsic.icon.proto import matrix_pb2
from intrinsic.icon.proto import part_status_pb2
from intrinsic.icon.proto import streaming_output_pb2
from intrinsic.icon.proto.v1 import service_pb2
from intrinsic.icon.proto.v1 import types_pb2
from absl.testing import absltest
# pylint: enable=unused-import

if __name__ == '__main__':
  absltest.main()
