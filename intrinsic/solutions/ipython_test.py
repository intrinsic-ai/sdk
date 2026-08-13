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

import io
import sys
from unittest import mock

from absl.testing import absltest

from intrinsic.solutions import ipython


class IpythonTest(absltest.TestCase):

  def test_no_html_display_outside_ipython(self):
    stdout_mock = io.StringIO()
    with mock.patch.object(sys, 'stdout', stdout_mock):
      ipython.display_html_if_ipython('<span>Some html</span>')

    self.assertEqual(
        stdout_mock.getvalue(), 'Display only executed in IPython.\n'
    )

  def test_no_display_outside_ipython(self):
    stdout_mock = io.StringIO()
    with mock.patch.object(sys, 'stdout', stdout_mock):
      ipython.display_if_ipython('<span>Some html</span>')

    self.assertEqual(
        stdout_mock.getvalue(), 'Display only executed in IPython.\n'
    )


if __name__ == '__main__':
  absltest.main()
