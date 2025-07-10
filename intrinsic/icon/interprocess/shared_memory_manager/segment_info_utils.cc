// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/interprocess/shared_memory_manager/segment_info_utils.h"

#include <stddef.h>

#include <cstdint>
#include <string>
#include <vector>

#include "absl/status/statusor.h"
#include "intrinsic/icon/flatbuffers/fixed_string.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/segment_info.fbs.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::icon {

namespace {

inline static absl::StatusOr<std::string> InterfaceNameFromSegment(
    const intrinsic_fbs::SegmentName& name) {
  auto status_or_name = intrinsic_fbs::StringView(name.value());
  if (!status_or_name.ok()) {
    return status_or_name.status();
  }
  return std::string(status_or_name.value());
}

}  // namespace

absl::StatusOr<std::vector<std::string>> GetNamesFromSegmentInfo(
    const intrinsic_fbs::SegmentInfo& segment_info) {
  std::vector<std::string> names;
  names.reserve(segment_info.size());

  for (uint32_t i = 0; i < segment_info.size(); ++i) {
    INTR_ASSIGN_OR_RETURN(
        auto name, InterfaceNameFromSegment(*segment_info.names()->Get(i)));
    names.emplace_back(name);
  }

  return names;
}

absl::StatusOr<std::vector<std::string>> GetNamesFromFileDescriptorNames(
    const intrinsic_fbs::FileDescriptorNames& file_descriptor_names) {
  std::vector<std::string> names;
  names.reserve(file_descriptor_names.size());

  for (uint32_t i = 0; i < file_descriptor_names.size(); ++i) {
    INTR_ASSIGN_OR_RETURN(
        auto name,
        InterfaceNameFromSegment(*file_descriptor_names.names()->Get(i)));
    names.emplace_back(name);
  }

  return names;
}

absl::StatusOr<std::vector<std::string>>
GetRequiredInterfaceNamesFromSegmentInfo(
    const intrinsic_fbs::SegmentInfo& segment_info) {
  std::vector<std::string> names;
  names.reserve(segment_info.size());

  for (uint32_t i = 0; i < segment_info.size(); ++i) {
    const intrinsic_fbs::SegmentName* name = segment_info.names()->Get(i);
    if (!name->must_be_used()) {
      continue;
    }
    INTR_ASSIGN_OR_RETURN(auto interface_name, InterfaceNameFromSegment(*name));
    names.emplace_back(interface_name);
  }

  return names;
}

}  // namespace intrinsic::icon
