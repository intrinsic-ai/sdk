// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_PROTO_DESCRIPTORS_H_
#define INTRINSIC_UTIL_PROTO_DESCRIPTORS_H_

#include <type_traits>

#include "absl/status/status.h"
#include "google/protobuf/descriptor.h"
#include "google/protobuf/descriptor.pb.h"
#include "google/protobuf/descriptor_database.h"
#include "google/protobuf/message.h"

namespace intrinsic {

// Generates a FileDescriptorSet given a Descriptor. The returned
// FileDescriptorSet includes all transitive dependencies of the descriptor.
google::protobuf::FileDescriptorSet GenFileDescriptorSet(
    const google::protobuf::Descriptor& descriptor);

// Generates a FileDescriptorSet for message type ProtoT. The returned
// FileDescriptorSet includes all transitive dependencies of the ProtoT.
template <class ProtoT, typename = std::enable_if_t<std::is_base_of_v<
                            google::protobuf::Message, ProtoT>>>
google::protobuf::FileDescriptorSet GenFileDescriptorSet() {
  return GenFileDescriptorSet(*ProtoT::GetDescriptor());
}

// Populates the given descriptor database with all file descriptors in the
// given FileDescriptorSet.
absl::Status PopulateDescriptorDatabase(
    google::protobuf::SimpleDescriptorDatabase* db,
    const google::protobuf::FileDescriptorSet& file_descriptor_set);

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_PROTO_DESCRIPTORS_H_
