// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_WORLD_HASHING_HASHING_TEST_CC_
#define INTRINSIC_WORLD_HASHING_HASHING_TEST_CC_

#include "intrinsic/world/hashing/hashing.h"

#include <gtest/gtest.h>

#include "intrinsic/world/entity_id.h"
#include "intrinsic/world/labels.h"

namespace intrinsic {
namespace {

TEST(EntityIdHasher, Works) {
  WorldHasher<EntityId> hasher;
  EXPECT_EQ(hasher(EntityId(1)), hasher(EntityId(1)));
  EXPECT_NE(hasher(EntityId(1)), hasher(EntityId(2)));
}

TEST(AttachmentEntityIdHasher, Works) {
  WorldHasher<AttachmentEntityId> hasher;
  EXPECT_EQ(hasher(AttachmentEntityId(1)), hasher(AttachmentEntityId(1)));
  EXPECT_NE(hasher(AttachmentEntityId(1)), hasher(AttachmentEntityId(2)));
}

TEST(LabelIdHasher, Works) {
  WorldHasher<LabelId> hasher;
  EXPECT_EQ(hasher(LabelId("label1")), hasher(LabelId("label1")));
  EXPECT_NE(hasher(LabelId("label1")), hasher(LabelId("label2")));
}

}  // namespace
}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_HASHING_HASHING_TEST_CC_
