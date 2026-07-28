// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_PROTO_REPEATED_FIELD_UTIL_H_
#define INTRINSIC_UTIL_PROTO_REPEATED_FIELD_UTIL_H_

#include <algorithm>
#include <string>

#include "google/protobuf/repeated_ptr_field.h"

namespace intrinsic {
// Remove all elements from the input array for which the input predicate
// pred is true. Returns number of erased elements.
template <typename T, typename Predicate>
int RemoveIf(google::protobuf::RepeatedPtrField<T>* array,
             const Predicate& pred) {
  int i = 0, end = array->size();
  while (i < end && !pred(&array->Get(i))) ++i;
  if (i == end) return 0;
  // 'i' is positioned at first element to be removed.
  for (int j = i + 1; j < end; ++j) {
    if (!pred(&array->Get(j))) array->SwapElements(j, i++);
  }
  array->DeleteSubrange(i, end - i);
  return end - i;
}
template <typename T, typename LessThan>
inline void Sort(google::protobuf::RepeatedPtrField<T>* array,
                 const LessThan& lt) {
  std::sort(array->pointer_begin(), array->pointer_end(), lt);
}

inline void Sort(google::protobuf::RepeatedPtrField<std::string>* array) {
  Sort(array,
       [](const std::string* x, const std::string* y) { return *x < *y; });
}

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_PROTO_REPEATED_FIELD_UTIL_H_
