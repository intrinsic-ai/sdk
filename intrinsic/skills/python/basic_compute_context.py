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

"""Defines a context for basic compute methods."""

import abc

from intrinsic.world.python import object_world_client


class BasicComputeContext(abc.ABC):
  """Provides extra metadata and functionality to a basic compute method.

  It is provided by the code execution service to the compute function of a
  Python code execution task in a task node and allows, e.g., access to the
  world.

  Attributes:
    object_world: A client for interacting with the object world.
  """

  @property
  @abc.abstractmethod
  def object_world(self) -> object_world_client.ObjectWorldClient:
    pass
