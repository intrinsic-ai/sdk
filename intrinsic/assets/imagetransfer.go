// Copyright 2023 Intrinsic Innovation LLC

// Package imagetransfer contains a Transferer interface and implementation
// that reads and writes images to a container registry.
package imagetransfer

import (
	"fmt"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/google/go-containerregistry/pkg/name"
	containerregistry "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/pkg/errors"
)

// Number of times to try uploading a container image if we get retriable errors.
const remoteWriteTries = 5

// Transferer provides methods to read and write images to a container registry.
type Transferer interface {
	Write(ref name.Reference, img containerregistry.Image) error
}

type remoteImage struct {
	Opts []remote.Option
}

// Write pushes an image to a container registry.
func (r remoteImage) Write(ref name.Reference, img containerregistry.Image) error {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), remoteWriteTries)
	if err := backoff.Retry(func() error {
		err := remote.Write(ref, img, r.Opts...)
		if err, ok := err.(*transport.Error); ok && err.StatusCode >= 500 {
			// Retry server errors like 504 Gateway Timeout.
			return err
		}
		if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}, b); err != nil {
		return errors.Wrapf(err, "remote.Write to %q", ref)
	}
	return nil
}

// RemoteTransferer returns a new Transferer using the passed-in options.
func RemoteTransferer(opts ...remote.Option) Transferer {
	return remoteImage{
		Opts: opts,
	}
}

// NoOpTransferer errors if any attempt is made to read or write an image.
type NoOpTransferer struct{}

func (NoOpTransferer) Write(ref name.Reference, img containerregistry.Image) error {
	return fmt.Errorf("NoOpTransferer forbids writing an image")
}
