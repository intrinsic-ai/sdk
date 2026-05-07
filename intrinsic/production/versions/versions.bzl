# Copyright 2023 Intrinsic Innovation LLC

"""Constants for Intrinsic OS versions."""

# NOTE: careful! the OS versions here are *without* the prefix `xfa.` which
# they do use in versions.go. Make sure not to add `xfa.` here, as it will
# break the fleet-manager.
PREVIOUS_OS_VERSION = "20260326.RC01"

# Version that is currently running with the intrinsic stack
STABLE_OS_VERSION = "20260416.RC02"

CANARY_OS_VERSION = "20260507.RC00"
