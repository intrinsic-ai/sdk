# Copyright 2023 Intrinsic Innovation LLC

"""Internal module for world queries.

Do not import this module directly - use the forwarded names defined in
intrinsic.solutions.behavior_tree.
"""

from __future__ import annotations

import enum
from typing import Any as AnyType
from typing import List
from typing import Optional
from typing import Union

from intrinsic.executive.proto import any_with_assignments_pb2
from intrinsic.executive.proto import world_query_pb2
from intrinsic.solutions import blackboard_value
from intrinsic.solutions import errors as solutions_errors
from intrinsic.solutions import utils
from intrinsic.world.proto import object_world_refs_pb2
from intrinsic.world.python import object_world_resources

WorldQueryObject = Union[
    object_world_resources.WorldObject,
    object_world_refs_pb2.ObjectReference,
    blackboard_value.BlackboardValue,
]


class WorldQuery:
  """Wrapper for WorldQuery proto for easier construction and conversion."""

  _proto: world_query_pb2.WorldQuery
  _assignments: List[any_with_assignments_pb2.AnyWithAssignments.Assignment]

  def __init__(self, proto: Optional[world_query_pb2.WorldQuery] = None):
    self._proto = world_query_pb2.WorldQuery()
    if proto is not None:
      self._proto.CopyFrom(proto)
    self._assignments = []

  def _object_to_reference(
      self, obj: Optional[WorldQueryObject]
  ) -> Optional[object_world_refs_pb2.ObjectReference]:
    """Converts an object to a reference (or returns as-is if reference given).

    Args:
      obj: object to convert

    Returns:
      An object reference, either retrieved from the WorldObject or just the
      reference that was passed in.

    Raises:
      TypeError: if the passed in object is neither a WorldObject nor an
        ObjectReference.
    """
    if obj is None:
      return None

    if isinstance(obj, blackboard_value.BlackboardValue):
      return object_world_refs_pb2.ObjectReference()

    if isinstance(obj, object_world_resources.WorldObject):
      return obj.reference

    if isinstance(obj, object_world_refs_pb2.ObjectReference):
      return obj

    raise TypeError(
        'Invalid type for object, cannot convert to ObjectReference'
    )

  @utils.protoenum(
      proto_enum_type=world_query_pb2.WorldQuery.Order.Criterion,
      unspecified_proto_enum_map_to_none=(
          world_query_pb2.WorldQuery.Order.Criterion.SORT_ORDER_UNSPECIFIED
      ),
  )
  class OrderCriterion(enum.Enum):
    """Specifies what to sort returned values by."""

  @utils.protoenum(proto_enum_type=world_query_pb2.WorldQuery.Order.Direction)
  class OrderDirection(enum.Enum):
    """Specifies sort order for returned items."""

  def _handle_blackboard_assignments(
      self, path: str, assigned: AnyType
  ) -> None:
    """Handle assignments that might come from the blackboard.

    If the assigned object is a BlackboardValue an assignment to path will be
    added. Otherwise nothing happens.

    Args:
      path: The field in the WorldQuery that would be set.
      assigned: Either a BlackboardValue or something else.
    """
    if isinstance(assigned, blackboard_value.BlackboardValue):
      self._assignments.append(
          any_with_assignments_pb2.AnyWithAssignments.Assignment(
              path=path,
              cel_expression=assigned.value_access_path(),
          )
      )

  def select(
      self,
      *,
      child_frames_of: Optional[WorldQueryObject] = None,
      child_objects_of: Optional[WorldQueryObject] = None,
      children_of: Optional[WorldQueryObject] = None,
  ) -> WorldQuery:
    """Sets the query of the world query.

    Set only one of the possible arguments.

    Args:
      child_frames_of: the object of which to retrieve child frames of
      child_objects_of: the object of which to retrieve child objects of
      children_of: the object of which to retrieve children of

    Returns:
      Self (for chaining in a builder pattern)

    Raises:
      InvalidArgumentError: if zero or more than 1 input argument is set
    """
    num_inputs = 0

    if child_frames_of is not None:
      num_inputs += 1
      self._proto.select.child_frames_of.CopyFrom(
          self._object_to_reference(child_frames_of)
      )
      self._handle_blackboard_assignments(
          'select.child_frames_of', child_frames_of
      )
    if child_objects_of is not None:
      num_inputs += 1
      self._proto.select.child_objects_of.CopyFrom(
          self._object_to_reference(child_objects_of)
      )
      self._handle_blackboard_assignments(
          'select.child_objects_of', child_objects_of
      )
    if children_of is not None:
      num_inputs += 1
      self._proto.select.children_of.CopyFrom(
          self._object_to_reference(children_of)
      )
      self._handle_blackboard_assignments('select.children_of', children_of)

    if num_inputs != 1:
      raise solutions_errors.InvalidArgumentError(
          'Data node for create or update requires exactly 1 input'
          f' element, got {num_inputs}'
      )

    return self

  def filter(
      self, *, name_regex: Union[str, blackboard_value.BlackboardValue]
  ) -> WorldQuery:
    """Sets the filter of the world query.

    Args:
      name_regex: RE2 regular expression that names must fully match to be
        returned.

    Returns:
      Self (for chaining in a builder pattern)
    """
    if isinstance(name_regex, blackboard_value.BlackboardValue):
      self._handle_blackboard_assignments('filter.name_regex', name_regex)
    else:
      self._proto.filter.name_regex = name_regex
    return self

  def order(
      self,
      *,
      by: OrderCriterion,
      direction: OrderDirection = OrderDirection.ASCENDING,
  ) -> WorldQuery:
    """Sets the ordering of the world query.

    Args:
      by: criterion identifying what to order by
      direction: ordering direction, ascending or descending

    Returns:
      Self (for chaining in a builder pattern)
    """
    self._proto.order.by = by.value
    self._proto.order.direction = direction.value
    return self

  @property
  def proto(self) -> world_query_pb2.WorldQuery:
    return self._proto

  @property
  def assignments(
      self,
  ) -> List[any_with_assignments_pb2.AnyWithAssignments.Assignment]:
    return self._assignments

  @classmethod
  def create_from_proto(
      cls, proto_object: world_query_pb2.WorldQuery
  ) -> WorldQuery:
    return cls(proto_object)

  def __str__(self) -> str:
    return (
        f'WorldQuery(text_format.Parse("""{self._proto}""",'
        ' intrinsic_proto.executive.WorldQuery()))'
    )
