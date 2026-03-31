// Copyright 2023 Intrinsic Innovation LLC

// Package recordingdatacategories provides embedded text protobuf data
// for the recording data categories single source of truth.
package recordingdatacategories

import (
	_ "embed"
	"sync"

	recordingcategoriespb "intrinsic/logging/proto/recording_data_categories_go_proto"

	"google.golang.org/protobuf/encoding/prototext"
)

//go:embed recording_data_categories.pbtxt
var recordingDataCategoriesBytes []byte

var (
	recordingDataCategories *recordingcategoriespb.RecordingDataCategories
	once                    sync.Once
)

// GetRecordingDataCategories returns the unmarshaled RecordingDataCategories.
// The result is cached after the first call.
func GetRecordingDataCategories() *recordingcategoriespb.RecordingDataCategories {
	once.Do(func() {
		recordingDataCategories = &recordingcategoriespb.RecordingDataCategories{}
		if err := prototext.Unmarshal(recordingDataCategoriesBytes, recordingDataCategories); err != nil {
			panic("failed to unmarshal recording data categories: " + err.Error())
		}
	})
	return recordingDataCategories
}
