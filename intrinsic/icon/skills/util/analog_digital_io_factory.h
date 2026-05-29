// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_SKILLS_UTIL_ANALOG_DIGITAL_IO_FACTORY_H_
#define INTRINSIC_ICON_SKILLS_UTIL_ANALOG_DIGITAL_IO_FACTORY_H_
#include <memory>
#include <string>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/equipment/channel_factory.h"
#include "intrinsic/icon/skills/util/analog_digital_io.h"
#include "intrinsic/resources/proto/resource_handle.pb.h"
#include "intrinsic/skills/cc/equipment_pack.h"
#include "intrinsic/skills/proto/equipment.pb.h"

namespace intrinsic::skills {
absl::StatusOr<std::unique_ptr<AnalogDigitalIOInterface>> CreateAnalogDigitalIO(
    absl::string_view equipment_slot,
    const google::protobuf::Map<std::string,
                                intrinsic_proto::resources::ResourceHandle>&
        resource_handles,
    const icon::ChannelFactory* channel_factory = nullptr);

absl::StatusOr<std::unique_ptr<AnalogDigitalIOInterface>> CreateAnalogDigitalIO(
    absl::string_view equipment_slot, const EquipmentPack& equipment_pack,
    const icon::ChannelFactory* channel_factory = nullptr);

absl::StatusOr<std::unique_ptr<AnalogDigitalIOInterface>> CreateAnalogDigitalIO(
    const intrinsic_proto::resources::ResourceHandle& resource_handle,
    const icon::ChannelFactory* channel_factory = nullptr);

}  // namespace intrinsic::skills

#endif  // INTRINSIC_ICON_SKILLS_UTIL_ANALOG_DIGITAL_IO_FACTORY_H_
