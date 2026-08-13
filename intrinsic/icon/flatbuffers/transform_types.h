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


#ifndef INTRINSIC_ICON_FLATBUFFERS_TRANSFORM_TYPES_H_
#define INTRINSIC_ICON_FLATBUFFERS_TRANSFORM_TYPES_H_

#include <cstddef>
#include <vector>

#include "Eigen/Dense"
#include "flatbuffers/detached_buffer.h"
#include "intrinsic/icon/flatbuffers/transform_types.fbs.h"

namespace intrinsic_fbs {

// Creates a detached buffer that can store a fixed size Vector of type
// intrinsic_fbs::VectorNd. The Vector is initialized from a reference to
// std::vector<double> input.
flatbuffers::DetachedBuffer CreateVectorNdBuffer(
    const std::vector<double>& data);

// Creates a detached buffer that can store a fixed size Vector of type
// intrinsic_fbs::VectorNd. The Vector length is passed as an input.
flatbuffers::DetachedBuffer CreateVectorNdBuffer(size_t length);

// Copies an std::vector<double> reference to a an existing buffer of a given
// size. Returns a reference to intrinsic_fbs::VectorNd that points to memory
// inside the provided buffer, or return nullptr when buffer is null or too
// small.
VectorNd* CopyToVectorNdBuffer(const std::vector<double>& data, void* buffer,
                               size_t size);

// Creates a detached buffer that can store a fixed size Vector of type
// intrinsic_fbs::Wrench.
flatbuffers::DetachedBuffer CreateWrenchBuffer();

}  // namespace intrinsic_fbs

#endif  // INTRINSIC_ICON_FLATBUFFERS_TRANSFORM_TYPES_H_
