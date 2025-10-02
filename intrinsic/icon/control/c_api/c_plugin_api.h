// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_CONTROL_C_API_C_PLUGIN_API_H_
#define INTRINSIC_ICON_CONTROL_C_API_C_PLUGIN_API_H_

#include <stddef.h>
#include <stdint.h>

#include "intrinsic/icon/control/c_api/c_realtime_status.h"
#include "intrinsic/icon/control/c_api/c_rtcl_action.h"
#include "intrinsic/icon/control/c_api/c_types.h"

#ifdef __cplusplus
extern "C" {
#endif

// Plugins call this to register one or more Action types.
// `action_type_name` and `action_signature_proto` are owned by the caller.
// `action_signature_proto` is a serialized
// intrinsic_proto::icon::v1::ActionSignature proto.
//
// Returns OkStatus on success.
// Returns AlreadyExists if `action_type_name` is not unique (i.e. there is
// already an Action Type of that name).
// Returns InvalidArgument if `icon_api_version` does not match the server's API
// version.
typedef IntrinsicIconRealtimeStatus (*IntrinsicIconRegisterActionType)(
    int64_t icon_api_version, IntrinsicIconStringView action_type_name,
    IntrinsicIconStringView action_signature_proto,
    IntrinsicIconRtclActionVtable action_vtable);

// Entrypoint into Action plugins. The host process calls this upon loading the
// plugin.
// Returns OkStatus on success.
typedef IntrinsicIconRealtimeStatus (*IntrinsicIconRegisterActionTypes)(
    IntrinsicIconRegisterActionType register_action_type_fn);

// This is the canonical name for the entry point into an IntrinsicIcon custom
// Action Plugin. A Plugin must have a function with this name and the signature
// defined above (IntrinsicIconRegisterActionTypes).
//
// That is, plugins should declare the entrypoint like this:
//
// IntrinsicIconRealtimeStatus INTRINSIC_ICON_PLUGIN_ENTRY_POINT(
//     IntrinsicIconRegisterActionType register_action_type_fn) {
//   ...
// }
#define INTRINSIC_ICON_ACTION_PLUGIN_ENTRY_POINT \
  IntrinsicIconRegisterActionTypesForPlugin

#ifdef __cplusplus
}
#endif

#endif  // INTRINSIC_ICON_CONTROL_C_API_C_PLUGIN_API_H_
