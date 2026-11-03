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

package workcellinfo

import (
	"testing"

	"google.golang.org/protobuf/proto"

	workcellinfopb "intrinsic/platform/common/proto/workcell_info_go_proto"
	"intrinsic/platform/pubsub/golang/fakekvstore"
)

const (
	testWorkcell      = "test_workcell"
	testCloudEndpoint = "test_cloud_endpoint"
	testProject       = "test_project"
)

func TestSetGet(t *testing.T) {
	kv := fakekvstore.New()

	originalWorkcellInfo := &workcellinfopb.WorkcellInfo{
		WorkcellName:  testWorkcell,
		CloudEndpoint: testCloudEndpoint,
		Project:       testProject,
	}
	if err := doSet(kv, originalWorkcellInfo); err != nil {
		t.Fatalf("Failed to set workcell info: %v", err)
	}

	workcellInfo, err := doGet(kv)
	if err != nil {
		t.Fatalf("Failed to get workcell info: %v", err)
	}
	if !proto.Equal(workcellInfo, originalWorkcellInfo) {
		t.Errorf("Read unexpected value. Got %v, want %v", workcellInfo, originalWorkcellInfo)
	}
}

type WorkcellInfoCRUDOperationResult struct {
	op           string
	workcellInfo *workcellinfopb.WorkcellInfo
	err          error
}

func TestGetBlocksUntilValueIsSet(t *testing.T) {
	kv := fakekvstore.New()
	ch := make(chan WorkcellInfoCRUDOperationResult, 1)

	go func() {
		originalWorkcellInfo := &workcellinfopb.WorkcellInfo{
			WorkcellName:  testWorkcell,
			CloudEndpoint: testCloudEndpoint,
			Project:       testProject,
		}
		err := doSet(kv, originalWorkcellInfo)
		ch <- WorkcellInfoCRUDOperationResult{
			op:           "set",
			workcellInfo: originalWorkcellInfo,
			err:          err,
		}
	}()

	go func() {
		workcellInfo, err := doGet(kv)
		ch <- WorkcellInfoCRUDOperationResult{
			op:           "get",
			workcellInfo: workcellInfo,
			err:          err,
		}
		close(ch)
	}()

	results := []WorkcellInfoCRUDOperationResult{}
	for result := range ch {
		results = append(results, result)
	}

	if len(results) != 2 {
		t.Fatalf("Got %d results, want 2", len(results))
	}
	if results[0].op != "set" || results[1].op != "get" {
		t.Fatalf(
			"Wrong order of operations. Got [%v, %v], want [set, get]",
			results[0].op,
			results[1].op)
	}
	if !proto.Equal(results[0].workcellInfo, results[1].workcellInfo) {
		t.Errorf("Read unexpected value. Got %v, want %v", results[1].workcellInfo, results[0].workcellInfo)
	}
}
