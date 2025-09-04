// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_CONFIG_ENVIRONMENTS_H_
#define INTRINSIC_CONFIG_ENVIRONMENTS_H_
#include <array>
#include <string>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/types/span.h"

namespace environments {

// Environment constants
inline constexpr char kProd[] = "prod";
inline constexpr char kStaging[] = "staging";
inline constexpr char kDev[] = "dev";

// Accounts project constants
inline constexpr char kAccountsProjectDev[] = "intrinsic-accounts-dev";
inline constexpr char kAccountsProjectStaging[] = "intrinsic-accounts-staging";
inline constexpr char kAccountsProjectProd[] = "intrinsic-accounts-prod";

// Accounts domain constants
inline constexpr char kAccountsDomainDev[] = "accounts-dev.intrinsic.ai";
inline constexpr char kAccountsDomainStaging[] = "accounts-qa.intrinsic.ai";
inline constexpr char kAccountsDomainProd[] = "accounts.intrinsic.ai";

// Portal project constants
inline constexpr char kPortalProjectDev[] = "intrinsic-portal-dev";
inline constexpr char kPortalProjectStaging[] = "intrinsic-portal-staging";
inline constexpr char kPortalProjectProd[] = "intrinsic-portal-prod";

// Portal domain constants
inline constexpr char kPortalDomainDev[] = "flowstate-dev.intrinsic.ai";
inline constexpr char kPortalDomainStaging[] = "flowstate-qa.intrinsic.ai";
inline constexpr char kPortalDomainProd[] = "flowstate.intrinsic.ai";

// Assets project constants
inline constexpr char kAssetsProjectDev[] = "intrinsic-assets-dev";
inline constexpr char kAssetsProjectStaging[] = "intrinsic-assets-staging";
inline constexpr char kAssetsProjectProd[] = "intrinsic-assets-prod";

// Assets domain constants
inline constexpr char kAssetsDomainDev[] = "assets-dev.intrinsic.ai";
inline constexpr char kAssetsDomainStaging[] = "assets-qa.intrinsic.ai";
inline constexpr char kAssetsDomainProd[] = "assets.intrinsic.ai";

namespace internal {
// absl::Span, does not own data, so we need to declare a global constexpr array
// for the span to point to. Make it internal to avoid accidental hard-coded
// deps on the literal size of the array, which we may not want as part of our
// API contract; absl::Span forces users to call absl::Span::size() to get the
// size, which allows us to change the size without breaking existing uses.
inline constexpr std::array<absl::string_view, 3> kAll = {kProd, kStaging,
                                                          kDev};
}  // namespace internal

inline constexpr absl::Span<const absl::string_view> kAll =
    absl::MakeConstSpan(internal::kAll);

absl::StatusOr<std::string> FromDomain(absl::string_view domain);
absl::StatusOr<std::string> FromProject(absl::string_view project);
std::string FromComputeProject(absl::string_view project);

std::string PortalDomain(absl::string_view env);
std::string PortalProject(absl::string_view env);
std::string AccountsDomain(absl::string_view env);
std::string AccountsProjectFromEnv(absl::string_view env);
std::string AccountsProjectFromProject(absl::string_view project);
std::string AssetsDomain(absl::string_view env);
std::string AssetsProject(absl::string_view env);

}  // namespace environments

#endif  // INTRINSIC_CONFIG_ENVIRONMENTS_H_
