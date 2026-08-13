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

#include "intrinsic/world/objects/transform_node.h"

#include <memory>
#include <utility>

#include "intrinsic/math/pose3.h"
#include "intrinsic/world/objects/object_world_ids.h"
#include "intrinsic/world/proto/object_world_refs.pb.h"

namespace intrinsic {
namespace world {

TransformNode::TransformNode(std::shared_ptr<const Data> data)
    : data_(std::move(data)) {}

TransformNode TransformNode::AsTransformNode() { return TransformNode(data_); }

ObjectWorldResourceId TransformNode::Id() const { return data_->Id(); }

Pose3d TransformNode::ParentTThis() const { return data_->ParentTThis(); }

intrinsic_proto::world::TransformNodeReference
TransformNode::TransformNodeReference() const {
  return data_->TransformNodeReference();
}

const TransformNode::Data& TransformNode::GetData() const { return *data_; }

}  // namespace world
}  // namespace intrinsic
