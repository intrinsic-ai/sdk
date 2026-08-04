// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_PROTO_DESCRIPTOR_POOLS_H_
#define INTRINSIC_UTIL_PROTO_DESCRIPTOR_POOLS_H_

#include <memory>
#include <string_view>

#include "absl/status/statusor.h"
#include "google/protobuf/descriptor.h"
#include "google/protobuf/message.h"

namespace intrinsic {

// Creates a new, empty proto instance of the message type with the given name
// using the prototype message of the given descriptor pool and message factory.
//
// The returned message is owned by the caller and must be destroyed before the
// given pool and factory are destroyed (the message has references to pool and
// factory).
absl::StatusOr<std::unique_ptr<google::protobuf::Message>>
CreateProtoInstanceFromDescriptorPool(
    std::string_view message_type_name,
    const google::protobuf::DescriptorPool* desc_pool,
    google::protobuf::MessageFactory* msg_factory);

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_PROTO_DESCRIPTOR_POOLS_H_
