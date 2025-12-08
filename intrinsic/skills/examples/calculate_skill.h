// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SKILLS_EXAMPLES_CALCULATE_SKILL_H_
#define INTRINSIC_SKILLS_EXAMPLES_CALCULATE_SKILL_H_

#include <memory>

#include "absl/status/statusor.h"
#include "google/protobuf/message.h"
#include "intrinsic/skills/cc/skill_interface.h"
#include "intrinsic/skills/proto/skill_service.pb.h"

namespace intrinsic {
namespace skills {

// This Skill interacts with a Calculator service.
class CalculateSkill : public SkillInterface {
 public:
  absl::StatusOr<std::unique_ptr<google::protobuf::Message>> Execute(
      const ExecuteRequest& request, ExecuteContext& context) override;

  static std::unique_ptr<SkillInterface> CreateSkill() {
    return std::make_unique<CalculateSkill>();
  }
};

}  // namespace skills
}  // namespace intrinsic

#endif  // INTRINSIC_SKILLS_EXAMPLES_CALCULATE_SKILL_H_
