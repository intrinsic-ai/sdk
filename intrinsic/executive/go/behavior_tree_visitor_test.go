// Copyright 2023 Intrinsic Innovation LLC

package behaviortree_test

import (
	"context"
	"slices"
	"strings"
	"testing"

	"intrinsic/executive/go/behaviortree"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"

	btpb "intrinsic/executive/proto/behavior_tree_go_proto"
)

func conditionName(cond *btpb.BehaviorTree_Condition) string {
	switch cond.GetConditionType().(type) {
	case *btpb.BehaviorTree_Condition_Blackboard:
		return "Condition_Blackboard"
	case *btpb.BehaviorTree_Condition_BehaviorTree:
		return "Condition_BehaviorTree"
	case *btpb.BehaviorTree_Condition_AllOf:
		return "Condition_AllOf"
	case *btpb.BehaviorTree_Condition_AnyOf:
		return "Condition_AnyOf"
	case *btpb.BehaviorTree_Condition_Not:
		return "Condition_Not"
	case *btpb.BehaviorTree_Condition_StatusMatch:
		return "Condition_StatusMatch"
	default:
		return "Unknown"
	}
}

func elementName(elem behaviortree.VisitElement) string {
	if tree := elem.Tree(); tree != nil {
		return tree.GetTreeId()
	}
	if node := elem.Node(); node != nil {
		return node.GetName()
	}
	if cond := elem.Condition(); cond != nil {
		return conditionName(cond)
	}
	panic("invalid visit element")
}

type nodeNameCollector struct {
	names      []string
	stopOnName string
}

func (c *nodeNameCollector) Visit(ctx context.Context, element behaviortree.VisitElement) error {
	namePath := []string{
		elementName(element),
	}
	for ancestor := range element.Ancestors() {
		namePath = append(namePath, elementName(ancestor))
	}
	// Reverse from element>root to root>element.
	slices.Reverse(namePath)
	c.names = append(c.names, strings.Join(namePath, "/"))
	if c.stopOnName != "" && elementName(element) == c.stopOnName {
		return behaviortree.Stop
	}
	return nil
}

func (c *nodeNameCollector) VisitCondition(ctx context.Context, cond *btpb.BehaviorTree_Condition) error {
	return nil
}

func TestNodes(t *testing.T) {
	tree := &btpb.BehaviorTree{
		TreeId: proto.String("Tree_A"),
		Root: &btpb.BehaviorTree_Node{
			Name: proto.String("Sequence_1"),
			NodeType: &btpb.BehaviorTree_Node_Sequence{
				Sequence: &btpb.BehaviorTree_SequenceNode{
					Children: []*btpb.BehaviorTree_Node{
						{Name: proto.String("Sequence_1")},
						{
							Name: proto.String("Parallel_1"),
							NodeType: &btpb.BehaviorTree_Node_Parallel{
								Parallel: &btpb.BehaviorTree_ParallelNode{
									Children: []*btpb.BehaviorTree_Node{
										{Name: proto.String("Child_1")},
										{Name: proto.String("Child_2")},
									},
								},
							},
						},
						{
							Name: proto.String("Fallback_1"),
							NodeType: &btpb.BehaviorTree_Node_Fallback{
								Fallback: &btpb.BehaviorTree_FallbackNode{
									Children: []*btpb.BehaviorTree_Node{
										{Name: proto.String("Child_1")},
										{Name: proto.String("Child_2")},
									},
								},
							},
						},
						{
							Name: proto.String("Fallback_2"),
							NodeType: &btpb.BehaviorTree_Node_Fallback{
								Fallback: &btpb.BehaviorTree_FallbackNode{
									Tries: []*btpb.BehaviorTree_FallbackNode_Try{
										{Node: &btpb.BehaviorTree_Node{Name: proto.String("Try_1")}},
										{
											Condition: &btpb.BehaviorTree_Condition{
												ConditionType: &btpb.BehaviorTree_Condition_AnyOf{
													AnyOf: &btpb.BehaviorTree_Condition_LogicalCompound{},
												},
											},
											Node: &btpb.BehaviorTree_Node{Name: proto.String("Try_2")},
										},
										{
											Condition: &btpb.BehaviorTree_Condition{
												ConditionType: &btpb.BehaviorTree_Condition_BehaviorTree{
													BehaviorTree: &btpb.BehaviorTree{
														TreeId: proto.String("Tree_B"),
														Root: &btpb.BehaviorTree_Node{
															Name: proto.String("Sequence_2"),
															NodeType: &btpb.BehaviorTree_Node_Sequence{
																Sequence: &btpb.BehaviorTree_SequenceNode{},
															},
														},
													},
												},
											},
											Node: &btpb.BehaviorTree_Node{Name: proto.String("Try_3")},
										},
									},
								},
							},
						},
						{
							Name: proto.String("Selector_1"),
							NodeType: &btpb.BehaviorTree_Node_Selector{
								Selector: &btpb.BehaviorTree_SelectorNode{
									Children: []*btpb.BehaviorTree_Node{
										{Name: proto.String("Child_1")},
										{Name: proto.String("Child_2")},
									},
								},
							},
						},
						{
							Name: proto.String("Selector_2"),
							NodeType: &btpb.BehaviorTree_Node_Selector{
								Selector: &btpb.BehaviorTree_SelectorNode{
									Branches: []*btpb.BehaviorTree_SelectorNode_Branch{
										{Node: &btpb.BehaviorTree_Node{Name: proto.String("Branch_1")}},
										{
											Condition: &btpb.BehaviorTree_Condition{
												ConditionType: &btpb.BehaviorTree_Condition_AnyOf{
													AnyOf: &btpb.BehaviorTree_Condition_LogicalCompound{},
												},
											},
											Node: &btpb.BehaviorTree_Node{Name: proto.String("Branch_2")},
										},
										{
											Condition: &btpb.BehaviorTree_Condition{
												ConditionType: &btpb.BehaviorTree_Condition_BehaviorTree{
													BehaviorTree: &btpb.BehaviorTree{
														TreeId: proto.String("Tree_C"),
														Root: &btpb.BehaviorTree_Node{
															Name: proto.String("Sequence_3"),
															NodeType: &btpb.BehaviorTree_Node_Sequence{
																Sequence: &btpb.BehaviorTree_SequenceNode{},
															},
														},
													},
												},
											},
											Node: &btpb.BehaviorTree_Node{Name: proto.String("Branch_3")},
										},
									},
								},
							},
						},
						{
							Name: proto.String("Retry_1"),
							NodeType: &btpb.BehaviorTree_Node_Retry{
								Retry: &btpb.BehaviorTree_RetryNode{
									MaxTries: *proto.Uint32(3),
									Child:    &btpb.BehaviorTree_Node{Name: proto.String("Child_1")},
									Recovery: &btpb.BehaviorTree_Node{Name: proto.String("Recovery_1")},
								},
							},
						},
						{
							Name: proto.String("Subtree_1"),
							NodeType: &btpb.BehaviorTree_Node_SubTree{
								SubTree: &btpb.BehaviorTree_SubtreeNode{
									Tree: &btpb.BehaviorTree{
										TreeId: proto.String("Tree_D"),
										Root:   &btpb.BehaviorTree_Node{Name: proto.String("Child_1")},
									},
								},
							},
						},
						{
							Name: proto.String("Task_1"),
							NodeType: &btpb.BehaviorTree_Node_Task{
								Task: &btpb.BehaviorTree_TaskNode{
									CalledTreeState: &btpb.BehaviorTree{
										TreeId: proto.String("Tree_E"),
										Root:   &btpb.BehaviorTree_Node{Name: proto.String("Child_1")},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name string
		opts []behaviortree.Option
		want []string
	}{
		{
			name: "normal",
			want: []string{
				"Tree_A",
				"Tree_A/Sequence_1",
				"Tree_A/Sequence_1/Sequence_1",
				"Tree_A/Sequence_1/Parallel_1",
				"Tree_A/Sequence_1/Parallel_1/Child_1",
				"Tree_A/Sequence_1/Parallel_1/Child_2",
				"Tree_A/Sequence_1/Fallback_1",
				"Tree_A/Sequence_1/Fallback_1/Child_1",
				"Tree_A/Sequence_1/Fallback_1/Child_2",
				"Tree_A/Sequence_1/Fallback_2",
				"Tree_A/Sequence_1/Fallback_2/Try_1",
				"Tree_A/Sequence_1/Fallback_2/Condition_AnyOf",
				"Tree_A/Sequence_1/Fallback_2/Try_2",
				"Tree_A/Sequence_1/Fallback_2/Condition_BehaviorTree",
				"Tree_A/Sequence_1/Fallback_2/Condition_BehaviorTree/Tree_B",
				"Tree_A/Sequence_1/Fallback_2/Condition_BehaviorTree/Tree_B/Sequence_2",
				"Tree_A/Sequence_1/Fallback_2/Try_3",
				"Tree_A/Sequence_1/Selector_1",
				"Tree_A/Sequence_1/Selector_1/Child_1",
				"Tree_A/Sequence_1/Selector_1/Child_2",
				"Tree_A/Sequence_1/Selector_2",
				"Tree_A/Sequence_1/Selector_2/Branch_1",
				"Tree_A/Sequence_1/Selector_2/Condition_AnyOf",
				"Tree_A/Sequence_1/Selector_2/Branch_2",
				"Tree_A/Sequence_1/Selector_2/Condition_BehaviorTree",
				"Tree_A/Sequence_1/Selector_2/Condition_BehaviorTree/Tree_C",
				"Tree_A/Sequence_1/Selector_2/Condition_BehaviorTree/Tree_C/Sequence_3",
				"Tree_A/Sequence_1/Selector_2/Branch_3",
				"Tree_A/Sequence_1/Retry_1",
				"Tree_A/Sequence_1/Retry_1/Child_1",
				"Tree_A/Sequence_1/Retry_1/Recovery_1",
				"Tree_A/Sequence_1/Subtree_1",
				"Tree_A/Sequence_1/Subtree_1/Tree_D",
				"Tree_A/Sequence_1/Subtree_1/Tree_D/Child_1",
				"Tree_A/Sequence_1/Task_1",
			},
		},
		{
			name: "visit_called_tree_state",
			opts: []behaviortree.Option{behaviortree.VisitCalledTreeState()},
			want: []string{
				"Tree_A",
				"Tree_A/Sequence_1",
				"Tree_A/Sequence_1/Sequence_1",
				"Tree_A/Sequence_1/Parallel_1",
				"Tree_A/Sequence_1/Parallel_1/Child_1",
				"Tree_A/Sequence_1/Parallel_1/Child_2",
				"Tree_A/Sequence_1/Fallback_1",
				"Tree_A/Sequence_1/Fallback_1/Child_1",
				"Tree_A/Sequence_1/Fallback_1/Child_2",
				"Tree_A/Sequence_1/Fallback_2",
				"Tree_A/Sequence_1/Fallback_2/Try_1",
				"Tree_A/Sequence_1/Fallback_2/Condition_AnyOf",
				"Tree_A/Sequence_1/Fallback_2/Try_2",
				"Tree_A/Sequence_1/Fallback_2/Condition_BehaviorTree",
				"Tree_A/Sequence_1/Fallback_2/Condition_BehaviorTree/Tree_B",
				"Tree_A/Sequence_1/Fallback_2/Condition_BehaviorTree/Tree_B/Sequence_2",
				"Tree_A/Sequence_1/Fallback_2/Try_3",
				"Tree_A/Sequence_1/Selector_1",
				"Tree_A/Sequence_1/Selector_1/Child_1",
				"Tree_A/Sequence_1/Selector_1/Child_2",
				"Tree_A/Sequence_1/Selector_2",
				"Tree_A/Sequence_1/Selector_2/Branch_1",
				"Tree_A/Sequence_1/Selector_2/Condition_AnyOf",
				"Tree_A/Sequence_1/Selector_2/Branch_2",
				"Tree_A/Sequence_1/Selector_2/Condition_BehaviorTree",
				"Tree_A/Sequence_1/Selector_2/Condition_BehaviorTree/Tree_C",
				"Tree_A/Sequence_1/Selector_2/Condition_BehaviorTree/Tree_C/Sequence_3",
				"Tree_A/Sequence_1/Selector_2/Branch_3",
				"Tree_A/Sequence_1/Retry_1",
				"Tree_A/Sequence_1/Retry_1/Child_1",
				"Tree_A/Sequence_1/Retry_1/Recovery_1",
				"Tree_A/Sequence_1/Subtree_1",
				"Tree_A/Sequence_1/Subtree_1/Tree_D",
				"Tree_A/Sequence_1/Subtree_1/Tree_D/Child_1",
				"Tree_A/Sequence_1/Task_1",
				"Tree_A/Sequence_1/Task_1/Tree_E",
				"Tree_A/Sequence_1/Task_1/Tree_E/Child_1",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			visitor := &nodeNameCollector{}
			err := behaviortree.Walk(context.Background(), tree, visitor, tc.opts...)
			if err != nil {
				t.Errorf("Tree walker failed on \n%v\ngot %v", tree, err)
			}
			if !cmp.Equal(visitor.names, tc.want) {
				t.Errorf("Failed on \n%v\ngot %v, want %v", tree, visitor.names, tc.want)
			}
		})
	}
}

func TestConditions(t *testing.T) {
	tree := &btpb.BehaviorTree{
		TreeId: proto.String("Tree_A"),
		Root: &btpb.BehaviorTree_Node{
			Name: proto.String("Sequence_1"),
			Decorators: &btpb.BehaviorTree_Node_Decorators{
				Condition: &btpb.BehaviorTree_Condition{
					ConditionType: &btpb.BehaviorTree_Condition_AllOf{
						AllOf: &btpb.BehaviorTree_Condition_LogicalCompound{
							Conditions: []*btpb.BehaviorTree_Condition{
								{
									ConditionType: &btpb.BehaviorTree_Condition_BehaviorTree{
										BehaviorTree: &btpb.BehaviorTree{
											TreeId: proto.String("Tree_B"),
											Root: &btpb.BehaviorTree_Node{
												Name: proto.String("Sequence_2"),
												NodeType: &btpb.BehaviorTree_Node_Sequence{
													Sequence: &btpb.BehaviorTree_SequenceNode{},
												},
											},
										},
									},
								},
								{
									ConditionType: &btpb.BehaviorTree_Condition_BehaviorTree{
										BehaviorTree: &btpb.BehaviorTree{
											TreeId: proto.String("Tree_C"),
											Root: &btpb.BehaviorTree_Node{
												Name: proto.String("Sequence_3"),
												NodeType: &btpb.BehaviorTree_Node_Sequence{
													Sequence: &btpb.BehaviorTree_SequenceNode{},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			NodeType: &btpb.BehaviorTree_Node_Sequence{
				Sequence: &btpb.BehaviorTree_SequenceNode{
					Children: []*btpb.BehaviorTree_Node{
						{
							Name: proto.String("Child_1"),
							Decorators: &btpb.BehaviorTree_Node_Decorators{
								Condition: &btpb.BehaviorTree_Condition{
									ConditionType: &btpb.BehaviorTree_Condition_AnyOf{
										AnyOf: &btpb.BehaviorTree_Condition_LogicalCompound{
											Conditions: []*btpb.BehaviorTree_Condition{
												{
													ConditionType: &btpb.BehaviorTree_Condition_BehaviorTree{
														BehaviorTree: &btpb.BehaviorTree{
															TreeId: proto.String("Tree_D"),
															Root: &btpb.BehaviorTree_Node{
																Name: proto.String("Sequence_4"),
																NodeType: &btpb.BehaviorTree_Node_Sequence{
																	Sequence: &btpb.BehaviorTree_SequenceNode{},
																},
															},
														},
													},
												},
												{
													ConditionType: &btpb.BehaviorTree_Condition_BehaviorTree{
														BehaviorTree: &btpb.BehaviorTree{
															TreeId: proto.String("Tree_E"),
															Root: &btpb.BehaviorTree_Node{
																Name: proto.String("Sequence_5"),
																NodeType: &btpb.BehaviorTree_Node_Sequence{
																	Sequence: &btpb.BehaviorTree_SequenceNode{},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Name: proto.String("Child_2"),
							Decorators: &btpb.BehaviorTree_Node_Decorators{
								Condition: &btpb.BehaviorTree_Condition{
									ConditionType: &btpb.BehaviorTree_Condition_Not{
										Not: &btpb.BehaviorTree_Condition{
											ConditionType: &btpb.BehaviorTree_Condition_BehaviorTree{
												BehaviorTree: &btpb.BehaviorTree{
													TreeId: proto.String("Tree_F"),
													Root: &btpb.BehaviorTree_Node{
														Name: proto.String("Sequence_6"),
														NodeType: &btpb.BehaviorTree_Node_Sequence{
															Sequence: &btpb.BehaviorTree_SequenceNode{},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Name: proto.String("Child_3"),
							NodeType: &btpb.BehaviorTree_Node_Branch{
								Branch: &btpb.BehaviorTree_BranchNode{
									If: &btpb.BehaviorTree_Condition{
										ConditionType: &btpb.BehaviorTree_Condition_BehaviorTree{
											BehaviorTree: &btpb.BehaviorTree{
												TreeId: proto.String("Tree_G"),
												Root: &btpb.BehaviorTree_Node{
													Name: proto.String("Sequence_7"),
													NodeType: &btpb.BehaviorTree_Node_Sequence{
														Sequence: &btpb.BehaviorTree_SequenceNode{},
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Name: proto.String("Child_4"),
							NodeType: &btpb.BehaviorTree_Node_Loop{
								Loop: &btpb.BehaviorTree_LoopNode{
									LoopType: &btpb.BehaviorTree_LoopNode_While{
										While: &btpb.BehaviorTree_Condition{
											ConditionType: &btpb.BehaviorTree_Condition_BehaviorTree{
												BehaviorTree: &btpb.BehaviorTree{
													TreeId: proto.String("Tree_H"),
													Root: &btpb.BehaviorTree_Node{
														Name: proto.String("Sequence_8"),
														NodeType: &btpb.BehaviorTree_Node_Sequence{
															Sequence: &btpb.BehaviorTree_SequenceNode{},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	visitor := &nodeNameCollector{}
	want := []string{
		"Tree_A",
		"Tree_A/Sequence_1",
		"Tree_A/Sequence_1/Condition_AllOf",
		"Tree_A/Sequence_1/Condition_AllOf/Condition_BehaviorTree",
		"Tree_A/Sequence_1/Condition_AllOf/Condition_BehaviorTree/Tree_B",
		"Tree_A/Sequence_1/Condition_AllOf/Condition_BehaviorTree/Tree_B/Sequence_2",
		"Tree_A/Sequence_1/Condition_AllOf/Condition_BehaviorTree",
		"Tree_A/Sequence_1/Condition_AllOf/Condition_BehaviorTree/Tree_C",
		"Tree_A/Sequence_1/Condition_AllOf/Condition_BehaviorTree/Tree_C/Sequence_3",
		"Tree_A/Sequence_1/Child_1",
		"Tree_A/Sequence_1/Child_1/Condition_AnyOf",
		"Tree_A/Sequence_1/Child_1/Condition_AnyOf/Condition_BehaviorTree",
		"Tree_A/Sequence_1/Child_1/Condition_AnyOf/Condition_BehaviorTree/Tree_D",
		"Tree_A/Sequence_1/Child_1/Condition_AnyOf/Condition_BehaviorTree/Tree_D/Sequence_4",
		"Tree_A/Sequence_1/Child_1/Condition_AnyOf/Condition_BehaviorTree",
		"Tree_A/Sequence_1/Child_1/Condition_AnyOf/Condition_BehaviorTree/Tree_E",
		"Tree_A/Sequence_1/Child_1/Condition_AnyOf/Condition_BehaviorTree/Tree_E/Sequence_5",
		"Tree_A/Sequence_1/Child_2",
		"Tree_A/Sequence_1/Child_2/Condition_Not",
		"Tree_A/Sequence_1/Child_2/Condition_Not/Condition_BehaviorTree",
		"Tree_A/Sequence_1/Child_2/Condition_Not/Condition_BehaviorTree/Tree_F",
		"Tree_A/Sequence_1/Child_2/Condition_Not/Condition_BehaviorTree/Tree_F/Sequence_6",
		"Tree_A/Sequence_1/Child_3",
		"Tree_A/Sequence_1/Child_3/Condition_BehaviorTree",
		"Tree_A/Sequence_1/Child_3/Condition_BehaviorTree/Tree_G",
		"Tree_A/Sequence_1/Child_3/Condition_BehaviorTree/Tree_G/Sequence_7",
		"Tree_A/Sequence_1/Child_4",
		"Tree_A/Sequence_1/Child_4/Condition_BehaviorTree",
		"Tree_A/Sequence_1/Child_4/Condition_BehaviorTree/Tree_H",
		"Tree_A/Sequence_1/Child_4/Condition_BehaviorTree/Tree_H/Sequence_8",
	}
	err := behaviortree.Walk(context.Background(), tree, visitor)
	if err != nil {
		t.Errorf("Tree walker failed on \n%v\ngot %v", tree, err)
	}
	if !cmp.Equal(visitor.names, want) {
		t.Errorf("Failed on \n%v\ngot %v, want %v", tree, visitor.names, want)
	}
}

func TestContextCancellation(t *testing.T) {
	tree := &btpb.BehaviorTree{
		Root: &btpb.BehaviorTree_Node{
			Name: proto.String("A"),
			NodeType: &btpb.BehaviorTree_Node_Sequence{
				Sequence: &btpb.BehaviorTree_SequenceNode{
					Children: []*btpb.BehaviorTree_Node{
						{Name: proto.String("B")},
					},
				},
			},
		},
	}

	visitor := &nodeNameCollector{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := behaviortree.Walk(ctx, tree, visitor)
	if err != context.Canceled {
		t.Errorf("behaviortree.Walk() got error %v, want %v", err, context.Canceled)
	}
}

func TestWalk_Stop(t *testing.T) {
	tree := &btpb.BehaviorTree{
		TreeId: proto.String("Tree_A"),
		Root: &btpb.BehaviorTree_Node{
			Name: proto.String("Sequence_1"),
			NodeType: &btpb.BehaviorTree_Node_Sequence{
				Sequence: &btpb.BehaviorTree_SequenceNode{
					Children: []*btpb.BehaviorTree_Node{
						{Name: proto.String("Child_1")},
						{Name: proto.String("Child_2")},
					},
				},
			},
		},
	}
	wantNames := []string{
		"Tree_A",
		"Tree_A/Sequence_1",
		"Tree_A/Sequence_1/Child_1",
	}

	visitor := &nodeNameCollector{
		stopOnName: "Child_1",
	}
	err := behaviortree.Walk(context.Background(), tree, visitor)

	if err != nil {
		t.Errorf("behaviortree.Walk() got unexpected error %v, want nil", err)
	}
	if !cmp.Equal(visitor.names, wantNames) {
		t.Errorf("behaviortree.Walk() got names %v, want %v", visitor.names, wantNames)
	}
}

func TestVisitElement_Ancestors(t *testing.T) {
	tree := &btpb.BehaviorTree{TreeId: proto.String("Tree_A")}
	node1 := &btpb.BehaviorTree_Node{Name: proto.String("Node_1")}
	node2 := &btpb.BehaviorTree_Node{Name: proto.String("Node_2")}
	cond1 := &btpb.BehaviorTree_Condition{
		ConditionType: &btpb.BehaviorTree_Condition_Blackboard{},
	}

	tests := []struct {
		name      string
		element   behaviortree.VisitElement
		wantNames []string
	}{
		{
			name:      "single_element",
			element:   behaviortree.VisitElementFromTree(tree),
			wantNames: nil,
		},
		{
			name: "two_elements",
			element: behaviortree.VisitElementFromTree(tree).
				AsAncestorForNode(node1),
			wantNames: []string{"Tree_A"},
		},
		{
			name: "three_elements",
			element: behaviortree.VisitElementFromTree(tree).
				AsAncestorForNode(node1).
				AsAncestorForNode(node2),
			wantNames: []string{"Node_1", "Tree_A"},
		},
		{
			name: "with_condition",
			element: behaviortree.VisitElementFromTree(tree).
				AsAncestorForNode(node1).
				AsAncestorForCondition(cond1).
				AsAncestorForNode(node2),
			wantNames: []string{"Condition_Blackboard", "Node_1", "Tree_A"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotNames []string
			for ancestor := range tc.element.Ancestors() {
				gotNames = append(gotNames, elementName(ancestor))
			}
			if diff := cmp.Diff(tc.wantNames, gotNames); diff != "" {
				t.Errorf("Ancestors() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
