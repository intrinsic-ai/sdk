#!/bin/sh

# Print workspace status as described on https://bazel.build/docs/user-manual#workspace-status-command
# The actual workspace status variables are expected to be present in the environment and if undefined, this command will error, blocking the build.

# Fail if any environment variables used for status are missing.
set -u

echo STABLE_SDK_VERSION ${CANDIDATE_NAME}
echo STABLE_SEMVER ${SEMVER}
