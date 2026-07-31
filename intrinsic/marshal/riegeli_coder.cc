// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/marshal/riegeli_coder.h"

#include <cstddef>
#include <string>

#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/util/status/status_macros.h"
#include "riegeli/bytes/string_reader.h"
#include "riegeli/bytes/string_writer.h"
#include "riegeli/endian/endian_reading.h"
#include "riegeli/endian/endian_writing.h"
#include "riegeli/records/record_reader.h"
#include "riegeli/records/record_writer.h"

namespace intrinsic {

namespace {
template <typename T>
absl::StatusOr<std::string> EncodeLittleEndianToRecord(const T& value) {
  std::string result;
  riegeli::StringWriter writer(&result);
  if (!riegeli::WriteLittleEndian<T>(value, writer)) {
    return absl::InternalError(
        absl::StrCat("Failed to write ", RiegeliCoder<T>::TypeName()));
  }
  if (!writer.Close()) {
    return writer.status();
  }
  return result;
}

template <typename T>
absl::StatusOr<T> DecodeLittleEndianFromRecord(absl::string_view record) {
  riegeli::StringReader reader((std::string(record)));
  T value;
  if (!riegeli::ReadLittleEndian<T>(reader, value)) {
    return absl::InternalError(
        absl::StrCat("Failed to read ", RiegeliCoder<T>::TypeName()));
  }
  if (!reader.Close()) {
    return reader.status();
  }
  return value;
}

}  // namespace

template <>
absl::Status RiegeliCoder<std::string>::Encode(
    const std::string& value, riegeli::RecordWriterBase& writer) {
  if (!writer.WriteRecord(value)) {
    return writer.status();
  }
  return absl::OkStatus();
}

template <>
absl::StatusOr<std::string> RiegeliCoder<std::string>::Decode(
    riegeli::RecordReaderBase& reader) {
  absl::string_view value;
  if (!reader.ReadRecord(value)) {
    return reader.status();
  }
  return std::string(value);
}

template <>
absl::Status RiegeliCoder<int>::Encode(const int& value,
                                       riegeli::RecordWriterBase& writer) {
  INTR_ASSIGN_OR_RETURN(auto record, EncodeLittleEndianToRecord<int>(value));
  return RiegeliCoder<std::string>::Encode(record, writer);
}

template <>
absl::StatusOr<int> RiegeliCoder<int>::Decode(
    riegeli::RecordReaderBase& reader) {
  INTR_ASSIGN_OR_RETURN(auto record, RiegeliCoder<std::string>::Decode(reader));
  INTR_ASSIGN_OR_RETURN(auto result, DecodeLittleEndianFromRecord<int>(record));
  return result;
}

template <>
absl::Status RiegeliCoder<size_t>::Encode(const size_t& value,
                                          riegeli::RecordWriterBase& writer) {
  INTR_ASSIGN_OR_RETURN(auto record, EncodeLittleEndianToRecord<size_t>(value));
  return RiegeliCoder<std::string>::Encode(record, writer);
}

template <>
absl::StatusOr<size_t> RiegeliCoder<size_t>::Decode(
    riegeli::RecordReaderBase& reader) {
  INTR_ASSIGN_OR_RETURN(auto record, RiegeliCoder<std::string>::Decode(reader));
  INTR_ASSIGN_OR_RETURN(auto result,
                        DecodeLittleEndianFromRecord<size_t>(record));
  return result;
}

template <>
absl::Status RiegeliCoder<double>::Encode(const double& value,
                                          riegeli::RecordWriterBase& writer) {
  INTR_ASSIGN_OR_RETURN(auto record, EncodeLittleEndianToRecord<double>(value));
  return RiegeliCoder<std::string>::Encode(record, writer);
}

template <>
absl::StatusOr<double> RiegeliCoder<double>::Decode(
    riegeli::RecordReaderBase& reader) {
  INTR_ASSIGN_OR_RETURN(auto record, RiegeliCoder<std::string>::Decode(reader));
  INTR_ASSIGN_OR_RETURN(auto result,
                        DecodeLittleEndianFromRecord<double>(record));
  return result;
}

template <>
absl::Status RiegeliCoder<float>::Encode(const float& value,
                                         riegeli::RecordWriterBase& writer) {
  INTR_ASSIGN_OR_RETURN(auto record, EncodeLittleEndianToRecord<float>(value));
  return RiegeliCoder<std::string>::Encode(record, writer);
}

template <>
absl::StatusOr<float> RiegeliCoder<float>::Decode(
    riegeli::RecordReaderBase& reader) {
  INTR_ASSIGN_OR_RETURN(auto record, RiegeliCoder<std::string>::Decode(reader));
  INTR_ASSIGN_OR_RETURN(auto result,
                        DecodeLittleEndianFromRecord<float>(record));
  return result;
}

template <>
absl::Status RiegeliCoder<bool>::Encode(const bool& value,
                                        riegeli::RecordWriterBase& writer) {
  INTR_RETURN_IF_ERROR(RiegeliCoder<int>::Encode(value ? 1 : 0, writer));
  return absl::OkStatus();
}

template <>
absl::StatusOr<bool> RiegeliCoder<bool>::Decode(
    riegeli::RecordReaderBase& reader) {
  INTR_ASSIGN_OR_RETURN(auto value, RiegeliCoder<int>::Decode(reader));
  return value == 1;
}

}  // namespace intrinsic
