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

#ifndef INTRINSIC_WORLD_GZFILE_CHUNK_ENTRY_H_
#define INTRINSIC_WORLD_GZFILE_CHUNK_ENTRY_H_

#include <cstdint>
#include <string>
#include <utility>

#include "absl/strings/string_view.h"

namespace intrinsic {

enum class CompressionMode {
  kStore = 0,       /* disable compression, eg if already compressed */
  kCompressDefault, /* use default compression */
};

class ChunkEntry {
 public:
  ChunkEntry(
      uint32_t version, absl::string_view data,
      CompressionMode compression_mode = CompressionMode::kCompressDefault)
      : data_(std::string(data)),
        data_version_(version),
        compression_mode_(compression_mode) {}
  ChunkEntry(
      uint32_t version, std::string&& data,
      CompressionMode compression_mode = CompressionMode::kCompressDefault)
      : data_(std::move(data)),
        data_version_(version),
        compression_mode_(compression_mode) {}
  ChunkEntry(
      uint32_t version, const char* data,
      CompressionMode compression_mode = CompressionMode::kCompressDefault)
      : ChunkEntry(version, std::string(data), compression_mode) {}

  // Returns the version number associated with this chunk data.
  uint32_t GetDataVersion() const { return data_version_; }

  // Returns the data contained in this chunk after decompression.
  absl::string_view GetUncompressedData() const { return data_; }

  CompressionMode GetCompressionMode() const { return compression_mode_; }

 private:
  std::string data_;
  uint32_t data_version_;
  CompressionMode compression_mode_;
};

}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_GZFILE_CHUNK_ENTRY_H_
