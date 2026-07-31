// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/demangle.h"

#include <cxxabi.h>

#include "absl/cleanup/cleanup.h"

namespace intrinsic::details {

std::string Demangle(const char* mangled) {
  // https://github.com/llvm/llvm-project/blob/main/libcxxabi/src/cxa_demangle.cpp
  int status = 0;
  char* demangled = abi::__cxa_demangle(mangled, nullptr, nullptr, &status);

  absl::Cleanup free_mem = [demangled]() {
    if (demangled != nullptr) {
      free(demangled);
    }
  };

  if (status == 0 && demangled != nullptr) {
    return std::string(demangled);
  }
  return {};
}

}  // namespace intrinsic::details
