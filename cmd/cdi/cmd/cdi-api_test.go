/*
   Copyright 2026 The CDI Authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package cmd

import (
	"testing"

	oci "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/stretchr/testify/require"
)

func TestCollectCDIDevicesFromOCISpec(t *testing.T) {
	for _, tc := range []struct {
		name        string
		spec        *oci.Spec
		devices     []string
		linuxDevs   []oci.LinuxDevice
		emptyDevice bool
	}{
		{
			name: "collect CDI devices and keep non-CDI devices",
			spec: &oci.Spec{
				Linux: &oci.Linux{
					Devices: []oci.LinuxDevice{
						{Path: "vendor.com/device=dev0"},
						{Path: "/dev/null", Type: "c", Major: 1, Minor: 3},
						{Path: "/dev/null", Type: "c", Major: 1, Minor: 5},
						{Path: "vendor.com/device=dev1"},
					},
				},
			},
			devices: []string{
				"vendor.com/device=dev0",
				"vendor.com/device=dev1",
			},
			linuxDevs: []oci.LinuxDevice{
				{Path: "/dev/null", Type: "c", Major: 1, Minor: 5},
			},
		},
		{
			name: "all CDI devices leave an empty device slice",
			spec: &oci.Spec{
				Linux: &oci.Linux{
					Devices: []oci.LinuxDevice{
						{Path: "vendor.com/device=dev0"},
					},
				},
			},
			devices:     []string{"vendor.com/device=dev0"},
			linuxDevs:   []oci.LinuxDevice{},
			emptyDevice: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			devices := collectCDIDevicesFromOCISpec(tc.spec)

			require.Equal(t, tc.devices, devices)
			require.Equal(t, tc.linuxDevs, tc.spec.Linux.Devices)
			if tc.emptyDevice {
				require.NotNil(t, tc.spec.Linux.Devices)
			}
		})
	}
}
