# Copyright 2023 Intrinsic Innovation LLC

"""Tests behavior_tree_visitor module."""

from absl.testing import absltest
from absl.testing import parameterized
from google.protobuf import text_format

from intrinsic.executive.proto import behavior_tree_pb2
from intrinsic.executive.py import behavior_tree_visitor


class BehaviorTreeVisitorTest(parameterized.TestCase):

  def test_walk(self):
    behavior_tree = behavior_tree_pb2.BehaviorTree()
    text_format.Parse(
        r"""
      tree_id: "T1"
      root {
        id: 1
        sequence {
          children: [
            {
              id: 11,
              decorators {
                condition {
                  all_of {
                    conditions: [
                      { blackboard { cel_expression: "c1" } },
                      { blackboard { cel_expression: "c2" } }
                    ]
                  }
                }
              }
              parallel {
                children: [
                  {
                    id: 111
                    fail {}
                  },
                  {
                    id: 112
                    fail {}
                  }
                ]
              }
            },
            {
              decorators {
                condition {
                  any_of {
                    conditions: [
                      { blackboard { cel_expression: "c3" } },
                      { blackboard { cel_expression: "c4" } }
                    ]
                  }
                }
              }
              id: 2
              fail {}
            },
            {
              id: 3
              sub_tree {
                tree {
                  tree_id: "T2"
                  root {
                    id: 31
                    decorators {
                      condition {
                        not { blackboard { cel_expression: "c5"} }
                      }
                    }
                    selector {
                      children: [
                        {
                          id: 32
                          fail {}
                        },
                        {
                          id:33
                          fail {}
                        }
                      ]
                    }
                  }
                }
              }
            },
            {
              id: 4
              fallback {
                children: [
                  {
                    id: 41
                    fail {}
                  },
                  {
                    id: 42
                    fail {}
                  }
                ]
              }
            },
            {
              id: 5
              branch {
                if {
                  behavior_tree {
                    tree_id: "T3"
                    root {
                      id: 53
                      decorators {
                        condition {
                          blackboard { cel_expression: "c6" }
                        }
                      }
                      sequence {}
                    }
                  }
                }
                then {
                  id: 51
                  fail {}
                }
                else {
                  id: 52
                  fail {}
                }
              }
            },
            {
              id: 6
              loop {
                do {
                  id: 61
                  fail {}
                }
              }
            },
            {
              id: 7
              retry {
                child {
                  id: 71
                  fail {}
                }
                recovery {
                  id: 72
                  fail {}
                }
              }
            },
            {
              id: 8
              selector {
                branches {
                  node {
                    id: 81
                    fail {}
                  }
                }
                branches {
                  condition { any_of{} } 
                  node {
                    id: 82
                    fail {}
                  }
                }
                branches {
                  condition {
                    behavior_tree {
                      tree_id: "T4"
                      root {
                        id: 801
                        fail {}
                      }
                    }
                  } 
                  node{
                    id: 83
                    fail {}
                  }
                }
              }
            },
            {
              id: 9
              fallback {
                tries {
                  node {
                    id: 91
                    fail {}
                  }
                }
                tries {
                  condition { any_of{} } 
                  node {
                    id: 92
                    fail {}
                  }
                }
                tries {
                  condition {
                    behavior_tree {
                      tree_id: "T5"
                      root {
                        id: 901
                        fail {}
                      }
                    }
                  } 
                  node{
                    id: 93
                    fail {}
                  }
                }
              }
            }
          ]
        }
      }
    """,
        behavior_tree,
    )
    tree_ids = []
    node_ids = []
    conds = []

    def visit_tree(tree: behavior_tree_pb2.BehaviorTree):
      tree_ids.append(tree.tree_id)

    def visit_node(
        tree: behavior_tree_pb2.BehaviorTree,
        node: behavior_tree_pb2.BehaviorTree.Node,
    ):
      node_ids.append((tree.tree_id, node.id))

    def visit_cond(
        tree: behavior_tree_pb2.BehaviorTree,
        cond: behavior_tree_pb2.BehaviorTree.Condition,
    ):
      conds.append((tree.tree_id, text_format.MessageToString(cond)))

    behavior_tree_visitor.walk(
        behavior_tree,
        tree_visitor=visit_tree,
        node_visitor=visit_node,
        condition_visitor=visit_cond,
    )

    self.assertEqual(tree_ids, ["T1", "T2", "T3", "T4", "T5"])
    self.assertEqual(
        node_ids,
        [("T1", n) for n in [1, 11, 111, 112, 2, 3]]
        + [("T2", n) for n in [31, 32, 33]]
        + [("T1", n) for n in [4, 41, 42, 5]]
        + [("T3", n) for n in [53]]
        + [("T1", n) for n in [51, 52, 6, 61, 7, 71, 72, 8, 81, 82]]
        + [("T4", n) for n in [801]]
        + [("T1", n) for n in [83, 9, 91, 92]]
        + [("T5", n) for n in [901]]
        + [("T1", n) for n in [93]],
    )
    self.assertEqual(
        conds,
        [
            (
                "T1",
                """all_of {
  conditions {
    blackboard {
      cel_expression: "c1"
    }
  }
  conditions {
    blackboard {
      cel_expression: "c2"
    }
  }
}
""",
            ),
            ("T1", 'blackboard {\n  cel_expression: "c1"\n}\n'),
            ("T1", 'blackboard {\n  cel_expression: "c2"\n}\n'),
            (
                "T1",
                """any_of {
  conditions {
    blackboard {
      cel_expression: "c3"
    }
  }
  conditions {
    blackboard {
      cel_expression: "c4"
    }
  }
}
""",
            ),
            ("T1", 'blackboard {\n  cel_expression: "c3"\n}\n'),
            ("T1", 'blackboard {\n  cel_expression: "c4"\n}\n'),
            ("T2", 'not {\n  blackboard {\n    cel_expression: "c5"\n  }\n}\n'),
            ("T2", 'blackboard {\n  cel_expression: "c5"\n}\n'),
            (
                "T1",
                """behavior_tree {
  root {
    sequence {
    }
    decorators {
      condition {
        blackboard {
          cel_expression: "c6"
        }
      }
    }
    id: 53
  }
  tree_id: "T3"
}
""",
            ),
            ("T3", 'blackboard {\n  cel_expression: "c6"\n}\n'),
            ("T1", "any_of {\n}\n"),
            (
                "T1",
                """behavior_tree {
  root {
    fail {
    }
    id: 801
  }
  tree_id: "T4"
}
""",
            ),
            ("T1", "any_of {\n}\n"),
            (
                "T1",
                """behavior_tree {
  root {
    fail {
    }
    id: 901
  }
  tree_id: "T5"
}
""",
            ),
        ],
    )


if __name__ == "__main__":
  absltest.main()
