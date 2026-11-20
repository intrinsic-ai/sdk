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

#include <cstddef>
#include <cstdint>
#include <cstring>
#include <string>
#include <utility>

#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "google/protobuf/any.pb.h"
#include "google/protobuf/wrappers.pb.h"
#include "intrinsic/middleware/buffers/buffer_factory.h"
#include "intrinsic/middleware/buffers/immutable_buffer.h"
#include "intrinsic/middleware/buffers/mutable_buffer.h"
#include "internal/testing.h"
#include "intrinsic/platform/pubsub/kvstore.h"

namespace intrinsic {
namespace {

using ::absl_testing::StatusIs;
using ::google::protobuf::StringValue;
using ::testing::HasSubstr;
using ::testing::Test;

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

  // Negative and zero timeout boundaries on Get and GetRaw.
  EXPECT_THAT(kv_store_.Get<StringValue>("some/key", -absl::Seconds(1)),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(kv_store_.Get<StringValue>("some/key", absl::ZeroDuration()),
              StatusIs(absl::StatusCode::kNotFound));
  EXPECT_THAT(kv_store_.GetRaw("some/key", absl::ZeroDuration()),
              StatusIs(absl::StatusCode::kNotFound));
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

}  // namespace
}  // namespace intrinsic
