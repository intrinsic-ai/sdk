// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_PROTO_STATUS_SPECS_H_
#define INTRINSIC_UTIL_PROTO_STATUS_SPECS_H_

#include <cstdint>

#include "absl/strings/string_view.h"

namespace intrinsic::util::proto {

constexpr absl::string_view kExtendedStatusComponent =
    "ai.intrinsic.util.proto.parsed_type_url";

constexpr uint32_t kInvalidUrlCode = 12001;
constexpr char kInvalidUrlTitle[] = "Invalid type URL";
constexpr char kInvalidUrlInstructions[] =
    "Intrinsic type URLs must start with 'type.intrinsic.ai/' and conform to "
    "the Intrinsic URL specification.";

}  // namespace intrinsic::util::proto

#endif  // INTRINSIC_UTIL_PROTO_STATUS_SPECS_H_
