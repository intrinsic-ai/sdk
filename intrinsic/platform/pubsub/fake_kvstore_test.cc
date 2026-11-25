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

#include "intrinsic/platform/pubsub/fake_kvstore.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <atomic>
#include <cstddef>
#include <cstdint>
#include <cstring>
#include <memory>
#include <optional>
#include <string>
#include <thread>
#include <type_traits>
#include <utility>
#include <vector>

#include "absl/container/flat_hash_map.h"
#include "absl/flags/flag.h"
#include "absl/flags/reflection.h"
#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/notification.h"
#include "absl/time/time.h"
#include "google/protobuf/any.pb.h"
#include "google/protobuf/wrappers.pb.h"
#include "intrinsic/middleware/buffers/buffer_factory.h"
#include "intrinsic/middleware/buffers/immutable_buffer.h"
#include "intrinsic/middleware/buffers/mutable_buffer.h"
#include "internal/testing.h"
#include "intrinsic/platform/pubsub/kvstore.h"
#include "intrinsic/platform/pubsub/subscription.h"
#include "intrinsic/platform/pubsub/topic_config.h"

namespace intrinsic {
namespace {

using ::absl_testing::IsOkAndHolds;
using ::absl_testing::StatusIs;
using ::google::protobuf::StringValue;
using ::testing::ElementsAre;
using ::testing::HasSubstr;
using ::testing::IsEmpty;
using ::testing::Key;
using ::testing::Pair;
using ::testing::Test;
using ::testing::UnorderedElementsAre;

template <typename Store, typename Arg>
concept CanSetRaw = requires(Store& store, Arg&& arg) {
  store.SetRaw("k", std::forward<Arg>(arg));
};

template <typename Store, typename Arg>
concept CanSetRawWithVerification =
    requires(Store& store, Arg&& arg,
             const KeyValueStore::SetWithVerificationOptions& opts) {
      store.SetRawWithVerification("k", std::forward<Arg>(arg), opts);
    };

// Compile-time checks verifying that SetRaw and SetRawWithVerification
// only accept rvalue MutableBuffer arguments and reject lvalues, while
// accepting ImmutableBuffer by value (lvalue or rvalue).
template <typename Store>
consteval void CheckRawOverloads() {
  static_assert(CanSetRaw<Store, MutableBuffer>);
  static_assert(!CanSetRaw<Store, MutableBuffer&>);
  static_assert(!CanSetRaw<Store, const MutableBuffer&>);
  static_assert(!CanSetRaw<Store, const MutableBuffer>);
  static_assert(CanSetRaw<Store, ImmutableBuffer>);
  static_assert(CanSetRaw<Store, ImmutableBuffer&>);
  static_assert(CanSetRaw<Store, const ImmutableBuffer&>);
  static_assert(CanSetRaw<Store, const ImmutableBuffer>);

  static_assert(CanSetRawWithVerification<Store, MutableBuffer>);
  static_assert(!CanSetRawWithVerification<Store, MutableBuffer&>);
  static_assert(!CanSetRawWithVerification<Store, const MutableBuffer&>);
  static_assert(!CanSetRawWithVerification<Store, const MutableBuffer>);
  static_assert(CanSetRawWithVerification<Store, ImmutableBuffer>);
  static_assert(CanSetRawWithVerification<Store, ImmutableBuffer&>);
  static_assert(CanSetRawWithVerification<Store, const ImmutableBuffer&>);
  static_assert(CanSetRawWithVerification<Store, const ImmutableBuffer>);
}

static_assert((CheckRawOverloads<KeyValueStore>(),
               CheckRawOverloads<FakeKeyValueStore>(), true));

absl::StatusOr<MutableBuffer> MakeMutableBuffer(absl::string_view data) {
  absl::StatusOr<MutableBuffer> buf = AllocateMutableBuffer(data.size());
  if (!buf.ok()) {
    return buf.status();
  }
  if (!data.empty()) {
    std::memcpy(buf->data(), data.data(), data.size());
  }
  return buf;
}

absl::StatusOr<ImmutableBuffer> MakeImmutableBuffer(absl::string_view data) {
  absl::StatusOr<MutableBuffer> buf = MakeMutableBuffer(data);
  if (!buf.ok()) {
    return buf.status();
  }
  return FreezeMutableBuffer(std::move(*buf));
}

absl::string_view AsStringView(const ImmutableBuffer& buf) {
  return absl::string_view(reinterpret_cast<const char*>(buf.data()),
                           buf.size());
}

class FakeKeyValueStoreTest : public Test {
 protected:
  FakeKeyValueStore kv_store_;
};

TEST_F(FakeKeyValueStoreTest, SetAndGetSuccess) {
  StringValue val;
  val.set_value("test_value");

  EXPECT_OK(kv_store_.Set("my/key", val));

  ASSERT_OK_AND_ASSIGN(StringValue retrieved,
                       kv_store_.Get<StringValue>("my/key"));
  EXPECT_EQ(retrieved.value(), "test_value");
}

TEST_F(FakeKeyValueStoreTest, SetWithVerificationAndGetSuccess) {
  StringValue val;
  val.set_value("test_value");

  EXPECT_OK(kv_store_.SetWithVerification("my/key", val, {}));

  ASSERT_OK_AND_ASSIGN(StringValue retrieved,
                       kv_store_.Get<StringValue>("my/key"));
  EXPECT_EQ(retrieved.value(), "test_value");
}

TEST_F(FakeKeyValueStoreTest, GetNotFound) {
  EXPECT_THAT(kv_store_.Get<StringValue>("nonexistent"),
              StatusIs(absl::StatusCode::kNotFound,
                       HasSubstr("Key not found: nonexistent")));
}

TEST_F(FakeKeyValueStoreTest, DeleteSuccess) {
  StringValue val;
  val.set_value("test_value");

  EXPECT_OK(kv_store_.Set("my/key", val));
  EXPECT_OK(kv_store_.Delete("my/key"));

  EXPECT_THAT(kv_store_.Get<StringValue>("my/key"),
              StatusIs(absl::StatusCode::kNotFound));
}

TEST_F(FakeKeyValueStoreTest, DeleteNotFound) {
  EXPECT_THAT(kv_store_.Delete("nonexistent"),
              StatusIs(absl::StatusCode::kNotFound,
                       HasSubstr("Key not found: nonexistent")));
}

TEST_F(FakeKeyValueStoreTest, KeyPrefixDefaultsToKvStore) {
  EXPECT_EQ(kv_store_.key_prefix(), "kv_store");
}

// FakeKeyValueStore forwards an explicit prefix_override to KeyValueStore so
// callers can exercise non-default storage prefixes in unit tests.
TEST(FakeKeyValueStorePrefixOverrideTest, KeyPrefixReturnsConfiguredOverride) {
  const FakeKeyValueStore custom_store("in_mem_kv_store");

  EXPECT_EQ(custom_store.key_prefix(), "in_mem_kv_store");
}

TEST_F(FakeKeyValueStoreTest, SetRawAndGetRawRoundTripPreservesAlignedPayload) {
  ASSERT_OK_AND_ASSIGN(MutableBuffer buf, AllocateMutableBuffer(128));
  EXPECT_TRUE(buf.valid());
  EXPECT_EQ(buf.size(), 128);
  EXPECT_EQ(buf.alignment(), 256);
  EXPECT_EQ(reinterpret_cast<uintptr_t>(buf.data()) % 256, 0);

  std::memset(buf.data(), 0xAB, 128);

  EXPECT_OK(kv_store_.SetRaw("raw/data", std::move(buf)));
  // Moving buf into SetRaw consumes it.
  EXPECT_FALSE(buf.valid());

  ASSERT_OK_AND_ASSIGN(ImmutableBuffer retrieved, kv_store_.GetRaw("raw/data"));
  EXPECT_TRUE(retrieved.valid());
  EXPECT_EQ(retrieved.size(), 128);
  EXPECT_EQ(retrieved.alignment(), 256);
  EXPECT_EQ(reinterpret_cast<uintptr_t>(retrieved.data()) % 256, 0);

  const uint8_t* bytes = retrieved.data();
  for (size_t i = 0; i < 128; ++i) {
    EXPECT_EQ(bytes[i], 0xAB);
  }
}

TEST_F(FakeKeyValueStoreTest, SetRawWithVerificationSuccess) {
  ASSERT_OK_AND_ASSIGN(MutableBuffer buf, AllocateMutableBuffer(64));
  std::memset(buf.data(), 0x42, 64);

  KeyValueStore::SetWithVerificationOptions opts{.timeout = absl::Seconds(5)};
  EXPECT_OK(
      kv_store_.SetRawWithVerification("raw/verified", std::move(buf), opts));

  ASSERT_OK_AND_ASSIGN(ImmutableBuffer retrieved,
                       kv_store_.GetRaw("raw/verified"));
  EXPECT_EQ(retrieved.size(), 64);
  EXPECT_EQ(retrieved.data()[0], 0x42);
}

TEST_F(FakeKeyValueStoreTest, ZeroByteRawBufferSupported) {
  ASSERT_OK_AND_ASSIGN(MutableBuffer buf, AllocateMutableBuffer(0));
  EXPECT_TRUE(buf.valid());
  EXPECT_EQ(buf.size(), 0);
  EXPECT_NE(buf.data(), nullptr);
  EXPECT_EQ(buf.alignment(), 256);

  EXPECT_OK(kv_store_.SetRaw("raw/zero", std::move(buf)));

  ASSERT_OK_AND_ASSIGN(ImmutableBuffer retrieved, kv_store_.GetRaw("raw/zero"));
  EXPECT_TRUE(retrieved.valid());
  EXPECT_EQ(retrieved.size(), 0);
  EXPECT_NE(retrieved.data(), nullptr);
  EXPECT_EQ(retrieved.alignment(), 256);
}

TEST_F(FakeKeyValueStoreTest, InvalidInputsReturnInvalidArgument) {
  // Invalid MutableBuffer.
  MutableBuffer invalid_buf;
  EXPECT_FALSE(invalid_buf.valid());
  EXPECT_THAT(kv_store_.SetRaw("valid/key", std::move(invalid_buf)),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // Invalid keys preserve buffer validity across SetRaw and
  // SetRawWithVerification, and fail for Set, Delete, and GetRaw.
  ASSERT_OK_AND_ASSIGN(MutableBuffer valid_mut, MakeMutableBuffer("payload8"));
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer valid_imm_buf,
                       FreezeMutableBuffer(std::move(valid_mut)));
  for (absl::string_view bad_key : {"", "/bad", "bad/", "bad//key"}) {
    ASSERT_OK_AND_ASSIGN(MutableBuffer buf, MakeMutableBuffer("payload8"));
    EXPECT_THAT(kv_store_.SetRaw(bad_key, std::move(buf)),
                StatusIs(absl::StatusCode::kInvalidArgument));
    EXPECT_TRUE(buf.valid());

    EXPECT_THAT(kv_store_.SetRawWithVerification(
                    bad_key, std::move(buf),
                    {.mode = KeyValueStore::SetWithVerificationOptions::
                         VerificationMode::kHighConsistency,
                     .timeout = absl::Seconds(1)}),
                StatusIs(absl::StatusCode::kInvalidArgument));
    EXPECT_TRUE(buf.valid());
    EXPECT_THAT(kv_store_.Set(bad_key, google::protobuf::Any()),
                StatusIs(absl::StatusCode::kInvalidArgument));
    EXPECT_THAT(kv_store_.SetRaw(bad_key, valid_imm_buf),
                StatusIs(absl::StatusCode::kInvalidArgument));
    EXPECT_THAT(kv_store_.Delete(bad_key),
                StatusIs(absl::StatusCode::kInvalidArgument));
    EXPECT_THAT(kv_store_.GetRaw(bad_key),
                StatusIs(absl::StatusCode::kInvalidArgument));
  }

  // Negative timeout on GetRaw.
  EXPECT_THAT(kv_store_.GetRaw("some/key", absl::Seconds(-1)),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // Non-positive timeout on SetRawWithVerification preserves buffer validity.
  ASSERT_OK_AND_ASSIGN(MutableBuffer b4, MakeMutableBuffer("payload8"));
  KeyValueStore::SetWithVerificationOptions non_pos_opts{
      .timeout = absl::ZeroDuration()};
  EXPECT_THAT(kv_store_.SetRawWithVerification("valid/key", std::move(b4),
                                               non_pos_opts),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_TRUE(b4.valid());

  // Non-positive timeout on SetWithVerification.
  StringValue pv;
  pv.set_value("v");
  EXPECT_THAT(kv_store_.SetWithVerification("valid/key", pv, non_pos_opts),
              StatusIs(absl::StatusCode::kInvalidArgument));
  KeyValueStore::SetWithVerificationOptions neg_opts{.timeout =
                                                         -absl::Seconds(1)};
  EXPECT_THAT(kv_store_.SetWithVerification("valid/key", pv, neg_opts),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // Negative and zero timeout boundaries on Get, GetRaw, GetAllSynchronous,
  // and GetAllRawSynchronous.
  EXPECT_THAT(kv_store_.Get<StringValue>("some/key", -absl::Seconds(1)),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(kv_store_.Get<StringValue>("some/key", absl::ZeroDuration()),
              StatusIs(absl::StatusCode::kNotFound));
  EXPECT_THAT(kv_store_.GetRaw("some/key", absl::ZeroDuration()),
              StatusIs(absl::StatusCode::kNotFound));
  EXPECT_THAT(kv_store_.GetAllSynchronous("some/**", -absl::Seconds(1)),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(kv_store_.GetAllRawSynchronous("some/**", -absl::Seconds(1)),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(kv_store_.GetAllSynchronous("some/**", absl::ZeroDuration()),
              IsOkAndHolds(IsEmpty()));
  EXPECT_THAT(kv_store_.GetAllRawSynchronous("some/**", absl::ZeroDuration()),
              IsOkAndHolds(IsEmpty()));
}

TEST_F(FakeKeyValueStoreTest, SetRawImmutableBufferSharesStorageAcrossKeys) {
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer frozen,
                       MakeImmutableBuffer("0123456789"));

  EXPECT_OK(kv_store_.SetRaw("fanout/k1", frozen));
  EXPECT_OK(kv_store_.SetRaw("fanout/k2", frozen));
  EXPECT_TRUE(frozen.valid());

  ASSERT_OK_AND_ASSIGN(ImmutableBuffer read1, kv_store_.GetRaw("fanout/k1"));
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer read2, kv_store_.GetRaw("fanout/k2"));
  EXPECT_EQ(read1.data(), frozen.data());
  EXPECT_EQ(read2.data(), frozen.data());
  EXPECT_EQ(AsStringView(read1), "0123456789");
}

TEST_F(FakeKeyValueStoreTest, SetRawRealignsUnalignedImmutableBufferSlice) {
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer frozen,
                       MakeImmutableBuffer("0123456789"));
  ImmutableBuffer slice = frozen.Slice(2, 4);
  ASSERT_TRUE(slice.valid());
  EXPECT_LT(slice.alignment(), ImmutableBuffer::kAlignment);

  EXPECT_OK(kv_store_.SetRaw("fanout/slice", slice));
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer read_slice,
                       kv_store_.GetRaw("fanout/slice"));
  EXPECT_EQ(read_slice.alignment(), ImmutableBuffer::kAlignment);
  EXPECT_EQ(AsStringView(read_slice), "2345");
}

TEST_F(FakeKeyValueStoreTest,
       SetRawWithVerificationMovesImmutableBufferWithoutCopying) {
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer frozen,
                       MakeImmutableBuffer("0123456789"));
  const uint8_t* expected_ptr = frozen.data();
  KeyValueStore::SetWithVerificationOptions valid_opts{.timeout =
                                                           absl::Seconds(5)};

  EXPECT_OK(kv_store_.SetRawWithVerification("fanout/moved", std::move(frozen),
                                             valid_opts));
  EXPECT_FALSE(frozen.valid());
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer read_moved,
                       kv_store_.GetRaw("fanout/moved"));
  EXPECT_EQ(read_moved.data(), expected_ptr);
  EXPECT_EQ(AsStringView(read_moved), "0123456789");
}

TEST_F(FakeKeyValueStoreTest, SetRawRoutesThroughBaseKeyValueStoreReference) {
  KeyValueStore& base_ref = kv_store_;
  ASSERT_OK_AND_ASSIGN(MutableBuffer poly_mut, MakeMutableBuffer("abc"));
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer poly_immut, MakeImmutableBuffer("2345"));
  KeyValueStore::SetWithVerificationOptions valid_opts{.timeout =
                                                           absl::Seconds(5)};

  EXPECT_OK(base_ref.SetRaw("poly/mut", std::move(poly_mut)));
  EXPECT_FALSE(poly_mut.valid());
  EXPECT_OK(base_ref.SetRaw("poly/immut", poly_immut));
  EXPECT_OK(
      base_ref.SetRawWithVerification("poly/immut_v", poly_immut, valid_opts));

  ASSERT_OK_AND_ASSIGN(ImmutableBuffer poly_read,
                       kv_store_.GetRaw("poly/immut_v"));
  EXPECT_EQ(AsStringView(poly_read), "2345");
}

TEST_F(FakeKeyValueStoreTest,
       SetRawImmutableBufferRejectsInvalidBufferAndTimeouts) {
  ImmutableBuffer invalid_immut;
  KeyValueStore::SetWithVerificationOptions valid_opts{.timeout =
                                                           absl::Seconds(5)};
  EXPECT_THAT(kv_store_.SetRaw("guard/invalid", invalid_immut),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(kv_store_.SetRawWithVerification("guard/invalid", invalid_immut,
                                               valid_opts),
              StatusIs(absl::StatusCode::kInvalidArgument));

  ASSERT_OK_AND_ASSIGN(ImmutableBuffer valid_immut,
                       MakeImmutableBuffer("2345"));
  KeyValueStore::SetWithVerificationOptions bad_opts{.timeout =
                                                         absl::ZeroDuration()};
  EXPECT_THAT(kv_store_.SetRawWithVerification("guard/bad_timeout", valid_immut,
                                               bad_opts),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_TRUE(valid_immut.valid());

  bad_opts.timeout = -absl::Seconds(1);
  EXPECT_THAT(kv_store_.SetRawWithVerification("guard/bad_timeout", valid_immut,
                                               bad_opts),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_TRUE(valid_immut.valid());
}

TEST_F(FakeKeyValueStoreTest, SetRawZeroByteImmutableBufferSucceeds) {
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer zero_immut, MakeImmutableBuffer(""));
  EXPECT_OK(kv_store_.SetRaw("guard/zero", std::move(zero_immut)));
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer read_zero,
                       kv_store_.GetRaw("guard/zero"));
  EXPECT_TRUE(read_zero.valid());
  EXPECT_EQ(read_zero.size(), 0);
}

TEST_F(FakeKeyValueStoreTest, ProtobufReadsAndAdminCopyRejectRawKey) {
  ASSERT_OK_AND_ASSIGN(MutableBuffer buf, AllocateMutableBuffer(16));
  EXPECT_OK(kv_store_.SetRaw("mixed/key", std::move(buf)));

  EXPECT_THAT(
      kv_store_.Get<StringValue>("mixed/key"),
      StatusIs(absl::StatusCode::kInvalidArgument,
               HasSubstr("Key 'kv_store/mixed/key' stores raw bytes, not a "
                         "Protobuf message")));
  EXPECT_THAT(kv_store_.Get<google::protobuf::Any>("mixed/key"),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("Key 'kv_store/mixed/key' stores raw bytes")));
  EXPECT_THAT(
      kv_store_.AdminCloudCopy("mixed/key", "mixed/dest", absl::Seconds(1)),
      StatusIs(absl::StatusCode::kInvalidArgument,
               HasSubstr("Key 'kv_store/mixed/key' stores raw bytes")));
}

TEST_F(FakeKeyValueStoreTest, GetRawRejectsProtobufKey) {
  StringValue proto_val;
  proto_val.set_value("now_proto");
  EXPECT_OK(kv_store_.Set("mixed/key", proto_val));

  EXPECT_THAT(kv_store_.GetRaw("mixed/key"),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("Key 'kv_store/mixed/key' stores a Protobuf "
                                 "message, not raw bytes")));
}

TEST_F(FakeKeyValueStoreTest, SameKeyCanOverwriteBetweenRawAndProtobuf) {
  ASSERT_OK_AND_ASSIGN(MutableBuffer buf1, AllocateMutableBuffer(16));
  EXPECT_OK(kv_store_.SetRaw("mixed/key", std::move(buf1)));

  StringValue proto_val;
  proto_val.set_value("now_proto");
  EXPECT_OK(kv_store_.Set("mixed/key", proto_val));
  ASSERT_OK_AND_ASSIGN(StringValue read_proto,
                       kv_store_.Get<StringValue>("mixed/key"));
  EXPECT_EQ(read_proto.value(), "now_proto");

  ASSERT_OK_AND_ASSIGN(MutableBuffer buf2, AllocateMutableBuffer(32));
  EXPECT_OK(kv_store_.SetRaw("mixed/key", std::move(buf2)));
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer read_raw, kv_store_.GetRaw("mixed/key"));
  EXPECT_EQ(read_raw.size(), 32);
}

TEST_F(FakeKeyValueStoreTest,
       GetAllAndGetAllSynchronousMatchWildcardsAndSkipRawEntries) {
  StringValue v1;
  v1.set_value("one");
  StringValue v2;
  v2.set_value("two");
  StringValue v3;
  v3.set_value("three");

  EXPECT_OK(kv_store_.Set("sensors/camera/left", v1));
  EXPECT_OK(kv_store_.Set("sensors/camera/right", v2));
  EXPECT_OK(kv_store_.Set("sensors/lidar/front", v3));

  // Store a raw value under the same parent to test isolation in GetAll.
  ASSERT_OK_AND_ASSIGN(MutableBuffer raw_buf, AllocateMutableBuffer(10));
  EXPECT_OK(kv_store_.SetRaw("sensors/camera/raw_frame", std::move(raw_buf)));

  // Test GetAllSynchronous with single-segment wildcard `*`.
  EXPECT_THAT(
      kv_store_.GetAllSynchronous("sensors/camera/*", absl::Seconds(5)),
      IsOkAndHolds(UnorderedElementsAre(Key("kv_store/sensors/camera/left"),
                                        Key("kv_store/sensors/camera/right"))));

  // Test GetAll with multi-segment wildcard `**` and callbacks.
  std::vector<std::string> keys_called;
  std::string done_key;
  ASSERT_OK_AND_ASSIGN(
      KVQuery query,
      kv_store_.GetAll(
          "sensors/**",
          [&keys_called](absl::string_view key,
                         std::unique_ptr<google::protobuf::Any> val) {
            keys_called.emplace_back(key);
          },
          [&done_key](absl::string_view key) { done_key = std::string(key); }));

  EXPECT_THAT(keys_called,
              UnorderedElementsAre("kv_store/sensors/camera/left",
                                   "kv_store/sensors/camera/right",
                                   "kv_store/sensors/lidar/front"));
  EXPECT_EQ(done_key, "kv_store/sensors/**");
}

TEST_F(FakeKeyValueStoreTest,
       GetAllRawAndGetAllRawSynchronousMatchWildcardsAndSkipProtobufEntries) {
  ASSERT_OK_AND_ASSIGN(MutableBuffer b1, AllocateMutableBuffer(10));
  ASSERT_OK_AND_ASSIGN(MutableBuffer b2, AllocateMutableBuffer(20));
  ASSERT_OK_AND_ASSIGN(MutableBuffer b3, AllocateMutableBuffer(30));

  EXPECT_OK(kv_store_.SetRaw("data/chunk_1", std::move(b1)));
  EXPECT_OK(kv_store_.SetRaw("data/chunk_2", std::move(b2)));
  EXPECT_OK(kv_store_.SetRaw("data/meta_info", std::move(b3)));

  // Store a Protobuf value to test raw query isolation.
  StringValue pv;
  pv.set_value("proto");
  EXPECT_OK(kv_store_.Set("data/chunk_proto", pv));

  // Test GetAllRawSynchronous with `$*` pattern.
  using RawMap = absl::flat_hash_map<std::string, ImmutableBuffer>;
  ASSERT_OK_AND_ASSIGN(
      RawMap sync_results,
      kv_store_.GetAllRawSynchronous("data/chunk_$*", absl::Seconds(5)));
  EXPECT_EQ(sync_results.size(), 2);
  EXPECT_TRUE(sync_results.contains("kv_store/data/chunk_1"));
  EXPECT_TRUE(sync_results.contains("kv_store/data/chunk_2"));
  EXPECT_EQ(sync_results["kv_store/data/chunk_1"].size(), 10);
  EXPECT_EQ(sync_results["kv_store/data/chunk_2"].size(), 20);

  // Test GetAllRaw streaming callback.
  std::vector<std::string> keys_called;
  std::string done_key;
  ASSERT_OK_AND_ASSIGN(
      KVQuery query,
      kv_store_.GetAllRaw(
          "data/**",
          [&keys_called](absl::string_view key, ImmutableBuffer buf) {
            keys_called.emplace_back(key);
          },
          [&done_key](absl::string_view key) { done_key = std::string(key); }));

  EXPECT_THAT(keys_called, UnorderedElementsAre("kv_store/data/chunk_1",
                                                "kv_store/data/chunk_2",
                                                "kv_store/data/meta_info"));
  EXPECT_EQ(done_key, "kv_store/data/**");
}

TEST_F(FakeKeyValueStoreTest,
       WildcardQueriesHandleOverlappingSubpatternsAndDepthBounds) {
  auto put_raw = [this](absl::string_view key) {
    ASSERT_OK_AND_ASSIGN(MutableBuffer buf, AllocateMutableBuffer(4));
    EXPECT_OK(kv_store_.SetRaw(key, std::move(buf)));
  };

  put_raw("x/aba");
  put_raw("x/abba");
  put_raw("x/ab123ba");
  put_raw("x/abc");
  put_raw("x/axbyc");
  put_raw("x/axcyb");
  put_raw("x/abbbc");
  put_raw("x/axbbybc");
  put_raw("x/abbc");
  put_raw("x/abcd");
  put_raw("x/a1b2c3d");
  put_raw("x/a1c2b3d");
  put_raw("x/a1bccde");
  put_raw("x/a1bcde");
  put_raw("a");
  put_raw("a/only");
  put_raw("a/b");
  put_raw("a/m/n/b");
  put_raw("a/two/deep");

  // Overlapping prefix "ab" and suffix "ba" must not match "x/aba"; also
  // exercises the 1-arg default-timeout GetAllRawSynchronous forwarding
  // overload.
  EXPECT_THAT(kv_store_.GetAllRawSynchronous("x/ab$*ba"),
              IsOkAndHolds(UnorderedElementsAre(Key("kv_store/x/abba"),
                                                Key("kv_store/x/ab123ba"))));

  // Multi-$* ordered sub-patterns match empty or populated gaps in order only.
  EXPECT_THAT(kv_store_.GetAllRawSynchronous("x/a$*b$*c"),
              IsOkAndHolds(UnorderedElementsAre(
                  Key("kv_store/x/abc"), Key("kv_store/x/axbyc"),
                  Key("kv_store/x/abbbc"), Key("kv_store/x/axbbybc"),
                  Key("kv_store/x/abbc"))));

  // Intermediate sub-pattern with length > 1 must advance past the full match.
  EXPECT_THAT(kv_store_.GetAllRawSynchronous("x/a$*bb$*bc"),
              IsOkAndHolds(UnorderedElementsAre(Key("kv_store/x/abbbc"),
                                                Key("kv_store/x/axbbybc"))));
  EXPECT_THAT(kv_store_.GetAllRawSynchronous("x/a$*bc$*cd$*e"),
              IsOkAndHolds(UnorderedElementsAre(Key("kv_store/x/a1bccde"))));

  // Multiple intermediate sub-patterns must all be matched in sequence.
  EXPECT_THAT(kv_store_.GetAllRawSynchronous("x/a$*b$*c$*d"),
              IsOkAndHolds(UnorderedElementsAre(Key("kv_store/x/abcd"),
                                                Key("kv_store/x/a1b2c3d"))));

  // `**` matches 0 or more intermediate or trailing segments.
  EXPECT_THAT(kv_store_.GetAllRawSynchronous("a/**/b"),
              IsOkAndHolds(UnorderedElementsAre(Key("kv_store/a/b"),
                                                Key("kv_store/a/m/n/b"))));

  EXPECT_THAT(
      kv_store_.GetAllRawSynchronous("a/**"),
      IsOkAndHolds(UnorderedElementsAre(
          Key("kv_store/a"), Key("kv_store/a/only"), Key("kv_store/a/b"),
          Key("kv_store/a/m/n/b"), Key("kv_store/a/two/deep"))));

  // `*` matches strictly 1 segment.
  EXPECT_THAT(kv_store_.GetAllRawSynchronous("a/*"),
              IsOkAndHolds(UnorderedElementsAre(Key("kv_store/a/only"),
                                                Key("kv_store/a/b"))));

  // 2-arg GetAllRaw and empty keyexpr validation.
  EXPECT_THAT(
      kv_store_.GetAllRaw("", [](absl::string_view, ImmutableBuffer) {}),
      StatusIs(absl::StatusCode::kInvalidArgument));
  std::vector<std::string> two_arg_keys;
  ASSERT_OK_AND_ASSIGN(
      KVQuery two_arg_query,
      kv_store_.GetAllRaw("x/a$*bc$*cd$*e", [&two_arg_keys](absl::string_view k,
                                                            ImmutableBuffer) {
        two_arg_keys.emplace_back(k);
      }));
  EXPECT_THAT(two_arg_keys, UnorderedElementsAre("kv_store/x/a1bccde"));
}

TEST_F(FakeKeyValueStoreTest, RawSubscriptionReceivesUpdatesAndDeletion) {
  std::vector<std::pair<std::string, size_t>> updates;
  std::vector<std::string> deletions;
  ASSERT_OK_AND_ASSIGN(
      Subscription sub,
      kv_store_.CreateRawSubscription(
          "device/**", TopicConfig(),
          [&updates](absl::string_view key, const ImmutableBuffer& buf) {
            updates.emplace_back(std::string(key), buf.size());
          },
          [&deletions](absl::string_view key) {
            deletions.emplace_back(std::string(key));
          }));
  EXPECT_EQ(sub.TopicName(), "kv_store/device/**");

  // 1. SetRaw and 0-byte SetRaw trigger update callback, not deletion.
  ASSERT_OK_AND_ASSIGN(MutableBuffer status_buf, AllocateMutableBuffer(42));
  EXPECT_OK(kv_store_.SetRaw("device/status", std::move(status_buf)));
  ASSERT_OK_AND_ASSIGN(MutableBuffer heartbeat_buf, AllocateMutableBuffer(0));
  EXPECT_OK(kv_store_.SetRaw("device/heartbeat", std::move(heartbeat_buf)));

  // 2. Protobuf writes do not trigger raw updates, though Delete does.
  StringValue pv;
  pv.set_value("ignored");
  EXPECT_OK(kv_store_.Set("device/status", pv));
  EXPECT_OK(kv_store_.Set("device/proto_only", pv));
  EXPECT_OK(kv_store_.Delete("device/status"));
  EXPECT_OK(kv_store_.Delete("device/proto_only"));
  EXPECT_THAT(updates, ElementsAre(Pair("kv_store/device/status", 42),
                                   Pair("kv_store/device/heartbeat", 0)));
  EXPECT_THAT(deletions, ElementsAre("kv_store/device/status",
                                     "kv_store/device/proto_only"));

  // 3. Unsubscribe stops subsequent deliveries.
  sub.Unsubscribe();
  EXPECT_TRUE(sub.TopicName().empty());
  ASSERT_OK_AND_ASSIGN(MutableBuffer after_unsub_buf,
                       AllocateMutableBuffer(10));
  EXPECT_OK(kv_store_.SetRaw("device/status", std::move(after_unsub_buf)));
  EXPECT_EQ(updates.size(), 2);
}

TEST_F(FakeKeyValueStoreTest,
       ProtoSubscriptionIgnoresRawWritesAndDeliversDeletions) {
  std::vector<std::pair<std::string, std::string>> updates;
  std::vector<std::string> deletions;
  ASSERT_OK_AND_ASSIGN(
      Subscription sub,
      kv_store_.CreateSubscription(
          "proto/topic", TopicConfig(),
          [&updates](absl::string_view key, const google::protobuf::Any& any) {
            StringValue unpacked;
            EXPECT_TRUE(any.UnpackTo(&unpacked));
            updates.emplace_back(std::string(key), unpacked.value());
          },
          [&deletions](absl::string_view key) {
            deletions.emplace_back(std::string(key));
          }));
  ASSERT_OK_AND_ASSIGN(MutableBuffer raw_buf, AllocateMutableBuffer(16));
  EXPECT_OK(kv_store_.SetRaw("proto/topic", std::move(raw_buf)));
  EXPECT_TRUE(updates.empty());
  EXPECT_OK(kv_store_.Delete("proto/topic"));
  EXPECT_THAT(deletions, ElementsAre("kv_store/proto/topic"));
  StringValue val;
  val.set_value("hello");
  EXPECT_OK(kv_store_.Set("proto/topic", val));
  EXPECT_THAT(updates, ElementsAre(Pair("kv_store/proto/topic", "hello")));
}

TEST_F(FakeKeyValueStoreTest, SubscriptionOutlivingFakeKeyValueStoreIsSafe) {
  Subscription sub;
  {
    FakeKeyValueStore short_lived_kv;
    ASSERT_OK_AND_ASSIGN(sub,
                         short_lived_kv.CreateRawSubscription(
                             "short/**", TopicConfig(),
                             [](absl::string_view, const ImmutableBuffer&) {},
                             [](absl::string_view) {}));
    EXPECT_EQ(sub.TopicName(), "kv_store/short/**");
  }
  sub.Unsubscribe();
  EXPECT_TRUE(sub.TopicName().empty());
}

TEST_F(FakeKeyValueStoreTest, SubscriptionCallbacksSupportReentrancy) {
  auto noop = [](absl::string_view, const ImmutableBuffer&) {};
  ASSERT_OK_AND_ASSIGN(MutableBuffer transient_buf, AllocateMutableBuffer(4));
  EXPECT_OK(kv_store_.SetRaw("reentrant/transient", std::move(transient_buf)));
  int update_calls = 0;
  int nested_calls = 0;
  std::optional<Subscription> nested_sub;
  ASSERT_OK_AND_ASSIGN(
      Subscription update_sub,
      kv_store_.CreateRawSubscription(
          "reentrant/src", TopicConfig(),
          [this, &update_calls, &nested_calls, &nested_sub](
              absl::string_view key, const ImmutableBuffer& buf) {
            ++update_calls;
            EXPECT_EQ(key, "kv_store/reentrant/src");
            if (!nested_sub.has_value()) {
              ASSERT_OK_AND_ASSIGN(nested_sub,
                                   kv_store_.CreateRawSubscription(
                                       "reentrant/derived", TopicConfig(),
                                       [&nested_calls](absl::string_view,
                                                       const ImmutableBuffer&) {
                                         ++nested_calls;
                                       }));
              EXPECT_OK(kv_store_.Delete("reentrant/transient"));
            }
            EXPECT_OK(kv_store_.SetRaw("reentrant/derived", buf));
          },
          /*deletion_callback=*/nullptr));
  ASSERT_OK_AND_ASSIGN(MutableBuffer src_buf, MakeMutableBuffer("reentr"));
  EXPECT_OK(kv_store_.SetRaw("reentrant/src", std::move(src_buf)));
  EXPECT_EQ(update_calls, 1);
  EXPECT_EQ(nested_calls, 1);
  EXPECT_THAT(kv_store_.GetRaw("reentrant/transient"),
              StatusIs(absl::StatusCode::kNotFound));
  update_sub.Unsubscribe();
  EXPECT_TRUE(update_sub.TopicName().empty());
  nested_sub->Unsubscribe();
  ASSERT_OK_AND_ASSIGN(MutableBuffer after_unsub_buf, AllocateMutableBuffer(3));
  EXPECT_OK(kv_store_.SetRaw("reentrant/src", std::move(after_unsub_buf)));
  EXPECT_EQ(update_calls, 1);
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer derived_read,
                       kv_store_.GetRaw("reentrant/derived"));
  EXPECT_EQ(AsStringView(derived_read), "reentr");

  int delete_calls = 0;
  ASSERT_OK_AND_ASSIGN(
      Subscription delete_sub,
      kv_store_.CreateRawSubscription(
          "reentrant/del", TopicConfig(), noop,
          [this, &delete_calls](absl::string_view key) {
            ++delete_calls;
            EXPECT_EQ(key, "kv_store/reentrant/del");
            ASSERT_OK_AND_ASSIGN(MutableBuffer ok_buf, MakeMutableBuffer("ok"));
            EXPECT_OK(
                kv_store_.SetRaw("reentrant/after_del", std::move(ok_buf)));
          }));
  ASSERT_OK_AND_ASSIGN(MutableBuffer del_buf, AllocateMutableBuffer(4));
  EXPECT_OK(kv_store_.SetRaw("reentrant/del", std::move(del_buf)));
  EXPECT_OK(kv_store_.Delete("reentrant/del"));
  EXPECT_EQ(delete_calls, 1);
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer after_del_read,
                       kv_store_.GetRaw("reentrant/after_del"));
  EXPECT_EQ(AsStringView(after_del_read), "ok");
  delete_sub.Unsubscribe();
  EXPECT_TRUE(delete_sub.TopicName().empty());
}

TEST_F(FakeKeyValueStoreTest,
       CreateRawSubscriptionRejectsEmptyKeyAndSupportsDeleteOnly) {
  auto noop = [](absl::string_view, const ImmutableBuffer&) {};
  EXPECT_THAT(kv_store_.CreateRawSubscription("", TopicConfig(), noop),
              StatusIs(absl::StatusCode::kInvalidArgument));

  std::vector<std::string> del_only;
  ASSERT_OK_AND_ASSIGN(
      Subscription del_sub,
      kv_store_.CreateRawSubscription(
          "reentrant/d0", TopicConfig(), /*value_callback=*/nullptr,
          [&del_only](absl::string_view k) { del_only.emplace_back(k); }));
  // Also register a matching subscription with no deletion callback to verify
  // `Delete()` skips null deletion callbacks safely.
  ASSERT_OK_AND_ASSIGN(
      Subscription value_only_sub,
      kv_store_.CreateRawSubscription("reentrant/d0", TopicConfig(), noop));
  ASSERT_OK_AND_ASSIGN(MutableBuffer del_buf, MakeMutableBuffer("x"));
  EXPECT_OK(kv_store_.SetRaw("reentrant/d0", std::move(del_buf)));
  EXPECT_OK(kv_store_.Delete("reentrant/d0"));
  EXPECT_THAT(del_only, ElementsAre("kv_store/reentrant/d0"));
  value_only_sub.Unsubscribe();
}

TEST_F(FakeKeyValueStoreTest,
       UnsubscribeDestroyingCapturedSiblingSubscriptionIsSafe) {
  auto noop = [](absl::string_view, const ImmutableBuffer&) {};
  ASSERT_OK_AND_ASSIGN(
      Subscription sibling_sub,
      kv_store_.CreateRawSubscription("reentrant/s", TopicConfig(), noop));
  auto sib = std::make_shared<Subscription>(std::move(sibling_sub));
  ASSERT_OK_AND_ASSIGN(
      Subscription outer_sub,
      kv_store_.CreateRawSubscription(
          "reentrant/o", TopicConfig(),
          [s = std::move(sib)](absl::string_view, const ImmutableBuffer&) {}));
  outer_sub.Unsubscribe();
  EXPECT_TRUE(outer_sub.TopicName().empty());
}

TEST_F(FakeKeyValueStoreTest, UnsubscribeDestroyingSelfOwnerIsSafe) {
  class Owner {
   public:
    explicit Owner(bool* destroyed) : destroyed_(destroyed) {}
    Owner(const Owner&) = delete;
    Owner& operator=(const Owner&) = delete;
    ~Owner() { *destroyed_ = true; }

    Subscription& sub() { return sub_; }

   private:
    Subscription sub_;
    bool* destroyed_;
  };
  bool owner_destroyed = false;
  auto owner = std::make_shared<Owner>(&owner_destroyed);
  ASSERT_OK_AND_ASSIGN(
      owner->sub(), kv_store_.CreateRawSubscription(
                        "reentrant/self_owner", TopicConfig(),
                        [owner](absl::string_view, const ImmutableBuffer&) {}));
  Owner* raw_owner = owner.get();
  owner = nullptr;
  ASSERT_FALSE(owner_destroyed);
  raw_owner->sub().Unsubscribe();
  EXPECT_TRUE(owner_destroyed);
}

TEST_F(FakeKeyValueStoreTest,
       UnsubscribeCoMatchingSiblingInsideCallbackSkipsAndDestroysCaptures) {
  class DestroyFlag {
   public:
    explicit DestroyFlag(bool* destroyed) : destroyed_(destroyed) {}
    DestroyFlag(const DestroyFlag&) = delete;
    DestroyFlag& operator=(const DestroyFlag&) = delete;
    ~DestroyFlag() { *destroyed_ = true; }

   private:
    bool* destroyed_;
  };
  bool a_destroyed = false;
  bool b_destroyed = false;
  int total_calls = 0;
  Subscription sub_a;
  Subscription sub_b;
  auto sentinel_a = std::make_shared<DestroyFlag>(&a_destroyed);
  auto sentinel_b = std::make_shared<DestroyFlag>(&b_destroyed);

  ASSERT_OK_AND_ASSIGN(
      sub_a, kv_store_.CreateRawSubscription(
                 "reentrant/co_match", TopicConfig(),
                 [sentinel = std::move(sentinel_a), &sub_b, &b_destroyed,
                  &total_calls](absl::string_view, const ImmutableBuffer&) {
                   ++total_calls;
                   sub_b.Unsubscribe();
                   EXPECT_TRUE(b_destroyed);
                 }));
  ASSERT_OK_AND_ASSIGN(
      sub_b, kv_store_.CreateRawSubscription(
                 "reentrant/co_match", TopicConfig(),
                 [sentinel = std::move(sentinel_b), &sub_a, &a_destroyed,
                  &total_calls](absl::string_view, const ImmutableBuffer&) {
                   ++total_calls;
                   sub_a.Unsubscribe();
                   EXPECT_TRUE(a_destroyed);
                 }));

  ASSERT_OK_AND_ASSIGN(MutableBuffer buf, MakeMutableBuffer("ping"));
  EXPECT_OK(kv_store_.SetRaw("reentrant/co_match", std::move(buf)));
  EXPECT_EQ(total_calls, 1);
  EXPECT_TRUE(a_destroyed || b_destroyed);
}

TEST_F(FakeKeyValueStoreTest,
       SubscriptionCallbacksInvokeInPlaceWithoutCopying) {
  class CopyTracker {
   public:
    CopyTracker(int* copies, int* destructions)
        : copies_(copies), destructions_(destructions) {}
    CopyTracker(const CopyTracker& other)
        : copies_(other.copies_), destructions_(other.destructions_) {
      ++(*copies_);
    }
    CopyTracker(CopyTracker&& other) noexcept
        : copies_(other.copies_),
          destructions_(other.destructions_),
          is_live_(std::exchange(other.is_live_, false)) {}
    CopyTracker& operator=(const CopyTracker&) = delete;
    CopyTracker& operator=(CopyTracker&&) = delete;
    ~CopyTracker() {
      if (is_live_) {
        ++(*destructions_);
      }
    }

   private:
    int* copies_;
    int* destructions_;
    bool is_live_ = true;
  };

  int raw_copies = 0;
  int raw_destructions = 0;
  int del_copies = 0;
  int del_destructions = 0;
  std::vector<int> raw_seqs;
  std::vector<int> del_seqs;
  ASSERT_OK_AND_ASSIGN(
      Subscription raw_sub,
      kv_store_.CreateRawSubscription(
          "stateful/raw", TopicConfig(),
          [this, tracker = CopyTracker(&raw_copies, &raw_destructions), seq = 0,
           &raw_seqs](absl::string_view, const ImmutableBuffer& buf) mutable {
            static_cast<void>(tracker);
            ++seq;
            raw_seqs.push_back(seq);
            if (seq == 1) {
              EXPECT_OK(kv_store_.SetRaw("stateful/raw", buf));
            }
          },
          [tracker = CopyTracker(&del_copies, &del_destructions), del_seq = 0,
           &del_seqs](absl::string_view) mutable {
            static_cast<void>(tracker);
            del_seqs.push_back(++del_seq);
          }));

  // Reset counters after registration so we strictly assert zero copies or
  // destructions on notification dispatch.
  raw_copies = 0;
  raw_destructions = 0;
  del_copies = 0;
  del_destructions = 0;

  ASSERT_OK_AND_ASSIGN(MutableBuffer buf1, MakeMutableBuffer("v1"));
  EXPECT_OK(kv_store_.SetRaw("stateful/raw", std::move(buf1)));
  ASSERT_OK_AND_ASSIGN(MutableBuffer buf2, MakeMutableBuffer("v2"));
  EXPECT_OK(kv_store_.SetRaw("stateful/raw", std::move(buf2)));
  EXPECT_OK(kv_store_.Delete("stateful/raw"));
  ASSERT_OK_AND_ASSIGN(MutableBuffer buf3, MakeMutableBuffer("v3"));
  EXPECT_OK(kv_store_.SetRaw("stateful/raw", std::move(buf3)));
  EXPECT_OK(kv_store_.Delete("stateful/raw"));

  EXPECT_THAT(raw_seqs, ElementsAre(1, 2, 3, 4));
  EXPECT_THAT(del_seqs, ElementsAre(1, 2));
  EXPECT_EQ(raw_copies, 0);
  EXPECT_EQ(raw_destructions, 0);
  EXPECT_EQ(del_copies, 0);
  EXPECT_EQ(del_destructions, 0);

  raw_sub.Unsubscribe();
  EXPECT_EQ(raw_copies, 0);
  EXPECT_EQ(raw_destructions, 1);
  EXPECT_EQ(del_copies, 0);
  EXPECT_EQ(del_destructions, 1);

  // Verify pointer stability across `flat_hash_map` rehashes using a small
  // trivially-copyable lambda stored inline in `std::function`'s SBO buffer.
  struct SboRehashContext {
    FakeKeyValueStore* store;
    std::vector<int> seqs;
    std::vector<Subscription> extra_subs;
  };
  SboRehashContext sbo_ctx{.store = &kv_store_};
  auto sbo_cb = [&sbo_ctx, seq = 0](absl::string_view,
                                    const ImmutableBuffer&) mutable {
    if (seq == 0) {
      for (int i = 0; i < 64; ++i) {
        ASSERT_OK_AND_ASSIGN(
            Subscription extra,
            sbo_ctx.store->CreateRawSubscription(
                "stateful/rehash_other", TopicConfig(),
                [](absl::string_view, const ImmutableBuffer&) {}));
        sbo_ctx.extra_subs.push_back(std::move(extra));
      }
    }
    sbo_ctx.seqs.push_back(++seq);
  };
  static_assert(std::is_trivially_copyable_v<decltype(sbo_cb)> &&
                sizeof(sbo_cb) <= 16);
  ASSERT_OK_AND_ASSIGN(Subscription sbo_sub,
                       kv_store_.CreateRawSubscription("stateful/sbo_rehash",
                                                       TopicConfig(), sbo_cb));
  ASSERT_OK_AND_ASSIGN(MutableBuffer sbo_buf1, MakeMutableBuffer("r1"));
  EXPECT_OK(kv_store_.SetRaw("stateful/sbo_rehash", std::move(sbo_buf1)));
  ASSERT_OK_AND_ASSIGN(MutableBuffer sbo_buf2, MakeMutableBuffer("r2"));
  EXPECT_OK(kv_store_.SetRaw("stateful/sbo_rehash", std::move(sbo_buf2)));
  EXPECT_THAT(sbo_ctx.seqs, ElementsAre(1, 2));
  sbo_sub.Unsubscribe();

  std::vector<int> proto_seqs;
  ASSERT_OK_AND_ASSIGN(
      Subscription proto_sub,
      kv_store_.CreateSubscription(
          "stateful/proto", TopicConfig(),
          [seq = 0, &proto_seqs](absl::string_view,
                                 const google::protobuf::Any&) mutable {
            proto_seqs.push_back(++seq);
          },
          /*deletion_callback=*/nullptr));
  StringValue val;
  val.set_value("p");
  EXPECT_OK(kv_store_.Set("stateful/proto", val));
  EXPECT_OK(kv_store_.Set("stateful/proto", val));
  EXPECT_THAT(proto_seqs, ElementsAre(1, 2));
  proto_sub.Unsubscribe();
}

TEST_F(FakeKeyValueStoreTest,
       ReentrantSameSubscriptionCompletesNestedInvocations) {
  int calls = 0;
  ASSERT_OK_AND_ASSIGN(
      Subscription sub,
      kv_store_.CreateRawSubscription(
          "reentrant/same_key", TopicConfig(),
          [this, &calls](absl::string_view, const ImmutableBuffer& buf) {
            ++calls;
            if (calls == 1) {
              EXPECT_OK(kv_store_.SetRaw("reentrant/same_key", buf));
            }
          }));

  ASSERT_OK_AND_ASSIGN(MutableBuffer buf, MakeMutableBuffer("nested"));
  EXPECT_OK(kv_store_.SetRaw("reentrant/same_key", std::move(buf)));
  EXPECT_EQ(calls, 2);
  sub.Unsubscribe();
  EXPECT_TRUE(sub.TopicName().empty());
}

TEST_F(FakeKeyValueStoreTest,
       UnsubscribeWaitsForInFlightCallbackOnOtherThread) {
  absl::Notification entered;
  absl::Notification proceed;
  absl::Notification unsub_done;
  class DestroySentinel {
   public:
    explicit DestroySentinel(std::atomic<bool>* destroyed)
        : destroyed_(destroyed) {}
    DestroySentinel(const DestroySentinel&) = delete;
    DestroySentinel& operator=(const DestroySentinel&) = delete;
    ~DestroySentinel() { destroyed_->store(true); }

    int value() const { return value_; }

   private:
    int value_ = 42;
    std::atomic<bool>* destroyed_;
  };
  std::atomic<bool> capture_destroyed{false};
  auto sentinel = std::make_shared<DestroySentinel>(&capture_destroyed);
  ASSERT_OK_AND_ASSIGN(Subscription drain,
                       kv_store_.CreateRawSubscription(
                           "reentrant/drain", TopicConfig(),
                           [&entered, &proceed, &capture_destroyed,
                            sentinel = std::move(sentinel)](
                               absl::string_view, const ImmutableBuffer&) {
                             entered.Notify();
                             proceed.WaitForNotification();
                             EXPECT_FALSE(capture_destroyed.load());
                             EXPECT_EQ(sentinel->value(), 42);
                           }));
  ASSERT_OK_AND_ASSIGN(MutableBuffer drain_buf, MakeMutableBuffer("d"));
  std::thread writer([this, buf = std::move(drain_buf)]() mutable {
    EXPECT_OK(kv_store_.SetRaw("reentrant/drain", std::move(buf)));
  });
  entered.WaitForNotification();
  std::thread unsubscriber([&drain, &unsub_done] {
    drain.Unsubscribe();
    unsub_done.Notify();
  });
  EXPECT_FALSE(
      unsub_done.WaitForNotificationWithTimeout(absl::Milliseconds(20)));
  EXPECT_FALSE(capture_destroyed.load());
  proceed.Notify();
  writer.join();
  unsubscriber.join();
  EXPECT_TRUE(unsub_done.HasBeenNotified());
  EXPECT_TRUE(capture_destroyed.load());
}

#ifdef GTEST_HAS_DEATH_TEST
TEST(FakeKeyValueStoreDeathTest, SelfUnsubscribeInsideCallbackAborts) {
  EXPECT_DEATH(
      {
        FakeKeyValueStore kv;
        Subscription s;
        ASSERT_OK_AND_ASSIGN(
            s, kv.CreateRawSubscription(
                   "death/put", TopicConfig(),
                   [&s](absl::string_view, const ImmutableBuffer&) {
                     s.Unsubscribe();
                   }));
        ASSERT_OK_AND_ASSIGN(MutableBuffer buf, AllocateMutableBuffer(4));
        kv.SetRaw("death/put", std::move(buf)).IgnoreError();
      },
      "Cannot Unsubscribe");
  EXPECT_DEATH(
      {
        FakeKeyValueStore kv;
        Subscription s;
        ASSERT_OK_AND_ASSIGN(
            s, kv.CreateRawSubscription(
                   "death/del", TopicConfig(),
                   [](absl::string_view, const ImmutableBuffer&) {},
                   [&s](absl::string_view) { s.Unsubscribe(); }));
        ASSERT_OK_AND_ASSIGN(MutableBuffer buf, AllocateMutableBuffer(4));
        kv.SetRaw("death/del", std::move(buf)).IgnoreError();
        kv.Delete("death/del").IgnoreError();
      },
      "Cannot Unsubscribe");
  EXPECT_DEATH(
      {
        FakeKeyValueStore kv;
        Subscription s;
        int depth = 0;
        ASSERT_OK_AND_ASSIGN(
            s, kv.CreateRawSubscription(
                   "death/nested", TopicConfig(),
                   [&kv, &s, &depth](absl::string_view,
                                     const ImmutableBuffer& buf) {
                     if (++depth == 1) {
                       kv.SetRaw("death/nested", buf).IgnoreError();
                       s.Unsubscribe();
                     }
                   }));
        ASSERT_OK_AND_ASSIGN(MutableBuffer buf, AllocateMutableBuffer(4));
        kv.SetRaw("death/nested", std::move(buf)).IgnoreError();
      },
      "Cannot Unsubscribe");
}
#endif

TEST_F(FakeKeyValueStoreTest,
       ListAllKeysIncludesBothEncodingsAndValidatesReplicationMode) {
  StringValue proto_val;
  proto_val.set_value("v");
  EXPECT_OK(kv_store_.Set("list/proto", proto_val));
  ASSERT_OK_AND_ASSIGN(MutableBuffer raw_buf, AllocateMutableBuffer(4));
  EXPECT_OK(kv_store_.SetRaw("list/raw", std::move(raw_buf)));
  EXPECT_OK(kv_store_.Set("list/deleted", proto_val));
  EXPECT_OK(kv_store_.Delete("list/deleted"));

  EXPECT_THAT(kv_store_.ListAllKeys(),
              IsOkAndHolds(UnorderedElementsAre("kv_store/list/proto",
                                                "kv_store/list/raw")));
  EXPECT_THAT(kv_store_.ListAllKeys(absl::ZeroDuration()),
              IsOkAndHolds(UnorderedElementsAre("kv_store/list/proto",
                                                "kv_store/list/raw")));
  EXPECT_THAT(kv_store_.ListAllKeys(-absl::Milliseconds(1)),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(kv_store_.ListAllGlobalKeys(),
              StatusIs(absl::StatusCode::kUnimplemented));
  EXPECT_THAT(kv_store_.ListAllOnpremKeys("wc1"),
              StatusIs(absl::StatusCode::kUnimplemented));

  absl::FlagSaver flag_saver;
  absl::SetFlag(&FLAGS_use_replicated_kv_store, true);
  FakeKeyValueStore repl_kv;
  EXPECT_OK(repl_kv.Set("global/cfg", proto_val));
  ASSERT_OK_AND_ASSIGN(MutableBuffer repl_raw, AllocateMutableBuffer(4));
  EXPECT_OK(repl_kv.SetRaw("global/blob", std::move(repl_raw)));
  EXPECT_OK(repl_kv.Set("global/old", proto_val));
  EXPECT_OK(repl_kv.Delete("global/old"));
  EXPECT_OK(repl_kv.Set("wc1/state", proto_val));
  ASSERT_OK_AND_ASSIGN(MutableBuffer wc2_raw, AllocateMutableBuffer(4));
  EXPECT_OK(repl_kv.SetRaw("wc2/frame", std::move(wc2_raw)));

  EXPECT_THAT(repl_kv.ListAllGlobalKeys(),
              IsOkAndHolds(UnorderedElementsAre("kv_store_repl/global/cfg",
                                                "kv_store_repl/global/blob")));
  EXPECT_THAT(repl_kv.ListAllOnpremKeys("wc1"),
              IsOkAndHolds(UnorderedElementsAre("kv_store_repl/wc1/state")));
  EXPECT_THAT(repl_kv.ListAllOnpremKeys("wc2"),
              IsOkAndHolds(UnorderedElementsAre("kv_store_repl/wc2/frame")));
  EXPECT_THAT(repl_kv.ListAllOnpremKeys(""),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(repl_kv.ListAllGlobalKeys(-absl::Milliseconds(1)),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(repl_kv.ListAllOnpremKeys("wc1", -absl::Milliseconds(1)),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(repl_kv.ListAllKeys(),
              StatusIs(absl::StatusCode::kUnimplemented));
  EXPECT_THAT(
      repl_kv.GetAllRawSynchronous("global/**"),
      IsOkAndHolds(UnorderedElementsAre(Key("kv_store_repl/global/blob"))));
  std::vector<std::string> repl_raw_keys;
  std::string repl_done_key;
  ASSERT_OK_AND_ASSIGN(
      KVQuery repl_raw_query,
      repl_kv.GetAllRaw(
          "wc2/**",
          [&repl_raw_keys](absl::string_view k, ImmutableBuffer) {
            repl_raw_keys.emplace_back(k);
          },
          [&repl_done_key](absl::string_view k) { repl_done_key = k; }));
  EXPECT_THAT(repl_raw_keys, UnorderedElementsAre("kv_store_repl/wc2/frame"));
  EXPECT_EQ(repl_done_key, "kv_store_repl/wc2/**");
}

TEST_F(FakeKeyValueStoreTest, ReplicatedFlagUsesReplicationKeyPrefix) {
  absl::FlagSaver flag_saver;
  FakeKeyValueStore custom_kv("custom_kv");
  std::string ckey;
  ASSERT_OK_AND_ASSIGN(
      Subscription custom_sub,
      custom_kv.CreateRawSubscription(
          "item", TopicConfig(),
          [&ckey](absl::string_view k, const ImmutableBuffer&) { ckey = k; }));
  ASSERT_OK_AND_ASSIGN(MutableBuffer cbuf, MakeMutableBuffer("v"));
  EXPECT_OK(custom_kv.SetRaw("item", std::move(cbuf)));
  EXPECT_EQ(ckey, "custom_kv/item");

  absl::SetFlag(&FLAGS_use_replicated_kv_store, true);
  FakeKeyValueStore repl_kv;
  std::vector<std::string> update_keys, del_keys;
  ASSERT_OK_AND_ASSIGN(
      Subscription sub,
      repl_kv.CreateRawSubscription(
          "repl/**", TopicConfig(),
          [&update_keys](absl::string_view k, const ImmutableBuffer&) {
            update_keys.emplace_back(k);
          },
          [&del_keys](absl::string_view k) { del_keys.emplace_back(k); }));
  EXPECT_EQ(sub.TopicName(), "kv_store_repl/repl/**");
  ASSERT_OK_AND_ASSIGN(MutableBuffer buf, MakeMutableBuffer("value"));
  EXPECT_OK(repl_kv.SetRaw("repl/item", std::move(buf)));
  EXPECT_THAT(update_keys, ElementsAre("kv_store_repl/repl/item"));
  EXPECT_OK(repl_kv.Delete("repl/item"));
  EXPECT_THAT(del_keys, ElementsAre("kv_store_repl/repl/item"));
}

}  // namespace
}  // namespace intrinsic
