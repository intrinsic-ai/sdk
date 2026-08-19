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

"""
Exec Properties used for executors.
"""

# @unsorted-dict-items
DEFAULT = {
    "container-image": "docker://us-central1-docker.pkg.dev/intrinsic-mirror/intrinsic-build-images/bazel-rbe-executor@sha256:c2d50e5f5a3bbea4c47ffcdc1b2755168d1c682ec2b3ba620fb3134eca4bab0d",
    "OSFamily": "Linux",
}

REQUIRES_NETWORK = {
    "dockerNetwork": "standard",
}
