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

"""Utility functions for skills related proto conversions."""

from typing import TypeVar

from google.protobuf import any_pb2
from google.protobuf import message

T = TypeVar("T", bound=message.Message)


class ProtoMismatchTypeError(TypeError):
  """An unexpected proto type was specified."""


def unpack_any(any_message: any_pb2.Any, proto_message: T) -> T:
  """Unpacks a proto Any into a message.

  The message passed must be the same type as is contained in the proto Any.

  Args:
    any_message: a proto Any message.
    proto_message: A proto message of the type contained in the any to unpack
      the message into.

  Returns:
    The unpacked proto message.

  Raises:
    ProtoMismatchTypeError: If the type of `proto_message` does not match that
      of the specified Any proto.
  """

  if not any_message.Unpack(proto_message):
    any_type_name = any_message.TypeName()
    any_type_msg = (
        f"of type {any_type_name}" if any_type_name else "with no contents"
    )

    raise ProtoMismatchTypeError(
        f"Cannot unpack Any {any_type_msg} into message of type "
        f"{proto_message.DESCRIPTOR.name}."
    )

  return proto_message
