// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_STORAGE_GZF_STORAGE_H_
#define INTRINSIC_GEOMETRY_STORAGE_GZF_STORAGE_H_

#include <memory>

#include "absl/base/attributes.h"
#include "intrinsic/geometry/storage/geometry_library.h"
#include "intrinsic/world/gzfile/gzfile.h"

namespace intrinsic {

// Returns a GeometryLibrary instance that is backed by `gzfile`.
//
// DOES NOT take ownership of `gzfile`. `gzfile` MUST outlive the
// GeometryLibrary.
std::unique_ptr<GeometryLibrary> GetGzfGeometryLibrary(
    GZFile& gzfile ABSL_ATTRIBUTE_LIFETIME_BOUND);

// Returns a read-only GeometryLibrary instance that is backed by `gzfile`.
//
// DOES NOT take ownership of `gzfile`. `gzfile` MUST outlive the
// GeometryLibrary.
//
// The Library's Serializer() returns UnimplementedError.
std::unique_ptr<GeometryLibrary> GetReadOnlyGzfGeometryLibrary(
    const GZFile& gzfile ABSL_ATTRIBUTE_LIFETIME_BOUND);

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_STORAGE_GZF_STORAGE_H_
