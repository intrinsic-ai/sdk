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

#ifndef INTRINSIC_LOGGING_UTILS_DOWNSAMPLER_PROTO_CONVERSION_H_
#define INTRINSIC_LOGGING_UTILS_DOWNSAMPLER_PROTO_CONVERSION_H_

#include "absl/status/statusor.h"
#include "google/protobuf/duration.pb.h"
#include "google/protobuf/timestamp.pb.h"
#include "intrinsic/logging/proto/downsampler.pb.h"
#include "intrinsic/logging/utils/downsampler/downsampler.h"

// All conversions from data_logger Intrinsic protos to their respective C++
// types should be declared in namespace intrinsic_proto::data_logger.
//
// This makes it possible to make unqualified calls to FromProto() throughout
// our code base via argument-dependent name lookup (ADL).
namespace intrinsic_proto::data_logger {

absl::StatusOr<intrinsic::data_logger::DownsamplerOptions> FromProto(
    const intrinsic_proto::data_logger::DownsamplerOptions& proto);

absl::StatusOr<intrinsic::data_logger::DownsamplerEventSourceState> FromProto(
    const intrinsic_proto::data_logger::DownsamplerEventSourceState& proto);

absl::StatusOr<intrinsic::data_logger::DownsamplerState> FromProto(
    const intrinsic_proto::data_logger::DownsamplerState& proto);

}  // namespace intrinsic_proto::data_logger

// All conversions to data_logger Intrinsic protos from their respective C++
// types should be declared in namespace intrinsic::data_logger.
//
// This makes it possible to make unqualified calls to ToProto() throughout our
// code base via argument-dependent name lookup (ADL).
namespace intrinsic::data_logger {

absl::StatusOr<intrinsic_proto::data_logger::DownsamplerOptions> ToProto(
    const intrinsic::data_logger::DownsamplerOptions& options);

absl::StatusOr<intrinsic_proto::data_logger::DownsamplerEventSourceState>
ToProto(const intrinsic::data_logger::DownsamplerEventSourceState& state);

absl::StatusOr<intrinsic_proto::data_logger::DownsamplerState> ToProto(
    const intrinsic::data_logger::DownsamplerState& state);

}  // namespace intrinsic::data_logger

#endif  // INTRINSIC_LOGGING_UTILS_DOWNSAMPLER_PROTO_CONVERSION_H_
