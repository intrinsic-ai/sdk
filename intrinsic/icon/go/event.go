// Copyright 2023 Intrinsic Innovation LLC

package icon

// Event is a session event, originated server-side.
type Event interface {
	isEvent()
}

// ReactionEvent signifies that a reaction has triggered on the server.
type ReactionEvent struct {
	// EventSignal is a marker that identifies the reaction.
	EventSignal EventSignal
}

func (*ReactionEvent) isEvent() {}
