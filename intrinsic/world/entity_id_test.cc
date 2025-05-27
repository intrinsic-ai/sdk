// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/world/entity_id.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <map>
#include <set>
#include <tuple>
#include <type_traits>
#include <utility>
#include <vector>

#include "absl/hash/hash_testing.h"
#include "intrinsic/util/aggregate_type.h"
#include "intrinsic/util/testing/gtest_wrapper.h"
#include "intrinsic/world/hashing/hashing.h"

namespace intrinsic {
namespace {

using ::testing::Pair;
using ::testing::UnorderedElementsAre;

TEST(EntityId, SupportsAbslHash) {
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      EntityId(),
      EntityId(1),
      EntityId(2),
      EntityId(4),
  }));
}

TEST(EntityHandle, SupportsAbslHash) {
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      PhysicalEntityId(),
      PhysicalEntityId(1),
      PhysicalEntityId(2),
      PhysicalEntityId(4),
  }));
}

TEST(EntityHandle, SupportsTupleAbslHash) {
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      std::make_tuple(PhysicalEntityId(10), PhysicalEntityId(0)),
      std::make_tuple(PhysicalEntityId(9), PhysicalEntityId(1)),
      std::make_tuple(PhysicalEntityId(8), PhysicalEntityId(2)),
      std::make_tuple(PhysicalEntityId(7), PhysicalEntityId(4)),
      std::make_tuple(PhysicalEntityId(6), PhysicalEntityId(5)),
  }));
}

TEST(EntityHandle, HashSetWithTuple) {
  WorldHashSet<std::tuple<PhysicalEntityId, PhysicalEntityId>> set;
  set.insert(std::make_tuple(PhysicalEntityId(2), PhysicalEntityId(4)));
  EXPECT_THAT(set, ::testing::SizeIs(1));
}

TEST(EntityHandle, HashMapWithTuple) {
  WorldHashMap<std::tuple<PhysicalEntityId, PhysicalEntityId>, bool> map;
  map[std::make_tuple(PhysicalEntityId(2), PhysicalEntityId(4))] = true;
  EXPECT_THAT(map, ::testing::SizeIs(1));
}

TEST(EntityId, Constants) {
  EXPECT_GT(kFirstEntityId, kInvalidEntityId);
  EXPECT_GT(kFirstEntityId, kRootEntityId.id);
}

TEST(EntityHandle, TypeTests) {
  struct EmptyType {};

  struct ComplexType {
    int a;

   private:
    int b;
  };

  // using InvalidType = TypedEntityId<ComplexType>;
  EXPECT_FALSE(std::is_empty<ComplexType>::value);
  EXPECT_TRUE(std::is_empty<EmptyType>::value);
}

TEST(EntityHandle, ImplicitCast) {
  auto func1 = [](EntityId id) { return id.value(); };
  PhysicalEntityId entity(3);
  EXPECT_EQ(entity.id.value(), func1(entity));

  auto func2 = [](EntityHandle id) { return id.value(); };
  EXPECT_EQ(entity.id.value(), func2(entity));

  auto func3 = [](AttachmentEntityId id) { return id.value(); };
  EXPECT_EQ(entity.id.value(), func3(entity));
}

TEST(EntityHandle, ConstImplicitCast) {
  auto func1 = [](EntityId id) { return id.value(); };
  const PhysicalEntityId entity(3);
  EXPECT_EQ(entity.id.value(), func1(entity));

  auto func2 = [](EntityHandle id) { return id.value(); };
  EXPECT_EQ(entity.id.value(), func2(entity));

  auto func3 = [](AttachmentEntityId id) { return id.value(); };
  EXPECT_EQ(entity.id.value(), func3(entity));
}

TEST(EntityHandle, CopyConstruction) {
  EntityHandle entity(EntityId(1));
  EntityHandle other = entity;
  EXPECT_EQ(other, entity);
}

TEST(EntityHandle, Construction_With_EntityId) {
  EntityHandle entity(EntityId(1));
  RobotEntityId robot(EntityId(2));
  JointEntityId joint(EntityId(3));
  LinkEntityId link(EntityId(4));
  PhysicalEntityId object(EntityId(5));

  EXPECT_EQ(entity.id.value(), 1);
  EXPECT_EQ(robot.id.value(), 2);
  EXPECT_EQ(joint.id.value(), 3);
  EXPECT_EQ(link.id.value(), 4);
  EXPECT_EQ(object.id.value(), 5);

  EXPECT_EQ(entity.value(), 1);
  EXPECT_EQ(robot.value(), 2);
  EXPECT_EQ(joint.value(), 3);
  EXPECT_EQ(link.value(), 4);
  EXPECT_EQ(object.value(), 5);

  AttachmentEntityId attachment(EntityId(6));
  CollisionEntityId collision(EntityId(7));
  GeometryEntityId geometry(EntityId(8));
  KinematicsEntityId kinematics(EntityId(9));
  CollectionsEntityId collections(EntityId(10));
  CollectionsMemberEntityId collections_member(EntityId(11));
  PhysicsEntityId physics(EntityId(12));

  EXPECT_EQ(attachment.id.value(), 6);
  EXPECT_EQ(collision.id.value(), 7);
  EXPECT_EQ(geometry.id.value(), 8);
  EXPECT_EQ(kinematics.id.value(), 9);
  EXPECT_EQ(collections.id.value(), 10);
  EXPECT_EQ(collections_member.id.value(), 11);
  EXPECT_EQ(physics.id.value(), 12);

  EXPECT_EQ(attachment.value(), 6);
  EXPECT_EQ(collision.value(), 7);
  EXPECT_EQ(geometry.value(), 8);
  EXPECT_EQ(kinematics.value(), 9);
  EXPECT_EQ(collections.value(), 10);
  EXPECT_EQ(collections_member.value(), 11);
  EXPECT_EQ(physics.value(), 12);
}

TEST(EntityHandle, Construction_With_Value) {
  EntityHandle entity(1);
  RobotEntityId robot(2);
  JointEntityId joint(3);
  LinkEntityId link(4);
  PhysicalEntityId object(5);

  EXPECT_EQ(entity.value(), 1);
  EXPECT_EQ(robot.value(), 2);
  EXPECT_EQ(joint.value(), 3);
  EXPECT_EQ(link.value(), 4);
  EXPECT_EQ(object.value(), 5);

  AttachmentEntityId attachment(6);
  CollisionEntityId collision(7);
  GeometryEntityId geometry(8);
  KinematicsEntityId kinematics(9);
  PhysicsEntityId physics(11);
  CollectionsMemberEntityId collections_member(10);

  EXPECT_EQ(attachment.value(), 6);
  EXPECT_EQ(collision.value(), 7);
  EXPECT_EQ(geometry.value(), 8);
  EXPECT_EQ(kinematics.value(), 9);
  EXPECT_EQ(physics.value(), 11);
  EXPECT_EQ(collections_member.value(), 10);
}

TEST(EntityHandle, Assignment) {
  EntityHandle entity(1);
  RobotEntityId robot(2);
  JointEntityId joint(3);
  LinkEntityId link(4);
  PhysicalEntityId object(5);

  EXPECT_EQ(entity.value(), 1);
  entity = robot;
  EXPECT_EQ(entity.value(), 2);
  entity = joint;
  EXPECT_EQ(entity.value(), 3);
  entity = link;
  EXPECT_EQ(entity.value(), 4);
  entity = object;
  EXPECT_EQ(entity.value(), 5);

  AttachmentEntityId attachment(6);
  CollisionEntityId collision(7);
  GeometryEntityId geometry(8);
  KinematicsEntityId kinematics(9);
  CollectionsEntityId collections(10);
  CollectionsMemberEntityId collections_member(11);
  PhysicsEntityId physics(12);

  EXPECT_EQ(attachment.value(), 6);
  EXPECT_EQ(collision.value(), 7);
  EXPECT_EQ(geometry.value(), 8);
  EXPECT_EQ(kinematics.value(), 9);
  EXPECT_EQ(collections.value(), 10);
  EXPECT_EQ(collections_member.value(), 11);
  EXPECT_EQ(physics.value(), 12);

  attachment = joint;
  EXPECT_EQ(attachment.value(), 3);
  attachment = link;
  EXPECT_EQ(attachment.value(), 4);
  attachment = object;
  EXPECT_EQ(attachment.value(), 5);

  collision = link;
  EXPECT_EQ(collision.value(), 4);
  collision = object;
  EXPECT_EQ(collision.value(), 5);

  geometry = link;
  EXPECT_EQ(geometry.value(), 4);
  geometry = object;
  EXPECT_EQ(geometry.value(), 5);

  kinematics = joint;
  EXPECT_EQ(kinematics.value(), 3);

  physics = link;
  EXPECT_EQ(physics.value(), 4);

  collections_member = joint;
  EXPECT_EQ(collections_member.value(), 3);
  collections_member = link;
  EXPECT_EQ(collections_member.value(), 4);
}

TEST(EntityHandleContainers, Vector) {
  EntityHandle entity(1);
  RobotEntityId robot(2);
  JointEntityId joint(3);
  LinkEntityId link(4);
  PhysicalEntityId object(5);

  std::vector<EntityHandle> container;
  container.push_back(entity);
  container.emplace_back(entity);

  container.push_back(robot);
  container.emplace_back(robot);

  container.push_back(joint);
  container.emplace_back(joint);

  container.push_back(link);
  container.emplace_back(link);

  container.push_back(object);
  container.emplace_back(object);

  EXPECT_EQ(container.size(), 10);
}

TEST(EntityHandleContainers, Set) {
  EntityHandle entity(1);
  RobotEntityId robot(2);
  JointEntityId joint(3);
  LinkEntityId link(4);
  PhysicalEntityId object(5);

  std::set<EntityHandle> container;
  container.insert(entity);
  container.emplace(entity);

  container.insert(robot);
  container.emplace(robot);

  container.insert(joint);
  container.emplace(joint);

  container.insert(link);
  container.emplace(link);

  container.insert(object);
  container.emplace(object);

  EXPECT_THAT(container, UnorderedElementsAre(EntityHandle(1), EntityHandle(2),
                                              EntityHandle(3), EntityHandle(4),
                                              EntityHandle(5)));
}

TEST(EntityHandleContainers, HashSet) {
  EntityHandle entity(1);
  RobotEntityId robot(2);
  JointEntityId joint(3);
  LinkEntityId link(4);
  PhysicalEntityId object(5);

  WorldHashSet<EntityHandle> container;
  container.insert(entity);
  container.emplace(entity);

  container.insert(robot);
  container.emplace(robot);

  container.insert(joint);
  container.emplace(joint);

  container.insert(link);
  container.emplace(link);

  container.insert(object);
  container.emplace(object);

  container.emplace(link);
  container.emplace(joint);

  EXPECT_THAT(container, UnorderedElementsAre(EntityHandle(1), EntityHandle(2),
                                              EntityHandle(3), EntityHandle(4),
                                              EntityHandle(5)));
}

TEST(EntityHandleContainers, Map_Key) {
  EntityHandle entity(1);
  RobotEntityId robot(2);
  JointEntityId joint(3);
  LinkEntityId link(4);
  PhysicalEntityId object(5);

  std::map<EntityHandle, int> container;
  container.insert(std::pair(entity, 1));
  container.emplace(entity, 2);
  container[entity] = 3;

  container.insert(std::pair(robot, 1));
  container.emplace(robot, 4);
  container[robot] = 5;

  container.insert(std::pair(joint, 1));
  container.emplace(joint, 6);
  container[joint] = 7;

  container.insert(std::pair(link, 1));
  container.emplace(link, 8);
  container[link] = 9;

  container.insert(std::pair(object, 1));
  container[object] = 10;

  EXPECT_THAT(container, UnorderedElementsAre(
                             Pair(EntityHandle(1), 3), Pair(EntityHandle(2), 5),
                             Pair(EntityHandle(3), 7), Pair(EntityHandle(4), 9),
                             Pair(EntityHandle(5), 10)));
}

TEST(EntityHandleContainers, HashMap_Key) {
  EntityHandle entity(1);
  RobotEntityId robot(2);
  JointEntityId joint(3);
  LinkEntityId link(4);
  PhysicalEntityId object(5);

  WorldHashMap<EntityHandle, int> container;
  container.insert(std::pair(entity, 1));
  container.emplace(entity, 2);
  container[entity] = 3;

  container.insert(std::pair(robot, 1));
  container.emplace(robot, 4);
  container[robot] = 5;

  container.insert(std::pair(joint, 1));
  container.emplace(joint, 6);
  container[joint] = 7;

  container.insert(std::pair(link, 1));
  container.emplace(link, 8);
  container[link] = 9;

  container.insert(std::pair(object, 1));
  container[object] = 10;

  EXPECT_THAT(container, UnorderedElementsAre(
                             Pair(EntityHandle(1), 3), Pair(EntityHandle(2), 5),
                             Pair(EntityHandle(3), 7), Pair(EntityHandle(4), 9),
                             Pair(EntityHandle(5), 10)));
}

}  // namespace
}  // namespace intrinsic
