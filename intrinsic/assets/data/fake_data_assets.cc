// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/assets/data/fake_data_assets.h"

#include <algorithm>
#include <cstddef>
#include <memory>
#include <string>
#include <vector>

#include "absl/container/flat_hash_map.h"
#include "absl/log/log.h"
#include "absl/memory/memory.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/types/span.h"
#include "grpcpp/security/server_credentials.h"
#include "grpcpp/server.h"
#include "grpcpp/server_builder.h"
#include "grpcpp/server_context.h"
#include "grpcpp/support/channel_arguments.h"
#include "grpcpp/support/status.h"
#include "intrinsic/assets/data/proto/v1/data_assets.grpc.pb.h"
#include "intrinsic/assets/data/proto/v1/data_assets.pb.h"
#include "intrinsic/assets/id_utils.h"
#include "intrinsic/util/status/ret_check.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/util/status/status_macros_grpc.h"

namespace intrinsic::assets {
namespace {
constexpr int kDefaultPageSize = 20;
}  // namespace

absl::StatusOr<std::unique_ptr<FakeDataAssetsService>>
FakeDataAssetsService::Create(
    absl::Span<const intrinsic_proto::data::v1::DataAsset> data_assets,
    int port) {
  absl::flat_hash_map<std::string, intrinsic_proto::data::v1::DataAsset>
      asset_map;
  for (const auto& asset : data_assets) {
    INTR_ASSIGN_OR_RETURN(const std::string asset_id,
                          IdFromProto(asset.metadata().id_version().id()));

    INTR_RET_CHECK(!asset_map.contains(asset_id))
        << "Duplicate DataAsset id found: " << asset_id;
    asset_map[asset_id] = asset;
  }
  return absl::WrapUnique(new FakeDataAssetsService(asset_map, port));
}

FakeDataAssetsService::FakeDataAssetsService(
    const absl::flat_hash_map<
        std::string, intrinsic_proto::data::v1::DataAsset>& data_assets,
    int port)
    : data_assets_(data_assets),
      port_(port),
      address_(absl::StrCat("dns:///localhost:", port_)) {
  server_ =
      grpc::ServerBuilder()
          .RegisterService(this)
          .AddListeningPort(
              address_, grpc::InsecureServerCredentials())  // NOLINT (insecure)
          .BuildAndStart();
}

grpc::Status FakeDataAssetsService::ListDataAssets(
    grpc::ServerContext* context,
    const intrinsic_proto::data::v1::ListDataAssetsRequest* request,
    intrinsic_proto::data::v1::ListDataAssetsResponse* response) {
  std::vector<intrinsic_proto::data::v1::DataAsset> filtered_assets;
  for (const auto& [id, asset] : data_assets_) {
    if (request->has_strict_filter()) {
      if (request->strict_filter().has_proto_name()) {
        absl::string_view type_url = asset.data().type_url();
        size_t last_slash = type_url.find_last_of('/');
        absl::string_view proto_name = (last_slash == absl::string_view::npos)
                                           ? type_url
                                           : type_url.substr(last_slash + 1);
        if (request->strict_filter().proto_name() != proto_name) {
          continue;
        }
      }
    }
    filtered_assets.push_back(asset);
  }

  // Sort by ID for consistent pagination.
  std::sort(filtered_assets.begin(), filtered_assets.end(),
            [](const intrinsic_proto::data::v1::DataAsset& a,
               const intrinsic_proto::data::v1::DataAsset& b) {
              const auto& id_a = a.metadata().id_version().id();
              const auto& id_b = b.metadata().id_version().id();
              if (id_a.package() != id_b.package()) {
                return id_a.package() < id_b.package();
              }
              return id_a.name() < id_b.name();
            });

  auto it = filtered_assets.begin();
  if (!request->page_token().empty()) {
    it = std::lower_bound(
        filtered_assets.begin(), filtered_assets.end(), request->page_token(),
        [](const intrinsic_proto::data::v1::DataAsset& asset,
           absl::string_view token) {
          return IdFromProto(asset.metadata().id_version().id()).value() <
                 token;
        });
    // If the token matches an asset, start from the *next* one.
    if (it != filtered_assets.end()) {
      absl::StatusOr<std::string> current_asset_id_str =
          IdFromProto(it->metadata().id_version().id());
      if (current_asset_id_str.ok() &&
          *current_asset_id_str == request->page_token()) {
        ++it;
      }
    }
  }

  int page_size =
      request->page_size() > 0 ? request->page_size() : kDefaultPageSize;
  int count = 0;
  while (it != filtered_assets.end() && count < page_size) {
    *response->add_data_assets() = *it;
    ++it;
    ++count;
  }

  if (it != filtered_assets.end()) {
    absl::StatusOr<std::string> next_page_token = IdFromProto(
        response->data_assets().rbegin()->metadata().id_version().id());
    if (next_page_token.ok()) {
      response->set_next_page_token(*next_page_token);
    }
  }

  return grpc::Status::OK;
}

grpc::Status FakeDataAssetsService::GetDataAsset(
    grpc::ServerContext* context,
    const intrinsic_proto::data::v1::GetDataAssetRequest* request,
    intrinsic_proto::data::v1::DataAsset* response) {
  INTR_ASSIGN_OR_RETURN_GRPC(const std::string asset_id,
                             IdFromProto(request->id()));
  auto it = data_assets_.find(asset_id);
  if (it == data_assets_.end()) {
    return ToGrpcStatus(absl::NotFoundError(
        absl::StrCat("DataAsset with id '", asset_id, "' not found.")));
  }
  *response = it->second;
  return grpc::Status::OK;
}

grpc::Status FakeDataAssetsService::ListDataAssetMetadata(
    grpc::ServerContext* context,
    const intrinsic_proto::data::v1::ListDataAssetMetadataRequest* request,
    intrinsic_proto::data::v1::ListDataAssetMetadataResponse* response) {
  intrinsic_proto::data::v1::ListDataAssetsRequest list_assets_request;
  if (request->has_strict_filter()) {
    *list_assets_request.mutable_strict_filter() = request->strict_filter();
  }
  list_assets_request.set_page_size(request->page_size());
  list_assets_request.set_page_token(request->page_token());

  intrinsic_proto::data::v1::ListDataAssetsResponse list_assets_response;
  INTR_RETURN_IF_ERROR_GRPC(
      ListDataAssets(context, &list_assets_request, &list_assets_response));

  for (const auto& asset : list_assets_response.data_assets()) {
    *response->add_metadata() = asset.metadata();
  }
  if (!list_assets_response.next_page_token().empty()) {
    response->set_next_page_token(list_assets_response.next_page_token());
  }
  return grpc::Status::OK;
}

grpc::Status FakeDataAssetsService::StreamReferencedData(
    grpc::ServerContext* context,
    const intrinsic_proto::data::v1::StreamReferencedDataRequest* request,
    grpc::ServerWriter<intrinsic_proto::data::v1::StreamReferencedDataResponse>*
        writer) {
  return ToGrpcStatus(absl::UnimplementedError(
      "StreamReferencedData is not implemented in FakeDataAssetsService."));
}

std::unique_ptr<intrinsic_proto::data::v1::DataAssets::Stub>
FakeDataAssetsService::NewInternalStub() {
  return intrinsic_proto::data::v1::DataAssets::NewStub(
      server_->InProcessChannel(grpc::ChannelArguments()));
}

}  // namespace intrinsic::assets
