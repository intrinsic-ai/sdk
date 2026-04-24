// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/interfaces/io_controller_utils.h"

#include <cstddef>
#include <string>
#include <vector>

#include "flatbuffers/buffer.h"
#include "flatbuffers/detached_buffer.h"
#include "flatbuffers/flatbuffer_builder.h"
#include "intrinsic/icon/hal/interfaces/adio.fbs.h"
#include "intrinsic/icon/hal/interfaces/io_controller.fbs.h"

namespace intrinsic_fbs {

flatbuffers::DetachedBuffer BuildAIOStatus(
    const std::vector<std::string>& descriptions) {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);
  std::vector<flatbuffers::Offset<AnalogInputStatus>> tmp_signals;

  for (size_t i = 0; i < descriptions.size(); i++) {
    auto tmp_description = builder.CreateString(descriptions.at(i));
    auto tmp_signal = CreateAnalogInputStatus(builder, tmp_description,
                                              AnalogInputUnit::kUnknown);

    tmp_signals.push_back(tmp_signal);
  }

  auto status = CreateAIOStatusDirect(builder, &tmp_signals);
  builder.Finish(status);
  return builder.Release();
}

flatbuffers::DetachedBuffer BuildAIOCommand(
    const std::vector<std::string>& descriptions) {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);
  std::vector<flatbuffers::Offset<AnalogOutputCommand>> tmp_signals;

  for (size_t i = 0; i < descriptions.size(); i++) {
    auto tmp_description = builder.CreateString(descriptions.at(i));
    auto tmp_signal = CreateAnalogOutputCommand(builder, tmp_description,
                                                AnalogInputUnit::kUnknown);

    tmp_signals.push_back(tmp_signal);
  }

  auto status = CreateAIOCommandDirect(builder, &tmp_signals);
  builder.Finish(status);
  return builder.Release();
}

flatbuffers::DetachedBuffer BuildDIOStatus(
    const std::vector<std::string>& descriptions) {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);
  std::vector<flatbuffers::Offset<DigitalInputStatus>> tmp_signals;

  int bit_number = 0;
  for (const auto& desc : descriptions) {
    auto tmp_name = builder.CreateString(desc);
    // We only want populated fields in the flatbuffer
    DigitalInputStatusBuilder di_builder(builder);
    di_builder.add_name(tmp_name);
    di_builder.add_bit_number(bit_number++);
    di_builder.add_value(false);
    tmp_signals.push_back(di_builder.Finish());
  }

  auto status = CreateDIOStatusDirect(builder, &tmp_signals);
  builder.Finish(status);
  return builder.Release();
}

flatbuffers::DetachedBuffer BuildDIOCommand(
    const std::vector<std::string>& descriptions) {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);
  std::vector<flatbuffers::Offset<DigitalOutputCommand>> tmp_signals;

  int bit_number = 0;
  for (const auto& desc : descriptions) {
    auto tmp_name = builder.CreateString(desc);
    // We only want populated fields in the flatbuffer
    DigitalOutputCommandBuilder do_builder(builder);
    do_builder.add_name(tmp_name);
    do_builder.add_bit_number(bit_number++);
    do_builder.add_value(false);
    tmp_signals.push_back(do_builder.Finish());
  }

  auto status = CreateDIOCommandDirect(builder, &tmp_signals);
  builder.Finish(status);
  return builder.Release();
}

}  // namespace intrinsic_fbs
