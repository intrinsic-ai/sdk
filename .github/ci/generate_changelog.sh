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
    local dev_image_base_name
    local dev_image_tag

    if [[ -z "$version" ]]; then # Default if called directly with no version
        version="latest" # Or align with generate_changelog's default
    fi

    if [[ "$version" == "main" || "$version" == "latest" ]]; then
        dev_image_base_name="ghcr.io/intrinsic-ai/intrinsic-dev-img"
        dev_image_tag="$version"
    elif [[ "$version" =~ ^candidate/intrinsic\.platform\.([0-9]{8})\.(RC[0-9]+)$ ]]; then
        # New candidate format: candidate/intrinsic.platform.20250512.RC05
        # Dev image: ghcr.io/intrinsic-ai/intrinsic-dev-img-staging:0.20250512.0-RC05
        dev_image_base_name="ghcr.io/intrinsic-ai/intrinsic-dev-img-staging"
        local date_part="${BASH_REMATCH[1]}" # 20250512
        local rc_part="${BASH_REMATCH[2]}"   # RC05
        dev_image_tag="0.${date_part}.0-${rc_part}"
    else
        # Standard versions like v1.2.3 or 1.2.3
        dev_image_base_name="ghcr.io/intrinsic-ai/intrinsic-dev-img"
        dev_image_tag="${version#v}" # Remove leading 'v' if present (e.g. v1.2.3 -> 1.2.3)
    fi

    cat <<EOF
{
    "name": "intrinsic-flowstate-devcontainer",
    "image": "${dev_image_base_name}:${dev_image_tag}",
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

  # Check if version is blank, "main", "latest", a valid semver tag, or the new candidate format
  if [[ -z "$version" || \
        "$version" == "main" || \
        "$version" == "latest" || \
        "$version" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+$ || \
        "$version" =~ ^candidate/intrinsic\.platform\.[0-9]{8}\.RC[0-9]+$ ]]; then
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
    generate_changelog "latest" # Default to "latest"
else
    if validate_version "$1"; then
        generate_changelog "$1"
    else
        echo "Error: Invalid version format."
        echo "Supported formats are: blank (defaults to 'latest'), 'main', 'latest',"
        echo "semantic versions like 'v1.2.3' or '1.2.3',"
        echo "or candidate tags like 'candidate/intrinsic.platform.YYYYMMDD.RCXX'."
        exit 1
    fi
fi
