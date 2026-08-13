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

#ifndef INTRINSIC_ICON_CONTROL_C_API_WRAPPERS_STRING_WRAPPER_H_
#define INTRINSIC_ICON_CONTROL_C_API_WRAPPERS_STRING_WRAPPER_H_

#include <string>

#include "absl/strings/string_view.h"
#include "intrinsic/icon/control/c_api/c_types.h"

namespace intrinsic::icon {

// Destroys `str`, freeing both the memory for the IntrinsicIconString struct
// itself *and* the memory for its character buffer.
void DestroyString(IntrinsicIconString* str);

// Creates a new IntrinsicIconString on the heap, including a buffer to move the
// contents of `str` into. The result can be passed to API functions for them to
// keep (and eventually destroy using DestroyString() above).
IntrinsicIconString* Wrap(absl::string_view str);

// Wraps a string_view into an IntrinsicIconString that can be passed to API
// functions as an immutable, non-owned parameter.
IntrinsicIconStringView WrapView(absl::string_view str);

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_CONTROL_C_API_WRAPPERS_STRING_WRAPPER_H_
