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
  EXPECT_EQ(FromDomain(PortalDomainProd).value(), Prod);
  EXPECT_EQ(FromDomain(AccountsDomainProd).value(), Prod);
  EXPECT_EQ(FromDomain(AssetsDomainProd).value(), Prod);
}

TEST(EnvironmentsTest, FromDomainStaging) {
  EXPECT_EQ(FromDomain(PortalDomainStaging).value(), Staging);
  EXPECT_EQ(FromDomain(AccountsDomainStaging).value(), Staging);
  EXPECT_EQ(FromDomain(AssetsDomainStaging).value(), Staging);
}

TEST(EnvironmentsTest, FromDomainDev) {
  EXPECT_EQ(FromDomain(PortalDomainDev).value(), Dev);
  EXPECT_EQ(FromDomain(AccountsDomainDev).value(), Dev);
  EXPECT_EQ(FromDomain(AssetsDomainDev).value(), Dev);
}

TEST(EnvironmentsTest, FromDomainInvalid) {
  absl::StatusOr<std::string> result = FromDomain("invalid_domain");
  EXPECT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
  EXPECT_THAT(result.status().message(),
              HasSubstr("unknown domain: invalid_domain"));
}

TEST(EnvironmentsTest, FromProjectProd) {
  EXPECT_EQ(FromProject(PortalProjectProd).value(), Prod);
  EXPECT_EQ(FromProject(AccountsProjectProd).value(), Prod);
  EXPECT_EQ(FromProject(AssetsProjectProd).value(), Prod);
}

TEST(EnvironmentsTest, FromProjectStaging) {
  EXPECT_EQ(FromProject(PortalProjectStaging).value(), Staging);
  EXPECT_EQ(FromProject(AccountsProjectStaging).value(), Staging);
  EXPECT_EQ(FromProject(AssetsProjectStaging).value(), Staging);
}

TEST(EnvironmentsTest, FromProjectDev) {
  EXPECT_EQ(FromProject(PortalProjectDev).value(), Dev);
  EXPECT_EQ(FromProject(AccountsProjectDev).value(), Dev);
  EXPECT_EQ(FromProject(AssetsProjectDev).value(), Dev);
}

TEST(EnvironmentsTest, FromProjectInvalid) {
  absl::StatusOr<std::string> result = FromProject("invalid_project");
  EXPECT_FALSE(result.ok());
  EXPECT_EQ(result.status().code(), absl::StatusCode::kInvalidArgument);
  EXPECT_THAT(result.status().message(),
              HasSubstr("unknown project: invalid_project"));
}

TEST(EnvironmentsTest, FromComputeProjectProd) {
  EXPECT_EQ(FromComputeProject("some_other_project"), Prod);
}

TEST(EnvironmentsTest, PortalDomain) {
  EXPECT_EQ(PortalDomain(Prod), PortalDomainProd);
  EXPECT_EQ(PortalDomain(Staging), PortalDomainStaging);
  EXPECT_EQ(PortalDomain(Dev), PortalDomainDev);
  EXPECT_EQ(PortalDomain("invalid_env"), "");
}

TEST(EnvironmentsTest, AccountsDomain) {
  EXPECT_EQ(AccountsDomain(Prod), AccountsDomainProd);
  EXPECT_EQ(AccountsDomain(Staging), AccountsDomainStaging);
  EXPECT_EQ(AccountsDomain(Dev), AccountsDomainDev);
  EXPECT_EQ(AccountsDomain("invalid_env"), "");
}

TEST(EnvironmentsTest, AccountsProjectFromEnv) {
  EXPECT_EQ(AccountsProjectFromEnv(Prod), AccountsProjectProd);
  EXPECT_EQ(AccountsProjectFromEnv(Staging), AccountsProjectStaging);
  EXPECT_EQ(AccountsProjectFromEnv(Dev), AccountsProjectDev);
  EXPECT_EQ(AccountsProjectFromEnv("invalid_env"), "");
}

TEST(EnvironmentsTest, AccountsProjectFromProject) {
  EXPECT_EQ(AccountsProjectFromProject(PortalProjectProd), AccountsProjectProd);
  EXPECT_EQ(AccountsProjectFromProject(PortalProjectStaging),
            AccountsProjectStaging);
  EXPECT_EQ(AccountsProjectFromProject(PortalProjectDev), AccountsProjectDev);

  EXPECT_EQ(AccountsProjectFromProject(AccountsProjectProd),
            AccountsProjectProd);
  EXPECT_EQ(AccountsProjectFromProject(AccountsProjectStaging),
            AccountsProjectStaging);
  EXPECT_EQ(AccountsProjectFromProject(AccountsProjectDev), AccountsProjectDev);

  EXPECT_EQ(AccountsProjectFromProject(AssetsProjectProd), AccountsProjectProd);
  EXPECT_EQ(AccountsProjectFromProject(AssetsProjectStaging),
            AccountsProjectStaging);
  EXPECT_EQ(AccountsProjectFromProject(AssetsProjectDev), AccountsProjectDev);

  EXPECT_EQ(AccountsProjectFromProject("some_other_project"),
            AccountsProjectProd);
}

TEST(EnvironmentsTest, AssetsDomain) {
  EXPECT_EQ(AssetsDomain(Prod), AssetsDomainProd);
  EXPECT_EQ(AssetsDomain(Staging), AssetsDomainStaging);
  EXPECT_EQ(AssetsDomain(Dev), AssetsDomainDev);
  EXPECT_EQ(AssetsDomain("invalid_env"), "");
}

TEST(EnvironmentsTest, AssetsProject) {
  EXPECT_EQ(AssetsProject(Prod), AssetsProjectProd);
  EXPECT_EQ(AssetsProject(Staging), AssetsProjectStaging);
  EXPECT_EQ(AssetsProject(Dev), AssetsProjectDev);
  EXPECT_EQ(AssetsProject("invalid_env"), "");
}
}  // namespace environments
