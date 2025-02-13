// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_CONFIG_ENVIRONMENTS_H_
#define INTRINSIC_CONFIG_ENVIRONMENTS_H_
#include <string>
#include <vector>

#include "absl/status/statusor.h"

namespace environments {

// Environment constants
constexpr char Prod[] = "prod";
constexpr char Staging[] = "staging";
constexpr char Dev[] = "dev";

// Accounts project constants
constexpr char AccountsProjectDev[] = "intrinsic-accounts-dev";
constexpr char AccountsProjectStaging[] = "intrinsic-accounts-staging";
constexpr char AccountsProjectProd[] = "intrinsic-accounts-prod";

// Accounts domain constants
constexpr char AccountsDomainDev[] = "accounts-dev.intrinsic.ai";
constexpr char AccountsDomainStaging[] = "accounts-qa.intrinsic.ai";
constexpr char AccountsDomainProd[] = "accounts.intrinsic.ai";

// Portal project constants
constexpr char PortalProjectDev[] = "intrinsic-portal-dev";
constexpr char PortalProjectStaging[] = "intrinsic-portal-staging";
constexpr char PortalProjectProd[] = "intrinsic-portal-prod";

// Portal domain constants
constexpr char PortalDomainDev[] = "flowstate-dev.intrinsic.ai";
constexpr char PortalDomainStaging[] = "flowstate-qa.intrinsic.ai";
constexpr char PortalDomainProd[] = "flowstate.intrinsic.ai";

// Assets project constants
constexpr char AssetsProjectDev[] = "intrinsic-assets-dev";
constexpr char AssetsProjectStaging[] = "intrinsic-assets-staging";
constexpr char AssetsProjectProd[] = "intrinsic-assets-prod";

// Assets domain constants
constexpr char AssetsDomainDev[] = "assets-dev.intrinsic.ai";
constexpr char AssetsDomainStaging[] = "assets-qa.intrinsic.ai";
constexpr char AssetsDomainProd[] = "assets.intrinsic.ai";

// All environments
extern const std::vector<std::string> All;

absl::StatusOr<std::string> FromDomain(const std::string& domain);
absl::StatusOr<std::string> FromProject(const std::string& project);
std::string FromComputeProject(const std::string& project);

std::string PortalDomain(const std::string& env);
std::string PortalProject(const std::string& env);
std::string AccountsDomain(const std::string& env);
std::string AccountsProjectFromEnv(const std::string& env);
std::string AccountsProjectFromProject(const std::string& project);
std::string AssetsDomain(const std::string& env);
std::string AssetsProject(const std::string& env);

}  // namespace environments

#endif  // INTRINSIC_CONFIG_ENVIRONMENTS_H_
