// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/interprocess/remote_trigger/remote_trigger_client.h"

#include <atomic>
#include <optional>
#include <string>
#include <utility>

#include "absl/cleanup/cleanup.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "intrinsic/icon/interprocess/binary_futex.h"
#include "intrinsic/icon/interprocess/remote_trigger/remote_trigger_constants.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/domain_socket_utils.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/memory_segment.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_macro.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::icon {

RemoteTriggerClient::AsyncRequest::AsyncRequest(
    ReadOnlyMemorySegment<BinaryFutex>* response_futex,
    std::atomic<bool>* request_started)
    : response_futex_(response_futex), request_started_(request_started) {}

RemoteTriggerClient::AsyncRequest::AsyncRequest(AsyncRequest&& other) noexcept
    : response_futex_(other.response_futex_),
      request_started_(std::exchange(other.request_started_, nullptr)) {}

RemoteTriggerClient::AsyncRequest& RemoteTriggerClient::AsyncRequest::operator=(
    AsyncRequest&& other) noexcept {
  if (this != &other) {
    response_futex_ = std::move(other.response_futex_);
    request_started_ = std::exchange(other.request_started_, nullptr);
  }

  return *this;
}

RemoteTriggerClient::AsyncRequest::~AsyncRequest() noexcept {
  if (request_started_ != nullptr) {
    request_started_->store(false);
  }
}

bool RemoteTriggerClient::AsyncRequest::Valid() const {
  if (request_started_ == nullptr) {
    return false;
  }
  return request_started_->load();
}

bool RemoteTriggerClient::AsyncRequest::Ready() const {
  if (response_futex_ == nullptr) {
    return false;
  }
  return response_futex_->GetValue().Value() > 0;
}

RealtimeStatus RemoteTriggerClient::AsyncRequest::WaitUntil(
    absl::Time deadline) {
  if (!Valid()) {
    return FailedPreconditionError("async request no longer valid");
  }

  auto response = response_futex_->GetValue().WaitUntil(deadline);
  request_started_->store(false);
  // We have to set the pointer to null, indicating in the destructor that we've
  // completed the async request and already cleared the `request_started_`
  // flag.
  request_started_ = nullptr;
  return response;
}

absl::StatusOr<RemoteTriggerClient> RemoteTriggerClient::Create(
    const SegmentNameToFileDescriptorMap& segment_name_to_file_descriptor_map,
    absl::string_view server_name) {
  RemoteTriggerClient client(server_name, segment_name_to_file_descriptor_map);
  INTR_RETURN_IF_ERROR(client.Connect());
  return client;
}

RemoteTriggerClient::RemoteTriggerClient(
    absl::string_view server_name,
    const SegmentNameToFileDescriptorMap& segment_name_to_file_descriptor_map)
    : server_name_(server_name),
      segment_name_to_file_descriptor_map_(
          segment_name_to_file_descriptor_map) {}

RemoteTriggerClient::RemoteTriggerClient(RemoteTriggerClient&& other) noexcept
    : server_name_(std::exchange(other.server_name_, "")),
      segment_name_to_file_descriptor_map_(
          std::exchange(other.segment_name_to_file_descriptor_map_, {})),
      request_futex_(std::exchange(other.request_futex_,
                                   ReadWriteMemorySegment<BinaryFutex>())),
      response_futex_(std::exchange(other.response_futex_,
                                    ReadOnlyMemorySegment<BinaryFutex>())),
      request_started_(false) {}

RemoteTriggerClient& RemoteTriggerClient::operator=(
    RemoteTriggerClient&& other) noexcept {
  if (this != &other) {
    server_name_ = std::exchange(other.server_name_, "");
    segment_name_to_file_descriptor_map_ =
        std::exchange(other.segment_name_to_file_descriptor_map_, {});
    request_futex_ = std::exchange(other.request_futex_,
                                   ReadWriteMemorySegment<BinaryFutex>());
    response_futex_ = std::exchange(other.response_futex_,
                                    ReadOnlyMemorySegment<BinaryFutex>());
    request_started_.store(false);
  }
  return *this;
}

absl::Status RemoteTriggerClient::Connect() {
  if (IsConnected()) {
    return absl::OkStatus();
  }
  std::string request_memory_name =
      absl::StrCat(server_name_, kSemRequestSuffix);
  std::string response_memory_name =
      absl::StrCat(server_name_, kSemResponseSuffix);
  INTR_ASSIGN_OR_RETURN(
      request_futex_,
      ReadWriteMemorySegment<BinaryFutex>::Get(
          segment_name_to_file_descriptor_map_, request_memory_name));
  INTR_ASSIGN_OR_RETURN(
      response_futex_,
      ReadOnlyMemorySegment<BinaryFutex>::Get(
          segment_name_to_file_descriptor_map_, response_memory_name));
  return absl::OkStatus();
}

bool RemoteTriggerClient::IsConnected() const {
  return request_futex_.IsValid() && response_futex_.IsValid();
}

RealtimeStatus RemoteTriggerClient::Trigger(absl::Time deadline) {
  if (!IsConnected()) {
    return InvalidArgumentError("client not connected");
  }
  if (absl::Now() > deadline) {
    return DeadlineExceededError("specified deadline is in the past");
  }
  if (bool expected = false;
      !request_started_.compare_exchange_strong(expected, true)) {
    return AlreadyExistsError("request already triggered");
  }

  // Clear the `request_started_` flag in any case when we leave the scope of
  // this function.
  absl::Cleanup clear_request_flag = [this] { request_started_.store(false); };

  // Signal the server to start the execution.
  INTRINSIC_RT_RETURN_IF_ERROR(request_futex_.GetValue().Post());
  // Wait for the response from the server.
  return response_futex_.GetValue().WaitUntil(deadline);
}

RealtimeStatusOr<RemoteTriggerClient::AsyncRequest>
RemoteTriggerClient::TriggerAsync() {
  if (!IsConnected()) {
    return InvalidArgumentError("client not connected");
  }

  if (bool expected = false;
      !request_started_.compare_exchange_strong(expected, true)) {
    return AlreadyExistsError("request already triggered");
  }

  // In case the response is already set before we do the request, it is a left
  // over from a previous call. TryWait resets the futex state to 0 in case it
  // was already 1.
  if (response_futex_.GetValue().TryWait() == std::nullopt) {
    return FailedPreconditionError(
        "The server is gone, cannot trigger any more requests");
  }

  INTRINSIC_RT_RETURN_IF_ERROR(request_futex_.GetValue().Post());
  return AsyncRequest(&response_futex_, &request_started_);
}

}  // namespace intrinsic::icon
