// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/internal/legacy/mesh/io/ai_scene_to_mesh.h"

#include <cmath>
#include <utility>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "assimp/mesh.h"
#include "assimp/scene.h"
#include "assimp/vector3.h"
#include "intrinsic/geometry/internal/legacy/mesh/concatenate_meshes.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"

namespace intrinsic {
namespace geometry_legacy {
absl::StatusOr<Mesh> AiSceneToMesh(const aiScene& scene) {
  if (scene.mRootNode == nullptr) {
    return absl::InvalidArgumentError("Input scene has no root node!");
  }

  if (!scene.mRootNode->mTransformation.IsIdentity()) {
    return absl::InvalidArgumentError(
        "We do not support transformations on the scene. b/180961517");
  }

  std::vector<Mesh> meshes;
  for (int mesh_index = 0; mesh_index < scene.mNumMeshes; mesh_index++) {
    aiMesh* ai_scene_mesh = scene.mMeshes[mesh_index];

    Mesh::VertexCollection vertices;
    vertices.reserve(ai_scene_mesh->mNumVertices);
    // Add new rows to the vertex matrix and fill them in.
    for (int vdx = 0; vdx < ai_scene_mesh->mNumVertices; vdx++) {
      const aiVector3D& p = ai_scene_mesh->mVertices[vdx];
      if (!(isfinite(p.x) && isfinite(p.y) && isfinite(p.z))) {
        return absl::InvalidArgumentError("Mesh contains invalid vertices.");
      }
      vertices.emplace_back(p.x, p.y, p.z);
    }

    Mesh::FaceCollection faces;
    faces.reserve(ai_scene_mesh->mNumFaces);
    // Add new rows to the face matrix and fill them in.
    for (int fdx = 0; fdx < ai_scene_mesh->mNumFaces; fdx++) {
      const aiFace& face = ai_scene_mesh->mFaces[fdx];
      if (face.mNumIndices != 3) {
        return absl::InternalError("Unsupported number of vertices in a face");
      }

      faces.emplace_back(face.mIndices[0], face.mIndices[1], face.mIndices[2]);
    }
    meshes.emplace_back(std::move(vertices), std::move(faces));
  }

  return ConcatenateMeshes(meshes);
}

}  // namespace geometry_legacy
}  // namespace intrinsic
