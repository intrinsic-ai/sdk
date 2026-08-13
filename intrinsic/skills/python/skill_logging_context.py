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

"""SkillLoggingContext for storing skill information useful in logs."""

import dataclasses

from intrinsic.logging.proto import context_pb2


@dataclasses.dataclass(frozen=True)
class SkillLoggingContext:
  """Provides logging information for a skill.

  Attributes:
    data_logger_context: The logging context of the execution.
    skill_id: The id of the skill.
  """

  data_logger_context: context_pb2.Context
  skill_id: str
