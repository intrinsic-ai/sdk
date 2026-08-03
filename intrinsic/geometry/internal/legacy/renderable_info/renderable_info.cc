// Copyright 2023 Intrinsic Innovation LLC

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
