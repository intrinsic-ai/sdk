# Copyright 2023 Intrinsic Innovation LLC

"""Lightweight Python wrappers around actions."""

import abc
import datetime
from typing import Optional

from intrinsic.executive.proto import behavior_call_pb2


class ActionBase(abc.ABC):
  """Abstract base class of an action.

  Derived classes need to override the getter for self.proto.
  """

  def __init__(self):
    self._project_timeout: Optional[datetime.timedelta] = None
    self._execute_timeout: Optional[datetime.timedelta] = None

  @property
  @abc.abstractmethod
  def proto(self) -> behavior_call_pb2.BehaviorCall:
    """Proto representation of action.

    Needs to be overridden by subclasses.

    Returns:
      Proto representation of action as behavior_call_pb2.BehaviorCall.

    Raises:
      NoImplementedError if the class fails to override method.
    """
    raise NotImplementedError

  @property
  def execute_timeout(self) -> Optional[datetime.timedelta]:
    """Timeout after which execution should be considered failed."""
    return self._execute_timeout

  @execute_timeout.setter
  def execute_timeout(self, timeout: datetime.timedelta) -> None:
    self._execute_timeout = timeout

  @property
  def project_timeout(self) -> Optional[datetime.timedelta]:
    """Timeout after which projection should be considered failed."""
    return self._project_timeout

  @project_timeout.setter
  def project_timeout(self, timeout: datetime.timedelta) -> None:
    self._project_timeout = timeout
