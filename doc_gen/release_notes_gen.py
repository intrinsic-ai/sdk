import argparse
import git
from git.repo import Repo
from git.exc import InvalidGitRepositoryError
import os
import re
import sys

def validate_version(version):
    """Validates the version string."""
    if not version or version == "main" or version == "latest" or re.match(r"^[v]?\d+\.\d+\.\d+$", version):
        return True
    return False

def generate_bazel_dep(version="main"):
    """Generates the Bazel dependency and archive override."""
    if version == "main" or version == "latest":
        url = "https://github.com/intrinsic-ai/sdk/archive/refs/heads/main.tar.gz"
        strip_prefix = "sdk-main/"
    elif version.startswith("candidate/"):
        strip_prefix = f"sdk-{version.replace('/', '-')}"
        url = f"https://github.com/intrinsic-ai/sdk/archive/refs/tags/{version}.tar.gz"
    else:
        if not version.startswith("v"):
            version = f"v{version}"
        strip_prefix = f"sdk-{version[1:]}/"
        url = f"https://github.com/intrinsic-ai/sdk/archive/refs/tags/{version}.tar.gz"

    bazel_code = f"""bazel_dep(name = "ai_intrinsic_sdks")
archive_override(
    module_name = "ai_intrinsic_sdks",
    urls = ["{url}"],
    strip_prefix = "{strip_prefix}"
)"""
    return bazel_code

def generate_devcontainer_config(version="latest"):
    """Generates the devcontainer configuration."""
    version = version.lstrip("v")
    devcontainer_config = f"""{{
    "name": "intrinsic-flowstate-devcontainer",
    "image": "ghcr.io/intrinsic-ai/intrinsic-dev-img:{version}",
    "runArgs": [
        "--network=host"
    ],
    "customizations": {{
        "vscode": {{
            "settings": {{
                "intrinsic.defaultSdkRepository": "https://github.com/intrinsic-ai/sdk.git"
            }}
        }}
    }}
}}"""
    return devcontainer_config

def generate_release_notes(repo_path, last_release_tag, version="latest"):
    """
    Reads commit messages from a Git repository after a specific release tag,
    skipping commits with messages starting with "SDK Update", and generates
    a markdown document with bullet points, joining subsequent non-empty lines with ". ",
    and removing the last line if it starts with "GitOrigin-RevId:",
    and prepends Bazel and devcontainer configuration.

    Args:
        repo_path (str): The path to the Git repository.
        last_release_tag (str): The name of the last release tag.
        version (str, optional): The release version for Bazel and devcontainer. Defaults to "latest".
        output_file (str, optional): The name of the output markdown file.
            Defaults to "release_notes.md".
    """
    if not validate_version(version):
        print("Error: Invalid version format. Must be blank, 'main', 'latest', or a tag like v1.2.3 or 1.2.3")
        sys.exit(1)

    bazel_code = generate_bazel_dep(version)
    bazel_info = "## bazelmod configuration\n\n"
    bazel_info += "Update your MODULE.bazel file to use the newest release archive:\n\n"
    bazel_info += f"""```\n{bazel_code}\n````\n\n"""

    devcontainer_config = generate_devcontainer_config(version)
    devcontainer_info = "## devcontainer configuration\n\n"
    devcontainer_info += "Update your devcontainer configuration to use the latest base image:\n\n"
    devcontainer_info += f"""```json\n{devcontainer_config}\n````\n\n"""

    commit_messages = ""
    try:
        # Load the repository
        repo = Repo(repo_path)

        # Get the last release tag
        try:
            tag = repo.tags[last_release_tag]
        except IndexError:
            print(f"Error: Tag '{last_release_tag}' not found in the repository.")
            return

        # Get all commits after the specified tag
        commits = list(repo.iter_commits(f'{tag.commit}..HEAD'))
        commits.reverse()  # Order commits from oldest to newest

        if not commits:
            print(f"No commits found after tag '{last_release_tag}'.")
            return bazel_info + devcontainer_info

        # Prepare the markdown content for commit messages
        commit_messages += "## Commit History\n\n"
        for commit in commits:
            message_lines = commit.message.strip().split('\n')
            if message_lines and not message_lines[0].startswith("SDK update"):
                processed_lines = []
                for line in message_lines:
                    stripped_line = line.strip()
                    if stripped_line:  # If the line is not empty after stripping
                        if processed_lines:
                            processed_lines.append(". " + stripped_line)
                        else:
                            processed_lines.append(stripped_line)

                if processed_lines and len(processed_lines) > 1 and processed_lines[-1].startswith("GitOrigin-RevId:"):
                    processed_lines = processed_lines[:-1]

                single_line_message = "".join(processed_lines)
                commit_messages += f"* {single_line_message}\n"

    except InvalidGitRepositoryError:
        print(f"Error: '{repo_path}' is not a valid Git repository.")
    except Exception as e:
        print(f"An error occurred: {e}")

    # Combine all parts
    return bazel_info + devcontainer_info + commit_messages

def write_output_file (file_path, file_output):
    # Write the markdown content to a file
    with open(file_path, 'w') as f:
        f.write(file_output)

    print(f"Release notes generated successfully in '{file_path}'.")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Generate release notes for a Git repository, including Bazel and Devcontainer configurations.")

    parser.add_argument(
        "--repo-path",
        type=str,
        required=True,
        help="The path to the Git repository (e.g., /path/to/your/repo or C:\\path\\to\\your\\repo)."
    )
    parser.add_argument(
        "--last-release-tag",
        type=str,
        required=True,
        help="The name of the last release tag (e.g., v1.0.0)."
    )
    parser.add_argument(
        "--version",
        type=str,
        default="latest",
        help="The release version for Bazel and Devcontainer configurations (e.g., v1.1.0, main, latest). Defaults to 'latest'."
    )
    parser.add_argument(
        "--output-file",
        type=str,
        default="RELEASE_NOTES.md",
        help="The desired output filename for the release notes. Defaults to 'release_notes.md'."
    )

    args = parser.parse_args()

    final_output = generate_release_notes(
        repo_path=args.repo_path,
        last_release_tag=args.last_release_tag,
        version=args.version
    )

    write_output_file(args.output_file, final_output)
