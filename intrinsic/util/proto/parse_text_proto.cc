// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/proto/parse_text_proto.h"

#include <string>
#include <string_view>

#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_replace.h"
#include "absl/strings/string_view.h"
#include "google/protobuf/descriptor.h"
#include "google/protobuf/message.h"
#include "google/protobuf/text_format.h"
#include "intrinsic/util/proto/error_collector.h"
#include "re2/re2.h"

namespace intrinsic {
namespace {

constexpr absl::string_view kSlashReplacement = "_SLSH_";
constexpr absl::string_view kDotReplacement = "_DOT_";
constexpr absl::string_view kHyphenReplacement = "_HYPH_";
constexpr absl::string_view kPlusReplacement = "_PLUS_";

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

// Rewrites all Any type URLs in the given text proto so that they don't get
// rejected by TextFormat::Parser. This workaround is currently necessary to
// support Intrinsic-style type URLs. The performed replacements are easy to
// revert so that a custom TextFormat::Finder can reconstruct the original type
// URLs from the rewritten ones.
//
// Example:
//     value: {
//         [type.foo.com/bar/0.1/intrinsic_proto.Pose3d]
//         ...
//     }
//   becomes:
//     value: {
//         [type.foo.com_SLASH_bar_SLASH_0_DOT_1/intrinsic_proto.Pose3d]
//         ...
//     }
std::string RewriteAnyTypeUrls(std::string_view asciipb) {
  std::string asciipb_rewritten;
  asciipb_rewritten.reserve(asciipb.size());

  std::string_view unconsumed(asciipb);
  std::string_view consumed_text, consumed_type_url_prefix, consumed_type_name;

  // Expanded Anys in text protos can be unambiguously identified by a
  // "[<url>]"-expression following a "{" where <url> contains at least one "/".
  // See go/textformat-spec#any.
  static constexpr LazyRE2 kAnyTypeUrlRegex = {
      R"re((?s)(.*?\{\s*\[)([^\]]*)(/[^\]/]*)\])re"};

  // Repeatedly consume text up to and including the next Any type URL, perform
  // replacements on the type URL and append the consumed and replaced text to
  // the result.
  while (RE2::Consume(&unconsumed, *kAnyTypeUrlRegex, &consumed_text,
                      &consumed_type_url_prefix, &consumed_type_name)) {
    // After each match, the capture variables hold, e.g., the  following:
    //   consumed_text:              "... { ["
    //   consumed_type_url_prefix:   "type.foo.com/bar"
    //   consumed_type_name:         "/intrinsic_proto.Pose3d"
    //   <not captured in variable>: "]"

    // Some characters in the type URL *prefix* are generally not supported -
    // replace all of them.
    std::string replaced_type_url_prefix = absl::StrReplaceAll(
        consumed_type_url_prefix, {{"/", kSlashReplacement},
                                   {"+", kPlusReplacement},
                                   {"-", kHyphenReplacement}});

    // Dots may not be followed by a digit. Replace, e.g., ".3" -> "_DOT_3".
    static constexpr LazyRE2 kDotFollowedByDigitRegex = {R"re(\.(\d))re"};
    RE2::GlobalReplace(&replaced_type_url_prefix, *kDotFollowedByDigitRegex,
                       absl::StrCat(kDotReplacement, "\\1"));

    absl::StrAppend(&asciipb_rewritten, consumed_text, replaced_type_url_prefix,
                    consumed_type_name, "]");
  }

  // Append the remaining text (everything after the last Any type URL or
  // simply everything if there are no Any type URLs).
  asciipb_rewritten.reserve(asciipb_rewritten.size() + unconsumed.size());
  absl::StrAppend(&asciipb_rewritten, unconsumed);

  return asciipb_rewritten;
}

}  // namespace

absl::Status ParseTextProtoInto(std::string_view asciipb,
                                google::protobuf::Message* message) {
  google::protobuf::TextFormat::Parser parser;
  Finder finder;
  parser.SetFinder(&finder);
  SimpleErrorCollector error_collector;
  parser.RecordErrorsTo(&error_collector);

  std::string asciipb_rewritten = RewriteAnyTypeUrls(asciipb);

  if (!parser.ParseFromString(asciipb_rewritten, message)) {
    return absl::InvalidArgumentError(absl::StrCat(
        "Cannot parse protobuf ", message->GetDescriptor()->full_name(),
        " from text: ", error_collector.str()));
  }
  return absl::OkStatus();
}

}  // namespace intrinsic
