// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/tools/inbuild/integration_tests/inbuild_skill.h"

#include <memory>

#include "absl/log/log.h"
#include "absl/status/statusor.h"
#include "google/protobuf/descriptor.h"
#include "google/protobuf/message.h"
#include "intrinsic/skills/cc/skill_interface.h"
#include "intrinsic/skills/cc/skill_interface_utils.h"
#include "intrinsic/skills/proto/equipment.pb.h"
#include "intrinsic/skills/proto/skill_service.pb.h"
#include "intrinsic/tools/inbuild/integration_tests/inbuild_skill.pb.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::skills {

std::unique_ptr<SkillInterface> InbuildSkill::CreateSkill() {
  return std::make_unique<InbuildSkill>();
}

absl::StatusOr<std::unique_ptr<google::protobuf::Message>>
InbuildSkill::Execute(const ExecuteRequest& request, ExecuteContext& context) {
  INTR_ASSIGN_OR_RETURN(
      auto params,
      request.params<intrinsic_proto::skills::InbuildSkillParams>());

  LOG(INFO) << "Hello from InbuildSkill::Execute: " << params.foo();
  return nullptr;
}

absl::StatusOr<std::unique_ptr<google::protobuf::Message>>
InbuildSkill::Preview(const PreviewRequest& request, PreviewContext& context) {
  return PreviewViaExecute(*this, request, context);
}

}  // namespace intrinsic::skills
