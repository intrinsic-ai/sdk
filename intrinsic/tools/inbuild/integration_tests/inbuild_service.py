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
