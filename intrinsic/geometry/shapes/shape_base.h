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

#ifndef INTRINSIC_GEOMETRY_SHAPES_SHAPE_BASE_H_
#define INTRINSIC_GEOMETRY_SHAPES_SHAPE_BASE_H_

#include <memory>

#include "absl/log/check.h"
#include "intrinsic/math/pose3.h"

namespace intrinsic {
namespace shapes {

/** Supported shapes */
enum class ShapeType {
  AXES, /** A group of cylinders + arrowheads representing X, Y, and Z axes */
  TRIANGLE_MESH, /** A concave or convex triangle mesh */
  MESHFILE,      /** A concave or convex mesh stored in a file */
  CONVEX_HULL,   /** A convex mesh, defined by a set of 3D points */
  CYLINDER,      /** A cylinder, with length and radius */
  CAPSULE,       /** A capsule, with length and radius */
  SPHERE,        /** A sphere with a radius */
  SPHERES,       /** A set of spheres */
  BOX,           /** A box, defined by its 3 dimensions */
  ELLIPSOID,     /** A ellipsoid, defined by its 3 radii */
  POINT_CLOUD,   /** A point cloud */
  DISCS,         /** A collection of discs */
  LINE_STRIP,    /** A series of vertices connected by straight lines */
  FRUSTUM,  // A pyramid frustum that extends in the +z direction, defined by
            //   its angle around x-axis and y-axis, the min and max z distances
};

/** A base class for geometry shapes */
class ShapeBase {
 public:
  virtual ~ShapeBase() = default;

  /**
   * Gets the shape type
   * @returns the shape type
   */
  ShapeType getType() const { return type_; }

  /**
   * Gets a specific shape type
   * Downcasts the shape to the specific type provided as a template parameter
   * @tparam T the shape type to which to downcast
   * @returns a reference to this object, downcasted to the specific type.
   * Throws if the current shape is not of the provided type
   */
  template <typename T>
  T const& get() const {
    CHECK(T::type == getType()) << "Requesting incompatible Geometry type";
    return static_cast<T const&>(*this);
  }

  virtual std::unique_ptr<ShapeBase> clone() const = 0;

 protected:
  explicit ShapeBase(ShapeType type) : type_(type) {}

  ShapeType type_;
};

// Returns a string representation of `shape_type`.
inline const char* ToString(ShapeType shape_type) {
  switch (shape_type) {
    case ShapeType::AXES:
      return "ShapeType::AXES";
    case ShapeType::TRIANGLE_MESH:
      return "ShapeType::TRIANGLE_MESH";
    case ShapeType::MESHFILE:
      return "ShapeType::MESHFILE";
    case ShapeType::CONVEX_HULL:
      return "ShapeType::CONVEX_HULL";
    case ShapeType::CYLINDER:
      return "ShapeType::CYLINDER";
    case ShapeType::CAPSULE:
      return "ShapeType::CAPSULE";
    case ShapeType::SPHERE:
      return "ShapeType::SPHERE";
    case ShapeType::SPHERES:
      return "ShapeType::SPHERES";
    case ShapeType::BOX:
      return "ShapeType::BOX";
    case ShapeType::ELLIPSOID:
      return "ShapeType::ELLIPSOID";
    case ShapeType::POINT_CLOUD:
      return "ShapeType::POINT_CLOUD";
    case ShapeType::DISCS:
      return "ShapeType::DISCS";
    case ShapeType::LINE_STRIP:
      return "ShapeType::LINE_STRIP";
    case ShapeType::FRUSTUM:
      return "ShapeType::FRUSTUM";
  }
  return "Invalid ShapeType value";
}

}  //  namespace shapes

}  //  namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_SHAPES_SHAPE_BASE_H_
