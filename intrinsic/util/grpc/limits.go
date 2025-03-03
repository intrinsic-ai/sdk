// Copyright 2023 Intrinsic Innovation LLC

// Package grpclimits provides constants for recommended gRPC limits.
package grpclimits

// GrpcRecommendedMaxMetadataSoftLimit sets a metadata limit, this includes, for
// example, the size of the information gathered from an absl::Status and
// ExtendedStatus on error. Default is 8KB.  We use a rather large limit because
// this in particular contains ExtendedStatus information which collects traces
// and can be signifantly larger than the default. Any request with metadata
// larger than the hard limit is rejected. Between the soft limit and the hard
// limit, some requests will be rejected.
const GrpcRecommendedMaxMetadataSoftLimit = 512 * 1024

// GrpcRecommendedMaxMetadataHardLimit defines the hard limit (see above).
const GrpcRecommendedMaxMetadataHardLimit = GrpcRecommendedMaxMetadataSoftLimit +
	GrpcRecommendedMaxMetadataSoftLimit/4
