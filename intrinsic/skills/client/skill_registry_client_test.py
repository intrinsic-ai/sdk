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

from intrinsic.skills.client import skill_registry_client
from intrinsic.skills.proto import skill_registry_pb2
from intrinsic.skills.proto import skills_pb2


class SkillRegistryTest(absltest.TestCase):
  """Tests all public methods of the SkillRegistry gRPC wrapper class."""

  def setUp(self):
    super().setUp()

    self._skill_registry_stub = mock.MagicMock()
    self._client = skill_registry_client.SkillRegistryClient(
        self._skill_registry_stub
    )

  def test_get_skills_works(self):
    self._skill_registry_stub.GetSkills.return_value = (
        skill_registry_pb2.GetSkillsResponse(
            skills=[
                skills_pb2.Skill(id='ai.intrinsic.throw_ball'),
                skills_pb2.Skill(id='ai.intrinsic.catch_ball'),
            ]
        )
    )

    result = self._client.get_skills()

    self.assertLen(result, 2)
    self.assertEqual(result[0].id, 'ai.intrinsic.throw_ball')
    self.assertEqual(result[1].id, 'ai.intrinsic.catch_ball')

  def test_get_skill_works(self):
    self._skill_registry_stub.GetSkill.return_value = (
        skill_registry_pb2.GetSkillResponse(
            skill=skills_pb2.Skill(id='ai.intrinsic.throw_ball')
        )
    )

    result = self._client.get_skill('ai.intrinsic.throw_ball')

    self.assertEqual(result.id, 'ai.intrinsic.throw_ball')
    self._skill_registry_stub.GetSkill.assert_called_once_with(
        skill_registry_pb2.GetSkillRequest(id='ai.intrinsic.throw_ball')
    )


if __name__ == '__main__':
  absltest.main()
