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
