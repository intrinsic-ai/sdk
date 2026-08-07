// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/scene/util/scene_object_gzf.h"

#include <stdint.h>

#include <optional>
#include <string>
#include <utility>

#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "google/protobuf/io/coded_stream.h"
#include "google/protobuf/io/zero_copy_stream_impl_lite.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"
#include "intrinsic/world/gzfile/chunk_entry.h"
#include "intrinsic/world/gzfile/gzfile.h"

namespace intrinsic {
namespace scene_object {
namespace {

// The chunk id for the top-level SceneObject proto.
constexpr ChunkId::ValueType kSceneObjectProtoChunkId = 'SOBJ';

// Version 0 == The data is a SceneObject proto in binary format.
constexpr uint32_t kSceneObjectProtoChunkVersionNumber = 0;

}  // namespace

absl::Status AddSceneObjectToGzf(
    const intrinsic_proto::scene_object::v1::SceneObject& scene_object,
    GZFile& gzfile) {
  std::string pbtxt;
  {
    google::protobuf::io::StringOutputStream output_stream(&pbtxt);
    google::protobuf::io::CodedOutputStream coded_stream(&output_stream);
    // Enable deterministic serialization (see b/172338083).
    coded_stream.SetSerializationDeterministic(true);
    if (!scene_object.SerializeToCodedStream(&coded_stream)) {
      return absl::DataLossError("SceneObject cannot be serialized to chunk");
    }
  }
  ChunkEntry chunk(kSceneObjectProtoChunkVersionNumber, std::move(pbtxt));
  return gzfile.SetChunk(ChunkId(kSceneObjectProtoChunkId), chunk);
}

absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
GetSceneObjectFromGzf(const GZFile& gzfile) {
  std::optional<ChunkEntry> chunk =
      gzfile.GetChunk(ChunkId(kSceneObjectProtoChunkId));
  if (!chunk.has_value()) {
    return absl::NotFoundError("GZFile has no SceneObject chunk");
  }

  if (chunk->GetDataVersion() != kSceneObjectProtoChunkVersionNumber) {
    return absl::DataLossError(
        absl::StrCat("SceneObject chunk has incompatible data version: ",
                     chunk->GetDataVersion()));
  }

  intrinsic_proto::scene_object::v1::SceneObject scene_object;
  if (!scene_object.ParseFromString(chunk->GetUncompressedData())) {
    return absl::DataLossError("Could not parse the SceneObject proto");
  }

  return scene_object;
}

}  // namespace scene_object
}  // namespace intrinsic
