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

#include "intrinsic/icon/control/c_api/wrappers/string_wrapper.h"

#include <cstring>

#include "absl/strings/string_view.h"
#include "intrinsic/icon/control/c_api/c_types.h"

namespace intrinsic::icon {

void DestroyString(IntrinsicIconString* str) {
  if (str == nullptr) {
    return;
  }
  delete[] str->data;
  delete str;
}

IntrinsicIconString* Wrap(absl::string_view str) {
  char* data = new char[str.size()];
  std::memcpy(data, str.data(), str.size());
  return new IntrinsicIconString({.data = data, .size = str.size()});
}

IntrinsicIconStringView WrapView(absl::string_view str) {
  return {.data = str.data(), .size = str.size()};
}

}  // namespace intrinsic::icon
