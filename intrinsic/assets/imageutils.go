// Copyright 2023 Intrinsic Innovation LLC

// Package imageutils contains docker image utility functions.
package imageutils

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	installerpb "intrinsic/kubernetes/workcell_spec/proto/installer_go_grpc_proto"
	installerservicegrpcpb "intrinsic/kubernetes/workcell_spec/proto/installer_go_grpc_proto"
)

var (
	buildCommand    = "bazel"
	build           = buildExec // Stubbed out for testing.
	buildConfigArgs = []string{
		"-c", "opt",
	}
)

const (
	// Domain used for images managed by Intrinsic.
	registryDomain = "gcr.io"

	// Number of times to try uploading a container image if we get retriable errors.
	remoteWriteTries = 5

	maxImageTagLength = 128
)

// TargetType determines how the "target" target command-line argument will be
// used.
type TargetType string

const (
	// Build mode builds the docker container image using the associated build
	// target name
	Build TargetType = "build"
	// Archive mode assumes the given target points to an already-built image
	Archive TargetType = "archive"
	// ID mode assumes the target is the skill id (only used for stop)
	ID TargetType = "id"
)

// buildExec runs the build command and captures its output.
func buildExec(buildCommand string, buildArgs ...string) ([]byte, error) {
	buildCmd := exec.Command(buildCommand, buildArgs...)
	out, err := buildCmd.Output() // Ignore stderr
	if err != nil {
		return nil, fmt.Errorf("could not build docker image: %v\n%s", err, out)
	}
	return out, nil
}

func getOutputFiles(target string) ([]string, error) {
	buildArgs := []string{"cquery"}
	buildArgs = append(buildArgs, buildConfigArgs...)
	buildArgs = append(buildArgs, "--output=files", target)
	out, err := build(buildCommand, buildArgs...)
	if err != nil {
		return nil, fmt.Errorf("could not get output files: %v\n%s", err, out)
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n"), nil
}

// buildImage builds the given target. The built image's file path is returned.
func buildImage(target string) (string, error) {
	log.Printf("Building image %q using build command %q", target, buildCommand)
	buildArgs := []string{"build"}
	buildArgs = append(buildArgs, buildConfigArgs...)
	buildArgs = append(buildArgs, target)
	out, err := build(buildCommand, buildArgs...)
	if err != nil {
		return "", fmt.Errorf("could not build docker image: %v\n%s", err, out)
	}

	outputs, err := getOutputFiles(target)
	if err != nil {
		return "", fmt.Errorf("could not determine output files: %v", err)
	}
	// Assume rule has a single output file - the built image
	if len(outputs) != 1 {
		return "", fmt.Errorf("could not determine image from target [%s] output files\n%v", target, outputs)
	}
	tarFile := outputs[0]
	if !strings.HasSuffix(tarFile, ".tar") {
		return "", fmt.Errorf("output file did not have .tar extension\n%s", tarFile)
	}
	log.Printf("Finished building and the output filepath is %q", tarFile)
	return string(tarFile), nil
}

// GetImagePath returns the image path.
func GetImagePath(target string, targetType TargetType) (string, error) {
	switch targetType {
	case Build:
		if !strings.HasSuffix(target, ".tar") {
			return "", fmt.Errorf("target should end with .tar")
		}
		return buildImage(target)
	case Archive:
		return target, nil
	default:
		return "", fmt.Errorf("unimplemented target type: %v", targetType)
	}
}

// RemoveContainerParams holds parameters for RemoveContainer.
type RemoveContainerParams struct {
	Address    string
	Connection *grpc.ClientConn
	Request    *installerpb.RemoveContainerAddonRequest
}

// RemoveContainer uses the installer service to remove a new container.
func RemoveContainer(ctx context.Context, params *RemoveContainerParams) error {

	client := installerservicegrpcpb.NewInstallerServiceClient(params.Connection)
	_, err := client.RemoveContainerAddon(ctx, params.Request)
	if status.Code(err) == codes.Unimplemented {
		return fmt.Errorf("installer service not implemented at server side (is it running and accessible at %s?): %v", params.Address, err)
	} else if err != nil {
		return fmt.Errorf("RemoveContainerAddon failed: %v", err)
	}

	return nil
}

// GetRegistry returns the registry to use for images in the specified project.
func GetRegistry(project string) string {
	return fmt.Sprintf("%s/%s", registryDomain, project)
}
