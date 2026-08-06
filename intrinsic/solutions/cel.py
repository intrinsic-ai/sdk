# Copyright 2023 Intrinsic Innovation LLC

"""Helpers for the Common Expression Language (CEL)."""

from intrinsic.solutions.blackboard_value import BlackboardValue


class CelExpression:
  """A CEL expression.

  Represents an expression to be evaluated in lieu of a skill parameter value or
  a BlackboardValue.
  """

  _expression: str

  def __init__(self, expression: str | BlackboardValue):
    self._expression = str(expression)

  def __str__(self) -> str:
    return self._expression
