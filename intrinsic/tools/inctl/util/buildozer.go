// Copyright 2023 Intrinsic Innovation LLC

// package buildozer contains utilities for using buildozer
package buildozer

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bazelbuild/buildtools/edit"
)

// ExecuteBuildozerCommands runs the given package-level buildozer commands
// (which have the form "buildozer ... //my_package:__pkg__"). Modify or
// fork this function if you need to run non-package-level buildozer commands.
func ExecuteBuildozerCommands(cmds []string, bazelWorkspaceDir string, bazelPackage []string) error {
	opts := edit.NewOpts()
	opts.RootDir = bazelWorkspaceDir
	packageLabel := fmt.Sprintf("//%s:__pkg__", strings.Join(bazelPackage, "/"))

	for _, cmd := range cmds {
		// Capture and suppress output (buildozer uses stdout/stderr by default)
		// and only print it in case of an error (see below).
		var out, err bytes.Buffer
		opts.OutWriter = &out
		opts.ErrWriter = &err

		args := []string{cmd, packageLabel}
		result := edit.Buildozer(opts, args)

		// Buildozer return codes:
		// 0: success
		// 3: no error, but no files were modified
		if result != 0 && result != 3 {
			return fmt.Errorf("command %q returned with error code %d:\n%s",
				strings.Join(args, " "), result, err.String())
		}
	}

	return nil
}
