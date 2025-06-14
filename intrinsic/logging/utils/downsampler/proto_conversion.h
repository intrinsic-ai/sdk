// Copyright 2023 Intrinsic Innovation LLC

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
