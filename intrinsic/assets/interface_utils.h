// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ASSETS_INTERFACE_UTILS_H_
#define INTRINSIC_ASSETS_INTERFACE_UTILS_H_

#include "absl/status/status.h"
#include "absl/strings/string_view.h"

namespace intrinsic {
namespace assets {

// The prefix used for gRPC service dependencies.
inline constexpr absl::string_view kGrpcUriPrefix = "grpc://";
// The prefix used for proto-based data dependencies.
inline constexpr absl::string_view kDataUriPrefix = "data://";

// Validates an interface name with a protocol prefix.
absl::Status ValidateInterfaceName(absl::string_view uri);

}  // namespace assets
}  // namespace intrinsic

#endif  // INTRINSIC_ASSETS_INTERFACE_UTILS_H_
