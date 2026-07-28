// Copyright 2023 Intrinsic Innovation LLC

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
