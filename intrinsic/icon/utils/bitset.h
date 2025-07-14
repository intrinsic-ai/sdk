// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_UTILS_BITSET_H_
#define INTRINSIC_ICON_UTILS_BITSET_H_

#include <bitset>
#include <concepts>
#include <type_traits>

namespace intrinsic {

namespace internal {

// Helper template to get the size of a bitset, avoids template specialization
// (for bool) of an alias template, which isn't allowed.
template <std::unsigned_integral T>
struct bitset_traits {
  using type = std::bitset<sizeof(T) * 8>;
};

template <>
struct bitset_traits<bool> {
  using type = std::bitset<1>;
};

}  // namespace internal

// An intrinsic::bitset is a std::bitset with the size taken from the type T.
template <std::unsigned_integral T>
using bitset = internal::bitset_traits<T>::type;

// Returns the value of the bitset, cast to the type T, which is type-safe
template <std::unsigned_integral T>
T GetValue(const bitset<T>& t) {
  return static_cast<T>(t.to_ullong());
}

// Returns an intrinsic::bitset with the size of the underlying type of the
// enum.
template <typename Enum>
using from_enum_t = bitset<std::underlying_type_t<Enum>>;

// Returns an intrinsic::bitset with the value of the enum.
template <typename Enum>
from_enum_t<Enum> FromEnum(const Enum& enum_value) {
  return from_enum_t<Enum>(
      static_cast<std::underlying_type_t<Enum>>(enum_value));
}

// Returns the enum value of the bitset.
// Note, that this might return a value that is not in the enum (which is valid
// in C++).
template <typename Enum>
Enum ToEnum(const from_enum_t<Enum>& bitset) {
  return static_cast<Enum>(bitset.to_ullong());
}

}  // namespace intrinsic

#endif  // INTRINSIC_ICON_UTILS_BITSET_H_
