# Copyright 2023 Intrinsic Innovation LLC

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
