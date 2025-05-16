# Copyright 2023 Intrinsic Innovation LLC

from absl.testing import absltest
from intrinsic.config import environments


class EnvironmentsTest(absltest.TestCase):

  def test_all_environments_constant(self):
    self.assertEqual(environments.ALL, ["prod", "staging", "dev"])
    self.assertIn(environments.PROD, environments.ALL)
    self.assertIn(environments.STAGING, environments.ALL)
    self.assertIn(environments.DEV, environments.ALL)

  def test_from_domain_valid(self):
    self.assertEqual(
        environments.from_domain(environments.ACCOUNTS_DOMAIN_PROD),
        environments.PROD,
    )
    self.assertEqual(
        environments.from_domain(environments.PORTAL_DOMAIN_STAGING),
        environments.STAGING,
    )
    self.assertEqual(
        environments.from_domain(environments.ASSETS_DOMAIN_DEV),
        environments.DEV,
    )

  def test_from_domain_invalid(self):
    with self.assertRaisesRegex(
        ValueError, "Unknown domain: invalid.domain.com"
    ):
      environments.from_domain("invalid.domain.com")

  def test_from_project_valid(self):
    self.assertEqual(
        environments.from_project(environments.ACCOUNTS_PROJECT_PROD),
        environments.PROD,
    )
    self.assertEqual(
        environments.from_project(environments.PORTAL_PROJECT_STAGING),
        environments.STAGING,
    )
    self.assertEqual(
        environments.from_project(environments.ASSETS_PROJECT_DEV),
        environments.DEV,
    )

  def test_from_project_invalid(self):
    with self.assertRaisesRegex(ValueError, "Unknown project: invalid-project"):
      environments.from_project("invalid-project")

  def test_from_compute_project(self):
    # from_compute_project currently delegates to from_project
    self.assertEqual(
        environments.from_compute_project(environments.PORTAL_PROJECT_PROD),
        environments.PROD,
    )
    self.assertEqual(
        environments.from_compute_project("invalid-compute-project"),
        environments.PROD,
    )

  def test_portal_domain(self):
    self.assertEqual(
        environments.portal_domain(environments.PROD),
        environments.PORTAL_DOMAIN_PROD,
    )
    self.assertEqual(
        environments.portal_domain(environments.STAGING),
        environments.PORTAL_DOMAIN_STAGING,
    )
    self.assertEqual(
        environments.portal_domain(environments.DEV),
        environments.PORTAL_DOMAIN_DEV,
    )
    with self.assertRaisesRegex(ValueError, "Unknown environment: invalid_env"):
      environments.portal_domain("invalid_env")

  def test_portal_project(self):
    self.assertEqual(
        environments.portal_project(environments.PROD),
        environments.PORTAL_PROJECT_PROD,
    )
    self.assertEqual(
        environments.portal_project(environments.STAGING),
        environments.PORTAL_PROJECT_STAGING,
    )
    self.assertEqual(
        environments.portal_project(environments.DEV),
        environments.PORTAL_PROJECT_DEV,
    )
    with self.assertRaisesRegex(ValueError, "Unknown environment: invalid_env"):
      environments.portal_project("invalid_env")

  def test_accounts_domain(self):
    self.assertEqual(
        environments.accounts_domain(environments.PROD),
        environments.ACCOUNTS_DOMAIN_PROD,
    )
    self.assertEqual(
        environments.accounts_domain(environments.STAGING),
        environments.ACCOUNTS_DOMAIN_STAGING,
    )
    self.assertEqual(
        environments.accounts_domain(environments.DEV),
        environments.ACCOUNTS_DOMAIN_DEV,
    )
    with self.assertRaisesRegex(ValueError, "Unknown environment: invalid_env"):
      environments.accounts_domain("invalid_env")

  def test_accounts_project_from_env(self):
    self.assertEqual(
        environments.accounts_project_from_env(environments.PROD),
        environments.ACCOUNTS_PROJECT_PROD,
    )
    self.assertEqual(
        environments.accounts_project_from_env(environments.STAGING),
        environments.ACCOUNTS_PROJECT_STAGING,
    )
    self.assertEqual(
        environments.accounts_project_from_env(environments.DEV),
        environments.ACCOUNTS_PROJECT_DEV,
    )
    with self.assertRaisesRegex(ValueError, "Unknown environment: invalid_env"):
      environments.accounts_project_from_env("invalid_env")

  def test_accounts_project_from_project(self):
    self.assertEqual(
        environments.accounts_project_from_project(
            environments.PORTAL_PROJECT_PROD
        ),
        environments.ACCOUNTS_PROJECT_PROD,
    )
    self.assertEqual(
        environments.accounts_project_from_project(
            environments.ASSETS_PROJECT_STAGING
        ),
        environments.ACCOUNTS_PROJECT_STAGING,
    )
    self.assertEqual(
        environments.accounts_project_from_project("invalid-project"),
        environments.ACCOUNTS_PROJECT_PROD,
    )

  def test_assets_domain(self):
    self.assertEqual(
        environments.assets_domain(environments.PROD),
        environments.ASSETS_DOMAIN_PROD,
    )
    self.assertEqual(
        environments.assets_domain(environments.STAGING),
        environments.ASSETS_DOMAIN_STAGING,
    )
    self.assertEqual(
        environments.assets_domain(environments.DEV),
        environments.ASSETS_DOMAIN_DEV,
    )
    with self.assertRaisesRegex(ValueError, "Unknown environment: invalid_env"):
      environments.assets_domain("invalid_env")

  def test_assets_project(self):
    self.assertEqual(
        environments.assets_project(environments.PROD),
        environments.ASSETS_PROJECT_PROD,
    )
    self.assertEqual(
        environments.assets_project(environments.STAGING),
        environments.ASSETS_PROJECT_STAGING,
    )
    self.assertEqual(
        environments.assets_project(environments.DEV),
        environments.ASSETS_PROJECT_DEV,
    )
    with self.assertRaisesRegex(ValueError, "Unknown environment: invalid_env"):
      environments.assets_project("invalid_env")


if __name__ == "__main__":
  absltest.main()
