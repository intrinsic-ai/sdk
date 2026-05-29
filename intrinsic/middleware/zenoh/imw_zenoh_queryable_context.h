// Copyright 2023 Intrinsic Innovation LLC

#ifndef MIDDLEWARE_ZENOH_IMW_ZENOH_QUERYABLE_CONTEXT_H_
#define MIDDLEWARE_ZENOH_IMW_ZENOH_QUERYABLE_CONTEXT_H_

#include <string>

namespace intrinsic {

class IMWZenoh;

class IMWZenohQueryableContext {
 public:
  IMWZenohQueryableContext(IMWZenoh* imw_zenoh_instance,
                           const std::string& queryable_keyexpr);
  ~IMWZenohQueryableContext() = default;

  // Move constructor is fine
  IMWZenohQueryableContext(IMWZenohQueryableContext&& other) = default;

  // No need for copy constructors in the intended usage of this class.
  IMWZenohQueryableContext(const IMWZenohQueryableContext&) = delete;
  IMWZenohQueryableContext& operator=(const IMWZenohQueryableContext&) = delete;

  IMWZenoh* get_imw_zenoh_instance() { return imw_zenoh_instance_; }
  const std::string& get_queryable_keyexpr() { return queryable_keyexpr_; }

 private:
  IMWZenoh* imw_zenoh_instance_;
  std::string queryable_keyexpr_;
};

}  // namespace intrinsic

#endif  // MIDDLEWARE_ZENOH_IMW_ZENOH_QUERYABLE_CONTEXT_H_
