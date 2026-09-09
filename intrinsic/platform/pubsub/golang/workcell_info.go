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
	"fmt"
	"time"

	"intrinsic/platform/pubsub/golang/kvstore"
	"intrinsic/platform/pubsub/golang/pubsub"

	log "github.com/golang/glog"

	workcellinfopb "intrinsic/platform/common/proto/workcell_info_go_proto"
)

const (
	workcellInfoKey = "workcell_info"
)

// Set stores the given workcell info in the workcell's local KV store.
func Set(workcellInfo *workcellinfopb.WorkcellInfo) error {
	ps, err := pubsub.NewPubSub()
	if err != nil {
		return fmt.Errorf("failed to create PubSub: %w", err)
	}
	defer ps.Close()
	kv := ps.KVStore()

	return doSet(kv, workcellInfo)
}

// doSet stores the workcell info in the given KV store.
// Can be called from unit tests that use the fake KV store.
func doSet(kv kvstore.KVStore, workcellInfo *workcellinfopb.WorkcellInfo) error {
	for {
		if err := kv.Set(workcellInfoKey, workcellInfo, true); err != nil {
			log.Errorf("Error setting workcell info in kv store: %v", err)
			continue
		}
		log.Info("Successfully set workcell info in kv store")
		return nil
	}
}

// Get reads workcell info from the workcell's local KV store.
// It blocks the calling thread until the workcell info is read.
// The calling thread may be blocked for a long time if the workcell info
// hasn't been written yet. It may happen, for example, when intrinsic-base
// has just been installed, and the container that stores the workcell info
// is still starting.
func Get() (*workcellinfopb.WorkcellInfo, error) {
	ps, err := pubsub.NewPubSub()
	if err != nil {
		return nil, fmt.Errorf("failed to create PubSub: %w", err)
	}
	defer ps.Close()
	kvstore := ps.KVStore()

	return doGet(kvstore)
}

// doGet reads workcell info from the given KV store.
// Can be called from unit tests that use the fake KV store.
func doGet(kv kvstore.KVStore) (*workcellinfopb.WorkcellInfo, error) {
	log.Info("Reading workcell info from the KV store")
	var workcellInfo = &workcellinfopb.WorkcellInfo{}
	for {
		valueAny, err := kv.Get(workcellInfoKey, nil /*timeout*/)
		if err != nil {
			log.Infof("Failed to read info, will retry: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		err = valueAny.UnmarshalTo(workcellInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to parse workcell info: %w", err)
		}

		log.Infof("Got workcell info: %v", workcellInfo)
		return workcellInfo, nil
	}
}
