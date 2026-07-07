// Copyright 2023 Intrinsic Innovation LLC

package typeurl

const (
	// DefaultPrefix is the default type URL prefix used by Protobuf.
	DefaultPrefix = "type.googleapis.com/"

	// IntrinsicPrefix is the prefix for Intrinsic type URLs which can be resolved
	// by the proto registry.
	IntrinsicPrefix = "type.intrinsic.ai/"

	// IntrinsicAreaSkills is the area (=the first path element) used in Intrinsic
	// type URLs for skills.
	IntrinsicAreaSkills = "skills"

	// IntrinsicAreaAssets is the area (=the first path element) used in Intrinsic
	// type URLs for assets.
	IntrinsicAreaAssets = "assets"

	// IntrinsicAreaCommon is the area (=the first path element) used in Intrinsic
	// type URLs for common types.
	IntrinsicAreaCommon = "common"

	// IntrinsicAreaWellKnown is the area (=the first path element) used in
	// Intrinsic type URLs for well known types.
	//
	// Deprecated: Use [IntrinsicAreaCommon] instead.
	IntrinsicAreaWellKnown = "well-known"

	// Separator is the top-level separator used in type URLs.
	Separator = "/"
)
