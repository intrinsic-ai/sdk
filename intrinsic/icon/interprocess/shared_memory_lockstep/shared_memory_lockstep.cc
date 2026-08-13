// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#include "intrinsic/icon/interprocess/shared_memory_lockstep/shared_memory_lockstep.h"

#include <utility>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/domain_socket_utils.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/memory_segment.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/shared_memory_manager.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/util/thread/lockstep.h"

namespace intrinsic::icon {

namespace {

absl::StatusOr<SharedMemoryLockstep> GetSharedMemoryLockstep(
    const SegmentNameToFileDescriptorMap& segment_name_to_file_descriptor_map,
    absl::string_view memory_name) {
  INTR_ASSIGN_OR_RETURN(auto segment,
                        ReadWriteMemorySegment<Lockstep>::Get(
                            segment_name_to_file_descriptor_map, memory_name));
  return SharedMemoryLockstep(std::move(segment));
}

}  // namespace

bool SharedMemoryLockstep::Connected() const {
  if (!memory_segment_.IsValid()) {
    return false;
  }
  return memory_segment_.Header().WriterRefCount() == 2;
}

absl::StatusOr<SharedMemoryLockstep> CreateSharedMemoryLockstep(
    SharedMemoryManager& manager, absl::string_view memory_name) {
  INTR_RETURN_IF_ERROR(manager.AddSegment(memory_name, false, Lockstep()));

  INTR_ASSIGN_OR_RETURN(
      auto segment, manager.Get<ReadWriteMemorySegment<Lockstep>>(memory_name));

  return SharedMemoryLockstep(std::move(segment));
}
absl::StatusOr<SharedMemoryLockstep> GetSharedMemoryLockstep(
    const SharedMemoryManager& manager, absl::string_view memory_name) {
  return GetSharedMemoryLockstep(manager.SegmentNameToFileDescriptorMap(),
                                 memory_name);
}

}  // namespace intrinsic::icon
