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

#include "intrinsic/platform/pubsub/fake_kvstore.h"

#include <pybind11/pybind11.h>
#include <pybind11/stl.h>

#include <memory>
#include <optional>
#include <string>

#include "intrinsic/platform/pubsub/kvstore.h"

namespace intrinsic::pubsub {

PYBIND11_MODULE(fake_kvstore, m) {
  pybind11::module_::import("intrinsic.platform.pubsub.python.pubsub");
  pybind11::class_<FakeKeyValueStore, KeyValueStore,
                   std::shared_ptr<FakeKeyValueStore>>(
      m, "FakeKeyValueStore", "In-memory fake KeyValueStore for unit tests.")
      .def(pybind11::init<std::optional<std::string>>(),
           pybind11::arg("prefix_override") = std::nullopt);
}

}  // namespace intrinsic::pubsub
