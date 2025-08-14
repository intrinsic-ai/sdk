// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_INTERFACES_ADIO_UTILS_H_
#define INTRINSIC_ICON_HAL_INTERFACES_ADIO_UTILS_H_

#include <cstdint>
#include <string>

#include "flatbuffers/detached_buffer.h"
#include "intrinsic/icon/hal/interfaces/adio.fbs.h"

namespace intrinsic_fbs {

flatbuffers::DetachedBuffer BuildDigitalOutputCommand(std::string name,
                                                      bool value,
                                                      uint32_t bit_number);

flatbuffers::DetachedBuffer BuildDigitalInputStatus(std::string name,
                                                    bool value,
                                                    uint32_t bit_number);

flatbuffers::DetachedBuffer BuildAnalogInputStatus(std::string name,
                                                   AnalogInputUnit unit,
                                                   double value,
                                                   bool is_enabled = false);

flatbuffers::DetachedBuffer BuildAnalogOutputCommand(std::string name,
                                                     AnalogInputUnit unit,
                                                     double value,
                                                     bool is_enabled = false);

}  // namespace intrinsic_fbs

#endif  // INTRINSIC_ICON_HAL_INTERFACES_ADIO_UTILS_H_
