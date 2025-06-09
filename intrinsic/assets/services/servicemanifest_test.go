// Copyright 2023 Intrinsic Innovation LLC

package servicemanifest

import (
	"testing"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	idpb "intrinsic/assets/proto/id_go_proto"
	vpb "intrinsic/assets/proto/vendor_go_proto"
	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	svpb "intrinsic/assets/services/proto/service_volume_go_proto"
)

func TestValidateServiceManifest(t *testing.T) {
	m := &smpb.ServiceManifest{
		Metadata: &smpb.ServiceMetadata{
			Id: &idpb.Id{
				Name:    "test",
				Package: "package.some",
			},
			DisplayName: "Test Service",
			Vendor: &vpb.Vendor{
				DisplayName: "vendor",
			},
		},
		ServiceDef: &smpb.ServiceDef{
			ServiceProtoPrefixes:  []string{"/intrinsic_proto.services.Calculator/"},
			ConfigMessageFullName: "intrinsic_proto.services.CalculatorConfig",
			SimSpec:               &smpb.ServicePodSpec{},
		},
	}

	mNoService := proto.Clone(m).(*smpb.ServiceManifest)
	mNoService.ServiceDef = nil
	mInvalidName := proto.Clone(m).(*smpb.ServiceManifest)
	mInvalidName.GetMetadata().GetId().Name = "_invalid_name"
	mInvalidPackage := proto.Clone(m).(*smpb.ServiceManifest)
	mInvalidPackage.GetMetadata().GetId().Package = "_invalid_package"
	mNoDisplayName := proto.Clone(m).(*smpb.ServiceManifest)
	mNoDisplayName.GetMetadata().Vendor.DisplayName = ""
	mNoVendor := proto.Clone(m).(*smpb.ServiceManifest)
	mNoVendor.GetMetadata().Vendor = nil
	mNoSimSpec := proto.Clone(m).(*smpb.ServiceManifest)
	mNoSimSpec.ServiceDef = &smpb.ServiceDef{}
	mInvalidServiceProtoPrefix := proto.Clone(m).(*smpb.ServiceManifest)
	mInvalidServiceProtoPrefix.GetServiceDef().ServiceProtoPrefixes = []string{"intrinsic_proto.services.Calculator"}
	mInvalidConfigMessageFullName := proto.Clone(m).(*smpb.ServiceManifest)
	mInvalidConfigMessageFullName.GetServiceDef().ConfigMessageFullName = "invalid.config.Message"

	fds := &descriptorpb.FileDescriptorSet{
		File: []*descriptorpb.FileDescriptorProto{
			&descriptorpb.FileDescriptorProto{
				Name:    proto.String("intrinsic_proto.services.Calculator"),
				Package: proto.String("intrinsic_proto.services"),
				MessageType: []*descriptorpb.DescriptorProto{
					&descriptorpb.DescriptorProto{
						Name: proto.String("Calculator"),
					},
				},
			},
			&descriptorpb.FileDescriptorProto{
				Name:    proto.String("intrinsic_proto.services.CalculatorConfig"),
				Package: proto.String("intrinsic_proto.services"),
				MessageType: []*descriptorpb.DescriptorProto{
					&descriptorpb.DescriptorProto{
						Name: proto.String("CalculatorConfig"),
					},
				},
			},
		},
	}
	files, err := protodesc.NewFiles(fds)
	if err != nil {
		t.Fatalf("failed to find file descriptor: %v", err)
	}

	tests := []struct {
		desc    string
		given   *smpb.ServiceManifest
		opts    []ValidateServiceManifestOption
		wantErr bool
	}{
		{
			desc:    "valid service manifest, but no descriptors given",
			given:   m,
			wantErr: true,
		},
		{
			desc:  "valid service manifest without service def",
			given: mNoService,
		},
		{
			desc:    "empty service manifest",
			given:   &smpb.ServiceManifest{},
			wantErr: true,
		},
		{
			desc:    "invalid name",
			given:   mInvalidName,
			wantErr: true,
		},
		{
			desc:    "invalid package",
			given:   mInvalidPackage,
			wantErr: true,
		},
		{
			desc:    "no display name",
			given:   mNoDisplayName,
			wantErr: true,
		},
		{
			desc:    "no vendor",
			given:   mNoVendor,
			wantErr: true,
		},
		{
			desc:    "no sim spec",
			given:   mNoSimSpec,
			wantErr: true,
		},
		{
			desc:    "missing service descriptor",
			given:   m,
			opts:    []ValidateServiceManifestOption{WithFiles(&protoregistry.Files{})},
			wantErr: true,
		},
		{
			desc:  "valid service descriptor",
			given: m,
			opts: []ValidateServiceManifestOption{
				WithFiles(files),
			},
		},
		{
			desc:    "invalid service proto prefix",
			given:   mInvalidServiceProtoPrefix,
			wantErr: true,
		},
		{
			desc:    "invalid config message full name",
			given:   mInvalidConfigMessageFullName,
			wantErr: true,
		},
		{
			desc: "valid volumes",
			given: serviceManifestWithVolumes(
				[]*svpb.Volume{
					&svpb.Volume{
						Name: "volume1",
						Source: &svpb.Volume_HostPath{
							HostPath: &svpb.HostPathVolumeSource{
								Path: "/path/to/volume1",
							},
						},
					},
				},
				[]*svpb.VolumeMount{
					&svpb.VolumeMount{
						Name:      "volume1",
						MountPath: "/path/to/volume1",
					},
				},
			),
		},
		{
			desc: "invalid volume name",
			given: serviceManifestWithVolumes(
				[]*svpb.Volume{
					&svpb.Volume{
						Name: "volume_one",
						Source: &svpb.Volume_HostPath{
							HostPath: &svpb.HostPathVolumeSource{
								Path: "/path/to/volume_one",
							},
						},
					},
				},
				[]*svpb.VolumeMount{
					&svpb.VolumeMount{
						Name:      "volume_one",
						MountPath: "/path/to/volume_one",
					},
				},
			),
			wantErr: true,
		},
		{
			desc: "invalid volume path",
			given: serviceManifestWithVolumes(
				[]*svpb.Volume{
					&svpb.Volume{
						Name: "volume1",
						Source: &svpb.Volume_HostPath{
							HostPath: &svpb.HostPathVolumeSource{
								Path: "/path/to/'bad_folder'",
							},
						},
					},
				},
				[]*svpb.VolumeMount{
					&svpb.VolumeMount{
						Name:      "volume1",
						MountPath: "/path/to/volume1/",
					},
				},
			),
			wantErr: true,
		},
		{
			desc: "duplicate volume name",
			given: serviceManifestWithVolumes(
				[]*svpb.Volume{
					&svpb.Volume{
						Name: "volume1",
						Source: &svpb.Volume_HostPath{
							HostPath: &svpb.HostPathVolumeSource{
								Path: "/path/to/volume1",
							},
						},
					},
					&svpb.Volume{
						Name: "volume1",
						Source: &svpb.Volume_HostPath{
							HostPath: &svpb.HostPathVolumeSource{
								Path: "/path/to/volume_one",
							},
						},
					},
				},
				[]*svpb.VolumeMount{
					&svpb.VolumeMount{
						Name:      "volume1",
						MountPath: "/path/to/volume1",
					},
				},
			),
			wantErr: true,
		},
		{
			desc: "volume mount references non-existent volume",
			given: serviceManifestWithVolumes(
				[]*svpb.Volume{
					&svpb.Volume{
						Name: "volume1",
						Source: &svpb.Volume_HostPath{
							HostPath: &svpb.HostPathVolumeSource{
								Path: "/path/to/volume1",
							},
						},
					},
				},
				[]*svpb.VolumeMount{
					&svpb.VolumeMount{
						Name:      "volume2",
						MountPath: "/path/to/volume2",
					},
				},
			),
			wantErr: true,
		},
		{
			desc: "invalid volume mount path",
			given: serviceManifestWithVolumes(
				[]*svpb.Volume{
					&svpb.Volume{
						Name: "volume1",
						Source: &svpb.Volume_HostPath{
							HostPath: &svpb.HostPathVolumeSource{
								Path: "/path/to/volume1",
							},
						},
					},
				},
				[]*svpb.VolumeMount{
					&svpb.VolumeMount{
						Name:      "volume1",
						MountPath: "/path/to/'bad_folder'",
					},
				},
			),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			if err := ValidateServiceManifest(tc.given, tc.opts...); (err != nil) != tc.wantErr {
				t.Fatalf("ValidateServiceManifest(%v) returned unexpected error, got: %v, want: %v", tc.given, err, tc.wantErr)
			}
		})
	}
}

func TestValidateVolume(t *testing.T) {
	tests := []struct {
		desc    string
		given   *svpb.Volume
		wantErr bool
	}{
		{
			desc: "valid host path volume",
			given: &svpb.Volume{
				Name: "volume1",
				Source: &svpb.Volume_HostPath{
					HostPath: &svpb.HostPathVolumeSource{
						Path: "/path/to/volume1",
					},
				},
			},
		},
		{
			desc: "valid empty dir volume",
			given: &svpb.Volume{
				Name: "volume1",
				Source: &svpb.Volume_EmptyDir{
					EmptyDir: &svpb.EmptyDirVolumeSource{
						Medium: svpb.EmptyDirMedium_EMPTY_DIR_MEDIUM_MEMORY,
					},
				},
			},
		},
		{
			desc: "no volume source",
			given: &svpb.Volume{
				Name: "volume1",
			},
			wantErr: true,
		},
		{
			desc: "invalid volume name",
			given: &svpb.Volume{
				Name: "_invalid_name",
				Source: &svpb.Volume_HostPath{
					HostPath: &svpb.HostPathVolumeSource{
						Path: "/path/to/volume1",
					},
				},
			},
			wantErr: true,
		},
		{
			desc: "invalid host path volume",
			given: &svpb.Volume{
				Name: "volume1",
				Source: &svpb.Volume_HostPath{
					HostPath: &svpb.HostPathVolumeSource{
						Path: "/path/to/'bad_folder'",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			if err := ValidateVolume(tc.given); (err != nil) != tc.wantErr {
				t.Fatalf("ValidateVolume(%v) returned unexpected error, got: %v, want: %v", tc.given, err, tc.wantErr)
			}
		})
	}
}

func TestValidateVolumeMount(t *testing.T) {
	tests := []struct {
		desc    string
		given   *svpb.VolumeMount
		wantErr bool
	}{
		{
			desc: "valid volume mount",
			given: &svpb.VolumeMount{
				Name:      "volume1",
				MountPath: "/path/to/volume1",
			},
		},
		{
			desc: "invalid volume name",
			given: &svpb.VolumeMount{
				Name:      "_invalid_name",
				MountPath: "/path/to/volume1",
			},
			wantErr: true,
		},
		{
			desc: "invalid mount path",
			given: &svpb.VolumeMount{
				Name:      "volume1",
				MountPath: "/path/to/'bad_folder'",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			if err := ValidateVolumeMount(tc.given); (err != nil) != tc.wantErr {
				t.Fatalf("ValidateVolumeMount(%v) returned unexpected error, got: %v, want: %v", tc.given, err, tc.wantErr)
			}
		})
	}
}

func serviceManifestWithVolumes(volumes []*svpb.Volume, volumeMounts []*svpb.VolumeMount) *smpb.ServiceManifest {
	return &smpb.ServiceManifest{
		Metadata: &smpb.ServiceMetadata{
			Id: &idpb.Id{
				Name:    "test",
				Package: "package.some",
			},
			DisplayName: "Test Service",
			Vendor: &vpb.Vendor{
				DisplayName: "vendor",
			},
		},
		ServiceDef: &smpb.ServiceDef{
			SimSpec: &smpb.ServicePodSpec{
				Image: &smpb.ServiceImage{
					Settings: &smpb.ServiceImageSettings{
						VolumeMounts: volumeMounts,
					},
				},
				Settings: &smpb.ServicePodSettings{
					Volumes: volumes,
				},
			},
		},
	}
}
