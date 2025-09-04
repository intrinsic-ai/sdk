// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/config/environments.h"

#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"

namespace environments {

absl::StatusOr<std::string> FromDomain(absl::string_view domain) {
  if (domain == kPortalDomainProd || domain == kAccountsDomainProd ||
      domain == kAssetsDomainProd) {
    return kProd;
  } else if (domain == kPortalDomainStaging ||
             domain == kAccountsDomainStaging ||
             domain == kAssetsDomainStaging) {
    return kStaging;
  } else if (domain == kPortalDomainDev || domain == kAccountsDomainDev ||
             domain == kAssetsDomainDev) {
    return kDev;
  } else {
    return absl::InvalidArgumentError(absl::StrCat("unknown domain: ", domain));
  }
}

absl::StatusOr<std::string> FromProject(absl::string_view project) {
  if (project == kPortalProjectProd || project == kAccountsProjectProd ||
      project == kAssetsProjectProd) {
    return kProd;
  } else if (project == kPortalProjectStaging ||
             project == kAccountsProjectStaging ||
             project == kAssetsProjectStaging) {
    return kStaging;
  } else if (project == kPortalProjectDev || project == kAccountsProjectDev ||
             project == kAssetsProjectDev) {
    return kDev;
  } else {
    return absl::InvalidArgumentError(
        absl::StrCat("unknown project: ", project));
  }
}

std::string FromComputeProject(absl::string_view project) {
    return kProd;
}

std::string PortalDomain(absl::string_view env) {
  if (env == kProd) {
    return kPortalDomainProd;
  } else if (env == kStaging) {
    return kPortalDomainStaging;
  } else if (env == kDev) {
    return kPortalDomainDev;
  } else {
    return "";
  }
}

std::string AccountsDomain(absl::string_view env) {
  if (env == kProd) {
    return kAccountsDomainProd;
  } else if (env == kStaging) {
    return kAccountsDomainStaging;
  } else if (env == kDev) {
    return kAccountsDomainDev;
  } else {
    return "";
  }
}

std::string AccountsProjectFromEnv(absl::string_view env) {
  if (env == kProd) {
    return kAccountsProjectProd;
  } else if (env == kStaging) {
    return kAccountsProjectStaging;
  } else if (env == kDev) {
    return kAccountsProjectDev;
  } else {
    return "";
  }
}

std::string AccountsProjectFromProject(absl::string_view project) {
  auto result = FromProject(project);
  if (result.ok()) {
    return AccountsProjectFromEnv(result.value());
  } else {
    return AccountsProjectFromEnv(FromComputeProject(project));
  }
}

std::string AssetsDomain(absl::string_view env) {
  if (env == kProd) {
    return kAssetsDomainProd;
  } else if (env == kStaging) {
    return kAssetsDomainStaging;
  } else if (env == kDev) {
    return kAssetsDomainDev;
  } else {
    return "";
  }
}

std::string AssetsProject(absl::string_view env) {
  if (env == kProd) {
    return kAssetsProjectProd;
  } else if (env == kStaging) {
    return kAssetsProjectStaging;
  } else if (env == kDev) {
    return kAssetsProjectDev;
  } else {
    return "";
  }
}

}  // namespace environments
