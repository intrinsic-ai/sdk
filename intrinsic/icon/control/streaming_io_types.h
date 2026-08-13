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

#ifndef INTRINSIC_ICON_CONTROL_STREAMING_IO_TYPES_H_
#define INTRINSIC_ICON_CONTROL_STREAMING_IO_TYPES_H_

#include <any>
#include <cstdint>
#include <functional>

#include "absl/status/statusor.h"
#include "google/protobuf/any.pb.h"
#include "ortools/base/strong_int.h"

namespace intrinsic::icon {

// Realtime Actions use StreamingInputIds to access streaming inputs. An Action
// factory saves IDs for any inputs. The Action instance can then use those IDs
// to access streaming inputs via the StreamingIoRealtimeAccess object that we
// pass to its Sense() method.
DEFINE_STRONG_INT_TYPE(StreamingInputId, int64_t);

// These definitions are used for storing streaming input parsers / output
// converters. Users do not interact with these definitions directly, but rather
// use wrappers that automatically convert from concrete proto message types to
// google::protobuf::Any and from realtime types to absl::any (and vice versa).
using GenericStreamingInputParser = std::function<absl::StatusOr<std::any>(
    const google::protobuf::Any& streaming_input)>;
using GenericStreamingOutputConverter =
    std::function<absl::StatusOr<google::protobuf::Any>(
        const std::any& streaming_output_any)>;

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_CONTROL_STREAMING_IO_TYPES_H_
