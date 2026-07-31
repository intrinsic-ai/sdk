// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_MARSHAL_RIEGELI_CODER_H_
#define INTRINSIC_MARSHAL_RIEGELI_CODER_H_

#include <cstddef>
#include <cstdint>
#include <string>
#include <utility>
#include <vector>

#include "absl/container/flat_hash_map.h"
#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/util/demangle.h"
#include "intrinsic/util/status/status.pb.h"
#include "intrinsic/util/status/status_conversion_proto.h"
#include "intrinsic/util/status/status_macros.h"
#include "riegeli/records/record_reader.h"
#include "riegeli/records/record_writer.h"

namespace intrinsic {

// Base definition of a Riegeli Coder. It uses riegeli::RecordReaderBase and
// Writer to serialize the given type.
//
// The default Encode and Decode implementations are a no-op and should be
// overridden by specializations.
// The default TypeName implementation uses Demangle on the type.
template <typename T>
struct RiegeliCoder {
  static std::string TypeName() { return Demangle<T>(); }
  static absl::Status Encode(const T& value,
                             riegeli::RecordWriterBase& writer) = delete;
  static absl::StatusOr<T> Decode(riegeli::RecordReaderBase& reader) = delete;
};

template <>
absl::Status RiegeliCoder<std::string>::Encode(
    const std::string& value, riegeli::RecordWriterBase& writer);

template <>
absl::StatusOr<std::string> RiegeliCoder<std::string>::Decode(
    riegeli::RecordReaderBase& reader);

template <>
absl::Status RiegeliCoder<int>::Encode(const int& value,
                                       riegeli::RecordWriterBase& writer);

template <>
absl::StatusOr<int> RiegeliCoder<int>::Decode(
    riegeli::RecordReaderBase& reader);

template <>
absl::Status RiegeliCoder<size_t>::Encode(const size_t& value,
                                          riegeli::RecordWriterBase& writer);

template <>
absl::StatusOr<size_t> RiegeliCoder<size_t>::Decode(
    riegeli::RecordReaderBase& reader);

template <>
absl::Status RiegeliCoder<double>::Encode(const double& value,
                                          riegeli::RecordWriterBase& writer);

template <>
absl::StatusOr<double> RiegeliCoder<double>::Decode(
    riegeli::RecordReaderBase& reader);

template <>
absl::Status RiegeliCoder<float>::Encode(const float& value,
                                         riegeli::RecordWriterBase& writer);

template <>
absl::StatusOr<float> RiegeliCoder<float>::Decode(
    riegeli::RecordReaderBase& reader);

template <>
absl::Status RiegeliCoder<bool>::Encode(const bool& value,
                                        riegeli::RecordWriterBase& writer);

template <>
absl::StatusOr<bool> RiegeliCoder<bool>::Decode(
    riegeli::RecordReaderBase& reader);

template <typename T>
struct RiegeliCoder<std::vector<T>> {
  static std::string TypeName() { return Demangle<std::vector<T>>(); }
  static absl::Status Encode(const std::vector<T>& value,
                             riegeli::RecordWriterBase& writer) {
    INTR_RETURN_IF_ERROR(RiegeliCoder<size_t>::Encode(value.size(), writer));
    for (const T& element : value) {
      INTR_RETURN_IF_ERROR(RiegeliCoder<T>::Encode(element, writer));
    }
    return absl::OkStatus();
  }
  static absl::StatusOr<std::vector<T>> Decode(
      riegeli::RecordReaderBase& reader) {
    INTR_ASSIGN_OR_RETURN(size_t size, RiegeliCoder<size_t>::Decode(reader));
    std::vector<T> result;
    result.reserve(size);
    for (size_t i = 0; i < size; ++i) {
      INTR_ASSIGN_OR_RETURN(T element, RiegeliCoder<T>::Decode(reader));
      result.push_back(element);
    }
    return result;
  }
};

template <typename T1, typename T2>
struct RiegeliCoder<std::pair<T1, T2>> {
  static std::string TypeName() { return Demangle<std::pair<T1, T2>>(); }
  static absl::Status Encode(const std::pair<T1, T2>& value,
                             riegeli::RecordWriterBase& writer) {
    INTR_RETURN_IF_ERROR(RiegeliCoder<T1>::Encode(value.first, writer));
    INTR_RETURN_IF_ERROR(RiegeliCoder<T2>::Encode(value.second, writer));
    return absl::OkStatus();
  }
  static absl::StatusOr<std::pair<T1, T2>> Decode(
      riegeli::RecordReaderBase& reader) {
    INTR_ASSIGN_OR_RETURN(T1 first, RiegeliCoder<T1>::Decode(reader));
    INTR_ASSIGN_OR_RETURN(T2 second, RiegeliCoder<T2>::Decode(reader));
    return std::make_pair(std::move(first), std::move(second));
  }
};

template <typename T1, typename T2>
struct RiegeliCoder<absl::flat_hash_map<T1, T2>> {
  static std::string TypeName() {
    return Demangle<absl::flat_hash_map<T1, T2>>();
  }
  static absl::Status Encode(const absl::flat_hash_map<T1, T2>& value,
                             riegeli::RecordWriterBase& writer) {
    INTR_RETURN_IF_ERROR(RiegeliCoder<size_t>::Encode(value.size(), writer));
    std::vector<T1> keys;
    keys.reserve(value.size());
    for (const auto& [k, v] : value) {
      keys.push_back(k);
    }
    std::sort(keys.begin(), keys.end());
    for (const T1& key : keys) {
      INTR_RETURN_IF_ERROR(RiegeliCoder<T1>::Encode(key, writer));
      INTR_RETURN_IF_ERROR(RiegeliCoder<T2>::Encode(value.at(key), writer));
    }
    return absl::OkStatus();
  }
  static absl::StatusOr<absl::flat_hash_map<T1, T2>> Decode(
      riegeli::RecordReaderBase& reader) {
    INTR_ASSIGN_OR_RETURN(size_t size, RiegeliCoder<size_t>::Decode(reader));
    absl::flat_hash_map<T1, T2> result;
    result.reserve(size);
    for (size_t i = 0; i < size; ++i) {
      INTR_ASSIGN_OR_RETURN(T1 key, RiegeliCoder<T1>::Decode(reader));
      INTR_ASSIGN_OR_RETURN(T2 value, RiegeliCoder<T2>::Decode(reader));
      result.emplace(std::move(key), std::move(value));
    }
    return result;
  }
};

template <typename T>
struct RiegeliCoder<absl::StatusOr<T>> {
  static std::string TypeName() { return Demangle<absl::StatusOr<T>>(); }
  static absl::Status Encode(const absl::StatusOr<T>& value,
                             riegeli::RecordWriterBase& writer) {
    INTR_RETURN_IF_ERROR(RiegeliCoder<bool>::Encode(value.ok(), writer));
    if (value.ok()) {
      INTR_RETURN_IF_ERROR(RiegeliCoder<T>::Encode(value.value(), writer));
    } else {
      intrinsic_proto::StatusProto proto;
      intrinsic::SaveStatusToProto(value.status(), &proto);
      std::string serialized;
      CHECK(proto.SerializeToString(&serialized));
      return intrinsic::RiegeliCoder<std::string>::Encode(serialized, writer);
    }
    return absl::OkStatus();
  }
  static absl::StatusOr<absl::StatusOr<T>> Decode(
      riegeli::RecordReaderBase& reader) {
    INTR_ASSIGN_OR_RETURN(bool ok, RiegeliCoder<bool>::Decode(reader));
    if (ok) {
      INTR_ASSIGN_OR_RETURN(T value, RiegeliCoder<T>::Decode(reader));
      return value;
    } else {
      INTR_ASSIGN_OR_RETURN(std::string serialized,
                            RiegeliCoder<std::string>::Decode(reader));
      intrinsic_proto::StatusProto proto;
      CHECK(proto.ParseFromString(serialized));
      return intrinsic::MakeStatusFromProto(proto);
    }
  }
};

}  // namespace intrinsic

#endif  // INTRINSIC_MARSHAL_RIEGELI_CODER_H_
