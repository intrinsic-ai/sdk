// Copyright 2023 Intrinsic Innovation LLC

package icon

// This file contains internal helpers for tracking relations between
// ActionHandles, ReactionIDs and EventSignals.

import (
	"fmt"
	"sync"
)

var (
	errInvalidReactionID           = fmt.Errorf("want a non-zero ReactionID")
	errAlreadyExists               = fmt.Errorf("already exists")
	errReactionIDAlreadyAssociated = fmt.Errorf("ReactionID is already associated")
)

// reactionDataEntrySet is a set of reactionDataEntry objects.
type reactionDataEntrySet map[reactionDataEntry]struct{}

// reactionDataEntry is an (ActionHandle, ReactionID, EventSignal) tuple.
type reactionDataEntry struct {
	// Action identifies a server-side action instance. If nil, the reaction is free-standing.
	Action *ActionHandle
	// ReactID identifies a server-side reaction instance.
	ReactID ReactionID
	// Signal identifies a client-side reaction event marker that should be
	// emitted when reaction with ReactID occurs.
	Signal EventSignal
}

// reactionData stores relations between ActionHandles, ReactionIDs and
// EventSignals. This data structure is used to support: (i) efficient
// insertion, (ii) lookup by ReactionID and (iii) removal by ActionHandle. The
// following constraints are enforced: (1) Each (Action, ReactID, Signal) tuple
// is unique; (2) Each ReactID is associated with at most 1 Action. All methods
// are safe to use concurrently.
type reactionData struct {
	// entries ensures no duplicate entries.
	entries reactionDataEntrySet
	// actionByReaction ensures that a ReactionID is not associated with
	// multiple ActionHandles.
	actionByReaction map[ReactionID]*ActionHandle
	// entriesByAction supports efficient removal of entries matching an
	// ActionHandle.
	entriesByAction map[ActionHandle][]reactionDataEntry
	// signalsByReaction supports efficient lookup of EventSignals by
	// ReactionID.
	signalsByReaction map[ReactionID][]EventSignal
	mu                sync.Mutex
}

// newReactionData creates a new reactionData data store.
func newReactionData() *reactionData {
	return &reactionData{
		entries:           make(reactionDataEntrySet),
		actionByReaction:  make(map[ReactionID]*ActionHandle),
		entriesByAction:   make(map[ActionHandle][]reactionDataEntry),
		signalsByReaction: make(map[ReactionID][]EventSignal),
	}
}

// Insert inserts all entries, or nothing on failure.
func (d *reactionData) Insert(entries ...reactionDataEntry) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	// Detect duplicates among passed-in entries.
	es := make(reactionDataEntrySet)
	for _, e := range entries {
		if _, exists := es[e]; exists {
			return fmt.Errorf("reactionData entry %v: %w", e, errAlreadyExists)
		}
		es[e] = struct{}{}
	}
	// Detect already associated ReactIDs among passed-in entries.
	ar := make(map[ReactionID]*ActionHandle)
	for _, e := range entries {
		if _, exists := ar[e.ReactID]; exists {
			return fmt.Errorf("cannot insert reactionData entry %v: %w", e, errReactionIDAlreadyAssociated)
		}
		ar[e.ReactID] = e.Action

	}
	// Further validity checks, including against stored data.
	for _, e := range entries {
		if e.Action != nil && e.Action.IsZero() {
			return fmt.Errorf("invalid reactionData entry %v: %w", e, ErrInvalidActionHandle)
		}
		if e.ReactID == 0 {
			return fmt.Errorf("invalid reactionData entry %v: %w", e, errInvalidReactionID)
		}
		if e.Signal.IsEmpty() {
			return fmt.Errorf("invalid reactionData entry %v: %w", e, ErrInvalidEventSignal)
		}
		if _, exists := d.entries[e]; exists {
			return fmt.Errorf("reactionData entry %v: %w", e, errAlreadyExists)
		}
		if _, exists := d.actionByReaction[e.ReactID]; exists {
			return fmt.Errorf("cannot insert reactionData entry %v: %w", e, errReactionIDAlreadyAssociated)
		}
	}
	for _, e := range entries {
		d.entries[e] = struct{}{}
		d.actionByReaction[e.ReactID] = e.Action
		if e.Action != nil {
			d.entriesByAction[*e.Action] = append(d.entriesByAction[*e.Action], e)
		}
		d.signalsByReaction[e.ReactID] = append(d.signalsByReaction[e.ReactID], e.Signal)
	}
	return nil
}

// RemoveActions removes all entries matching actions.
func (d *reactionData) RemoveActions(actions ...ActionHandle) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, a := range actions {
		for _, e := range d.entriesByAction[a] {
			delete(d.entries, e)
			delete(d.actionByReaction, e.ReactID)
			delete(d.signalsByReaction, e.ReactID)
		}
		delete(d.entriesByAction, a)
	}
}

// Clear resets d, clearing all data.
func (d *reactionData) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.entries = make(reactionDataEntrySet)
	d.actionByReaction = make(map[ReactionID]*ActionHandle)
	d.entriesByAction = make(map[ActionHandle][]reactionDataEntry)
	d.signalsByReaction = make(map[ReactionID][]EventSignal)
}

// ListSignals lists all EventSignals from entries matching r. Returns nil if
// none are found.
func (d *reactionData) ListSignals(r ReactionID) []EventSignal {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.signalsByReaction[r]
}
