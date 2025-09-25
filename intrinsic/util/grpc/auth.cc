// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/grpc/auth.h"

#include <cstdlib>
#include <fstream>
#include <map>
#include <string>
#include <string_view>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "nlohmann/json.hpp"

namespace intrinsic {
namespace auth {

using json = nlohmann::json;

constexpr std::string_view kStoreDirectory = "intrinsic/projects";
constexpr std::string_view kAuthConfigExtension = ".user-token";
constexpr std::string_view kAliasDefaultToken = "default";

absl::StatusOr<std::multimap<std::string, std::string>> GetRequestMetadata(
    std::string_view project_name) {
  const char* home_dir = std::getenv("HOME");
  if (home_dir == nullptr) {
    return absl::NotFoundError("$HOME environment variable not set.");
  }

  std::string file_name = absl::StrCat(home_dir, "/.config/", kStoreDirectory,
                                       "/", project_name, kAuthConfigExtension);

  std::ifstream f(file_name);
  if (!f.is_open()) {
    return absl::NotFoundError(absl::StrCat("Could not open ", file_name));
  }

  json data = json::parse(f, /*cb=*/nullptr, /*allow_exceptions=*/false);
  if (data.is_discarded()) {
    return absl::FailedPreconditionError(
        absl::StrCat("Could not parse ", file_name));
  }

  if (!data.contains("tokens") ||
      !data["tokens"].contains(kAliasDefaultToken) ||
      !data["tokens"][kAliasDefaultToken].contains("apiKey")) {
    return absl::FailedPreconditionError(absl::StrCat(
        "Could not find default token in ", file_name,
        ". Please run 'inctl auth login --project ", project_name, "'"));
  }

  if (!data["tokens"][kAliasDefaultToken]["apiKey"].is_string()) {
    return absl::FailedPreconditionError(
        absl::StrCat("apiKey in default token is not a string in ", file_name));
  }

  std::string api_key = data["tokens"][kAliasDefaultToken]["apiKey"];
  std::multimap<std::string, std::string> metadata;
  metadata.insert({"authorization", absl::StrCat("Bearer ", api_key)});
  return metadata;
}

}  // namespace auth
}  // namespace intrinsic
