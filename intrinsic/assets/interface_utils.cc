// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/assets/interface_utils.h"

#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "re2/re2.h"

namespace intrinsic {
namespace assets {
namespace {
static constexpr LazyRE2 kUriRegex = {
    R"(^(grpc://|data://)([A-Za-z_][A-Za-z0-9_]*\.)+[A-Za-z_][A-Za-z0-9_]*$)"};
}

// Validates an interface name with a protocol prefix.
absl::Status ValidateInterfaceName(absl::string_view uri) {
  if (!RE2::FullMatch(uri, *kUriRegex)) {
    return absl::InvalidArgumentError(
        absl::StrCat("Expected URI to be formatted as "
                     "'<protocol>://<package>.<message>', got '",
                     uri, "'"));
  }
  return absl::OkStatus();
}

}  // namespace assets
}  // namespace intrinsic
