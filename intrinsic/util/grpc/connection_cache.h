// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_GRPC_CONNECTION_CACHE_H_
#define INTRINSIC_UTIL_GRPC_CONNECTION_CACHE_H_

#include <memory>
#include <utility>

#include "absl/base/thread_annotations.h"
#include "absl/container/flat_hash_map.h"
#include "absl/status/statusor.h"
#include "absl/synchronization/mutex.h"
#include "intrinsic/util/grpc/channel.h"
#include "intrinsic/util/grpc/connection_params.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {

// A thread-safe cache for multiple gRPC connections allowing stubs and channels
// to be re-used, see also https://grpc.io/docs/guides/performance/.
template <typename ServiceType>
class ConnectionCache {
 public:
  struct Connection {
    std::shared_ptr<Channel> channel;
    std::unique_ptr<typename ServiceType::Stub> stub;
  };

  absl::StatusOr<std::shared_ptr<Connection>> Get(
      const ConnectionParams& connection_params) {
    {
      absl::MutexLock lock(&mutex_);
      if (const auto it = connection_by_params_.find(connection_params);
          it != connection_by_params_.end()) {
        return it->second;
      }
    }

    INTR_ASSIGN_OR_RETURN(std::shared_ptr<Channel> channel,
                          Channel::Make(connection_params));
    std::unique_ptr<typename ServiceType::Stub> stub =
        ServiceType::NewStub(channel->GetChannel());
    if (stub == nullptr) {
      return InternalErrorBuilder()
             << "Cannot connect to service: " << connection_params;
    }

    absl::MutexLock lock(&mutex_);
    return connection_by_params_
        .emplace(connection_params, std::make_shared<Connection>(
                                        std::move(channel), std::move(stub)))
        .first->second;
  }

 private:
  absl::Mutex mutex_;
  absl::flat_hash_map<ConnectionParams, std::shared_ptr<Connection>>
      connection_by_params_ ABSL_GUARDED_BY(mutex_);
};

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_GRPC_CONNECTION_CACHE_H_
