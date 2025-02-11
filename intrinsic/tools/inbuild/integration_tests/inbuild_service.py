# Copyright 2023 Intrinsic Innovation LLC

"""A python service built using inbuild."""

import logging
import sys
import time

from intrinsic.resources.proto import runtime_context_pb2
from intrinsic.tools.inbuild.integration_tests import inbuild_service_pb2


def main():
  logging.info('----------------------------------')
  logging.info('-- Inbuild Python service starting')
  logging.info('----------------------------------')

  with open('/etc/intrinsic/runtime_config.pb', 'rb') as fin:
    context = runtime_context_pb2.RuntimeContext.FromString(fin.read())

  # Parse the configuration
  config = inbuild_service_pb2.InbuildServiceConfig()
  context.config.Unpack(config)

  logging.info('Hello from Python InbuildService: %s', config.bar)

  while True:
    time.sleep(5)


if __name__ == '__main__':
  logging.basicConfig(stream=sys.stderr, level=logging.INFO)
  main()
