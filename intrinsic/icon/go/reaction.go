// Copyright 2023 Intrinsic Innovation LLC

package icon

import (
	"errors"
	"fmt"

	conditiontypespb "intrinsic/icon/proto/v1/condition_types_go_proto"
	typespb "intrinsic/icon/proto/v1/types_go_proto"
)

var (
	// ErrMissingResponse occurs if a Reaction has no Responses.
	ErrMissingResponse = errors.New("a Reaction needs at least one Response")

	// ErrInvalidActionHandle occurs if a zero-valued ActionHandle is used
	// when a valid ActionHandle is expected.
	ErrInvalidActionHandle = errors.New("invalid ActionHandle; want a handle created with session.MakeActionHandle()")

	// ErrInvalidEventSignal occurs if a zero-valued EventSignal is used
	// when a valid EventSignal is expected.
	ErrInvalidEventSignal = errors.New("invalid EventSignal; want a signal created with session.MakeEventSignal()")
)

// ReactionID is an int64 used to identify a server-side reaction instance.
type ReactionID int64

// EventSignal is a marker used to watch for a reaction event. Create a unique
// marker with session.MakeEventSignal(). Use EmitEventSignal to describe a
// reaction that will emit an event with the specified marker, when it
// triggers. Finally, call session.NextEvent to watch for the marker.
//
// EventSignal objects may be safely compared for equality (==, !=), copied by
// value, and used as keys in a map, but must not be deep copied.
type EventSignal struct {
	// s is the session that created this EventSignal.
	session *Session
	// unique is a value that is meaningless other than to provide uniqueness.
	unique int64
}

// Equal compares two EventSignal objects.
func (e EventSignal) Equal(o EventSignal) bool {
	return e == o
}

// IsEmpty reports whether e is zero-valued, which represents an invalid
// EventSignal.
func (e EventSignal) IsEmpty() bool {
	return e == EventSignal{}
}

// Response describes a reaction's response.
// This is one of the following concrete classes:
// - StartActionInRealTimeResponse
// - EmitEventSignalResponse
type Response interface {
	isResponse()
}

// StartActionInRealTimeResponse describes a reaction response which starts
// another action before the next control cycle.
type StartActionInRealTimeResponse struct {
	// Handle identifies the action to start in real-time.
	Handle ActionHandle
}

func (StartActionInRealTimeResponse) isResponse() {}

// StartActionInRealTime creates a new StartActionInRealTimeResponse.
func StartActionInRealTime(h ActionHandle) Response {
	return &StartActionInRealTimeResponse{h}
}

// StartParallelActionInRealTimeResponse describes a reaction response which starts
// another action in parallel before the next control cycle.
type StartParallelActionInRealTimeResponse struct {
	// Handle identifies the action to start in real-time.
	Handle ActionHandle
}

func (StartParallelActionInRealTimeResponse) isResponse() {}

// StartParallelActionInRealTime creates a new StartParallelActionInRealTimeResponse.
func StartParallelActionInRealTime(h ActionHandle) Response {
	return &StartParallelActionInRealTimeResponse{h}
}

// EmitEventSignalResponse describes a reaction response which emits a
// non-real-time event marker that is observable by calling
// session.NextEvent().
type EmitEventSignalResponse struct {
	// Signal identifies the event marker.
	Signal EventSignal
}

func (EmitEventSignalResponse) isResponse() {}

// EmitEventSignal creates a new EmitEventSignalResponse.
func EmitEventSignal(es EventSignal) Response {
	return &EmitEventSignalResponse{es}
}

// Reaction describes a reaction, which is a real-time condition along with
// one or more responses. On the server, all currently-running actions'
// reaction conditions are evaluated every control cycle. When such a
// condition evaluates to True, the corresponding responses occur. `stopAssociatedAction`
// configures whether the associated action of this reaction stops or
// continues to run when the reaction triggers an action change.
type Reaction struct {
	Condition            *conditiontypespb.Condition
	Responses            []Response
	stopAssociatedAction bool
}

// NewReaction creates a new Reaction with the provided condition and
// responses.
func NewReaction(c *conditiontypespb.Condition, rs ...Response) *Reaction {
	return &Reaction{
		Condition: c,
		Responses: rs,
	}
}

// extractEventSignals extracts a list of all EventSignals that are in r's
// EmitEventSignal responses. Returns an error if any of those EventSignals are
// zero-valued.
func (r *Reaction) extractEventSignals() ([]EventSignal, error) {
	var rhs []EventSignal
	for _, response := range r.Responses {
		switch resp := response.(type) {
		case *EmitEventSignalResponse:
			if resp.Signal.IsEmpty() {
				return nil, ErrInvalidEventSignal
			}
			rhs = append(rhs, resp.Signal)
		}
	}
	return rhs, nil
}

// proto makes a reaction Proto from the receiver, with the given action handle
// and reaction ID.
func (r *Reaction) proto(ah *ActionHandle, id ReactionID) (proto *typespb.Reaction, err error) {
	var responseProto *typespb.Response
	if ah != nil && ah.IsZero() {
		return nil, ErrInvalidActionHandle
	}
	var stopAssociatedAction bool = false
	if len(r.Responses) == 0 {
		return nil, ErrMissingResponse
	}
	for _, response := range r.Responses {
		switch resp := response.(type) {
		case *StartActionInRealTimeResponse:
			stopAssociatedAction = true
			if responseProto != nil {
				return nil, fmt.Errorf("at most one Start(Parallel)ActionInRealTime response is allowed per reaction; r=%v", r)
			}
			if resp.Handle.IsZero() {
				return nil, ErrInvalidActionHandle
			}
			responseProto = &typespb.Response{
				StartActionInstanceId: int64(resp.Handle.ID()),
			}
		case *StartParallelActionInRealTimeResponse:
			stopAssociatedAction = false
			if responseProto != nil {
				return nil, fmt.Errorf("at most one Start(Parallel)ActionInRealTime response is allowed per reaction; r=%v", r)
			}
			if resp.Handle.IsZero() {
				return nil, ErrInvalidActionHandle
			}
			responseProto = &typespb.Response{
				StartActionInstanceId: int64(resp.Handle.ID()),
			}
		case *EmitEventSignalResponse:
			// EmitEventSignal responses do not affect the request payload.
		default:
			return nil, fmt.Errorf("unsupported response type %T", resp)
		}
	}
	var actionAssociation *typespb.Reaction_ActionAssociation = nil
	if ah != nil {
		actionAssociation = &typespb.Reaction_ActionAssociation{
			ActionInstanceId:     int64(ah.ID()),
			StopAssociatedAction: stopAssociatedAction,
		}
	}

	return &typespb.Reaction{
		ReactionInstanceId: int64(id),
		ActionAssociation:  actionAssociation,
		Condition:          r.Condition,
		Response:           responseProto,
	}, nil
}
