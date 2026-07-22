// Copyright 2023 Intrinsic Innovation LLC

package icon

// This file contains helpers for describing Conditions, which are expressions
// that reference an Action's State Variables by string name.  On the ICON
// server, an Action instance may have associated Reactions. Each Reaction has
// a Condition and Responses. An active Action's Conditions are evaluated every
// control cycle. If a Condition evaluates to True, the corresponding Responses
// occur.

import (
	conditiontypespb "intrinsic/icon/proto/v1/condition_types_go_proto"
)

// boolComparison describes a condition comparing a boolean state variable's
// value with v.
func boolComparison(stateVar string, op conditiontypespb.Comparison_OpEnum, v bool) *conditiontypespb.Condition {
	return &conditiontypespb.Condition{
		Condition: &conditiontypespb.Condition_Comparison{
			Comparison: &conditiontypespb.Comparison{
				StateVariableName: stateVar,
				Operation:         op,
				Value: &conditiontypespb.Comparison_BoolValue{
					BoolValue: v,
				},
			},
		},
	}
}

// IsTrue describes a condition which is satisfied when a boolean state
// variable's value is True.
func IsTrue(stateVar string) *conditiontypespb.Condition {
	return EqualBool(stateVar, true)
}

// IsFalse describes a condition which is satisfied when a boolean state
// variable's value is False.
func IsFalse(stateVar string) *conditiontypespb.Condition {
	return EqualBool(stateVar, false)
}

// EqualBool describes a condition which is satisfied when a boolean state
// variable's value equals v.
func EqualBool(stateVar string, v bool) *conditiontypespb.Condition {
	return boolComparison(stateVar, conditiontypespb.Comparison_EQUAL, v)
}

// floatComparison describes a condition that approximately compares a
// floating-point state variable's value with v, using a tolerance of epsilon.
func floatComparison(stateVar string, op conditiontypespb.Comparison_OpEnum, v, epsilon float64) *conditiontypespb.Condition {
	return &conditiontypespb.Condition{
		Condition: &conditiontypespb.Condition_Comparison{
			Comparison: &conditiontypespb.Comparison{
				StateVariableName: stateVar,
				Operation:         op,
				Value: &conditiontypespb.Comparison_DoubleValue{
					DoubleValue: v,
				},
				MaxAbsError: epsilon,
			},
		},
	}
}

// int64Comparison describes a condition that compares an integer state variable's value with v.
func int64Comparison(stateVar string, op conditiontypespb.Comparison_OpEnum, v int64) *conditiontypespb.Condition {
	return &conditiontypespb.Condition{
		Condition: &conditiontypespb.Condition_Comparison{
			Comparison: &conditiontypespb.Comparison{
				StateVariableName: stateVar,
				Operation:         op,
				Value: &conditiontypespb.Comparison_Int64Value{
					Int64Value: v,
				},
			},
		},
	}
}

// exactFloatComparison describes a condition comparing a floating-point state
// variable's value with v, with zero tolerance.
func exactFloatComparison(stateVar string, op conditiontypespb.Comparison_OpEnum, v float64) *conditiontypespb.Condition {
	return floatComparison(stateVar, op, v, 0)
}

// EqualInt64 describes a condition which is satisfied when an integer
// state variable's value is equal to v.
func EqualInt64(stateVar string, v int64) *conditiontypespb.Condition {
	return int64Comparison(stateVar, conditiontypespb.Comparison_EQUAL, v)
}

// NotEqualInt64 describes a condition which is satisfied when an integer
// state variable's value is not equal to v.
func NotEqualInt64(stateVar string, v int64) *conditiontypespb.Condition {
	return int64Comparison(stateVar, conditiontypespb.Comparison_NOT_EQUAL, v)
}

// ApproxEqual describes a condition which is satisfied when a float-point
// state variable's value is within epsilon of v.
func ApproxEqual(stateVar string, v, epsilon float64) *conditiontypespb.Condition {
	return floatComparison(stateVar, conditiontypespb.Comparison_APPROX_EQUAL, v, epsilon)
}

// ApproxNotEqual describes a condition which is satisfied when a float-point
// state variable's value is not within epsilon of v.
func ApproxNotEqual(stateVar string, v, epsilon float64) *conditiontypespb.Condition {
	return floatComparison(stateVar, conditiontypespb.Comparison_APPROX_NOT_EQUAL, v, epsilon)
}

// GreaterThan describes a condition which is satisfied when a floating-point
// state variable's value is greater than v.
func GreaterThan(stateVar string, v float64) *conditiontypespb.Condition {
	return exactFloatComparison(stateVar, conditiontypespb.Comparison_GREATER_THAN, v)
}

// GreaterThanInt64 describes a condition which is satisfied when an integer
// state variable's value is greater than v.
func GreaterThanInt64(stateVar string, v int64) *conditiontypespb.Condition {
	return int64Comparison(stateVar, conditiontypespb.Comparison_GREATER_THAN, v)
}

// GreaterThanOrEqual describes a condition which is satisfied when a
// floating-point state variable's value is greater than or equal to v.
func GreaterThanOrEqual(stateVar string, v float64) *conditiontypespb.Condition {
	return exactFloatComparison(stateVar, conditiontypespb.Comparison_GREATER_THAN_OR_EQUAL, v)
}

// GreaterThanOrEqualToInt64 describes a condition which is satisfied when an integer
// state variable's value is greater than or equal to v.
func GreaterThanOrEqualToInt64(stateVar string, v int64) *conditiontypespb.Condition {
	return int64Comparison(stateVar, conditiontypespb.Comparison_GREATER_THAN_OR_EQUAL, v)
}

// LessThan describes a condition which is satisfied when a floating-point
// state variable's value is less than v.
func LessThan(stateVar string, v float64) *conditiontypespb.Condition {
	return exactFloatComparison(stateVar, conditiontypespb.Comparison_LESS_THAN, v)
}

// LessThanInt64 describes a condition which is satisfied when an integer
// state variable's value is less than to v.
func LessThanInt64(stateVar string, v int64) *conditiontypespb.Condition {
	return int64Comparison(stateVar, conditiontypespb.Comparison_LESS_THAN, v)
}

// LessThanOrEqual describes a condition which is satisfied when a
// floating-point state variable's value is less than or equal to v.
func LessThanOrEqual(stateVar string, v float64) *conditiontypespb.Condition {
	return exactFloatComparison(stateVar, conditiontypespb.Comparison_LESS_THAN_OR_EQUAL, v)
}

// LessThanOrEqualToInt64 describes a condition which is satisfied when an integer
// state variable's value is less than or equal to v.
func LessThanOrEqualToInt64(stateVar string, v int64) *conditiontypespb.Condition {
	return int64Comparison(stateVar, conditiontypespb.Comparison_LESS_THAN_OR_EQUAL, v)
}

// conjunction describes a conjunction condition (i.e. "any of" or "all of").
func conjunction(op conditiontypespb.ConjunctionCondition_OpEnum, conditions ...*conditiontypespb.Condition) *conditiontypespb.Condition {
	return &conditiontypespb.Condition{
		Condition: &conditiontypespb.Condition_ConjunctionCondition{
			ConjunctionCondition: &conditiontypespb.ConjunctionCondition{
				Operation:  op,
				Conditions: conditions,
			},
		},
	}
}

// AnyOf describes a condition which is satisfied when any of its child
// conditions are satisfied.
func AnyOf(conditions ...*conditiontypespb.Condition) *conditiontypespb.Condition {
	return conjunction(conditiontypespb.ConjunctionCondition_ANY_OF, conditions...)
}

// AllOf describes a condition which is satisfied when all of its child
// conditions are satisfied.
func AllOf(conditions ...*conditiontypespb.Condition) *conditiontypespb.Condition {
	return conjunction(conditiontypespb.ConjunctionCondition_ALL_OF, conditions...)
}

// Not describes a condition which is satisfied when its child condition is not satisfied.
func Not(condition *conditiontypespb.Condition) *conditiontypespb.Condition {
	return &conditiontypespb.Condition{
		Condition: &conditiontypespb.Condition_NegatedCondition{
			NegatedCondition: &conditiontypespb.NegatedCondition{
				Condition: condition,
			},
		},
	}
}
