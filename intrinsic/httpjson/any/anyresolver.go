// Copyright 2023 Intrinsic Innovation LLC

// Package anyresolver resolves Any messages on a best-effort basis given only the type_url.
package anyresolver

import (
	"context"
	"sync"
	"time"

	"intrinsic/httpjson/any/greedyresolver"
	"intrinsic/httpjson/any/installedassetsresolver"
	"intrinsic/httpjson/any/protoregistryresolver"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const (
	// Update cached data regularly for Solution Builders who are sideloading
	// assets with new versions of proto messages as they build a solution.
	cacheRefreshInterval = 30 * time.Second
)

// AnyResolver implements the MessageTypeResolver and ExtensionTypeResolver interfaces.
type AnyResolver struct {
	greedyResolver *greedyresolver.GreedyResolver
	iaResolver     *installedassetsresolver.InstalledAssetsResolver
	cancel         context.CancelFunc
	wg             sync.WaitGroup
}

// NewAnyResolver creates and returns a new instance of AnyResolver.
func NewAnyResolver(protoRegistryAddress string, installedAssetsAddress string) (*AnyResolver, error) {
	protoRegistryResolver, err := protoregistryresolver.NewProtoRegistryResolver(protoRegistryAddress)
	if err != nil {
		return nil, err
	}
	installedAssetsResolver, err := installedassetsresolver.NewInstalledAssetsResolver(installedAssetsAddress)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	r := &AnyResolver{
		greedyResolver: greedyresolver.NewGreedyResolver([]protoregistry.MessageTypeResolver{
			// Try the ProtoRegistry resolver first because it has a unique type url prefix,
			// and it returns NotFound immediately if the URL lacks that prefix.
			protoRegistryResolver,
			// Try InstalledAssets service next because many Any types in our APIs are
			// meant to hold parameters and configs for installed assets.
			installedAssetsResolver,
			// Last result, see if it is already in our global types
			protoregistry.GlobalTypes,
		}),
		iaResolver: installedAssetsResolver,
		cancel:     cancel,
	}

	r.wg.Add(1)
	go r.backgroundLoop(ctx)

	return r, nil
}

func (a *AnyResolver) Close() {
	a.cancel()
	a.wg.Wait()
}

func (a *AnyResolver) backgroundLoop(ctx context.Context) {
	defer a.wg.Done()
	ticker := time.NewTicker(cacheRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.iaResolver.RefreshInstalledAssets()
		}
	}
}

// FindExtensionByName looks up a extension field by the field's full name.
// Note that this is the full name of the field as determined by
// where the extension is declared and is unrelated to the full name of the
// message being extended.
//
// This returns (nil, NotFound) if not found.
func (a *AnyResolver) FindExtensionByName(field protoreflect.FullName) (protoreflect.ExtensionType, error) {
	// Nothing to add, but must implement it
	return protoregistry.GlobalTypes.FindExtensionByName(field)
}

// FindExtensionByNumber looks up a extension field by the field number
// within some parent message, identified by full name.
//
// This returns (nil, NotFound) if not found.
func (a *AnyResolver) FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error) {
	// Nothing to add, but must implement it
	return protoregistry.GlobalTypes.FindExtensionByNumber(message, field)
}

// FindMessageByName looks up a message by its full name.
// E.g., "google.protobuf.Any"
//
// This return (nil, NotFound) if not found.
func (a *AnyResolver) FindMessageByName(message protoreflect.FullName) (protoreflect.MessageType, error) {
	return a.greedyResolver.FindMessageByName(message)
}

// FindMessageByURL looks up a message by a URL identifier.
// See documentation on google.protobuf.Any.type_url for the URL format.
//
// This returns (nil, NotFound) if not found.
func (a *AnyResolver) FindMessageByURL(url string) (protoreflect.MessageType, error) {
	return a.greedyResolver.FindMessageByURL(url)
}
