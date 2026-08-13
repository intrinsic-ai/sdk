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

#include "intrinsic/icon/cc_client/client_utils.h"

#include <string>

#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/cc_client/client.h"
#include "intrinsic/icon/cc_client/condition.h"
#include "intrinsic/icon/common/builtins.h"
#include "intrinsic/icon/common/part_properties.h"

namespace intrinsic {
namespace icon {

Comparison IsDone() { return IsTrue(kIsDone); }

absl::Status SetPartProperty(const Client& client, absl::string_view part_name,
                             absl::string_view property_name,
                             double double_value) {
  return client.SetPartProperties(PartPropertyMap{
      .properties = {{std::string(part_name),
                      {{std::string(property_name), double_value}}}}});
}

absl::Status SetPartProperty(const Client& client, absl::string_view part_name,
                             absl::string_view property_name, bool bool_value) {
  return client.SetPartProperties(PartPropertyMap{
      .properties = {{std::string(part_name),
                      {{std::string(property_name), bool_value}}}}});
}

}  // namespace icon
}  // namespace intrinsic
