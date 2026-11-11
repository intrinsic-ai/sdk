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
Module extension for non-module dependencies
"""

def _non_module_deps_impl(ctx):
    # The set of generated repos is deterministic and does not depend on the environment, so
    # there is no need to store this extension in the lockfile. This statement also ensures
    # that the function stays valid if copybara strips everything else out of it.
    return ctx.extension_metadata(reproducible = True)

non_module_deps_ext = module_extension(implementation = _non_module_deps_impl)
