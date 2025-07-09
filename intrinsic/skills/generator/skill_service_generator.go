// Copyright 2023 Intrinsic Innovation LLC

// Package main provides a skill service generator command line tool.
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"flag"
	log "github.com/golang/glog"
	"google.golang.org/protobuf/proto"
	intrinsic "intrinsic/production/intrinsic"
	gen "intrinsic/skills/generator/gen"

	manifestpb "intrinsic/skills/proto/skill_manifest_go_proto"
)

type stringArray []string

func (i *stringArray) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *stringArray) Set(value string) error {
	if len(*i) > 0 {
		return errors.New("flag already set")
	}

	if value == "" {
		return errors.New("empty value provided for flag")
	}

	for _, s := range strings.Split(value, ",") {
		*i = append(*i, s)
	}
	return nil
}

func (i *stringArray) Get() any {
	return []string(*i)
}

var (
	out           = flag.String("out", "", "The path for the generated file.")
	manifestPath  = flag.String("manifest", "", "The path to the protobin file containing the intrinsic_proto.skills.SkillManifest.")
	lang          = flag.String("lang", "", "The language the skill is implemented in; should be one of: {cpp, python}.")
	ccHeaderPaths = func() *stringArray {
		p := new(stringArray)
		flag.Var(p, "cc_headers", "The comma-separated list of paths to the cpp proto header files for the skill's cpp deps.")
		return p
	}()
)

func main() {
	intrinsic.Init()

	manifestBinary, err := os.ReadFile(*manifestPath)
	if err != nil {
		log.Exitf("cannot read file: %v", err)
	}
	manifest := &manifestpb.SkillManifest{}
	if err := proto.Unmarshal(manifestBinary, manifest); err != nil {
		log.Exitf("cannot unmarshal binary to proto: %v", err)
	}

	switch *lang {
	case "cpp":
		if err := gen.WriteSkillServiceCC(manifest, *ccHeaderPaths, *out); err != nil {
			log.Exitf("Cannot write cc skill service file: %v.", err)
		}
		return
	case "python":
		if err := gen.WriteSkillServicePy(manifest, *out); err != nil {
			log.Exitf("Cannot write py skill service file: %v.", err)
		}
		return
	default:
		log.Exitf("Invalid language selection for skill. lang=%s; should be one of {cpp, python}", *lang)
	}
}
