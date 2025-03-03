# Copyright 2023 Intrinsic Innovation LLC

"""Limit recommendations for GRPC connections."""

# Metadata limit, this includes, for example, the size of the information
# gathered from an absl::Status and ExtendedStatus on error. Default is 8KB.
# We use a rather large limit because this in particular contains
# ExtendedStatus information which collects traces and can be signifantly
# larger than the default. Any request with metadata larger than the hard
# limit is rejected. Between the soft limit and the hard limit, some requests
# will be rejected.

GRPC_RECOMMENDED_MAX_METADATA_SOFT_LIMIT = 512 * 1024  # Bytes
GRPC_RECOMMENDED_MAX_METADATA_HARD_LIMIT = (
    GRPC_RECOMMENDED_MAX_METADATA_SOFT_LIMIT
    + GRPC_RECOMMENDED_MAX_METADATA_SOFT_LIMIT / 4
)
