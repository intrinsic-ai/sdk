// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/proto/parse_text_proto.h"

#include <string>
#include <string_view>

#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "google/protobuf/descriptor.h"
#include "google/protobuf/message.h"
#include "google/protobuf/text_format.h"
#include "intrinsic/util/proto/error_collector.h"

namespace intrinsic {
namespace internal {

namespace {

// Exactly the same as the default Finder implementation except that it does not
// error out if the type URL prefix of an Any proto is not one of
// 'type.googleapis.com' or 'type.googleprod.com'. This enables parsing text
// protos with Anys having Intrinsic-style type URL prefixes.
class Finder : public google::protobuf::TextFormat::Finder {
  const google::protobuf::Descriptor* FindAnyType(
      const google::protobuf::Message& message,
      [[maybe_unused]] const std::string& prefix,
      const std::string& name) const override {
    return message.GetDescriptor()->file()->pool()->FindMessageTypeByName(name);
  }
};

}  // namespace

absl::Status ParseTextProtoImpl(std::string_view asciipb,
                                google::protobuf::Message& message) {
  google::protobuf::TextFormat::Parser parser;
  Finder finder;
  parser.SetFinder(&finder);
  SimpleErrorCollector error_collector;
  parser.RecordErrorsTo(&error_collector);

  if (!parser.ParseFromString(asciipb, &message)) {
    return absl::InvalidArgumentError(absl::StrCat(
        "Cannot parse protobuf ", message.GetDescriptor()->full_name(),
        " from text: ", error_collector.str()));
  }
  return absl::OkStatus();
}

}  // namespace internal
}  // namespace intrinsic
