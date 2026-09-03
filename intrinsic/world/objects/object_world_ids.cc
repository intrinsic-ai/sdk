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

#include "intrinsic/world/objects/object_world_ids.h"

namespace intrinsic {

const ObjectWorldResourceId& RootObjectId() {
  static const ObjectWorldResourceId* kRootObjectId =
      new ObjectWorldResourceId("root");
  return *kRootObjectId;
}

const ObjectWorldResourceId& RootEntityId() {
  static const ObjectWorldResourceId* kRootEntityId =
      new ObjectWorldResourceId("eid_root");
  return *kRootEntityId;
}

const WorldObjectName& RootObjectName() {
  static const WorldObjectName* kRootObjectName = new WorldObjectName("root");
  return *kRootObjectName;
}

const FrameName& FlangeFrameName() {
  static const FrameName* kFlangeFrameName = new FrameName("flange");
  return *kFlangeFrameName;
}

const FrameName& SensorFrameName() {
  static const FrameName* kSensorFrameName = new FrameName("sensor");
  return *kSensorFrameName;
}

}  // namespace intrinsic
