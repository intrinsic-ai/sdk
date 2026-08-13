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

from intrinsic.skills.internal import error_bindings

from pybind11_abseil import status  # isort: skip


class ErrorBindingsPyTest(absltest.TestCase):

  def test_raise_status(self):
    # The return_status function should convert a non-ok status to an exception.
    with self.assertRaises(status.StatusNotOk) as cm:
      error_bindings.raise_status(status.StatusCode.CANCELLED, 'test')
    self.assertEqual(cm.exception.status.code(), status.StatusCode.CANCELLED)
    self.assertEqual(cm.exception.status.message(), 'test')

  def test_raise_status_with_raise(self):
    # While the 'raise' is not necessary, it helps the python type check in some
    # cases, so we'll confirm that we can get reasonable results with this code.
    with self.assertRaises(status.StatusNotOk) as cm:
      raise error_bindings.raise_status(status.StatusCode.CANCELLED, 'test')
    self.assertEqual(cm.exception.status.code(), status.StatusCode.CANCELLED)
    self.assertEqual(cm.exception.status.message(), 'test')


if __name__ == '__main__':
  absltest.main()
