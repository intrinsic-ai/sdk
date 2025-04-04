// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_CC_CLIENT_SESSION_H_
#define INTRINSIC_ICON_CC_CLIENT_SESSION_H_

#include <atomic>
#include <cstdint>
#include <functional>
#include <list>
#include <memory>
#include <optional>
#include <string>
#include <utility>
#include <variant>
#include <vector>

#include "absl/base/thread_annotations.h"
#include "absl/container/flat_hash_map.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/mutex.h"
#include "absl/time/time.h"
#include "absl/types/span.h"
#include "google/protobuf/any.pb.h"
#include "google/rpc/status.pb.h"
#include "grpcpp/client_context.h"
#include "grpcpp/support/sync_stream.h"
#include "intrinsic/icon/cc_client/condition.h"
#include "intrinsic/icon/cc_client/stream.h"
#include "intrinsic/icon/common/id_types.h"
#include "intrinsic/icon/common/slot_part_map.h"
#include "intrinsic/icon/proto/joint_space.pb.h"
#include "intrinsic/icon/proto/streaming_output.pb.h"
#include "intrinsic/icon/proto/v1/service.grpc.pb.h"
#include "intrinsic/icon/proto/v1/service.pb.h"
#include "intrinsic/icon/proto/v1/types.pb.h"
#include "intrinsic/icon/release/source_location.h"
#include "intrinsic/logging/proto/context.pb.h"
#include "intrinsic/production/external/intops/strong_int.h"
#include "intrinsic/util/atomic_sequence_num.h"
#include "intrinsic/util/grpc/channel_interface.h"
#include "intrinsic/util/thread/thread.h"

namespace intrinsic {
namespace icon {

// Client-side identifier for a Reaction.
DEFINE_STRONG_INT_TYPE(ReactionHandle, int64_t);

// Describes a reaction consisting of a condition that is evaluated on the
// robot, and possible events that are triggered when the condition is true.
// A reaction is triggered if
// 1) the condition is true when the reaction becomes active
// or 2) on a rising edge when it is already active.
//
// A reaction is active when its associated action is active or when it is
// added as a free-standing reaction.
class ReactionDescriptor {
 public:
  // Constructs a `ReactionDescriptor` with the given `condition`. Conditions
  // are evaluated in real-time on the robot.
  explicit ReactionDescriptor(const Condition& condition);

  // Associates `handle` with the `ReactionDescriptor`. `handle` can then be
  // used to refer to the reaction and its associated callback. `handle` must be
  // unique in the surrounding `Session`.
  //
  // `loc` should usually not be provided directly, except when defining helper
  // functions – it is used to provide a code location in error messages when a
  // non-unique `ReactionHandle` is encountered.
  ReactionDescriptor& WithHandle(
      ReactionHandle reaction_handle,
      intrinsic::SourceLocation loc = intrinsic::SourceLocation::current());

  // Adds an action to switch to in the real-time context on the robot once the
  // `condition` is fulfilled. Only one `action_id` may be switched to,
  // subsequent calls to WithRealtimeActionOnCondition() or
  // WithParallelRealtimeActionOnCondition() replace the previous `action_id`.
  ReactionDescriptor& WithRealtimeActionOnCondition(ActionInstanceId action_id);

  // Adds an action to start in parallel in the real-time context on the robot
  // when `condition` is fulfilled. Only one `action_id` may be started,
  // subsequent calls to WithParallelRealtimeActionOnCondition() or
  // WithRealtimeActionOnCondition() replace the previous `action_id`.
  //
  // To start multiple actions, add multiple reactions with the same condition
  // but different action id given to WithParallelRealtimeActionOnCondition().
  //
  // The action referenced by `action_id` and the action, to which this reaction
  // is bound, must use a non-overlapping part set since in this case the action
  // with `action_id` cannot start in parallel and would preempt the active
  // action (will result in an error when adding the reaction).
  ReactionDescriptor& WithParallelRealtimeActionOnCondition(
      ActionInstanceId action_id);

  // Adds a callback to be invoked each time the `condition` occurs. Only one
  // `on_condition` may be added, subsequent calls to WithWatcherOnCondition()
  // replace the previous `on_condition`.
  ReactionDescriptor& WithWatcherOnCondition(
      std::function<void()> on_condition);

  // Configures the behavior of *repeated* reaction triggering. The initial
  // triggering is described above.
  //
  // If `enable` is true, the reaction will only trigger once as long
  // as the associated action is active. It can trigger again if the action is
  // executed again. If the reaction is free-standing, it will only trigger
  // once.
  //
  // If `enable` is false, the reaction will trigger again on every (subsequent)
  // rising edge.
  //
  // FireOnce(false) is the default behavior if FireOnce() is not called.
  ReactionDescriptor& FireOnce(bool enable = true);

  // Triggers `signal_name` the first time `condition` is met. The signal
  // remains true thereafter. `signal_name` must be a real-time signal declared
  // by the signature of the `action_id_` which this reaction is attached to. If
  // it is not, ICON will return an error when the Reaction is added.
  //
  // Repeated calls will override the previously set real-time signal. Each
  // reaction may only trigger a single signal.
  ReactionDescriptor& WithRealtimeSignalOnCondition(
      absl::string_view realtime_signal_name);

  // Creates a reaction from `reaction_descriptor`, applied to the action
  // identified by `action_id` or as a free-standing reaction if `action_id` is
  // not set.
  static intrinsic_proto::icon::v1::Reaction ToProto(
      const ReactionDescriptor& reaction_descriptor, ReactionId reaction_id,
      const std::optional<ActionInstanceId>& action_id);

 private:
  friend class Session;

  const Condition condition_;
  std::optional<ActionInstanceId> action_id_;
  std::optional<std::function<void()>> on_condition_;
  std::optional<std::pair<ReactionHandle, intrinsic::SourceLocation>>
      reaction_handle_;
  std::optional<std::string> realtime_signal_name_;
  bool fire_once_ = false;
  bool stop_associated_action_ = false;
};

// Describes an action to be built.
class ActionDescriptor {
 public:
  // Constructs an `ActionDescriptor` for the `action_type_name` with the
  // `action_id`. The `action_type_name` must exist in the ICON server, and the
  // `action_id` must be unique within this Session.
  ActionDescriptor(absl::string_view action_type_name,
                   ActionInstanceId action_id,
                   const SlotPartMap& slot_part_map);
  // Same as above, except that the SlotPartMap is inferred based on the given
  // Action type's signature upon calling AddAction(s). This works only for
  // Actions that have a single Slot.
  //
  // If either of those two conditions is not met, an AddAction(s) call that
  // includes this ActionDescriptor will fail.
  ActionDescriptor(absl::string_view action_type_name,
                   ActionInstanceId action_id, absl::string_view part_name);

  // Adds fixed parameters to the action. No references to `fixed_params` are
  // retained beyond this call. Only one `fixed_params` may be associated
  // with each `ActionDescriptor`, subsequent calls to `WithFixedParams()`
  // replace the previous `fixed_params`.
  ActionDescriptor& WithFixedParams(
      const ::google::protobuf::Message& fixed_params);

  // Adds a reaction to the action. While an action is active, reactions
  // associated with it trigger events. Multiple reactions may be added to a
  // single action, Multiple `ReactionDescriptor`s  may be associated with each
  // `ActionDescriptor`, subsequent calls to `WithReaction()` append to the
  // existing `ReactionDescriptor`s.
  ActionDescriptor& WithReaction(const ReactionDescriptor& reaction_descriptor);

  ActionInstanceId Id() const { return action_id_; }

  std::string ActionTypeName() const { return action_type_name_; }

 private:
  friend class Session;

  const std::string action_type_name_;
  const ActionInstanceId action_id_;
  // Holds either a SlotPartMap, the name of a single Part to infer one from.
  const std::variant<SlotPartMap, std::string> slot_data_;
  std::optional<google::protobuf::Any> fixed_params_;
  std::vector<ReactionDescriptor> reaction_descriptors_;
};

// Provides a handle to the user for an already-created action.
class Action {
 public:
  // Returns the `id` of this action. This corresponds to the `action_id` given
  // to the `ActionDescriptor`.
  ActionInstanceId id() const { return id_; }

 private:
  friend class Session;
  explicit Action(ActionInstanceId id);

  ActionInstanceId id_;
};

// A `Session` scopes control of a set of parts to a single session. The
// `Session` provides the ability to manipulate those parts by adding actions
// and/or reactions.
// All functions in this class are thread-safe. Often, `QuitWatcherLoop()` or
// `End()` are called from other threads to implement skill cancellation.
class Session {
 public:
  // Creates a Session for the `parts` and starts it.
  //
  // The `context` is used to tag the part status. If it is empty the part
  // status is only tagged with the session and action id.
  //
  // The factory returned by `icon_channel.GetClientContextFactory()` is invoked
  // before each gRPC request to obtain a ::grpc::ClientContext.  This provides
  // an opportunity to set client metadata, or other ClientContext settings, for
  // all ICON API requests made by the Session.
  //
  // `deadline` is an optional deadline for establishing the session. If given,
  // it overrides any deadline set by the ClientContext factory.
  static absl::StatusOr<std::unique_ptr<Session>> Start(
      std::shared_ptr<ChannelInterface> icon_channel,
      absl::Span<const std::string> parts,
      const intrinsic_proto::data_logger::Context& context = {},
      std::optional<absl::Time> deadline = std::nullopt);

  // Creates a Session for the `parts` and starts it.
  //
  // The resulting session uses default-constructed ::grpc::ClientContext
  // objects.
  //
  // `deadline` is an optional deadline for establishing the session. If given,
  // it overrides any deadline set by the ClientContext factory.
  static absl::StatusOr<std::unique_ptr<Session>> Start(
      std::unique_ptr<intrinsic_proto::icon::v1::IconApi::StubInterface> stub,
      absl::Span<const std::string> parts,
      const ClientContextFactory& client_context_factory =
          DefaultClientContextFactory,
      const intrinsic_proto::data_logger::Context& context = {},
      std::optional<absl::Time> deadline = std::nullopt);

  // Disallow move.
  Session(Session&&) = delete;
  Session& operator=(Session&&) = delete;

  // Disallow copy.
  Session(const Session&) = delete;
  Session& operator=(const Session&) = delete;

  // Ends the session and logs errors.
  ~Session();

  // Adds the action described by `action_descriptor` to the session. Returns
  // an aborted error if the session ended. Other errors may be returned due to
  // the `action_descriptor` specified, such as an invalid `action_type_name`,
  // an `action_id` that's already in use, etc.
  absl::StatusOr<Action> AddAction(const ActionDescriptor& action_descriptor);

  // Adds the actions described by `action_descriptors` to the session. Returns
  // an aborted error if the session ended. Other errors may be returned due to
  // the `action_descriptors` specified, such as an invalid `action_type_name`,
  // an `action_id` that's already in use, etc.
  absl::StatusOr<std::vector<Action>> AddActions(
      absl::Span<const ActionDescriptor> action_descriptors);

  // Adds the reaction described by `reaction_descriptor` to the session as
  // a free-standing reaction. This reaction is not attached to a specific
  // action but is active as long as the session is active.
  absl::Status AddFreestandingReaction(
      const ReactionDescriptor& reaction_descriptor);

  // Adds the reactions described by `reaction_descriptors` to the session as
  // free-standing reactions. Those reactions are not attached to a specific
  // action but are active as long as the session is active.
  absl::Status AddFreestandingReactions(
      absl::Span<const ReactionDescriptor> reaction_descriptors);

  // Removes the action identified by the `action_id`, as well as any Reactions
  // that originate from or switch to that Action.
  //
  // If the deleted Action is active at the time this command is handled by the
  // realtime system, ICON will switch to the default Action for the Part the
  // Action was running on.
  //
  // N.B. This does not "recycle" `action_id` – no new Action can be added with
  // the same ID.
  absl::Status RemoveAction(ActionInstanceId action_id);

  // Removes the actions identified by the `action_ids`, as well as any
  // Reactions that originate from or switch to those Actions.
  //
  // If any of the deleted Actions are active at the time this command is
  // handled by the realtime system, ICON will switch to the default Action(s)
  // for the respective Part(s).
  //
  // N.B. This does not "recycle" `action_ids` – no new Actions can be added
  // with a previously used ID.
  absl::Status RemoveActions(const std::vector<ActionInstanceId>& action_ids);

  // Removes all Actions and Reactions from the Session. ICON will fall back to
  // the default Action(s), which normally stops the robot.
  //
  // N.B. This essentially invalidates all Action and ReactionHandle objects
  // obtained from this Session.
  absl::Status ClearAllActionsAndReactions();

  // Starts the given actions on the server.
  //
  // All `actions` must have non overlapping part sets. Otherwise this function
  // returns an error.
  //
  // If `stop_active_actions` is true, all active actions will be stopped. All
  // unused parts will fall back to their default action (usually a stop
  // action).
  //
  // If `stop_active_actions` is false, `actions` will be started in parallel.
  // If the part set used by `actions` overlaps with an active action, this
  // active action will be deactivated. All parts that are now unused (in case
  // this active action used parts that are not used by `actions`) will
  // fall back to the default action of that part.
  //
  // Returns AbortedError if the session ended.
  // Other errors may be returned due to the `actions` specified.
  absl::Status StartActions(absl::Span<const Action> actions,
                            bool stop_active_actions = true);
  absl::Status StartAction(const Action& action,
                           bool stop_active_actions = true);

  // Stops all active actions in the session. All parts will fall back to the
  // default action (usually a stop action).
  absl::Status StopAllActions();

  // Runs watchers associated with this `Session` from added reactions.  Blocks
  // until QuitWatcherLoop() is called, the session ends, or the deadline is
  // reached.  All associated callbacks are invoked on the calling thread. If
  // the session dies due to an error during execution, returns an aborted
  // error. If deadline is reached, returns a kDeadlineExceeded error. If the
  // deadline is in the past, this processes queued events before returning
  // kDeadlineExceeded.
  absl::Status RunWatcherLoop(absl::Time deadline = absl::InfiniteFuture());

  // Runs watchers associated with this `Session` from added reactions.
  // Blocks until QuitWatcherLoop() is called, the session ends, the deadline is
  // reached, or after running watchers associated with `reaction_handle`. All
  // associated callbacks are invoked on the calling thread. If the session dies
  // due to an error during execution, returns an aborted error. If the deadline
  // is in the past, this processes queued events before returning
  // kDeadlineExceeded.
  absl::Status RunWatcherLoopUntilReaction(
      ReactionHandle reaction_handle,
      absl::Time deadline = absl::InfiniteFuture());

  // Stops running watchers after the current event is finished processing, if
  // they are running. Watchers may be restarted again by calling RunWatchers().
  void QuitWatcherLoop() ABSL_LOCKS_EXCLUDED(reactions_queue_mutex_);

  // Creates a StreamWriter for the given `input_name` of the given action.
  //
  // Returns an aborted error if the session ended. Other errors may be returned
  // due to mismatched types, an input already in use, etc.
  template <typename T>
  absl::StatusOr<std::unique_ptr<StreamWriterInterface<T>>> StreamWriter(
      const Action& action, absl::string_view input_name) {
    return intrinsic::icon::internal::StreamWriter<T>::Open(
        session_id_, action.id(), input_name, stub_.get(),
        channel_ ? channel_->GetClientContextFactory() : nullptr);
  }

  // Returns the latest output of the Action with `id`. Blocks until `deadline`
  // if that Action is active, but has not published an output value yet.
  absl::StatusOr<::intrinsic_proto::icon::StreamingOutput> GetLatestOutput(
      ActionInstanceId id, absl::Time deadline);

  absl::StatusOr<::intrinsic_proto::icon::JointTrajectoryPVA>
  GetPlannedTrajectory(ActionInstanceId id);

  // Ends the session and returns the session end status. Returns a precondition
  // failed status if the session has already ended.
  absl::Status End()
      ABSL_LOCKS_EXCLUDED(session_mutex_, reactions_queue_mutex_);

  SessionId Id() const { return session_id_; }

 private:
  // Common implementation of Start.
  static absl::StatusOr<std::unique_ptr<Session>> StartImpl(
      const intrinsic_proto::data_logger::Context& context,
      std::shared_ptr<ChannelInterface> icon_channel,
      std::unique_ptr<intrinsic_proto::icon::v1::IconApi::StubInterface> stub,
      absl::Span<const std::string> parts,
      const ClientContextFactory& client_context_factory,
      std::optional<absl::Time> deadline);

  Session(
      std::shared_ptr<ChannelInterface> icon_channel,
      std::unique_ptr<grpc::ClientContext> session_context,
      std::unique_ptr<grpc::ClientReaderWriterInterface<
          intrinsic_proto::icon::v1::OpenSessionRequest,
          intrinsic_proto::icon::v1::OpenSessionResponse>>
          session_stream,
      std::unique_ptr<grpc::ClientContext> watcher_context,
      std::unique_ptr<::grpc::ClientReaderInterface<
          intrinsic_proto::icon::v1::WatchReactionsResponse>>
          watcher_stream,
      std::unique_ptr<intrinsic_proto::icon::v1::IconApi::StubInterface> stub,
      SessionId session_id,
      const intrinsic_proto::data_logger::Context& context,
      ClientContextFactory client_context_factory);

  // Creates a vector of actions from the `action_descriptors`.
  static std::vector<Action> MakeActionVector(
      absl::Span<const ActionDescriptor> action_descriptors);

  // Returns AlreadyExistsError if `reaction_descriptors` contains any
  // ReactionHandle that appears more than once *across both*
  // `reaction_descriptors` and `reaction_handle_to_id_and_loc_`.
  absl::Status CheckReactionHandlesUnique(
      absl::Span<const ReactionDescriptor> reaction_descriptors) const
      ABSL_LOCKS_EXCLUDED(session_mutex_);

  // Saves any callbacks and ReactionHandles contained in
  // `reaction_descriptors_by_id` to `reaction_callback_map_` and
  // `reaction_handle_to_id_and_loc_`.
  void SaveReactionData(
      const absl::flat_hash_map<ReactionId, ReactionDescriptor>&
          reaction_descriptors_by_id) ABSL_LOCKS_EXCLUDED(session_mutex_);

  absl::StatusOr<intrinsic_proto::icon::v1::OpenSessionResponse>
  GetResponseOrEnd(const intrinsic_proto::icon::v1::OpenSessionRequest& request)
      ABSL_LOCKS_EXCLUDED(session_mutex_, reactions_queue_mutex_);

  // Triggers the reaction callbacks for the given `reaction`.
  void TriggerReactionCallbacks(
      const intrinsic_proto::icon::v1::WatchReactionsResponse& reaction)
      ABSL_LOCKS_EXCLUDED(session_mutex_, watch_reactions_mutex_,
                          reactions_queue_mutex_);

  // Reads out the reaction watcher buffer and finishes the call. Logs any
  // additional reactions received, and any call errors.
  void CleanUpWatcherCall() ABSL_LOCKS_EXCLUDED(reactions_queue_mutex_);

  // Converts the `status` to an absl::Status. Ends the session if its code is
  // absl::Status::Code::kAborted and logs the session call status.
  // Returns AbortedError if the Session was closed (either successfully or with
  // an error).
  // Returns the absl::Status version of `status` otherwise. In this case, the
  // Session remains active.
  absl::Status EndAndLogOnAbort(const ::google::rpc::Status& status)
      ABSL_LOCKS_EXCLUDED(session_mutex_, reactions_queue_mutex_);

  // Reads from watcher stream, and queues new reactions into the
  // `reactions_queue_`.
  void WatchReactionsThreadBody();

  bool ShouldWakeWatcherLoop() const
      ABSL_EXCLUSIVE_LOCKS_REQUIRED(reactions_queue_mutex_) {
    return quit_watcher_loop_ || reactions_stream_closed_ ||
           !reactions_queue_.empty();
  }

  // We need multiple mutexes to run some functions in parallel. Define some
  // total order for acquiring to avoid deadlocks.
  mutable absl::Mutex session_mutex_
      ABSL_ACQUIRED_BEFORE(watch_reactions_mutex_);
  absl::Mutex watch_reactions_mutex_
      ABSL_ACQUIRED_BEFORE(reactions_queue_mutex_);
  absl::Mutex reactions_queue_mutex_;

  // Hold onto the channel, if any, so that callers do not need to worry about
  // its lifetime. May be nullptr depending on the version of Start used to
  // construct this session.
  std::shared_ptr<ChannelInterface> channel_;

  // Indicates whether the call is already dead. If so, the `session_stream_`
  // and `watcher_stream_` should no longer be accessed.
  std::atomic<bool> session_ended_ = false;

  // Must outlive `session_stream_`.
  std::unique_ptr<grpc::ClientContext> session_context_
      ABSL_GUARDED_BY(session_mutex_);
  std::unique_ptr<grpc::ClientReaderWriterInterface<
      intrinsic_proto::icon::v1::OpenSessionRequest,
      intrinsic_proto::icon::v1::OpenSessionResponse>>
      session_stream_ ABSL_GUARDED_BY(session_mutex_);

  // Must outlive `watcher_stream_`.
  std::unique_ptr<grpc::ClientContext> watcher_context_
      ABSL_GUARDED_BY(watch_reactions_mutex_);
  // Only call watcher_stream_::Read() on `watcher_read_thread_` until
  // `watcher_read_thread_` has joined.
  std::unique_ptr<grpc::ClientReaderInterface<
      intrinsic_proto::icon::v1::WatchReactionsResponse>>
      watcher_stream_ ABSL_GUARDED_BY(watch_reactions_mutex_);

  // A map of callbacks registered to reactions, keyed by the reaction id.
  absl::flat_hash_map<ReactionId, std::function<void()>> reaction_callback_map_
      ABSL_GUARDED_BY(session_mutex_);

  // Reaction events are written to the `reactions_queue_` from the
  // `watcher_read_thread_`, and read during `RunWatcherLoop()` on the calling
  // thread.
  std::list<absl::StatusOr<intrinsic_proto::icon::v1::WatchReactionsResponse>>
      reactions_queue_ ABSL_GUARDED_BY(reactions_queue_mutex_);
  // Only needed for waking `RunWatcherLoop()` when the stream is closed from
  // server side.
  bool reactions_stream_closed_ ABSL_GUARDED_BY(reactions_queue_mutex_) = false;
  // Signals RunWatcherLoop* to quit after processing the reaction queue.
  bool quit_watcher_loop_ ABSL_GUARDED_BY(reactions_queue_mutex_) = false;

  // Used to read reaction events in the background. `watcher_stream_::Read()`
  // calls should only be made on this thread. It is ok for other
  // watcher_stream_ methods to be invoked on another thread.
  Thread watcher_read_thread_;

  std::unique_ptr<intrinsic_proto::icon::v1::IconApi::StubInterface> stub_;

  // Used to generate unique ReactionIds.
  SequenceNumber<ReactionId> reaction_id_sequence_
      ABSL_GUARDED_BY(session_mutex_);
  absl::flat_hash_map<ReactionHandle,
                      std::pair<ReactionId, intrinsic::SourceLocation>>
      reaction_handle_to_id_and_loc_ ABSL_GUARDED_BY(session_mutex_);

  const SessionId session_id_;

  // Factory function that produces ::grpc::ClientContext objects before each
  // gRPC request. This is required to make new grpc calls on the fly since we
  // need to propagate the original icon connection parameters stored in here.
  const ClientContextFactory client_context_factory_;
};

}  // namespace icon
}  // namespace intrinsic

#endif  // INTRINSIC_ICON_CC_CLIENT_SESSION_H_
