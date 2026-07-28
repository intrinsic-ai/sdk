// Copyright 2023 Intrinsic Innovation LLC

#include "third_party/strings_numbers/sort.h"

#include <gtest/gtest.h>

#include <algorithm>
#include <optional>
#include <string_view>
#include <vector>

#include "absl/container/btree_set.h"

namespace intrinsic {
namespace {

using ::testing::Values;

struct SortParam {
  std::string_view left;
  std::string_view right;
  int cmp;
  std::optional<int> strict_cmp;
};

class SortTest : public ::testing::TestWithParam<SortParam> {};

TEST_P(SortTest, TestAutoDigitStrCmp) {
  const SortParam param = GetParam();
  // Non-strict tests.
  EXPECT_EQ(AutoDigitLessThan(param.left, param.right), param.cmp < 0);
  EXPECT_EQ(AutoDigitStrCmp(param.left, param.right, false), param.cmp);
  // Strict tests.
  const int strict_cmp = param.strict_cmp.value_or(param.cmp);
  EXPECT_EQ(StrictAutoDigitLessThan(param.left, param.right), strict_cmp < 0);
  EXPECT_EQ(AutoDigitStrCmp(param.left, param.right, true), strict_cmp);
}

INSTANTIATE_TEST_SUITE_P(
    SortTest, SortTest,
    Values(SortParam{.left = "15", .right = "15", .cmp = 0},
           SortParam{.left = "", .right = "1", .cmp = -1},
           SortParam{.left = "2", .right = "3", .cmp = -1},
           SortParam{.left = "3", .right = "11", .cmp = -1},
           SortParam{.left = "03", .right = "011", .cmp = -1},
           SortParam{.left = "c", .right = "cd", .cmp = -1},
           SortParam{.left = "cd", .right = "d", .cmp = -1},
           SortParam{.left = "pre2", .right = "pre10", .cmp = -1},
           SortParam{.left = "pre2-ext", .right = "pre10-ext", .cmp = -1},
           // Leading zeroes in strict mode make a string smaller.
           SortParam{
               .left = "00005", .right = "5", .cmp = 0, .strict_cmp = -1}));

TEST(SortTest, TestStdSort) {
  const std::vector<std::string_view> input = {
      "file00", "file1",   "file02",    "file010",   "file20", "filea01",
      "filea1", "file022", "file00100", "file00001", "file1a"};

  const std::vector<std::string_view> ascending = {
      "file00", "file00001", "file1",     "file1a",  "file02", "file010",
      "file20", "file022",   "file00100", "filea01", "filea1"};
  const std::vector<std::string_view> descending{ascending.rbegin(),
                                                 ascending.rend()};
  {
    std::vector<std::string_view> sorted_input = input;
    std::sort(sorted_input.begin(), sorted_input.end(),
              strict_autodigit_less());
    EXPECT_EQ(ascending, sorted_input);
  }
  {
    std::vector<std::string_view> sorted_input = input;
    std::sort(sorted_input.begin(), sorted_input.end(),
              strict_autodigit_greater());
    EXPECT_EQ(descending, sorted_input);
  }
}

TEST(SortTest, TestStdStableSort) {
  const std::vector<std::string_view> input = {
      "file00", "file1",   "file02",    "file010",   "file20", "filea01",
      "filea1", "file022", "file00100", "file00001", "file1a"};
  const std::vector<std::string_view> ascending = {
      "file00", "file1",   "file00001", "file1a",  "file02", "file010",
      "file20", "file022", "file00100", "filea01", "filea1"};
  const std::vector<std::string_view> descending = {
      "filea01", "filea1", "file00100", "file022",   "file20", "file010",
      "file02",  "file1a", "file1",     "file00001", "file00"};
  {
    std::vector<std::string_view> sorted_input = input;
    std::stable_sort(sorted_input.begin(), sorted_input.end(),
                     autodigit_less());
    EXPECT_EQ(ascending, sorted_input);
  }
  {
    std::vector<std::string_view> sorted_input = input;
    std::stable_sort(sorted_input.begin(), sorted_input.end(),
                     autodigit_greater());
    EXPECT_EQ(descending, sorted_input);
  }
}

using AutodigitTypes =
    ::testing::Types<autodigit_less, autodigit_greater, strict_autodigit_less,
                     strict_autodigit_greater>;
template <typename T>
class AutodigitTypesSortTest : public ::testing::Test {};
TYPED_TEST_SUITE(AutodigitTypesSortTest, AutodigitTypes);

TYPED_TEST(AutodigitTypesSortTest, TestBtreeSet) {
  const absl::btree_set<std::string_view, TypeParam> file_set{"file3", "file4",
                                                              "file40"};
  constexpr std::string_view file3 = "file3";
  constexpr std::string_view file03 = "file03";

  EXPECT_TRUE(file_set.contains(file3));
  if (std::is_same_v<TypeParam, autodigit_less> ||
      std::is_same_v<TypeParam, autodigit_greater>) {
    EXPECT_TRUE(file_set.contains(file03));
  } else {
    EXPECT_FALSE(file_set.contains(file03));
  }
}

}  // namespace
}  // namespace intrinsic
