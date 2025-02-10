// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/hardware_module_runtime.h"

#include <atomic>
#include <cstddef>
#include <functional>
#include <list>
#include <memory>
#include <optional>
#include <string_view>
#include <utility>
#include <vector>

#include "absl/base/thread_annotations.h"
#include "absl/log/check.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/mutex.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "absl/types/span.h"
#include "grpcpp/server_builder.h"
#include "intrinsic/icon/hal/hardware_interface_handle.h"
#include "intrinsic/icon/hal/hardware_interface_registry.h"
#include "intrinsic/icon/hal/hardware_interface_traits.h"
#include "intrinsic/icon/hal/hardware_module_init_context.h"
#include "intrinsic/icon/hal/hardware_module_interface.h"
#include "intrinsic/icon/hal/hardware_module_util.h"
#include "intrinsic/icon/hal/icon_state_register.h"  // IWYU pragma: keep
#include "intrinsic/icon/hal/interfaces/hardware_module_state.fbs.h"
#include "intrinsic/icon/hal/interfaces/hardware_module_state_utils.h"
#include "intrinsic/icon/hal/interfaces/icon_state.fbs.h"
#include "intrinsic/icon/interprocess/remote_trigger/remote_trigger_server.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/domain_socket_server.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/domain_socket_utils.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/shared_memory_manager.h"
#include "intrinsic/icon/testing/realtime_annotations.h"
#include "intrinsic/icon/utils/async_buffer.h"
#include "intrinsic/icon/utils/async_request.h"
#include "intrinsic/icon/utils/clock.h"
#include "intrinsic/icon/utils/fixed_string.h"
#include "intrinsic/icon/utils/log.h"
#include "intrinsic/icon/utils/metrics_logger.h"
#include "intrinsic/icon/utils/realtime_metrics.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_macro.h"
#include "intrinsic/platform/common/buffers/rt_promise.h"
#include "intrinsic/platform/common/buffers/rt_queue.h"
#include "intrinsic/platform/common/buffers/rt_queue_multi_writer.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/util/thread/rt_thread.h"
#include "intrinsic/util/thread/thread.h"
#include "intrinsic/util/thread/thread_options.h"

namespace intrinsic::icon {

namespace hardware_interface_traits {
INTRINSIC_ADD_HARDWARE_INTERFACE(intrinsic_fbs::HardwareModuleState,
                                 intrinsic_fbs::BuildHardwareModuleState,
                                 "intrinsic_fbs.HardwareModuleState")
}  // namespace hardware_interface_traits

static constexpr absl::Duration kMetricsExportInterval = absl::Seconds(1);

class HardwareModuleRuntime::CallbackHandler final {
  struct AsyncRequestData {
    intrinsic_fbs::StateCode from;
    intrinsic_fbs::StateCode to;
    icon::FixedString<256> message;
    Clock::time_point timestamp;
  };
  using AsyncRequest =
      intrinsic::icon::AsyncRequest<AsyncRequestData, icon::RealtimeStatus>;

 public:
  explicit CallbackHandler(
      absl::string_view name, HardwareModuleInterface* instance,
      intrinsic_fbs::HardwareModuleState* hardware_module_state) noexcept
      : instance_(instance),
        shared_memory_hardware_module_state_(hardware_module_state),
        request_queue_(10) {
    SetStateDirectly(intrinsic_fbs::StateCode::kDeactivated, "",
                     /*force=*/true);
  }
  ~CallbackHandler() {
    Shutdown();
    QCHECK(action_lock_.TryLock())
        << "CallbackHandler destroyed while an action is still "
           "ongoing - this is likely a bug in the "
           "HardwareModuleRuntime shutdown logic.";
  }

  void SetMetricsLogger(MetricsLogger* metrics_logger) {
    metrics_logger_ = metrics_logger;
  }

  void SetCycleTimeMetricsHelper(CycleTimeMetricsHelper* metrics_helper) {
    metrics_helper_ = metrics_helper;
  }

  // Server callback to trigger `Prepare` on the hardware module.
  void OnPrepare() {
    switch (hardware_module_state_code_) {
      case intrinsic_fbs::StateCode::kActivated:
      case intrinsic_fbs::StateCode::kMotionEnabled:
      case intrinsic_fbs::StateCode::kMotionEnabling:
      case intrinsic_fbs::StateCode::kMotionDisabling:
      case intrinsic_fbs::StateCode::kFaulted:
      case intrinsic_fbs::StateCode::kClearingFaults:
      case intrinsic_fbs::StateCode::kPreparing:
        OnDeactivate();
        break;
      default:
        break;
    }

    if (!SetStateDirectly(intrinsic_fbs::StateCode::kPreparing)) {
      return;
    }
    CancelPendingRequests("Request cancelled by a call to Prepare()");
    if (auto ret = instance_->Prepare(); !ret.ok()) {
      INTRINSIC_RT_LOG_THROTTLED(ERROR)
          << "PUBLIC: Call to 'Prepare' failed: " << ret.message();
      SetStateDirectly(intrinsic_fbs::StateCode::kFatallyFaulted,
                       ret.message());
    } else {
      SetStateDirectly(intrinsic_fbs::StateCode::kPrepared, "");
    }
  }

  // Server callback for trigger `Activate` on the hardware module.
  void OnActivate() {
    // The ICON main loop shall not be running yet, so we can and must set the
    // shared memory state directly.
    if (!SetStateDirectly(intrinsic_fbs::StateCode::kActivating, "")) {
      return;
    }
    CancelPendingRequests("Request cancelled due to activation");
    if (auto ret = instance_->Activate(); !ret.ok()) {
      INTRINSIC_RT_LOG_THROTTLED(ERROR)
          << "PUBLIC: Call to 'Activate' failed: " << ret.message();
      SetStateDirectly(intrinsic_fbs::StateCode::kFatallyFaulted,
                       ret.message());
    } else {
      SetStateDirectly(intrinsic_fbs::StateCode::kActivated, "");
    }
    reject_new_requests_ = false;
  }

  // Server callback for trigger `Deactivate` on the hardware module.
  void OnDeactivate() {
    // The ICON main loop shall not be running anymore, so we
    // can and must set the shared memory state directly.
    if (!SetStateDirectly(intrinsic_fbs::StateCode::kDeactivating, "")) {
      return;
    }
    // It is possible that ongoing calls (e.g. EnableMotion) might miss the
    // `CancelPendingRequests()`, but we cannot get the `non_rt_buffer_lock_`
    // here. The worst that could happen is that unlucky requests will time
    // out and the remaining data will be cleaned up safely in the destructor or
    // on the next call to `SetStateAndWait()`.
    reject_new_requests_ = true;
    CancelPendingRequests("Request cancelled due to deactivation");

    if (auto ret = instance_->Deactivate(); !ret.ok()) {
      INTRINSIC_RT_LOG_THROTTLED(ERROR)
          << "PUBLIC: Call to 'Deactivate' failed: " << ret.message();
      SetStateDirectly(intrinsic_fbs::StateCode::kFatallyFaulted,
                       ret.message());
    } else {
      SetStateDirectly(intrinsic_fbs::StateCode::kDeactivated, "");
    }
  }

  // Server callback for trigger `EnableMotion` on the hardware module.
  void OnEnableMotion() {
    absl::MutexLock lock(&action_lock_);
    if (!SetStateAndWait(hardware_module_state_code_,
                         intrinsic_fbs::StateCode::kMotionEnabling)) {
      return;
    }
    if (auto ret = instance_->EnableMotion(); !ret.ok()) {
      INTRINSIC_RT_LOG_THROTTLED(ERROR)
          << "PUBLIC: Call to 'EnableMotion' failed: " << ret.message();
      SetStateAndWait(intrinsic_fbs::StateCode::kMotionEnabling,
                      ret.code() == absl::StatusCode::kAborted
                          ? intrinsic_fbs::StateCode::kFatallyFaulted
                          : intrinsic_fbs::StateCode::kFaulted,
                      ret.message());
    } else {
      SetStateAndWait(intrinsic_fbs::StateCode::kMotionEnabling,
                      intrinsic_fbs::StateCode::kMotionEnabled, "");
    }
  }

  // Server callback for trigger `DisableMotion` on the hardware module.
  void OnDisableMotion() {
    absl::MutexLock lock(&action_lock_);

    if (!SetStateAndWait(hardware_module_state_code_,
                         intrinsic_fbs::StateCode::kMotionDisabling)) {
      return;
    }
    INTRINSIC_RT_LOG_THROTTLED(INFO) << "PUBLIC: 'DisableMotion' called.";
    if (auto ret = instance_->DisableMotion(); !ret.ok()) {
      INTRINSIC_RT_LOG_THROTTLED(ERROR)
          << "PUBLIC: Call to 'DisableMotion' failed: " << ret.message();
      SetStateAndWait(intrinsic_fbs::StateCode::kMotionDisabling,
                      ret.code() == absl::StatusCode::kAborted
                          ? intrinsic_fbs::StateCode::kFatallyFaulted
                          : intrinsic_fbs::StateCode::kFaulted,
                      ret.message());
    } else {
      SetStateAndWait(intrinsic_fbs::StateCode::kMotionDisabling,
                      intrinsic_fbs::StateCode::kActivated, "");
    }
  }

  // Server callback for trigger `ClearFaults` on the hardware module.
  void OnClearFaults() {
    absl::MutexLock lock(&action_lock_);
    if (!SetStateAndWait(hardware_module_state_code_,
                         intrinsic_fbs::StateCode::kClearingFaults)) {
      return;
    }
    if (auto ret = instance_->ClearFaults(); !ret.ok()) {
      INTRINSIC_RT_LOG_THROTTLED(ERROR)
          << "PUBLIC: Call to 'ClearFaults' failed: " << ret.message();
      SetStateAndWait(intrinsic_fbs::StateCode::kClearingFaults,
                      ret.code() == absl::StatusCode::kAborted
                          ? intrinsic_fbs::StateCode::kFatallyFaulted
                          : intrinsic_fbs::StateCode::kFaulted,
                      ret.message());
    } else {
      SetStateAndWait(intrinsic_fbs::StateCode::kClearingFaults,
                      intrinsic_fbs::StateCode::kActivated, "");
    }
  }

  // Server callback for trigger `ReadStatus` on the hardware module.
  void OnReadStatus() INTRINSIC_CHECK_REALTIME_SAFE {
    // The HWM state must only be written in the RT thread, when the HWM is
    // activated. Therefore, the processing of requests must take place in this
    // function, which is always called when the HWM is activated.
    ProcessNextPendingRequest();

    // Trigger the transition hook in the first cycle where the current state is
    // `kMotionEnabled` (set by the Runtime in the previous cycle), so that ICON
    // calls`ApplyCommand()` as well in this cycle. Don't call this function
    // directly after setting `kMotionEnabled` since ICON won't know about this
    // state change until the next cycle and won't call `ApplyCommand()` in this
    // cycle yet.
    CheckAndTriggerEnabledTransitionHook(
        previous_cycle_hardware_module_state_code_,
        hardware_module_state_code_);

    // Logs realtime metrics on creation and destruction.
    // After the potential metrics reset
    // in`CheckAndTriggerEnabledTransitionHook`.
    ReadStatusScope read_status_scope(
        metrics_helper_, /*is_active=*/hardware_module_state_code_ ==
                             intrinsic_fbs::StateCode::kMotionEnabled);

    if (auto ret = instance_->ReadStatus();
        !ret.ok() && hardware_module_state_code_ !=
                         intrinsic_fbs::StateCode::kClearingFaults) {
      INTRINSIC_RT_LOG_THROTTLED(ERROR)
          << "PUBLIC: Call to 'ReadStatus' failed: " << ret.message();
      if (SetStateDirectly(ret.code() == absl::StatusCode::kAborted
                               ? intrinsic_fbs::StateCode::kFatallyFaulted
                               : intrinsic_fbs::StateCode::kFaulted,
                           ret.message())) {
        // Cancel all requests since we got a new error.
        CancelPendingRequests("Request cancelled due to error in ReadStatus");
      }
    }

    // Export the current metrics to non-realtime every Second.
    if (absl::Now() >= next_metrics_export_) {
      if (metrics_logger_ && metrics_helper_ &&
          !metrics_logger_->AddCycleTimeMetrics(metrics_helper_->Metrics())) {
        INTRINSIC_RT_LOG_THROTTLED(WARNING)
            << "Failed to add cycle time metrics to "
               "metrics logger. Is the queue full?";
      }
      next_metrics_export_ = absl::Now() + kMetricsExportInterval;
    }
  }

  // Server callback for trigger `ApplyCommand` on the hardware module.
  void OnApplyCommand() INTRINSIC_CHECK_REALTIME_SAFE {
    if (hardware_module_state_code_ ==
        intrinsic_fbs::StateCode::kMotionDisabling) {
      // This happens in the first cycle after disabling the motion.
      return;
    } else if (hardware_module_state_code_ !=
               intrinsic_fbs::StateCode::kMotionEnabled) [[unlikely]] {
      auto message = "PUBLIC: 'ApplyCommand' called while not enabled.";
      INTRINSIC_RT_LOG_THROTTLED(WARNING) << message;
      if (SetStateDirectly(intrinsic_fbs::StateCode::kFaulted, message)) {
        // Cancel all requests since we got a new error.
        CancelPendingRequests("Request cancelled due to error in ApplyCommand");
      }
      return;
    }

    // Logs realtime metrics on creation and destruction.
    ApplyCommandScope apply_command_scope(
        metrics_helper_, /*is_active=*/hardware_module_state_code_ ==
                             intrinsic_fbs::StateCode::kMotionEnabled);

    if (auto ret = instance_->ApplyCommand(); !ret.ok()) {
      INTRINSIC_RT_LOG_THROTTLED(ERROR)
          << "PUBLIC: Call to 'ApplyCommand' failed: " << ret.message();
      if (SetStateDirectly(ret.code() == absl::StatusCode::kAborted
                               ? intrinsic_fbs::StateCode::kFatallyFaulted
                               : intrinsic_fbs::StateCode::kFaulted,
                           ret.message())) {
        // Cancel all requests since we got a new error.
        CancelPendingRequests("Request cancelled due to error in ApplyCommand");
      }
    }

  }

  // Sets the internal state *and* the state in shared memory directly. Only
  // call this, when you *know* that no other thread/process might be reading
  // the state in shared memory at the same time.
  // Returns true, if the state was set and different than before.
  bool SetStateDirectly(intrinsic_fbs::StateCode state,
                        absl::string_view fault_reason = "", bool force = false,
                        bool silent = false) INTRINSIC_CHECK_REALTIME_SAFE {
    if (auto result =
            HardwareModuleTransitionGuard(hardware_module_state_code_, state);
        !force && result != TransitionGuardResult::kAllowed) {
      if (!silent && result == TransitionGuardResult::kProhibited) {
        INTRINSIC_RT_LOG_THROTTLED(ERROR)
            << "Switching from "
            << EnumNameStateCode(hardware_module_state_code_) << " to "
            << EnumNameStateCode(state) << " is prohibited!";
      }
      return false;
    }
    if (!silent && hardware_module_state_code_ != state) {
      if (fault_reason.empty()) {
        INTRINSIC_RT_LOG(INFO) << "Switching from "
                               << EnumNameStateCode(hardware_module_state_code_)
                               << " to " << EnumNameStateCode(state);
      } else {
        INTRINSIC_RT_LOG(INFO) << "Switching from "
                               << EnumNameStateCode(hardware_module_state_code_)
                               << " to " << EnumNameStateCode(state)
                               << " with message '" << fault_reason << "'";
      }
    }
    if (hardware_module_state_code_ == state &&
        GetMessage(shared_memory_hardware_module_state_) == fault_reason) {
      // Don't update timestamp when state and message is the same.
      return false;
    }

    const bool state_changed = hardware_module_state_code_ != state;
    // We need to copy the string, since we potentially create a new string in
    // the if-block.
    RealtimeStatus::MessageType fault_reason_string = fault_reason;
    if (state_changed) {
      if (auto status = CheckAndTriggerDisabledTransitionHook(
              hardware_module_state_code_);
          !status.ok()) {
        if (state != intrinsic_fbs::StateCode::kInitFailed &&
            state != intrinsic_fbs::StateCode::kFatallyFaulted &&
            state != intrinsic_fbs::StateCode::kFaulted) {
          // Only set the state to faulted, if it is not already faulted. So
          // that we do not overwrite the original fault reason.
          INTRINSIC_RT_LOG(ERROR)
              << "PUBLIC: Disabled() failed: " << status.message();
          fault_reason_string =
              RealtimeStatus::StrCat("Disabled() failed: ", status.message());
          state = intrinsic_fbs::StateCode::kFaulted;
        }
      }
    }
    hardware_module_state_code_ = state;
    hardware_module_state_update_time_ = Clock::Now();
    intrinsic_fbs::SetState(shared_memory_hardware_module_state_, state,
                            fault_reason_string);
    // Publish the state for non-rt threads. We can use it here without a lock
    // since this function should not be called in parallel.
    auto& hwm_state_buffer = ABSL_TS_UNCHECKED_READ(hwm_state_buffer_);
    intrinsic_fbs::SetState(hwm_state_buffer.GetFreeBuffer(), state,
                            fault_reason_string);
    hwm_state_buffer.CommitFreeBuffer();
    return state_changed;
  }

  // Cancels all pending requests. Call only from ICON lockstep thread or when
  // you *know* that the ICON lockstep thread is not running.
  void CancelPendingRequests(absl::string_view cancel_reason)
      INTRINSIC_CHECK_REALTIME_SAFE {
    while (!request_queue_.reader()->Empty()) {
      AsyncRequest async_request = std::move(
          *request_queue_.reader()
               ->Front());  // Needs to be move so that it will get
                            // destroyed when leaving this scope. Otherwise
                            // the future in the other thread will wait
                            // forever for this promise to get destroyed.
      INTRINSIC_RT_LOG(INFO)
          << "Canceling request to switch to "
          << intrinsic_fbs::EnumNameStateCode(async_request.GetRequest().to)
          << ": " << cancel_reason;
      // Do not use `async_request.Cancel()` so that we convey a message.
      (void)async_request.SetResponse(CancelledError(cancel_reason));
      request_queue_.reader()->DropFront();
    }
  }

  void Shutdown() {
    absl::MutexLock lock(
        &non_rt_buffer_lock_);  // Lock the write buffer so that no new request
                                // can be inserted in the meanwhile.
    reject_new_requests_ = true;
    CancelPendingRequests("Request cancelled due to shutdown");
  }

  intrinsic_fbs::HardwareModuleState GetHardwareModuleState() {
    absl::MutexLock lock(&non_rt_buffer_lock_);
    intrinsic_fbs::HardwareModuleState* state;
    hwm_state_buffer_.GetActiveBuffer(&state);
    return *state;
  }

 private:
  // Sets the state to `to` and waits until the state has been processed.
  // The functions queues a new rt-promise as a request to change the state in
  // the rt thread and waits for the completion of the promise using a future.
  // Attaches `fault_reason` to the new state.
  //
  // Returns true if the state was set successfully.
  // It checks whether the transition `from` to `to` is allowed. If not, returns
  // false.
  // If the transition would be a no-op, it returns false.
  // Returns false on any error as well.
  bool SetStateAndWait(intrinsic_fbs::StateCode from,
                       intrinsic_fbs::StateCode to,
                       absl::string_view fault_reason = "")
      ABSL_EXCLUSIVE_LOCKS_REQUIRED(
          action_lock_)  // Should only be called from non-rt actions such as
                         // EnableMotion, so the lock should be held.
      INTRINSIC_NON_REALTIME_ONLY {
    if (auto result = HardwareModuleTransitionGuard(from, to);
        result != TransitionGuardResult::kAllowed) {
      if (result == TransitionGuardResult::kProhibited) {
        INTRINSIC_RT_LOG(ERROR)
            << "Switching from " << EnumNameStateCode(from) << " to "
            << EnumNameStateCode(to) << " is prohibited!";
      }
      return false;
    }

    {
      // Check all abandoned futures if they can be destroyed. This can only
      // happen, when Deactivate() is called while another Action is active and
      // the timing is very unlucky.
      absl::MutexLock lock(&non_rt_buffer_lock_);
      for (auto it = future_hospice_.begin(); it != future_hospice_.end();) {
        if ((**it).CanBeDestroyed()) {
          it = future_hospice_.erase(it);
        } else {
          it++;
        }
      }
      constexpr size_t kAbandonedFutureWarnLimit = 100;
      if (future_hospice_.size() >= kAbandonedFutureWarnLimit) {
        LOG_EVERY_N_SEC(WARNING, 180)
            << "Found " << future_hospice_.size()
            << " abandoned futures. This indicates a bug in the "
               "HardwareModuleRuntime::CallbackHandler.";
      }
    }

    auto state_change_status = [&]() -> absl::Status {
      auto future = std::make_unique<
          intrinsic::NonRealtimeFuture<icon::RealtimeStatus>>();
      INTR_ASSIGN_OR_RETURN(auto promise, future->GetPromise());
      {
        absl::MutexLock lock(&non_rt_buffer_lock_);
        if (reject_new_requests_) {
          return absl::FailedPreconditionError(
              "Request cancelled due to deactivation");
        }
        INTR_RETURN_IF_ERROR(request_queue_writer_.Insert(
            AsyncRequest(AsyncRequestData{from, to, fault_reason, Clock::Now()},
                         std::move(promise))));
      }

      // Timeout until the state should have been processed. The state is
      // processed in every realtime cycle, so 10 seconds should be sufficient
      // and never be reached.
      constexpr absl::Duration kStatechangeRequestTimeout = absl::Seconds(10);
      absl::StatusOr<RealtimeStatus> status =
          future->GetWithTimeout(kStatechangeRequestTimeout);

      // If we get a timeout, it is likely that the future can also not be
      // destroyed. This can happen, when Deactivate() is
      // called while another action is active and the timing is very unlucky.
      // To not block until shutdown, we move the future away.
      if (!future->CanBeDestroyed()) {
        absl::MutexLock lock(&non_rt_buffer_lock_);
        future_hospice_.push_back(std::move(future));
      }

      if (!status.ok()) {
        return status.status();
      }
      return *status;
    }();
    if (!state_change_status.ok()) {
      INTRINSIC_RT_LOG_THROTTLED(ERROR)
          << "State change request to " << EnumNameStateCode(to)
          << " failed: " << state_change_status.message();
      return false;
    }
    return true;
  }

  // Processes the next pending request in the queue to change the HWM state.
  // Reports the result back via the rt-promise.
  //
  // If the request is outdated, cancels the request.
  // Checks if the request is valid. If the request is valid, applies it to the
  // HWM state. Otherwise, reports the error via the rt-promise.
  void ProcessNextPendingRequest() INTRINSIC_CHECK_REALTIME_SAFE {
    previous_cycle_hardware_module_state_code_ =
        hardware_module_state_code_.load();
    if (!request_queue_.reader()->Empty()) {
      // We need to move the promise (contained in AsyncRequest) so that it will
      // get destroyed when leaving this scope. Otherwise the future in the
      // other thread will wait forever for this promise to get destroyed.
      AsyncRequest item = std::move(*request_queue_.reader()->Front());
      AsyncRequestData newest_data = item.GetRequest();

      request_queue_.reader()->DropFront();

      absl::Status status;
      if (newest_data.timestamp >= hardware_module_state_update_time_ &&
          newest_data.from == hardware_module_state_code_) {
        const bool allowed = SetStateDirectly(
            newest_data.to, absl::string_view(newest_data.message),
            /*force=*/false, /*silent=*/false);
        status = item.SetResponse(
            allowed
                ? OkStatus()
                : FailedPreconditionError(RealtimeStatus::StrCat(
                      "Transition from ",
                      EnumNameStateCode(hardware_module_state_code_), " to ",
                      EnumNameStateCode(newest_data.to), " is prohibited!")));
      } else {
        status = item.SetResponse(
            CancelledError("Request cancelled due to newer request"));
      }
      if (!status.ok()) {
        INTRINSIC_RT_LOG_THROTTLED(ERROR)
            << "Failed to set reply to non rt-call: " << status.message();
      }
    }
  }

  // Checks if the transition from `from` to `to` is `kMotionEnabling` to
  // `kMotionEnabled` and, if so, calls `Enabled()` on the hardware
  // module. Must be called from the rt thread and just after
  // `SetStateDirectly()`.
  void CheckAndTriggerEnabledTransitionHook(intrinsic_fbs::StateCode from,
                                            intrinsic_fbs::StateCode to)
      INTRINSIC_CHECK_REALTIME_SAFE {
    if (from == intrinsic_fbs::StateCode::kMotionEnabling &&
        to == intrinsic_fbs::StateCode::kMotionEnabled) {
      {
        if (auto status = instance_->Enabled(); !status.ok()) {
          SetStateDirectly(intrinsic_fbs::StateCode::kFaulted,
                           RealtimeStatus::StrCat("Enabled() callback failed: ",
                                                  status.message()),
                           /*force=*/false, /*silent=*/true);
        } else {
          INTRINSIC_RT_LOG(INFO) << "Motion Enabled";
          if (metrics_helper_) {
            // Only gather metrics when enabled, because we expect the
            // clock to be ticked correctly only then.
            //
            // Consider using ResetReadStatusStart and logging in a
            // different callback for extended use cases.
            metrics_helper_->Reset();
          }
        }
      }
    }
  }

  // Checks if `from` state is `kMotionEnabled` and, if so, calls `Disabled()`
  // on the hardware module. It only needs to check the `from` state, since
  // `Disabled()` needs to be called for every transition from `kMotionEnabled`.
  RealtimeStatus CheckAndTriggerDisabledTransitionHook(
      intrinsic_fbs::StateCode from) INTRINSIC_CHECK_REALTIME_SAFE {
    if (from == intrinsic_fbs::StateCode::kMotionEnabled) {
      {
        INTRINSIC_RT_RETURN_IF_ERROR(instance_->Disabled());
        INTRINSIC_RT_LOG(INFO) << "Motion Disabled";
      }
    }
    return OkStatus();
  }

  HardwareModuleInterface* instance_;
  absl::Mutex action_lock_;
  // Current state of the HWM that should only be used from the RT thread and
  // resides in the shared memory. ICON reads this state.
  intrinsic_fbs::HardwareModuleState* shared_memory_hardware_module_state_;
  // Current state of the HWM that can be used from multiple threads.
  std::atomic<intrinsic_fbs::StateCode> hardware_module_state_code_ =
      intrinsic_fbs::StateCode::kDeactivated;
  // State of the HWM from the previous cycle. `ProcessNextPendingRequest()`
  // updates this variable before it calls `ReadStatus()` on the HWM. Therefore,
  // this variable only updates while the ICON main loop is running.
  std::atomic<intrinsic_fbs::StateCode>
      previous_cycle_hardware_module_state_code_ =
          intrinsic_fbs::StateCode::kDeactivated;
  // Provides the HWM state from rt threads to non-rt threads.
  AsyncBuffer<intrinsic_fbs::HardwareModuleState> hwm_state_buffer_
      ABSL_GUARDED_BY(non_rt_buffer_lock_);

  // The thread safe queue of pending requests that. The rt thread will read
  // from the queue in `OnReadStatus()`.
  intrinsic::RealtimeQueue<AsyncRequest> request_queue_;
  // The thread safe writer for the `request_queue_` that can be written from
  // all non rt-callbacks (onEnable, etc.).
  absl::Mutex non_rt_buffer_lock_;
  intrinsic::RealtimeQueueMultiWriter<AsyncRequest> request_queue_writer_
      ABSL_GUARDED_BY(non_rt_buffer_lock_){*request_queue_.writer()};
  // While this is true, HardwareModuleRuntime rejects any new non-rt requests
  // (e.g. EnableMotion()).
  std::atomic_bool reject_new_requests_ = false;

  // Timestamp of when the last update of the hwm state was executed. Used to
  // prevent applying updates that are outdated.
  Clock::time_point hardware_module_state_update_time_ = Clock::Now();
  // Container that holds futures that were not ready to be destroyed when the
  // request ended.
  // `SetStateAndWait()` and the destructor will clean up this container if the
  // futures are ready to be destroyed.
  std::list<std::unique_ptr<intrinsic::NonRealtimeFuture<icon::RealtimeStatus>>>
      future_hospice_ ABSL_GUARDED_BY(non_rt_buffer_lock_);

  intrinsic::icon::MetricsLogger* metrics_logger_ = nullptr;
  CycleTimeMetricsHelper* metrics_helper_ = nullptr;
  absl::Time next_metrics_export_ = absl::InfinitePast();
};

absl::StatusOr<HardwareModuleRuntime> HardwareModuleRuntime::Create(
    std::unique_ptr<SharedMemoryManager> shared_memory_manager,
    HardwareModule hardware_module) {
  // Locks the name used by this module. Ensures only a single instance can
  // run at a time. Fails if the lock can't be acquired within the timeout.
  INTR_ASSIGN_OR_RETURN(auto domain_socket_server,
                        DomainSocketServer::Create(
                            /*socket_directory=*/SocketDirectoryFromNamespace(
                                shared_memory_manager->SharedMemoryNamespace()),
                            /*module_name=*/shared_memory_manager->ModuleName(),
                            DomainSocketServer::kDefaultLockAcquireTimeout));

  // C++ allows evaluating function parameters in any order, and we move
  // shared_memory_manager into one of the parameters of the
  // HardwareModuleRuntime constructor below. Because we can't be sure when this
  // happens in relation to the other parameters, we can't create
  // HardwareInterfaceRegistry "in-place", or its constructor might access the
  // moved-from unique_ptr that shared_memory_manager leaves behind.
  auto registry = HardwareInterfaceRegistry(*shared_memory_manager);
  HardwareModuleRuntime runtime(std::move(hardware_module), std::move(registry),
                                std::move(shared_memory_manager),
                                std::move(domain_socket_server));
  INTR_RETURN_IF_ERROR(runtime.Connect());
  return runtime;
}

HardwareModuleRuntime::HardwareModuleRuntime(
    HardwareModule hardware_module,
    HardwareInterfaceRegistry interface_registry,
    std::unique_ptr<SharedMemoryManager> shared_memory_manager,
    std::unique_ptr<DomainSocketServer> domain_socket_server)
    : interface_registry_(std::move(interface_registry)),
      shared_memory_manager_(std::move(shared_memory_manager)),
      hardware_module_(std::move(hardware_module)),
      domain_socket_server_(std::move(domain_socket_server)),
      callback_handler_(nullptr),
      activate_server_(nullptr),
      deactivate_server_(nullptr),
      prepare_server_(nullptr),
      enable_motion_server_(nullptr),
      disable_motion_server_(nullptr),
      read_status_server_(nullptr),
      apply_command_server_(nullptr),
      stop_requested_(std::make_unique<std::atomic<bool>>(false)) {}

HardwareModuleRuntime::~HardwareModuleRuntime() {
  if (callback_handler_) {
    callback_handler_->Shutdown();
  }

  if (stop_requested_) {
    stop_requested_->store(true);
  }
  if (state_change_thread_.joinable()) {
    LOG(INFO)
        << "Joining state change thread - this could be blocked by frozen "
           "callbacks such as EnableMotion";
    state_change_thread_.join();
  }
}

HardwareModuleRuntime::HardwareModuleRuntime(HardwareModuleRuntime&& other)
    : interface_registry_(std::move(other.interface_registry_)),
      shared_memory_manager_(std::move(other.shared_memory_manager_)),
      hardware_module_(std::exchange(other.hardware_module_, HardwareModule())),
      domain_socket_server_(std::move(other.domain_socket_server_)),
      callback_handler_(std::exchange(other.callback_handler_, nullptr)),
      activate_server_(std::exchange(other.activate_server_, nullptr)),
      deactivate_server_(std::exchange(other.deactivate_server_, nullptr)),
      prepare_server_(std::exchange(other.prepare_server_, nullptr)),
      enable_motion_server_(
          std::exchange(other.enable_motion_server_, nullptr)),
      disable_motion_server_(
          std::exchange(other.disable_motion_server_, nullptr)),
      clear_faults_server_(std::exchange(other.clear_faults_server_, nullptr)),
      read_status_server_(std::exchange(other.read_status_server_, nullptr)),
      apply_command_server_(
          std::exchange(other.apply_command_server_, nullptr)),
      hardware_module_state_interface_(
          std::move(other.hardware_module_state_interface_)),
      stop_requested_(std::exchange(other.stop_requested_, nullptr)),
      state_change_thread_(std::move(other.state_change_thread_)) {}

HardwareModuleRuntime& HardwareModuleRuntime::operator=(
    HardwareModuleRuntime&& other) {
  if (&other == this) {
    return *this;
  }
  interface_registry_ = std::move(other.interface_registry_);
  shared_memory_manager_ = std::move(other.shared_memory_manager_);
  hardware_module_ = std::exchange(other.hardware_module_, HardwareModule());
  domain_socket_server_ = std::move(other.domain_socket_server_),
  callback_handler_ = std::exchange(other.callback_handler_, nullptr);
  activate_server_ = std::exchange(other.activate_server_, nullptr);
  deactivate_server_ = std::exchange(other.deactivate_server_, nullptr);
  prepare_server_ = std::exchange(other.prepare_server_, nullptr);
  enable_motion_server_ = std::exchange(other.enable_motion_server_, nullptr);
  disable_motion_server_ = std::exchange(other.disable_motion_server_, nullptr);
  clear_faults_server_ = std::exchange(other.clear_faults_server_, nullptr);
  read_status_server_ = std::exchange(other.read_status_server_, nullptr);
  apply_command_server_ = std::exchange(other.apply_command_server_, nullptr);
  hardware_module_state_interface_ =
      std::move(other.hardware_module_state_interface_);
  stop_requested_ = std::exchange(other.stop_requested_, nullptr);
  state_change_thread_ = std::move(other.state_change_thread_);
  return *this;
}

absl::Status HardwareModuleRuntime::Connect() {
  // Adds an "inbuilt" status segment for the hardware module state.
  INTR_ASSIGN_OR_RETURN(
      hardware_module_state_interface_,
      interface_registry_
          .AdvertiseMutableInterface<intrinsic_fbs::HardwareModuleState>(
              "hardware_module_state",
              intrinsic_fbs::BuildHardwareModuleState()));

  callback_handler_ = std::make_unique<CallbackHandler>(
      hardware_module_.config.GetName(), hardware_module_.instance.get(),
      *hardware_module_state_interface_);

  // Adds an "inbuilt" status segment for ICON to publish its state (e.g.
  // current cycle).
  INTR_ASSIGN_OR_RETURN(
      icon_state_interface_,
      interface_registry_.AdvertiseInterface<intrinsic_fbs::IconState>(
          kIconStateInterfaceName));

  INTR_ASSIGN_OR_RETURN(
      auto activate_server,
      RemoteTriggerServer::Create(
          *shared_memory_manager_, "activate",
          std::bind(&HardwareModuleRuntime::CallbackHandler::OnActivate,
                    callback_handler_.get())));
  activate_server_ =
      std::make_unique<RemoteTriggerServer>(std::move(activate_server));

  INTR_ASSIGN_OR_RETURN(
      auto deactivate_server,
      RemoteTriggerServer::Create(
          *shared_memory_manager_, "deactivate",
          std::bind(&HardwareModuleRuntime::CallbackHandler::OnDeactivate,
                    callback_handler_.get())));
  deactivate_server_ =
      std::make_unique<RemoteTriggerServer>(std::move(deactivate_server));

  INTR_ASSIGN_OR_RETURN(
      auto prepare_server,
      RemoteTriggerServer::Create(
          *shared_memory_manager_, "prepare",
          std::bind(&HardwareModuleRuntime::CallbackHandler::OnPrepare,
                    callback_handler_.get())));
  prepare_server_ =
      std::make_unique<RemoteTriggerServer>(std::move(prepare_server));

  INTR_ASSIGN_OR_RETURN(
      auto enable_motion_server,
      RemoteTriggerServer::Create(
          *shared_memory_manager_, "enable_motion",
          std::bind(&HardwareModuleRuntime::CallbackHandler::OnEnableMotion,
                    callback_handler_.get())));
  enable_motion_server_ =
      std::make_unique<RemoteTriggerServer>(std::move(enable_motion_server));

  INTR_ASSIGN_OR_RETURN(
      auto disable_motion_server,
      RemoteTriggerServer::Create(
          *shared_memory_manager_, "disable_motion",
          std::bind(&HardwareModuleRuntime::CallbackHandler::OnDisableMotion,
                    callback_handler_.get())));
  disable_motion_server_ =
      std::make_unique<RemoteTriggerServer>(std::move(disable_motion_server));

  INTR_ASSIGN_OR_RETURN(
      auto clear_faults_server,
      RemoteTriggerServer::Create(
          *shared_memory_manager_, "clear_faults",
          std::bind(&HardwareModuleRuntime::CallbackHandler::OnClearFaults,
                    callback_handler_.get())));
  clear_faults_server_ =
      std::make_unique<RemoteTriggerServer>(std::move(clear_faults_server));

  INTR_ASSIGN_OR_RETURN(
      auto read_status_server,
      RemoteTriggerServer::Create(
          *shared_memory_manager_, "read_status",
          std::bind(&HardwareModuleRuntime::CallbackHandler::OnReadStatus,
                    callback_handler_.get())));
  read_status_server_ =
      std::make_unique<RemoteTriggerServer>(std::move(read_status_server));

  INTR_ASSIGN_OR_RETURN(
      auto apply_command_server,
      RemoteTriggerServer::Create(
          *shared_memory_manager_, "apply_command",
          std::bind(&HardwareModuleRuntime::CallbackHandler::OnApplyCommand,
                    callback_handler_.get())));
  apply_command_server_ =
      std::make_unique<RemoteTriggerServer>(std::move(apply_command_server));

  return absl::OkStatus();
}

absl::Status HardwareModuleRuntime::Run(grpc::ServerBuilder& server_builder,
                                        bool is_realtime,
                                        const std::vector<int>& cpu_affinity) {
  if (activate_server_ == nullptr) {
    return absl::InternalError(
        "PUBLIC: Hardware module does not seem to be connected. Did you call "
        "`Connect()`?");
  }

  // Helper lambda to set the state to kInitFailed if any of the
  // initialization steps below fail.
  auto set_init_failed_on_error = [this](absl::Status status) -> absl::Status {
    if (!status.ok()) {
      callback_handler_->SetStateDirectly(intrinsic_fbs::StateCode::kInitFailed,
                                          status.message());
    }
    return status;
  };
  HardwareModuleInitContext init_context(interface_registry_, server_builder,
                                         hardware_module_.config);
  const auto init_status =
      set_init_failed_on_error(hardware_module_.instance->Init(init_context));
  if (!init_status.ok()) {
    LOG(ERROR) << "Initializing the module failed with: " << init_status;
  }

  if (!domain_socket_server_) {
    return absl::InternalError(
        "Run was called, but domain_socket_server_ is nullptr. This should "
        "never happen.");
  }
  // Segments added after this call will not be visible to DomainSocketServer
  // and its clients (like ICON).
  INTR_RETURN_IF_ERROR(domain_socket_server_->AddSegmentInfoServeShmDescriptors(
      *shared_memory_manager_));

  // Ensures that no methods on the uninitialized module can be called.
  INTR_RETURN_IF_ERROR(init_status);

  if (const absl::Duration cycle_duration =
          init_context.GetCycleDurationForCycleTimeMetrics();
      cycle_duration != absl::ZeroDuration()) {
    auto cycle_time_metrics_helper =
        intrinsic::icon::CycleTimeMetricsHelper::Create(
            cycle_duration,
            /*log_cycle_time_warnings=*/init_context
                .AreCycleTimeWarningsEnabled());

    if (cycle_time_metrics_helper.ok()) {
      cycle_time_metrics_helper_ = std::make_unique<CycleTimeMetricsHelper>(
          std::move(*cycle_time_metrics_helper));
    } else {
      LOG(ERROR) << "Failed to create cycle time metrics helper: "
                 << cycle_time_metrics_helper.status();
    }

    if (cycle_time_metrics_helper_ != nullptr) {
      metrics_logger_ =
          std::make_unique<MetricsLogger>(hardware_module_.config.GetName());
      if (const auto status = metrics_logger_->Start(); !status.ok()) {
        LOG(WARNING) << "Failed to start metrics logger: " << status;
        metrics_logger_ = nullptr;
      } else {
        callback_handler_->SetMetricsLogger(metrics_logger_.get());
        callback_handler_->SetCycleTimeMetricsHelper(
            cycle_time_metrics_helper_.get());
        LOG(INFO) << "Cycle time metrics gathering is enabled with a cycle "
                     "duration of "
                  << cycle_duration << ". Cycle time warnings are "
                  << (init_context.AreCycleTimeWarningsEnabled() ? "enabled"
                                                                 : "disabled ")
                  << ".";
      }
    }
  }

  intrinsic::ThreadOptions state_change_thread_options;
  state_change_thread_options.SetName("StateChange");

  intrinsic::ThreadOptions activate_thread_options;
  activate_thread_options.SetName("Activate");

  intrinsic::ThreadOptions read_status_thread_options;
  read_status_thread_options.SetName("ReadStatus");

  intrinsic::ThreadOptions apply_command_thread_options;
  apply_command_thread_options.SetName("ApplyCommand");

  if (is_realtime) {
    state_change_thread_options.SetRealtimeLowPriorityAndScheduler();
    state_change_thread_options.SetAffinity(cpu_affinity);
    activate_thread_options.SetRealtimeLowPriorityAndScheduler();
    activate_thread_options.SetAffinity(cpu_affinity);
    read_status_thread_options.SetRealtimeHighPriorityAndScheduler();
    read_status_thread_options.SetAffinity(cpu_affinity);
    apply_command_thread_options.SetRealtimeHighPriorityAndScheduler();
    apply_command_thread_options.SetAffinity(cpu_affinity);
  }

  intrinsic::ThreadOptions deactivate_thread_options = activate_thread_options;
  deactivate_thread_options.SetName("Deactivate");

  auto state_change_query = [](std::atomic<bool>* stop_requested,
                               RemoteTriggerServer* prepare_server,
                               RemoteTriggerServer* enable_motion_server,
                               RemoteTriggerServer* disable_motion_server,
                               RemoteTriggerServer* clear_faults_server) {
    while (!stop_requested->load()) {
      prepare_server->Query();
      enable_motion_server->Query();
      disable_motion_server->Query();
      clear_faults_server->Query();
    }
  };

  if (auto thread = CreateRealtimeCapableThread(
          state_change_thread_options, state_change_query,
          stop_requested_.get(), prepare_server_.get(),
          enable_motion_server_.get(), disable_motion_server_.get(),
          clear_faults_server_.get());
      thread.ok()) {
    state_change_thread_ = *std::move(thread);
    INTR_RETURN_IF_ERROR(set_init_failed_on_error(absl::OkStatus()));
  } else {
    INTR_RETURN_IF_ERROR(set_init_failed_on_error(thread.status()));
  }

  INTR_RETURN_IF_ERROR(set_init_failed_on_error(
      activate_server_->StartAsync(activate_thread_options)));
  INTR_RETURN_IF_ERROR(set_init_failed_on_error(
      deactivate_server_->StartAsync(deactivate_thread_options)));
  INTR_RETURN_IF_ERROR(set_init_failed_on_error(
      read_status_server_->StartAsync(read_status_thread_options)));
  INTR_RETURN_IF_ERROR(set_init_failed_on_error(
      apply_command_server_->StartAsync(apply_command_thread_options)));

  return absl::OkStatus();
}

absl::Status HardwareModuleRuntime::Stop() {
  callback_handler_->Shutdown();
  apply_command_server_->RequestStop();
  read_status_server_->RequestStop();
  deactivate_server_->RequestStop();
  enable_motion_server_->RequestStop();
  clear_faults_server_->RequestStop();
  prepare_server_->RequestStop();
  activate_server_->RequestStop();
  stop_requested_->store(true);
  auto status = hardware_module_.instance->Shutdown();
  apply_command_server_->JoinAsyncThread();
  read_status_server_->JoinAsyncThread();
  if (state_change_thread_.joinable()) {
    state_change_thread_.join();
  }
  return status;
}

bool HardwareModuleRuntime::IsStarted() const {
  bool started = state_change_thread_.joinable();
  started &= read_status_server_->IsStarted();
  started &= apply_command_server_->IsStarted();

  return started;
}

const HardwareModule& HardwareModuleRuntime::GetHardwareModule() const {
  return hardware_module_;
}

absl::StatusOr<const intrinsic_fbs::HardwareModuleState>
HardwareModuleRuntime::GetHardwareModuleState() const {
  if (callback_handler_ == nullptr) {
    return absl::InternalError(
        "Hardware Module Runtime callback_handler is null");
  }

  return callback_handler_->GetHardwareModuleState();
}

void HardwareModuleRuntime::SetStateTestOnly(intrinsic_fbs::StateCode state,
                                             std::string_view fault_reason) {
  callback_handler_->SetStateDirectly(state, fault_reason, /*force=*/true);
}

}  // namespace intrinsic::icon
