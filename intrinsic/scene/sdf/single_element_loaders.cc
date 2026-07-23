// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/scene/sdf/single_element_loaders.h"

#include <memory>
#include <string>

#include "absl/log/check.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_join.h"
#include "absl/strings/string_view.h"
#include "absl/strings/substitute.h"
#include "intrinsic/util/status/status_macros.h"
#include "sdf/Collision.hh"
#include "sdf/Element.hh"
#include "sdf/Link.hh"
#include "sdf/Model.hh"
#include "sdf/Root.hh"
#include "sdf/Types.hh"
#include "sdf/Visual.hh"
#include "sdf/parser.hh"

namespace intrinsic {
namespace sdf {

// Converts a single element xml string into a sdf element pointer.
absl::StatusOr<::sdf::ElementPtr> SdfElementFromString(
    absl::string_view element_string, absl::string_view schema_file) {
  auto element = std::make_shared<::sdf::Element>();
  if (schema_file.empty()) {
    LOG(WARNING) << "SDF schema file is empty, please provide a valid schema "
                    "file for element validation. Fallback to root.sdf.";
  }
  const std::string schema_file_with_fallback =
      schema_file.empty() ? "root.sdf" : std::string(schema_file);
  if (!::sdf::initFile(schema_file_with_fallback, element)) {
    return absl::InvalidArgumentError(absl::Substitute(
        "Failed to load SDF schema file '$0'", schema_file_with_fallback));
  }

  std::string sdf_string =
      absl::StrCat("<sdf version=\"1.10\">", element_string, "</sdf>");

  ::sdf::Errors errors;
  if (!::sdf::readString(sdf_string, element, errors)) {
    if (!errors.empty()) {
      return absl::InvalidArgumentError(absl::Substitute(
          "Failed to parse SDF string $0 with $1 errors. Errors are: $2",
          sdf_string, errors.size(),
          absl::StrJoin(errors, "\n", absl::StreamFormatter())));
    } else {
      return absl::InvalidArgumentError(absl::Substitute(
          "Failed to parse SDF string $0 with 0 errors. ", element_string));
    }
  }

  // When schema_file is empty, the element is parsed with the root.sdf schema,
  // but we want the actual child element parsed.
  return schema_file.empty() ? element->GetFirstElement() : element;
}

template <>
absl::StatusOr<::sdf::Link> TypedSdfElementFromString<::sdf::Link>(
    absl::string_view element_string) {
  // Link needs a wrapper model for link so that semantic pose is evaluated
  // properly.
  std::string sdf_string = absl::Substitute(R"(
    <sdf version="1.10" xmlns:intrinsic="https://intrinsic.ai/">
      <model name="wrapper_model">
        $0
      </model>
    </sdf>)",
                                            element_string);
  auto root = std::make_unique<::sdf::Root>();
  auto errors = root->LoadSdfString(sdf_string);
  if (!errors.empty()) {
    return absl::InvalidArgumentError(absl::Substitute(
        "Failed to parse SDF string with $0 errors. Errors are: $1",
        errors.size(), absl::StrJoin(errors, "\n", absl::StreamFormatter())));
  }
  if (!root->Model()) {
    return absl::InvalidArgumentError(absl::Substitute(
        "Failed to parse SDF string \n'$0'\n as a single link.",
        element_string));
  }
  const auto* model = root->Model();

  if (model->LinkCount() != 1) {
    return absl::InvalidArgumentError(absl::Substitute(
        "Failed to parse SDF string \n'$0'\n as a single link.",
        element_string));
  }

  return *model->LinkByIndex(0);
}

template <>
absl::StatusOr<::sdf::Visual> TypedSdfElementFromString<::sdf::Visual>(
    absl::string_view element_string) {
  // Visual need a wrapper model and link so that semantic pose is evaluated
  // properly.
  std::string link_string =
      absl::StrCat("<link name=\"wrapper_link\">", element_string, "</link>");
  INTR_ASSIGN_OR_RETURN(auto link,
                        TypedSdfElementFromString<::sdf::Link>(link_string));
  if (link.VisualCount() != 1) {
    return absl::InvalidArgumentError(absl::Substitute(
        "Failed to parse SDF string \n'$0'\n as a single visual.",
        element_string));
  }

  return *link.VisualByIndex(0);
}

template <>
absl::StatusOr<::sdf::Collision> TypedSdfElementFromString<::sdf::Collision>(
    absl::string_view element_string) {
  // Collision need a wrapper model and link so that semantic pose is evaluated
  // properly.
  std::string link_string =
      absl::StrCat("<link name=\"wrapper_link\">", element_string, "</link>");
  INTR_ASSIGN_OR_RETURN(auto link,
                        TypedSdfElementFromString<::sdf::Link>(link_string));
  if (link.CollisionCount() != 1) {
    return absl::InvalidArgumentError(absl::Substitute(
        "Failed to parse SDF string \n'$0'\n as a single collision.",
        element_string));
  }

  return *link.CollisionByIndex(0);
}

}  // namespace sdf
}  // namespace intrinsic
