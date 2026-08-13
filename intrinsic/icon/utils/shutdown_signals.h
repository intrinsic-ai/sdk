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

#ifndef INTRINSIC_ICON_UTILS_SHUTDOWN_SIGNALS_H_
#define INTRINSIC_ICON_UTILS_SHUTDOWN_SIGNALS_H_

namespace intrinsic::icon {

enum class ShutdownType {
  // No shutdown was requested.
  kNotRequested = 0,
  // A signal requested shutdown, i.e. by Kubernetes.
  kSignalledRequest,
  // A user requested the shutdown, i.e. over grpc.
  kUserRequest,
};

// Initiates a shutdown. This function is signal-safe.
void ShutdownSignalHandler(int sig);

// Returns if and what kind of shutdown was requested.
ShutdownType IsShutdownRequested();

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_UTILS_SHUTDOWN_SIGNALS_H_
