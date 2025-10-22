// Copyright 2023 Intrinsic Innovation LLC

package slogattrs

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	// TagLogLevel is the tag key for the log level.
	TagLogLevel = tag.MustNewKey("log_level")

	// MLogCount is the metric for the count of logs. Log levels are added as tags.
	MLogCount = stats.Int64("log_count", "Count of logs", stats.UnitDimensionless)

	logCountView = view.View{
		Name:        "slogattrs/log_count",
		Measure:     MLogCount,
		Description: "Count of logs, broken down by log level",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{TagLogLevel},
	}

	// Views is the list of views for slogattrs.
	Views = []*view.View{&logCountView}
)
