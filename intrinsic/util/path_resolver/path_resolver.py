# Copyright 2023 Intrinsic Innovation LLC

"""Generic resolver that will attempt to resolve runfiles or external directories if available."""

import os

from python.runfiles import runfiles

_repo_name = 'ai_intrinsic_sdks'
def rlocation(path: str) -> str:
  """Resolves the runfiles path for the given rlocationpath."""
  return os.path.normpath(
      runfiles.Create().Rlocation(path),
  )
def resolve_runfiles_path(path: str) -> str:
  """Returns the runfiles path for the given path."""
  if os.path.isabs(path):
    return path

  return rlocation(os.path.join(_repo_name, path))
