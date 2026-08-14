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

package main

import (
	"bytes"
	"testing"

	"google.golang.org/protobuf/encoding/prototext"

	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	svpb "intrinsic/assets/services/proto/service_volume_go_proto"
)

func TestServiceManifestTemplate_RunningEthercatOss(t *testing.T) {
	partialManifest := `
metadata {
  id {
    package: "ai.intrinsic"
    name: "test_module"
  }
  vendor {
    display_name: "Intrinsic"
  }
  documentation {
    description: "Test"
  }
  display_name: "Test"
}
`

	tests := []struct {
		desc               string
		runningEthercatOss bool
		wantCharDev        bool
	}{
		{
			desc:               "running ethercat oss includes etherlab volume with CharDevice type",
			runningEthercatOss: true,
			wantCharDev:        true,
		},
		{
			desc:               "not running ethercat oss excludes etherlab volume",
			runningEthercatOss: false,
			wantCharDev:        false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			data := struct {
				PartialManifest           string
				Image                     string
				ImageSim                  string
				ProvidesServiceInspection bool
				RequiresRTPC              bool
				RequiresAtemsys           bool
				RunningEthercatOss        bool
				ServiceProtoPrefixes      []string
				IntrinsicIconPath         string
			}{
				PartialManifest:    partialManifest,
				Image:              "real_image.tar",
				ImageSim:           "sim_image.tar",
				RunningEthercatOss: tc.runningEthercatOss,
				IntrinsicIconPath:  "/tmp/intrinsic_icon",
			}

			var b bytes.Buffer
			if err := serviceManifestTemplate.Execute(&b, &data); err != nil {
				t.Fatalf("serviceManifestTemplate.Execute failed: %v", err)
			}

			if err := validateFull(b.Bytes()); err != nil {
				t.Fatalf("validateFull failed: %v\nManifest:\n%s", err, b.String())
			}

			sm := &smpb.ServiceManifest{}
			if err := prototext.Unmarshal(b.Bytes(), sm); err != nil {
				t.Fatalf("prototext.Unmarshal failed: %v", err)
			}

			volumes := sm.GetServiceDef().GetRealSpec().GetSettings().GetVolumes()
			var foundEtherlab bool
			for _, v := range volumes {
				if v.GetName() == "etherlab" {
					foundEtherlab = true
					if gotPath := v.GetHostPath().GetPath(); gotPath != "/dev/EtherCAT0" {
						t.Errorf("etherlab host_path path = %q, want /dev/EtherCAT0", gotPath)
					}
					if gotType := v.GetHostPath().GetType(); gotType != svpb.HostPathVolumeSourceType_HOST_PATH_VOLUME_SOURCE_TYPE_CHAR_DEVICE {
						t.Errorf("etherlab host_path type = %v, want %v", gotType, svpb.HostPathVolumeSourceType_HOST_PATH_VOLUME_SOURCE_TYPE_CHAR_DEVICE)
					}
				}
			}

			if foundEtherlab != tc.wantCharDev {
				t.Errorf("found etherlab volume = %v, want %v", foundEtherlab, tc.wantCharDev)
			}
		})
	}
}
