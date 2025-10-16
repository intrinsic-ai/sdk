// Copyright 2023 Intrinsic Innovation LLC

// Package behaviortree provides utilities to operate on Behavior Trees.
//
// Features are:
// - Enables to walk and execute code for each node and condition in the tree.
package behaviortree

import (
	"context"

	btpb "intrinsic/executive/proto/behavior_tree_go_proto"
)

// The Visitor defines requirements for visitor implementations for the
// BehaviorTree proto walker.
type Visitor interface {
	// Visit a specific condition in the tree.
	VisitCondition(ctx context.Context, cond *btpb.BehaviorTree_Condition) error
	// Visit a specific node in the tree.
	VisitNode(ctx context.Context, node *btpb.BehaviorTree_Node) error
}

func walkCondition(ctx context.Context, cond *btpb.BehaviorTree_Condition, visitor Visitor) error {
	if cond == nil {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	err := visitor.VisitCondition(ctx, cond)
	if err != nil {
		return err
	}
	switch cond.ConditionType.(type) {
	case *btpb.BehaviorTree_Condition_BehaviorTree:
		err := Walk(ctx, cond.GetBehaviorTree(), visitor)
		if err != nil {
			return err
		}

	case *btpb.BehaviorTree_Condition_AllOf:
		for _, c := range cond.GetAllOf().GetConditions() {
			err := walkCondition(ctx, c, visitor)
			if err != nil {
				return err
			}
		}

	case *btpb.BehaviorTree_Condition_AnyOf:
		for _, c := range cond.GetAnyOf().GetConditions() {
			err := walkCondition(ctx, c, visitor)
			if err != nil {
				return err
			}
		}

	case *btpb.BehaviorTree_Condition_Not:
		err = walkCondition(ctx, cond.GetNot(), visitor)
		if err != nil {
			return err
		}
	}

	return nil
}

func walkNode(ctx context.Context, node *btpb.BehaviorTree_Node, visitor Visitor) error {
	if node == nil {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	err := visitor.VisitNode(ctx, node)
	if err != nil {
		return err
	}

	if node.GetDecorators().GetCondition() != nil {
		err := walkCondition(ctx, node.GetDecorators().GetCondition(), visitor)
		if err != nil {
			return err
		}
	}

	switch node.NodeType.(type) {
	case *btpb.BehaviorTree_Node_Sequence:
		for _, c := range node.GetSequence().GetChildren() {
			err := walkNode(ctx, c, visitor)
			if err != nil {
				return err
			}
		}

	case *btpb.BehaviorTree_Node_Parallel:
		for _, c := range node.GetParallel().GetChildren() {
			err := walkNode(ctx, c, visitor)
			if err != nil {
				return err
			}
		}

	case *btpb.BehaviorTree_Node_Selector:
		for _, c := range node.GetSelector().GetChildren() {
			err := walkNode(ctx, c, visitor)
			if err != nil {
				return err
			}
		}
		for _, c := range node.GetSelector().GetBranches() {
			err := walkCondition(ctx, c.GetCondition(), visitor)
			if err != nil {
				return err
			}
			err = walkNode(ctx, c.GetNode(), visitor)
			if err != nil {
				return err
			}
		}

	case *btpb.BehaviorTree_Node_Fallback:
		for _, c := range node.GetFallback().GetChildren() {
			err := walkNode(ctx, c, visitor)
			if err != nil {
				return err
			}
		}
		for _, c := range node.GetFallback().GetTries() {
			err := walkCondition(ctx, c.GetCondition(), visitor)
			if err != nil {
				return err
			}
			err = walkNode(ctx, c.GetNode(), visitor)
			if err != nil {
				return err
			}
		}

	case *btpb.BehaviorTree_Node_Branch:
		err := walkCondition(ctx, node.GetBranch().GetIf(), visitor)
		if err != nil {
			return err
		}
		err = walkNode(ctx, node.GetBranch().GetThen(), visitor)
		if err != nil {
			return err
		}
		err = walkNode(ctx, node.GetBranch().GetElse(), visitor)
		if err != nil {
			return err
		}

	case *btpb.BehaviorTree_Node_Loop:
		err := walkCondition(ctx, node.GetLoop().GetWhile(), visitor)
		if err != nil {
			return err
		}
		err = walkNode(ctx, node.GetLoop().GetDo(), visitor)
		if err != nil {
			return err
		}

	case *btpb.BehaviorTree_Node_Retry:
		err := walkNode(ctx, node.GetRetry().GetChild(), visitor)
		if err != nil {
			return err
		}

	case *btpb.BehaviorTree_Node_SubTree:
		err := Walk(ctx, node.GetSubTree().GetTree(), visitor)
		if err != nil {
			return err
		}
	}

	return nil
}

// Walk walks the given Behavior Tree and invokes the given visitor for nodes
// and conditions of the tree.
func Walk(ctx context.Context, tree *btpb.BehaviorTree, visitor Visitor) error {
	return walkNode(ctx, tree.Root, visitor)
}
