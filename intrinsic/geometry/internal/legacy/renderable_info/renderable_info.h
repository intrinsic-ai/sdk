// Copyright 2023 Intrinsic Innovation LLC

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
