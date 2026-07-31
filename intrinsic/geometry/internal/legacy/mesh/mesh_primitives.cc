// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/internal/legacy/mesh/mesh_primitives.h"

#include <algorithm>
#include <cmath>
#include <cstddef>
#include <cstdint>
#include <map>
#include <utility>
#include <vector>

#include "absl/log/check.h"
#include "absl/numeric/bits.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"

namespace intrinsic {
namespace geometry_legacy {

namespace {

void CreateSphericalVertex(int i, int j, int N, double radius, int nb_cap_rings,
                           size_t index, Mesh::VertexCollection* vertices) {
  double theta = i * 2 * M_PI / N;
  double phi = -M_PI / 2 + M_PI * j / (2 * nb_cap_rings);
  (*vertices)[index] = eigenmath::Vector3d(
      radius * std::cos(phi) * std::cos(theta),
      radius * std::cos(phi) * std::sin(theta), radius * std::sin(phi));
}

Mesh CreateBaseOctahedron() {
  constexpr int n = 0;
  constexpr int s = 1;
  constexpr int a = 2;
  constexpr int b = 3;
  constexpr int c = 4;
  constexpr int d = 5;

  Mesh::VertexCollection vertices(6);
  vertices[n] = eigenmath::Vector3d(0.0, 0.0, 1.0);
  vertices[s] = eigenmath::Vector3d(0.0, 0.0, -1.0);
  vertices[a] = eigenmath::Vector3d(0.0, 1.0, 0.0);
  vertices[b] = eigenmath::Vector3d(1.0, 0.0, 0.0);
  vertices[c] = eigenmath::Vector3d(0.0, -1.0, 0.0);
  vertices[d] = eigenmath::Vector3d(-1.0, 0.0, 0.0);

  Mesh::FaceCollection faces = {{n, b, a}, {n, c, b}, {n, d, c}, {n, a, d},
                                {s, a, b}, {s, b, c}, {s, c, d}, {s, d, a}};
  return Mesh(std::move(vertices), std::move(faces));
}

Mesh SubdivideLinearOneIteration(const Mesh& mesh) {
  const Mesh::FaceCollection& old_faces = mesh.faces();
  Mesh::VertexCollection vertices = mesh.vertices();
  vertices.reserve(vertices.size() + old_faces.size() * 2);
  Mesh::FaceCollection new_faces;
  new_faces.reserve(old_faces.size() * 4);

  // Cache for midpoints: edge (v0, v1) -> midpoint vertex index
  std::map<std::pair<int, int>, int> midpoint_cache;

  auto get_midpoint = [&](int v0, int v1) {
    if (v0 > v1) std::swap(v0, v1);
    auto edge = std::make_pair(v0, v1);
    auto it = midpoint_cache.find(edge);
    if (it != midpoint_cache.end()) {
      return it->second;
    }
    // Compute midpoint
    eigenmath::Vector3d midpoint = (vertices[v0] + vertices[v1]) * 0.5;
    int new_idx = vertices.size();
    vertices.push_back(midpoint);
    midpoint_cache[edge] = new_idx;
    return new_idx;
  };

  for (const auto& face : old_faces) {
    int v0 = face[0];
    int v1 = face[1];
    int v2 = face[2];

    int m01 = get_midpoint(v0, v1);
    int m12 = get_midpoint(v1, v2);
    int m20 = get_midpoint(v2, v0);

    new_faces.emplace_back(v0, m01, m20);
    new_faces.emplace_back(v1, m12, m01);
    new_faces.emplace_back(v2, m20, m12);
    new_faces.emplace_back(m01, m12, m20);
  }
  return Mesh(std::move(vertices), std::move(new_faces));
}

Mesh ProjectToUnitSphere(Mesh&& mesh) {
  Mesh::VertexCollection vertices = std::move(mesh).vertices();
  for (auto& v : vertices) {
    v.normalize();
  }
  return Mesh(std::move(vertices), std::move(mesh).faces());
}

Mesh CreateUnitSphereOctahedron(int recursion_level) {
  Mesh octahedron = CreateBaseOctahedron();
  for (int level = 0; level < recursion_level; ++level) {
    octahedron = SubdivideLinearOneIteration(octahedron);
    octahedron = ProjectToUnitSphere(std::move(octahedron));
  }
  return octahedron;
}

}  // namespace

// Interface implementation
Mesh CreateCylinder(const int N, const double radius, const double length) {
  CHECK(N >= 3) << "N should be at least 3 to approximate a cylinder.";

  Mesh::VertexCollection vertices(2 * N);
  for (int i = 0; i < N; ++i) {
    double theta = i * 2 * M_PI / N;
    vertices[i] = eigenmath::Vector3d(radius * std::cos(theta),
                                      radius * std::sin(theta), length / 2);
    vertices[N + i] = eigenmath::Vector3d(
        radius * std::cos(theta), radius * std::sin(theta), -length / 2);
  }

  Mesh::FaceCollection faces;
  faces.reserve(4 * N - 4);

  // Caps
  for (int i = 1; i <= N - 2; ++i) {
    faces.emplace_back(0, i, i + 1);
    faces.emplace_back(N, N + i + 1, N + i);
  }
  // Sides
  for (int i = 0; i < N; ++i) {
    int next_i = (i + 1) % N;
    faces.emplace_back(i, N + i, next_i);
    faces.emplace_back(next_i, N + i, N + next_i);
  }

  return Mesh(std::move(vertices), std::move(faces));
}

Mesh CreateCapsule(int resolution, double radius, double height) {
  CHECK_GE(resolution, 0) << "Resolution must be non-negative.";
  CHECK_LE(resolution, 28) << "Resolution must be less than or equal to 28 in "
                              "the current implementation.";
  Mesh::VertexCollection vertices;
  Mesh::FaceCollection faces;

  const int N = 4 << resolution;
  const int nb_cap_rings = N / 4;

  const size_t nb_vertex = (2 * nb_cap_rings) * N + 2;

  vertices.resize(nb_vertex);

  // double theta, phi;
  size_t index = 0;

  // Bottom hemisphere tip vertex
  vertices[index] = eigenmath::Vector3d(0, 0, -radius - height / 2);
  index++;

  // The rest of the bottom hemisphere
  for (int j = 1; j <= nb_cap_rings; ++j) {
    for (int i = 0; i < N; ++i) {
      CreateSphericalVertex(i, j, N, radius, nb_cap_rings, index, &vertices);
      vertices[index][2] -= height / 2;
      index++;
    }
  }

  // Top hemisphere
  for (int j = nb_cap_rings; j <= 2 * nb_cap_rings - 1; ++j) {
    for (int i = 0; i < N; ++i) {
      CreateSphericalVertex(i, j, N, radius, nb_cap_rings, index, &vertices);
      vertices[index][2] += height / 2;
      index++;
    }
  }

  // Finish with the tip of the top hemisphere
  vertices[index] = eigenmath::Vector3d(0, 0, radius + height / 2);
  index++;

  CHECK_EQ(index, nb_vertex);

  // Create faces
  const size_t nb_faces =
      (2 * nb_cap_rings - 1) * (2 * (N - 1) + 2) + 2 * (N - 1) + 2;
  faces.resize(nb_faces, Mesh::Face::Zero());
  index = 0;
  int a;
  int b;
  int c;
  int d;
  for (int j = 1; j < 2 * nb_cap_rings; ++j) {
    for (int i = 0; i < N - 1; ++i) {
      a = j * N + i;
      b = j * N + (i + 1);
      c = (j + 1) * N + (i + 1);
      d = (j + 1) * N + i;

      faces[index] << a, b, c;
      faces[index + 1] << a, c, d;

      index += 2;
    }

    // Close the surface along the circumference
    int i = N - 1;
    a = j * N + i;
    b = j * N;
    c = j * N + (i + 1);
    d = (j + 1) * N + i;
    faces[index] << a, b, c;
    faces[index + 1] << a, c, d;
    index += 2;
  }

  // Subtracts face indexes to take into account that the bottom tip vertex
  // is not duplicated;
  for (auto& f : faces) {
    f[0] -= (N - 1);
    f[1] -= (N - 1);
    f[2] -= (N - 1);
  }

  // Finish closing the sphere on the single tip vertex
  for (int i = 0; i < N - 1; ++i) {
    a = 0;
    c = N + (i + 1) - (N - 1);
    d = N + i - (N - 1);
    faces[index] << a, c, d;

    int j = 2 * nb_cap_rings - 1;
    a = nb_vertex - 1;
    c = (j + 1) * N + i + 1 - (N - 1);
    d = (j + 1) * N + i - (N - 1);
    faces[index + 1] << a, d, c;

    index += 2;
  }

  // Lastly, finish with the last two faces that close along the circumference
  // and the tips.
  a = 0;
  c = 1;
  d = N;
  faces[index] << a, c, d;

  a = nb_vertex - 1;
  c = (2 * nb_cap_rings * N) - (N - 1);
  d = 2 * nb_cap_rings * N;
  faces[index + 1] << a, d, c;

  return Mesh{vertices, faces};
}

Mesh CreateSphere(int recursion_level, double radius) {
  return CreateEllipsoid(recursion_level, radius, radius, radius);
}

Mesh CreateEllipsoid(int recursion_level, double radius_x, double radius_y,
                     double radius_z) {
  Mesh ellipsoid = CreateUnitSphereOctahedron(recursion_level);

  eigenmath::Vector3d radius_scaling;
  radius_scaling << radius_x, radius_y, radius_z;
  ellipsoid.Scale(radius_scaling);
  return ellipsoid;
}

Mesh CreateCuboid(const Eigen::Vector3d& half_edge) {
  Mesh::VertexCollection vertices;

  const double& x = half_edge(0);
  const double& y = half_edge(1);
  const double& z = half_edge(2);

  vertices.reserve(8);
  vertices.emplace_back(-x, -y, -z);
  vertices.emplace_back(-x, -y, z);
  vertices.emplace_back(-x, y, z);
  vertices.emplace_back(-x, y, -z);
  vertices.emplace_back(x, -y, -z);
  vertices.emplace_back(x, -y, z);
  vertices.emplace_back(x, y, z);
  vertices.emplace_back(x, y, -z);

  static Mesh::FaceCollection* faces = []() {
    auto* faces = new Mesh::FaceCollection();
    faces->reserve(12);
    faces->emplace_back(0, 1, 2);
    faces->emplace_back(2, 3, 0);
    faces->emplace_back(4, 6, 5);
    faces->emplace_back(4, 7, 6);
    faces->emplace_back(0, 5, 1);
    faces->emplace_back(0, 4, 5);
    faces->emplace_back(2, 7, 3);
    faces->emplace_back(2, 6, 7);
    faces->emplace_back(3, 4, 0);
    faces->emplace_back(3, 7, 4);
    faces->emplace_back(1, 6, 2);
    faces->emplace_back(1, 5, 6);
    return faces;
  }();

  return Mesh(std::move(vertices), Mesh::FaceCollection(*faces));
}

Mesh CreatePartialCuboid(const Eigen::Vector3d& half_edge, uint8_t sides) {
  Mesh::VertexCollection vertices;
  Mesh::FaceCollection faces;

  const double& x = half_edge(0);
  const double& y = half_edge(1);
  const double& z = half_edge(2);

  vertices.reserve(8);

  vertices.emplace_back(-x, -y, -z);
  vertices.emplace_back(-x, -y, z);
  vertices.emplace_back(-x, y, z);
  vertices.emplace_back(-x, y, -z);
  vertices.emplace_back(x, -y, -z);
  vertices.emplace_back(x, -y, z);
  vertices.emplace_back(x, y, z);
  vertices.emplace_back(x, y, -z);

  faces.reserve(absl::popcount(sides) * 2);

  // -X
  if (sides & 1) {
    faces.emplace_back(0, 1, 2);
    faces.emplace_back(2, 3, 0);
  }

  // +X
  if (sides & (1 << 1)) {
    faces.emplace_back(4, 6, 5);
    faces.emplace_back(4, 7, 6);
  }

  // -Y
  if (sides & (1 << 2)) {
    faces.emplace_back(0, 5, 1);
    faces.emplace_back(0, 4, 5);
  }

  // +Y
  if (sides & (1 << 3)) {
    faces.emplace_back(2, 7, 3);
    faces.emplace_back(2, 6, 7);
  }

  // -Z
  if (sides & (1 << 4)) {
    faces.emplace_back(3, 4, 0);
    faces.emplace_back(3, 7, 4);
  }

  // +Z
  if (sides & (1 << 5)) {
    faces.emplace_back(1, 6, 2);
    faces.emplace_back(1, 5, 6);
  }

  CHECK(((sides & (1 << 6)) == 0) && ((sides & (1 << 7))) == 0)
      << "Illegal bits set on sides mask. Only bits 0-5 are legal.";

  return Mesh(std::move(vertices), std::move(faces));
}

absl::StatusOr<Mesh> CreateFrustumMesh(double x_angle, double y_angle,
                                       double min_z_distance,
                                       double max_z_distance) {
  // Validate frustum parameters
  if (x_angle < 0.0) {
    return absl::InvalidArgumentError(
        "Frustum x angle must be greater than or equal to 0");
  }
  if (x_angle >= M_PI * 0.5) {
    return absl::InvalidArgumentError(
        "Frustum x angle must be less than 90 degrees or PI/2 radians");
  }

  if (y_angle < 0.0) {
    return absl::InvalidArgumentError(
        "Frustum y angle must be greater than or equal to 0");
  }
  if (y_angle >= M_PI * 0.5) {
    return absl::InvalidArgumentError(
        "Frustum y angle must be less than 90 degrees or PI/2 radians");
  }
  if (min_z_distance < 0.0) {
    return absl::InvalidArgumentError(
        "Frustum min z distance must be greater than or equal to 0");
  }
  if (max_z_distance < min_z_distance) {
    return absl::InvalidArgumentError(
        "Frustum max z distance must be greater than or equal min z distance");
  }
  const double tan_x = tan(x_angle);
  const double tan_y = tan(y_angle);

  const double& near_z = min_z_distance;
  const double near_y = min_z_distance * tan_x;
  const double near_x = min_z_distance * tan_y;
  const double& far_z = max_z_distance;
  const double far_y = max_z_distance * tan_x;
  const double far_x = max_z_distance * tan_y;

  Mesh::VertexCollection vertices;
  vertices.reserve(8);
  // Near plane points
  vertices.emplace_back(-near_x, -near_y, near_z);
  vertices.emplace_back(-near_x, near_y, near_z);
  vertices.emplace_back(near_x, near_y, near_z);
  vertices.emplace_back(near_x, -near_y, near_z);
  // Far plane points
  vertices.emplace_back(-far_x, -far_y, far_z);
  vertices.emplace_back(-far_x, far_y, far_z);
  vertices.emplace_back(far_x, far_y, far_z);
  vertices.emplace_back(far_x, -far_y, far_z);

  Mesh::FaceCollection faces;
  faces.reserve(12);
  // Near z face
  faces.emplace_back(0, 1, 2);
  faces.emplace_back(2, 3, 0);

  // Far z face
  faces.emplace_back(4, 6, 5);
  faces.emplace_back(4, 7, 6);

  // -x face
  faces.emplace_back(0, 5, 1);
  faces.emplace_back(0, 4, 5);

  // +x face
  faces.emplace_back(2, 7, 3);
  faces.emplace_back(2, 6, 7);

  // -y face
  faces.emplace_back(3, 4, 0);
  faces.emplace_back(3, 7, 4);

  // +y face
  faces.emplace_back(1, 6, 2);
  faces.emplace_back(1, 5, 6);

  return Mesh(std::move(vertices), std::move(faces));
}
}  // namespace geometry_legacy
}  // namespace intrinsic
