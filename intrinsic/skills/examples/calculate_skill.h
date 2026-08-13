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
