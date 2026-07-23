// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/scene/sdf/xml_utils.h"

#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/util/status/status_builder.h"
#include "sdf/Element.hh"
#include "tinyxml2.h"

namespace intrinsic::sdf {
namespace {

struct XmlOptions {
  bool compact = false;
  // If true, collapse whitespace. This is useful for normalizing XML strings
  // for comparison in tests.
  bool collapse_whitespace = false;
};

absl::StatusOr<std::string> FormatXml(absl::string_view xml,
                                      const XmlOptions& options) {
  tinyxml2::XMLDocument doc(
      /*processEntities=*/true,
      options.collapse_whitespace ? tinyxml2::Whitespace::COLLAPSE_WHITESPACE
                                  : tinyxml2::Whitespace::PRESERVE_WHITESPACE);
  if (doc.Parse(xml.data(), xml.size()) != tinyxml2::XML_SUCCESS) {
    return intrinsic::InvalidArgumentErrorBuilder()
           << "Failed to parse XML: " << doc.ErrorStr();
  }
  tinyxml2::XMLPrinter printer(/*file=*/nullptr, options.compact);
  doc.Print(&printer);
  return std::string(printer.CStr());
}

absl::StatusOr<std::string> FormatXml(const ::sdf::ElementConstPtr& element,
                                      const XmlOptions& options) {
  if (element == nullptr) {
    return absl::InvalidArgumentError("Cannot format xml from a nullptr");
  }
  return FormatXml(element->ToString(""), options);
}

}  // namespace

absl::StatusOr<std::string> NormalizeXml(absl::string_view xml) {
  return FormatXml(xml, {.collapse_whitespace = true});
}

absl::StatusOr<std::string> GetCompactXml(
    const ::sdf::ElementConstPtr& element) {
  return FormatXml(element, {.compact = true});
}

absl::StatusOr<std::string> GetPrettifiedXml(
    const ::sdf::ElementConstPtr& element) {
  return FormatXml(element, {.compact = false});
}

std::string EscapeXml(absl::string_view s) {
  std::string result;
  result.reserve(s.size());
  for (char c : s) {
    switch (c) {
      case '&':
        result += "&amp;";
        break;
      case '<':
        result += "&lt;";
        break;
      case '>':
        result += "&gt;";
        break;
      case '\"':
        result += "&quot;";
        break;
      case '\'':
        result += "&apos;";
        break;
      default:
        result += c;
        break;
    }
  }
  return result;
}

}  // namespace intrinsic::sdf
