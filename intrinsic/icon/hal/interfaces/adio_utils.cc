// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/interfaces/adio_utils.h"

#include <cstdint>
#include <string>

#include "flatbuffers/buffer.h"
#include "flatbuffers/detached_buffer.h"
#include "flatbuffers/flatbuffer_builder.h"
#include "intrinsic/icon/hal/interfaces/adio.fbs.h"

namespace intrinsic_fbs {

flatbuffers::DetachedBuffer BuildDigitalOutputCommand(std::string name,
                                                      bool value,
                                                      uint32_t bit_number) {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);
  auto name_fbs = builder.CreateString(name);
  builder.Finish(CreateDigitalOutputCommand(
      builder, /*name=*/name_fbs, /*value=*/value, /*bit_number=*/bit_number));
  return builder.Release();
}

flatbuffers::DetachedBuffer BuildDigitalInputStatus(std::string name,
                                                    bool value,
                                                    uint32_t bit_number) {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);
  auto name_fbs = builder.CreateString(name);
  builder.Finish(CreateDigitalInputStatus(
      builder, /*name=*/name_fbs, /*value=*/value, /*bit_number=*/bit_number));
  return builder.Release();
}

flatbuffers::DetachedBuffer BuildAnalogOutputCommand(std::string name,
                                                     AnalogInputUnit unit,
                                                     double value,
                                                     bool is_enabled) {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);
  auto name_fbs = builder.CreateString(name);
  builder.Finish(CreateAnalogOutputCommand(builder, /*name=*/name_fbs,
                                           /*unit=*/unit, /*value=*/value,
                                           /*is_enabled=*/is_enabled));
  return builder.Release();
}

flatbuffers::DetachedBuffer BuildAnalogInputStatus(std::string name,
                                                   AnalogInputUnit unit,
                                                   double value,
                                                   bool is_enabled) {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);
  auto name_fbs = builder.CreateString(name);
  builder.Finish(CreateAnalogInputStatus(builder, /*name=*/name_fbs,
                                         /*unit=*/unit, /*value=*/value,
                                         /*is_enabled=*/is_enabled));
  return builder.Release();
}

}  // namespace intrinsic_fbs
