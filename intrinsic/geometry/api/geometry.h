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

#ifndef INTRINSIC_GEOMETRY_API_GEOMETRY_H_
#define INTRINSIC_GEOMETRY_API_GEOMETRY_H_

#include <memory>
#include <optional>
#include <string>
#include <variant>

#include "absl/base/thread_annotations.h"
#include "absl/status/statusor.h"
#include "absl/synchronization/mutex.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/renderable.h"
#include "intrinsic/geometry/proto/v1/material.pb.h"

namespace intrinsic::geo {

// Geometry is an abstraction for use by the intrinsic codebase that allows the
// geometry package as a whole to evolve its internal representations without
// changing call sites. It represents a piece of geometry that can be concrete
// like a mesh or it may be lazily evaluated like a swept volume being
// represented by a generator and path.
class Geometry {
 public:
  // Creates an immutable empty Geometry.
  Geometry();

  // Standard Move/Copy constructors
  Geometry(const Geometry& other);
  Geometry(Geometry&& other) noexcept;

  // Provenance is a way to track the history of a geometry object. It is
  // intended to be used to allow the user to revert to a previous state if
  // some operation was not as expected, or no longer needed.
  struct Provenance {
    // The human readable reason for the update that led from the previous
    // geometry to this geometry.
    std::string human_readable_update_reason;

    // A pointer to the previous geometry that was used to create this geometry.
    // This can be a string representing a URI to the previous geometry proto
    // stored somewhere else, or a shared pointer to the previous geometry if it
    // is available.
    std::variant<std::string, std::shared_ptr<const Geometry>>
        previous_geometry;
  };

  // Standard constructors
  explicit Geometry(const ExactGeometry& exact_geometry,
                    std::optional<Provenance> provenance);
  Geometry(const ExactGeometry& exact_geometry,
           std::shared_ptr<const Renderable> renderable,
           std::optional<Provenance> provenance);

  // If keep_renderable is true we will keep the given renderable as user
  // specified and not generated unless it is nullptr, we will also store the
  // renderable when saving. If keep_renderable is false or the given renderable
  // is nullptr we may generate a renderable in the future and not store it.
  Geometry(const ExactGeometry& exact_geometry,
           std::shared_ptr<const Renderable> renderable,
           bool keep_renderable = false,
           std::optional<intrinsic_proto::geometry::v1::MaterialProperties>
               material_properties = std::nullopt,
           std::optional<Provenance> provenance = std::nullopt);

  // Standard Move/Copy assignment operators
  Geometry& operator=(const Geometry& other);
  Geometry& operator=(Geometry&& other) noexcept;

  bool operator==(const Geometry& other) const;
  bool operator!=(const Geometry& other) const;

  // Returns the optional renderable representing this Geometry object if it
  // exists, nullptr otherwise.
  std::shared_ptr<const Renderable> GetRenderable() const;

  // Returns the ExactGeometry represented by this Geometry.
  const ExactGeometry& GetExactGeometry() const;

  // Returns true if the Renderable should be kept during serialization.
  bool KeepRenderableForSerialization() const;

  std::optional<intrinsic_proto::geometry::v1::MaterialProperties>
  material_properties() const;

  std::optional<Provenance> provenance() const;

 private:
  friend absl::StatusOr<std::string> GenerateFingerprint(
      const Geometry& geometry);
  mutable absl::Mutex fingerprint_mutex_;
  mutable std::string fingerprint_ ABSL_GUARDED_BY(fingerprint_mutex_);

 private:
  std::shared_ptr<const ExactGeometry> exact_geometry_;

  // Optional Renderable that was precomputed for this instance.
  std::shared_ptr<const Renderable> renderable_;

  // If set we will keep the renderable on serialization.
  bool keep_renderable_ = false;
  std::optional<intrinsic_proto::geometry::v1::MaterialProperties>
      material_properties_;

  // The provenance of this geometry, aka where did it come from. This is
  // optional and may be lost during some operations that do not support it yet.
  // It is intended to allow the user to revert to a previous state if some
  // operation was not as expected, or no longer needed.
  std::optional<Provenance> provenance_;
};

}  // namespace intrinsic::geo

namespace intrinsic {
using ::intrinsic::geo::Geometry;
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_GEOMETRY_H_
