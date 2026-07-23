// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_WORLD_GZFILE_GZFILE_H_
#define INTRINSIC_WORLD_GZFILE_GZFILE_H_

#include <cstdint>
#include <cstdio>
#include <memory>
#include <optional>
#include <set>
#include <string>
#include <utility>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/types/optional.h"
#include "intrinsic/world/gzfile/chunk_entry.h"
#include "ortools/base/strong_int.h"
#include "zip.h"
#include "zipconf.h"

namespace intrinsic {

// ChunkId in little endian order. The bytes should be printable characters. It
// represents the TYPE of the data within the chunk.
DEFINE_STRONG_INT_TYPE(ChunkId, uint32_t);

// ChunkKey is an arbitrary string.
using ChunkKey = std::string;

// GZFile is an internal format to store arbitrary contents.
//
// WARNING: Users should not rely on the stability of the format and it can
// change without notice. There is no backwards compatibility guarantee.
class GZFile {
 public:
  ~GZFile();

  // Opens the file for read/write and leave the contents as is.
  static absl::StatusOr<std::unique_ptr<GZFile>> Open(absl::string_view path);

  // Opens the file for write and ignores the on disk data (if any).
  static absl::StatusOr<std::unique_ptr<GZFile>> Create(absl::string_view path);

  // Gets a ChunkEntry corresponding to the chunk_id if it exists. Uses a
  // default key of 0.
  std::optional<ChunkEntry> GetChunk(ChunkId chunk_id) const;

  // Gets a ChunkEntry corresponding to the chunk_id if it exists.
  std::optional<ChunkEntry> GetChunk(ChunkId chunk_id, ChunkKey key) const;

  // Gets an existing ChunkEntry corresponding to the chunk_id or dies if
  // it does not exist. Uses a default key of 0.
  ChunkEntry GetChunkOrDie(ChunkId chunk_id) const;

  // Gets an existing ChunkEntry corresponding to the chunk_id or dies if
  // it does not exist
  ChunkEntry GetChunkOrDie(ChunkId chunk_id, ChunkKey key) const;

  // Returns true if the chunk id is part of this file. Uses a default key of 0.
  bool HasChunk(ChunkId chunk_id) const;

  // Returns true if the chunk id is part of this file.
  bool HasChunk(ChunkId chunk_id, ChunkKey key) const;

  // Returns the set of chunk-ids contained in this file.
  std::set<std::pair<ChunkId, ChunkKey>> GetAllChunkIds() const;

  // Returns the set of keys contained in this file for a given chunk.
  std::set<ChunkKey> GetAllChunkKeys(ChunkId chunk_id) const;

  // Update the chunk data for one of the chunks. Uses a default key of 0.
  absl::Status SetChunk(ChunkId chunk_id, const ChunkEntry& data);

  // Update the chunk data for one of the chunks.
  absl::Status SetChunk(ChunkId chunk_id, ChunkKey key, const ChunkEntry& data);

  // Flush any updates that were specified with SetChunk but not yet written.
  // This should be called before ~GZFile() because, although the destructor
  // will flush the updates, errors when writing the output file may not be
  // noticed.
  absl::Status Flush();

  // Clears updates that have not been flushed yet.
  void ClearUpdates();

 private:
  explicit GZFile(absl::string_view archive_path)
      : archive_path_(archive_path) {}

  absl::Status OpenImpl(int flags);
  absl::Status Close();

  std::optional<ChunkEntry> GetChunkByIx(zip_int64_t ix) const;
  std::optional<ChunkEntry> GetChunkById(const std::string& name) const;

  struct FileZipSourceContext {
    zip_error_t error;
    zip_int64_t supports;

    // reading
    std::string fname;
    FILE* fin;
    struct zip_stat stat;
    zip_error_t stat_error;

    // writing
    std::string tmpname;
    FILE* fout;
  };

  static zip_int64_t CreateTempOutput(FileZipSourceContext* ctx);
  static zip_int64_t FileZipSourceCallback(void* userdata, void* data,
                                           zip_uint64_t len,
                                           zip_source_cmd_t cmd);

  std::unique_ptr<FileZipSourceContext> context_;
  zip_source_t* zip_source_ = nullptr;
  zip_t* archive_ = nullptr;
  std::string archive_path_;
};

inline bool operator<(const ChunkId& lhs, const ChunkId& rhs) {
  return lhs.value() < rhs.value();
}

}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_GZFILE_GZFILE_H_
