// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/api/renderable.h"

#include <string>
#include <utility>

namespace intrinsic {

Renderable::Renderable(std::string glb_string)
    : glb_string_(std::move(glb_string)) {}

std::string Renderable::GetGLBString() const { return glb_string_; }

bool Renderable::operator==(const Renderable& other) const {
  return glb_string_ == other.glb_string_;
}

bool Renderable::operator!=(const Renderable& other) const {
  return !(*this == other);
}

}  // namespace intrinsic
