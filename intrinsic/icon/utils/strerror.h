// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_UTILS_STRERROR_H_
#define INTRINSIC_ICON_UTILS_STRERROR_H_

#include <cstddef>
#include <cstdio>

#include "intrinsic/icon/utils/fixed_str_cat.h"
#include "intrinsic/icon/utils/fixed_string.h"

namespace intrinsic::icon {

// Provides a thread-safe, allocation-free alternative to ::strerror.
//
// Unlike ::strerror, this function avoids dynamic memory allocation, making it
// safe for use in real-time systems and async-signal handlers. It also bypasses
// ::strerror_r, which suffers from conflicting and incompatible POSIX and GNU
// definitions.
//
// For more details see:
// https://www.club.cc.cmu.edu/~cmccabe/blog_strerror.html.
template <std::size_t N = 128>
icon::FixedString<N> StrError(const int err) {
  if (err >= 0 && err < ::sys_nerr) {
    return icon::FixedStrCat<N>(::sys_errlist[err]);
  }
  return icon::FixedStrCat<N>("Unknown error ", err);
}

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_UTILS_STRERROR_H_
