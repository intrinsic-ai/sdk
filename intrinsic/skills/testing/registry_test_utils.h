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

#ifndef INTRINSIC_SKILLS_TESTING_REGISTRY_TEST_UTILS_H_
#define INTRINSIC_SKILLS_TESTING_REGISTRY_TEST_UTILS_H_

#include "absl/status/statusor.h"
#include "intrinsic/skills/proto/skill_manifest.pb.h"
#include "intrinsic/skills/proto/skills.pb.h"

namespace intrinsic::skills {

// Gets a Skill proto that's ready for use in, for instance, the SkillRegistry
// config. Does not populate any default parameters. Specifies `kTestVersion` as
// the semver version.
absl::StatusOr<intrinsic_proto::skills::Skill> BuildTestSkillProto(
    const intrinsic_proto::skills::SkillManifest& manifest,
    const google::protobuf::FileDescriptorSet& param_type_file_descriptor_set,
    const google::protobuf::FileDescriptorSet& return_type_file_descriptor_set);

}  // namespace intrinsic::skills

#endif  // INTRINSIC_SKILLS_TESTING_REGISTRY_TEST_UTILS_H_
