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


#include "intrinsic/icon/flatbuffers/transform_types.h"

#include <cstddef>
#include <cstring>
#include <vector>

#include "flatbuffers/buffer.h"
#include "flatbuffers/detached_buffer.h"
#include "flatbuffers/flatbuffer_builder.h"
#include "intrinsic/icon/flatbuffers/transform_types.fbs.h"

namespace intrinsic_fbs {

flatbuffers::DetachedBuffer CreateVectorNdBuffer(
    const std::vector<double>& data) {
  flatbuffers::FlatBufferBuilder builder;
  builder.Finish(CreateVectorNdDirect(builder, &data));
  return builder.Release();
}

flatbuffers::DetachedBuffer CreateVectorNdBuffer(size_t length) {
  std::vector<double> vector(length, 0);
  flatbuffers::FlatBufferBuilder builder;
  builder.Finish(CreateVectorNdDirect(builder, &vector));
  return builder.Release();
}

intrinsic_fbs::VectorNd* CopyToVectorNdBuffer(const std::vector<double>& data,
                                              void* buffer, size_t size) {
  if (!buffer) {
    return nullptr;
  }

  flatbuffers::FlatBufferBuilder builder;
  builder.Finish(CreateVectorNdDirect(builder, &data));

  if (size < builder.GetSize()) {
    return nullptr;
  }

  memcpy(buffer, builder.GetBufferPointer(), builder.GetSize());
  return flatbuffers::GetMutableRoot<VectorNd>(buffer);
}

flatbuffers::DetachedBuffer CreateWrenchBuffer() {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);
  builder.Finish(builder.CreateStruct(Wrench()));
  return builder.Release();
}

}  // namespace intrinsic_fbs
