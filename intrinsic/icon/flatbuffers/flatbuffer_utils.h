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

#ifndef INTRINSIC_ICON_FLATBUFFERS_FLATBUFFER_UTILS_H_
#define INTRINSIC_ICON_FLATBUFFERS_FLATBUFFER_UTILS_H_

#include <cstddef>
#include <cstdint>
#include <cstring>

#include "flatbuffers/array.h"
#include "flatbuffers/vector.h"
#include "intrinsic/icon/utils/realtime_status.h"
namespace intrinsic_fbs {

// Helper function to get the number of elements of a flatbuffer struct's array
// member at compile time.
//
// This function takes a parameter `array_getter_ptr` of type
//                                             ╲
// const ArrayType* (FlatbufferStructT::*array_getter_ptr)() const
// └───────┬──────┘  └─────────┬────────┘                      │
//         │                   ├───────────────────────────────┘
//         │       pointer to const member function of FlatbufferStructT...
//         │
//  ...that returns const ArrayType*
//
// Usage:
//
// constexpr kMyArraySize =
//     FlatbufferArrayNumElements(
//         &::my::fbs::namespace::MyStruct::array_member);
//
// NOTE: You do *not* need to specify any of the template parameters, the
//       compiler can infer all three from the function pointer you pass into
//       this function.
template <typename FlatbufferStructT, typename ArrayT, uint16_t array_length>
constexpr size_t FlatbufferArrayNumElements(
    const flatbuffers::Array<ArrayT, array_length>* (
        FlatbufferStructT::*array_getter_ptr)() const) {
  return array_length;
}

// Performs a deep copy of a flatbuffers::Vector.
//
// The template parameter `T` must be a simple flatbuffer type that can be
// copied trough std::copy.
// Returns OutOfRangeError if the `from` vector and the
// `to`vector have different sizes.
template <typename T>
intrinsic::icon::RealtimeStatus CopyFbsVector(
    const flatbuffers::Vector<T>& from, flatbuffers::Vector<T>& to) {
  if (from.size() != to.size()) {
    return intrinsic::icon::OutOfRangeError(
        intrinsic::icon::RealtimeStatus::StrCat(
            "Vector sizes are not equal: ", from.size(), " != ", to.size()));
  }
  std::copy(from.data(), from.data() + from.size(), to.data());
  return intrinsic::icon::OkStatus();
}

}  // namespace intrinsic_fbs

#endif  // INTRINSIC_ICON_FLATBUFFERS_FLATBUFFER_UTILS_H_
