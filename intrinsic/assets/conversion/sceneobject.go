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

// Package sceneobject provides conversion utilities for SceneObject Assets.
package sceneobject

import (
	"fmt"

	sompb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"
)

// ManifestFromAsset extracts from the given ProcessedSceneObjectManifest a SceneObjectManifest
// which is suitable for creating a corresponding bundle from it.
func ManifestFromAsset(pm *sompb.ProcessedSceneObjectManifest) (*sompb.SceneObjectManifest, error) {
	if pm == nil {
		return nil, fmt.Errorf("ProcessedSceneObjectManifest must not be nil")
	}

	return &sompb.SceneObjectManifest{
		Metadata: pm.GetMetadata(),
	}, nil
}
