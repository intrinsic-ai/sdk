<!--
Copyright 2023 Intrinsic Innovation LLC
-->
# HTTP/JSON APIs for gRPC services tests

The Intrinsic Platform will offer a Bazel macro `intrinsic_http_service` to generate a Service Asset that translates HTTP/JSON requests into gRPC service calls.
This macro does not yet exist, but some of the code for it does.
This folder tests that HTTP/JSON endpoints for a fictional "Inventory Service" work as expected.
