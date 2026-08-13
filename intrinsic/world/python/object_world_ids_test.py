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

from intrinsic.world.python import object_world_ids


class ObjectWorldIdsTest(absltest.TestCase):

  def test_root_object_name(self):
    self.assertEqual(
        object_world_ids.ROOT_OBJECT_NAME,
        object_world_ids.WorldObjectName('root'),
    )

  def test_root_object_id(self):
    self.assertEqual(
        object_world_ids.ROOT_OBJECT_ID,
        object_world_ids.ObjectWorldResourceId('root'),
    )


if __name__ == '__main__':
  absltest.main()
