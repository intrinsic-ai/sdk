# Copyright 2023 Intrinsic Innovation LLC

"""Do-nothing service for testing."""

from collections.abc import Sequence

from absl import app
from absl import flags
from absl import logging
from intrinsic.resources.proto import runtime_context_pb2
from intrinsic.tools.inbuild.cmd.service.test_data import example_service_pb2

_RUNTIME_CONTEXT_FILE = flags.DEFINE_string(
    'runtime_context_file',
    '/etc/intrinsic/runtime_config.pb',
    (
        'Path to the runtime context file containing'
        ' intrinsic_proto.config.RuntimeContext binary proto.'
    ),
)


def main(argv: Sequence[str]) -> None:
  del argv  # unused

  with open(_RUNTIME_CONTEXT_FILE.value, 'rb') as f:
    runtime_context = runtime_context_pb2.RuntimeContext.FromString(f.read())
    logging.info('Runtime context level: %d', runtime_context.level)

  config = example_service_pb2.ExampleConfig()
  if not runtime_context.config.Unpack(config):
    raise RuntimeError('Failed to unpack config.')

  logging.info('--------------------------------')
  logging.info('-- Example service --')
  logging.info('--------------------------------')
  logging.info('Hello world: %s', config.hello_world)

  input('Shall we play a game?')


if __name__ == '__main__':
  app.run(main)
