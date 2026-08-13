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

#ifndef INTRINSIC_ICON_COMMON_PART_PROPERTIES_H_
#define INTRINSIC_ICON_COMMON_PART_PROPERTIES_H_

#include <cstddef>
#include <string>
#include <variant>

#include "absl/container/flat_hash_map.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/time/time.h"
#include "intrinsic/icon/proto/v1/service.pb.h"
#include "ortools/base/strong_int.h"

namespace intrinsic::icon {

// Identifier for a Part Property.
DEFINE_STRONG_INT_TYPE(PartPropertyId, size_t);

// A PartPropertyValue is the in-memory representation of a part property. It is
// stored in the AsyncBuffers that move information about part properties
// between the realtime and non-realtime thread.
// NOTE: We initialize this to one of the supported types at startup time, and
// after that point it should never change types!
using PartPropertyValue = std::variant<bool, double>;

struct TimestampedPartProperties {
  absl::Duration timestamp_control;
  absl::Time timestamp_wall;
  absl::flat_hash_map<std::string,
                      absl::flat_hash_map<std::string, PartPropertyValue>>
      properties;
};

struct PartPropertyMap {
  absl::flat_hash_map<std::string,
                      absl::flat_hash_map<std::string, PartPropertyValue>>
      properties;
};

::intrinsic_proto::icon::v1::PartPropertyValue ToProto(
    const PartPropertyValue& value);

absl::StatusOr<PartPropertyValue> FromProto(
    const ::intrinsic_proto::icon::v1::PartPropertyValue& value);

// Visitor to assign one PropertyValue (variant<double, bool>) to another
// without changing the type of the held value. Returns an error if `dst` holds
// a different type than `src`.
struct AssignPropertyValue {
  // These two implement the happy case and return OkStatus.
  absl::Status operator()(bool src, bool& dst);
  absl::Status operator()(double src, double& dst);
  // These two implement the error case and return InvalidArgumentError.
  absl::Status operator()(double src, bool& dst);
  absl::Status operator()(bool src, double& dst);

  // This is used for the error message in case src has a different type than
  // dst.
  std::string property_name;
};

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_COMMON_PART_PROPERTIES_H_
