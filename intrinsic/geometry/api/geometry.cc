// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/api/geometry.h"

#include <memory>
#include <optional>
#include <string>
#include <utility>

#include "absl/log/die_if_null.h"
#include "absl/strings/str_cat.h"
#include "absl/synchronization/mutex.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/renderable.h"
#include "intrinsic/geometry/proto/v1/material.pb.h"
#include "intrinsic/util/object_store/object_ref.h"

namespace intrinsic {
namespace {

std::shared_ptr<const ExactGeometry> SharedFromRef(
    const ExactGeometry& lazy_exact_geo) {
  auto weak_ptr = lazy_exact_geo.weak_from_this();
  if (weak_ptr.expired()) {
    return std::make_shared<ExactGeometry>(lazy_exact_geo);
  }

  return weak_ptr.lock();
}

}  // namespace

Geometry::Geometry()
    : exact_geometry_(
          std::make_shared<ExactGeometry>(ExactGeometry::CreateEmpty())) {}

Geometry::Geometry(const Geometry& other)
    : exact_geometry_(other.exact_geometry_),
      renderable_(other.renderable_),
      keep_renderable_(other.keep_renderable_),
      material_properties_(other.material_properties_),
      provenance_(other.provenance_) {
  std::string fingerprint;
  {
    absl::MutexLock lock(other.fingerprint_mutex_);
    fingerprint = other.fingerprint_;
  }

  absl::MutexLock lock(fingerprint_mutex_);
  fingerprint_ = std::move(fingerprint);
}

Geometry::Geometry(Geometry&& other) noexcept
    : exact_geometry_(std::move(other.exact_geometry_)),
      renderable_(std::move(other.renderable_)),
      keep_renderable_(std::move(other.keep_renderable_)),
      material_properties_(std::move(other.material_properties_)),
      provenance_(std::move(other.provenance_)) {
  std::string fingerprint;
  {
    absl::MutexLock lock(other.fingerprint_mutex_);
    fingerprint = std::move(other.fingerprint_);
  }

  absl::MutexLock lock(fingerprint_mutex_);
  fingerprint_ = std::move(fingerprint);
}

Geometry::Geometry(const ExactGeometry& exact_geometry,
                   std::optional<Provenance> provenance)
    : exact_geometry_(ABSL_DIE_IF_NULL(SharedFromRef(exact_geometry))),
      provenance_(std::move(provenance)) {}

Geometry::Geometry(const ExactGeometry& exact_geometry,
                   std::shared_ptr<const Renderable> renderable,
                   std::optional<Provenance> provenance)
    : exact_geometry_(ABSL_DIE_IF_NULL(SharedFromRef(exact_geometry))),
      renderable_(std::move(renderable)),
      keep_renderable_(renderable_ != nullptr),
      provenance_(std::move(provenance)) {}

Geometry::Geometry(
    const ExactGeometry& exact_geometry,
    std::shared_ptr<const Renderable> renderable, bool keep_renderable,
    std::optional<intrinsic_proto::geometry::v1::MaterialProperties>
        material_properties,
    std::optional<Provenance> provenance)
    : exact_geometry_(ABSL_DIE_IF_NULL(SharedFromRef(exact_geometry))),
      renderable_(std::move(renderable)),
      keep_renderable_(keep_renderable && renderable_ != nullptr),
      material_properties_(std::move(material_properties)),
      provenance_(std::move(provenance)) {}

Geometry& Geometry::operator=(const Geometry& other) {
  exact_geometry_ = other.exact_geometry_;
  renderable_ = other.renderable_;
  keep_renderable_ = other.keep_renderable_;
  material_properties_ = other.material_properties_;
  provenance_ = other.provenance_;
  std::string fingerprint;
  {
    absl::MutexLock lock(other.fingerprint_mutex_);
    fingerprint = other.fingerprint_;
  }

  absl::MutexLock lock(fingerprint_mutex_);
  fingerprint_ = std::move(fingerprint);
  return *this;
}

Geometry& Geometry::operator=(Geometry&& other) noexcept {
  exact_geometry_ = std::move(other.exact_geometry_);
  renderable_ = std::move(other.renderable_);
  keep_renderable_ = std::move(other.keep_renderable_);
  material_properties_ = std::move(other.material_properties_);
  provenance_ = std::move(other.provenance_);

  std::string fingerprint;
  {
    absl::MutexLock lock(other.fingerprint_mutex_);
    fingerprint = std::move(other.fingerprint_);
  }

  absl::MutexLock lock(fingerprint_mutex_);
  fingerprint_ = std::move(fingerprint);
  return *this;
}

bool Geometry::operator==(const Geometry& other) const {
  // If this and the other are the exact same instance immediately return, we
  // would otherwise deadlock below.
  if (this == &other) {
    return true;
  }

  if (keep_renderable_ != other.keep_renderable_) {
    return false;
  }

  if (keep_renderable_) {
    // If one has a renderable and the other doesn't or vice versa then they
    // are not equal.
    if ((renderable_ == nullptr) != (other.renderable_ == nullptr)) {
      return false;
    }

    if (renderable_ != nullptr && *renderable_ != *other.renderable_) {
      return false;
    }
  }

  return *exact_geometry_ == *other.exact_geometry_;
}

bool Geometry::operator!=(const Geometry& other) const {
  return !(*this == other);
}

std::shared_ptr<const Renderable> Geometry::GetRenderable() const {
  return renderable_;
}

const ExactGeometry& Geometry::GetExactGeometry() const {
  return *exact_geometry_;
}

bool Geometry::KeepRenderableForSerialization() const {
  return keep_renderable_;
}

std::optional<intrinsic_proto::geometry::v1::MaterialProperties>
Geometry::material_properties() const {
  return material_properties_;
}

std::optional<Geometry::Provenance> Geometry::provenance() const {
  return provenance_;
}

}  // namespace intrinsic
