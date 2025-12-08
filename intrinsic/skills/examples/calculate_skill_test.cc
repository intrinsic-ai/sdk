// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/skills/examples/calculate_skill.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <memory>

#include "google/protobuf/message_lite.h"
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

  CalculatorServiceImpl calculator_service(CalculatorConfig{});
  intrinsic_proto::assets::v1::ResolvedDependency::Interface interface =
      skill_test_factory.RunService(&calculator_service, "calculator");
  params.mutable_calculator()->mutable_interfaces()->insert(
      {"grpc://intrinsic_proto.services.Calculator", interface});

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
