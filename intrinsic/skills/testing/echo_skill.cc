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

#include "intrinsic/skills/testing/echo_skill.h"

#include <memory>

#include "absl/status/statusor.h"
#include "google/protobuf/any.pb.h"
#include "google/protobuf/message.h"
#include "intrinsic/skills/cc/skill_interface.h"
#include "intrinsic/skills/cc/skill_interface_utils.h"
#include "intrinsic/skills/proto/equipment.pb.h"
#include "intrinsic/skills/proto/skill_service.pb.h"
#include "intrinsic/skills/testing/echo_skill.pb.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::skills {

std::unique_ptr<SkillInterface> EchoSkill::CreateSkill() {
  return std::make_unique<EchoSkill>();
}

absl::StatusOr<std::unique_ptr<google::protobuf::Message>> EchoSkill::Execute(
    const ExecuteRequest& request, ExecuteContext& context) {
  INTR_ASSIGN_OR_RETURN(
      auto params, request.params<intrinsic_proto::skills::EchoSkillParams>());

  auto result = std::make_unique<intrinsic_proto::skills::EchoSkillReturn>();
  result->set_foo(params.foo());

  return result;
}

absl::StatusOr<std::unique_ptr<::google::protobuf::Message>> EchoSkill::Preview(
    const PreviewRequest& request, PreviewContext& context) {
  return PreviewViaExecute(*this, request, context);
}

}  // namespace intrinsic::skills
