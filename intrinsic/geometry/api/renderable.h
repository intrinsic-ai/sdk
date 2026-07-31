// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_API_RENDERABLE_H_
#define INTRINSIC_GEOMETRY_API_RENDERABLE_H_

#include <string>

namespace intrinsic {

class Renderable {
 public:
  explicit Renderable(std::string glb_string);

  // Returns the serialized GLB representation of the renderable.
  std::string GetGLBString() const;

  bool operator==(const Renderable& other) const;
  bool operator!=(const Renderable& other) const;

 private:
  std::string glb_string_;
};

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_RENDERABLE_H_
