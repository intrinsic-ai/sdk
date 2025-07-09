// Copyright 2023 Intrinsic Innovation LLC

// Package bundleio contains a function that reads a bundle archive file.
package bundleio

import (
	"errors"
	"fmt"
	"io"
	"os"

	"archive/tar"
	"google.golang.org/protobuf/proto"

	idpb "intrinsic/assets/proto/id_go_proto"
	ipb "intrinsic/kubernetes/workcell_spec/proto/image_go_proto"
)

type handler func(io.Reader) error
type fallbackHandler func(string, io.Reader) error

// ImageProcessor is a closure that pushes an image and returns the resulting
// pointer to the container registry.  It is provided the id of the bundle being
// processed as well as the name of the specific image.  It is expected to
// upload the image and produce a usable image spec.  The reader points to an
// image archive.  This may be invoked multiple times.  Images are ignored if it
// is not specified.
type ImageProcessor func(idProto *idpb.Id, filename string, r io.Reader) (*ipb.Image, error)

// walkTarFile walks through a tar file and invokes handlers on specific
// filenames.  fallback can be nil.  Returns an error if all handlers in
// handlers are not invoked.  It ignores all non-regular files.
func walkTarFile(t *tar.Reader, handlers map[string]handler, fallback fallbackHandler) error {
	for len(handlers) > 0 || fallback != nil {
		hdr, err := t.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("getting next file failed: %v", err)
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}

		n := hdr.Name
		if h, ok := handlers[n]; ok {
			delete(handlers, n)
			if err := h(t); err != nil {
				return fmt.Errorf("error processing file %q: %v", n, err)
			}
		} else if fallback != nil {
			if err := fallback(n, t); err != nil {
				return fmt.Errorf("error processing file %q: %v", n, err)
			}
		}
	}
	if len(handlers) != 0 {
		keys := make([]string, 0, len(handlers))
		for k := range handlers {
			keys = append(keys, k)
		}
		return fmt.Errorf("missing expected files %s", keys)
	}
	return nil
}

// ignoreHandler is a function that can be used as a handler to ignore specific
// files.
func ignoreHandler(r io.Reader) error {
	return nil
}

// alwaysErrorAsUnexpected can be used as a fallback handler that will always
// trigger an unexpected file error.  This forces all files to be handled
// explicitly.
func alwaysErrorAsUnexpected(n string, r io.Reader) error {
	return fmt.Errorf("unexpected file %q", n)
}

// makeBinaryProtoHandler creates a handler that reads a binary proto file and
// unmarshals it into a file.  The proto must not be nil.
func makeBinaryProtoHandler(p proto.Message) handler {
	return func(r io.Reader) error {
		b, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("error reading: %v", err)
		}
		if err := proto.Unmarshal(b, p); err != nil {
			return fmt.Errorf("error parsing proto: %v", err)
		}
		return nil
	}
}

// makeCollectInlinedFallbackHandler constructs a default handler that collects
// all of the unknown files and reads their bytes into a map.  The key of the
// map is the filename, and the value is the file contents.
func makeCollectInlinedFallbackHandler() (map[string][]byte, fallbackHandler) {
	inlined := map[string][]byte{}
	fallback := func(n string, r io.Reader) error {
		b, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("error reading: %v", err)
		}
		inlined[n] = b
		return nil
	}
	return inlined, fallback
}

// readBinaryProto reads a binary proto from a reader and unmarshals it into a proto.
func readBinaryProto(r io.Reader, p proto.Message) error {
	if b, err := io.ReadAll(r); err != nil {
		return fmt.Errorf("error reading: %v", err)
	} else if err := proto.Unmarshal(b, p); err != nil {
		return fmt.Errorf("error parsing proto: %v", err)
	}

	return nil
}

// BundleType is used to return the type of a bundle file.
type BundleType int

// The different bundle types that can be detected from a file.
const (
	BundleTypeSkill BundleType = iota
	BundleTypeService
)

var (
	errNoValidTypeDetected   = errors.New("no recognized manifest detected")
	errMultipleTypesDetected = errors.New("unsupported bundle")
)

// DetectBundleType will return the type of bundle a file represents.  It does
// not do any validation of the particular file, just provides an indication
// what sort of processing should be done on the file.
func DetectBundleType(path string) (BundleType, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("could not open %q: %v", path, err)
	}
	defer f.Close()

	lookup := map[string]BundleType{
		serviceManifestPathInTar:       BundleTypeService,
		skillManifestPathInTar:         BundleTypeSkill,
	}

	var bt BundleType
	var found int
	if err := walkTarFile(tar.NewReader(f), map[string]handler{}, func(path string, _ io.Reader) error {
		if val, ok := lookup[path]; ok {
			found++
			bt = val
		}
		return nil
	}); err != nil {
		return bt, err
	}
	switch found {
	case 0:
		return 0, errNoValidTypeDetected
	case 1:
		return bt, nil
	default:
		return 0, errMultipleTypesDetected
	}
}
