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

// Copyright 2026 Intrinsic Innovation LLC. All rights reserved.

#include <pybind11/eigen.h>
#include <pybind11/pybind11.h>
#include <pybind11/stl.h>

#include <optional>
#include <string>
#include <utility>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/eigenmath/python/pose3_caster.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/scene/processing/joint_editing.h"
#include "intrinsic/scene/processing/orient_scene_object.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"
#include "intrinsic/world/proto/kinematics_component.pb.h"
#include "pybind11_abseil/absl_casters.h"
#include "pybind11_abseil/status_casters.h"
#include "pybind11_protobuf/native_proto_caster.h"

namespace intrinsic::scene_object {
namespace {

using ::intrinsic::eigenmath::Vector3d;
using ::intrinsic_proto::world::KinematicsComponent;

PYBIND11_MODULE(scene_object_editing_py, m) {
  pybind11::google::ImportStatusModule();
  pybind11_protobuf::ImportNativeProtoCasters();

  m.doc() = "Python bindings for scene object editing.";

  pybind11::enum_<KinematicsComponent::MotionType>(m, "MotionType")
      .value("UNDEFINED", KinematicsComponent::MOTION_TYPE_UNDEFINED)
      .value("FIXED", KinematicsComponent::MOTION_TYPE_FIXED)
      .value("REVOLUTE", KinematicsComponent::MOTION_TYPE_REVOLUTE)
      .value("PRISMATIC", KinematicsComponent::MOTION_TYPE_PRISMATIC)
      .export_values();

  pybind11::class_<AddJointParams>(
      m, "AddJointParams",
      R"doc(Parameters for creating an articulated joint in a SceneObject.

Example:
  params = scene_object_editing_py.AddJointParams(
      parent_link="base_link",
      child_link="link_1",
      joint_name="j1",
      motion_type=scene_object_editing_py.MotionType.REVOLUTE,
      origin_root_frame=np.array([0.0, 0.0, 0.5]),
      axis_root_frame=np.array([0.0, 0.0, 1.0]),
      limits=(-1.57, 1.57),
  )

Attributes:
  parent_link: Name of the parent link entity.
  child_link: Name of the child link entity.
  joint_name: Unique name for the new joint entity.
  motion_type: Joint motion type (`MotionType.REVOLUTE`, `PRISMATIC`, or
    `FIXED`).
  origin_root_frame: 3D joint origin `(x, y, z)` in meters in the SceneObject's
    root entity frame.
  axis_root_frame: Non-zero 3D joint axis `(x, y, z)` in the SceneObject's root
    entity frame.
  limits: Optional `(lower, upper)` bounds in radians for revolute joints or
    meters for prismatic joints.
)doc")
      .def(pybind11::init<std::string, std::string, std::string,
                          KinematicsComponent::MotionType, Vector3d, Vector3d,
                          std::optional<std::pair<double, double>>>(),
           pybind11::arg("parent_link") = "", pybind11::arg("child_link") = "",
           pybind11::arg("joint_name") = "",
           pybind11::arg("motion_type") =
               KinematicsComponent::MOTION_TYPE_UNDEFINED,
           pybind11::arg("origin_root_frame") = Vector3d(0.0, 0.0, 0.0),
           pybind11::arg("axis_root_frame") = Vector3d(0.0, 0.0, 1.0),
           pybind11::arg("limits") = std::nullopt)
      .def_readwrite("parent_link", &AddJointParams::parent_link)
      .def_readwrite("child_link", &AddJointParams::child_link)
      .def_readwrite("joint_name", &AddJointParams::joint_name)
      .def_readwrite("motion_type", &AddJointParams::motion_type)
      .def_readwrite("origin_root_frame", &AddJointParams::origin_root_frame)
      .def_readwrite("axis_root_frame", &AddJointParams::axis_root_frame)
      .def_readwrite("limits", &AddJointParams::limits);

  pybind11::class_<UpdateJointParams>(
      m, "UpdateJointParams",
      R"doc(Parameters for partially updating an existing joint.

Example:
  params = scene_object_editing_py.UpdateJointParams(
      joint_name="j1",
      origin_root_frame=np.array([0.0, 0.0, 0.6]),
      limits=(-0.5, 0.75),
  )

Attributes:
  joint_name: Name of the existing joint entity to update.
  parent_link: Optional replacement parent link name (`None` preserves
    current).
  child_link: Optional replacement child link name (`None` preserves current).
  motion_type: Optional replacement `MotionType` (`None` preserves current).
  origin_root_frame: Optional 3D joint origin `(x, y, z)` in meters in the
    SceneObject's root entity frame (`None` preserves current).
  axis_root_frame: Optional non-zero 3D joint axis `(x, y, z)` in the
    SceneObject's root entity frame (`None` preserves current).
  limits: Optional `(lower, upper)` bounds (`None` preserves existing fixed
    limits on non-fixed joints).
)doc")
      .def(pybind11::init<std::string, std::optional<std::string>,
                          std::optional<std::string>,
                          std::optional<KinematicsComponent::MotionType>,
                          std::optional<Vector3d>, std::optional<Vector3d>,
                          std::optional<std::pair<double, double>>>(),
           pybind11::arg("joint_name") = "",
           pybind11::arg("parent_link") = std::nullopt,
           pybind11::arg("child_link") = std::nullopt,
           pybind11::arg("motion_type") = std::nullopt,
           pybind11::arg("origin_root_frame") = std::nullopt,
           pybind11::arg("axis_root_frame") = std::nullopt,
           pybind11::arg("limits") = std::nullopt)
      .def_readwrite("joint_name", &UpdateJointParams::joint_name)
      .def_readwrite("parent_link", &UpdateJointParams::parent_link)
      .def_readwrite("child_link", &UpdateJointParams::child_link)
      .def_readwrite("motion_type", &UpdateJointParams::motion_type)
      .def_readwrite("origin_root_frame", &UpdateJointParams::origin_root_frame)
      .def_readwrite("axis_root_frame", &UpdateJointParams::axis_root_frame)
      .def_readwrite("limits", &UpdateJointParams::limits);

  m.def("add_joint", &AddJoint, pybind11::arg("scene_object"),
        pybind11::arg("params"),
        R"doc(Adds a joint connecting parent_link and child_link.

Example:
  updated = scene_object_editing_py.add_joint(
      scene_object,
      scene_object_editing_py.AddJointParams(
          parent_link="base_link",
          child_link="link_1",
          joint_name="j1",
          motion_type=scene_object_editing_py.MotionType.REVOLUTE,
      ),
  )

Args:
  scene_object: Input `SceneObject` proto to copy and update.
  params: `AddJointParams` specifying link names, motion type, root-frame
    origin/axis, and optional limits.

Returns:
  Updated `SceneObject` proto containing the new joint entity.

Raises:
  status.StatusNotOk (RuntimeError): With `INVALID_ARGUMENT` if joint
    parameters, axis, limits, or entity types are invalid, or `NOT_FOUND` if
    `parent_link` or `child_link` does not exist.
)doc");
  m.def("update_joint", &UpdateJoint, pybind11::arg("scene_object"),
        pybind11::arg("params"),
        R"doc(Updates an existing joint while preserving None fields.

Example:
  updated = scene_object_editing_py.update_joint(
      scene_object,
      scene_object_editing_py.UpdateJointParams(
          joint_name="j1",
          limits=(-1.0, 1.0),
      ),
  )

Args:
  scene_object: Input `SceneObject` proto containing `params.joint_name`.
  params: `UpdateJointParams` specifying `joint_name` and optional overrides.

Returns:
  Updated `SceneObject` proto with the modified joint.

Raises:
  status.StatusNotOk (RuntimeError): With `INVALID_ARGUMENT` if `joint_name` is
    empty, names a non-joint entity, or specifies invalid kinematics, or
    `NOT_FOUND` if `joint_name` or target links do not exist.
)doc");
  m.def("delete_joint", &DeleteJoint, pybind11::arg("scene_object"),
        pybind11::arg("joint_name"),
        R"doc(Deletes a joint and reparents its children to its parent link.

Example:
  updated = scene_object_editing_py.delete_joint(scene_object, "j1")

Args:
  scene_object: Input `SceneObject` proto containing `joint_name`.
  joint_name: Name of the joint entity to delete.

Returns:
  Updated `SceneObject` proto with `joint_name` removed and children
  reparented.

Raises:
  status.StatusNotOk (RuntimeError): With `INVALID_ARGUMENT` if `joint_name` is
    empty or names a non-joint entity, or `NOT_FOUND` if `joint_name` does not
    exist.
)doc");

  pybind11::enum_<AxisDirection>(m, "AxisDirection",
                                 R"doc(A cardinal axis direction.

The direction is expressed in the SceneObject's root entity frame.

Example:
  up_axis = scene_object_editing_py.AxisDirection.POS_Y
)doc")
      .value("POS_X", AxisDirection::kPosX)
      .value("NEG_X", AxisDirection::kNegX)
      .value("POS_Y", AxisDirection::kPosY)
      .value("NEG_Y", AxisDirection::kNegY)
      .value("POS_Z", AxisDirection::kPosZ)
      .value("NEG_Z", AxisDirection::kNegZ)
      .export_values();

  pybind11::class_<OrientAndRecenterResult>(
      m, "OrientAndRecenterResult",
      R"doc(Result of orient_and_recenter_scene_object.

Example:
  result = scene_object_editing_py.orient_and_recenter_scene_object(
      scene_object,
      up_axis=scene_object_editing_py.AxisDirection.POS_Y,
      origin_offset=np.array([0.0, 0.0, 0.0]),
  )
  oriented_scene_object = result.scene_object

Attributes:
  scene_object: The oriented and recentered SceneObject proto.
  applied_transform: The new_root_t_old_root data_types.Pose3 applied to all
    children of the root entity and to the root entity's geometry.
)doc")
      .def_readonly("scene_object", &OrientAndRecenterResult::scene_object)
      .def_readonly("applied_transform",
                    &OrientAndRecenterResult::applied_transform);

  m.def("orient_and_recenter_scene_object", &OrientAndRecenterSceneObject,
        pybind11::arg("scene_object"), pybind11::arg("up_axis"),
        pybind11::arg("origin_offset"),
        R"doc(Orients a SceneObject to +Z-up, then recenters it.

The two steps run in this order:
  1. Rotation first: rotates scene_object about the root origin so that
     up_axis maps to +Z (p_rotated = R * p).
  2. Recentering second, in the post-rotation +Z-up frame: shifts the rotated
     root frame by -origin_offset (p_new = (R * p) - origin_offset). The
     post-rotation point origin_offset moves to the new root origin (0, 0, 0).

For example, a +Y-up SceneObject (AxisDirection.POS_Y) has a point at
(0, 2, 0), and origin_offset = (0, 0, 0.5):
  (1) Rotating +Y to +Z (+90 deg about +X) maps (0, 2, 0) to (0, 0, 2) in the
      +Z-up frame.
  (2) Translating by -origin_offset = (0, 0, -0.5) maps (0, 0, 2) to
      (0, 0, 1.5), and places (0, 0, 0.5) at the new root origin (0, 0, 0).

Example:
  result = scene_object_editing_py.orient_and_recenter_scene_object(
      scene_object,
      up_axis=scene_object_editing_py.AxisDirection.POS_Y,
      origin_offset=np.array([0.0, 0.0, 0.5]),
  )

Args:
  scene_object: The SceneObject proto to transform.
  up_axis: The axis of scene_object that points up. It is mapped to +Z.
  origin_offset: The point, in meters in the post-rotation +Z-up frame, that
    becomes the new root origin.

Returns:
  An OrientAndRecenterResult with the transformed SceneObject and the applied
  data_types.Pose3 transform.

Raises:
  RuntimeError: INVALID_ARGUMENT if origin_offset has a non-finite
    component. If the applied transform is not the identity, also the error of
    the underlying SceneObject transform (for example, INVALID_ARGUMENT if
    scene_object has no unique root entity).
)doc");
}

}  // namespace
}  // namespace intrinsic::scene_object
