// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_SKILLS_UTIL_ADIO_UTIL_H_
#define INTRINSIC_ICON_SKILLS_UTIL_ADIO_UTIL_H_

#include <string>

#include "absl/container/flat_hash_map.h"
#include "absl/status/statusor.h"
#include "intrinsic/hardware/gpio/gpio_service_equipment.pb.h"
#include "intrinsic/icon/cc_client/robot_config.h"
#include "intrinsic/icon/equipment/icon_equipment.pb.h"
#include "intrinsic/resources/proto/resource_handle.pb.h"
#include "intrinsic/skills/proto/equipment.pb.h"

namespace intrinsic::skills {

struct BlockTypeOptions {
  bool digital_in = false;
  bool digital_out = false;
  bool analog_in = false;
  bool analog_out = false;
};

// GetBlockToPartNameMap returns a map of ADIO block names to ICON part names
// based on the provided equipment config and robot config.
//
// The function takes an Icon2AdioPart equipment config, a RobotConfig, and a
// BlockTypeOptions struct as input. The BlockTypeOptions struct specifies which
// types of blocks to include in the map (digital input, digital output, analog
// input, analog output).
//
// The function returns a map of block names to part names. The block names are
// the names of the ADIO blocks as configured in the ICON part configuration.
// The part names are the names of the ICON parts that contain the ADIO blocks.
//
// If the equipment config contains any ICON parts that are not present in the
// robot config, or if the ICON parts do not have configured IO signals, the
// function will return an error.
//
// If the equipment config contains any ADIO blocks that have the same name as
// another ADIO block, the function will return an error.
absl::StatusOr<absl::flat_hash_map<std::string, std::string>>
GetBlockToPartNameMap(
    const intrinsic_proto::icon::Icon2AdioPart& adio_equipment_config,
    const icon::RobotConfig& robot_config,
    const BlockTypeOptions& block_type_options);

}  // namespace intrinsic::skills

#endif  // INTRINSIC_ICON_SKILLS_UTIL_ADIO_UTIL_H_
