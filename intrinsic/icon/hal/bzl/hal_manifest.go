// Copyright 2023 Intrinsic Innovation LLC

// hal_manifest validates and completes a partially filled out manifest for
// hardware modules.
package main

import (
	"bytes"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"text/template"

	"flag"
	log "github.com/golang/glog"
	"google.golang.org/protobuf/encoding/prototext"
	intrinsic "intrinsic/production/intrinsic"
	intrinsicflag "intrinsic/util/flag"

	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"

	_ "embed"
)
const intrinsicIconPath = "/tmp/intrinsic_icon"

var (
	image                = flag.String("image", "", "The image archive file to be included")
	imageSim             = flag.String("image_sim", "", "The image archive file to be included for sim, if applicable")
	manifest             = flag.String("manifest", "", "A textproto file containing the manifest for the HAL asset")
	manifestType         = flag.String("manifest_type", "", "The input manifest type")
	output               = flag.String("output", "-", "Output file name")
	serviceProtoPrefixes = intrinsicflag.MultiString("service_proto_prefix", nil, "Output file name")
	requiresRTPCNode     = flag.Bool("requires_rtpc_node", false, "Whether or not the hardware module requires an RTPC Node")
	requiresAtemsys      = flag.Bool("requires_atemsys", false, "Whether or not the hardware module requires an Atemsys Ethercat device")
	runningEthercatOss   = flag.Bool("running_ethercat_oss", false, "Whether or not the hardware module is running ethercat oss")

	//go:embed hal_service_manifest.textproto.tmpl
	serviceManifestTemplateText string
	serviceManifestTemplate     = template.Must(template.New("manifest").Parse(serviceManifestTemplateText))

	validFamilyIDs = []string{
		"hardware_module",
		"hardware_module_without_geometry",
	}
	generators = map[string]generator{
		"service":  &serviceGen{},
	}
)

type generator interface {
	// template provides a Template used to generate a full manifest using
	// partial manifest data.
	template() *template.Template
	// validatePartial does validation on the partial manifest provided by the
	// user.  This can ensure that user provided parts won't be unexpectedly
	// overwritten or duplicated or that some fields match expectations for the
	// manifest type.
	validatePartial([]byte) error
	// validateFull does validation on the full manifest after it is created
	// using the template.
	validateFull([]byte) error
}

type serviceGen struct{}

func (*serviceGen) template() *template.Template {
	return serviceManifestTemplate
}

func (*serviceGen) validatePartial(partial []byte) error {
	sm := &smpb.ServiceManifest{}
	if err := prototext.Unmarshal(partial, sm); err != nil {
		return fmt.Errorf("unable to parse manifest: %v", err)
	}
	if sm.GetServiceDef() != nil {
		return fmt.Errorf("manifest specifies a service_def, but that will be overwritten")
	}
	if sm.GetAssets() != nil {
		return fmt.Errorf("manifest specifies assets, but overwritten")
	}
	if sm.GetMetadata() == nil {
		return fmt.Errorf("manifest does not specify metadata")
	}
	return nil
}

func (*serviceGen) validateFull(full []byte) error {
	// Parse the completed template to ensure we're always generating a valid
	// manifest.
	sm := &smpb.ServiceManifest{}
	if err := prototext.Unmarshal(full, sm); err != nil {
		return fmt.Errorf("unable to parse: %v", err)
	}
	return nil
}

// fileBaseLeaveEmpty returns the basename of the file, unless the string is
// empty, in which case it's left unchanged.  Otherwise an unspecified flag may
// be interpreted as the path ".".
func fileBaseLeaveEmpty(filename string) string {
	if filename == "" {
		return ""
	}
	return filepath.Base(filename)
}

func halManifest() error {
	if *image == "" {
		return fmt.Errorf("image file not specified")
	}
	if *manifest == "" {
		return fmt.Errorf("manifest file not specified")
	}
	if *output == "" {
		return fmt.Errorf("output file not specified")
	}

	partialManifest, err := os.ReadFile(*manifest)
	if err != nil {
		return fmt.Errorf("unable to load metadata: %v", err)
	}

	var gen generator
	if g, ok := generators[*manifestType]; !ok {
		return fmt.Errorf("unrecognized --manifest_type=%q, allowed values are: %q", *manifestType, slices.Collect(maps.Keys(generators)))
	} else {
		gen = g
	}
	if err := gen.validatePartial(partialManifest); err != nil {
		return fmt.Errorf("manifest was invalid: %v\n%v", err, string(partialManifest))
	}

	var outFile *os.File
	if *output == "-" {
		outFile = os.Stdout
	} else {
		outFile, err = os.Create(*output)
		if err != nil {
			return fmt.Errorf("unable to open %q: %v", *output, err)
		}
	}
	defer outFile.Close()

	data := struct {
		PartialManifest      string
		Image                string
		ImageSim             string
		RequiresRTPC         bool
		RequiresAtemsys      bool
		RunningEthercatOss   bool
		ServiceProtoPrefixes []string
		IntrinsicIconPath    string
	}{
		// The bundle rule always puts the files into the root of the archive,
		// so just take the base image the filename here.
		Image:                fileBaseLeaveEmpty(*image),
		ImageSim:             fileBaseLeaveEmpty(*imageSim),
		PartialManifest:      string(partialManifest),
		RequiresRTPC:         *requiresRTPCNode,
		RequiresAtemsys:      *requiresAtemsys,
		RunningEthercatOss:   *runningEthercatOss,
		ServiceProtoPrefixes: *serviceProtoPrefixes,
		IntrinsicIconPath:    intrinsicIconPath,
	}

	var b bytes.Buffer
	if err = gen.template().Execute(&b, &data); err != nil {
		return fmt.Errorf("unable to format template: %v", err)
	}

	if err := gen.validateFull(b.Bytes()); err != nil {
		return fmt.Errorf("invalid full manifest, file a bug with assets: %v\n%v", err, b.String())
	}

	if _, err := outFile.Write(b.Bytes()); err != nil {
		return fmt.Errorf("unable to write out manifest: %v", err)
	}
	return nil
}

func main() {
	intrinsic.Init()
	if err := halManifest(); err != nil {
		log.Exitf("Unable to create manifest: %v", err)
	}
}
