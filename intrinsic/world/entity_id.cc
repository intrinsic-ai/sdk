// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/world/entity_id.h"

#include <algorithm>

#include "intrinsic/production/external/intops/strong_int.h"

namespace intrinsic {

const AttachmentEntityId kRootEntityId = AttachmentEntityId(1);
const EntityId kFirstEntityId =
    std::max(kInvalidEntityId, kRootEntityId.id) + EntityId(1);

}  // namespace intrinsic
