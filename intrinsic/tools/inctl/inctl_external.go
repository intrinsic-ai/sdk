// Copyright 2023 Intrinsic Innovation LLC

package main

import (
	"intrinsic/tools/inctl/cmd/root"

	_ "intrinsic/assets/inctl/assetcmd"
	_ "intrinsic/assets/services/inctl/service"
	_ "intrinsic/tools/inctl/cmd/auth/auth"
	_ "intrinsic/tools/inctl/cmd/bazel/bazel"
	_ "intrinsic/tools/inctl/cmd/cluster/cluster"
	_ "intrinsic/tools/inctl/cmd/customer/customer"
	_ "intrinsic/tools/inctl/cmd/device/device"
	_ "intrinsic/tools/inctl/cmd/doctor/doctor" 
	_ "intrinsic/tools/inctl/cmd/logs/logs"
	_ "intrinsic/tools/inctl/cmd/markdown"
	_ "intrinsic/tools/inctl/cmd/notebook/notebook"
	_ "intrinsic/tools/inctl/cmd/process/process"
	_ "intrinsic/tools/inctl/cmd/recordings/recordings"
	_ "intrinsic/tools/inctl/cmd/skill"
	_ "intrinsic/tools/inctl/cmd/solution/solution"
	_ "intrinsic/tools/inctl/cmd/version/version"
)

func main() {
	root.Inctl()
}
