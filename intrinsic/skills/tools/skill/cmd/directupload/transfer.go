// Copyright 2023 Intrinsic Innovation LLC

package directupload

import (
	"context"
	"fmt"
	"io"

	backoff "github.com/cenkalti/backoff/v4"
	log "github.com/golang/glog"
	"github.com/google/go-containerregistry/pkg/name"
	crv1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"intrinsic/assets/imagetransfer"
	"intrinsic/storage/artifacts/client/client"

	artifactgrpcpb "intrinsic/storage/artifacts/proto/v1/artifact_go_grpc_proto"
)

// Option allows setting direct upload transferer options.
type Option func(transfer *directTransfer)

// WithMaxRetries allows setting max retries for the upload
func WithMaxRetries(maxRetries int) Option {
	return func(transfer *directTransfer) {
		transfer.maxRetries = maxRetries
	}
}

// WithClient allows caller to set client side implementation. If this option
// is specified, the client will be used to create an uploader instance,
// ignoring discovery strategy set by WithDiscovery
func WithClient(client artifactgrpcpb.ArtifactServiceApiClient) Option {
	return func(transfer *directTransfer) {
		transfer.client = client
	}
}

// WithOutput allows adding simple progress monitor with w as its output.
func WithOutput(w io.Writer) Option {
	return func(transfer *directTransfer) {
		transfer.ctx = client.SetProgressMonitor(transfer.ctx, newMonitor(w))
	}
}

// WithDiscovery allows setting a TargetDiscovery implementation to discover
// the most suitable client path. One of WithClient or WithDiscovery have to be
// used in order to specify upload target.
func WithDiscovery(discovery TargetDiscovery) Option {
	return func(transfer *directTransfer) {
		transfer.discovery = discovery
	}
}

// WithFailOver allows to set fail-over transferer in case direct upload
// is not possible.
func WithFailOver(failOver imagetransfer.Transferer) Option {
	return func(transfer *directTransfer) {
		transfer.failOver = failOver
	}
}

// NewTransferer create a new instance of direct upload Transferer implementation
// and applies options if specified.
func NewTransferer(ctx context.Context, opts ...Option) imagetransfer.Transferer {
	transfer := &directTransfer{
		maxRetries: 5,
		ctx:        ctx,
	}

	for _, opt := range opts {
		opt(transfer)
	}

	if transfer.client == nil {
		if transfer.discovery == nil {
			// this is programmer error...
			panic("cannot obtain client, use WithDiscovery or WithClient options")
		}
	}

	return transfer
}

type directTransfer struct {
	maxRetries int
	failOver   imagetransfer.Transferer
	uploader   client.Uploader
	client     artifactgrpcpb.ArtifactServiceApiClient
	ctx        context.Context
	discovery  TargetDiscovery
}

func (dt *directTransfer) Write(ref name.Reference, img crv1.Image) error {

	if dt.uploader == nil {
		apiClient, err := dt.getClient()
		if err != nil {
			return fmt.Errorf("cannot connect: %w", err)
		}
		dt.uploader, err = client.NewUploader(apiClient, client.WithSequentialUpload(),
			// To mitigate b/330747118; this is not full fix, but should help.
			// This setting is forcing to run only 1 upload task at a time.
			// We are taking significant performance penalty.
			client.WithUploadParallelism(1))
		if err != nil {
			return fmt.Errorf("cannot create uploader: %w", err)
		}
	}

	numAttempts := atomic.NewUint32(0)
	// The initial attempt is not counted as a retry.
	maxAttempts := 1 + dt.maxRetries
	err := backoff.Retry(func() error {
		if dt.ctx.Err() != nil {
			return backoff.Permanent(dt.ctx.Err())
		}
		attempt := numAttempts.Inc()
		err := dt.uploader.UploadImage(dt.ctx, ref.String(), img)
		if err != nil {
			log.Errorf("attempt %d/%d: failed to upload image (%s): %s", attempt, maxAttempts, ref, err)
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return backoff.Permanent(err)
			}
			// todo: evaluate other permanent errors, such as 500
		}
		return err
	}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), uint64(dt.maxRetries)))
	if err != nil {
		if dt.failOver != nil {
			if foErr := dt.failOver.Write(ref, img); foErr != nil {
				return fmt.Errorf("image write failed (direct: %s): %w", err, foErr)
			}
			log.Warningf("fail over succeeded with prior direct upload failure: %s", err)
			return nil
		}
		return fmt.Errorf("image write failed: %w", err)
	}
	return nil
}

func (dt *directTransfer) getClient() (artifactgrpcpb.ArtifactServiceApiClient, error) {
	if dt.client != nil {
		return dt.client, nil
	}

	apiClient, err := dt.discovery.GetClient(dt.ctx)
	if err != nil {
		return nil, err
	}

	dt.client = apiClient
	return dt.client, nil
}
