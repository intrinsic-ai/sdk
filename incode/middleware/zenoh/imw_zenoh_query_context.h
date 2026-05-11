// Copyright 2023 Intrinsic Innovation LLC

#ifndef MIDDLEWARE_ZENOH_IMW_ZENOH_QUERY_CONTEXT_H_
#define MIDDLEWARE_ZENOH_IMW_ZENOH_QUERY_CONTEXT_H_

#include "incode/middleware/imw.h"

namespace intrinsic {

class IMWZenoh;

class IMWZenohQueryContext {
 public:
  IMWZenohQueryContext(IMWZenoh* imw_zenoh_instance, const char* keyexpr,
                       imw_query_callback_fn* callback,
                       imw_query_on_done_callback_fn* on_done,
                       void* user_context, imw_query_options_t* options)
      : imw_zenoh_instance_(imw_zenoh_instance),
        keyexpr_(keyexpr),
        callback_(callback),
        on_done_(on_done),
        user_context_(user_context) {
    if (options != nullptr) {
      options_ = *options;
    }
  }

  IMWZenoh* imw_zenoh_instance_;
  const char* keyexpr_;
  imw_query_callback_fn* callback_;
  imw_query_on_done_callback_fn* on_done_;
  void* user_context_;
  imw_query_options_t options_;
};

}  // namespace intrinsic

#endif  // MIDDLEWARE_ZENOH_IMW_ZENOH_QUERY_CONTEXT_H_
