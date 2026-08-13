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

"""Do-nothing skill for testing."""

from collections.abc import Sequence

from absl import app
from absl import logging


def main(argv: Sequence[str]) -> None:
  logging.info('--------------------------------')
  logging.info('-- Example skill --')
  logging.info('--------------------------------')
  logging.info('Hello world: %s', argv)

  input('Shall we play a game?')


if __name__ == '__main__':
  app.run(main)
