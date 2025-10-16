# Copyright 2023 Intrinsic Innovation LLC

"""Provides behavior trees from a solution."""

from typing import Iterator
import grpc
from intrinsic.executive.proto import behavior_tree_pb2
from intrinsic.frontend.solution_service.proto import solution_service_pb2
from intrinsic.frontend.solution_service.proto import solution_service_pb2_grpc
from intrinsic.solutions import behavior_tree
from intrinsic.solutions import providers


_DEFAULT_PAGE_SIZE = 50


class Processes(providers.ProcessProvider):
  """Provides the processes (= behavior trees) from a solution."""

  _solution: solution_service_pb2_grpc.SolutionServiceStub

  def __init__(self, solution: solution_service_pb2_grpc.SolutionServiceStub):
    self._solution = solution

  def _list_all_behavior_trees(self) -> list[behavior_tree_pb2.BehaviorTree]:
    response = self._solution.ListBehaviorTrees(
        solution_service_pb2.ListBehaviorTreesRequest(
            page_size=_DEFAULT_PAGE_SIZE
        )
    )
    bts = []
    bts.extend(response.behavior_trees)
    while response.next_page_token:
      response = self._solution.ListBehaviorTrees(
          solution_service_pb2.ListBehaviorTreesRequest(
              page_size=_DEFAULT_PAGE_SIZE,
              page_token=response.next_page_token,
          )
      )
      bts.extend(response.behavior_trees)
    return bts

  def keys(self) -> list[str]:
    return [bt.name for bt in self._list_all_behavior_trees()]

  def __iter__(self) -> Iterator[str]:
    for bt in self._list_all_behavior_trees():
      yield bt.name

  def items(self) -> Iterator[tuple[str, behavior_tree.BehaviorTree]]:
    for bt in self._list_all_behavior_trees():
      yield bt.name, behavior_tree.BehaviorTree(bt=bt)

  def values(self) -> Iterator[behavior_tree.BehaviorTree]:
    for bt in self._list_all_behavior_trees():
      yield behavior_tree.BehaviorTree(bt=bt)

  def __contains__(self, identifier: str) -> bool:
    try:
      self._solution.GetBehaviorTree(
          solution_service_pb2.GetBehaviorTreeRequest(name=identifier)
      )
      return True
    except grpc.RpcError as e:
      if hasattr(e, 'code') and e.code() == grpc.StatusCode.NOT_FOUND:
        return False
      raise e

  def __getitem__(self, identifier: str) -> behavior_tree.BehaviorTree:
    try:
      return behavior_tree.BehaviorTree(
          bt=self._solution.GetBehaviorTree(
              solution_service_pb2.GetBehaviorTreeRequest(name=identifier)
          )
      )
    except Exception as e:
      raise KeyError(f'Process "{identifier}" not available') from e

  def __setitem__(self, identifier: str, value: behavior_tree.BehaviorTree):
    value.name = identifier
    self._solution.UpdateBehaviorTree(
        solution_service_pb2.UpdateBehaviorTreeRequest(
            behavior_tree=value.proto,
            allow_missing=True,
        )
    )

  def __delitem__(self, name: str):
    try:
      self._solution.DeleteBehaviorTree(
          solution_service_pb2.DeleteBehaviorTreeRequest(name=name)
      )
    except Exception as e:
      raise KeyError(f"Failed to delete behavior tree '{name}'") from e
