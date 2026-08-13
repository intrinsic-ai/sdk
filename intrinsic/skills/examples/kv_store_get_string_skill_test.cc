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

#include "intrinsic/skills/examples/kv_store_get_string_skill.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <memory>
#include <string>

#include "absl/container/flat_hash_map.h"
#include "absl/strings/str_cat.h"
#include "google/protobuf/any.pb.h"
#include "google/protobuf/message.h"
#include "google/protobuf/wrappers.pb.h"
#include "grpcpp/server_context.h"
#include "grpcpp/support/status.h"
#include "intrinsic/assets/interface_utils.h"
#include "intrinsic/assets/proto/v1/resolved_dependency.pb.h"
#include "intrinsic/platform/pubsub/kvstore_grpc/kvstore.grpc.pb.h"
#include "intrinsic/platform/pubsub/kvstore_grpc/kvstore.pb.h"
#include "intrinsic/skills/cc/skill_interface.h"
#include "intrinsic/skills/examples/kv_store_get_string_skill.pb.h"
#include "intrinsic/skills/testing/skill_test_utils.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

namespace intrinsic {
namespace skills {
namespace {

class FakeKvStore : public intrinsic_proto::kvstore::KVStore::Service {
 public:
  grpc::Status Get(grpc::ServerContext* context,
                   const intrinsic_proto::kvstore::GetRequest* request,
                   intrinsic_proto::kvstore::GetResponse* response) override {
    auto it = store_.find(request->key());
    if (it == store_.end()) {
      return grpc::Status(grpc::StatusCode::NOT_FOUND, "Key not found");
    }
    *response->mutable_value() = it->second;
    return grpc::Status::OK;
  }

  void SetString(const std::string& key, const std::string& value) {
    google::protobuf::StringValue str_val;
    str_val.set_value(value);
    store_[key].PackFrom(str_val);
  }

 private:
  absl::flat_hash_map<std::string, google::protobuf::Any> store_;
};

TEST(KvStoreGetStringSkillTest, GetsValueFromStore) {
  auto skill_test_factory = SkillTestFactory();

  auto skill = KvStoreGetStringSkill::CreateSkill();

  FakeKvStore kv_store_service;
  kv_store_service.SetString("my_key", "my_value");

  intrinsic_proto::assets::v1::ResolvedDependency::Interface interface =
      skill_test_factory.RunService(&kv_store_service, "intrinsic_runtime");

  intrinsic_proto::skills::KvStoreGetStringParams params;
  const std::string kv_store_interface_uri =
      absl::StrCat(::intrinsic::assets::kGrpcUriPrefix,
                   intrinsic_proto::kvstore::KVStore::service_full_name());
  (*params.mutable_intrinsic_runtime()
        ->mutable_interfaces())[kv_store_interface_uri] = interface;

  params.set_key("my_key");

  auto request = skill_test_factory.MakeExecuteRequest(params);
  auto context = skill_test_factory.MakeExecuteContext({});
  ASSERT_OK_AND_ASSIGN(std::unique_ptr<google::protobuf::Message> result,
                       skill->Execute(request, *context));

  auto return_value = google::protobuf::DownCastMessage<
      intrinsic_proto::skills::KvStoreGetStringResult>(result.get());
  ASSERT_NE(return_value, nullptr);
  EXPECT_EQ(return_value->value().value(), "my_value");
}

}  // namespace
}  // namespace skills
}  // namespace intrinsic
