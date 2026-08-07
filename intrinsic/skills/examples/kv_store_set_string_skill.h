// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SKILLS_EXAMPLES_KV_STORE_SET_STRING_SKILL_H_
#define INTRINSIC_SKILLS_EXAMPLES_KV_STORE_SET_STRING_SKILL_H_

#include <memory>

#include "absl/status/statusor.h"
#include "google/protobuf/message.h"
#include "intrinsic/skills/cc/skill_interface.h"
#include "intrinsic/skills/proto/skill_service.pb.h"

namespace intrinsic {
namespace skills {

// This Skill interacts with the KVStore service of the Intrinsic Runtime to
// set a string value in the store.
class KvStoreSetStringSkill : public SkillInterface {
 public:
  absl::StatusOr<std::unique_ptr<google::protobuf::Message>> Execute(
      const ExecuteRequest& request, ExecuteContext& context) override;

  static std::unique_ptr<SkillInterface> CreateSkill() {
    return std::make_unique<KvStoreSetStringSkill>();
  }
};

}  // namespace skills
}  // namespace intrinsic

#endif  // INTRINSIC_SKILLS_EXAMPLES_KV_STORE_SET_STRING_SKILL_H_
