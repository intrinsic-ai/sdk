// Copyright 2023 Intrinsic Innovation LLC

// Package generatemain contains the entry point for inbuild httpservice generate.
package generatemain

import (
	_ "embed"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var (
	flagServiceGoImportPath string
	flagGrpcService         string
	flagOpenAPIPath         string
	flagOutput              string
)

//go:embed templates/main.go.template
var mainTemplateContent string

// GenerateMainCmd creates files for http bridges to services
var GenerateMainCmd *cobra.Command

// The init function establishes command line flags for `inbuild httpjson generatemain`
func init() {
	resetGenerateCommand()
}

// Reset global variables so unit tests don't interfere with each other.
func resetGenerateCommand() {
	GenerateMainCmd = &cobra.Command{
		Use:   "generatemain",
		Short: "Creates a main.go that offers HTTP/JSON endpoints for a gRPC service.",
		Long:  "Creates a main.go that offers HTTP/JSON endpoints for a gRPC service.",
		RunE:  run,
	}

	// Updated flags to match current variable declarations
	GenerateMainCmd.Flags().StringVar(&flagServiceGoImportPath, "service_go_importpath", "", "Import path for the generated golang code for a gRPC service.")
	GenerateMainCmd.Flags().StringVar(&flagGrpcService, "grpc_service", "", "Fully qualified name of a gRPC service to bridge.")
	GenerateMainCmd.Flags().StringVar(&flagOpenAPIPath, "openapi_path", "", "Path to a generated OpenAPI specification.")
	GenerateMainCmd.Flags().StringVar(&flagOutput, "output", "main.go", "Name of the golang file to generate.")
}

func run(cmd *cobra.Command, args []string) error {
	err := validateFlags()
	if err != nil {
		return err
	}

	fmt.Printf("Generating HTTP/JSON golang binary for %s...\n", flagServiceGoImportPath)

	// Read and Base64 encode the OpenAPI spec file
	openAPIBytes, err := os.ReadFile(flagOpenAPIPath)
	if err != nil {
		return fmt.Errorf("failed to read openapi file: %w", err)
	}
	openAPIB64 := base64.StdEncoding.EncodeToString(openAPIBytes)

	grpcServiceType, err := parseFullyQualifiedName(flagGrpcService)
	if err != nil {
		return err
	}

	// Use all collected information to generate the a main file
	err = generateMainDotGo(flagOutput, flagServiceGoImportPath, openAPIB64, grpcServiceType)
	if err != nil {
		return err
	}

	return nil
}

// validateFlags makes sure CLI arguments have usable values.
// It operates on the global flag* variables.
func validateFlags() error {
	// Validate Required Flags
	if flagServiceGoImportPath == "" {
		return errors.New("--service_go_importpath is required")
	}
	if flagGrpcService == "" {
		return errors.New("--grpc_service is required")
	}
	if flagOpenAPIPath == "" {
		return errors.New("--openapi_path is required")
	}
	if flagOutput == "" {
		return errors.New("--output is required")
	}

	pathsToConvert := []string{
		flagOpenAPIPath,
		flagOutput,
	}

	absPaths, err := makeAbsolutePaths(pathsToConvert)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute paths: %w", err)
	}

	flagOpenAPIPath = absPaths[0]
	flagOutput = absPaths[1]

	return nil
}

// Generate output_dir/main.go for the http bridge
func generateMainDotGo(outputPath string, serviceGoImportpath string, openAPIB64 string, serviceType *ProtoType) error {
	// Parse the template content
	tmpl, err := template.New("main").Parse(mainTemplateContent)
	if err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Evaluate the template using an anonymous struct with field names matching those expected by the temlate
	data := struct {
		ServiceGoImportpath string
		OpenAPIB64          string
		ServiceName         string
	}{
		OpenAPIB64:          openAPIB64,
		ServiceName:         serviceType.Name,
		ServiceGoImportpath: serviceGoImportpath,
	}

	return tmpl.Execute(f, data)
}

func makeAbsolutePaths(paths []string) ([]string, error) {
	var absolutePaths []string
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		absolutePaths = append(absolutePaths, absPath)
	}
	return absolutePaths, nil
}

func parseFullyQualifiedName(fqn string) (*ProtoType, error) {
	// Find the last index of the dot to separate the package from the name
	lastDotIndex := strings.LastIndex(fqn, ".")
	if lastDotIndex == -1 {
		return nil, errors.New("invalid FQN: no dot separator found")
	}

	serviceName := fqn[lastDotIndex+1:]
	packageName := fqn[:lastDotIndex]

	// Handle edge case where string ends in a dot (e.g., "com.example.")
	if serviceName == "" {
		return nil, errors.New("invalid FQN: name component is empty")
	}

	return &ProtoType{
		Name:    serviceName,
		Package: packageName,
	}, nil
}

type ProtoType struct {
	Name    string
	Package string
}

func (pt ProtoType) FullyQualifiedName() string {
	return pt.Package + "." + pt.Name
}
