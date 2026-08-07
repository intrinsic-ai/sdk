// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/api/apply_transform.h"

#include <memory>
#include <optional>
#include <string>
#include <utility>
#include <vector>

#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/api/geometry_options.h"
#include "intrinsic/geometry/api/renderable.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"
#include "intrinsic/geometry/internal/legacy/utils/export_as_gltf.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/util/object_store/object_ref.h"
#include "intrinsic/util/object_store/object_store.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {
namespace {

using ::intrinsic::geometry_legacy::Mesh;

class ApplyTransformFunctor {
 public:
  explicit ApplyTransformFunctor(const Pose3d& ref_t_geo)
      : ref_t_geo_(ref_t_geo.matrix()) {}

  explicit ApplyTransformFunctor(const eigenmath::Matrix4d& ref_t_geo)
      : ref_t_geo_(ref_t_geo) {}

  template <typename Geo>
  absl::StatusOr<ObjectRef<Mesh>> operator()(const Geo& geo) const {
    LOG(WARNING)
        << "WARNING: intrinsic::ApplyTransformFunctor does not support\n"
        << "operator()(const Geo&)) const, with \n"
        << "Geo = " << Demangle<Geo>();
    return absl::InvalidArgumentError(
        "ApplyTransformFunctor does not support argument type.");
  }

  template <>
  absl::StatusOr<ObjectRef<Mesh>> operator()(const ObjectRef<Mesh>& geo) const {
    Mesh mesh = geo.Value().Clone();
    mesh.Transform(ref_t_geo_);
    return DeDuplicate(std::move(mesh));
  }

  eigenmath::AffineTransform3d ref_t_geo_;
};

}  // namespace

absl::StatusOr<Geometry> ApplyTransform(const Geometry& geo,
                                        const Pose3d& ref_t_geo) {
  return ApplyTransform(geo, ref_t_geo.matrix());
}

absl::StatusOr<Geometry> ApplyTransform(const Geometry& geo,
                                        const eigenmath::Matrix4d& ref_t_geo) {
  if (ref_t_geo.isApprox(eigenmath::Matrix4d::Identity())) {
    return geo;
  }

  std::optional<ExactGeometry> exact_geo = std::nullopt;

  const std::vector<TransformedPrimitiveShapePtr>& primitive_shapes =
      geo.GetExactGeometry().GetPrimitiveShapes();
  if (!primitive_shapes.empty()) {
    std::vector<TransformedPrimitiveShapePtr> updated_shapes;
    updated_shapes.reserve(primitive_shapes.size());
    for (const auto& shape : primitive_shapes) {
      updated_shapes.emplace_back(shape.shape(),
                                  ref_t_geo * shape.ref_t_shape());
    }
    INTR_ASSIGN_OR_RETURN(
        exact_geo, ExactGeometry::Create(std::move(updated_shapes),
                                         geo.GetExactGeometry().options()));
  } else {
    INTR_ASSIGN_OR_RETURN(
        ObjectRef<Mesh> mesh,
        geo.GetExactGeometry().visit(ApplyTransformFunctor(ref_t_geo)));
    INTR_ASSIGN_OR_RETURN(
        exact_geo, ExactGeometry::Create(std::move(mesh),
                                         geo.GetExactGeometry().options()));
  }

  std::shared_ptr<const Renderable> renderable;
  auto original_renderable = geo.GetRenderable();
  if (geo.KeepRenderableForSerialization() && original_renderable != nullptr) {
    INTR_ASSIGN_OR_RETURN(std::string updated_glb_string,
                          geometry_legacy::ExportAsGltf(
                              original_renderable->GetGLBString(), ref_t_geo));
    renderable = std::make_shared<Renderable>(std::move(updated_glb_string));
  }

  Geometry::Provenance provenance = {
      .human_readable_update_reason = "Apply transform",
      .previous_geometry = std::make_shared<const Geometry>(geo),
  };

  return Geometry(std::move(exact_geo).value(), std::move(renderable),
                  geo.KeepRenderableForSerialization(),
                  geo.material_properties(), std::move(provenance));
}

}  // namespace intrinsic
