// Copyright 2023 Intrinsic Innovation LLC

#include "incode/middleware/zenoh/imw_zenoh_data_callback_context.h"

using std::string;

namespace intrinsic {

IMWZenohDataCallbackContext::IMWZenohDataCallbackContext(
    IMWZenoh* imw_zenoh_instance, const std::string& subscription_keyexpr)
    : imw_zenoh_instance_(imw_zenoh_instance),
      subscription_keyexpr_(subscription_keyexpr) {}

IMWZenohDataCallbackContext::~IMWZenohDataCallbackContext() {}

}  // namespace intrinsic
