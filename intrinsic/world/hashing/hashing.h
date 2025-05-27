// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_WORLD_HASHING_HASHING_H_
#define INTRINSIC_WORLD_HASHING_HASHING_H_

#include <algorithm>
#include <cstddef>
#include <functional>
#include <string>
#include <string_view>
#include <tuple>
#include <type_traits>
#include <utility>
#include <vector>

#include "absl/container/flat_hash_map.h"
#include "absl/container/flat_hash_set.h"
#include "intrinsic/util/string_type.h"
#include "intrinsic/world/entity_id.h"

namespace intrinsic {

namespace internal {
template <typename T>
auto hash_world_object(const T& s) {
  using ValueType = std::decay_t<decltype(s.value())>;
  return std::hash<ValueType>{}(s.value());
}

// A simple hash combine function, used to combine multiple hashes into one.
inline std::size_t hash_combine(std::size_t a, std::size_t b) {
  return a ^ (b + 0x9e3779b9 + (a << 6) + (a >> 2));
}
}  // namespace internal

// This file contains a collection of hash functions for the various types
// used in the world library.
// The goal is to provide a workaround for b/380032603.

template <typename T>
struct WorldHasher {
  std::size_t operator()(const T& s) const noexcept {
    return std::hash<T>{}(s);
  }
};

template <>
struct WorldHasher<std::string> {
  using is_transparent = void;  // Enables heterogeneous lookup

  std::size_t operator()(const std::string& s) const noexcept {
    return std::hash<std::string>{}(s);
  }
  std::size_t operator()(std::string_view s) const noexcept {
    return std::hash<std::string_view>{}(s);
  }
  size_t operator()(const char* cstr) const {
    return std::hash<std::string_view>{}(cstr);
  }
};

template <typename T>
struct WorldHasher<StringType<T>> {
  std::size_t operator()(const StringType<T>& s) const noexcept {
    return internal::hash_world_object(s);
  }
};

template <>
struct WorldHasher<TypedEntityId<>> {
  std::size_t operator()(const TypedEntityId<>& s) const noexcept {
    return internal::hash_world_object(s);
  }
};

template <typename... Types>
struct WorldHasher<TypedEntityId<Types...>> {
  std::size_t operator()(const TypedEntityId<Types...>& s) const noexcept {
    return internal::hash_world_object(s);
  }
};

template <typename T, typename U>
struct WorldHasher<std::pair<T, U>> {
  std::size_t operator()(const std::pair<T, U>& s) const noexcept {
    return internal::hash_combine(WorldHasher<T>{}(s.first),
                                  WorldHasher<U>{}(s.second));
  }
};

template <typename X, typename Y>
struct WorldHasher<std::tuple<X, Y>> {
  std::size_t operator()(const std::tuple<X, Y>& s) const noexcept {
    return internal::hash_combine(WorldHasher<X>{}(std::get<0>(s)),
                                  WorldHasher<Y>{}(std::get<1>(s)));
  }
};

template <>
struct WorldHasher<EntityId> {
  std::size_t operator()(const EntityId& s) const noexcept {
    return internal::hash_world_object(s);
  }
};

// Hash function for absl::flat_hash_set. This function sorts the elements of
// the set before hashing them. This is done to ensure that the hash of a set is
// independent of the order of the elements.
template <typename T, typename H>
struct WorldHasher<absl::flat_hash_set<T, H>> {
  std::size_t operator()(const absl::flat_hash_set<T, H>& s) const noexcept {
    if (s.empty()) {
      return 0;
    }
    std::vector<size_t> hashes;
    for (const auto& t : s) {
      hashes.push_back(H{}(t));
    }
    std::sort(hashes.begin(), hashes.end());
    auto hash = hashes[0];
    for (size_t i = 1; i < hashes.size(); ++i) {
      hash = internal::hash_combine(hash, hashes[i]);
    }
    return hash;
  }
};

// Forward declaration of the WorldHasher template to prevent circular
// dependencies.
template <typename T>
struct WorldHasher;

template <typename T>
using WorldHashSet = absl::flat_hash_set<T, WorldHasher<T>>;

template <typename K, typename V>
using WorldHashMap = absl::flat_hash_map<K, V, WorldHasher<K>>;

}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_HASHING_HASHING_H_
