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

#ifndef INTRINSIC_GEOMETRY_GLTF_GLTF_BUFFER_TEST_UTIL_H_
#define INTRINSIC_GEOMETRY_GLTF_GLTF_BUFFER_TEST_UTIL_H_

// Helpers for the glTF buffer and validation tests.

#include <array>
#include <cstddef>

#include "absl/types/span.h"
#include "intrinsic/geometry/gltf/gltf_buffer_writer.h"
#include "tiny_gltf.h"

namespace intrinsic::geo {

// The 4 corners of a unit quad, as float VEC3s.
inline constexpr std::array<float, 12> kQuadPositions = {
    0.0f, 0.0f, 0.0f, 1.0f, 0.0f, 0.0f, 0.0f, 1.0f, 0.0f, 1.0f, 1.0f, 0.0f,
};

// Appends `data` to buffer 0 as a bufferView and an accessor over it, and
// returns the accessor index.
template <typename Container>
int AddAccessor(tinygltf::Model& model, const Container& data,
                int component_type, int type, size_t count,
                bool normalized = false) {
  const auto span = absl::MakeConstSpan(data);
  tinygltf::Accessor acc;
  acc.bufferView =
      AppendGltfBufferView(model, /*buffer_index=*/0, span.data(),
                           span.size() * sizeof(typename Container::value_type),
                           /*target=*/0)
          .value();
  acc.componentType = component_type;
  acc.type = type;
  acc.count = count;
  acc.normalized = normalized;
  model.accessors.push_back(acc);
  return static_cast<int>(model.accessors.size() - 1);
}

// Returns a model with one quad primitive (4 vertices, no indices) in mesh 0.
inline tinygltf::Model MakeQuadModel() {
  tinygltf::Model model;
  model.buffers.emplace_back();
  tinygltf::Primitive prim;
  prim.attributes["POSITION"] =
      AddAccessor(model, kQuadPositions, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 4);
  model.meshes.emplace_back().primitives.push_back(prim);
  return model;
}

}  // namespace intrinsic::geo

#endif  // INTRINSIC_GEOMETRY_GLTF_GLTF_BUFFER_TEST_UTIL_H_
