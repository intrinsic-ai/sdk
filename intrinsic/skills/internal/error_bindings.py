# Copyright 2023 Intrinsic Innovation LLC

"""Utility function that raises a pybind11_abseil StatusNotOk execption."""

from pybind11_abseil import status


def raise_status(code: status.StatusCode, text: str) -> None:
  raise status.BuildStatusNotOk(code, text)
SKILL_SERVICE_COMPONENT = 'ai.intrinsic.skill'
SKILL_SERVICE_WAIT_TIMEOUT_CODE = 11010
