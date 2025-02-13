// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/config/environments.h"

#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"

namespace environments {

absl::StatusOr<std::string> FromDomain(const std::string& domain) {
  if (domain == PortalDomainProd || domain == AccountsDomainProd ||
      domain == AssetsDomainProd) {
    return Prod;
  } else if (domain == PortalDomainStaging || domain == AccountsDomainStaging ||
             domain == AssetsDomainStaging) {
    return Staging;
  } else if (domain == PortalDomainDev || domain == AccountsDomainDev ||
             domain == AssetsDomainDev) {
    return Dev;
  } else {
    return absl::InvalidArgumentError("unknown domain: " + domain);
  }
}

absl::StatusOr<std::string> FromProject(const std::string& project) {
  if (project == PortalProjectProd || project == AccountsProjectProd ||
      project == AssetsProjectProd) {
    return Prod;
  } else if (project == PortalProjectStaging ||
             project == AccountsProjectStaging ||
             project == AssetsProjectStaging) {
    return Staging;
  } else if (project == PortalProjectDev || project == AccountsProjectDev ||
             project == AssetsProjectDev) {
    return Dev;
  } else {
    return absl::InvalidArgumentError("unknown project: " + project);
  }
}

std::string FromComputeProject(const std::string& project) {
    return Prod;
}

std::string PortalDomain(const std::string& env) {
  if (env == Prod) {
    return PortalDomainProd;
  } else if (env == Staging) {
    return PortalDomainStaging;
  } else if (env == Dev) {
    return PortalDomainDev;
  } else {
    return "";
  }
}

std::string AccountsDomain(const std::string& env) {
  if (env == Prod) {
    return AccountsDomainProd;
  } else if (env == Staging) {
    return AccountsDomainStaging;
  } else if (env == Dev) {
    return AccountsDomainDev;
  } else {
    return "";
  }
}

std::string AccountsProjectFromEnv(const std::string& env) {
  if (env == Prod) {
    return AccountsProjectProd;
  } else if (env == Staging) {
    return AccountsProjectStaging;
  } else if (env == Dev) {
    return AccountsProjectDev;
  } else {
    return "";
  }
}

std::string AccountsProjectFromProject(const std::string& project) {
  auto result = FromProject(project);
  if (result.ok()) {
    return AccountsProjectFromEnv(result.value());
  } else {
    return AccountsProjectFromEnv(FromComputeProject(project));
  }
}

std::string AssetsDomain(const std::string& env) {
  if (env == Prod) {
    return AssetsDomainProd;
  } else if (env == Staging) {
    return AssetsDomainStaging;
  } else if (env == Dev) {
    return AssetsDomainDev;
  } else {
    return "";
  }
}

std::string AssetsProject(const std::string& env) {
  if (env == Prod) {
    return AssetsProjectProd;
  } else if (env == Staging) {
    return AssetsProjectStaging;
  } else if (env == Dev) {
    return AssetsProjectDev;
  } else {
    return "";
  }
}

}  // namespace environments
