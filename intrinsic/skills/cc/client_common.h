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

#ifndef INTRINSIC_SKILLS_CC_CLIENT_COMMON_H_
#define INTRINSIC_SKILLS_CC_CLIENT_COMMON_H_

#include "absl/time/time.h"

namespace intrinsic {
namespace skills {

// Constant `kClientDefaultTimeout` is the default timeout for GRPC requests
// made by client libraries.
constexpr absl::Duration kClientDefaultTimeout = absl::Seconds(180);

}  // namespace skills
}  // namespace intrinsic

#endif  // INTRINSIC_SKILLS_CC_CLIENT_COMMON_H_
