// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/config/environments.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

namespace environments {
using ::testing::HasSubstr;

TEST(EnvironmentsTest, FromDomainProd) {
  EXPECT_EQ(FromDomain(kPortalDomainProd).value(), kProd);
  EXPECT_EQ(FromDomain(kAccountsDomainProd).value(), kProd);
  EXPECT_EQ(FromDomain(kAssetsDomainProd).value(), kProd);
}

TEST(EnvironmentsTest, FromDomainStaging) {
  EXPECT_EQ(FromDomain(kPortalDomainStaging).value(), kStaging);
  EXPECT_EQ(FromDomain(kAccountsDomainStaging).value(), kStaging);
  EXPECT_EQ(FromDomain(kAssetsDomainStaging).value(), kStaging);
}

TEST(EnvironmentsTest, FromDomainDev) {
  EXPECT_EQ(FromDomain(kPortalDomainDev).value(), kDev);
  EXPECT_EQ(FromDomain(kAccountsDomainDev).value(), kDev);
  EXPECT_EQ(FromDomain(kAssetsDomainDev).value(), kDev);
}

TEST(EnvironmentsTest, FromDomainInvalid) {
  absl::StatusOr<std::string> result = FromDomain("invalid_domain");
  EXPECT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
  EXPECT_THAT(result.status().message(),
              HasSubstr("unknown domain: invalid_domain"));
}

TEST(EnvironmentsTest, FromProjectProd) {
  EXPECT_EQ(FromProject(kPortalProjectProd).value(), kProd);
  EXPECT_EQ(FromProject(kAccountsProjectProd).value(), kProd);
  EXPECT_EQ(FromProject(kAssetsProjectProd).value(), kProd);
}

TEST(EnvironmentsTest, FromProjectStaging) {
  EXPECT_EQ(FromProject(kPortalProjectStaging).value(), kStaging);
  EXPECT_EQ(FromProject(kAccountsProjectStaging).value(), kStaging);
  EXPECT_EQ(FromProject(kAssetsProjectStaging).value(), kStaging);
}

TEST(EnvironmentsTest, FromProjectDev) {
  EXPECT_EQ(FromProject(kPortalProjectDev).value(), kDev);
  EXPECT_EQ(FromProject(kAccountsProjectDev).value(), kDev);
  EXPECT_EQ(FromProject(kAssetsProjectDev).value(), kDev);
}

TEST(EnvironmentsTest, FromProjectInvalid) {
  absl::StatusOr<std::string> result = FromProject("invalid_project");
  EXPECT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
  EXPECT_THAT(result.status().message(),
              HasSubstr("unknown project: invalid_project"));
}

TEST(EnvironmentsTest, FromComputeProjectProd) {
  EXPECT_EQ(FromComputeProject("some_other_project"), kProd);
}

TEST(EnvironmentsTest, PortalDomain) {
  EXPECT_EQ(PortalDomain(kProd), kPortalDomainProd);
  EXPECT_EQ(PortalDomain(kStaging), kPortalDomainStaging);
  EXPECT_EQ(PortalDomain(kDev), kPortalDomainDev);
  EXPECT_EQ(PortalDomain("invalid_env"), "");
}

TEST(EnvironmentsTest, AccountsDomain) {
  EXPECT_EQ(AccountsDomain(kProd), kAccountsDomainProd);
  EXPECT_EQ(AccountsDomain(kStaging), kAccountsDomainStaging);
  EXPECT_EQ(AccountsDomain(kDev), kAccountsDomainDev);
  EXPECT_EQ(AccountsDomain("invalid_env"), "");
}

TEST(EnvironmentsTest, AccountsProjectFromEnv) {
  EXPECT_EQ(AccountsProjectFromEnv(kProd), kAccountsProjectProd);
  EXPECT_EQ(AccountsProjectFromEnv(kStaging), kAccountsProjectStaging);
  EXPECT_EQ(AccountsProjectFromEnv(kDev), kAccountsProjectDev);
  EXPECT_EQ(AccountsProjectFromEnv("invalid_env"), "");
}

TEST(EnvironmentsTest, AccountsProjectFromProject) {
  EXPECT_EQ(AccountsProjectFromProject(kPortalProjectProd),
            kAccountsProjectProd);
  EXPECT_EQ(AccountsProjectFromProject(kPortalProjectStaging),
            kAccountsProjectStaging);
  EXPECT_EQ(AccountsProjectFromProject(kPortalProjectDev), kAccountsProjectDev);

  EXPECT_EQ(AccountsProjectFromProject(kAccountsProjectProd),
            kAccountsProjectProd);
  EXPECT_EQ(AccountsProjectFromProject(kAccountsProjectStaging),
            kAccountsProjectStaging);
  EXPECT_EQ(AccountsProjectFromProject(kAccountsProjectDev),
            kAccountsProjectDev);

  EXPECT_EQ(AccountsProjectFromProject(kAssetsProjectProd),
            kAccountsProjectProd);
  EXPECT_EQ(AccountsProjectFromProject(kAssetsProjectStaging),
            kAccountsProjectStaging);
  EXPECT_EQ(AccountsProjectFromProject(kAssetsProjectDev), kAccountsProjectDev);

  EXPECT_EQ(AccountsProjectFromProject("some_other_project"),
            kAccountsProjectProd);
}

TEST(EnvironmentsTest, AssetsDomain) {
  EXPECT_EQ(AssetsDomain(kProd), kAssetsDomainProd);
  EXPECT_EQ(AssetsDomain(kStaging), kAssetsDomainStaging);
  EXPECT_EQ(AssetsDomain(kDev), kAssetsDomainDev);
  EXPECT_EQ(AssetsDomain("invalid_env"), "");
}

TEST(EnvironmentsTest, AssetsProject) {
  EXPECT_EQ(AssetsProject(kProd), kAssetsProjectProd);
  EXPECT_EQ(AssetsProject(kStaging), kAssetsProjectStaging);
  EXPECT_EQ(AssetsProject(kDev), kAssetsProjectDev);
  EXPECT_EQ(AssetsProject("invalid_env"), "");
}

TEST(EnvironmentsTest, All) {
  EXPECT_THAT(kAll, testing::UnorderedElementsAre(kProd, kStaging, kDev));
}

}  // namespace environments
