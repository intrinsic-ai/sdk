// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/actions/action_utils.h"

#include <string>
#include <utility>

#include "absl/container/flat_hash_set.h"
#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_join.h"
#include "absl/strings/string_view.h"
#include "google/protobuf/descriptor.pb.h"
#include "intrinsic/icon/proto/v1/types.pb.h"
#include "intrinsic/icon/release/source_location.h"

namespace intrinsic {
namespace icon {

namespace {
struct FeatureInterfaceNameFormatter {
  void operator()(
      std::string* out,
      const intrinsic_proto::icon::v1::FeatureInterfaceTypes& t) const {
    absl::StrAppend(out,
                    intrinsic_proto::icon::v1::FeatureInterfaceTypes_Name(t));
  }
};
}  // namespace

absl::Status ActionSignatureBuilder::SetFixedParametersTypeImpl(
    absl::string_view fixed_parameters_message_type,
    const google::protobuf::FileDescriptorSet& fixed_parameters_descriptor_set,
    intrinsic::SourceLocation loc,
    intrinsic_proto::icon::v1::ActionSignature& dest_signature) {
  if (!dest_signature.fixed_parameters_message_type().empty()) {
    return absl::AlreadyExistsError(
        absl::StrCat(loc.file_name(), ":", loc.line(),
                     " Fixed parameters type already set to \"",
                     dest_signature.fixed_parameters_message_type(), "\""));
  }
  dest_signature.set_fixed_parameters_message_type(
      fixed_parameters_message_type);
  *dest_signature.mutable_fixed_parameters_descriptor_set() =
      fixed_parameters_descriptor_set;
  return absl::OkStatus();
}

absl::Status ActionSignatureBuilder::AddPartSlot(
    absl::string_view slot_name, absl::string_view slot_description,
    absl::flat_hash_set<intrinsic_proto::icon::v1::FeatureInterfaceTypes>
        required_feature_interfaces,
    absl::flat_hash_set<intrinsic_proto::icon::v1::FeatureInterfaceTypes>
        optional_feature_interfaces,
    intrinsic::SourceLocation loc) {
  if (bool inserted = part_slot_names_.emplace(slot_name).second; !inserted) {
    return absl::AlreadyExistsError(
        absl::StrCat(loc.file_name(), ":", loc.line(),
                     " Duplicate Part Slot name \"", slot_name, "\""));
  }
  absl::flat_hash_set<intrinsic_proto::icon::v1::FeatureInterfaceTypes>
      feature_interfaces_intersection;
  for (const auto& interface : required_feature_interfaces) {
    if (optional_feature_interfaces.contains(interface)) {
      feature_interfaces_intersection.insert(interface);
    }
  }
  if (!feature_interfaces_intersection.empty()) {
    return absl::InvalidArgumentError(absl::StrCat(
        "The following Feature interfaces were listed as both required and "
        "optional for Slot '",
        slot_name,
        "', please ensure each Feature Interface only appears once: [",
        absl::StrJoin(feature_interfaces_intersection, ", ",
                      FeatureInterfaceNameFormatter())));
  }
  intrinsic_proto::icon::v1::ActionSignature::PartSlotInfo info;
  info.set_description(std::string(slot_description));
  *(info.mutable_required_feature_interfaces()) = {
      required_feature_interfaces.begin(), required_feature_interfaces.end()};
  *(info.mutable_optional_feature_interfaces()) = {
      optional_feature_interfaces.begin(), optional_feature_interfaces.end()};
  signature_.mutable_part_slot_infos()->insert({std::string(slot_name), info});

  return absl::OkStatus();
}

absl::Status ActionSignatureBuilder::AddRealtimeSignal(
    absl::string_view signal_name, absl::string_view signal_description,
    intrinsic::SourceLocation loc) {
  if (realtime_signal_names_.contains(signal_name)) {
    return absl::AlreadyExistsError(
        absl::StrCat(loc.file_name(), ":", loc.line(),
                     " Duplicate Realtime Signal name \"", signal_name, "\""));
  }
  intrinsic_proto::icon::v1::ActionSignature::RealtimeSignalInfo* signal_info =
      signature_.add_realtime_signal_infos();
  signal_info->set_signal_name(signal_name);
  signal_info->set_text_description(std::string(signal_description));
  realtime_signal_names_.insert(std::string(signal_name));
  return absl::OkStatus();
}

}  // namespace icon
}  // namespace intrinsic
