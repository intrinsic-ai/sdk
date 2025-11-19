// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ASSETS_DEPENDENCIES_UTILS_H_
#define INTRINSIC_ASSETS_DEPENDENCIES_UTILS_H_

#include <memory>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "google/protobuf/any.pb.h"
#include "google/protobuf/descriptor.h"
#include "grpcpp/client_context.h"
#include "intrinsic/assets/data/proto/v1/data_assets.grpc.pb.h"
#include "intrinsic/assets/proto/v1/resolved_dependency.pb.h"

namespace intrinsic::assets::dependencies {

// Creates a gRPC channel for communicating with the provider of the specified
// interface.
//
// The context will be updated with any needed metadata for communicating with
// the provider.
absl::StatusOr<std::shared_ptr<grpc::Channel>> Connect(
    grpc::ClientContext& context,
    const intrinsic_proto::assets::v1::ResolvedDependency& dep,
    absl::string_view iface);

// Retrieves the payload for the specified data interface.
absl::StatusOr<google::protobuf::Any> GetDataPayload(
    const intrinsic_proto::assets::v1::ResolvedDependency& dep,
    absl::string_view iface,
    intrinsic_proto::data::v1::DataAssets::StubInterface* data_assets_client =
        nullptr);

// Options for HasResolvedDependency.
struct ResolvedDepsIntrospectionOptions {
  bool check_dependency_annotation;
  bool check_skill_annotations;
};

// Checks if the given proto has any ResolvedDependency fields.
//
// If additional introspection options are provided, the method returns true
// only if all of the options are satisfied.
bool HasResolvedDependency(const google::protobuf::Descriptor& descriptor,
                           const ResolvedDepsIntrospectionOptions& options);

}  // namespace intrinsic::assets::dependencies

#endif  // INTRINSIC_ASSETS_DEPENDENCIES_UTILS_H_
