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

"""Example of retrieving metadata from the ICON Application Layer.

Showcases using the Python client API for the ICON  gRPC service by listing
the available parts and action signatures.
"""

from absl import app
from absl import flags

from intrinsic.icon.python import icon_api
from intrinsic.util.grpc import connection

_HOST = flags.DEFINE_string(
    'host', 'xfa.lan', 'ICON server gRPC connection host'
)
_PORT = flags.DEFINE_integer('port', 17080, 'ICON server gRPC connection port')
_INSTANCE = flags.DEFINE_string(
    'instance',
    None,
    'The instance of ICON, if behind an ingress, you intend to talk to',
)


def main(argv):
  if len(argv) > 1:
    raise app.UsageError('Too many command-line arguments.')

  icon_client = icon_api.Client.connect_with_params(
      connection.ConnectionParams(
          f'{_HOST.value}:{_PORT.value}', _INSTANCE.value, None
      )
  )

  parts = icon_client.list_parts()
  print('Parts:', parts)

  action_signatures = icon_client.list_action_signatures()
  print('Action Signatures:')
  for action_signature in action_signatures:
    print(
        '-',
        action_signature.action_type_name,
        ':',
        action_signature.text_description,
    )


if __name__ == '__main__':
  app.run(main)
