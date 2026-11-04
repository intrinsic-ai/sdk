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

package shared

import (
	"fmt"
	"strings"
	"testing"
)

func TestInterfacePrint(t *testing.T) {
	// Verify String() does not panic and returns a human readable representation of the DHCP6 pointer
	tv := true
	v := Interface{
		DHCP6: &tv,
	}
	got := fmt.Sprintf("%v", v)
	want := `{"dhcp4":false,"gateway4":"","dhcp6":true,"gateway6":"","mtu":0,"nameservers":{"search":null,"addresses":null},"addresses":null,"realtime":false,"ether_type":0,"display_name":""}`
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestDNSRangeJSON(t *testing.T) {
	t.Run("empty DNSRange is omitted from JSON", func(t *testing.T) {
		iface := Interface{
			DHCP4: true,
		}
		got := fmt.Sprintf("%v", iface)
		if strings.Contains(got, "dns_range") {
			t.Errorf("expected dns_range to be omitted from JSON when empty, got: %s", got)
		}
	})

	t.Run("populated DNSRange is serialized to JSON", func(t *testing.T) {
		iface := Interface{
			DHCP4: true,
			DNSRange: &DNSRange{
				Pool:         "192.168.1.100-192.168.1.200",
				Subnet:       "192.168.1.0/24",
				SearchDomain: "lab.intrinsic.ai",
			},
		}
		got := fmt.Sprintf("%v", iface)
		wantSub := `"dns_range":{"pool":"192.168.1.100-192.168.1.200","subnet":"192.168.1.0/24","search_domain":"lab.intrinsic.ai"}`
		if !strings.Contains(got, wantSub) {
			t.Errorf("expected JSON to contain %s, got: %s", wantSub, got)
		}
	})
}
