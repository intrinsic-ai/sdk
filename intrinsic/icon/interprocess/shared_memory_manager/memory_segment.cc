// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/interprocess/shared_memory_manager/memory_segment.h"

#include <errno.h>
#include <fcntl.h>
#include <stddef.h>
#include <stdint.h>
#include <string.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>

#include <cerrno>
#include <string>
#include <utility>

#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_join.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/domain_socket_utils.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/segment_header.h"

namespace intrinsic::icon {

bool MemorySegment::IsValid() const { return value_ != nullptr; }

std::string MemorySegment::Name() const { return name_; }

const SegmentHeader& MemorySegment::Header() const {
  // The header will stay valid throughout the lifetime of the memory segment,
  // so it's safe to reference it.
  return *header_;
}

absl::StatusOr<MemorySegment::SegmentDescriptor> MemorySegment::Get(
    const SegmentNameToFileDescriptorMap& segment_name_to_file_descriptor_map,
    absl::string_view name) {
  int shm_fd = -1;
  if (auto it = segment_name_to_file_descriptor_map.find(name);
      it != segment_name_to_file_descriptor_map.end()) {
    shm_fd = it->second;
  } else {
    return absl::NotFoundError(
        absl::StrCat("No file descriptor found for segment: ", name,
                     ". Available segments: ",
                     absl::StrJoin(segment_name_to_file_descriptor_map, ", ",
                                   absl::PairFormatter("="))));
  }

  if (shm_fd == -1) {
    return absl::InternalError(absl::StrCat(
        "Invalid file descriptor for shared memory segment: ", name, "."));
  }

  struct stat shared_memory_stats;
  if (fstat(shm_fd, &shared_memory_stats) != 0) {
    // Return an error and forward errno
    return absl::InternalError(
        absl::StrCat("Failed to read size of segment '", name,
                     "'. 'fstat' failed with:", strerror(errno)));
  }
  SegmentDescriptor segment_info;

  segment_info.size = shared_memory_stats.st_size;

  // The segment needs to be at least the size of a SegmentHeader!
  if (segment_info.size <= sizeof(SegmentHeader)) {
    return absl::InternalError(
        absl::StrCat("Shared memory segment ", name,
                     " must be bigger than the SegmentHeader."));
  }

  // Note: This mapping survives closing the file descriptor.
  segment_info.segment_start = static_cast<uint8_t*>(
      mmap(nullptr, segment_info.size, PROT_WRITE | PROT_READ,
           MAP_SHARED | MAP_LOCKED, shm_fd, 0));
  if (segment_info.segment_start == nullptr ||
      segment_info.segment_start == MAP_FAILED) {
    return absl::InternalError(
        absl::StrCat("Unable to map shared memory segment: ", name, " [",
                     strerror(errno), "]"));
  }

  // Additionally locking the pages as recommended by
  // https://man7.org/linux/man-pages/man2/mmap.2.html, because major faults are
  // not acceptable after the initialization of the mapping.
  if (mlock(/*__addr=*/segment_info.segment_start,
            /*__len=*/segment_info.size) != 0) {
    return absl::InternalError(
        absl::StrCat("Unable to mlock shared memory segment \"", name,
                     "\" with error: ", strerror(errno), "."));
  }

  return segment_info;
}

MemorySegment::~MemorySegment() { CleanUpSharedMemory(); }

MemorySegment::MemorySegment(absl::string_view name, SegmentDescriptor segment,
                             MemorySegment::ReadWriteKind kind)
    : name_(name),
      header_(reinterpret_cast<SegmentHeader*>(segment.segment_start)),
      value_(segment.segment_start + sizeof(SegmentHeader)),
      size_(segment.size),
      read_write_kind_(kind) {
  if (header_ == nullptr) {
    return;
  }
  switch (read_write_kind_) {
    case ReadWriteKind::kReadOnly:
      header_->IncrementReaderRefCount();
      break;
    case ReadWriteKind::kReadWrite:
      header_->IncrementWriterRefCount();
      break;
    default:
      break;
  }
}

MemorySegment::MemorySegment(MemorySegment&& other) noexcept
    : name_(std::exchange(other.name_, "")),
      header_(std::exchange(other.header_, nullptr)),
      value_(std::exchange(other.value_, nullptr)),
      size_(std::exchange(other.size_, 0)),
      read_write_kind_(
          std::exchange(other.read_write_kind_, ReadWriteKind::kUnknown)) {}

MemorySegment& MemorySegment::operator=(MemorySegment&& other) noexcept {
  name_ = std::exchange(other.name_, "");
  CleanUpSharedMemory();
  header_ = std::exchange(other.header_, nullptr);
  value_ = std::exchange(other.value_, nullptr);
  size_ = std::exchange(other.size_, 0);
  read_write_kind_ =
      std::exchange(other.read_write_kind_, ReadWriteKind::kUnknown);
  return *this;
}

SegmentHeader* MemorySegment::HeaderPointer() { return header_; }

uint8_t* MemorySegment::Value() { return value_; }
const uint8_t* MemorySegment::Value() const { return value_; }

size_t MemorySegment::ValueSize() const {
  if (!IsValid() || (sizeof(SegmentHeader) > size_)) {
    return 0;
  }
  return size_ - sizeof(SegmentHeader);
}

void MemorySegment::CleanUpSharedMemory() noexcept {
  if (header_ != nullptr) {
    // We're about to drop our old header pointer, decrement its
    // reader/writer count accordingly.
    switch (read_write_kind_) {
      case MemorySegment::ReadWriteKind::kReadWrite:
        header_->DecrementWriterRefCount();
        break;
      case MemorySegment::ReadWriteKind::kReadOnly:
        header_->DecrementReaderRefCount();
        break;
      default:
        break;
    }
    // Also, unmap the memory. Since MemorySegment is move-only, nothing else
    // should be using this particular pointer.
    //
    // This automatically releases the mlock on that memory too.
    if (munmap(header_, size_) == -1) {
      LOG(WARNING) << "Failed to unmap memory for '" << name_
                   << "'. with error: " << strerror(errno)
                   << ". Continuing anyways.";
    }
  }
}

}  // namespace intrinsic::icon
