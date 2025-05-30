// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/common/part_properties.h"

#include <variant>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "intrinsic/icon/proto/v1/service.pb.h"

namespace intrinsic::icon {
namespace {
// Helper to convert from std::variant<bool, double> to a PartPropertyValue
// proto with std::visit
struct VariantToPartPropertyProto {
  ::intrinsic_proto::icon::v1::PartPropertyValue operator()(bool value) {
    ::intrinsic_proto::icon::v1::PartPropertyValue proto;
    proto.set_bool_value(value);
    return proto;
  }

  ::intrinsic_proto::icon::v1::PartPropertyValue operator()(double value) {
    ::intrinsic_proto::icon::v1::PartPropertyValue proto;
    proto.set_double_value(value);
    return proto;
  }
};
}  // namespace

::intrinsic_proto::icon::v1::PartPropertyValue ToProto(
    const PartPropertyValue& value) {
  return std::visit(VariantToPartPropertyProto(), value);
}

absl::StatusOr<PartPropertyValue> FromProto(
    const ::intrinsic_proto::icon::v1::PartPropertyValue& value) {
  switch (value.value_case()) {
    case ::intrinsic_proto::icon::v1::PartPropertyValue::ValueCase::kBoolValue:
      return PartPropertyValue(value.bool_value());
    case ::intrinsic_proto::icon::v1::PartPropertyValue::ValueCase::
        kDoubleValue:
      return PartPropertyValue(value.double_value());
    case ::intrinsic_proto::icon::v1::PartPropertyValue::ValueCase::
        VALUE_NOT_SET:
      return absl::InvalidArgumentError("Part property has no value set");
    default:
      return absl::InvalidArgumentError(
          "Part property has unknown value type – this could be due to "
          "version skew between ICON client and server");
  }
}

absl::Status AssignPropertyValue::operator()(bool src, bool& dst) {
  dst = src;
  return absl::OkStatus();
}

absl::Status AssignPropertyValue::operator()(double src, double& dst) {
  dst = src;
  return absl::OkStatus();
}

absl::Status AssignPropertyValue::operator()(double src, bool& dst) {
  return absl::InvalidArgumentError(absl::StrCat(
      "Cannot assign double value to boolean property '", property_name, "'"));
}

absl::Status AssignPropertyValue::operator()(bool src, double& dst) {
  return absl::InvalidArgumentError(absl::StrCat(
      "Cannot assign boolean value to double property '", property_name, "'"));
}

}  // namespace intrinsic::icon
