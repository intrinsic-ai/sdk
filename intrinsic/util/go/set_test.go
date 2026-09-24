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

package set

import (
	"cmp"
	"iter"
	"slices"
	"testing"

	gcmp "github.com/google/go-cmp/cmp"
)

func collect[T cmp.Ordered](seq iter.Seq[T]) []T {
	s := slices.Collect(seq)
	slices.Sort(s)
	return s
}

func TestSet(t *testing.T) {
	s := New[int]()
	if got := s.Size(); got != 0 {
		t.Errorf("New().Size() = %d, want 0", got)
	}

	s.Add(1)
	s.Add(2)
	s.Add(1)

	if got := s.Size(); got != 2 {
		t.Errorf("Size() = %d, want 2", got)
	}
	if !s.Contains(1) {
		t.Errorf("Contains(1) = false, want true")
	}
	if !s.Contains(2) {
		t.Errorf("Contains(2) = false, want true")
	}
	if s.Contains(3) {
		t.Errorf("Contains(3) = true, want false")
	}

	s.Remove(1)
	if s.Size() != 1 {
		t.Errorf("Size() = 1, want 1")
	}
	if s.Contains(1) {
		t.Errorf("Contains(1) = true, want false")
	}

	s.AddAll(3, 4, 3)
	if got := s.Size(); got != 3 {
		t.Errorf("Size() = %d, want 3", got)
	}
	got := collect(s.All())
	want := []int{2, 3, 4}
	if diff := gcmp.Diff(want, got); diff != "" {
		t.Errorf("All() returned diff (-want +got):\n%s", diff)
	}
}

func TestNewWithValues(t *testing.T) {
	s := New(1, 2, 1)
	if got := s.Size(); got != 2 {
		t.Errorf("New(1, 2, 1).Size() = %d, want 2", got)
	}
	if !s.Contains(1) {
		t.Errorf("New(1, 2, 1).Contains(1) = false, want true")
	}
	if !s.Contains(2) {
		t.Errorf("New(1, 2, 1).Contains(2) = false, want true")
	}
}

func TestEquals(t *testing.T) {
	tests := []struct {
		name  string
		left  Set[int]
		right Set[int]
		want  bool
	}{
		{
			name:  "empty equals empty",
			left:  New[int](),
			right: New[int](),
			want:  true,
		},
		{
			name:  "empty not equals non-empty",
			left:  New[int](),
			right: New[int](1),
			want:  false,
		},
		{
			name:  "same sets are equal",
			left:  New[int](1, 2, 3),
			right: New[int](1, 2, 3),
			want:  true,
		},
		{
			name:  "same sets with different order are equal",
			left:  New[int](1, 2, 3),
			right: New[int](3, 2, 1),
			want:  true,
		},
		{
			name:  "different sets are not equal",
			left:  New[int](1, 2, 3),
			right: New[int](1, 2, 4),
			want:  false,
		},
		{
			name:  "different size sets are not equal",
			left:  New[int](1, 2, 3),
			right: New[int](1, 2),
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.left.Equals(tc.right); got != tc.want {
				t.Errorf("(%v).Equals(%v) = %t, want %t", tc.left, tc.right, got, tc.want)
			}
		})
	}
}

func TestUnion(t *testing.T) {
	tests := []struct {
		name  string
		s1    []int
		s2    []int
		wantR []int
	}{
		{
			name: "both empty",
		},
		{
			name:  "s1 empty",
			s2:    []int{1, 2},
			wantR: []int{1, 2},
		},
		{
			name:  "s2 empty",
			s1:    []int{1, 2},
			wantR: []int{1, 2},
		},
		{
			name:  "no overlap",
			s1:    []int{1, 2},
			s2:    []int{3, 4},
			wantR: []int{1, 2, 3, 4},
		},
		{
			name:  "overlap",
			s1:    []int{1, 2, 3},
			s2:    []int{3, 4, 5},
			wantR: []int{1, 2, 3, 4, 5},
		},
		{
			name:  "s1 subset of s2",
			s1:    []int{1, 2},
			s2:    []int{1, 2, 3},
			wantR: []int{1, 2, 3},
		},
		{
			name:  "s2 subset of s1",
			s1:    []int{1, 2, 3},
			s2:    []int{1, 2},
			wantR: []int{1, 2, 3},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s1 := New(tc.s1...)
			s2 := New(tc.s2...)
			got := collect(s1.Union(s2))
			want := tc.wantR
			if diff := gcmp.Diff(want, got); diff != "" {
				t.Errorf("Union() returned diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	tests := []struct {
		name  string
		s1    []int
		s2    []int
		wantR []int
	}{
		{
			name: "both empty",
		},
		{
			name: "s1 empty",
			s2:   []int{1, 2},
		},
		{
			name: "s2 empty",
			s1:   []int{1, 2},
		},
		{
			name: "no overlap",
			s1:   []int{1, 2},
			s2:   []int{3, 4},
		},
		{
			name:  "overlap",
			s1:    []int{1, 2, 3},
			s2:    []int{3, 4, 5},
			wantR: []int{3},
		},
		{
			name:  "s1 subset of s2",
			s1:    []int{1, 2},
			s2:    []int{1, 2, 3},
			wantR: []int{1, 2},
		},
		{
			name:  "s2 subset of s1",
			s1:    []int{1, 2, 3},
			s2:    []int{1, 2},
			wantR: []int{1, 2},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s1 := New(tc.s1...)
			s2 := New(tc.s2...)
			got := collect(s1.Intersection(s2))
			want := tc.wantR
			if diff := gcmp.Diff(want, got); diff != "" {
				t.Errorf("Intersection() returned diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDifference(t *testing.T) {
	tests := []struct {
		name  string
		s1    []int
		s2    []int
		wantR []int
	}{
		{
			name: "both empty",
		},
		{
			name: "s1 empty",
			s2:   []int{1, 2},
		},
		{
			name:  "s2 empty",
			s1:    []int{1, 2},
			wantR: []int{1, 2},
		},
		{
			name:  "no overlap",
			s1:    []int{1, 2},
			s2:    []int{3, 4},
			wantR: []int{1, 2},
		},
		{
			name:  "overlap",
			s1:    []int{1, 2, 3},
			s2:    []int{3, 4, 5},
			wantR: []int{1, 2},
		},
		{
			name: "s1 subset of s2",
			s1:   []int{1, 2},
			s2:   []int{1, 2, 3},
		},
		{
			name:  "s2 subset of s1",
			s1:    []int{1, 2, 3},
			s2:    []int{1, 2},
			wantR: []int{3},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s1 := New(tc.s1...)
			s2 := New(tc.s2...)
			got := collect(s1.Difference(s2))
			want := tc.wantR
			if diff := gcmp.Diff(want, got); diff != "" {
				t.Errorf("Difference() returned diff (-want +got):\n%s", diff)
			}
		})
	}
}
