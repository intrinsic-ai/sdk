// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_CONTROL_C_API_EXTERNAL_ACTION_API_TESTING_ICON_ACTION_FACTORY_CONTEXT_FAKE_H_
#define INTRINSIC_ICON_CONTROL_C_API_EXTERNAL_ACTION_API_TESTING_ICON_ACTION_FACTORY_CONTEXT_FAKE_H_

#include <functional>
#include <memory>
#include <string>
#include <utility>

#include "intrinsic/icon/control/c_api/c_action_factory_context.h"
#include "intrinsic/icon/control/c_api/external_action_api/icon_action_factory_context.h"
#include "intrinsic/icon/control/c_api/external_action_api/testing/icon_realtime_signal_access_and_map_fake.h"
#include "intrinsic/icon/control/c_api/external_action_api/testing/icon_slot_map_fake.h"
#include "intrinsic/icon/control/c_api/external_action_api/testing/icon_streaming_io_registry_fake.h"
#include "intrinsic/icon/proto/v1/types.pb.h"

namespace intrinsic::icon {

// Fake implementation of IconActionFactoryContext that implements the C API.
//
// Use this to test Realtime Action classes by passing it to their Create()
// method via MakeIconActionFactoryContext().
class IconActionFactoryContextFake {
 public:
  // Returns `server_config` and the SlotInfo values from `slot_map` to
  // Action factory methods.
  // Keeps a reference to `streaming_io_registry`. You can test any streaming
  // input / output handlers by calling methods on `streaming_io_registry`.
  IconActionFactoryContextFake(
      intrinsic_proto::icon::v1::ServerConfig server_config,
      IconSlotMapFake& slot_map,
      IconStreamingIoRegistryFake& streaming_io_registry,
      IconRealtimeSignalAccessAndMapFake& realtime_signal_access_and_map)
      : server_config_(std::move(server_config)),
        slot_map_(slot_map),
        streaming_io_registry_(streaming_io_registry),
        realtime_signal_access_and_map_(realtime_signal_access_and_map) {}

  // Returns an IconActionFactoryContext that is backed by this fake. Make sure
  // that this IconActionFactoryContextFake outlives the
  // IconActionFactoryContext.
  IconActionFactoryContext MakeIconActionFactoryContext();

 private:
  static IntrinsicIconActionFactoryContextVtable GetCApiVtable();

  const intrinsic_proto::icon::v1::ServerConfig server_config_;

  IconSlotMapFake& slot_map_;
  IconStreamingIoRegistryFake& streaming_io_registry_;
  IconRealtimeSignalAccessAndMapFake& realtime_signal_access_and_map_;
};

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_CONTROL_C_API_EXTERNAL_ACTION_API_TESTING_ICON_ACTION_FACTORY_CONTEXT_FAKE_H_
