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

#include "intrinsic/geometry/storage/file_system_storage.h"

#include <cstddef>
#include <filesystem>
#include <fstream>
#include <memory>
#include <optional>
#include <string>
#include <system_error>
#include <utility>

#include "absl/functional/function_ref.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/api/geometry_fingerprint.h"
#include "intrinsic/geometry/api/io.h"
#include "intrinsic/geometry/api/renderable.h"
#include "intrinsic/geometry/api/renderable_generation.h"
#include "intrinsic/geometry/proto/v1/exact_geometry.pb.h"
#include "intrinsic/geometry/proto/v1/geometry_storage_refs.pb.h"
#include "intrinsic/geometry/proto/v1/material.pb.h"
#include "intrinsic/geometry/storage/geometry_deserializer.h"
#include "intrinsic/geometry/storage/geometry_library.h"
#include "intrinsic/geometry/storage/geometry_serializer.h"
#include "intrinsic/util/status/status_builder.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::geo {
namespace {

using ExactGeometryProto = ::intrinsic_proto::geometry::v1::ExactGeometry;
using ::intrinsic::DataLossErrorBuilder;
using ::intrinsic::InternalErrorBuilder;
using ::intrinsic::NotFoundErrorBuilder;

// Writes the bytes produced by `content_fn` to `file_path`, creating parent
// directories if absent. Returns OK without invoking `content_fn` or rewriting
// when `file_path` is already a regular file.
absl::Status WriteContentToFile(
    const std::filesystem::path& file_path,
    absl::FunctionRef<absl::StatusOr<std::string>()> content_fn) {
  std::error_code error_code;
  if (std::filesystem::is_regular_file(file_path, error_code)) {
    return absl::OkStatus();
  }
  INTR_ASSIGN_OR_RETURN(const std::string bytes, content_fn());
  if (const std::filesystem::path parent_dir = file_path.parent_path();
      !parent_dir.empty()) {
    std::error_code dir_error;
    std::filesystem::create_directories(parent_dir, dir_error);
    if (dir_error) {
      return InternalErrorBuilder()
             << "Failed to create directory '" << parent_dir.string()
             << "': " << dir_error.message();
    }
  }
  std::ofstream out_stream(file_path, std::ios::binary | std::ios::trunc);
  if (!out_stream.is_open()) {
    return InternalErrorBuilder()
           << "Failed to write file: " << file_path.string();
  }
  if (!bytes.empty()) {
    out_stream.write(bytes.data(), static_cast<std::streamsize>(bytes.size()));
  }
  out_stream.close();
  if (!out_stream) {
    std::error_code cleanup_error;
    std::filesystem::remove(file_path, cleanup_error);
    return InternalErrorBuilder()
           << "Failed to write file: " << file_path.string();
  }
  return absl::OkStatus();
}

// Reads the binary contents of `file_path` referenced by `raw_ref`.
absl::StatusOr<std::string> ReadFileBytes(
    const std::filesystem::path& file_path, absl::string_view raw_ref) {
  std::error_code error_code;
  if (!std::filesystem::exists(file_path, error_code) && !error_code) {
    return NotFoundErrorBuilder() << "Could not find '" << raw_ref
                                  << "' in file system: " << file_path.string();
  }
  if (error_code || !std::filesystem::is_regular_file(file_path, error_code)) {
    return InternalErrorBuilder() << "Could not read '" << raw_ref
                                  << "' in file system: " << file_path.string();
  }
  const size_t byte_size = std::filesystem::file_size(file_path, error_code);
  if (error_code) {
    return InternalErrorBuilder() << "Could not read '" << raw_ref
                                  << "' in file system: " << file_path.string();
  }
  std::ifstream in_stream(file_path, std::ios::binary);
  if (!in_stream.is_open()) {
    return InternalErrorBuilder() << "Could not read '" << raw_ref
                                  << "' in file system: " << file_path.string();
  }
  std::string contents(byte_size, '\0');
  if (!in_stream.read(contents.data(),
                      static_cast<std::streamsize>(contents.size()))) {
    return InternalErrorBuilder() << "Could not read '" << raw_ref
                                  << "' in file system: " << file_path.string();
  }
  return contents;
}

// Reads binary protobuf bytes from `file_path` for `raw_ref` and parses them
// into `ProtoT`.
template <typename ProtoT>
absl::StatusOr<ProtoT> ReadBinaryProto(const std::filesystem::path& file_path,
                                       absl::string_view raw_ref) {
  INTR_ASSIGN_OR_RETURN(const std::string bytes,
                        ReadFileBytes(file_path, raw_ref));
  ProtoT proto;
  if (!proto.ParseFromString(bytes)) {
    return DataLossErrorBuilder()
           << "Failed to parse binary proto for '" << raw_ref
           << "' from: " << file_path.string();
  }
  return proto;
}

class FileSystemGeometryLibrary : public GeometryLibrary,
                                  private GeometryDeserializer,
                                  private GeometrySerializer {
 public:
  explicit FileSystemGeometryLibrary(absl::string_view prefix)
      : prefix_(std::string(prefix)) {}

  GeometrySerializer& Serializer() override { return *this; }
  const GeometryDeserializer& Deserializer() const override { return *this; }

  absl::StatusOr<intrinsic_proto::geometry::v1::GeometryStorageRefs>
  SaveGeometry(const Geometry& geometry,
               const SaveGeometryOptions& options) override {
    if (options.HasCustomKeys()) {
      return absl::UnimplementedError(
          "Custom storage keys are not yet implemented in "
          "FileSystemGeometryLibrary.");
    }
    intrinsic_proto::geometry::v1::GeometryStorageRefs geo_storage_refs;
    INTR_ASSIGN_OR_RETURN(const std::string exact_geo_hash,
                          GenerateFingerprint(geometry.GetExactGeometry()));
    const std::filesystem::path exact_geo_file_path =
        prefix_ / std::filesystem::path(exact_geo_hash);
    INTR_RETURN_IF_ERROR(WriteContentToFile(
        exact_geo_file_path,
        [&geometry]() -> absl::StatusOr<std::string> {
          INTR_ASSIGN_OR_RETURN(const ExactGeometryProto exact_geo_proto,
                                ToProto(geometry.GetExactGeometry()));
          return exact_geo_proto.SerializeAsString();
        }))
        << "Failed to upload exact geometry for hash: " << exact_geo_hash;
    geo_storage_refs.set_exact_geometry_ref(exact_geo_hash);

    // Write out the renderable if we need to keep it.
    if (geometry.KeepRenderableForSerialization()) {
      INTR_ASSIGN_OR_RETURN(auto renderable, GetOrGenerateRenderable(geometry));

      std::string renderable_hash = GenerateFingerprint(*renderable);

      const std::filesystem::path renderable_file_path =
          prefix_ / std::filesystem::path(renderable_hash);
      INTR_RETURN_IF_ERROR(WriteContentToFile(
          renderable_file_path,
          [&renderable] { return renderable->GetGLBString(); }))
          << "Failed to upload renderable for hash: " << renderable_hash;
      geo_storage_refs.set_renderable_ref(renderable_hash);
    }
    geo_storage_refs.set_keep_renderable(
        geometry.KeepRenderableForSerialization());
    return geo_storage_refs;
  }

  absl::StatusOr<Geometry> GetGeometry(
      const intrinsic_proto::geometry::v1::GeometryStorageRefs&
          geo_storage_refs,
      std::optional<intrinsic_proto::geometry::v1::MaterialProperties>
          material_properties) const override {
    const std::filesystem::path exact_geo_path =
        prefix_ / std::filesystem::path(geo_storage_refs.exact_geometry_ref())
                      .relative_path();
    INTR_ASSIGN_OR_RETURN(
        const ExactGeometryProto exact_geometry_proto,
        ReadBinaryProto<ExactGeometryProto>(
            exact_geo_path, geo_storage_refs.exact_geometry_ref()));
    INTR_ASSIGN_OR_RETURN(ExactGeometry exact_geometry,
                          ToGeometry(exact_geometry_proto));
    std::shared_ptr<const Renderable> renderable_ptr = nullptr;
    if (geo_storage_refs.keep_renderable()) {
      const std::filesystem::path renderable_path =
          prefix_ / std::filesystem::path(geo_storage_refs.renderable_ref())
                        .relative_path();
      INTR_ASSIGN_OR_RETURN(
          std::string glb_string,
          ReadFileBytes(renderable_path, geo_storage_refs.renderable_ref()));
      renderable_ptr =
          std::make_shared<const Renderable>(std::move(glb_string));
    }

    return Geometry(std::move(exact_geometry), std::move(renderable_ptr),
                    geo_storage_refs.keep_renderable(), material_properties);
  }

 private:
  std::filesystem::path prefix_;
};

}  // namespace

std::unique_ptr<GeometryLibrary> GetFileSystemGeometryLibrary(
    absl::string_view prefix) {
  return std::make_unique<FileSystemGeometryLibrary>(prefix);
}

}  // namespace intrinsic::geo
