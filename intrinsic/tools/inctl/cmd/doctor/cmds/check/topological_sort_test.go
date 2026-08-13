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

package check

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type TestEntity struct {
	Name string
	Deps []string
}

func newMapFromEntities(entities *[]TestEntity) *map[string]*TestEntity {
	newMap := make(map[string]*TestEntity)
	for _, entity := range *entities {
		newMap[entity.Name] = &entity
	}
	return &newMap
}

func TestTopologicalSort(t *testing.T) {
	testCases := []struct {
		name      string
		entities  *map[string]*TestEntity
		want      []string
		wantError error
	}{
		{
			name:      "empty",
			entities:  newMapFromEntities(&[]TestEntity{}),
			want:      []string(nil),
			wantError: nil,
		},
		{
			name: "no dependencies",
			entities: newMapFromEntities(&[]TestEntity{
				{Name: "a"},
				{Name: "b"},
				{Name: "c"},
			}),
			// The iteration order of the keys of an ordered map is not guaranteed, but we know that the
			// keys are sorted alphabetically by the topological sort function.
			want:      []string{"a", "b", "c"},
			wantError: nil,
		},
		{
			name: "dependencies",
			entities: newMapFromEntities(&[]TestEntity{
				{Name: "a", Deps: []string{"b"}},
				{Name: "b", Deps: []string{"c"}},
				{Name: "c"},
			}),
			want:      []string{"c", "b", "a"},
			wantError: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sorted, err := topologicalSort[TestEntity](tc.entities, func(entity *TestEntity) *[]string { return &entity.Deps })
			var got []string
			for _, entity := range sorted {
				got = append(got, entity.Name)
			}
			if err != nil {
				t.Errorf("topologicalSort(%v, <deps func>) returned an unexpected error: %v", tc.entities, err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("topologicalSort(%v, <deps func>) -> %v returned an unexpected result (-want +got):\n%s", tc.entities, got, diff)
			}
		})
	}
}

func TestTopologicalSort_Cycle(t *testing.T) {
	testCases := []struct {
		name     string
		entities *map[string]*TestEntity
	}{
		{
			name: "cycle",
			entities: newMapFromEntities(&[]TestEntity{
				{Name: "a", Deps: []string{"b"}},
				{Name: "b", Deps: []string{"c"}},
				{Name: "c", Deps: []string{"a"}},
			}),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sorted, err := topologicalSort[TestEntity](tc.entities, func(entity *TestEntity) *[]string { return &entity.Deps })
			if err == nil {
				t.Errorf("topologicalSort(%v, <deps func>) did not return an error: %v", tc.entities, sorted)
			}
		})
	}
}

func TestTopologicalSort_MissingDependency(t *testing.T) {
	testCases := []struct {
		name     string
		entities *map[string]*TestEntity
	}{
		{
			name: "missing dependency",
			entities: newMapFromEntities(&[]TestEntity{
				{Name: "a", Deps: []string{"b"}},
			}),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sorted, err := topologicalSort[TestEntity](tc.entities, func(entity *TestEntity) *[]string { return &entity.Deps })
			if err == nil {
				t.Errorf("topologicalSort(%v, <deps func>) did not return an error: %v", tc.entities, sorted)
			}
		})
	}
}
