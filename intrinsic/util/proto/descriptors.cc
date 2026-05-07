// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/proto/descriptors.h"

#include <queue>
#include <string>

#include "absl/container/flat_hash_map.h"
#include "absl/container/flat_hash_set.h"
#include "absl/status/status.h"
#include "absl/strings/str_format.h"
#include "google/protobuf/descriptor.h"
#include "google/protobuf/descriptor.pb.h"
#include "google/protobuf/descriptor_database.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {

namespace {

void MergeFileDescriptorSetImpl(const google::protobuf::Descriptor& descriptor,
                                google::protobuf::FileDescriptorSet& set,
                                bool include_source_code_info) {
  absl::flat_hash_set<std::string> visited;
  for (const auto& file : set.file()) {
    visited.emplace(file.name());
  }
  std::queue<const google::protobuf::FileDescriptor*> queue;
  queue.push(descriptor.file());
  while (!queue.empty()) {
    const google::protobuf::FileDescriptor* current = queue.front();
    queue.pop();
    if (visited.contains(current->name())) {
      // There's already a file descriptor with that name, so assume the
      // dependencies will be there too.
      continue;
    }
    visited.emplace(current->name());
    google::protobuf::FileDescriptorProto* proto = set.add_file();
    current->CopyTo(proto);
    if (include_source_code_info) {
      current->CopySourceCodeInfoTo(proto);
    }
    for (int i = 0; i < current->dependency_count(); i++) {
      queue.push(current->dependency(i));
    }
  }
}

void MergeFileDescriptorSetsImpl(
    const google::protobuf::FileDescriptorSet& src_set,
    google::protobuf::FileDescriptorSet& dst_set,
    bool include_source_code_info) {
  absl::flat_hash_set<std::string> dst_files;
  for (const auto& dst_file : dst_set.file()) {
    dst_files.emplace(dst_file.name());
  }
  for (const google::protobuf::FileDescriptorProto& src_file : src_set.file()) {
    if (dst_files.contains(src_file.name())) {
      continue;
    }
    dst_files.insert(src_file.name());
    google::protobuf::FileDescriptorProto* proto = dst_set.add_file();
    *proto = src_file;
    if (!include_source_code_info) {
      proto->clear_source_code_info();
    }
  }
}

}  // namespace

google::protobuf::FileDescriptorSet GenFileDescriptorSet(
    const google::protobuf::Descriptor& descriptor) {
  google::protobuf::FileDescriptorSet out;
  MergeFileDescriptorSetImpl(descriptor, out,
                             /*include_source_code_info=*/false);
  return out;
}

google::protobuf::FileDescriptorSet GenFileDescriptorSetWithSourceInfo(
    const google::protobuf::Descriptor& descriptor) {
  google::protobuf::FileDescriptorSet out;
  MergeFileDescriptorSetImpl(descriptor, out,
                             /*include_source_code_info=*/true);
  return out;
}

void MergeFileDescriptorSet(const google::protobuf::Descriptor& descriptor,
                            google::protobuf::FileDescriptorSet& set) {
  MergeFileDescriptorSetImpl(descriptor, set,
                             /*include_source_code_info=*/false);
}

void MergeFileDescriptorSetWithSourceInfo(
    const google::protobuf::Descriptor& descriptor,
    google::protobuf::FileDescriptorSet& set) {
  MergeFileDescriptorSetImpl(descriptor, set,
                             /*include_source_code_info=*/true);
}

void MergeFileDescriptorSet(const google::protobuf::FileDescriptorSet& src_set,
                            google::protobuf::FileDescriptorSet& dst_set) {
  MergeFileDescriptorSetsImpl(src_set, dst_set,
                              /*include_source_code_info=*/false);
}

void MergeFileDescriptorSetWithSourceInfo(
    const google::protobuf::FileDescriptorSet& src_set,
    google::protobuf::FileDescriptorSet& dst_set) {
  MergeFileDescriptorSetsImpl(src_set, dst_set,
                              /*include_source_code_info=*/true);
}

absl::Status AddToDescriptorDatabase(
    google::protobuf::SimpleDescriptorDatabase* db,
    const google::protobuf::FileDescriptorProto& file_descriptor,
    absl::flat_hash_map<std::string, google::protobuf::FileDescriptorProto>&
        file_by_name) {
  for (const std::string& dependency : file_descriptor.dependency()) {
    if (auto fd_iter = file_by_name.find(dependency);
        fd_iter != file_by_name.end()) {
      // add dependency now and remove from elements to visit
      google::protobuf::FileDescriptorProto dependency_file_descriptor =
          fd_iter->second;
      file_by_name.erase(fd_iter);
      INTR_RETURN_IF_ERROR(AddToDescriptorDatabase(
          db, dependency_file_descriptor, file_by_name));
    }
  }
  if (!db->Add(file_descriptor)) {
    return absl::InvalidArgumentError(
        absl::StrFormat("Failed to add descriptor '%s' to descriptor database",
                        file_descriptor.name()));
  }
  return absl::OkStatus();
}

absl::Status PopulateDescriptorDatabase(
    google::protobuf::SimpleDescriptorDatabase* db,
    const google::protobuf::FileDescriptorSet& file_descriptor_set) {
  absl::flat_hash_map<std::string, google::protobuf::FileDescriptorProto>
      file_by_name;
  for (const google::protobuf::FileDescriptorProto& file_descriptor :
       file_descriptor_set.file()) {
    file_by_name[file_descriptor.name()] = file_descriptor;
  }
  while (!file_by_name.empty()) {
    google::protobuf::FileDescriptorProto file_descriptor =
        file_by_name.begin()->second;
    file_by_name.erase(file_by_name.begin());
    INTR_RETURN_IF_ERROR(
        AddToDescriptorDatabase(db, file_descriptor, file_by_name));
  }
  return absl::OkStatus();
}

}  // namespace intrinsic
