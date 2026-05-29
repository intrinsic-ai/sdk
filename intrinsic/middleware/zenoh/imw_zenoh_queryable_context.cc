// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/middleware/zenoh/imw_zenoh_queryable_context.h"

using std::string;

namespace intrinsic {

IMWZenohQueryableContext::IMWZenohQueryableContext(
    IMWZenoh* imw_zenoh_instance, const std::string& queryable_keyexpr)
    : imw_zenoh_instance_(imw_zenoh_instance),
      queryable_keyexpr_(queryable_keyexpr) {}

}  // namespace intrinsic
