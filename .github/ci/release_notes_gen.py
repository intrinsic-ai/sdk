import argparse
import git
from git.exc import InvalidGitRepositoryError
import re
import sys

from typing import Optional

def generate_bazel_dep(version="main") -> str:
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

    bazel_code = f"""## bazelmod configuration

Update your MODULE.bazel file to use the newest release archive:

```
bazel_dep(name = "ai_intrinsic_sdks")
archive_override(
    module_name = "ai_intrinsic_sdks",
    urls = ["{url}"],
    strip_prefix = "{strip_prefix}"
)
```"""
    return bazel_code

def generate_devcontainer_config(version="latest") -> str:
    """Generates the devcontainer configuration."""
    version = version.lstrip("v")
    devcontainer_config = f"""## devcontainer configuration

Update your devcontainer configuration to use the latest base image:

```json
{{
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
}}
```"""
    return devcontainer_config

def generate_changelog(repo: git.Repo, last_release_tag: git.TagReference) -> str:
    release_notes = ""
    commits = list(repo.iter_commits(f'{last_release_tag.commit}..HEAD'))
    commits.reverse()  # Order commits from oldest to newest
    print(f"Found {len(commits)} commits")

    if commits:
        messages = []
        for commit in commits:
            message_lines = commit.message.strip().split('\n')
            if message_lines and not message_lines[0].startswith("SDK update"):
                messages.append(f'* {message_lines[0]}')
        release_notes = '\n'.join(messages)
    return release_notes


def find_most_recent_release(repo: git.Repo) -> Optional[git.TagReference]:
    """
    Iterate through history to find the most recent release tag
    """

    commit_to_tags = {}
    for tag in repo.tags:
        if tag.commit in commit_to_tags:
            commit_to_tags[tag.commit].append(tag)
        else:
            commit_to_tags[tag.commit] = [tag]

    # Regex to match vX.Y.Z format
    tag_pattern = re.compile(r'^v\d+\.\d+\.\d+$')

    # Iterate backwards from HEAD
    for commit in repo.iter_commits('HEAD'):
        if commit in commit_to_tags:
            for tag in commit_to_tags[commit]:
                if tag_pattern.match(tag.name):
                    return tag

    return None


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
        help="The name of the last release tag (e.g., v1.0.0)."
    )
    parser.add_argument(
        "--version",
        type=str,
        help="The release version for Bazel and Devcontainer configurations (e.g., v1.1.0, main, latest). Defaults to 'latest'."
    )
    parser.add_argument(
        "--output-file",
        type=str,
        default="RELEASE_NOTES.md",
        help="The desired output filename for the release notes. Defaults to 'release_notes.md'."
    )

    args = parser.parse_args()

    try:
        repo = git.Repo(args.repo_path)
    except git.InvalidGitRepositoryError:
        print(f"Error: '{args.repo_path}' is not a valid Git repository.")
        sys.exit(-1)
    except git.NoSuchPathError:
        print(f"Error: The path '{args.repo_path}' does not exist.")
        sys.exit(-1)

    if args.last_release_tag:
        try:
            last_release = repo.tags[args.last_release_tag]
        except IndexError:
            print(f"Error: Could not find '{args.last_release_tag}' in '{args.repo_path}'")
            sys.exit(-1)
    else:
        last_release = find_most_recent_release(repo)

    print(f"Generating changelog from last release: {last_release}")

    changelog = generate_changelog(repo, last_release)
    bazel_config = generate_bazel_dep(args.version)
    devcontainer_config = generate_devcontainer_config(args.version)

    output = f"""## Changes since {last_release}

{changelog}

{bazel_config}

{devcontainer_config}"""

    with open(args.output_file, 'w') as f:
        f.write(output)

    repo.close()
    sys.exit(0)

