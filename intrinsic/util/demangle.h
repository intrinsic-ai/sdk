// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_DEMANGLE_H_
#define INTRINSIC_UTIL_DEMANGLE_H_

#include <string>
#include <typeinfo>

namespace intrinsic {

namespace details {

std::string Demangle(const char* mangled);

}  // namespace details

// Returns demangled string for the symbol.
// An empty string is returned if demangling fails.
template <typename T>
std::string Demangle() {
  const char* mangled = typeid(T).name();
  return details::Demangle(mangled);
}

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_DEMANGLE_H_
