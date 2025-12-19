// Copyright 2023 Intrinsic Innovation LLC

// Package bundleimages has utilities to push images from a resource bundle.
package bundleimages

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"intrinsic/assets/idutils"
	"intrinsic/assets/imageutils"
	"intrinsic/assets/services/readeropener"
	"intrinsic/kubernetes/workcell_spec/imagetags"

	crname "github.com/google/go-containerregistry/pkg/name"
	containerregistry "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/pkg/errors"

	idpb "intrinsic/assets/proto/id_go_proto"
	ipb "intrinsic/kubernetes/workcell_spec/proto/image_go_proto"
)

const (
	// maxInMemorySizeForPushArchive is set to a conservative 100MB for now.
	// Consider raising this value in the future if needed.
	maxInMemorySizeForPushArchive = 100 * 1024 * 1024
)

// CreateImageProcessor returns a closure to handle images within a bundle.  It
// pushes images to the registry using a default tag.  The image is named with
// the id of the resource with the basename image filename appended.
func CreateImageProcessor(reg RegistryOptions) imageutils.ImageProcessor {
	return func(ctx context.Context, idProto *idpb.Id, filename string, r io.Reader) (*ipb.Image, error) {
		id, err := idutils.IDFromProto(idProto)
		if err != nil {
			return nil, fmt.Errorf("unable to get tag for image: %v", err)
		}

		fileNoExt := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
		name := fmt.Sprintf("%s.%s", id, fileNoExt)

		// Some images can be quite large (>1GB) and cause out-of-memory issues when
		// read into a byte buffer. We use the readeropener utility to use an
		// in-memory buffer when the size is small and to write the contents to disk
		// when large. Note that having some buffer is necessary as PushArchive will
		// attempt to read the buffer more than once and tar files don't have a way
		// to seek backwards (tape only ran one direction after all).
		opener, cleanup, err := readeropener.New(r, maxInMemorySizeForPushArchive)
		if err != nil {
			return nil, fmt.Errorf("could not process tar file %q: %v", filename, err)
		}
		defer cleanup()
		img, err := tarball.Image(tarball.Opener(opener), nil)
		if err != nil {
			return nil, fmt.Errorf("could not create tarball image: %v", err)
		}
		return pushImage(ctx, img, name, reg)
	}
}

// writer is the interface required to push an Image to a particular reference.
type writer interface {
	Write(context.Context, crname.Reference, containerregistry.Image) error
}

// BasicAuth provides the necessary fields to perform basic authentication with
// a resource registry.
type BasicAuth struct {
	// User is the username used to access the registry.
	User string
	// Pwd is the password used to authenticate registry access.
	Pwd string
}

// RegistryOptions is used to configure Push to a specific registry
type RegistryOptions struct {
	// URI of the container registry
	URI string
	// The transferer performs the work to write the container to the registry.
	Transferer writer
	// The optional parameters required to perform basic authentication with
	// the registry.
	BasicAuth
}

// pushImage takes an image and pushes it to the specified registry with the
// given options.
func pushImage(ctx context.Context, img containerregistry.Image, name string, reg RegistryOptions) (*ipb.Image, error) {
	registry := strings.TrimSuffix(reg.URI, "/")
	if len(registry) == 0 {
		return nil, fmt.Errorf("registry is empty")
	}
	// Use the rapid candidate name if provided or a placeholder tag otherwise.
	// For Rapid workflows, the deployed chart references the image by
	// candidate name. For dev workflows, we reference by digest.
	tag, err := imagetags.DefaultTag()
	if err != nil {
		return nil, errors.Wrap(err, "generating tag")
	}

	// A tag is required for retention.  Infra uses an img being untagged as
	// a signal it can be removed.
	dst := fmt.Sprintf("%s/%s:%s", registry, name, tag)
	ref, err := crname.NewTag(dst)
	if err != nil {
		return nil, errors.Wrapf(err, "name.NewReference(%q)", dst)
	}

	digest, err := img.Digest()
	if err != nil {
		return nil, fmt.Errorf("could not get the sha256 of the image: %v", err)
	}

	if err := reg.Transferer.Write(ctx, ref, img); err != nil {
		return nil, fmt.Errorf("could not write image %q: %v", dst, err)
	}

	// Always provide a spec in terms of the digest, since that is
	// reproducible, while a tag may not be.
	return &ipb.Image{
		Registry:     registry,
		Name:         name,
		Tag:          "@" + digest.String(),
		AuthUser:     reg.User,
		AuthPassword: reg.Pwd,
	}, nil
}
