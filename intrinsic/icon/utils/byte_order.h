// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_UTILS_BYTE_ORDER_H_
#define INTRINSIC_ICON_UTILS_BYTE_ORDER_H_

#include <algorithm>
#include <bit>
#include <concepts>
#include <cstdint>
#include <type_traits>

#include "absl/base/attributes.h"
#include "absl/numeric/bits.h"
#include "absl/types/span.h"

namespace intrinsic::icon {

// Converts unsigned integer in native endianness to desired endianness by
// swapping bytes if necessary. Can also be used for reverting the desired
// endianness to the native one.
template <typename T>
  requires std::unsigned_integral<T>
[[nodiscard]]
inline constexpr auto SwapBytesIfNeeded(const T value,
                                        const absl::endian endian) noexcept {
  if (absl::endian::native != endian) {
    return absl::byteswap(value);
  } else {
    return value;
  }
}

// Determine the unsigned integer type for a given size in bytes.
template <std::size_t N>
struct UIntOfSize;
template <>
struct UIntOfSize<1> {
  using type = uint8_t;
};
template <>
struct UIntOfSize<2> {
  using type = uint16_t;
};
template <>
struct UIntOfSize<4> {
  using type = uint32_t;
};
template <>
struct UIntOfSize<8> {
  using type = uint64_t;
};

template <std::size_t N>
using UIntOfSize_t = typename UIntOfSize<N>::type;

namespace internal {

// Stores arithmetic types and enums as unsigned integers with a desired
// endianness in memory. This is convenient when creating messages for protocols
// that require network order.
//
// Example:
//
//  enum class PacketType: uint32_t { kStatusPacket = 0, kCommandPacket = 1 };
//  enum class Version: uint32_t { kNone = 0, kVersion1 = 1, kVersion2 = 2 };
//
//  struct ABSL_ATTRIBUTE_PACKED CommandPacket {
//    // BigEndian<T> can be initialized via copy-initialization:
//    BigEndian<PacketType> pkt_type = PacketType::kCommandPacket;
//    BigEndian<Version> version = Version::kVersion2;
//    BigEndian<uint32_t> sequence_no = 0;
//    BigEndian<double> command[9] = {0, 0, 0, 0, 0, 0, 0, 0, 0};
//  };
//
//  CommandPacket pkt;
//  // Conversion from T to BigEndian<T> via copy-assignment:
//  pkt.sequence_no = 3;
//  // For vectors and arrays:
//  const std::array<double, 9> command = {1, 2.1, 3, 4.2, 5, 6.3, 7, 8.4, 9};
//  bool success {
//      CopyTo(absl::MakeConstSpan(command), absl::MakeSpan(pkt.command))};
//
//  // Conversion back to T can be performed via static_cast<T> as well as the
//  // BigEndian<T>::Load() member function.
//  const auto version = static_cast<Version>(pkt.version);
//  const uint32_t sequence_no{pkt.sequence_no.Load()};
//  // For vectors and arrays:
//  std::vector<double> command_copy(9);
//  success =
//      CopyTo(absl::MakeConstSpan(pkt.command), absl::MakeSpan(command_copy));
//
//  // We can also compare values of type T with BigEndian<T>. This will
//  // internally convert the BigEndian<T> to T and use the comparison operator
//  // of T.
//  std::cout << std::boolalpha << (pkt.version == Version::kVersion2);
template <typename T, absl::endian E>
// In theory we could also enable this class template for trivially-copyable
// structs. This could be convenient for cases where a packed struct holds a
// single member and a few functions for manipulating it, e.g. a status with
// different status bits. Bool is explicitly excluded as in C++ it is a 1 byte
// type.
  requires((std::is_arithmetic_v<T> || std::is_enum_v<T>) &&
           std::is_trivially_copyable_v<T> && !std::is_same_v<T, bool> &&
           !std::is_const_v<T> && !std::is_volatile_v<T>)
class ABSL_ATTRIBUTE_PACKED EndianBase {
 public:
  using value_type = T;
  // Unsigned integer with the same size as the wrapped data type.
  using underlying_type = UIntOfSize_t<sizeof(T)>;

  // Converting constructor: We want to allow copy initialization e.g.
  // EndianBase<int> e = 1;
  constexpr EndianBase(const T t = T{}) noexcept {
    static_assert(sizeof(EndianBase) == sizeof(T),
                  "EndianBase must not add padding");
    // Re-use implementation of operator=.
    *this = t;
    return;
  }

  constexpr EndianBase& operator=(const T& other) noexcept {
    Store(other);
    return *this;
  }

  template <absl::endian E1>
  explicit constexpr EndianBase(const EndianBase<T, E1>& other) noexcept {
    // Re-use implementation of operator=.
    *this = other;
    return;
  }

  template <absl::endian E1>
  constexpr EndianBase& operator=(const EndianBase<T, E1>& other) noexcept {
    Store(other.Load());
    return *this;
  }

  // Get endianness of data type.
  static constexpr absl::endian Endian() noexcept { return E; }

  // Explicit conversion operator so that it can only be used with static_cast.
  explicit constexpr operator T() const noexcept { return Load(); }

  // Set the internally stored value to the provided value.
  constexpr void Store(const T val) noexcept {
    // Data types are first cast to an unsigned integer type of same length and
    // then the endianness of the resulting unsigned integer is changed.
    const auto tmp = std::bit_cast<underlying_type>(val);
    raw_value_ = SwapBytesIfNeeded(tmp, E);
    return;
  }

  // Get the internally stored value with native endianness and native data
  // type.
  [[nodiscard]]
  constexpr T Load() const noexcept {
    return std::bit_cast<T>(SwapBytesIfNeeded(raw_value_, E));
  }

  // Get the raw unsigned integer value.
  [[nodiscard]]
  constexpr underlying_type RawValue() const noexcept {
    return raw_value_;
  }

  friend constexpr bool operator==(const EndianBase& lhs,
                                   const T rhs) noexcept {
    return lhs.Load() == rhs;
  }
  friend constexpr auto operator<=>(const EndianBase& lhs,
                                    const T rhs) noexcept {
    return lhs.Load() <=> rhs;
  }

  template <absl::endian E1>
  friend constexpr bool operator==(const EndianBase& lhs,
                                   const EndianBase<T, E1>& rhs) noexcept {
    return lhs.Load() == rhs.Load();
  }
  template <absl::endian E1>
  friend constexpr auto operator<=>(const EndianBase& lhs,
                                    const EndianBase<T, E1>& rhs) noexcept {
    return lhs.Load() <=> rhs.Load();
  }

 private:
  underlying_type raw_value_;
};

template <typename T>
struct is_endian_base : std::false_type {};
template <typename T, absl::endian E>
struct is_endian_base<EndianBase<T, E>> : std::true_type {};

template <typename T>
inline constexpr bool is_endian_base_v = is_endian_base<T>::value;

template <class T>
concept endian_base = is_endian_base_v<T>;

}  // namespace internal

// Copy in between containers of type T and BigEndian<T>. Source and destination
// containers must have the same size.
template <typename T1, typename T2>
  requires((internal::endian_base<T1> || internal::endian_base<T2>) &&
           std::constructible_from<T2, T1>)
[[nodiscard]]
inline constexpr bool CopyTo(absl::Span<const T1> src,
                             absl::Span<T2> dest) noexcept {
  if (src.size() != dest.size()) [[unlikely]] {
    return false;
  }
  std::transform(src.begin(), src.end(), dest.begin(),
                 [](const auto& val) { return static_cast<T2>(val); });
  return true;
}

// Value is stored as little-endian.
template <typename T>
using LittleEndian = internal::EndianBase<T, absl::endian::little>;

// Value is stored as big-endian.
template <typename T>
using BigEndian = internal::EndianBase<T, absl::endian::big>;

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_UTILS_BYTE_ORDER_H_
