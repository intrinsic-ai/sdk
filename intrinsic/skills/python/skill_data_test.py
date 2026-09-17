# Copyright 2026 Intrinsic Innovation LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import threading
from typing import Any

from absl.testing import absltest

from intrinsic.skills.python import skill_data


class SkillDataTest(absltest.TestCase):

  def test_default_construction_and_validation(self):
    sd = skill_data.SkillData()
    self.assertIsNotNone(sd)

    with self.assertRaises(ValueError):
      skill_data.SkillData(-1)

  def test_max_contexts_zero_does_not_cache(self):
    sd = skill_data.SkillData(0)
    compute_count = 0

    def compute_fn() -> int:
      nonlocal compute_count
      compute_count += 1
      return compute_count

    self.assertEqual(sd.get_or_compute("ctx_1", "key_1", compute_fn), 1)
    self.assertEqual(sd.get_or_compute("ctx_1", "key_1", compute_fn), 2)
    self.assertIsNone(sd.get("ctx_1", "key_1"))
    self.assertFalse(sd.delete("ctx_1"))

  def test_caching_none_value(self):
    sd = skill_data.SkillData(10)
    compute_count = 0

    def compute_fn() -> None:
      nonlocal compute_count
      compute_count += 1
      return None

    first = sd.get_or_compute("ctx_1", "key_1", compute_fn)
    self.assertIsNone(first)
    self.assertEqual(compute_count, 1)

    second = sd.get_or_compute("ctx_1", "key_1", compute_fn)
    self.assertIsNone(second)
    self.assertEqual(compute_count, 1)

  def test_cache_miss_computes_and_stores(self):
    sd = skill_data.SkillData(10)
    compute_count = 0

    def compute_fn() -> str:
      nonlocal compute_count
      compute_count += 1
      return "trajectory_data_123"

    result = sd.get_or_compute("ctx_1", "plan", compute_fn)
    self.assertEqual(result, "trajectory_data_123")
    self.assertEqual(compute_count, 1)

  def test_cache_hit_returns_cached_without_recomputing(self):
    sd = skill_data.SkillData(10)
    compute_count = 0

    def compute_fn() -> int:
      nonlocal compute_count
      compute_count += 1
      return 42

    first_result = sd.get_or_compute("ctx_1", "plan", compute_fn)
    self.assertEqual(first_result, 42)
    self.assertEqual(compute_count, 1)

    second_result = sd.get_or_compute("ctx_1", "plan", compute_fn)
    self.assertEqual(second_result, 42)
    self.assertEqual(compute_count, 1)

  def test_get_cache_miss_returns_none(self):
    sd = skill_data.SkillData(10)

    # Missing context and key
    result1 = sd.get("ctx_missing", "key_missing")
    self.assertIsNone(result1)

    # Existing context, missing key
    sd.get_or_compute("ctx_1", "key_1", lambda: 10)
    result2 = sd.get("ctx_1", "key_missing")
    self.assertIsNone(result2)

  def test_get_cache_hit_returns_value(self):
    sd = skill_data.SkillData(10)
    sd.get_or_compute("ctx_1", "plan", lambda: "trajectory_xyz")

    result = sd.get("ctx_1", "plan")
    self.assertEqual(result, "trajectory_xyz")

  def test_heterogeneous_python_objects(self):
    sd = skill_data.SkillData(10)
    test_dict = {"waypoints": [1, 2, 3], "status": "OK"}
    test_tuple = (10, "axis", 3.14)

    sd.get_or_compute("ctx_1", "config", lambda: test_dict)
    sd.get_or_compute("ctx_1", "params", lambda: test_tuple)

    self.assertIs(sd.get("ctx_1", "config"), test_dict)
    self.assertIs(sd.get("ctx_1", "params"), test_tuple)

  def test_get_validate_fn_accepts_and_rejects(self):
    sd = skill_data.SkillData(10)
    sd.get_or_compute("ctx_1", "key_1", lambda: 42)

    # Validator accepts
    accept_result = sd.get("ctx_1", "key_1", validate_fn=lambda v: v == 42)
    self.assertEqual(accept_result, 42)

    # Validator rejects -> returns None
    reject_result = sd.get("ctx_1", "key_1", validate_fn=lambda v: v != 42)
    self.assertIsNone(reject_result)

  def test_get_validate_fn_error_propagates_immediately(self):
    sd = skill_data.SkillData(10)
    sd.get_or_compute("ctx_1", "key_1", lambda: 42)

    def failing_validate(v: Any) -> bool:
      del v
      raise RuntimeError("validation check failed")

    with self.assertRaisesRegex(RuntimeError, "validation check failed"):
      sd.get("ctx_1", "key_1", validate_fn=failing_validate)

  def test_get_refreshes_lru_recency(self):
    max_contexts = 2
    sd = skill_data.SkillData(max_contexts)

    # Populate ctx_1 and ctx_2
    sd.get_or_compute("ctx_1", "k", lambda: 1)
    sd.get_or_compute("ctx_2", "k", lambda: 2)

    # Calling get on ctx_1 should refresh its LRU recency to most recent
    get_res = sd.get("ctx_1", "k")
    self.assertEqual(get_res, 1)

    # Insert ctx_3 -> Should evict ctx_2 (least recently used), NOT ctx_1!
    sd.get_or_compute("ctx_3", "k", lambda: 3)

    # ctx_1 should still be present
    ctx1_res = sd.get("ctx_1", "k")
    self.assertEqual(ctx1_res, 1)

    # ctx_2 should have been evicted
    ctx2_res = sd.get("ctx_2", "k")
    self.assertIsNone(ctx2_res)

  def test_multiple_keys_under_same_context(self):
    sd = skill_data.SkillData(10)

    sd.get_or_compute("ctx_1", "plan", lambda: "trajectory")
    sd.get_or_compute("ctx_1", "ik_seed", lambda: 7)
    sd.get_or_compute("ctx_1", "timestamp", lambda: 123.45)

    plan = sd.get_or_compute("ctx_1", "plan", lambda: "fresh")
    self.assertEqual(plan, "trajectory")

    ik = sd.get_or_compute("ctx_1", "ik_seed", lambda: 99)
    self.assertEqual(ik, 7)

    ts = sd.get_or_compute("ctx_1", "timestamp", lambda: 999.99)
    self.assertEqual(ts, 123.45)

  def test_validate_fn_accepts_cached(self):
    sd = skill_data.SkillData(10)
    compute_count = 0

    def compute_fn() -> int:
      nonlocal compute_count
      compute_count += 1
      return 100

    sd.get_or_compute("ctx_1", "key_1", compute_fn)
    self.assertEqual(compute_count, 1)

    result = sd.get_or_compute(
        "ctx_1", "key_1", compute_fn, validate_fn=lambda v: v == 100
    )
    self.assertEqual(result, 100)
    self.assertEqual(compute_count, 1)

  def test_validate_fn_rejects_and_recomputes(self):
    sd = skill_data.SkillData(10)
    value_to_return = 100
    compute_count = 0

    def compute_fn() -> int:
      nonlocal compute_count
      compute_count += 1
      return value_to_return

    sd.get_or_compute("ctx_1", "key_1", compute_fn)
    self.assertEqual(compute_count, 1)

    value_to_return = 200
    result = sd.get_or_compute(
        "ctx_1", "key_1", compute_fn, validate_fn=lambda v: v == 200
    )
    self.assertEqual(result, 200)
    self.assertEqual(compute_count, 2)

  def test_validate_fn_error_propagates_immediately(self):
    sd = skill_data.SkillData(10)
    sd.get_or_compute("ctx_1", "key_1", lambda: 100)

    def failing_validate(v: Any) -> bool:
      del v
      raise RuntimeError("Hardware state corrupted.")

    with self.assertRaisesRegex(RuntimeError, "Hardware state corrupted"):
      sd.get_or_compute(
          "ctx_1", "key_1", lambda: 100, validate_fn=failing_validate
      )

  def test_compute_fn_none_raises_error(self):
    sd = skill_data.SkillData(10)
    with self.assertRaises(ValueError):
      sd.get_or_compute("ctx_1", "key_1", None)

  def test_compute_fn_error_not_cached(self):
    sd = skill_data.SkillData(10)
    compute_count = 0

    def compute_fn() -> int:
      nonlocal compute_count
      compute_count += 1
      if compute_count == 1:
        raise RuntimeError("Planning solver timeout")
      return 555

    with self.assertRaisesRegex(RuntimeError, "Planning solver timeout"):
      sd.get_or_compute("ctx_1", "key_1", compute_fn)

    second_result = sd.get_or_compute("ctx_1", "key_1", compute_fn)
    self.assertEqual(second_result, 555)

  def test_delete_context_purges_all_keys_for_that_context(self):
    sd = skill_data.SkillData(10)
    compute_fn = lambda: 1

    sd.get_or_compute("ctx_1", "key_1", compute_fn)
    sd.get_or_compute("ctx_1", "key_2", compute_fn)
    sd.get_or_compute("ctx_2", "key_1", compute_fn)

    self.assertTrue(sd.delete("ctx_1"))
    self.assertFalse(sd.delete("ctx_1"))  # second delete returns False

    # ctx_1 keys should be evicted and recomputed
    ctx1_recomputed = 0

    def compute_fresh() -> int:
      nonlocal ctx1_recomputed
      ctx1_recomputed += 1
      return 77

    res_k1 = sd.get_or_compute("ctx_1", "key_1", compute_fresh)
    self.assertEqual(res_k1, 77)
    self.assertEqual(ctx1_recomputed, 1)

    # ctx_2 remains cached
    compute_count = 0

    def compute_never() -> int:
      nonlocal compute_count
      compute_count += 1
      return 999

    result = sd.get_or_compute("ctx_2", "key_1", compute_never)
    self.assertEqual(result, 1)
    self.assertEqual(compute_count, 0)

  def test_lru_eviction_capacity_bounded_by_context_count(self):
    max_contexts = 2
    sd = skill_data.SkillData(max_contexts)

    # Populate ctx_1 (with 2 keys) and ctx_2 (with 1 key)
    sd.get_or_compute("ctx_1", "key_a", lambda: 10)
    sd.get_or_compute("ctx_1", "key_b", lambda: 11)
    sd.get_or_compute("ctx_2", "key_a", lambda: 20)

    # Insert ctx_3 -> Should evict oldest context (ctx_1) and all its keys!
    sd.get_or_compute("ctx_3", "key_a", lambda: 30)

    # ctx_2 should still be cached without recomputing
    ctx2_recomputed = 0

    def compute_ctx2() -> int:
      nonlocal ctx2_recomputed
      ctx2_recomputed += 1
      return 200

    res_ctx2 = sd.get_or_compute("ctx_2", "key_a", compute_ctx2)
    self.assertEqual(res_ctx2, 20)
    self.assertEqual(ctx2_recomputed, 0)

    # ctx_1 keys should have been evicted and need recomputing
    ctx1_recomputed = 0

    def compute_fresh() -> int:
      nonlocal ctx1_recomputed
      ctx1_recomputed += 1
      return 100

    res_a = sd.get_or_compute("ctx_1", "key_a", compute_fresh)
    self.assertEqual(res_a, 100)
    self.assertEqual(ctx1_recomputed, 1)

  def test_lru_access_refreshes_context_order(self):
    max_contexts = 2
    sd = skill_data.SkillData(max_contexts)

    # Insert ctx_1 and ctx_2 (LRU order: ctx_1 oldest, ctx_2 newest)
    sd.get_or_compute("ctx_1", "key_a", lambda: 1)
    sd.get_or_compute("ctx_2", "key_a", lambda: 2)

    def fail_if_called():
      self.fail("compute_fn should not have been called on a cache hit")

    # Access ctx_1 -> Moves ctx_1 to newest. Now ctx_2 is oldest!
    hit_ctx1 = sd.get_or_compute("ctx_1", "key_a", fail_if_called)
    self.assertEqual(hit_ctx1, 1)

    # Insert ctx_3 -> Should evict ctx_2, keeping ctx_1!
    sd.get_or_compute("ctx_3", "key_a", lambda: 3)

    # ctx_1 should still be present
    check_ctx1 = sd.get_or_compute("ctx_1", "key_a", fail_if_called)
    self.assertEqual(check_ctx1, 1)

    # ctx_2 was evicted
    ctx2_recomputed = 0

    def compute_ctx2() -> int:
      nonlocal ctx2_recomputed
      ctx2_recomputed += 1
      return 200

    check_ctx2 = sd.get_or_compute("ctx_2", "key_a", compute_ctx2)
    self.assertEqual(check_ctx2, 200)
    self.assertEqual(ctx2_recomputed, 1)

  def test_concurrent_thread_access(self):
    sd = skill_data.SkillData(10)
    num_threads = 8
    operations_per_thread = 50
    errors = []

    def worker(t: int):
      try:
        for i in range(operations_per_thread):
          context_id = f"ctx_{t % 5}"
          key = f"key_{i % 10}"
          if i % 2 == 0:
            _ = sd.get(context_id, key)
          res = sd.get_or_compute(context_id, key, lambda: t * 1000 + i)
          self.assertIsInstance(res, int)
          if i % 7 == 0:
            sd.delete(context_id)
      except Exception as e:
        errors.append(e)

    threads = [
        threading.Thread(target=worker, args=(t,)) for t in range(num_threads)
    ]
    for thread in threads:
      thread.start()
    for thread in threads:
      thread.join()

    self.assertEmpty(errors)

  def test_get_skill_data_returns_same_instance(self):
    instance1 = skill_data.get_skill_data()
    instance2 = skill_data.get_skill_data()
    self.assertIs(instance1, instance2)

    res = skill_data.get_skill_data().get_or_compute(
        "global_ctx", "key", lambda: "singleton_data"
    )
    self.assertEqual(res, "singleton_data")
    skill_data.get_skill_data().delete("global_ctx")

  def test_get_with_empty_context_id_returns_none(self):
    sd = skill_data.SkillData(10)
    self.assertIsNone(sd.get("", "key_1"))

  def test_get_or_compute_with_empty_context_id_computes_without_caching(self):
    sd = skill_data.SkillData(10)
    compute_count = 0

    def compute_fn() -> int:
      nonlocal compute_count
      compute_count += 1
      return compute_count * 10

    # First computation returns 10
    first = sd.get_or_compute("", "plan", compute_fn)
    self.assertEqual(first, 10)
    self.assertEqual(compute_count, 1)

    # Second computation recomputes to 20 without caching
    second = sd.get_or_compute("", "plan", compute_fn)
    self.assertEqual(second, 20)
    self.assertEqual(compute_count, 2)

    # get with empty context_id returns None
    self.assertIsNone(sd.get("", "plan"))

    # Non-empty context_id caches normally
    normal_first = sd.get_or_compute("ctx_1", "plan", compute_fn)
    self.assertEqual(normal_first, 30)
    self.assertEqual(compute_count, 3)

    normal_second = sd.get_or_compute("ctx_1", "plan", compute_fn)
    self.assertEqual(normal_second, 30)
    self.assertEqual(compute_count, 3)

  def test_delete_with_empty_context_id_returns_false(self):
    sd = skill_data.SkillData(10)
    self.assertFalse(sd.delete(""))


if __name__ == "__main__":
  absltest.main()
