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

#include "intrinsic/icon/actions/empty_action_signature.h"

#include "absl/log/check.h"
#include "absl/strings/str_cat.h"
#include "intrinsic/icon/actions/action_utils.h"
#include "intrinsic/icon/common/builtins.h"
#include "intrinsic/icon/proto/v1/types.pb.h"

namespace intrinsic::icon {

intrinsic_proto::icon::v1::ActionSignature GetEmptyActionSignature() {
  ActionSignatureBuilder builder(
      kEmptyActionTypeName,
      absl::StrCat("This Action takes no parameters, does nothing, and only has"
                   " a single state variable '",
                   kIsDone, "'"));
  CHECK_OK(builder.AddStateVariable<intrinsic_proto::icon::v1::ActionSignature::
                                        StateVariableInfo::TYPE_BOOL>(
      kIsDone, kIsDoneDescription));
  CHECK_OK(builder.AddPartSlot(
      kEmptyActionSlotName,
      "Any Part is technically compatible with this Action, since the Action "
      "does nothing. However, some Parts may raise errors if they do not "
      "receive a command.",
      /*required_feature_interfaces=*/{}));

  CHECK_OK(builder.AddSupportedBehaviorOverride(
      intrinsic_proto::icon::v1::BehaviorOverrideRequest::
          BEHAVIOR_OVERRIDE_REQUEST_PAUSE,
      "No change in behavior, the action continues doing nothing."));
  return builder.Finish();
}

}  // namespace intrinsic::icon
