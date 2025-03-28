// Copyright 2023 Intrinsic Innovation LLC

package behaviortree_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"intrinsic/executive/go/behaviortree"
	btpb "intrinsic/executive/proto/behavior_tree_go_proto"
)

type nodeNameCollector struct {
	Names []string
}

func (c *nodeNameCollector) VisitNode(node *btpb.BehaviorTree_Node) error {
	c.Names = append(c.Names, node.GetName())
	return nil
}
func (c *nodeNameCollector) VisitCondition(cond *btpb.BehaviorTree_Condition) error {
	return nil
}

func TestNodes(t *testing.T) {
	tree := &btpb.BehaviorTree{
		Root: &btpb.BehaviorTree_Node{
			Name: proto.String("A"),
			NodeType: &btpb.BehaviorTree_Node_Sequence{
				Sequence: &btpb.BehaviorTree_SequenceNode{
					Children: []*btpb.BehaviorTree_Node{
						{Name: proto.String("B")},
						{
							Name: proto.String("C"),
							NodeType: &btpb.BehaviorTree_Node_Parallel{
								Parallel: &btpb.BehaviorTree_ParallelNode{
									Children: []*btpb.BehaviorTree_Node{
										{Name: proto.String("D")},
										{Name: proto.String("E")},
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
	want := []string{"A", "B", "C", "D", "E"}
	err := behaviortree.Walk(tree, visitor)
	if err != nil {
		t.Errorf("Tree walker failed on \n%v\ngot %v", tree, err)
	}
	if !cmp.Equal(visitor.Names, want) {
		t.Errorf("Failed on \n%v\ngot %v, want %v", tree, visitor.Names, want)
	}
}

func TestConditions(t *testing.T) {
	tree := &btpb.BehaviorTree{
		Root: &btpb.BehaviorTree_Node{
			Name: proto.String("A"),
			Decorators: &btpb.BehaviorTree_Node_Decorators{
				Condition: &btpb.BehaviorTree_Condition{
					ConditionType: &btpb.BehaviorTree_Condition_AllOf{
						AllOf: &btpb.BehaviorTree_Condition_LogicalCompound{
							Conditions: []*btpb.BehaviorTree_Condition{
								{
									ConditionType: &btpb.BehaviorTree_Condition_BehaviorTree{
										BehaviorTree: &btpb.BehaviorTree{
											Root: &btpb.BehaviorTree_Node{
												Name: proto.String("B"),
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
											Root: &btpb.BehaviorTree_Node{
												Name: proto.String("C"),
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
							Name: proto.String("D"),
							Decorators: &btpb.BehaviorTree_Node_Decorators{
								Condition: &btpb.BehaviorTree_Condition{
									ConditionType: &btpb.BehaviorTree_Condition_AnyOf{
										AnyOf: &btpb.BehaviorTree_Condition_LogicalCompound{
											Conditions: []*btpb.BehaviorTree_Condition{
												{
													ConditionType: &btpb.BehaviorTree_Condition_BehaviorTree{
														BehaviorTree: &btpb.BehaviorTree{
															Root: &btpb.BehaviorTree_Node{
																Name: proto.String("E"),
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
															Root: &btpb.BehaviorTree_Node{
																Name: proto.String("F"),
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
							Name: proto.String("G"),
							Decorators: &btpb.BehaviorTree_Node_Decorators{
								Condition: &btpb.BehaviorTree_Condition{
									ConditionType: &btpb.BehaviorTree_Condition_Not{
										Not: &btpb.BehaviorTree_Condition{
											ConditionType: &btpb.BehaviorTree_Condition_BehaviorTree{
												BehaviorTree: &btpb.BehaviorTree{
													Root: &btpb.BehaviorTree_Node{
														Name: proto.String("H"),
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
							Name: proto.String("I"),
							NodeType: &btpb.BehaviorTree_Node_Branch{
								Branch: &btpb.BehaviorTree_BranchNode{
									If: &btpb.BehaviorTree_Condition{
										ConditionType: &btpb.BehaviorTree_Condition_BehaviorTree{
											BehaviorTree: &btpb.BehaviorTree{
												Root: &btpb.BehaviorTree_Node{
													Name: proto.String("J"),
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
							Name: proto.String("K"),
							NodeType: &btpb.BehaviorTree_Node_Loop{
								Loop: &btpb.BehaviorTree_LoopNode{
									LoopType: &btpb.BehaviorTree_LoopNode_While{
										While: &btpb.BehaviorTree_Condition{
											ConditionType: &btpb.BehaviorTree_Condition_BehaviorTree{
												BehaviorTree: &btpb.BehaviorTree{
													Root: &btpb.BehaviorTree_Node{
														Name: proto.String("L"),
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
	want := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L"}
	err := behaviortree.Walk(tree, visitor)
	if err != nil {
		t.Errorf("Tree walker failed on \n%v\ngot %v", tree, err)
	}
	if !cmp.Equal(visitor.Names, want) {
		t.Errorf("Failed on \n%v\ngot %v, want %v", tree, visitor.Names, want)
	}
}
