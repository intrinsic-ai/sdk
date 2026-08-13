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

"""Python equivalent to file_helpers.h."""

from google.protobuf import message


def load_binary_proto(path: str, msg: message.Message):
  """Reads a binary proto from a file to and loads it into a message.

  Args:
    path: The path to the binary proto file.
    msg: The message to load the binary proto file into.
  """
  with open(path, 'rb') as fb:
    msg.ParseFromString(fb.read())
