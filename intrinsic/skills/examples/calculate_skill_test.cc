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

#include "intrinsic/skills/examples/calculate_skill.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <memory>
#include <string>

#include "absl/strings/str_cat.h"
#include "google/protobuf/message_lite.h"
#include "intrinsic/assets/interface_utils.h"
#include "intrinsic/assets/proto/v1/resolved_dependency.pb.h"
#include "intrinsic/assets/services/examples/calcserver/calc_server.h"
#include "intrinsic/assets/services/examples/calcserver/calc_server.pb.h"
#include "intrinsic/skills/cc/skill_interface.h"
#include "intrinsic/skills/examples/calculate_skill.pb.h"
#include "intrinsic/skills/testing/skill_test_utils.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

namespace intrinsic {
namespace skills {
namespace {

using intrinsic_proto::services::CalculatorConfig;
using services::CalculatorServiceImpl;

TEST(CalculateSkillTest, CalculatesSum) {
  auto skill_test_factory = SkillTestFactory();

  auto skill = CalculateSkill::CreateSkill();
  intrinsic_proto::skills::CalculateParams params;
  params.set_operation(intrinsic_proto::services::CALCULATOR_OPERATION_ADD);
  params.set_x(7);
  params.set_y(3);

  CalculatorServiceImpl calculator_service(
      CalculatorConfig{},
      [](const auto&,
         absl::string_view) -> absl::StatusOr<google::protobuf::Any> {
        return absl::UnimplementedError("Not implemented");
      });
  intrinsic_proto::assets::v1::ResolvedDependency::Interface interface =
      skill_test_factory.RunService(&calculator_service, "calculator");
  const std::string calculator_interface_uri =
      absl::StrCat(::intrinsic::assets::kGrpcUriPrefix,
                   intrinsic_proto::services::Calculator::service_full_name());
  (*params.mutable_calculator()
        ->mutable_interfaces())[calculator_interface_uri] = interface;

  auto request = skill_test_factory.MakeExecuteRequest(params);
  auto context = skill_test_factory.MakeExecuteContext({});
  ASSERT_OK_AND_ASSIGN(std::unique_ptr<google::protobuf::Message> result,
                       skill->Execute(request, *context));

  auto return_value = google::protobuf::DownCastMessage<
      intrinsic_proto::skills::CalculateResult>(result.get());
  ASSERT_NE(return_value, nullptr);
  EXPECT_EQ(return_value->result(), 10);
}

}  // namespace
}  // namespace skills
}  // namespace intrinsic
