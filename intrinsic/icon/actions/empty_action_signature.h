// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_ACTIONS_EMPTY_ACTION_SIGNATURE_H_
#define INTRINSIC_ICON_ACTIONS_EMPTY_ACTION_SIGNATURE_H_

#include "intrinsic/icon/proto/v1/types.pb.h"

namespace intrinsic::icon {

static constexpr char kEmptyActionTypeName[] = "intrinsic.empty";
static constexpr char kEmptyActionSlotName[] = "any_part";

intrinsic_proto::icon::v1::ActionSignature GetEmptyActionSignature();

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_ACTIONS_EMPTY_ACTION_SIGNATURE_H_
