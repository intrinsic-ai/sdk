// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ASSETS_DATA_FAKE_DATA_ASSETS_H_
#define INTRINSIC_ASSETS_DATA_FAKE_DATA_ASSETS_H_

#include <memory>
#include <string>

#include "absl/container/flat_hash_map.h"
#include "absl/status/statusor.h"
#include "absl/types/span.h"
#include "grpcpp/security/server_credentials.h"
#include "grpcpp/server_builder.h"
#include "grpcpp/server_context.h"
#include "grpcpp/support/status.h"
#include "intrinsic/assets/data/proto/v1/data_assets.grpc.pb.h"
#include "intrinsic/assets/data/proto/v1/data_assets.pb.h"

namespace intrinsic::assets {

// A fake implementation of the DataAssets service for testing.
class FakeDataAssetsService
    : public intrinsic_proto::data::v1::DataAssets::Service {
 public:
  // Static factory method to create and initialize FakeDataAssetsService.
  static absl::StatusOr<std::unique_ptr<FakeDataAssetsService>> Create(
      absl::Span<const intrinsic_proto::data::v1::DataAsset> data_assets,
      int port = 0);

  // Lists all data assets.
  grpc::Status ListDataAssets(
      grpc::ServerContext* context,
      const intrinsic_proto::data::v1::ListDataAssetsRequest* request,
      intrinsic_proto::data::v1::ListDataAssetsResponse* response) override;

  // Gets a specific data asset by its ID.
  grpc::Status GetDataAsset(
      grpc::ServerContext* context,
      const intrinsic_proto::data::v1::GetDataAssetRequest* request,
      intrinsic_proto::data::v1::DataAsset* response) override;

  // Streams the bytes referenced by ReferencedData in an installed Data asset.
  grpc::Status StreamReferencedData(
      grpc::ServerContext* context,
      const intrinsic_proto::data::v1::StreamReferencedDataRequest* request,
      grpc::ServerWriter<
          intrinsic_proto::data::v1::StreamReferencedDataResponse>* writer)
      override;

  // Lists only the metadata of installed Data assets.
  grpc::Status ListDataAssetMetadata(
      grpc::ServerContext* context,
      const intrinsic_proto::data::v1::ListDataAssetMetadataRequest* request,
      intrinsic_proto::data::v1::ListDataAssetMetadataResponse* response)
      override;

  std::unique_ptr<intrinsic_proto::data::v1::DataAssets::Stub>
  NewInternalStub();

 private:
  // Private constructor, called only by the Create factory.
  explicit FakeDataAssetsService(
      const absl::flat_hash_map<
          std::string, intrinsic_proto::data::v1::DataAsset>& data_assets,
      int port);

  absl::flat_hash_map<std::string, intrinsic_proto::data::v1::DataAsset>
      data_assets_;
  std::unique_ptr<grpc::Server> server_;
  int port_;
  std::string address_;
};

}  // namespace intrinsic::assets

#endif  // INTRINSIC_ASSETS_DATA_FAKE_DATA_ASSETS_H_
