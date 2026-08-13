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

"""BasicComputeContextImpl implementation provided by the code execution service."""

from intrinsic.skills.python import basic_compute_context
from intrinsic.world.python import object_world_client


class BasicComputeContextImpl(basic_compute_context.BasicComputeContext):
  """BasicComputeContext implementation provided by the code execution service.

  Attributes:
    object_world: A client for interacting with the object world.
  """

  @property
  def object_world(self) -> object_world_client.ObjectWorldClient:
    return self._object_world

  def __init__(
      self,
      object_world: object_world_client.ObjectWorldClient,
  ):
    self._object_world = object_world
