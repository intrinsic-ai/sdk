# Copyright 2023 Intrinsic Innovation LLC

"""Helper functions to work with environments."""

# Environment constants
PROD = "prod"
STAGING = "staging"
DEV = "dev"

# Accounts project constants
ACCOUNTS_PROJECT_DEV = "intrinsic-accounts-dev"
ACCOUNTS_PROJECT_STAGING = "intrinsic-accounts-staging"
ACCOUNTS_PROJECT_PROD = "intrinsic-accounts-prod"

# Accounts domain constants
ACCOUNTS_DOMAIN_DEV = "accounts-dev.intrinsic.ai"
ACCOUNTS_DOMAIN_STAGING = "accounts-qa.intrinsic.ai"
ACCOUNTS_DOMAIN_PROD = "accounts.intrinsic.ai"

# Portal project constants
PORTAL_PROJECT_DEV = "intrinsic-portal-dev"
PORTAL_PROJECT_STAGING = "intrinsic-portal-staging"
PORTAL_PROJECT_PROD = "intrinsic-portal-prod"

# Portal domain constants
PORTAL_DOMAIN_DEV = "flowstate-dev.intrinsic.ai"
PORTAL_DOMAIN_STAGING = "flowstate-qa.intrinsic.ai"
PORTAL_DOMAIN_PROD = "flowstate.intrinsic.ai"

# Assets project constants
ASSETS_PROJECT_DEV = "intrinsic-assets-dev"
ASSETS_PROJECT_STAGING = "intrinsic-assets-staging"
ASSETS_PROJECT_PROD = "intrinsic-assets-prod"

# Assets domain constants
ASSETS_DOMAIN_DEV = "assets-dev.intrinsic.ai"
ASSETS_DOMAIN_STAGING = "assets-qa.intrinsic.ai"
ASSETS_DOMAIN_PROD = "assets.intrinsic.ai"

# All environments
ALL = [PROD, STAGING, DEV]


def from_domain(domain: str) -> str:
  """Returns the environment for a given domain."""
  if domain in (ACCOUNTS_DOMAIN_PROD, PORTAL_DOMAIN_PROD, ASSETS_DOMAIN_PROD):
    return PROD
  if domain in (
      ACCOUNTS_DOMAIN_STAGING,
      PORTAL_DOMAIN_STAGING,
      ASSETS_DOMAIN_STAGING,
  ):
    return STAGING
  if domain in (ACCOUNTS_DOMAIN_DEV, PORTAL_DOMAIN_DEV, ASSETS_DOMAIN_DEV):
    return DEV
  raise ValueError(f"Unknown domain: {domain}")


def from_project(project: str) -> str:
  """Returns the environment for a given project."""
  if project in (
      ACCOUNTS_PROJECT_PROD,
      PORTAL_PROJECT_PROD,
      ASSETS_PROJECT_PROD,
  ):
    return PROD
  if project in (
      ACCOUNTS_PROJECT_STAGING,
      PORTAL_PROJECT_STAGING,
      ASSETS_PROJECT_STAGING,
  ):
    return STAGING
  if project in (
      ACCOUNTS_PROJECT_DEV,
      PORTAL_PROJECT_DEV,
      ASSETS_PROJECT_DEV,
  ):
    return DEV
  raise ValueError(f"Unknown project: {project}")


def from_compute_project(project: str) -> str:
  """Returns the environment for a given compute project."""
  return PROD


def portal_domain(env: str) -> str:
  """Returns the portal domain for a given environment."""
  if env == PROD:
    return PORTAL_DOMAIN_PROD
  if env == STAGING:
    return PORTAL_DOMAIN_STAGING
  if env == DEV:
    return PORTAL_DOMAIN_DEV
  raise ValueError(f"Unknown environment: {env}")


def portal_project(env: str) -> str:
  """Returns the portal project for a given environment."""
  if env == PROD:
    return PORTAL_PROJECT_PROD
  if env == STAGING:
    return PORTAL_PROJECT_STAGING
  if env == DEV:
    return PORTAL_PROJECT_DEV
  raise ValueError(f"Unknown environment: {env}")


def accounts_domain(env: str) -> str:
  """Returns the accounts domain for a given environment."""
  if env == PROD:
    return ACCOUNTS_DOMAIN_PROD
  if env == STAGING:
    return ACCOUNTS_DOMAIN_STAGING
  if env == DEV:
    return ACCOUNTS_DOMAIN_DEV
  raise ValueError(f"Unknown environment: {env}")


def accounts_project_from_env(env: str) -> str:
  """Returns the accounts project for a given environment."""
  if env == PROD:
    return ACCOUNTS_PROJECT_PROD
  if env == STAGING:
    return ACCOUNTS_PROJECT_STAGING
  if env == DEV:
    return ACCOUNTS_PROJECT_DEV
  raise ValueError(f"Unknown environment: {env}")


def accounts_project_from_project(project: str) -> str:
  """Returns the accounts project for a given project."""
  try:
    env = from_project(project)
  except ValueError:
    env = from_compute_project(project)
  return accounts_project_from_env(env)


def assets_domain(env: str) -> str:
  """Returns the assets domain for a given environment."""
  if env == PROD:
    return ASSETS_DOMAIN_PROD
  if env == STAGING:
    return ASSETS_DOMAIN_STAGING
  if env == DEV:
    return ASSETS_DOMAIN_DEV
  raise ValueError(f"Unknown environment: {env}")


def assets_project(env: str) -> str:
  """Returns the assets project for a given environment."""
  if env == PROD:
    return ASSETS_PROJECT_PROD
  if env == STAGING:
    return ASSETS_PROJECT_STAGING
  if env == DEV:
    return ASSETS_PROJECT_DEV
  raise ValueError(f"Unknown environment: {env}")
