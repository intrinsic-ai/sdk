// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/utils/async_buffer.h"

#include <gtest/gtest.h>

#include <algorithm>
#include <cstdint>
#include <random>

#include "absl/strings/str_cat.h"
#include "intrinsic/icon/release/source_location.h"
#include "intrinsic/util/thread/thread.h"

namespace intrinsic::icon {
namespace {

class AsyncBufferTest : public testing::Test {
 protected:
  class Buffer {
   private:
    uint32_t seq_no_;
    uint32_t data_[4096 / 4 - 1];  // seq_no + data = 4K

   public:
    Buffer() { Fill(0); }

    explicit Buffer(uint32_t seq_no) { Fill(seq_no); }

    void Fill(uint32_t seq_no) {
      seq_no_ = seq_no;

      std::mt19937 rng(seq_no);
      for (uint32_t& value : data_) {
        value = rng();
      }
    }

    void Check(uint32_t seq_no, intrinsic::SourceLocation loc =
                                    intrinsic::SourceLocation::current()) {
      testing::ScopedTrace scope(loc.file_name(), loc.line(),
                                 absl::StrCat("Called: ", __func__));
      ASSERT_EQ(seq_no_, seq_no);

      std::mt19937 rng(seq_no);
      for (uint32_t value : data_) {
        ASSERT_EQ(value, rng());
      }
    }
  };

  static void VerifyActive(
      AsyncBuffer<Buffer>& async, uint32_t seq_no,
      intrinsic::SourceLocation loc = intrinsic::SourceLocation::current()) {
    testing::ScopedTrace scope(loc.file_name(), loc.line(),
                               absl::StrCat("Called: ", __func__));
    Buffer* active;
    async.GetActiveBuffer(&active);
    ASSERT_NE(active, nullptr);
    active->Check(seq_no);
  }
};

// Unit-tests for the AsyncBuffer<T> class.
TEST_F(AsyncBufferTest, FillCheck) {
  for (uint32_t i = 0; i < 0x1000; i++) {
    Buffer buff(i);
    buff.Check(i);
  }
}

TEST_F(AsyncBufferTest, WellOrdered) {
  AsyncBuffer<Buffer> async;

  Buffer* active = async.GetFreeBuffer();

  // Note: we use ASSERT here instead of EXPECT.  If any of these checks fail,
  // chances are all of them are going to fail, and there is no point in
  // spamming the test log with ~4000 reports that this test failed.
  ASSERT_TRUE(active != nullptr);
  VerifyActive(async, 0);

  for (uint32_t i = 1; i < 0x1000; i++) {
    Buffer* free_buff = async.GetFreeBuffer();
    free_buff->Fill(i);

    VerifyActive(async, i - 1);

    ASSERT_TRUE(async.CommitFreeBuffer());

    VerifyActive(async, i);
    VerifyActive(async, i);
  }
}

TEST_F(AsyncBufferTest, ReturnValueSemantics) {
  AsyncBuffer<Buffer> async;

  Buffer* active_buffer = nullptr;
  Buffer* other_buffer = nullptr;

  // the mailbox is empty; expect false return value
  ASSERT_FALSE(async.GetActiveBuffer(&active_buffer));
  ASSERT_NE(active_buffer, nullptr);

  // call not preceded by call to GetFreeBuffer();
  // expect false return value
  ASSERT_FALSE(async.CommitFreeBuffer());

  // nothing has changed; continue to expect false return value
  ASSERT_FALSE(async.GetActiveBuffer(&active_buffer));
  ASSERT_NE(active_buffer, nullptr);

  // commit a buffer into mailbox
  ASSERT_NE(async.GetFreeBuffer(), nullptr);
  ASSERT_TRUE(async.CommitFreeBuffer());

  // the mailbox is full; expect true return value
  ASSERT_TRUE(async.GetActiveBuffer(&other_buffer));
  ASSERT_NE(active_buffer, other_buffer);
  ASSERT_NE(other_buffer, nullptr);
}

TEST(AsyncBufferSimpleTest, GetLatest) {
  AsyncBuffer<int> buffer;
  int* free_buffer = buffer.GetFreeBuffer();
  EXPECT_NE(free_buffer, nullptr);
  *free_buffer = 2;
  buffer.CommitFreeBuffer();
  int* active_buffer = nullptr;
  EXPECT_TRUE(buffer.GetActiveBuffer(&active_buffer));
  EXPECT_EQ(*active_buffer, 2);
}

TEST(AsyncBufferSimpleTest, GetLatestAfterMultipleWrites) {
  AsyncBuffer<int> buffer;
  int* free_buffer = buffer.GetFreeBuffer();
  EXPECT_NE(free_buffer, nullptr);
  *free_buffer = 2;
  buffer.CommitFreeBuffer();
  free_buffer = buffer.GetFreeBuffer();
  EXPECT_NE(free_buffer, nullptr);
  *free_buffer = 3;
  buffer.CommitFreeBuffer();
  int* active_buffer = nullptr;
  EXPECT_TRUE(buffer.GetActiveBuffer(&active_buffer));
  EXPECT_EQ(*active_buffer, 3);
}

TEST(AsyncBufferSimpleTest, EmptyReadsDefaultValue) {
  struct TestStruct {
    int value = 2;
  };
  AsyncBuffer<TestStruct> buffer;
  TestStruct* result = nullptr;
  EXPECT_FALSE(buffer.GetActiveBuffer(&result));
  EXPECT_NE(result, nullptr);
  EXPECT_EQ(result->value, 2);
}

TEST(AsyncBufferSimpleTest, ThreadSafe) {
  struct TestStruct {
    int i = 0;
    double a = 3.14;
  };
  AsyncBuffer<TestStruct> buffer;
  Thread write_thread([&]() {
    for (int i = 0; i < 1000; ++i) {
      TestStruct* b = buffer.GetFreeBuffer();
      b->i = i;
      buffer.CommitFreeBuffer();
    }
  });
  int largest_i = 0;
  for (int i = 0; i < 1000; ++i) {
    TestStruct* s = nullptr;
    bool new_data = buffer.GetActiveBuffer(&s);
    EXPECT_GE(s->i, 0);
    EXPECT_LT(s->i, 1000);
    EXPECT_EQ(s->a, 3.14);
    if (new_data) {
      EXPECT_GT(s->i, largest_i)
          << "Producer writes increasing values so new data should be larger.";
      largest_i = std::max(largest_i, s->i);
    } else {
      EXPECT_EQ(s->i, largest_i) << "Consumer reads latest value.";
    }
    largest_i = std::max(largest_i, s->i);
  }
  write_thread.join();
  TestStruct* s = nullptr;
  buffer.GetActiveBuffer(&s);
  EXPECT_EQ(s->i, 999);
}

}  // namespace
}  // namespace intrinsic::icon
