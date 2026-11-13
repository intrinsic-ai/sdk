<!--
Copyright 2026 Intrinsic Innovation LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->
# Middleware Aligned Buffers

Robotics payloads such as camera frames, point clouds, and tensors need 256-byte
memory alignment for SIMD and zero-copy interop (such as DLPack and `mmap`
shared memory), and they are often too large to copy between middleware layers.
Standard containers like `std::vector<uint8_t>` or `std::string` do not
guarantee 256-byte alignment, and copying them duplicates the underlying bytes.

This package provides `MutableBuffer`, `ImmutableBuffer`, and factory helpers in
`namespace intrinsic` so producers can build payloads once and share them across
threads or subsystems without copying bytes or exposing `<zenoh.h>`.

## Two-Stage Lifecycle: Write Exclusively, Then Share

Buffer ownership follows a write-once, read-many flow:

1.   Allocate a writable buffer with `AllocateMutableBuffer(size)`. Because
    `MutableBuffer` is move-only, only one writer owns and mutates its bytes at
    a time.
2.   Fill the writable range (`data()` through `data() + size()`, or convert the
    buffer directly to `std::span<uint8_t>`).
3.   Call `FreezeMutableBuffer(std::move(buf))` to turn the writable buffer into
    an `ImmutableBuffer` in O(1) without copying any bytes. Freezing moves
    ownership out of the mutable buffer, leaving the original handle invalid.
4.   Copy the immutable buffer across threads or call `Slice(offset, length)` to
    pass read-only sub-views. Every copy and slice shares the same backing
    memory via an internal atomic reference count until the last handle drops.

If you already have bytes in memory and want a read-only buffer in one step,
call `CopyToImmutableBuffer(data, size)`.

```cpp
absl::StatusOr<intrinsic::MutableBuffer> mutable_buf =
    intrinsic::AllocateMutableBuffer(1024);
if (!mutable_buf.ok()) {
  return mutable_buf.status();
}
std::memset(mutable_buf->data(), 0xAB, mutable_buf->size());

absl::StatusOr<intrinsic::ImmutableBuffer> shared_buf =
    intrinsic::FreezeMutableBuffer(std::move(*mutable_buf));
if (!shared_buf.ok()) {
  return shared_buf.status();
}
// *mutable_buf is now invalid (`mutable_buf->valid() == false`).

intrinsic::ImmutableBuffer header = shared_buf->Slice(0, 64);
intrinsic::ImmutableBuffer payload = shared_buf->Slice(64, 960);
```

## Key Guarantees and Edge Cases

*   Valid zero-byte buffers still have a non-null `data()` pointer. Calling
    `AllocateMutableBuffer(0)` returns a valid, empty buffer (`valid() == true`,
    `empty() == true`) whose data pointer is non-null and 256-byte aligned, so
    callers passing `(buf.data(), buf.size())` to C or SIMD APIs never trip
    null-pointer checks on empty payloads.
*   Default-constructed, moved-from, and out-of-range slice handles are invalid.
    An invalid buffer reports `valid() == false`, `data() == nullptr`,
    `size() == 0`, and `alignment() == 0`.
*   Slicing checks bounds without integer overflow. `Slice(offset, length)`
    returns an invalid buffer if the source is invalid, `offset > size()`, or
    `length > size() - offset`.
*   Root buffers always align to 256 bytes (`kAlignment`), while slices report
    their effective byte alignment capped at 256 bytes. Slicing at a multiple of
    256 preserves 256-byte alignment; slicing at offset 4 reports 4-byte
    alignment.
*   A mutable buffer requires external synchronization if touched by multiple
    threads. Separate `ImmutableBuffer` copies or slices sharing the same
    allocation can be read, copied, and destroyed concurrently from any thread.

## How Ownership Works Internally

Both buffer classes hold a data pointer, a byte size, and a pointer to a shared
heap-allocated `buffer_internal::BufferControlBlock` base in
`buffer_internal.h`. Concrete backing implementations (such as the host heap
control block in `heap_buffer_allocator.cc`) subclass `BufferControlBlock` and
override its virtual destructor. While a buffer is mutable, its control block
has a reference count of 1, and `FreezeMutableBuffer` moves that control-block
pointer straight into the new immutable buffer without allocating another
control block.

Copying or slicing an immutable buffer increments the atomic reference count
(`memory_order_relaxed`), and destroying each handle decrements it
(`memory_order_acq_rel`) and invokes the derived destructor to free the backing
storage when the count drops to 0.

## Files in This Directory

*   `mutable_buffer.h`, `mutable_buffer.cc`: Move-only writable buffer class
    (`:mutable_buffer`).
*   `immutable_buffer.h`, `immutable_buffer.cc`: Reference-counted read-only
    buffer class with zero-copy `Slice()` views (`:immutable_buffer`).
*   `buffer_factory.h`, `buffer_factory.cc`: Public `AllocateMutableBuffer`,
    `CopyToImmutableBuffer`, and `FreezeMutableBuffer` helpers
    (`:buffer_factory`).
*   `buffer_internal.h`: Middleware-private header defining the polymorphic
    `BufferControlBlock` base and `PrivateAccessor` (`:buffer_internal`).
*   `buffer_test.cc`: 13 unit tests covering alignment, empty buffers, moves,
    copies, slicing, freeze transfers, polymorphic control-block destruction,
    and multi-threaded readers (`:buffer_test`).
*   `allocators/`: Backing-storage allocator interface and heap allocator. See
    [the allocators README](allocators/README.md).

## Building and Testing

Run all 16 buffer container and allocator unit tests with:

```bash
bazel test //intrinsic/middleware/buffers/...
```
