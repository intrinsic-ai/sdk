// Copyright 2023 Intrinsic Innovation LLC

// Package skillio contains utilities that process user-provided skills.
package skillio

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"intrinsic/assets/bundleio"
	"intrinsic/assets/idutils"
)

var (
	buildCommand    = "bazel"
	buildConfigArgs = []string{
		"--config", "intrinsic",
	}
)

// execute runs a command and captures its output.
func execute(buildCommand string, buildArgs ...string) ([]byte, error) {
	c := exec.Command(buildCommand, buildArgs...)
	out, err := c.Output() // Ignore stderr
	if err != nil {
		return nil, fmt.Errorf("exec command failed: %v\n%s", err, out)
	}
	return out, nil
}

func getOutputFiles(target string) ([]string, error) {
	buildArgs := []string{"cquery"}
	buildArgs = append(buildArgs, buildConfigArgs...)
	buildArgs = append(buildArgs, "--output=files", target)
	out, err := execute(buildCommand, buildArgs...)
	if err != nil {
		return nil, fmt.Errorf("could not get output files: %v\n%s", err, out)
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n"), nil
}

func buildTarget(target string) (string, error) {
	buildArgs := []string{"build"}
	buildArgs = append(buildArgs, buildConfigArgs...)
	buildArgs = append(buildArgs, target)
	out, err := execute(buildCommand, buildArgs...)
	if err != nil {
		return "", fmt.Errorf("could not build target %q: %v\n%s", target, err, out)
	}

	outputFiles, err := getOutputFiles(target)
	if err != nil {
		return "", fmt.Errorf("could not get output files of target %s: %v", target, err)
	}
	if len(outputFiles) == 0 {
		return "", fmt.Errorf("target %s did not have any output files", target)
	}
	if len(outputFiles) > 1 {
		log.Printf("Warning: Rule %s was expected to have only one output file, but it had %d", target, len(outputFiles))
	}
	return outputFiles[0], nil
}

// isValidSkillBundle returns true if the given file is a skill bundle file.
func isValidSkillBundle(path string) bool {
	if _, err := bundleio.ReadSkillManifest(path); err != nil {
		return false
	}
	return true
}

// SkillIDFromArchive extracts the skill ID from a skill archive file.
func SkillIDFromArchive(path string) (string, error) {
	if isValidSkillBundle(path) {
		manifest, err := bundleio.ReadSkillManifest(path)
		if err != nil {
			return "", fmt.Errorf("failed to read skill manifest from bundle: %v", err)
		}
		id, err := idutils.IDFromProto(manifest.GetId())
		if err != nil {
			return "", fmt.Errorf("invalid skill ID in manifest: %v", err)
		}
		return id, nil
	}
	return "", fmt.Errorf("%q does not appear to be a valid skill", path)
}

// SkillIDFromBuildTarget extracts the skill ID from a build target.
func SkillIDFromBuildTarget(target string) (string, error) {
	path, err := buildTarget(target)
	if err != nil {
		return "", err
	}
	return SkillIDFromArchive(path)
}
