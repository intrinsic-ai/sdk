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

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_RENDERABLE_INFO_RENDERABLE_INFO_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_RENDERABLE_INFO_RENDERABLE_INFO_H_

#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/marshal/riegeli_coder.h"
#include "riegeli/records/record_reader.h"
#include "riegeli/records/record_writer.h"

namespace intrinsic {
namespace geometry_legacy {

// RenderableInfoData is deprecated please do not use it in new code.
struct RenderableInfoData {
  // The filename the information originated from.
  std::string filename;

  // The scale to use for loading the renderable info.
  eigenmath::Vector3d scale;

  // The file contents for the renderable information. The format of the data is
  // not known, it is inferred by the filename extension.
  std::string content;
};

}  // namespace geometry_legacy

// RenderableInfoData
template <>
absl::Status
RiegeliCoder<intrinsic::geometry_legacy::RenderableInfoData>::Encode(
    const intrinsic::geometry_legacy::RenderableInfoData& value,
    riegeli::RecordWriterBase& writer);

template <>
absl::StatusOr<intrinsic::geometry_legacy::RenderableInfoData>
RiegeliCoder<intrinsic::geometry_legacy::RenderableInfoData>::Decode(
    riegeli::RecordReaderBase& reader);

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_RENDERABLE_INFO_RENDERABLE_INFO_H_
