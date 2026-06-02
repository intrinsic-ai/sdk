# Copyright 2023 Intrinsic Innovation LLC

"""
Exec Properties used for executors.
"""

DEFAULT = {
    "container-image": "docker://us-central1-docker.pkg.dev/intrinsic-mirror/intrinsic-build-images/bazel-rbe-executor@sha256:c2d50e5f5a3bbea4c47ffcdc1b2755168d1c682ec2b3ba620fb3134eca4bab0d",
    "OSFamily": "Linux",
}

REQUIRES_NETWORK = {
    "dockerNetwork": "standard",
}
