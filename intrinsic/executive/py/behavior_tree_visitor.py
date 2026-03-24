# Copyright 2023 Intrinsic Innovation LLC

"""Visitor implementation for Behavior Trees in Python."""

from collections.abc import Callable

from intrinsic.executive.proto import behavior_tree_pb2

TreeVisitorCallback = Callable[[behavior_tree_pb2.BehaviorTree], None]
NodeVisitorCallback = Callable[
    [behavior_tree_pb2.BehaviorTree, behavior_tree_pb2.BehaviorTree.Node], None
]
ConditionVisitorCallback = Callable[
    [behavior_tree_pb2.BehaviorTree, behavior_tree_pb2.BehaviorTree.Condition],
    None,
]


def _walk_tree(
    tree: behavior_tree_pb2.BehaviorTree,
    tree_visitor: TreeVisitorCallback | None,
    node_visitor: NodeVisitorCallback | None,
    condition_visitor: ConditionVisitorCallback | None,
    visit_called_tree_state: bool,
) -> None:
  """Recursively walks a given tree and invokes visitor.

  Args:
    tree: tree to call visitor for and walk
    tree_visitor: optional callback to invoke for trees
    node_visitor: optional callback to invoke for nodes
    condition_visitor: optional callback to invoke for conditions
    visit_called_tree_state: whether to visit called_tree_state in TaskNodes
  """
  if tree_visitor is not None:
    tree_visitor(tree)

  _walk_node(
      tree,
      tree.root,
      tree_visitor,
      node_visitor,
      condition_visitor,
      visit_called_tree_state,
  )


def _walk_node(
    tree: behavior_tree_pb2.BehaviorTree,
    node: behavior_tree_pb2.BehaviorTree.Node,
    tree_visitor: TreeVisitorCallback | None,
    node_visitor: NodeVisitorCallback | None,
    condition_visitor: ConditionVisitorCallback | None,
    visit_called_tree_state: bool,
) -> None:
  """Recursively walks a given node and invokes visitor.

  Args:
    tree: tree the node belongs to
    node: node to call visitor for and walk
    tree_visitor: optional callback to invoke for trees
    node_visitor: optional callback to invoke for nodes
    condition_visitor: optional callback to invoke for conditions
    visit_called_tree_state: whether to visit called_tree_state in TaskNodes
  """
  if node_visitor is not None:
    node_visitor(tree, node)

  if node.HasField("decorators"):
    if node.decorators.HasField("condition"):
      _walk_condition(
          tree,
          node.decorators.condition,
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )

  if node.HasField("sequence"):
    for c in node.sequence.children:
      _walk_node(
          tree,
          c,
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )
  elif node.HasField("parallel"):
    for c in node.parallel.children:
      _walk_node(
          tree,
          c,
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )
  elif node.HasField("selector"):
    for c in node.selector.children:
      _walk_node(
          tree,
          c,
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )
    for c in node.selector.branches:
      if c.HasField("condition"):
        _walk_condition(
            tree,
            c.condition,
            tree_visitor,
            node_visitor,
            condition_visitor,
            visit_called_tree_state,
        )
      _walk_node(
          tree,
          c.node,
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )
  elif node.HasField("fallback"):
    for c in node.fallback.children:
      _walk_node(
          tree,
          c,
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )
    for c in node.fallback.tries:
      if c.HasField("condition"):
        _walk_condition(
            tree,
            c.condition,
            tree_visitor,
            node_visitor,
            condition_visitor,
            visit_called_tree_state,
        )
      _walk_node(
          tree,
          c.node,
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )
  elif node.HasField("branch"):
    if node.branch.HasField("if"):
      _walk_condition(
          tree,
          getattr(node.branch, "if"),
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )
    if node.branch.HasField("then"):
      _walk_node(
          tree,
          node.branch.then,
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )
    if node.branch.HasField("else"):
      _walk_node(
          tree,
          getattr(node.branch, "else"),
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )
  elif node.HasField("loop"):
    if node.loop.HasField("while"):
      _walk_condition(
          tree,
          getattr(node.loop, "while"),
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )
    if node.loop.HasField("do"):
      _walk_node(
          tree,
          node.loop.do,
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )
  elif node.HasField("retry"):
    if node.retry.HasField("child"):
      _walk_node(
          tree,
          node.retry.child,
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )
    if node.retry.HasField("recovery"):
      _walk_node(
          tree,
          node.retry.recovery,
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )
  elif node.HasField("sub_tree"):
    if node.sub_tree.HasField("tree"):
      _walk_tree(
          node.sub_tree.tree,
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )
  elif node.HasField("task"):
    if visit_called_tree_state and node.task.HasField("called_tree_state"):
      _walk_tree(
          node.task.called_tree_state,
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )


def _walk_condition(
    tree: behavior_tree_pb2.BehaviorTree,
    cond: behavior_tree_pb2.BehaviorTree.Condition,
    tree_visitor: TreeVisitorCallback | None,
    node_visitor: NodeVisitorCallback | None,
    condition_visitor: ConditionVisitorCallback | None,
    visit_called_tree_state: bool,
) -> None:
  """Recursively walks a given condition and invokes visitor.

  Args:
    tree: tree the condition belongs to
    cond: condition to call visitor for and walk
    tree_visitor: optional callback to invoke for trees
    node_visitor: optional callback to invoke for nodes
    condition_visitor: optional callback to invoke for conditions
    visit_called_tree_state: whether to visit called_tree_state in TaskNodes
  """
  if condition_visitor is not None:
    condition_visitor(tree, cond)

  if cond.HasField("behavior_tree"):
    _walk_tree(
        cond.behavior_tree,
        tree_visitor,
        node_visitor,
        condition_visitor,
        visit_called_tree_state,
    )
  elif cond.HasField("all_of"):
    for c in cond.all_of.conditions:
      _walk_condition(
          tree,
          c,
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )
  elif cond.HasField("any_of"):
    for c in cond.any_of.conditions:
      _walk_condition(
          tree,
          c,
          tree_visitor,
          node_visitor,
          condition_visitor,
          visit_called_tree_state,
      )
  elif cond.HasField("not"):
    _walk_condition(
        tree,
        getattr(cond, "not"),
        tree_visitor,
        node_visitor,
        condition_visitor,
        visit_called_tree_state,
    )


def walk(
    tree: behavior_tree_pb2.BehaviorTree,
    *,
    tree_visitor: TreeVisitorCallback | None = None,
    node_visitor: NodeVisitorCallback | None = None,
    condition_visitor: ConditionVisitorCallback | None = None,
    visit_called_tree_state=False,
) -> None:
  """Recursively walks a given tree and invokes visitors.

  Args:
    tree: tree to walk
    tree_visitor: optional callback to invoke for trees
    node_visitor: optional callback to invoke for nodes
    condition_visitor: optional callback to invoke for conditions
    visit_called_tree_state: whether to visit called_tree_state in TaskNodes
  """
  _walk_tree(
      tree,
      tree_visitor,
      node_visitor,
      condition_visitor,
      visit_called_tree_state,
  )
