// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/skills/cc/skill_utils.h"

#include <memory>
#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_format.h"
#include "intrinsic/resources/proto/resource_handle.pb.h"
#include "intrinsic/skills/proto/equipment.pb.h"
#include "intrinsic/skills/proto/skills.pb.h"
#include "intrinsic/util/grpc/channel.h"
#include "intrinsic/util/grpc/connection_params.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::skills {

absl::StatusOr<intrinsic::ConnectionParams> GetConnectionParamsFromHandle(
    const intrinsic_proto::resources::ResourceHandle& handle) {
  if (!handle.connection_info().has_grpc()) {
    return absl::InvalidArgumentError(absl::StrFormat(
        "Resource handle \"%s\" does not specify grpc connection_info",
        handle.name()));
  }
  return intrinsic::ConnectionParams{
      .address =
          std::string(handle.connection_info().grpc().address()),  // NOLINT
      .instance_name = std::string(                                // NOLINT
          handle.connection_info().grpc().server_instance()),      // NOLINT
      .header =
          std::string(handle.connection_info().grpc().header()),  // NOLINT
  };
}

absl::StatusOr<std::shared_ptr<intrinsic::Channel>> CreateChannelFromHandle(
    const intrinsic_proto::resources::ResourceHandle& handle) {
  INTR_ASSIGN_OR_RETURN(const intrinsic::ConnectionParams connection_params,
                        GetConnectionParamsFromHandle(handle));

  return intrinsic::Channel::MakeFromAddress(connection_params);
}

}  // namespace intrinsic::skills
