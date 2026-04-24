// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/interfaces/force_sensor_utils.h"

#include "flatbuffers/detached_buffer.h"
#include "flatbuffers/flatbuffer_builder.h"
#include "intrinsic/icon/flatbuffers/transform_types.fbs.h"
#include "intrinsic/icon/hal/interfaces/force_sensor.fbs.h"

namespace intrinsic_fbs {

// Creates a detached flatbuffer that stores a message defined as
// ForceSensorStatus.
flatbuffers::DetachedBuffer CreateForceSensorStatusBuffer() {
  intrinsic_fbs::Wrench wrench;

  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);

  builder.Finish(CreateForceSensorStatus(builder, &wrench, &wrench, &wrench,
                                         &wrench,
                                         ForceSensorStatusCode::GenericError));

  return builder.Release();
}

// Creates a detached flatbuffer that stores a message defined as
// ForceSensorCommand.
flatbuffers::DetachedBuffer CreateForceSensorCommandBuffer() {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);
  builder.Finish(CreateForceSensorCommand(builder, /*tare_sensor=*/false,
                                          /*num_taring_cycles=*/1));
  return builder.Release();
}

}  // namespace intrinsic_fbs
