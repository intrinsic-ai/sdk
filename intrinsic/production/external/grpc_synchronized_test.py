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

"""Make sure versions of gppcio from PyPI and grpc from github are synchronized."""

import argparse
import pathlib


def parse_args():
  parser = argparse.ArgumentParser()
  parser.add_argument("--module-bazel", required=True, type=pathlib.Path)
  parser.add_argument("--requirements-in", required=True, type=pathlib.Path)
  parser.add_argument("--requirements-txt", required=True, type=pathlib.Path)
  args = parser.parse_args()

  assert args.module_bazel.exists()
  assert args.requirements_in.exists()
  assert args.requirements_txt.exists()

  return args


def main():
  args = parse_args()
  assert 'name = "grpc", version = "1.74.0"' in args.module_bazel.read_text()
  assert "grpcio==1.74.0" in args.requirements_in.read_text()
  assert "grpcio==1.74.0" in args.requirements_txt.read_text()


if __name__ == "__main__":
  main()
