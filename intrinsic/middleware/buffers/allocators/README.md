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
# Middleware Buffer Allocators

`MutableBuffer` and `ImmutableBuffer` separate *how callers hold and share
bytes* from *where those bytes live*. This directory defines `BufferAllocator`,
the interface that allocates 256-byte-aligned backing memory for new buffers,
and `HeapBufferAllocator`, the default implementation backed by host heap
memory.

Separating allocation behind an interface lets factory helpers such as
`AllocateMutableBuffer` and `CopyToImmutableBuffer` default to normal heap
memory today (`HeapBufferAllocator::Instance()`) while allowing future shared
memory or device allocators to plug into the exact same buffer types without
changing caller signatures.

## Why BufferAllocator Is Designed This Way

*   Only trusted in-tree allocators can subclass `BufferAllocator`. Every valid
    buffer depends on an internal reference-counted control block
    (`buffer_internal::BufferControlBlock`) and strict 256-byte alignment.
    Keeping the base constructor private with explicit `friend` declarations
    prevents outside code from subclassing the allocator with unaligned memory
    or broken lifetime hooks.
*   Callers pass allocators by `const BufferAllocator&` rather than owning them
    through heap pointers. Allocators live as stateless or long-lived singletons
    (such as `HeapBufferAllocator::Instance()`), so `~BufferAllocator()` is
    protected and non-virtual to block accidental `delete` calls while keeping
    stateless allocators trivially destructible.
*   Every allocator implementation must be thread-safe so multiple threads can
    call `Allocate(size)` concurrently without external locking.

## How HeapBufferAllocator Works

*   For non-empty requests (`size > 0`), `HeapBufferAllocator::Allocate`
    allocates 256-byte-aligned host memory (`is_shm() == false`) via
    non-throwing aligned `::operator new[]` paired with a matching aligned
    deleter inside a heap-specific `BufferControlBlock` subclass.
*   For zero-byte requests (`size == 0`), it avoids a heap allocation for the
    payload by pointing `data()` at a single shared 256-byte-aligned static byte
    (`kEmptyAlignedByte`), while still creating a control block so the returned
    `MutableBuffer` stays `valid()` and non-null.
*   If either the payload allocation or control-block allocation fails, it
    returns `absl::StatusCode::kResourceExhausted` instead of throwing.

## Files in This Directory

*   `buffer_allocator.h`: Header-only closed abstract base class declaring
    `BufferAllocator::Allocate(size)` (`:buffer_allocator`).
*   `heap_buffer_allocator.h`, `heap_buffer_allocator.cc`: Stateless host heap
    allocator and process-wide `HeapBufferAllocator::Instance()` singleton
    (`:heap_buffer_allocator`).
*   `heap_buffer_allocator_test.cc`: 3 unit tests covering closed base-class
    traits, 256-byte aligned allocations, and zero-byte sentinel sharing
    (`:heap_buffer_allocator_test`).

## Building and Testing

Run the 3 allocator unit tests with:

```bash
bazel test //intrinsic/middleware/buffers/allocators/...
```
