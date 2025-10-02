// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_CONTROL_C_API_C_RTCL_ACTION_H_
#define INTRINSIC_ICON_CONTROL_C_API_C_RTCL_ACTION_H_

#include "intrinsic/icon/control/c_api/c_action_factory_context.h"
#include "intrinsic/icon/control/c_api/c_feature_interfaces.h"
#include "intrinsic/icon/control/c_api/c_realtime_signal_access.h"
#include "intrinsic/icon/control/c_api/c_realtime_slot_map.h"
#include "intrinsic/icon/control/c_api/c_realtime_status.h"
#include "intrinsic/icon/control/c_api/c_streaming_io_realtime_access.h"
#include "intrinsic/icon/control/c_api/c_types.h"

#ifdef __cplusplus
extern "C" {
#endif

struct IntrinsicIconServerFunctions {
  IntrinsicIconActionFactoryContextVtable action_factory_context;
  IntrinsicIconRealtimeSlotMapVtable realtime_slot_map;
  IntrinsicIconFeatureInterfaceVtable feature_interfaces;
  IntrinsicIconStreamingIoRealtimeAccessVtable streaming_io_access;
  IntrinsicIconRealtimeSignalAccessVtable realtime_signal;
};

struct IntrinsicIconRtclAction;

struct IntrinsicIconStateVariableValue {
  enum IntrinsicIconStateVariableType { kDouble, kBool, kInt64, kNone };
  union Value {
    double double_value;
    bool bool_value;
    int64_t int64_value;
  };
  union Value value;
  IntrinsicIconStateVariableType type;
};

struct IntrinsicIconRtclActionVtable {
  // Creates an IntrinsicIconRtclAction instance.
  // `action_factory_context` and the storage behind `params_any_proto` are
  // owned by the caller.
  //
  // The pointers in `server_functions` are guaranteed to outlive the
  // IntrinsicIconRtclAction instance. This function should save
  // `server_functions` as part of the IntrinsicIconRtclAction instance, so that
  // the Action can call those functions on any objects received from the
  // server:
  // * The IntrinsicIconActionFactoryContext passed into this function
  // * The IntrinsicIconRealtimeSlotMap passed into on_enter, sense and control
  //   * Any FeatureInterface pointers retrieved from the
  //   IntrinsicIconRealtimeSlotMap
  // * The IntrinsicIconStreamingIoRealtimeAccess passed into sense
  //
  // Returns an IntrinsicIconRealtimeStatus to indicate success or failure.
  //
  // Writes a pointer to the newly created Action to `action_ptr_out` on
  // success. The caller assumes ownership of that pointer and is responsible
  // for calling IntrinsicIconRtclActionDestroy on it before the end of the
  // program.
  IntrinsicIconRealtimeStatus (*create)(
      IntrinsicIconServerFunctions server_functions,
      IntrinsicIconStringView params_any_proto,
      IntrinsicIconActionFactoryContext* action_factory_context,
      IntrinsicIconRtclAction** action_ptr_out);

  // Deletes an IntrinsicIconRtclAction instance.
  // `self` is owned by the caller.
  void (*destroy)(IntrinsicIconRtclAction* self);

  // (Re-)Sets any realtime state of the Action and prepares for cyclic
  // execution. Can use `slot_map` to read the state of any Slots the Action has
  // access to.
  //
  // The ICON realtime control layer calls this once each time the Action
  // becomes active, before the corresponding cycle's Sense().
  IntrinsicIconRealtimeStatus (*on_enter)(
      IntrinsicIconRtclAction* self,
      const IntrinsicIconRealtimeSlotMap* slot_map);

  // Updates the state of the Action at the beginning of a cycle, including:
  // * Reading information (via `slot_map`) from Slots it controls
  // * State Variables (which are used to evaluate conditions for Reactions)
  // * Action-specific data (for example, sampling a trajectory based on the
  //   number of ticks the Action has been active for)
  // * Reading from and writing to streaming I/O values via
  //   `streaming_io_access`
  //
  // (Not every Action needs to do all of the above)
  //
  // Timeslicer calls this at the beginning of each cycle, just after the Parts
  // read their current status from the HAL.
  //
  // `slot_map` is *guaranteed* to have the slots that are registered in the
  // Action's signature, with the same mapping between Slot names and
  // RealtimeSlotIds that was handed to the Action's Factory.
  //
  // Similarly, `streaming_io_access` is guaranteed to provide access to all
  // streaming I/Os that were registered in the factory, if any.
  IntrinsicIconRealtimeStatus (*sense)(
      IntrinsicIconRtclAction* self,
      const IntrinsicIconRealtimeSlotMap* slot_map,
      IntrinsicIconStreamingIoRealtimeAccess* streaming_io_access,
      IntrinsicIconRealtimeSignalAccess* signal_access);

  // Sends commands to the Slots the Action controls, via `slot_map`. Should not
  // modify the externally visible state of the Action (i.e. State Variables).
  //
  // This is called at the end of each cycle, just before Parts apply their
  // commands.
  //
  // `slot_map` is *guaranteed* to have the slots that are registered in the
  // Action's signature, with the same mapping between Slot names and
  // RealtimeSlotIds that was handed to the Action's Factory.
  IntrinsicIconRealtimeStatus (*control)(
      IntrinsicIconRtclAction* self, IntrinsicIconRealtimeSlotMap* slot_map);

  // Gets the value of the requested state variable. `name_size` is the string
  // length of `name`. `self`, `name` and `state_variable_out` are owned by the
  // caller.
  //
  // Returns an IntrinsicIconRealtimeStatus to indicate success or failure.
  // Populates `state_variable_out` on success.
  IntrinsicIconRealtimeStatus (*get_state_variable)(
      const IntrinsicIconRtclAction* self, const char* name, size_t name_size,
      IntrinsicIconStateVariableValue* state_variable_out);
};

#ifdef __cplusplus
}
#endif
#endif  // INTRINSIC_ICON_CONTROL_C_API_C_RTCL_ACTION_H_
