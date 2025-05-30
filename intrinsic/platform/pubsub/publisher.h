// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_PLATFORM_PUBSUB_PUBLISHER_H_
#define INTRINSIC_PLATFORM_PUBSUB_PUBLISHER_H_

#include <memory>
#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "google/protobuf/any.pb.h"
#include "google/protobuf/message.h"

namespace intrinsic {

struct PublisherData;

class Publisher {
 public:
  Publisher(absl::string_view topic_name,
            std::unique_ptr<PublisherData> publisher_data);

  ~Publisher();

  Publisher(const Publisher&) = delete;
  Publisher& operator=(const Publisher&) = delete;
  Publisher(Publisher&&);
  Publisher& operator=(Publisher&&);

  absl::Status Publish(const google::protobuf::Message& message,
                       absl::Time event_time) const;

  absl::Status Publish(const google::protobuf::Message& message) const {
    return Publish(message, absl::Now());
  }

  absl::Status Publish(google::protobuf::Any message,
                       absl::Time event_time) const;

  absl::Status Publish(google::protobuf::Any message) const {
    return Publish(message, absl::Now());
  }

  absl::string_view TopicName() const { return topic_name_; }

  absl::StatusOr<bool> HasMatchingSubscribers();

 private:
  std::string topic_name_ = {};
  std::unique_ptr<PublisherData> publisher_data_ = {};
};

}  // namespace intrinsic

#endif  // INTRINSIC_PLATFORM_PUBSUB_PUBLISHER_H_
