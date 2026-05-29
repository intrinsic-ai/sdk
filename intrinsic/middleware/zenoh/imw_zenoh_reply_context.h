// Copyright 2023 Intrinsic Innovation LLC

#ifndef MIDDLEWARE_ZENOH_IMW_ZENOH_REPLY_CONTEXT_H_
#define MIDDLEWARE_ZENOH_IMW_ZENOH_REPLY_CONTEXT_H_

#include <string>

struct imw_queryable_options_t;
struct z_loaned_query_t;

namespace intrinsic {

class IMWZenohReplyContext {
 public:
  IMWZenohReplyContext(const z_loaned_query_t* query,
                       const imw_queryable_options_t* options)
      : query_(query), options_(options) {}
  ~IMWZenohReplyContext() = default;

  // Move constructor is fine
  IMWZenohReplyContext(IMWZenohReplyContext&& other) = default;

  // No need for copy constructors in the intended usage of this class.
  IMWZenohReplyContext(const IMWZenohReplyContext&) = delete;
  IMWZenohReplyContext& operator=(const IMWZenohReplyContext&) = delete;

  const z_loaned_query_t* query_;
  const imw_queryable_options_t* options_;
};

}  // namespace intrinsic

#endif  // MIDDLEWARE_ZENOH_IMW_ZENOH_REPLY_CONTEXT_H_
