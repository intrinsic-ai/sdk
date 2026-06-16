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
	flagHttpServices        []string
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

	GenerateMainCmd.Flags().StringVar(&flagServiceGoImportPath, "service_go_importpath", "", "Import path for the generated golang code for a gRPC service. (deprecated: use --http_service)")
	GenerateMainCmd.Flags().StringVar(&flagGrpcService, "grpc_service", "", "Fully qualified name of a gRPC service to bridge. (deprecated: use --http_service)")
	GenerateMainCmd.Flags().StringVar(&flagOpenAPIPath, "openapi_path", "", "Path to a generated OpenAPI specification.")
	GenerateMainCmd.Flags().StringVar(&flagOutput, "output", "main.go", "Name of the golang file to generate.")
	GenerateMainCmd.Flags().StringSliceVar(&flagHttpServices, "http_service", nil, "Mapping of gRPC service FQN and Go proto import path, formatted as <service_fqn>:<import_path>. Can be specified multiple times.")
}

type HttpServiceMapping struct {
	GrpcService       string
	GoProtoImportPath string
}

func parseDeprecatedServiceFlags() (*HttpServiceMapping, error) {
	if flagGrpcService == "" && flagServiceGoImportPath == "" {
		return nil, nil
	}
	if flagGrpcService == "" || flagServiceGoImportPath == "" {
		return nil, errors.New("both --grpc_service and --service_go_importpath must be specified when using the deprecated flags")
	}

	// Print deprecation warning showing the correct syntax going forward
	fmt.Fprintf(os.Stderr, "WARNING: Flags --grpc_service and --service_go_importpath are deprecated. "+
		"Please use --http_service \"%s:%s\" instead.\n", flagGrpcService, flagServiceGoImportPath)

	return &HttpServiceMapping{
		GrpcService:       flagGrpcService,
		GoProtoImportPath: flagServiceGoImportPath,
	}, nil
}

func parseHttpServiceFlag(entry string) (*HttpServiceMapping, error) {
	parts := strings.SplitN(entry, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid --http_service mapping format %q; must be <service_fqn>:<import_path>", entry)
	}
	return &HttpServiceMapping{
		GrpcService:       parts[0],
		GoProtoImportPath: parts[1],
	}, nil
}

func run(cmd *cobra.Command, args []string) error {
	err := validateFlags()
	if err != nil {
		return err
	}

	var services []*HttpServiceMapping

	deprecatedService, err := parseDeprecatedServiceFlags()
	if err != nil {
		return err
	}
	if deprecatedService != nil {
		services = append(services, deprecatedService)
	}

	for _, entry := range flagHttpServices {
		s, err := parseHttpServiceFlag(entry)
		if err != nil {
			return err
		}
		services = append(services, s)
	}

	if len(services) == 0 {
		return errors.New("no services specified: use --http_service or (deprecated) --grpc_service and --service_go_importpath")
	}

	// Check for duplicate service arguments
	seen := make(map[string]bool)
	for _, s := range services {
		key := s.GrpcService + ":" + s.GoProtoImportPath
		if seen[key] {
			return fmt.Errorf("duplicate --http_service argument: %q", key)
		}
		seen[key] = true
	}

	fmt.Printf("Generating HTTP/JSON golang binary for %d services...\n", len(services))

	// Read and Base64 encode the OpenAPI spec file
	openAPIBytes, err := os.ReadFile(flagOpenAPIPath)
	if err != nil {
		return fmt.Errorf("failed to read openapi file: %w", err)
	}
	openAPIB64 := base64.StdEncoding.EncodeToString(openAPIBytes)

	// Use all collected information to generate a main file
	err = generateMainDotGo(flagOutput, openAPIB64, services)
	if err != nil {
		return err
	}

	return nil
}

// validateFlags makes sure CLI arguments have usable values.
// It operates on the global flag* variables.
func validateFlags() error {
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

type ServiceTemplateInfo struct {
	ImportPath  string
	ImportAlias string
	ServiceName string
}

func buildTemplateServices(mappings []*HttpServiceMapping) ([]ServiceTemplateInfo, error) {
	var services []ServiceTemplateInfo
	for i, mapping := range mappings {
		grpcServiceType, err := parseFullyQualifiedName(mapping.GrpcService)
		if err != nil {
			return nil, err
		}
		services = append(services, ServiceTemplateInfo{
			ImportPath:  mapping.GoProtoImportPath,
			ImportAlias: fmt.Sprintf("pb%d", i),
			ServiceName: grpcServiceType.Name,
		})
	}
	return services, nil
}

// Generate output_dir/main.go for the http bridge
func generateMainDotGo(outputPath string, openAPIB64 string, mappings []*HttpServiceMapping) error {
	// Parse the template content
	tmpl, err := template.New("main").Parse(mainTemplateContent)
	if err != nil {
		return err
	}

	services, err := buildTemplateServices(mappings)
	if err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Evaluate the template using an anonymous struct with field names matching those expected by the template
	data := struct {
		OpenAPIB64 string
		Services   []ServiceTemplateInfo
	}{
		OpenAPIB64: openAPIB64,
		Services:   services,
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
