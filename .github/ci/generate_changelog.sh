#!/usr/bin/env bash

# Function to generate the Bazel dependency and archive override
generate_bazel_dep() {
  local version="$1"
  local url
  local strip_prefix

  # If no argument is provided, default to "main"
  if [[ -z "$version" ]]; then
    version="main"
  fi

  # First case is user passes main/latest
  if [[ "$version" == "main" || "$version" == "latest" ]]; then
    url="https://github.com/intrinsic-ai/sdk/archive/refs/heads/main.tar.gz"
    strip_prefix="sdk-main/"
  else
    if [[ "$version" =~ ^candidate/ ]]; then
      strip_prefix="sdk-$(echo "$version" | sed 's|/|-|g')" # Replace / with -
      url="https://github.com/intrinsic-ai/sdk/archive/refs/tags/${version}.tar.gz"
    else
     # Ensure the version starts with "v"
      if [[ ! "$version" =~ ^v ]]; then
        version="v$version"
      fi

      strip_prefix="sdk-${version#v}/" # remove 'v' prefix
      url="https://github.com/intrinsic-ai/sdk/archive/refs/tags/${version}.tar.gz"
    fi
  fi

  # Output the Bazel code
  cat <<EOF
bazel_dep(name = "ai_intrinsic_sdks")
archive_override(
    module_name = "ai_intrinsic_sdks",
    urls = ["${url}"],
    strip_prefix = "${strip_prefix}"
)
EOF
}

generate_devcontainer_config() {
  local version="$1"

  # Remove leading 'v' if present
  version="${version#v}"

  cat <<EOF
{
    "name": "intrinsic-flowstate-devcontainer",
    "image": "ghcr.io/intrinsic-ai/intrinsic-dev-img:${version}",
    "runArgs": [
        "--network=host"
    ],
    "customizations": {
        "vscode": {
            "settings": {
                "intrinsic.defaultSdkRepository": "https://github.com/intrinsic-ai/sdk.git"
            }
        }
    }
}
EOF
}

validate_version() {
  local version="$1"

  # Check if version is blank, "main", "latest", or a valid tag
  if [[ -z "$version" || "$version" == "main" || "$version" == "latest" || "$version" =~ ^[v]?[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    return 0 # Valid version
  else
    return 1 # Invalid version
  fi
}


generate_changelog()
{
  local version="$1"

  echo "## bazelmod configuration"
  echo
  echo "Update your MODULE.bazel file to use the newest release archive:"
  echo
  echo "\`\`\`"
  generate_bazel_dep "$version"
  echo "\`\`\`"

  echo
  echo

  echo "## devcontainer configuration"
  echo
  echo "Update your devcontainer configuration to use the latest base image:"
  echo
  echo "\`\`\`"
  generate_devcontainer_config "$version"
  echo "\`\`\`"
}

if [[ -z "$1" ]]; then
  generate_changelog "latest" # or main, or a default version.
else
  if validate_version "$1"; then
    generate_changelog "$1"
  else
    echo "Error: Invalid version format. Must be blank, 'main', 'latest', or a tag like v1.2.3 or 1.2.3"
    exit 1
  fi
fi
