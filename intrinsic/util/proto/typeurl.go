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

	// Separator is the top-level separator used in type URLs.
	Separator = "/"
)
