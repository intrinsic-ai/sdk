<!--
Copyright 2026 Intrinsic Innovation LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->
# HTTP/JSON APIs for gRPC services tests

The Intrinsic Platform will offer a Bazel macro `intrinsic_http_service` to generate a Service Asset that translates HTTP/JSON requests into gRPC service calls.
This macro does not yet exist, but some of the code for it does.
This folder tests that HTTP/JSON endpoints for a fictional "Inventory Service" work as expected.
