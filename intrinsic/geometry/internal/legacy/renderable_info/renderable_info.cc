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

#include "intrinsic/geometry/internal/legacy/renderable_info/renderable_info.h"

#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/internal/util/eigen_serialization.h"
#include "intrinsic/marshal/riegeli_coder.h"
#include "intrinsic/util/status/status_macros.h"
#include "riegeli/records/record_reader.h"
#include "riegeli/records/record_writer.h"

namespace intrinsic {

// RenderableInfoData
template <>
absl::Status
RiegeliCoder<intrinsic::geometry_legacy::RenderableInfoData>::Encode(
    const intrinsic::geometry_legacy::RenderableInfoData& value,
    riegeli::RecordWriterBase& writer) {
  INTR_RETURN_IF_ERROR(
      RiegeliCoder<std::string>::Encode(value.filename, writer));
  INTR_RETURN_IF_ERROR(RiegeliCoder<intrinsic::eigenmath::Vector3d>::Encode(
      value.scale, writer));
  INTR_RETURN_IF_ERROR(
      RiegeliCoder<std::string>::Encode(value.content, writer));
  return absl::OkStatus();
}

template <>
absl::StatusOr<intrinsic::geometry_legacy::RenderableInfoData>
RiegeliCoder<intrinsic::geometry_legacy::RenderableInfoData>::Decode(
    riegeli::RecordReaderBase& reader) {
  intrinsic::geometry_legacy::RenderableInfoData result;
  INTR_ASSIGN_OR_RETURN(result.filename,
                        RiegeliCoder<std::string>::Decode(reader));
  INTR_ASSIGN_OR_RETURN(
      result.scale,
      RiegeliCoder<intrinsic::eigenmath::Vector3d>::Decode(reader));
  INTR_ASSIGN_OR_RETURN(result.content,
                        RiegeliCoder<std::string>::Decode(reader));
  return result;
}

}  // namespace intrinsic
