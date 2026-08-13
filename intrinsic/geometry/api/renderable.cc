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
