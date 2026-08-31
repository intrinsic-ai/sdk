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

#include "intrinsic/skills/cc/skill_data.h"

#include <string>
#include <thread>
#include <vector>

#include "absl/status/status.h"
#include "absl/strings/str_format.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::skills {
namespace {

using ::testing::Eq;
using ::testing::HasSubstr;

TEST(SkillDataTest, CacheMissComputesAndStores) {
  SkillData skill_data(10);
  int compute_count = 0;

  auto compute_fn = [&]() -> absl::StatusOr<std::string> {
    ++compute_count;
    return "trajectory_data_123";
  };

  auto result =
      skill_data.GetOrCompute<std::string>("ctx_1", "plan", compute_fn);
  ASSERT_TRUE(result.ok()) << result.status();
  EXPECT_THAT(*result, Eq("trajectory_data_123"));
  EXPECT_THAT(compute_count, Eq(1));
}

TEST(SkillDataTest, CacheHitReturnsCachedWithoutRecomputing) {
  SkillData skill_data(10);
  int compute_count = 0;

  auto compute_fn = [&]() -> absl::StatusOr<int> {
    ++compute_count;
    return 42;
  };

  auto first_result = skill_data.GetOrCompute<int>("ctx_1", "plan", compute_fn);
  ASSERT_TRUE(first_result.ok());
  EXPECT_THAT(*first_result, Eq(42));
  EXPECT_THAT(compute_count, Eq(1));

  auto second_result =
      skill_data.GetOrCompute<int>("ctx_1", "plan", compute_fn);
  ASSERT_TRUE(second_result.ok());
  EXPECT_THAT(*second_result, Eq(42));
  EXPECT_THAT(compute_count, Eq(1));
}

TEST(SkillDataTest, MultipleKeysUnderSameContext) {
  SkillData skill_data(10);

  ASSERT_TRUE(
      skill_data
          .GetOrCompute<std::string>(
              "ctx_1", "plan",
              []() -> absl::StatusOr<std::string> { return "trajectory"; })
          .ok());
  ASSERT_TRUE(skill_data
                  .GetOrCompute<int>("ctx_1", "ik_seed",
                                     []() -> absl::StatusOr<int> { return 7; })
                  .ok());
  ASSERT_TRUE(skill_data
                  .GetOrCompute<double>(
                      "ctx_1", "timestamp",
                      []() -> absl::StatusOr<double> { return 123.45; })
                  .ok());

  auto plan = skill_data.GetOrCompute<std::string>(
      "ctx_1", "plan", []() -> absl::StatusOr<std::string> { return "fresh"; });
  ASSERT_TRUE(plan.ok());
  EXPECT_THAT(*plan, Eq("trajectory"));

  auto ik = skill_data.GetOrCompute<int>(
      "ctx_1", "ik_seed", []() -> absl::StatusOr<int> { return 99; });
  ASSERT_TRUE(ik.ok());
  EXPECT_THAT(*ik, Eq(7));
}

TEST(SkillDataTest, ValidateFnAcceptsCached) {
  SkillData skill_data(10);
  int compute_count = 0;

  auto compute_fn = [&]() -> absl::StatusOr<int> {
    ++compute_count;
    return 100;
  };

  auto validate_fn = [](const int& value) -> absl::StatusOr<bool> {
    return value == 100;
  };

  ASSERT_TRUE(skill_data.GetOrCompute<int>("ctx_1", "key_1", compute_fn).ok());
  EXPECT_THAT(compute_count, Eq(1));

  auto result =
      skill_data.GetOrCompute<int>("ctx_1", "key_1", compute_fn, validate_fn);
  ASSERT_TRUE(result.ok());
  EXPECT_THAT(*result, Eq(100));
  EXPECT_THAT(compute_count, Eq(1));
}

TEST(SkillDataTest, ValidateFnRejectsAndRecomputes) {
  SkillData skill_data(10);
  int value_to_return = 100;
  int compute_count = 0;

  auto compute_fn = [&]() -> absl::StatusOr<int> {
    ++compute_count;
    return value_to_return;
  };

  ASSERT_TRUE(skill_data.GetOrCompute<int>("ctx_1", "key_1", compute_fn).ok());
  EXPECT_THAT(compute_count, Eq(1));

  value_to_return = 200;
  auto validate_fn = [](const int& value) -> absl::StatusOr<bool> {
    return value == 200;
  };

  auto result =
      skill_data.GetOrCompute<int>("ctx_1", "key_1", compute_fn, validate_fn);
  ASSERT_TRUE(result.ok());
  EXPECT_THAT(*result, Eq(200));
  EXPECT_THAT(compute_count, Eq(2));
}

TEST(SkillDataTest, ValidateFnErrorPropagatesImmediately) {
  SkillData skill_data(10);

  auto compute_fn = []() -> absl::StatusOr<int> { return 100; };
  ASSERT_TRUE(skill_data.GetOrCompute<int>("ctx_1", "key_1", compute_fn).ok());

  auto validate_fn = [](const int&) -> absl::StatusOr<bool> {
    return absl::FailedPreconditionError("Hardware state corrupted.");
  };

  auto result =
      skill_data.GetOrCompute<int>("ctx_1", "key_1", compute_fn, validate_fn);
  EXPECT_TRUE(absl::IsFailedPrecondition(result.status()));
  EXPECT_THAT(result.status().message(), HasSubstr("Hardware state corrupted"));
}

TEST(SkillDataTest, ComputeFnNullReturnsError) {
  SkillData skill_data(10);
  auto result = skill_data.GetOrCompute<int>("ctx_1", "key_1", nullptr);
  EXPECT_TRUE(absl::IsInvalidArgument(result.status()));
}

TEST(SkillDataTest, ComputeFnErrorNotCached) {
  SkillData skill_data(10);
  int compute_count = 0;

  auto compute_fn = [&]() -> absl::StatusOr<int> {
    ++compute_count;
    if (compute_count == 1) {
      return absl::InternalError("Planning solver timeout");
    }
    return 555;
  };

  auto first_result =
      skill_data.GetOrCompute<int>("ctx_1", "key_1", compute_fn);
  EXPECT_TRUE(absl::IsInternal(first_result.status()));

  auto second_result =
      skill_data.GetOrCompute<int>("ctx_1", "key_1", compute_fn);
  ASSERT_TRUE(second_result.ok());
  EXPECT_THAT(*second_result, Eq(555));
}

TEST(SkillDataTest, TypeMismatchReturnsError) {
  SkillData skill_data(10);

  auto compute_string = []() -> absl::StatusOr<std::string> {
    return "some_string";
  };
  auto compute_int = []() -> absl::StatusOr<int> { return 123; };

  ASSERT_TRUE(
      skill_data.GetOrCompute<std::string>("ctx_1", "key_1", compute_string)
          .ok());

  auto result = skill_data.GetOrCompute<int>("ctx_1", "key_1", compute_int);
  EXPECT_TRUE(absl::IsInvalidArgument(result.status()));
  EXPECT_THAT(result.status().message(), HasSubstr("Type mismatch"));
}

TEST(SkillDataTest, DeleteContextPurgesAllKeysForThatContext) {
  SkillData skill_data(10);
  auto compute_fn = []() -> absl::StatusOr<int> { return 1; };

  ASSERT_TRUE(skill_data.GetOrCompute<int>("ctx_1", "key_1", compute_fn).ok());
  ASSERT_TRUE(skill_data.GetOrCompute<int>("ctx_1", "key_2", compute_fn).ok());
  ASSERT_TRUE(skill_data.GetOrCompute<int>("ctx_2", "key_1", compute_fn).ok());

  EXPECT_TRUE(skill_data.Delete("ctx_1"));
  EXPECT_FALSE(skill_data.Delete("ctx_1"));  // second delete returns false

  // ctx_1 keys should be evicted and recomputed
  int ctx1_recomputed = 0;
  auto compute_fresh = [&]() -> absl::StatusOr<int> {
    ++ctx1_recomputed;
    return 77;
  };
  auto res_k1 = skill_data.GetOrCompute<int>("ctx_1", "key_1", compute_fresh);
  ASSERT_TRUE(res_k1.ok());
  EXPECT_THAT(*res_k1, Eq(77));
  EXPECT_THAT(ctx1_recomputed, Eq(1));

  // ctx_2 remains cached
  int compute_count = 0;
  auto compute_never = [&]() -> absl::StatusOr<int> {
    ++compute_count;
    return 999;
  };
  auto result = skill_data.GetOrCompute<int>("ctx_2", "key_1", compute_never);
  ASSERT_TRUE(result.ok());
  EXPECT_THAT(*result, Eq(1));
  EXPECT_THAT(compute_count, Eq(0));
}

TEST(SkillDataTest, LruEvictionCapacityBoundedByContextCount) {
  constexpr size_t kMaxContexts = 2;
  SkillData skill_data(kMaxContexts);

  auto compute_val = [](int v) {
    return [v]() -> absl::StatusOr<int> { return v; };
  };

  // Populate ctx_1 (with 2 keys) and ctx_2 (with 1 key)
  ASSERT_TRUE(
      skill_data.GetOrCompute<int>("ctx_1", "key_a", compute_val(10)).ok());
  ASSERT_TRUE(
      skill_data.GetOrCompute<int>("ctx_1", "key_b", compute_val(11)).ok());
  ASSERT_TRUE(
      skill_data.GetOrCompute<int>("ctx_2", "key_a", compute_val(20)).ok());

  // Insert ctx_3 -> Should evict oldest context (ctx_1) and all its keys!
  ASSERT_TRUE(
      skill_data.GetOrCompute<int>("ctx_3", "key_a", compute_val(30)).ok());

  // ctx_2 should still be cached without recomputing
  int ctx2_recomputed = 0;
  auto compute_ctx2 = [&]() -> absl::StatusOr<int> {
    ++ctx2_recomputed;
    return 200;
  };
  auto res_ctx2 = skill_data.GetOrCompute<int>("ctx_2", "key_a", compute_ctx2);
  ASSERT_TRUE(res_ctx2.ok());
  EXPECT_THAT(*res_ctx2, Eq(20));
  EXPECT_THAT(ctx2_recomputed, Eq(0));

  // ctx_1 keys should have been evicted and need recomputing
  int ctx1_recomputed = 0;
  auto compute_fresh = [&]() -> absl::StatusOr<int> {
    ++ctx1_recomputed;
    return 100;
  };
  auto res_a = skill_data.GetOrCompute<int>("ctx_1", "key_a", compute_fresh);
  ASSERT_TRUE(res_a.ok());
  EXPECT_THAT(*res_a, Eq(100));
  EXPECT_THAT(ctx1_recomputed, Eq(1));
}

TEST(SkillDataTest, LruAccessRefreshesContextOrder) {
  constexpr size_t kMaxContexts = 2;
  SkillData skill_data(kMaxContexts);

  auto compute_val = [](int v) {
    return [v]() -> absl::StatusOr<int> { return v; };
  };

  // Insert ctx_1 and ctx_2 (LRU order: ctx_1 oldest, ctx_2 newest)
  ASSERT_TRUE(
      skill_data.GetOrCompute<int>("ctx_1", "key_a", compute_val(1)).ok());
  ASSERT_TRUE(
      skill_data.GetOrCompute<int>("ctx_2", "key_a", compute_val(2)).ok());

  // Access ctx_1 -> Moves ctx_1 to newest. Now ctx_2 is oldest!
  auto hit_ctx1 = skill_data.GetOrCompute<int>(
      "ctx_1", "key_a", []() -> absl::StatusOr<int> {
        return absl::InternalError("Not called");
      });
  ASSERT_TRUE(hit_ctx1.ok());
  EXPECT_THAT(*hit_ctx1, Eq(1));

  // Insert ctx_3 -> Should evict ctx_2, keeping ctx_1!
  ASSERT_TRUE(
      skill_data.GetOrCompute<int>("ctx_3", "key_a", compute_val(3)).ok());

  // ctx_1 should still be present
  auto check_ctx1 = skill_data.GetOrCompute<int>(
      "ctx_1", "key_a",
      []() -> absl::StatusOr<int> { return absl::InternalError("Evicted!"); });
  ASSERT_TRUE(check_ctx1.ok());
  EXPECT_THAT(*check_ctx1, Eq(1));

  // ctx_2 was evicted
  int ctx2_recomputed = 0;
  auto compute_ctx2 = [&]() -> absl::StatusOr<int> {
    ++ctx2_recomputed;
    return 200;
  };
  auto check_ctx2 =
      skill_data.GetOrCompute<int>("ctx_2", "key_a", compute_ctx2);
  ASSERT_TRUE(check_ctx2.ok());
  EXPECT_THAT(ctx2_recomputed, Eq(1));
}

TEST(SkillDataTest, ConcurrentThreadAccess) {
  SkillData skill_data(10);
  constexpr int kNumThreads = 8;
  constexpr int kOperationsPerThread = 50;

  std::vector<std::thread> threads;
  threads.reserve(kNumThreads);

  for (int t = 0; t < kNumThreads; ++t) {
    threads.emplace_back([&skill_data, t]() {
      for (int i = 0; i < kOperationsPerThread; ++i) {
        std::string context_id = absl::StrFormat("ctx_%d", t % 5);
        std::string key = absl::StrFormat("key_%d", i % 10);
        auto result = skill_data.GetOrCompute<int>(
            context_id, key,
            [t, i]() -> absl::StatusOr<int> { return t * 1000 + i; });
        EXPECT_TRUE(result.ok());

        if (i % 7 == 0) {
          skill_data.Delete(context_id);
        }
      }
    });
  }

  for (auto& thread : threads) {
    thread.join();
  }
}

TEST(SkillDataTest, GetSkillDataReturnsSameInstance) {
  SkillData& instance1 = GetSkillData();
  SkillData& instance2 = GetSkillData();
  EXPECT_EQ(&instance1, &instance2);

  auto res = GetSkillData().GetOrCompute<std::string>(
      "global_ctx", "key",
      []() -> absl::StatusOr<std::string> { return "singleton_data"; });
  ASSERT_TRUE(res.ok());
  EXPECT_THAT(*res, Eq("singleton_data"));

  GetSkillData().Delete("global_ctx");
}

}  // namespace
}  // namespace intrinsic::skills
