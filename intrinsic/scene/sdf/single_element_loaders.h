// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SCENE_SDF_SINGLE_ELEMENT_LOADERS_H_
#define INTRINSIC_SCENE_SDF_SINGLE_ELEMENT_LOADERS_H_

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_join.h"
#include "absl/strings/string_view.h"
#include "absl/strings/substitute.h"
#include "intrinsic/util/status/status_macros.h"
#include "sdf/Element.hh"
#include "sdf/Geometry.hh"
#include "sdf/Link.hh"
#include "sdf/Sensor.hh"
#include "sdf/Surface.hh"
#include "sdf/Types.hh"
#include "sdf/Visual.hh"

namespace intrinsic {
namespace sdf {

// Converts a single element xml string into a sdf element pointer.
absl::StatusOr<::sdf::ElementPtr> SdfElementFromString(
    absl::string_view element_string, absl::string_view schema_file);

template <typename T>
struct SdfTypeToSchemaFile {
  static constexpr absl::string_view kSchemaFile;
};
template <>
constexpr absl::string_view SdfTypeToSchemaFile<::sdf::Contact>::kSchemaFile =
    "contact.sdf";
template <>
constexpr absl::string_view SdfTypeToSchemaFile<::sdf::Geometry>::kSchemaFile =
    "geometry.sdf";
template <>
constexpr absl::string_view SdfTypeToSchemaFile<::sdf::Sensor>::kSchemaFile =
    "sensor.sdf";
template <>
constexpr absl::string_view SdfTypeToSchemaFile<::sdf::Surface>::kSchemaFile =
    "surface.sdf";

// Converts a single element xml string into a sdf data type `T` corresponding
// to that element.
template <typename T>
absl::StatusOr<T> TypedSdfElementFromString(absl::string_view element_string) {
  INTR_ASSIGN_OR_RETURN(
      auto element, SdfElementFromString(element_string,
                                         SdfTypeToSchemaFile<T>::kSchemaFile));

  T typed_element;
  if (::sdf::Errors load_errors = typed_element.Load(element);
      !load_errors.empty()) {
    return absl::InvalidArgumentError(absl::Substitute(
        "Failed to load SDF element to typed data with $0 errors. Errors are: "
        "\n$1"
        "\nOriginal string is:"
        "\n$2",
        load_errors.size(),
        absl::StrJoin(load_errors, "\n", absl::StreamFormatter()),
        element_string));
  }
  return typed_element;
}

template <>
absl::StatusOr<::sdf::Link> TypedSdfElementFromString<::sdf::Link>(
    absl::string_view element_string);

template <>
absl::StatusOr<::sdf::Visual> TypedSdfElementFromString<::sdf::Visual>(
    absl::string_view element_string);

template <>
absl::StatusOr<::sdf::Collision> TypedSdfElementFromString<::sdf::Collision>(
    absl::string_view element_string);

}  // namespace sdf
}  // namespace intrinsic

#endif  // INTRINSIC_SCENE_SDF_SINGLE_ELEMENT_LOADERS_H_
