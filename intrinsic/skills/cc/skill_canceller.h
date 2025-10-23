// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SKILLS_CC_SKILL_CANCELLER_H_
#define INTRINSIC_SKILLS_CC_SKILL_CANCELLER_H_

#include <memory>
#include <string>

#include "absl/base/thread_annotations.h"
#include "absl/functional/any_invocable.h"
#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/mutex.h"
#include "absl/synchronization/notification.h"
#include "absl/time/time.h"

namespace intrinsic {
namespace skills {

// Supports cooperative cancellation of skills by the skill service.
//
// When a cancellation request is received, the skill should:
// 1) stop as soon as possible and leave resources in a safe and recoverable
//    state;
// 2) Return absl::CancelledError.
//
// The skill must call Ready() once it is ready to be cancelled.
//
// A skill can implement cancellation in one of three ways:
// 1) Poll cancelled(), and safely cancel if and when it becomes true.
// 2) Call Wait() and interrupt the skill if it returns true.
//    Of course, Wait() is blocking, so the typical pattern is to call it from
//    an extra thread. If no cancellation happens and the skill finishes, a
//    cleanup function should unblock the extra thread with StopWait().
//    Here is an example:
//    {
//      intrinsic::Thread thread([&context] {
//        if (!context.canceller().Wait(absl::InfiniteDuration())) {
//          // StopWait() was called (or timeout happened).
//          return;
//        }
//        // Cancel long-running work (for example, TryCancel() gRPC context).
//      });
//      absl::Cleanup stop_wait = [&context] {
//        context.canceller().StopWait();
//      };
//      // Do long-running work (for example, call gRPC server).
//    }
// 3) Register a callback via RegisterCallback(). This callback will be invoked
//    when the skill receives a cancellation request. Important: Any references
//    captured by the callback must outlive the SkillCanceller (longer
//    than the skill function). In practice, this means that only
//    std::shared_ptr can be captured.
class SkillCanceller {
 public:
  virtual ~SkillCanceller() = default;

  // True if the skill has received a cancellation request.
  virtual bool cancelled() const = 0;

  // Signals that the skill is ready to be cancelled.
  virtual void Ready() = 0;

  // Sets a callback that will be invoked when a cancellation is requested.
  //
  // Only one callback may be registered, and the callback will be called at
  // most once. It must be registered before calling Ready().
  // Important: Any references captured by the callback must outlive
  // SkillCanceller (longer than the skill). If this is not guaranteed, use one
  // of the other two methods Wait() or cancelled() instead.
  //
  // After a successful invocation, the callback should return absl::OkStatus().
  // Not returning absl::OkStatus() will indicate that the skill could not
  // be cancelled and the skill will be considered to be in an error state.
  // Only return non-OK status if, after cancellation, the skill was not able to
  // leave resources in a safe and recoverable state.
  virtual absl::Status RegisterCallback(
      absl::AnyInvocable<absl::Status() const> callback) = 0;

  // Waits for the skill to be cancelled.
  //
  // Returns true if the skill was cancelled.
  // Returns false if the timeout expired or StopWait() was called.
  virtual bool Wait(absl::Duration timeout) = 0;

  // Unblocks Wait() if it is waiting.
  virtual void StopWait() = 0;
};

// A SkillCanceller used by the skill service to cancel skills.
class SkillCancellationManager : public SkillCanceller {
 public:
  explicit SkillCancellationManager(
      absl::Duration ready_timeout,
      absl::string_view operation_name = "operation");

  bool cancelled() const override {
    absl::MutexLock lock(&mutex_);
    return cancelled_;
  };

  // Sets the cancelled flag, notifies all waiters, and calls the callback if
  // set.
  absl::Status Cancel() ABSL_LOCKS_EXCLUDED(mutex_);

  void Ready() override { ready_.Notify(); };

  absl::Status RegisterCallback(
      absl::AnyInvocable<absl::Status() const> callback) override;

  bool Wait(absl::Duration timeout) override;

  void StopWait() override {
    absl::MutexLock lock(&mutex_);
    stop_wait_ = true;
  };

  // Waits for the skill to be ready for cancellation.
  absl::Status WaitForReady() ABSL_LOCKS_EXCLUDED(mutex_);

 private:
  bool CancelledOrStopWait() const ABSL_EXCLUSIVE_LOCKS_REQUIRED(mutex_) {
    return cancelled_ || stop_wait_;
  };
  const absl::Duration ready_timeout_;
  const std::string operation_name_;
  absl::Notification ready_;
  mutable absl::Mutex mutex_;
  bool cancelled_ ABSL_GUARDED_BY(mutex_) = false;
  bool stop_wait_ ABSL_GUARDED_BY(mutex_) = false;
  std::unique_ptr<absl::AnyInvocable<absl::Status() const>> callback_
      ABSL_GUARDED_BY(mutex_);
};

}  // namespace skills
}  // namespace intrinsic

#endif  // INTRINSIC_SKILLS_CC_SKILL_CANCELLER_H_
