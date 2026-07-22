// Copyright 2023 Intrinsic Innovation LLC

package icon

import (
	"context"
	"errors"
	"fmt"

	"intrinsic/icon/go/intsequence"

	log "github.com/golang/glog"
	codespb "google.golang.org/grpc/codes"

	grpcpb "intrinsic/icon/proto/v1/service_go_proto"
	servicepb "intrinsic/icon/proto/v1/service_go_proto"
	typespb "intrinsic/icon/proto/v1/types_go_proto"
	contextpb "intrinsic/logging/proto/context_go_proto"
)

var (
	// ErrInvalidActionDescription occurs when a bad action description is provided.
	ErrInvalidActionDescription = errors.New("invalid action description")
	// ErrInvalidReactionDescription occurs when a bad reaction description is provided.
	ErrInvalidReactionDescription = errors.New("invalid reaction description")
	// ErrInternal signifies that an internal error has occurred.
	ErrInternal = errors.New("internal error")
	// ErrCanceled occurs when an operation is canceled.
	ErrCanceled = errors.New("operation canceled")
	// ErrSessionEnded occurs when the session is already ended.
	ErrSessionEnded = errors.New("session already ended")
)

// internalError is an error that wraps Err and Is(ErrInternal).
type internalError struct {
	err error
}

// Error implements the error interface.
func (e internalError) Error() string { return fmt.Sprintf("internal error: %v", e.err) }

// Is returns true if target is ErrInternal.
func (e internalError) Is(target error) bool { return target == ErrInternal }

// Unwrap returns the contained error.
func (e internalError) Unwrap() error { return e.err }

// Session claims exclusive control over one or more parts. It provides
// the ability to manipulate those parts by adding actions and/or reactions, and
// bounds the lifetime of these server-side objects.
type Session struct {
	id                     int64 // Session ID
	client                 *grpcClient
	ended                  bool
	session                grpcpb.IconApi_OpenSessionClient
	reactData              *reactionData
	events                 chan Event
	watchReactionsError    chan error
	actionIDs              intsequence.Generator
	reactionIDs            intsequence.Generator
	eventSignalUniqueValue intsequence.Generator
	writeStreams           map[*WriteStream]bool
}

// watchReactions receives reaction updates from the server, converts them to
// Event objects, and sends them to the events channel. This runs until ctx is
// canceled or client.Recv() returns an error, whichever occurs earlier.
func watchReactions(ctx context.Context, client grpcpb.IconApi_WatchReactionsClient, reactData *reactionData, events chan<- Event) error {
	defer close(events)
	for {
		resp, err := client.Recv()
		if err != nil {
			return err
		}
		if resp.ReactionEvent == nil {
			continue
		}
		if resp.ReactionEvent.ReactionId == 0 { // Ignore non-reaction events.
			continue
		}
		signals := reactData.ListSignals(ReactionID(resp.ReactionEvent.ReactionId))
		for _, s := range signals {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case events <- &ReactionEvent{EventSignal: s}:
			}
		}
	}
}

// newSession creates a new Session for a provided list of parts. The ICON logs
// are tagged with the provided `logContext`, which can be nil.
func newSession(ctx context.Context, c *grpcClient, parts []string, logContext *contextpb.Context) (*Session, error) {
	// OpenSession creates a handle to an OpenSessionClient
	// which allows sending/receiving messages from the ICON server.
	req := &servicepb.OpenSessionRequest{
		InitialSessionData: &servicepb.OpenSessionRequest_InitialSessionData{
			AllocateParts: &typespb.PartList{Part: parts},
		},
		LogContext: logContext,
	}
	session, err := c.client.OpenSession(ctx)
	if err != nil {
		return nil, err
	}
	if err := session.Send(req); err != nil {
		return nil, err
	}
	resp, err := session.Recv()
	if err != nil {
		return nil, err
	}
	if resp.Status.Code != int32(codespb.OK) {
		return nil, fmt.Errorf("session gRPC failed: %v", resp.Status.Message)
	}
	if resp.InitialSessionData == nil {
		return nil, fmt.Errorf("missing initial session data")
	}
	s := &Session{
		id:                  resp.InitialSessionData.SessionId,
		client:              c,
		ended:               false,
		session:             session,
		reactData:           newReactionData(),
		watchReactionsError: make(chan error),
		events:              make(chan Event),
		writeStreams:        make(map[*WriteStream]bool),
	}
	// Start watchReactions in a goroutine.
	watchClient, err := c.client.WatchReactions(ctx, &servicepb.WatchReactionsRequest{SessionId: s.id})
	if err != nil {
		return nil, fmt.Errorf("unable to watch for reactions: %w", err)
	}
	// This goroutine returns at the latest when ctx is canceled or the
	// session is closed by the server, whichever occurs earlier.
	go func() {
		// Errors from watchReactions are communicated to the user when
		// they call s.NextEvent.
		defer close(s.watchReactionsError)
		s.watchReactionsError <- watchReactions(ctx, watchClient, s.reactData, s.events)
	}()
	log.InfoContextf(ctx, "Session started with log context: %v", logContext)
	return s, nil
}

// End the session.
func (s *Session) End() error {
	// Close all streams
	for ws := range s.writeStreams {
		// Ignore any errors.
		ws.Close()
		delete(s.writeStreams, ws)
	}
	// Close the session.
	if err := s.session.CloseSend(); err != nil {
		return err
	}
	s.ended = true
	return nil
}

// MakeActionHandle returns a unique ActionHandle. The handle's ID will be
// unique among action handles created with this method with the same receiver
// s. IDs will be assigned sequentially starting from 1. This is safe to call
// concurrently.
func (s *Session) MakeActionHandle() ActionHandle {
	return ActionHandle{
		ActionID(s.actionIDs.Next()),
	}
}

// MakeEventSignal returns a new event marker, which may be used to identify
// that a reaction has occurred.
func (s *Session) MakeEventSignal() EventSignal {
	return EventSignal{
		session: s,
		unique:  s.eventSignalUniqueValue.Next(),
	}
}

// AddAction adds an action to the current session and returns ad.Handle.
func (s *Session) AddAction(ad *ActionDescription) (ActionHandle, error) {
	hs, err := s.AddActions(ad)
	if err != nil {
		return ActionHandle{}, err
	}
	return hs[0], nil
}

// AddActions adds a list of actions to the current session.
// This returns a list of handles to the added actions, in the same order as
// the actions appear in ads.
func (s *Session) AddActions(ads ...*ActionDescription) ([]ActionHandle, error) {
	if s.ended {
		return nil, ErrSessionEnded
	}
	actionProtos := make([]*typespb.ActionInstance, len(ads))
	reactionProtos := []*typespb.Reaction{}
	reactionDataEntries := []reactionDataEntry{}
	for i, ad := range ads {
		if ad == nil {
			return nil, fmt.Errorf("action description is nil: %w", ErrInvalidActionDescription)
		}
		if ad.Handle.IsZero() {
			return nil, fmt.Errorf("ad=%v: %w", ad, ErrInvalidActionHandle)
		}
		for _, r := range ad.Reactions {
			// Generate a unique reaction ID for this reaction.
			id := ReactionID(s.reactionIDs.Next())
			// Convert the reaction description to a proto.
			proto, err := r.proto(&ad.Handle, id)
			if err != nil {
				return nil, fmt.Errorf("ad=%v: %w", ad, err)
			}
			reactionProtos = append(reactionProtos, proto)
			// Determine mappings to event signals.
			ess, err := r.extractEventSignals()
			if err != nil {
				return nil, fmt.Errorf("ad=%v: %w", ad, err)
			}
			for _, es := range ess {
				reactionDataEntries = append(reactionDataEntries, reactionDataEntry{
					Action:  &ad.Handle,
					ReactID: id,
					Signal:  es,
				})
			}
		}
		a, err := ad.proto()
		if err != nil {
			return nil, fmt.Errorf("ad=%v: %w", ad, err)
		}
		actionProtos[i] = a
	}
	req := &servicepb.OpenSessionRequest{
		ActionRequest: &servicepb.OpenSessionRequest_AddActionsAndReactions{
			AddActionsAndReactions: &typespb.ActionsAndReactions{ActionInstances: actionProtos, Reactions: reactionProtos},
		},
	}
	if err := s.checkSendAndRecv(req); err != nil {
		return nil, fmt.Errorf("request to add actions failed: %w", err)
	}
	if err := s.reactData.Insert(reactionDataEntries...); err != nil {
		// The actions & reactions have been successfully created on the server,
		// but a client-side bookkeeping error occurred. This should never happen
		// in normal usage, but in case it does: close the session as a precaution.
		if endErr := s.End(); endErr != nil {
			return nil, internalError{err: fmt.Errorf("reaction bookkeeping error: %v (also: failed to end session as a precaution: %v)", err, endErr)}
		}
		return nil, internalError{err: fmt.Errorf("reaction bookkeeping error: %v (session has been ended as a precaution)", err)}
	}
	// Return action handles.
	hs := make([]ActionHandle, len(ads))
	for i, ad := range ads {
		hs[i] = ad.Handle
	}
	return hs, nil
}

// AddFreestandingReactions adds a list of free-standing reactions to the current session.
// Free-standing reactions are not associated with any action and are active as long as the session
// is active. Returns the error if any occurred.
func (s *Session) AddFreestandingReactions(reactions ...*Reaction) error {
	if s.ended {
		return ErrSessionEnded
	}
	var reactionProtos []*typespb.Reaction
	var reactionDataEntries []reactionDataEntry

	for _, r := range reactions {
		// Generate a unique reaction ID for this reaction.
		id := ReactionID(s.reactionIDs.Next())
		// Convert the reaction description to a proto.
		proto, err := r.proto(nil, id)
		if err != nil {
			return fmt.Errorf("reaction=%v: %w", r, err)
		}
		reactionProtos = append(reactionProtos, proto)
		// Determine mappings to event signals.
		ess, err := r.extractEventSignals()
		if err != nil {
			return fmt.Errorf("reaction=%v: %w", r, err)
		}
		for _, es := range ess {
			reactionDataEntries = append(reactionDataEntries, reactionDataEntry{
				Action:  nil,
				ReactID: id,
				Signal:  es,
			})
		}
	}
	req := &servicepb.OpenSessionRequest{
		ActionRequest: &servicepb.OpenSessionRequest_AddActionsAndReactions{
			AddActionsAndReactions: &typespb.ActionsAndReactions{Reactions: reactionProtos},
		},
	}
	if err := s.checkSendAndRecv(req); err != nil {
		return fmt.Errorf("request to add reactions failed: %w", err)
	}
	if err := s.reactData.Insert(reactionDataEntries...); err != nil {
		// The reactions have been successfully created on the server,
		// but a client-side bookkeeping error occurred. This should never happen
		// in normal usage, but in case it does: close the session as a precaution.
		if endErr := s.End(); endErr != nil {
			return internalError{err: fmt.Errorf("reaction bookkeeping error: %v (also: failed to end session as a precaution: %v)", err, endErr)}
		}
		return internalError{err: fmt.Errorf("reaction bookkeeping error: %v (session has been ended as a precaution)", err)}
	}
	return nil
}

// RemoveActions removes a list of Actions from the session.
func (s *Session) RemoveActions(ahs ...ActionHandle) error {
	if s.ended {
		return ErrSessionEnded
	}
	ids := make([]int64, len(ahs))
	for i, h := range ahs {
		if h.IsZero() {
			return ErrInvalidActionHandle
		}
		ids[i] = int64(h.ID())
	}
	req := &servicepb.OpenSessionRequest{
		ActionRequest: &servicepb.OpenSessionRequest_RemoveActionAndReactionIds{
			RemoveActionAndReactionIds: &typespb.ActionAndReactionIds{ActionInstanceIds: ids},
		},
	}
	s.reactData.RemoveActions(ahs...)
	return s.checkSendAndRecv(req)
}

// ClearAllActionsAndReactions removes all Actions and Reactions from the session.
func (s *Session) ClearAllActionsAndReactions() error {
	if s.ended {
		return ErrSessionEnded
	}
	req := &servicepb.OpenSessionRequest{
		ActionRequest: &servicepb.OpenSessionRequest_ClearAllActionsReactions{
			ClearAllActionsReactions: &servicepb.OpenSessionRequest_ClearAllActions{},
		},
	}
	s.reactData.Clear()
	return s.checkSendAndRecv(req)
}

// StartAction starts an Action and stops all active actions.
func (s *Session) StartAction(ah ActionHandle) error {
	if s.ended {
		return ErrSessionEnded
	}
	if ah.IsZero() {
		return ErrInvalidActionHandle
	}

	req := &servicepb.OpenSessionRequest{
		StartActionsRequest: &servicepb.OpenSessionRequest_StartActionsRequestData{
			ActionInstanceIds: []int64{int64(ah.ID())},
			StopActiveActions: true,
		},
	}
	return s.checkSendAndRecv(req)
}

// StartParallelAction starts an Action in parallel to active actions.
// It will preempt all active actions with an overlapping part set.
func (s *Session) StartParallelAction(ah ActionHandle) error {
	if s.ended {
		return ErrSessionEnded
	}
	if ah.IsZero() {
		return ErrInvalidActionHandle
	}

	req := &servicepb.OpenSessionRequest{
		StartActionsRequest: &servicepb.OpenSessionRequest_StartActionsRequestData{
			ActionInstanceIds: []int64{int64(ah.ID())},
			StopActiveActions: false,
		},
	}
	return s.checkSendAndRecv(req)
}

// StartActions starts multiple actions specified in `ahs`. Depending on `stop_active_actions`,
// all active actions remain active (if no new action has an overlapping part set) or will be stopped.
func (s *Session) StartActions(ahs []ActionHandle, stopActiveActions bool) error {
	if s.ended {
		return ErrSessionEnded
	}
	if len(ahs) == 0 {
		return ErrInvalidActionHandle
	}
	ids := make([]int64, len(ahs))
	for i, ah := range ahs {
		if ah.IsZero() {
			return ErrInvalidActionHandle
		}
		ids[i] = int64(ah.ID())
	}
	req := &servicepb.OpenSessionRequest{
		StartActionsRequest: &servicepb.OpenSessionRequest_StartActionsRequestData{
			ActionInstanceIds: ids,
			StopActiveActions: stopActiveActions,
		},
	}

	return s.checkSendAndRecv(req)
}

// OpenWriteStream opens a new stream to send messages to a running action.
func (s *Session) OpenWriteStream(ctx context.Context, ah ActionHandle, field string) (*WriteStream, error) {
	if s.ended {
		return nil, ErrSessionEnded
	}
	ctx = s.client.addOutgoingMetadata(ctx)
	sc, err := s.client.client.OpenWriteStream(ctx)
	if err != nil {
		return nil, err
	}
	req := &servicepb.OpenWriteStreamRequest{
		AddWriteStream: &servicepb.AddStreamRequest{ActionId: uint64(ah.ID()), FieldName: field},
		SessionId:      s.id,
	}
	if err = sc.Send(req); err != nil {
		return nil, err
	}
	resp, err := sc.Recv()
	if err != nil {
		return nil, err
	}
	if resp.GetAddStreamResponse().GetStatus().GetCode() != int32(codespb.OK) {
		return nil, fmt.Errorf("AddStream failed with code: %v", resp.GetAddStreamResponse().GetStatus().GetCode())
	}
	ws := &WriteStream{
		sessionID: s.id,
		action:    ah,
		client:    sc,
	}
	s.writeStreams[ws] = true
	return ws, nil
}

// OpenReadStream opens a new stream to receive messages from a running action.
func (s *Session) OpenReadStream(ah ActionHandle) (*ReadStream, error) {
	if s.ended {
		return nil, ErrSessionEnded
	}
	return &ReadStream{
		sessionID: s.id,
		action:    ah,
		client:    s.client,
	}, nil
}

// NextEvent returns the next event in the session's event queue. This blocks
// until either (i) an event is received, OR (ii) if cancel is non-nil, the
// cancel channel is written to or closed, OR (iii) the session ends.
func (s *Session) NextEvent(cancel <-chan struct{}) (Event, error) {
	if s.ended {
		return nil, ErrSessionEnded
	}
	select {
	case event, more := <-s.events:
		// Note: Immediately after the event stream is closed, an err
		// value is sent on watchReactionsError.
		if !more {
			err := <-s.watchReactionsError
			return nil, fmt.Errorf("the WatchReaction stream closed with error: %w", err)
		}
		return event, nil
	case <-cancel:
		return nil, ErrCanceled
	}
}

func (s *Session) checkSendAndRecv(req *servicepb.OpenSessionRequest) error {
	if err := s.session.Send(req); err != nil {
		return err
	}
	resp, err := s.session.Recv()
	if err != nil {
		return err
	}
	if resp.Status.Code != int32(codespb.OK) {
		return fmt.Errorf("session gRPC failed: %v", resp.Status.Message)
	}
	return nil
}
