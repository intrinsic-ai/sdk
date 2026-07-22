// Copyright 2023 Intrinsic Innovation LLC

package icon

import (
	"context"
	"errors"
	"fmt"
	"time"

	codespb "google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"

	grpcpb "intrinsic/icon/proto/v1/service_go_proto"
	servicepb "intrinsic/icon/proto/v1/service_go_proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
)

// ErrStreamClosed occurs when a write is attempted on a closed stream.
var ErrStreamClosed = errors.New("stream already closed")

// WriteStream wraps IconApi_OpenWriteStreamClient to be used for streaming
// messages to the server.  It also keeps track of what session and action it
// was created for.
type WriteStream struct {
	sessionID int64
	action    ActionHandle
	client    grpcpb.IconApi_OpenWriteStreamClient
	closed    bool
}

// Write sends proto message to the server and checks the response
func (w *WriteStream) Write(msg proto.Message) error {
	if w.closed {
		return ErrStreamClosed
	}
	a, err := anypb.New(msg)
	if err != nil {
		return err
	}
	req := &servicepb.OpenWriteStreamRequest{
		WriteValue: &servicepb.OpenWriteStreamRequest_WriteValue{
			Value: a,
		},
	}
	if err := w.client.Send(req); err != nil {
		return err
	}
	resp, err := w.client.Recv()
	if err != nil {
		return err
	}
	if resp.WriteValueResponse.Code != int32(codespb.OK) {
		return fmt.Errorf("stream write failed: %v", resp.WriteValueResponse.Message)
	}
	return nil
}

// Close closes the stream and prevents any more messages from being sent.
func (w *WriteStream) Close() error {
	if w.closed {
		return nil
	}
	err := w.client.CloseSend()
	if err != nil {
		return err
	}
	w.closed = true
	return nil
}

// ReadStream is created for a specific session and action belonging to the
// session and exposes methods forgetting the latest output from a running
// action.
type ReadStream struct {
	sessionID int64
	action    ActionHandle
	client    *grpcClient
}

// Read gets the timestamp (representing the time since the server started)
// and payload for the most recent output from a running action.
func (r *ReadStream) Read(ctx context.Context) (time.Duration, *anypb.Any, error) {
	ctx = r.client.addOutgoingMetadata(ctx)
	req := &servicepb.GetLatestStreamingOutputRequest{
		SessionId: r.sessionID,
		ActionId:  uint64(r.action.ID()),
	}
	resp, err := r.client.client.GetLatestStreamingOutput(ctx, req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read stream: %v", err)
	}
	return time.Duration(resp.Output.TimestampNs), resp.Output.Payload, nil
}

// ReadUnpacked accepts an empty proto message that receives the unmarshalled
// payload. It returns a timestamp (representing the time since the server
// started).
func (r *ReadStream) ReadUnpacked(ctx context.Context, dest proto.Message) (time.Duration, error) {
	t, any, err := r.Read(ctx)
	if err != nil {
		return 0, err
	}
	if err := anypb.UnmarshalTo(any, dest, proto.UnmarshalOptions{}); err != nil {
		return 0, err
	}
	return t, nil
}
