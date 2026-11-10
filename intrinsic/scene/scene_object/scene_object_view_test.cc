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

#include "intrinsic/scene/scene_object/scene_object_view.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <algorithm>
#include <iterator>
#include <random>
#include <string>
#include <type_traits>
#include <utility>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/strings/substitute.h"
#include "absl/types/span.h"
#include "internal/testing.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/scene/proto/v1/entity.pb.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"
#include "intrinsic/util/proto/parse_text_proto.h"
#include "intrinsic/world/proto/kinematics_component.pb.h"
#include "protobuf-matchers/protocol-buffer-matchers.h"

namespace intrinsic {
namespace scene_object {
namespace {

using ::absl_testing::IsOkAndHolds;
using ::absl_testing::StatusIs;
using ::intrinsic::ParseTextProtoOrDie;
using ::intrinsic_proto::scene_object::v1::Entity;
using ::intrinsic_proto::scene_object::v1::SceneObject;
using ::intrinsic_proto::world::KinematicsComponent;
using ::protobuf_matchers::EqualsProto;
using ::testing::AllOf;
using ::testing::ElementsAre;
using ::testing::ElementsAreArray;
using ::testing::HasSubstr;
using ::testing::IsEmpty;
using ::testing::IsNull;
using ::testing::Pair;
using ::testing::Pointee;
using ::testing::Property;
using ::testing::SizeIs;
using ::testing::UnorderedElementsAre;

// Matches a `const Entity*` that points to an entity named `name`.
auto EntityNamed(const std::string& name) {
  return Pointee(Property(&Entity::name, name));
}

// Matches a `Pose3d` that is approximately equal to `expected`.
MATCHER_P(IsApprox, expected, "") { return arg.isApprox(expected); }

// Returns the names of `entities`.
std::vector<std::string> Names(absl::Span<const Entity* const> entities) {
  std::vector<std::string> names;
  names.reserve(entities.size());
  std::ranges::transform(entities, std::back_inserter(names),
                         [](const Entity* entity) { return entity->name(); });
  return names;
}

// Appends an entity named `name` with parent `parent_name` to `scene_object`
// and returns it. The caller sets the entity type.
Entity& AddEntity(SceneObject& scene_object, absl::string_view name,
                  absl::string_view parent_name) {
  Entity& entity = *scene_object.add_entities();
  entity.set_name(name);
  entity.set_parent_name(parent_name);
  return entity;
}

// Appends a joint entity with the given `motion_type` to `scene_object`.
void AddJoint(SceneObject& scene_object, absl::string_view name,
              absl::string_view parent_name,
              KinematicsComponent::MotionType motion_type) {
  AddEntity(scene_object, name, parent_name)
      .mutable_joint()
      ->mutable_kinematics_component()
      ->set_motion_type(motion_type);
}

// Kinematic tree used by most tests, with entities intentionally not listed
// in tree order:
//
// base
// |-- joint_b -> link_b
// |               `-- tool_frame
// |-- joint_a -> link_a
// |               `-- sensor
// `-- base_frame
SceneObject MakeTreeSceneObject() {
  return ParseTextProtoOrDie(R"pb(
    name: "tree"
    entities {
      name: "link_b"
      parent_name: "joint_b"
      link {}
    }
    entities {
      name: "tool_frame"
      parent_name: "link_b"
      frame {}
    }
    entities {
      name: "joint_b"
      parent_name: "base"
      joint { kinematics_component { motion_type: MOTION_TYPE_FIXED } }
    }
    entities {
      name: "base"
      link {}
    }
    entities {
      name: "base_frame"
      parent_name: "base"
      frame {}
    }
    entities {
      name: "joint_a"
      parent_name: "base"
      joint { kinematics_component { motion_type: MOTION_TYPE_FIXED } }
    }
    entities {
      name: "sensor"
      parent_name: "link_a"
      sensor {}
    }
    entities {
      name: "link_a"
      parent_name: "joint_a"
      link {}
    }
  )pb");
}

TEST(SceneObjectViewTest, FailsForInvalidSceneObject) {
  const SceneObject no_root = ParseTextProtoOrDie(R"pb(
    name: "no_root"
    entities {
      name: "link"
      parent_name: "missing"
      link {}
    }
  )pb");
  EXPECT_THAT(SceneObjectView::Create(no_root),
              StatusIs(absl::StatusCode::kInvalidArgument));

  const SceneObject multiple_roots = ParseTextProtoOrDie(R"pb(
    name: "multiple_roots"
    entities {
      name: "link1"
      link {}
    }
    entities {
      name: "link2"
      link {}
    }
  )pb");
  EXPECT_THAT(SceneObjectView::Create(multiple_roots),
              StatusIs(absl::StatusCode::kInvalidArgument));

  const SceneObject duplicate_names = ParseTextProtoOrDie(R"pb(
    name: "duplicate_names"
    entities {
      name: "link"
      link {}
    }
    entities {
      name: "link"
      parent_name: "link"
      link {}
    }
  )pb");
  EXPECT_THAT(SceneObjectView::Create(duplicate_names),
              StatusIs(absl::StatusCode::kInvalidArgument));

  const SceneObject empty = ParseTextProtoOrDie(R"pb(name: "empty")pb");
  EXPECT_THAT(SceneObjectView::Create(empty),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(SceneObjectViewTest, SingleEntity) {
  const SceneObject scene_object = ParseTextProtoOrDie(R"pb(
    name: "single"
    entities {
      name: "link"
      link {}
    }
  )pb");
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(scene_object));

  EXPECT_THAT(view.scene_object(), EqualsProto(scene_object));
  EXPECT_EQ(&view.GetRootEntity(), &view.scene_object().entities(0));
  EXPECT_EQ(view.GetNumEntities(), 1);
  EXPECT_THAT(view.GetParent("link"), IsNull());
  EXPECT_THAT(view.GetChildren("link"), IsEmpty());
  EXPECT_THAT(view.GetEntitiesInTreeOrder(), ElementsAre(EntityNamed("link")));
  EXPECT_THAT(view.GetSubtreeInTreeOrder("link"),
              ElementsAre(EntityNamed("link")));
  EXPECT_TRUE(view.IsInSubtree("link", "link"));
  EXPECT_THAT(view.GetPathFromRoot("link"), ElementsAre(EntityNamed("link")));
}

TEST(SceneObjectViewTest, ReturnsRootEntity) {
  const SceneObject scene_object = MakeTreeSceneObject();
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(scene_object));

  EXPECT_EQ(view.GetRootEntity().name(), "base");
  EXPECT_EQ(view.GetNumEntities(), 8);
}

TEST(SceneObjectViewTest, FindsEntities) {
  const SceneObject scene_object = MakeTreeSceneObject();
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(scene_object));

  for (const Entity& entity : view.scene_object().entities()) {
    EXPECT_TRUE(view.HasEntity(entity.name()));
    EXPECT_EQ(view.FindEntity(entity.name()), &entity);
  }
  EXPECT_FALSE(view.HasEntity("unknown"));
  EXPECT_THAT(view.FindEntity("unknown"), IsNull());
}

TEST(SceneObjectViewTest, ReturnsParents) {
  const SceneObject scene_object = MakeTreeSceneObject();
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(scene_object));

  EXPECT_THAT(view.GetParent("base"), IsNull());
  EXPECT_THAT(view.GetParent("joint_a"), EntityNamed("base"));
  EXPECT_THAT(view.GetParent("link_a"), EntityNamed("joint_a"));
  EXPECT_THAT(view.GetParent("sensor"), EntityNamed("link_a"));
  EXPECT_THAT(view.GetParent("tool_frame"), EntityNamed("link_b"));
  EXPECT_THAT(view.GetParent("unknown"), IsNull());
}

TEST(SceneObjectViewTest, ReturnsChildrenOrderedByName) {
  const SceneObject scene_object = MakeTreeSceneObject();
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(scene_object));

  EXPECT_THAT(view.GetChildren("base"),
              ElementsAre(EntityNamed("base_frame"), EntityNamed("joint_a"),
                          EntityNamed("joint_b")));
  EXPECT_THAT(view.GetChildren("joint_b"), ElementsAre(EntityNamed("link_b")));
  EXPECT_THAT(view.GetChildren("base_frame"), IsEmpty());
  EXPECT_THAT(view.GetChildren("unknown"), IsEmpty());
}

TEST(SceneObjectViewTest, ReturnsEntitiesInTreeOrder) {
  const SceneObject scene_object = MakeTreeSceneObject();
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(scene_object));

  EXPECT_THAT(view.GetEntitiesInTreeOrder(),
              ElementsAre(EntityNamed("base"), EntityNamed("base_frame"),
                          EntityNamed("joint_a"), EntityNamed("link_a"),
                          EntityNamed("sensor"), EntityNamed("joint_b"),
                          EntityNamed("link_b"), EntityNamed("tool_frame")));
}

TEST(SceneObjectViewTest, TreeOrderIsIndependentOfEntityOrder) {
  const SceneObject scene_object = MakeTreeSceneObject();
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(scene_object));
  std::vector<std::string> expected_names;
  for (const Entity* entity : view.GetEntitiesInTreeOrder()) {
    expected_names.push_back(entity->name());
  }

  std::mt19937 rng(/*seed=*/42);
  for (int i = 0; i < 10; ++i) {
    SceneObject shuffled = scene_object;
    std::ranges::shuffle(*shuffled.mutable_entities(), rng);
    ASSERT_OK_AND_ASSIGN(const SceneObjectView shuffled_view,
                         SceneObjectView::Create(shuffled));

    std::vector<std::string> names;
    for (const Entity* entity : shuffled_view.GetEntitiesInTreeOrder()) {
      names.push_back(entity->name());
    }
    EXPECT_EQ(names, expected_names);
  }
}

TEST(SceneObjectViewTest, ReturnsSubtreeInTreeOrder) {
  const SceneObject scene_object = MakeTreeSceneObject();
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(scene_object));

  // The subtree of the root entity is the whole tree.
  EXPECT_THAT(view.GetSubtreeInTreeOrder("base"),
              ElementsAreArray(view.GetEntitiesInTreeOrder()));
  EXPECT_THAT(view.GetSubtreeInTreeOrder("joint_a"),
              ElementsAre(EntityNamed("joint_a"), EntityNamed("link_a"),
                          EntityNamed("sensor")));
  EXPECT_THAT(view.GetSubtreeInTreeOrder("joint_b"),
              ElementsAre(EntityNamed("joint_b"), EntityNamed("link_b"),
                          EntityNamed("tool_frame")));
  EXPECT_THAT(view.GetSubtreeInTreeOrder("tool_frame"),
              ElementsAre(EntityNamed("tool_frame")));
  EXPECT_THAT(view.GetSubtreeInTreeOrder("unknown"), IsEmpty());
}

TEST(SceneObjectViewTest, ChecksSubtreeMembership) {
  const SceneObject scene_object = MakeTreeSceneObject();
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(scene_object));

  for (const Entity& entity : view.scene_object().entities()) {
    EXPECT_TRUE(view.IsInSubtree(entity.name(), "base"));
    EXPECT_TRUE(view.IsInSubtree(entity.name(), entity.name()));
  }
  EXPECT_TRUE(view.IsInSubtree("link_a", "joint_a"));
  EXPECT_TRUE(view.IsInSubtree("sensor", "joint_a"));
  EXPECT_FALSE(view.IsInSubtree("base", "joint_a"));
  EXPECT_FALSE(view.IsInSubtree("base_frame", "joint_a"));
  EXPECT_FALSE(view.IsInSubtree("tool_frame", "joint_a"));
  EXPECT_FALSE(view.IsInSubtree("link_b", "tool_frame"));
  EXPECT_FALSE(view.IsInSubtree("unknown", "base"));
  EXPECT_FALSE(view.IsInSubtree("base", "unknown"));
  EXPECT_FALSE(view.IsInSubtree("unknown", "unknown"));
}

TEST(SceneObjectViewTest, ReturnsPathFromRoot) {
  const SceneObject scene_object = MakeTreeSceneObject();
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(scene_object));

  EXPECT_THAT(view.GetPathFromRoot("base"), ElementsAre(EntityNamed("base")));
  EXPECT_THAT(view.GetPathFromRoot("tool_frame"),
              ElementsAre(EntityNamed("base"), EntityNamed("joint_b"),
                          EntityNamed("link_b"), EntityNamed("tool_frame")));
  EXPECT_THAT(view.GetPathFromRoot("unknown"), IsEmpty());
}

TEST(SceneObjectViewTest, ReturnsRootPoses) {
  // joint_a is rotated 90 degrees about z, so link_a's offset along x becomes
  // an offset along y in the root frame. base_frame has no parent_t_this.
  const SceneObject scene_object = ParseTextProtoOrDie(R"pb(
    name: "poses"
    entities {
      name: "link_a"
      parent_name: "joint_a"
      parent_t_this {
        position { x: 1 }
        orientation { w: 1 }
      }
      link {}
    }
    entities {
      name: "base"
      link {}
    }
    entities {
      name: "joint_a"
      parent_name: "base"
      parent_t_this {
        position { x: 1 }
        orientation { z: 0.7071067811865476 w: 0.7071067811865476 }
      }
      joint { kinematics_component { motion_type: MOTION_TYPE_FIXED } }
    }
    entities {
      name: "base_frame"
      parent_name: "base"
      frame {}
    }
  )pb");
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(scene_object));

  ASSERT_OK_AND_ASSIGN(const auto root_poses, view.GetRootPoses());
  ASSERT_EQ(root_poses.size(), 4);
  const eigenmath::Quaterniond rotate_z_90(0.7071067811865476, 0, 0,
                                           0.7071067811865476);
  EXPECT_TRUE(root_poses.at("base").isApprox(Pose3d::Identity()));
  EXPECT_TRUE(root_poses.at("base_frame").isApprox(Pose3d::Identity()));
  EXPECT_TRUE(root_poses.at("joint_a").isApprox(
      Pose3d(rotate_z_90, eigenmath::Vector3d(1, 0, 0))));
  EXPECT_TRUE(root_poses.at("link_a").isApprox(
      Pose3d(rotate_z_90, eigenmath::Vector3d(1, 1, 0))));
}

// A revolute or prismatic joint that has `raw_value` set and no
// `parent_t_this` is posed by parent_t_inboard * motion(raw_value) *
// outboard_t_child. The child link has no `parent_t_this`, so it has the same
// root pose as the joint.
struct JointPoseCase {
  std::string name;
  std::string kinematics_component;
  Pose3d expected_root_t_joint;
};

class SceneObjectViewJointPoseTest
    : public ::testing::TestWithParam<JointPoseCase> {};

TEST_P(SceneObjectViewJointPoseTest,
       ReturnsRootPosesFromRawValueWhenParentTThisIsUnset) {
  const JointPoseCase& test_case = GetParam();
  const SceneObject scene_object = ParseTextProtoOrDie(absl::Substitute(
      R"pb(
        name: "single_joint"
        entities {
          name: "base"
          link {}
        }
        entities {
          name: "joint"
          parent_name: "base"
          joint { kinematics_component { $0 } }
        }
        entities {
          name: "link"
          parent_name: "joint"
          link {}
        }
      )pb",
      test_case.kinematics_component));
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(scene_object));

  EXPECT_THAT(view.GetRootPoses(),
              IsOkAndHolds(UnorderedElementsAre(
                  Pair("base", IsApprox(Pose3d::Identity())),
                  Pair("joint", IsApprox(test_case.expected_root_t_joint)),
                  Pair("link", IsApprox(test_case.expected_root_t_joint)))));
}

INSTANTIATE_TEST_SUITE_P(
    MovableJoints, SceneObjectViewJointPoseTest,
    ::testing::Values(
        // Rotating outboard_t_child by 90 degrees about z moves its offset
        // along x to an offset along y.
        JointPoseCase{"RevoluteJoint", R"pb(
                        motion_type: MOTION_TYPE_REVOLUTE
                        axis { z: 1 }
                        raw_value: 1.5707963267948966
                        parent_t_inboard {
                          position { x: 1 }
                          orientation { w: 1 }
                        }
                        outboard_t_child {
                          position { x: 1 }
                          orientation { w: 1 }
                        }
                      )pb",
                      Pose3d(eigenmath::Quaterniond(0.7071067811865476, 0, 0,
                                                    0.7071067811865476),
                             eigenmath::Vector3d(1, 1, 0))},
        JointPoseCase{"PrismaticJoint", R"pb(
                        motion_type: MOTION_TYPE_PRISMATIC
                        axis { y: 1 }
                        raw_value: 0.5
                        parent_t_inboard {
                          position { x: 1 }
                          orientation { w: 1 }
                        }
                      )pb",
                      Pose3d(eigenmath::Vector3d(1, 0.5, 0))}),
    [](const ::testing::TestParamInfo<JointPoseCase>& info) {
      return info.param.name;
    });

TEST(SceneObjectViewTest, GetRootPosesFailsWhenJointPoseCannotBeResolved) {
  // `Create()` accepts a zero axis when the joint has no `parent_t_this`, but
  // the joint pose cannot be computed from it.
  const SceneObject scene_object = ParseTextProtoOrDie(R"pb(
    name: "zero_axis"
    entities {
      name: "base"
      link {}
    }
    entities {
      name: "joint"
      parent_name: "base"
      joint {
        kinematics_component {
          motion_type: MOTION_TYPE_REVOLUTE
          axis {}
          raw_value: 0.5
        }
      }
    }
  )pb");
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(scene_object));

  EXPECT_THAT(view.GetRootPoses(),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       AllOf(HasSubstr("'joint'"), HasSubstr("axis"))));
}

TEST(SceneObjectViewTest, SupportsWideTrees) {
  constexpr int kNumChildren = 1000;
  SceneObject scene_object;
  scene_object.set_name("wide");
  for (int i = kNumChildren - 1; i >= 0; --i) {
    AddEntity(scene_object, absl::StrCat("frame_", i), "root").mutable_frame();
  }
  AddEntity(scene_object, "root", "").mutable_link();
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(std::move(scene_object)));

  const absl::Span<const Entity* const> children = view.GetChildren("root");
  ASSERT_THAT(children, SizeIs(kNumChildren));
  // Names are compared lexicographically, not numerically.
  EXPECT_THAT(children.first(5),
              ElementsAre(EntityNamed("frame_0"), EntityNamed("frame_1"),
                          EntityNamed("frame_10"), EntityNamed("frame_100"),
                          EntityNamed("frame_101")));
  EXPECT_THAT(children.last(1), ElementsAre(EntityNamed("frame_999")));
  EXPECT_TRUE(std::ranges::is_sorted(Names(children)));
  EXPECT_THAT(view.GetEntitiesInTreeOrder().front(), EntityNamed("root"));
  EXPECT_THAT(view.GetEntitiesInTreeOrder().subspan(1),
              ElementsAreArray(children));
  for (const Entity* child : children) {
    EXPECT_THAT(view.GetParent(child->name()), EntityNamed("root"));
    EXPECT_THAT(view.GetChildren(child->name()), IsEmpty());
    EXPECT_THAT(view.GetSubtreeInTreeOrder(child->name()), ElementsAre(child));
    EXPECT_THAT(view.GetPathFromRoot(child->name()),
                ElementsAre(EntityNamed("root"), child));
  }
}

TEST(SceneObjectViewTest, OrdersChildrenByBytewiseNameComparison) {
  SceneObject scene_object;
  scene_object.set_name("names");
  AddEntity(scene_object, "root", "").mutable_link();
  for (const absl::string_view name :
       {"ab", "\xC3\xA9", "a_b", "Z", "a", "a0", "B"}) {
    AddEntity(scene_object, name, "root").mutable_frame();
  }
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(std::move(scene_object)));

  // Uppercase letters precede lowercase ones, a name precedes all names it is a
  // prefix of, and multi-byte UTF-8 characters follow ASCII.
  EXPECT_THAT(Names(view.GetChildren("root")),
              ElementsAre("B", "Z", "a", "a0", "a_b", "ab", "\xC3\xA9"));
  EXPECT_THAT(
      Names(view.GetEntitiesInTreeOrder()),
      ElementsAre("root", "B", "Z", "a", "a0", "a_b", "ab", "\xC3\xA9"));
}

// Robot arm with six revolute joints and a parallel gripper, mixing all entity
// types. Entities are intentionally not listed in tree order:
//
// base_link
// |-- joint_1 -> link_1 -> joint_2 -> link_2 -> joint_3
// |   |-- joint_3_torque (sensor)
// |   `-- link_3 -> joint_4 -> link_4 -> joint_5 -> link_5
// |       |-- joint_6 -> link_6
// |       |   |-- flange_frame
// |       |   `-- gripper_mount (fixed) -> gripper_base
// |       |       |-- finger_left_joint (prismatic) -> finger_left
// |       |       |-- finger_right_joint (prismatic) -> finger_right
// |       |       `-- grasp_frame -> approach_frame
// |       `-- wrist_camera (sensor)
// `-- mount_frame
SceneObject MakeRobotArmSceneObject() {
  constexpr KinematicsComponent::MotionType kRevolute =
      KinematicsComponent::MOTION_TYPE_REVOLUTE;
  constexpr KinematicsComponent::MotionType kPrismatic =
      KinematicsComponent::MOTION_TYPE_PRISMATIC;
  SceneObject scene_object;
  scene_object.set_name("robot_arm");
  AddEntity(scene_object, "approach_frame", "grasp_frame").mutable_frame();
  AddEntity(scene_object, "finger_right", "finger_right_joint").mutable_link();
  AddJoint(scene_object, "finger_right_joint", "gripper_base", kPrismatic);
  AddEntity(scene_object, "grasp_frame", "gripper_base").mutable_frame();
  AddJoint(scene_object, "finger_left_joint", "gripper_base", kPrismatic);
  AddEntity(scene_object, "finger_left", "finger_left_joint").mutable_link();
  AddEntity(scene_object, "gripper_base", "gripper_mount").mutable_link();
  AddJoint(scene_object, "gripper_mount", "link_6",
           KinematicsComponent::MOTION_TYPE_FIXED);
  AddEntity(scene_object, "flange_frame", "link_6").mutable_frame();
  AddEntity(scene_object, "wrist_camera", "link_5").mutable_sensor();
  AddEntity(scene_object, "joint_3_torque", "joint_3").mutable_sensor();
  for (int i = 6; i >= 1; --i) {
    AddEntity(scene_object, absl::StrCat("link_", i), absl::StrCat("joint_", i))
        .mutable_link();
    AddJoint(scene_object, absl::StrCat("joint_", i),
             i == 1 ? "base_link" : absl::StrCat("link_", i - 1), kRevolute);
  }
  AddEntity(scene_object, "mount_frame", "base_link").mutable_frame();
  AddEntity(scene_object, "base_link", "").mutable_link();
  return scene_object;
}

TEST(SceneObjectViewTest, SupportsRobotArmWithGripper) {
  const SceneObject scene_object = MakeRobotArmSceneObject();
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(scene_object));

  const std::vector<std::string> expected_tree_order = {
      "base_link",      "joint_1",
      "link_1",         "joint_2",
      "link_2",         "joint_3",
      "joint_3_torque", "link_3",
      "joint_4",        "link_4",
      "joint_5",        "link_5",
      "joint_6",        "link_6",
      "flange_frame",   "gripper_mount",
      "gripper_base",   "finger_left_joint",
      "finger_left",    "finger_right_joint",
      "finger_right",   "grasp_frame",
      "approach_frame", "wrist_camera",
      "mount_frame"};
  EXPECT_EQ(view.GetRootEntity().name(), "base_link");
  EXPECT_THAT(expected_tree_order, SizeIs(view.GetNumEntities()));
  EXPECT_THAT(Names(view.GetEntitiesInTreeOrder()),
              ElementsAreArray(expected_tree_order));

  // Children of different entity types are ordered by name only.
  EXPECT_THAT(Names(view.GetChildren("base_link")),
              ElementsAre("joint_1", "mount_frame"));
  EXPECT_THAT(Names(view.GetChildren("joint_3")),
              ElementsAre("joint_3_torque", "link_3"));
  EXPECT_THAT(Names(view.GetChildren("link_5")),
              ElementsAre("joint_6", "wrist_camera"));
  EXPECT_THAT(
      Names(view.GetChildren("gripper_base")),
      ElementsAre("finger_left_joint", "finger_right_joint", "grasp_frame"));
  EXPECT_THAT(view.GetParent("joint_3_torque"), EntityNamed("joint_3"));
  EXPECT_THAT(view.GetParent("approach_frame"), EntityNamed("grasp_frame"));

  EXPECT_THAT(Names(view.GetSubtreeInTreeOrder("gripper_base")),
              ElementsAre("gripper_base", "finger_left_joint", "finger_left",
                          "finger_right_joint", "finger_right", "grasp_frame",
                          "approach_frame"));
  EXPECT_THAT(Names(view.GetSubtreeInTreeOrder("link_5")),
              ElementsAreArray(expected_tree_order.begin() + 11,
                               expected_tree_order.end() - 1));
  EXPECT_THAT(Names(view.GetPathFromRoot("approach_frame")),
              ElementsAre("base_link", "joint_1", "link_1", "joint_2", "link_2",
                          "joint_3", "link_3", "joint_4", "link_4", "joint_5",
                          "link_5", "joint_6", "link_6", "gripper_mount",
                          "gripper_base", "grasp_frame", "approach_frame"));

  EXPECT_TRUE(view.IsInSubtree("finger_left", "joint_4"));
  EXPECT_TRUE(view.IsInSubtree("wrist_camera", "link_5"));
  EXPECT_FALSE(view.IsInSubtree("wrist_camera", "joint_6"));
  EXPECT_FALSE(view.IsInSubtree("mount_frame", "joint_1"));
  EXPECT_FALSE(view.IsInSubtree("joint_3_torque", "link_3"));
}

TEST(SceneObjectViewTest, SupportsDeepTrees) {
  constexpr int kNumLinks = 2000;
  // Entities are listed root-to-leaf and leaf-to-root.
  for (const bool leaf_to_root : {false, true}) {
    SCOPED_TRACE(absl::StrCat("leaf_to_root: ", leaf_to_root));
    SceneObject scene_object;
    scene_object.set_name("chain");
    for (int j = 0; j < kNumLinks; ++j) {
      const int i = leaf_to_root ? kNumLinks - 1 - j : j;
      AddEntity(scene_object, absl::StrCat("link_", i),
                i > 0 ? absl::StrCat("link_", i - 1) : "")
          .mutable_link();
    }
    ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                         SceneObjectView::Create(std::move(scene_object)));

    EXPECT_EQ(view.GetRootEntity().name(), "link_0");
    EXPECT_THAT(view.GetEntitiesInTreeOrder(), SizeIs(kNumLinks));
    EXPECT_THAT(view.GetSubtreeInTreeOrder("link_500"),
                SizeIs(kNumLinks - 500));
    EXPECT_TRUE(
        view.IsInSubtree(absl::StrCat("link_", kNumLinks - 1), "link_500"));
    EXPECT_FALSE(view.IsInSubtree("link_499", "link_500"));
    EXPECT_THAT(view.GetPathFromRoot(absl::StrCat("link_", kNumLinks - 1)),
                SizeIs(kNumLinks));
  }
}

TEST(SceneObjectViewTest, OwnsSceneObject) {
  SceneObject scene_object = MakeTreeSceneObject();
  const SceneObject expected = scene_object;
  ASSERT_OK_AND_ASSIGN(const SceneObjectView view,
                       SceneObjectView::Create(std::move(scene_object)));

  // The view is independent of the passed proto.
  scene_object = SceneObject();

  EXPECT_THAT(view.scene_object(), EqualsProto(expected));
  EXPECT_EQ(view.GetRootEntity().name(), "base");
  EXPECT_THAT(view.GetChildren("joint_a"), ElementsAre(EntityNamed("link_a")));
}

TEST(SceneObjectViewTest, RemainsValidAfterMove) {
  ASSERT_OK_AND_ASSIGN(SceneObjectView view,
                       SceneObjectView::Create(MakeTreeSceneObject()));
  const SceneObject* const proto = &view.scene_object();
  const Entity* const root = &view.GetRootEntity();
  const Entity* const joint_a = view.FindEntity("joint_a");
  const absl::Span<const Entity* const> base_children =
      view.GetChildren("base");
  const absl::Span<const Entity* const> tree_order =
      view.GetEntitiesInTreeOrder();
  const absl::Span<const Entity* const> subtree =
      view.GetSubtreeInTreeOrder("joint_b");

  SceneObjectView moved = std::move(view);

  EXPECT_EQ(&moved.scene_object(), proto);
  EXPECT_EQ(&moved.GetRootEntity(), root);
  EXPECT_EQ(moved.FindEntity("joint_a"), joint_a);
  EXPECT_THAT(base_children,
              ElementsAre(EntityNamed("base_frame"), EntityNamed("joint_a"),
                          EntityNamed("joint_b")));
  EXPECT_THAT(tree_order, SizeIs(8));
  EXPECT_THAT(subtree,
              ElementsAre(EntityNamed("joint_b"), EntityNamed("link_b"),
                          EntityNamed("tool_frame")));

  ASSERT_OK_AND_ASSIGN(
      SceneObjectView reassigned,
      SceneObjectView::Create(ParseTextProtoOrDie(R"pb(name: "other"
                                                       entities {
                                                         name: "other_root"
                                                         link {}
                                                       })pb")));
  reassigned = std::move(moved);

  EXPECT_EQ(&reassigned.scene_object(), proto);
  EXPECT_EQ(&reassigned.GetRootEntity(), root);
  EXPECT_EQ(reassigned.FindEntity("joint_a"), joint_a);
  EXPECT_EQ(reassigned.GetChildren("base").data(), base_children.data());
  EXPECT_EQ(reassigned.GetEntitiesInTreeOrder().data(), tree_order.data());
  EXPECT_EQ(reassigned.GetSubtreeInTreeOrder("joint_b").data(), subtree.data());
  EXPECT_THAT(base_children,
              ElementsAre(EntityNamed("base_frame"), EntityNamed("joint_a"),
                          EntityNamed("joint_b")));
  EXPECT_THAT(subtree,
              ElementsAre(EntityNamed("joint_b"), EntityNamed("link_b"),
                          EntityNamed("tool_frame")));
}

static_assert(!std::is_copy_constructible_v<SceneObjectView>);
static_assert(!std::is_copy_assignable_v<SceneObjectView>);
static_assert(std::is_nothrow_move_constructible_v<SceneObjectView>);

}  // namespace
}  // namespace scene_object
}  // namespace intrinsic
