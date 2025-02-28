// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_GRPC_LIMITS_H_
#define INTRINSIC_UTIL_GRPC_LIMITS_H_

// Metadata limit, this includes, for example, the size of the information
// gathered from an absl::Status and ExtendedStatus on error. Default is 8KB.
// We use a rather large limit because this in particular contains
// ExtendedStatus information which collects traces and can be signifantly
// larger than the default. Any request with metadata larger than the hard
// limit is rejected. Between the soft limit and the hard limit, some requests
// will be rejected.
constexpr int kGrpcRecommendedMaxMetadataSoftLimit = 512 * 1024;  // Bytes
constexpr int kGrpcRecommendedMaxMetadataHardLimit =
    kGrpcRecommendedMaxMetadataSoftLimit +
    kGrpcRecommendedMaxMetadataSoftLimit / 4;

#endif  // INTRINSIC_UTIL_GRPC_LIMITS_H_
