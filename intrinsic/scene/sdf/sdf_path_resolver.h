// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#ifndef INTRINSIC_SCENE_SDF_SDF_PATH_RESOLVER_H_
#define INTRINSIC_SCENE_SDF_SDF_PATH_RESOLVER_H_

#include <functional>
#include <string>

#include "absl/status/statusor.h"

namespace intrinsic {
namespace sdf {

// An interface used for resolving URIs encountered during the reading process.
// Note that it is slightly unfortunate that the resolution must result in a
// local filename rather than an abstract input stream, but this is the result
// of some implementation constraints of the underlying APIs.
using UriResolver =
    std::function<absl::StatusOr<std::string>(const std::string&)>;

/**
 * Generic resolver that will attempt to resolve runfiles or external
 * directories if available.
 */
absl::StatusOr<std::string> SdfPathResolver(const std::string& uri);

}  // namespace sdf
}  // namespace intrinsic

#endif  // INTRINSIC_SCENE_SDF_SDF_PATH_RESOLVER_H_
