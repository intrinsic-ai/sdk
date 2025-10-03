// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_CONTROL_C_API_C_STREAMING_IO_REALTIME_ACCESS_H_
#define INTRINSIC_ICON_CONTROL_C_API_C_STREAMING_IO_REALTIME_ACCESS_H_

#include "intrinsic/icon/control/c_api/c_action_factory_context.h"
#include "intrinsic/icon/control/c_api/c_realtime_status.h"

#ifdef __cplusplus
extern "C" {
#endif

struct IntrinsicIconStreamingIoRealtimeAccess;

struct IntrinsicIconStreamingIoRealtimeAccessVtable {
  // Polls a streaming input and writes a status into `status_out`.
  //
  // Returns nullptr if nothing has been written to the streaming input for `id`
  // yet.
  // Sets `status_out` to NotFoundError and returns nullptr if there is no
  // streaming input for `id`.
  // Returns a pointer to the current streaming input value and sets
  // `status_out` to Ok on success. The returned pointer remains valid until the
  // next call to this function.
  const IntrinsicIconStreamingInputType* (*poll_input)(
      IntrinsicIconStreamingIoRealtimeAccess* self, uint64_t input_id,
      IntrinsicIconRealtimeStatus* status_out);

  // Copies `size` bytes, starting at `output`, into a buffer that is then made
  // available to the non-realtime thread. NOTE: Because this copies data, using
  // data with a large `size` (>>1024) can break realtime safety!
  //
  // Returns NotFound if there is no streaming output in our
  // RealtimeStreamingIoStorage.
  // Returns an error if `size` exceeds the maximum size for streaming
  // outputs supported by the ICON server.
  IntrinsicIconRealtimeStatus (*write_output)(
      IntrinsicIconStreamingIoRealtimeAccess* self,
      const IntrinsicIconStreamingOutputType* output, size_t size);
};

#ifdef __cplusplus
}
#endif

#endif  // INTRINSIC_ICON_CONTROL_C_API_C_STREAMING_IO_REALTIME_ACCESS_H_
