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

package icon

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"

	typespb "intrinsic/icon/proto/v1/types_go_proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
)

// ErrNotFound is used when a requested object is not found.
var ErrNotFound = errors.New("not found")

// ServerConfig holds the static configuration details for each part.
type ServerConfig struct {
	partConfigs map[string]*typespb.PartConfig
}

// ServerConfigFromProto creates a new ServerConfig from proto.
func ServerConfigFromProto(proto []*typespb.PartConfig) (*ServerConfig, error) {
	m := make(map[string]*typespb.PartConfig)
	for _, p := range proto {
		if _, exists := m[p.Name]; exists {
			return nil, internalError{err: fmt.Errorf("part %q listed twice; with two configs: %v AND %v", p.Name, m[p.Name], p.Config)}
		}
		m[p.Name] = p
	}
	return &ServerConfig{partConfigs: m}, nil
}

// PartConfigs returns a map containing the configuration for each part, keyed
// by part name. The configurations are returned as typespb.PartConfig objects.
func (s *ServerConfig) PartConfigs() map[string]*typespb.PartConfig {
	// Return a deep copy of the map to prevent the user from modifying
	// s.partConfigs.
	m := make(map[string]*typespb.PartConfig)
	for k, v := range s.partConfigs {
		m[k] = proto.Clone(v).(*typespb.PartConfig)
	}
	return m
}

// PartConfig returns the configuration for part, as an typespb.PartConfig.
func (s *ServerConfig) PartConfig(part string) (*typespb.PartConfig, error) {
	pc, exists := s.partConfigs[part]
	if !exists {
		return nil, fmt.Errorf("part %q: %w", part, ErrNotFound)
	}
	return proto.Clone(pc).(*typespb.PartConfig), nil
}

// PartConfigUnpacked obtains the configuration for part, unpacking it into dest.
// Dest must have a type matching the part-specific configuration type, as
// reported in the part's signature.
func (s *ServerConfig) PartConfigUnpacked(part string, dest proto.Message) error {
	pc, exists := s.partConfigs[part]
	if !exists {
		return fmt.Errorf("part %q: %w", part, ErrNotFound)
	}
	return anypb.UnmarshalTo(pc.Config, dest, proto.UnmarshalOptions{})
}
