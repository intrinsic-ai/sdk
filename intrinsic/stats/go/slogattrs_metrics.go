// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
