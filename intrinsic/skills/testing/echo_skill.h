// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SKILLS_TESTING_ECHO_SKILL_H_
#define INTRINSIC_SKILLS_TESTING_ECHO_SKILL_H_

#include <memory>

#include "absl/status/statusor.h"
#include "google/protobuf/message.h"
#include "intrinsic/skills/cc/skill_interface.h"
#include "intrinsic/skills/proto/skill_service.pb.h"

namespace intrinsic::skills {

class EchoSkill : public SkillInterface {
 public:
  static std::unique_ptr<SkillInterface> CreateSkill();

  absl::StatusOr<std::unique_ptr<google::protobuf::Message>> Execute(
      const ExecuteRequest& request, ExecuteContext& context) override;

  absl::StatusOr<std::unique_ptr<::google::protobuf::Message>> Preview(
      const PreviewRequest& request, PreviewContext& context) override;
};

}  // namespace intrinsic::skills

#endif  // INTRINSIC_SKILLS_TESTING_ECHO_SKILL_H_