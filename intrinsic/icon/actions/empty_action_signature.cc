// Copyright 2023 Intrinsic Innovation LLC

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
