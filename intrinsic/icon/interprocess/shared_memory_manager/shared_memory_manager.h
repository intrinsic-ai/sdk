// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_INTERPROCESS_SHARED_MEMORY_MANAGER_SHARED_MEMORY_MANAGER_H_
#define INTRINSIC_ICON_INTERPROCESS_SHARED_MEMORY_MANAGER_SHARED_MEMORY_MANAGER_H_

#include <stddef.h>
#include <stdint.h>

#include <memory>
#include <string>
#include <typeinfo>
#include <utility>
#include <vector>

#include "absl/container/flat_hash_map.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/domain_socket_utils.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/memory_segment.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/segment_header.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/segment_info.fbs.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::icon {

// A type `T` is suited for shared memory if it's trivially copyable (no heap
// allocation internally) and is not a pointer type.
template <class T>
inline void AssertSharedMemoryCompatibility() {
  static_assert(
      std::is_trivially_copyable_v<T>,
      "only trivially copyable data types are supported as shm segments");
  static_assert(!std::is_pointer_v<T>,
                "pointer types are not supported as shm segments");
}

// The `SharedMemoryManager` creates and administers a set of anonymous shared
// memory segments.
//
// Creates segments as anonymous files using `memfd_create` (see
// https://man7.org/linux/man-pages/man2/memfd_create.2.html).
//
// Each allocated segmented is prefixed with a `SegmentHeader` to store some
// meta information about the allocated segment such as a reference counting.
// The overall data layout of each segment looks thus like the following:
//
// [SegmentHeader][Payload T]
// ^              ^
// Header()       Value()
//
// The manager additionally maintains a map of all allocated segments for
// further introspection of the segments. Once the manager goes out of scope, it
// unlinks all allocated memory; The kernel then eventually deletes the shared
// memory files once there's no further process using them.
// Once a segment is added via `AddSegment` it is fully initialized with a
// default value or any given value.
class SharedMemoryManager final {
 public:
  struct MemorySegmentInfo {
    uint8_t* data = nullptr;
    // Required for unmapping the segment
    size_t length = 0;
    // A value of true indicates that this segment needs to be used
    // by ICON.
    bool must_be_used = false;
    // The file descriptor of the anonymous memory.
    int fd = -1;
  };

  // Creates a new `SharedMemoryManager`.
  // Returns an error of the module name is empty.
  // Returns unique_ptr because that allows taking stable references e.g. for
  // intrinsic/icon/hal/hardware_interface_registry.h.
  static absl::StatusOr<std::unique_ptr<SharedMemoryManager>> Create(
      absl::string_view shared_memory_namespace, absl::string_view module_name);

  SharedMemoryManager() = delete;

  // This class is move-only.
  SharedMemoryManager(const SharedMemoryManager& other) = delete;
  SharedMemoryManager& operator=(const SharedMemoryManager& other) = delete;
  // We need to clear other.memory_segments_ on moves in order to avoid
  // use-after-move bugs when accessing memory_segments_ in
  // ~SharedMemoryManager.
  SharedMemoryManager(SharedMemoryManager&& other) noexcept
      : memory_segments_(std::move(other.memory_segments_)) {
    other.memory_segments_.clear();
  }
  SharedMemoryManager& operator=(SharedMemoryManager&& other) noexcept {
    memory_segments_ = std::move(other.memory_segments_);
    other.memory_segments_.clear();
    return *this;
  }
  // Closes all shared memory segments.
  ~SharedMemoryManager();

  // Provides  access to the shared memory location specified by `segment_name`.
  // Returns NotFoundError if no such segment has been added.
  // Forwards mapping errors.
  template <class MemorySegmentT>
  absl::StatusOr<MemorySegmentT> Get(absl::string_view segment_name) const {
    static_assert(
        std::is_base_of_v<::intrinsic::icon::MemorySegment, MemorySegmentT>,
        "Template parameter for SharedMemoryManager::Get() must inherit from "
        "::intrinsic::icon::MemorySegment");
    return MemorySegmentT::Get(segment_name_to_file_descriptor_map_,
                               segment_name);
  }

  // Reference to the internal map of segment names to file descriptors.
  // Contains the names and file descriptors of all segments that are currently
  // registered with the SharedMemoryManager.
  // For tests where no HardwareModule Proxy is used.
  const SegmentNameToFileDescriptorMap& SegmentNameToFileDescriptorMap() const {
    return segment_name_to_file_descriptor_map_;
  }

  // Allocates a shared memory segment for the type `T` and initializes it with
  // the default value of `T`.
  // The type must be trivially copyable and not a pointer type; other types
  // fail to compile.
  // The name for the segment should be POSIX conform, in which the length is
  // not to exceed 255 characters.
  // The value of `must_be_used` indicates whether this segment needs to be
  // used by ICON.
  // Similarly, one can optionally pass in a type identifier string to uniquely
  // describe the type of the data segment. The string can't exceed a max length
  // of `SegmentHeader::TypeInfo::kMaxSize` and defaults to a compiler generated
  // typeid. Please note that the compiler generated default is not defined by
  // the C++ standard and thus may not conform across process boundaries with
  // different compilers.
  // Returns `absl::InvalidArgumentError` if the name is not valid.
  // Returns `absl::AlreadyExistsError` if the shared memory segment
  // with this name already exists
  // Returns `absl::InternalError` if the underlying POSIX call fails.
  // Returns `absl::OkStatus` is the shared memory segment was successfully
  // allocated.
  template <class T>
  absl::Status AddSegmentWithDefaultValue(absl::string_view name,
                                          bool must_be_used) {
    return AddSegmentWithDefaultValue<T>(name, must_be_used, typeid(T).name());
  }
  template <class T>
  absl::Status AddSegmentWithDefaultValue(absl::string_view name,
                                          bool must_be_used,
                                          const std::string& type_id) {
    AssertSharedMemoryCompatibility<T>();
    INTR_RETURN_IF_ERROR(InitSegment(name, must_be_used, sizeof(T), type_id));
    return SetSegmentValue(name, T());
  }

  // Allocates a shared memory segment for the type `T` and initializes it with
  // the specified value of `T`.
  // Besides the initialized value for the segment, this function behaves
  // exactly like `AddSegment` above.
  template <class T>
  absl::Status AddSegment(absl::string_view name, bool must_be_used,
                          const T& value) {
    return AddSegment<T>(name, must_be_used, value, typeid(T).name());
  }
  template <class T>
  absl::Status AddSegment(absl::string_view name, bool must_be_used,
                          const T& value, const std::string& type_id) {
    AssertSharedMemoryCompatibility<T>();
    INTR_RETURN_IF_ERROR(InitSegment(name, must_be_used, sizeof(T), type_id));
    return SetSegmentValue(name, value);
  }
  template <class T>
  absl::Status AddSegment(absl::string_view name, bool must_be_used,
                          T&& value) {
    return AddSegment<T>(name, must_be_used, std::forward<T>(value),
                         typeid(T).name());
  }
  template <class T>
  absl::Status AddSegment(absl::string_view name, bool must_be_used, T&& value,
                          const std::string& type_id) {
    INTR_RETURN_IF_ERROR(InitSegment(name, must_be_used, sizeof(T), type_id));
    return SetSegmentValue(name, std::forward<T>(value));
  }

  // Allocates a generic memory segment with a byte (uint8_t) array payload of
  // `payload_size` bytes.
  absl::Status AddSegment(absl::string_view name, bool must_be_used,
                          size_t payload_size) {
    return AddSegment(name, must_be_used, payload_size, typeid(uint8_t).name());
  }
  // Allocates a memory segment of type `type_id` with a payload of
  // `payload_size` bytes.
  absl::Status AddSegment(absl::string_view name, bool must_be_used,
                          size_t payload_size, const std::string& type_id) {
    INTR_RETURN_IF_ERROR(
        InitSegment(name, must_be_used, payload_size, type_id));
    return absl::OkStatus();
  }

  // Returns the `SegmentHeader` belonging to the shared memory segment
  // specified by the given name.
  // Returns null pointer if the segment with the given name does not exist.
  const SegmentHeader* GetSegmentHeader(absl::string_view name);

  // Returns the value belonging to the shared memory segment specified by the
  // given name.
  // Returns `nullptr` if the segment with the given name does
  // not exist.
  // Note, the type `T` has to match the type with which the segment was
  // originally created. This function leads to undefined behavior otherwise.
  template <class T>
  const T* GetSegmentValue(absl::string_view name) {
    return reinterpret_cast<T*>(GetRawValue(name));
  }

  // Copies the new value into an existing shared memory segment.
  // Returns `absl::NotFoundError` if the segment with the given name does
  // not exist.
  // Note, the type `T` has to match the type with which the segment was
  // originally created. This function leads to undefined behavior otherwise.
  template <class T>
  absl::Status SetSegmentValue(absl::string_view name, const T& new_value) {
    uint8_t* value = GetRawValue(name);
    if (value == nullptr) {
      return absl::NotFoundError(
          absl::StrCat("memory segment not found: ", name));
    }
    *reinterpret_cast<T*>(value) = new_value;
    return absl::OkStatus();
  }
  template <class T>
  absl::Status SetSegmentValue(absl::string_view name, T&& new_value) {
    uint8_t* value = GetRawValue(name);
    if (value == nullptr) {
      return absl::NotFoundError(
          absl::StrCat("memory segment not found: ", name));
    }
    *reinterpret_cast<T*>(value) = std::forward<T>(new_value);
    return absl::OkStatus();
  }

  // Returns a pointer to the untyped payload in the shared memory segment.
  // Memory layout is described in
  // intrinsic/icon/interprocess/shared_memory_manager/segment_header.h
  // This function might be used when access to the underlying generic memory
  // location is needed, e.g. via `std::memcpy`. One typical use case is to copy
  // a flatbuffer (or any other serialized data struct) into a shared memory
  // segment. Prefer accessing the values via `GetSegmentValue` or
  // `SetSegmentValue` for type safety.
  uint8_t* GetRawValue(absl::string_view name);

  // Returns a list of names for all registered shared memory segments.
  std::vector<std::string> GetRegisteredMemoryNames() const;

  // Returns a SegmentInfo struct containing the list of registered memory
  // segments.
  intrinsic_fbs::SegmentInfo GetSegmentInfo() const;

  // Name of the module owning this SharedMemoryManager.
  std::string ModuleName() const;

  // Namespace for the shared memory interfaces using this SharedMemoryManager.
  std::string SharedMemoryNamespace() const;

 private:
  explicit SharedMemoryManager(absl::string_view module_name,
                               absl::string_view shared_memory_namespace);

  // Creates an anonymous shared memory segment with the given name and the size
  // of SegmentHeader + `payload_size`.
  // The SegmentHeader is initialized with the given `type_id`.
  // Returns `absl::InternalError` if any of the underlying
  // POSIX calls fail.
  absl::Status InitSegment(absl::string_view name, bool must_be_used,
                           size_t payload_size, const std::string& type_id);

  // Returns a pointer to the start of the memory segment.
  // The SegmentHeader of this segment lives at the address this pointer
  // indicates.
  // Returns nullptr if the segment with the given name does not exist.
  uint8_t* GetRawSegment(absl::string_view name);

  // Can be generated from memory_segments_, but it's more efficient to simply
  // update the map when a new segment is added.
  icon::SegmentNameToFileDescriptorMap segment_name_to_file_descriptor_map_;

  // We not only store the name of the each initialized segment, but also a
  // pointer to its allocated memory. That way we can later on provide
  // introspection tools around all allocated memory in the system.
  absl::flat_hash_map<std::string, MemorySegmentInfo> memory_segments_;
  std::string module_name_;
  std::string shared_memory_namespace_;
};

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_INTERPROCESS_SHARED_MEMORY_MANAGER_SHARED_MEMORY_MANAGER_H_
