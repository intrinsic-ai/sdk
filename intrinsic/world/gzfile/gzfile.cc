// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/world/gzfile/gzfile.h"

#include <cerrno>
#include <cstdint>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <filesystem>  // NOLINT
#include <memory>
#include <optional>
#include <set>
#include <string>
#include <string_view>
#include <system_error>  // NOLINT
#include <utility>
#include <vector>

#include "absl/flags/flag.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/ascii.h"
#include "absl/strings/numbers.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_format.h"
#include "absl/strings/str_split.h"
#include "absl/strings/string_view.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/world/gzfile/chunk_entry.h"
#include "ortools/base/filesystem.h"
#include "ortools/base/options.h"
#include "ortools/base/path.h"
#include "zip.h"
#include "zipconf.h"

ABSL_FLAG(std::string, gzfile_tmpdir, "",
          "Directory for temporary files. Default: next to output file.");

/* The implementation uses libzip. Generated files can be inspected with
 * `zipinfo -v` but also used like a regular zip file (eg to extract a chunk for
 * debugging). Zip files allow individual entries to be stored using different
 * compressors. It is also possible to possible to extract them, or add new
 * ones, without applying compression or decompression to the entire archive.
 * See https://en.wikipedia.org/wiki/Zip_(file_format)#Design
 *
 * For this file format, we're storing the file version into the archive comment
 * and chunk versions into the file comments.
 * (e.g GZFF 1 for archive comment).
 *
 * TODO(ensonic): The zip library does not support empty zips files. We could
 * work around this by adding a marker chunk when creating new files, but hiding
 * this chunk from other operations, or make it an error to write an empty file.
 *
 * TODO(ensonic): By using the zip library all contained data need to fit in
 * memory. We need to strategically call Flush() to keep memory usage under
 * control. The penalty for calling Flush() is that when writing the 2nd batch
 * the previous data will be copied to a new file in order to do a safe atomic
 * write-and-rename.
 */

namespace intrinsic {
namespace {

// The version of the file header we currently support "natively"
constexpr uint32_t kCurrentVersion = 1;
// The fastest compression level, as defined by zip_set_file_compression.
constexpr uint32_t kFastestCompression = 1;

std::string ZipErrorStrError(zip_error_t* error) {
  const char* s = zip_error_strerror(error);
  std::string res(s);
  // s is shallow copied into error->s during zip_error_strerror. zip_error_fini
  // frees it and set error->s to null. This makes sure that there is no
  // unwrapped malloc buffer hanging around.
  zip_error_fini(error);
  return res;
}

absl::Status ZipErrorToStatus(absl::string_view name, zip_error_t* error) {
  std::string message = absl::StrCat(name, ": ", ZipErrorStrError(error));
  if (error->zip_err == ZIP_ER_NOENT) {
    return absl::NotFoundError(message);
  } else {
    return absl::InternalError(message);
  }
}

// The chunk_id is a human readable identifier.
std::string ChunkIdToName(ChunkId chunk_id, ChunkKey key) {
  uint32_t id = chunk_id.value();
  std::string chunk_id_str{static_cast<char>((id >> 24) & 0xff),
                           static_cast<char>((id >> 16) & 0xff),
                           static_cast<char>((id >> 8) & 0xff),
                           static_cast<char>(id & 0xff)};

  return absl::StrCat(chunk_id_str, "/", key);
}

std::optional<std::pair<ChunkId, ChunkKey>> NameToChunkId(const char* name) {
  if (!name || strlen(name) < 6) {
    LOG(WARNING) << "Can't make intrinsic chunk from '" << name << "'";
    return std::nullopt;
  }

  std::vector<absl::string_view> split_results = absl::StrSplit(name, '/');
  if (split_results.size() != 2) {
    LOG(WARNING) << "Can't make intrinsic chunk from '" << name << "'";
    return std::nullopt;
  }

  const auto& chunk_id_str = split_results[0];
  if (chunk_id_str.length() != 4) {
    LOG(WARNING) << "Can't make intrinsic chunk from '" << name << "'";
    return std::nullopt;
  }

  ChunkKey key(split_results[1]);

  return std::make_pair(ChunkId(static_cast<uint32_t>(chunk_id_str[0]) << 24 |
                                static_cast<uint32_t>(chunk_id_str[1]) << 16 |
                                static_cast<uint32_t>(chunk_id_str[2]) << 8 |
                                static_cast<uint32_t>(chunk_id_str[3])),
                        key);
}

absl::Status ChunkIdIsValid(ChunkId chunk_id) {
  uint32_t id = chunk_id.value();
  for (int i = 0; i < 4; i++) {
    if (!absl::ascii_isprint(id & 0xff)) {
      return absl::InvalidArgumentError(
          absl::StrCat("Non printable char at position ", (3 - i)));
    }
    id >>= 8;
  }

  return absl::OkStatus();
}

int StatusToErrno(absl::Status status) {
  switch (status.code()) {
    case absl::StatusCode::kOk:
      return 0;
    case absl::StatusCode::kInvalidArgument:
      return EINVAL;
    case absl::StatusCode::kNotFound:
      return ENOENT;
    case absl::StatusCode::kPermissionDenied:
      return EACCES;
    case absl::StatusCode::kUnauthenticated:
      return ENOKEY;
    case absl::StatusCode::kDeadlineExceeded:
      return ETIMEDOUT;
    case absl::StatusCode::kUnimplemented:
      return ENOSYS;
    case absl::StatusCode::kOutOfRange:
      return ERANGE;
    default:
      return EIO;
  }
}

int FileSeek(FILE* file, std::string_view filename, zip_int64_t offset,
             int whence, zip_error_t* error) {
  switch (whence) {
    case SEEK_SET:
    case SEEK_CUR:
    case SEEK_END: {
      if (int err = fseek(file, offset, whence); err != 0) {
        zip_error_set(error, ZIP_ER_SEEK, err);
        return -1;
      }
      return 0;
    }
    default: {
      zip_error_set(error, ZIP_ER_INVAL, 0);
      return -1;
    }
  }
}

// Rename a file, either using filesystem::rename, or if that fails (eg because
// the source and destination files are on different filesystems), fall back to
// a copy&remove. This will overwrite the destination file if it exists.
absl::Status RenameWithFallback(absl::string_view src, absl::string_view dst) {
  namespace fs = std::filesystem;
  std::error_code ec;
  fs::rename(src, dst, ec);
  if (!ec) {
    // If we were able to rename the file, we're done.
    return absl::OkStatus();
  }

  // If rename failed, we fall back to a copy+remove.
  if (!fs::copy_file(src, dst, fs::copy_options::overwrite_existing, ec)) {
    return absl::InternalError(ec.message());
  }

  if (!fs::remove(src, ec)) {
    return absl::InternalError(ec.message());
  }
  return absl::OkStatus();
}

}  // namespace

zip_int64_t GZFile::CreateTempOutput(FileZipSourceContext* ctx) {
  std::string tmpdir = absl::GetFlag(FLAGS_gzfile_tmpdir);
  if (tmpdir.empty()) {
    tmpdir = std::string(file::Dirname(ctx->fname));
  }
  auto filename = file::JoinPath(
      tmpdir, absl::StrFormat("%s-%s", file::Basename(ctx->fname),
                              file::Basename(std::tmpnam(nullptr))));
  ctx->tmpname = filename;
  ctx->fout = fopen(ctx->tmpname.c_str(), "w");
  if (ctx->fout == nullptr) {
    zip_error_set(&ctx->error, ZIP_ER_TMPOPEN, errno);
    return -1;
  }
  return 0;
}

// The semantics of commands follows closely with read_file() in
// http://third_party/libzip/lib/zip_source_filep.c
//
// Some error conditions are no longer ignored.
zip_int64_t GZFile::FileZipSourceCallback(void* userdata, void* data,
                                          zip_uint64_t len,
                                          zip_source_cmd_t cmd) {
  FileZipSourceContext* ctx = static_cast<FileZipSourceContext*>(userdata);
  switch (cmd) {
    case ZIP_SOURCE_BEGIN_WRITE: {
      if (ctx->fname.empty()) {
        zip_error_set(&ctx->error, ZIP_ER_OPNOTSUPP, 0);
        return -1;
      }
      return CreateTempOutput(ctx);
    }
    case ZIP_SOURCE_COMMIT_WRITE: {
      if (ctx->fout != nullptr) {
        if (int err = fclose(ctx->fout); err != 0) {
          zip_error_set(&ctx->error, ZIP_ER_WRITE, err);
          return -1;
        }
      }
      ctx->fout = nullptr;
      absl::Status status = RenameWithFallback(ctx->tmpname, ctx->fname);
      if (!status.ok()) {
        zip_error_set(&ctx->error, ZIP_ER_RENAME, StatusToErrno(status));
        return -1;
      }
      ctx->tmpname = "";
      return 0;
    }
    case ZIP_SOURCE_CLOSE: {
      if (!ctx->fname.empty()) {
        if (int err = fclose(ctx->fin); err != 0) {
          zip_error_set(&ctx->error, ZIP_ER_CLOSE, err);
          return -1;
        }
        ctx->fin = nullptr;
      }
      return 0;
    }
    case ZIP_SOURCE_ERROR: {
      return zip_error_to_data(&ctx->error, data, len);
    }
    case ZIP_SOURCE_FREE: {
      ctx->fname = "";
      ctx->tmpname = "";
      if (ctx->fin != nullptr) {
        if (int err = fclose(ctx->fin); err != 0) {
          zip_error_set(&ctx->error, ZIP_ER_CLOSE, err);
          return -1;
        }
      }
      return 0;
    }
    case ZIP_SOURCE_OPEN: {
      if (ctx->fname.empty()) {
        zip_error_set(&ctx->error, ZIP_ER_INVAL, 0);
        return -1;
      }

      std::string filename = ctx->fname;
      if (!file::IsAbsolutePath(filename)) {
        VLOG(1) << "We have a relative path, trying to resolve it.";
        filename = std::filesystem::absolute(filename);
      }
      ctx->fin = fopen(filename.c_str(), "r");
      if (ctx->fin == nullptr) {
        zip_error_set(&ctx->error, ZIP_ER_OPEN, errno);
        return -1;
      }
      return 0;
    }
    case ZIP_SOURCE_READ: {
      if (feof(ctx->fin) != 0) {
        return 0;
      }
      int nbytes_read = fread(data, sizeof(char), len, ctx->fin);
      if (nbytes_read < len) {
        zip_error_set(&ctx->error, ZIP_ER_READ, EIO);
        return -1;
      }
      return nbytes_read;
    }
    case ZIP_SOURCE_REMOVE: {
      absl::Status status = file::Delete(ctx->fname, file::Defaults());
      if (!status.ok()) {
        zip_error_set(&ctx->error, ZIP_ER_REMOVE, StatusToErrno(status));
        return -1;
      }
      return 0;
    }
    case ZIP_SOURCE_ROLLBACK_WRITE: {
      if (ctx->fout) {
        if (int err = fclose(ctx->fout); err != 0) {
          zip_error_set(&ctx->error, ZIP_ER_CLOSE, err);
          return -1;
        }
        ctx->fout = nullptr;
      }
      absl::Status status = file::Delete(ctx->tmpname, file::Defaults());
      if (!status.ok()) {
        zip_error_set(&ctx->error, ZIP_ER_REMOVE, StatusToErrno(status));
        return -1;
      }
      ctx->tmpname = "";
      return 0;
    }
    case ZIP_SOURCE_SEEK: {
      zip_source_args_seek_t* args =
          ZIP_SOURCE_GET_ARGS(zip_source_args_seek_t, data, len, nullptr);
      if (args == nullptr) {
        return -1;
      }
      return FileSeek(ctx->fin, ctx->fname, args->offset, args->whence,
                      &ctx->error);
    }
    case ZIP_SOURCE_SEEK_WRITE: {
      zip_source_args_seek_t* args =
          ZIP_SOURCE_GET_ARGS(zip_source_args_seek_t, data, len, &ctx->error);
      if (args == nullptr) {
        return -1;
      }
      return FileSeek(ctx->fout, ctx->tmpname, args->offset, args->whence,
                      &ctx->error);
    }
    case ZIP_SOURCE_STAT: {
      if (len < sizeof(ctx->stat)) {
        return -1;
      }
      if (zip_error_code_zip(&ctx->stat_error) != 0) {
        zip_error_set(&ctx->error, zip_error_code_zip(&ctx->stat_error),
                      zip_error_code_system(&ctx->stat_error));
        return -1;
      }
      memcpy(data, &ctx->stat, sizeof(ctx->stat));
      return sizeof(zip_stat_t);
    }
    case ZIP_SOURCE_SUPPORTS: {
      return ctx->supports;
    }
    case ZIP_SOURCE_TELL: {
      int64_t pos = ftell(ctx->fin);
      if (pos == -1) {
        zip_error_set(&ctx->error, ZIP_ER_TELL, errno);
        return -1;
      }
      return pos;
    }
    case ZIP_SOURCE_TELL_WRITE: {
      int64_t pos = ftell(ctx->fout);
      if (pos == -1) {
        zip_error_set(&ctx->error, ZIP_ER_TELL, errno);
        return -1;
      }
      return pos;
    }
    case ZIP_SOURCE_WRITE: {
      int64_t nbytes_written = fwrite(data, sizeof(char), len, ctx->fout);
      if (nbytes_written < len) {
        zip_error_set(&ctx->error, ZIP_ER_WRITE, EIO);
        return -1;
      }
      return nbytes_written;
    }
    default: {
      zip_error_set(&ctx->error, ZIP_ER_OPNOTSUPP, 0);
      return -1;
    }
  }
}

absl::Status GZFile::OpenImpl(int flags) {
  context_ = std::make_unique<FileZipSourceContext>();
  context_->fname = archive_path_;
  context_->supports = ZIP_SOURCE_SUPPORTS_WRITABLE;

  std::error_code ec;
  auto file_size = std::filesystem::file_size(archive_path_, ec);
  if (ec) {
    if (ec == std::errc::no_such_file_or_directory) {
      // The particular combination of zip_err and sys_err here is
      // crucial for libzip to handle ZIP_CREATE correctly.
      zip_error_set(&context_->stat_error, ZIP_ER_READ, ENOENT);
    } else {
      return absl::InternalError(ec.message());
    }
  } else {
    context_->stat.valid = ZIP_STAT_SIZE;
    context_->stat.size = file_size;
  }

  zip_error_t error;
  zip_error_init(&error);
  zip_source_ =
      zip_source_function_create(FileZipSourceCallback, context_.get(), &error);
  if (zip_source_ == nullptr) {
    return ZipErrorToStatus("zip_source_function_create", &error);
  }

  archive_ = zip_open_from_source(zip_source_, flags, &error);
  if (archive_ == nullptr) {
    return ZipErrorToStatus("zip_open_from_source", &error);
  }

  return absl::OkStatus();
}

absl::StatusOr<std::unique_ptr<GZFile>> GZFile::Open(absl::string_view path) {
  auto gzfile = std::unique_ptr<GZFile>(new GZFile(path));
  INTR_RETURN_IF_ERROR(gzfile->OpenImpl(/*flags=*/0))
      << " while opening " << path;

  int clen = 0;
  const char* cmt = zip_get_archive_comment(gzfile->archive_, &clen, 0);
  if (clen == 0) {
    return absl::InternalError("Missing file version");
  }
  uint32_t file_version;
  if (!absl::SimpleAtoi(cmt, &file_version)) {
    return absl::InternalError("Unparsable file version");
  }
  if (file_version != kCurrentVersion) {
    return absl::OutOfRangeError(
        absl::StrCat("Unsupported header version: ", file_version));
  }

  return gzfile;
}

absl::StatusOr<std::unique_ptr<GZFile>> GZFile::Create(absl::string_view path) {
  auto gzfile = std::unique_ptr<GZFile>(new GZFile(path));
  INTR_RETURN_IF_ERROR(gzfile->OpenImpl(ZIP_CREATE | ZIP_TRUNCATE));
  return gzfile;
}

GZFile::~GZFile() {
  auto res = Close();
  if (!res.ok()) {
    LOG(WARNING) << "Failed to close the intrinsic file: " << res;
  }
}

absl::Status GZFile::Close() {
  if (archive_ == nullptr) {
    if (zip_source_ != nullptr) {
      zip_source_free(zip_source_);
    }
    return absl::OkStatus();
  }

  auto res = absl::OkStatus();
  int num_entries = zip_get_num_entries(archive_, 0);

  if (num_entries > 0) {
    auto cmt = std::to_string(kCurrentVersion);
    int zerr = zip_set_archive_comment(archive_, cmt.c_str(), cmt.length());
    if (zerr < 0) {
      return absl::InternalError(
          absl::StrCat("Failed to store file version: ",
                       ZipErrorStrError(zip_get_error(archive_))));
    }
  }

  int zerr = zip_close(archive_);
  if (zerr < 0) {
    if (num_entries > 0) {
      // If the zip file is empty, it won't create a file.
      res = absl::InternalError(
          absl::StrCat("Closing the file on disk failed! ",
                       ZipErrorStrError(zip_get_error(archive_))));
    }
    zip_discard(archive_);
  }
  archive_ = nullptr;
  zip_source_ = nullptr;

  return res;
}

std::optional<ChunkEntry> GZFile::GetChunkByIx(zip_int64_t ix) const {
  zip_file_t* chunk = zip_fopen_index(archive_, ix, 0);
  if (!chunk) {
    return std::nullopt;
  }

  zip_uint32_t clen = 0;
  const char* cmt = zip_file_get_comment(archive_, ix, &clen, 0);
  if (clen == 0) {
    LOG(ERROR) << "Missing chunk version";
    zip_fclose(chunk);
    return std::nullopt;
  }
  uint32_t chunk_version;
  if (!absl::SimpleAtoi(cmt, &chunk_version)) {
    LOG(ERROR) << "Unparsable chunk version";
    zip_fclose(chunk);
    return std::nullopt;
  }

  // check size
  zip_stat_t sb;
  zip_stat_init(&sb);
  zip_stat_index(archive_, ix, ZIP_STAT_SIZE, &sb);

  // read data
  std::string buffer(sb.size, '\0');
  zip_int64_t ret = zip_fread(chunk, &buffer[0], sb.size);
  zip_fclose(chunk);
  if (ret < sb.size) {
    LOG(ERROR) << "Failed to read from chunk. Wanted " << sb.size << " got "
               << ret;
    return std::nullopt;
  }

  // create chunk entry and pass data to it
  return ChunkEntry(chunk_version, std::move(buffer));
}

std::optional<ChunkEntry> GZFile::GetChunkById(const std::string& name) const {
  zip_int64_t ix = zip_name_locate(archive_, name.c_str(), 0);
  if (ix < 0) {
    return std::nullopt;
  }
  return GetChunkByIx(ix);
}

std::optional<ChunkEntry> GZFile::GetChunk(ChunkId chunk_id) const {
  return GetChunk(chunk_id, "0");
}

std::optional<ChunkEntry> GZFile::GetChunk(ChunkId chunk_id,
                                           ChunkKey key) const {
  return GetChunkById(ChunkIdToName(chunk_id, key));
}

ChunkEntry GZFile::GetChunkOrDie(ChunkId chunk_id) const {
  return GetChunkOrDie(chunk_id, "0");
}

ChunkEntry GZFile::GetChunkOrDie(ChunkId chunk_id, ChunkKey key) const {
  auto optional_chunk = GetChunk(chunk_id, key);
  if (!optional_chunk.has_value()) {
    LOG(FATAL) << "Couldn't find chunk with id" << chunk_id;
  }
  return optional_chunk.value();
}

bool GZFile::HasChunk(ChunkId chunk_id) const {
  return HasChunk(chunk_id, "0");
}

bool GZFile::HasChunk(ChunkId chunk_id, ChunkKey key) const {
  auto name = ChunkIdToName(chunk_id, key);
  zip_int64_t ix = zip_name_locate(archive_, name.c_str(), 0);
  return ix > -1;
}

std::set<std::pair<ChunkId, ChunkKey>> GZFile::GetAllChunkIds() const {
  zip_int64_t num = zip_get_num_entries(archive_, 0);

  std::set<std::pair<ChunkId, ChunkKey>> chunks;
  for (zip_uint64_t i = 0; i < num; i++) {
    auto optional_chunk_id = NameToChunkId(zip_get_name(archive_, i, 0));
    if (optional_chunk_id.has_value()) {
      chunks.insert(optional_chunk_id.value());
    }
  }

  return chunks;
}

std::set<ChunkKey> GZFile::GetAllChunkKeys(ChunkId chunk_id) const {
  zip_int64_t num = zip_get_num_entries(archive_, 0);

  std::set<ChunkKey> chunks;
  for (zip_uint64_t i = 0; i < num; i++) {
    auto optional_chunk_id = NameToChunkId(zip_get_name(archive_, i, 0));
    if (optional_chunk_id.has_value() && optional_chunk_id->first == chunk_id) {
      chunks.insert(optional_chunk_id.value().second);
    }
  }

  return chunks;
}

absl::Status GZFile::SetChunk(ChunkId chunk_id, const ChunkEntry& data) {
  return SetChunk(chunk_id, "0", data);
}

absl::Status GZFile::SetChunk(ChunkId chunk_id, ChunkKey key,
                              const ChunkEntry& data) {
  auto chunk_id_check = ChunkIdIsValid(chunk_id);
  if (!chunk_id_check.ok()) {
    return absl::InvalidArgumentError(
        absl::StrCat("ChunkId is invalid :", chunk_id_check.ToString()));
  }
  const auto buf = data.GetUncompressedData();
  if (buf.empty()) {
    return absl::InvalidArgumentError("ChunkEntry must not be empty");
  }

  // Meh, this memory need to be retained until we close the zip file!
  char* copy = static_cast<char*>(malloc(sizeof(char) * buf.length()));
  if (copy == nullptr) {
    return absl::InternalError(
        absl::StrCat("Failed to alloc chunk buffer of size ", buf.length()));
  }
  memcpy(copy, buf.data(), buf.length());
  zip_source_t* s = zip_source_buffer(archive_, copy, buf.length(), 1);
  if (s == nullptr) {
    return absl::InternalError("Empty chunk data?");
  }

  auto name = ChunkIdToName(chunk_id, key);
  zip_int64_t ix = zip_file_add(archive_, name.c_str(), s, ZIP_FL_OVERWRITE);
  if (ix < 0) {
    zip_source_free(s);
    return absl::InternalError(absl::StrCat("Failed adding chunk: ", name, ": ",
                                            zip_strerror(archive_)));
  }
  // Nuke timestamps to get deterministic files
  zip_file_set_dostime(archive_, ix, 0, 0, 0);

  // apply compression, either none (store) or fast deflate (default)
  zip_int32_t cm = data.GetCompressionMode() == CompressionMode::kStore
                       ? ZIP_CM_STORE
                       : ZIP_CM_DEFLATE;
  int zerr = zip_set_file_compression(archive_, ix, cm, kFastestCompression);
  if (zerr < 0) {
    return absl::InternalError(
        absl::StrCat("Failed to set compression mode: ",
                     ZipErrorStrError(zip_get_error(archive_))));
  }

  // store version
  auto cmt = std::to_string(data.GetDataVersion());
  zerr = zip_file_set_comment(archive_, ix, cmt.c_str(), cmt.length(), 0);
  if (zerr < 0) {
    return absl::InternalError(
        absl::StrCat("Failed to store chunk version: ",
                     ZipErrorStrError(zip_get_error(archive_))));
  }

  return absl::OkStatus();
}

absl::Status GZFile::Flush() {
  // libzip has no flush method. The archive is written when we close the file,
  // so lets close and reopen.
  auto res = Close();
  if (!res.ok()) {
    return absl::InternalError(
        absl::StrCat("Failed to flush the file: ", res.ToString()));
  }

  return OpenImpl(ZIP_CREATE);
}

void GZFile::ClearUpdates() { zip_unchange_all(archive_); }

}  // namespace intrinsic
