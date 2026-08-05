// Copyright 2023 Intrinsic Innovation LLC

package auth

import (
	"context"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/metadata"
)

// Same as authtest.NewStoreForTest which is not accessible here as long as
// this test is in the 'auth' package (cyclic dependency).
func newStoreForTest(t *testing.T) *Store {
	configDir := t.TempDir()
	return &Store{
		GetConfigDirFx:      func() (string, error) { return configDir, nil },
		EnvironmentResolver: &defaultEnvironmentResolver{},
	}
}

func TestRFC3339Time_Marshaling(t *testing.T) {
	tests := []struct {
		input *RFC3339Time
		// ignores want, just ensures that it can get the same value
		biDirectional bool
		want          string
	}{
		{input: toRFC3339Time(time.Now().Truncate(time.Second)), biDirectional: true},
		{input: toRFC3339Time(time.Date(2000, 1, 10, 10, 9, 8, 0, time.UTC)), biDirectional: true},
		{input: toRFC3339Time(time.Date(2023, 4, 5, 0, 8, 47, 0, time.UTC)), want: "2023-04-05T00:08:47Z"},
	}

	for _, test := range tests {
		value, err := test.input.MarshalText()
		if err != nil {
			t.Errorf("cannot marshal time: %v", err)
		}
		if test.biDirectional {
			helper := new(RFC3339Time)
			if err = helper.UnmarshalText(value); err != nil {
				t.Errorf("cannot unmarshal time (%s): %v", value, err)
			}
			input := time.Time(*test.input)
			got := time.Time(*helper)
			if !input.Equal(got) {
				t.Errorf("output mismatch: got %s; wants: %s", got, input)
			}
		} else {
			// compares with fixed want value on string basis
			if string(value) != test.want {
				t.Errorf("output mismatch: got %s; want: %s", string(value), test.want)
			}
		}
	}
}

type countingEnvironmentResolver struct {
	calls       int
	envResolver EnvironmentResolver
}

func (r *countingEnvironmentResolver) HasProject(env, project string) (bool, error) {
	r.calls++
	return r.envResolver.HasProject(env, project)
}

func TestStore_GetConfiguration(t *testing.T) {
	tests := []struct {
		name              string
		projectName       string
		wantEnv           string
		setup             func(t *testing.T, s *Store, proj, env string)
		wantErr           bool
		wantResolverCalls int
	}{
		{
			name:        "a config was cached",
			projectName: "test-prod-project",
			wantEnv:     "prod",
			setup: func(t *testing.T, s *Store, p, e string) {
				if err := s.WriteProjectEnvironment(p, e); err != nil {
					t.Fatalf("WriteProjectInfo failed: %v", err)
				}
				cfg := NewConfiguration(e)
				cfg.Tokens[AliasDefaultToken] = &ProjectToken{
					APIKey:     "valid-key",
					ValidUntil: toRFC3339Time(time.Now().Add(24 * time.Hour)),
				}
				if err := s.WriteEnvConfiguration(cfg); err != nil {
					t.Fatalf("WriteEnvConfiguration failed: %v", err)
				}
			},
			wantResolverCalls: 0,
		},
		{
			name:        "a config wasn't cached and we had to make a network call",
			projectName: "test-prod-project",
			wantEnv:     "prod",
			setup: func(t *testing.T, s *Store, p, e string) {
				cfg := NewConfiguration(e)
				cfg.Tokens[AliasDefaultToken] = &ProjectToken{
					APIKey:     "valid-key",
					ValidUntil: toRFC3339Time(time.Now().Add(24 * time.Hour)),
				}
				if err := s.WriteEnvConfiguration(cfg); err != nil {
					t.Fatalf("WriteEnvConfiguration failed: %v", err)
				}
			},
			wantResolverCalls: 1,
		},
		{
			name:        "a matching environment was not found",
			projectName: "test-prod-project",
			wantEnv:     "dev",
			setup: func(t *testing.T, s *Store, p, e string) {
				cfg := NewConfiguration(e)
				cfg.Tokens[AliasDefaultToken] = &ProjectToken{
					APIKey:     "valid-key",
					ValidUntil: toRFC3339Time(time.Now().Add(24 * time.Hour)),
				}
				if err := s.WriteEnvConfiguration(cfg); err != nil {
					t.Fatalf("WriteEnvConfiguration failed: %v", err)
				}
			},
			wantErr:           true,
			wantResolverCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStoreForTest(t)

			resolver := &countingEnvironmentResolver{envResolver: s.EnvironmentResolver}
			s.EnvironmentResolver = resolver

			tt.setup(t, s, tt.projectName, tt.wantEnv)

			got, err := s.GetConfiguration(tt.projectName)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if got == nil || got.Name != tt.wantEnv {
					t.Errorf("GetConfiguration() got envName = %v, want %v", got, tt.wantEnv)
				}
			}

			if resolver.calls != tt.wantResolverCalls {
				t.Errorf("EnvironmentResolver called %d times, want %d", resolver.calls, tt.wantResolverCalls)
			}
		})
	}
}

func toRFC3339Time(time time.Time) *RFC3339Time {
	result := RFC3339Time(time)
	return &result
}

func TestStore_OrgInfoEquality(t *testing.T) {
	tests := []struct {
		name string
		want OrgInfo
	}{
		{
			name: "simple",
			want: OrgInfo{Organization: "org", Project: "project"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newStoreForTest(t)
			err := s.WriteOrgInfo(&tc.want)
			if err != nil {
				t.Fatalf("WriteOrgInfo returned an unexpected error: %v", err)
			}

			got, err := s.ReadOrgInfo(tc.want.Organization)
			if err != nil {
				t.Fatalf("OrgInfo returned an unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("OrgInfo returned an unexpected diff (-want +got): %v", diff)
			}
		})
	}
}

func TestStore_AuthorizeContext(t *testing.T) {
	projectName := "friendly-name"

	tests := []struct {
		name                      string
		givenProjectConfiguration *ProjectConfiguration
		ctx                       context.Context
		projectName               string
		envName                   string
		wantOutgoingMetadata      metadata.MD
		wantErr                   bool
	}{
		{
			name: "adds authorization header to the context",
			givenProjectConfiguration: &ProjectConfiguration{
				Name: "test-env",
				Tokens: map[string]*ProjectToken{
					AliasDefaultToken: {
						APIKey:     "abcdefg.xyz",
						ValidUntil: toRFC3339Time(time.Now().Add(24 * time.Hour)),
					},
				},
			},
			ctx:                  context.Background(),
			projectName:          projectName,
			envName:              "test-env",
			wantOutgoingMetadata: metadata.Pairs("authorization", "Bearer abcdefg.xyz"),
		},
		{
			name: "does not change an existing authorization header",
			givenProjectConfiguration: &ProjectConfiguration{
				Name: "test-env",
				Tokens: map[string]*ProjectToken{
					AliasDefaultToken: {
						APIKey:     "abcdefg.xyz",
						ValidUntil: toRFC3339Time(time.Now().Add(24 * time.Hour)),
					},
				},
			},
			ctx: metadata.NewOutgoingContext(
				context.Background(),
				metadata.Pairs("authorization", "Bearer existing.token"),
			),
			projectName:          projectName,
			envName:              "test-env",
			wantOutgoingMetadata: metadata.Pairs("authorization", "Bearer existing.token"),
		},
		{
			name: "fails if there is no authorization information for the project",
			givenProjectConfiguration: &ProjectConfiguration{
				Name: "test-env",
				Tokens: map[string]*ProjectToken{
					AliasDefaultToken: {
						APIKey:     "abcdefg.xyz",
						ValidUntil: toRFC3339Time(time.Now().Add(24 * time.Hour)),
					},
				},
			},
			ctx:         context.Background(),
			projectName: "other-project",
			envName:     "test-env",
			wantErr:     true,
		},
		{
			name: "fails if there is no default credential for the project",
			givenProjectConfiguration: &ProjectConfiguration{
				Name: "test-env",
				Tokens: map[string]*ProjectToken{
					"not-default": {
						APIKey:     "abcdefg.xyz",
						ValidUntil: toRFC3339Time(time.Now().Add(24 * time.Hour)),
					},
				},
			},
			ctx:         context.Background(),
			projectName: projectName,
			envName:     "test-env",
			wantErr:     true,
		},
		{
			name: "fails if the default credential for the project is invalid",
			givenProjectConfiguration: &ProjectConfiguration{
				Name: "test-env",
				Tokens: map[string]*ProjectToken{
					AliasDefaultToken: {
						APIKey:     "", // empty key is invalid
						ValidUntil: toRFC3339Time(time.Now().Add(24 * time.Hour)),
					},
				},
			},
			ctx:         context.Background(),
			projectName: projectName,
			envName:     "test-env",
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := newStoreForTest(t)

			if err := store.WriteProjectEnvironment(projectName, tc.envName); err != nil {
				t.Fatalf("WriteProjectEnvironment returned an unexpected error: %v", err)
			}
			if err := store.WriteEnvConfiguration(tc.givenProjectConfiguration); err != nil {
				t.Fatalf("WriteEnvConfiguration(%v) returned an unexpected error: %v", tc.givenProjectConfiguration, err)
			}

			got, err := store.AuthorizeContext(tc.ctx, tc.projectName)
			if !tc.wantErr && err != nil {
				t.Errorf("AuthorizeContext(%v) returned an unexpected error: %v", tc.projectName, err)
			}
			if tc.wantErr && err == nil {
				t.Errorf("AuthorizeContext(%v) returned no error, want error", tc.projectName)
			}

			gotOutgoingMetadata, _ := metadata.FromOutgoingContext(got)
			if diff := cmp.Diff(tc.wantOutgoingMetadata, gotOutgoingMetadata); diff != "" {
				t.Errorf("AuthorizeContext(%v) has unexpected metadata (-want +got): %v", tc.projectName, diff)
			}
		})
	}
}

func TestStore_RemoveOrganization(t *testing.T) {
	type fields struct {
		projects []ProjectConfiguration
		orgs     []OrgInfo
	}
	type args struct {
		name string
	}
	type wants struct {
		projects []string
		orgs     []string
		wantErr  bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wants  wants
	}{
		{
			name: "single-organization",
			fields: fields{
				orgs: []OrgInfo{
					{Organization: "first-org", Project: "first-project"},
				},
				projects: []ProjectConfiguration{
					{Name: "first-project"},
				},
			},
			args:  args{name: "first-org"},
			wants: wants{wantErr: false},
		},
		{
			name: "shared-project-not-removed",
			fields: fields{
				orgs: []OrgInfo{
					{Organization: "first-org", Project: "first-project"},
					{Organization: "second-org", Project: "first-project"},
				}, projects: []ProjectConfiguration{{Name: "first-project"}},
			},
			args: args{name: "first-org"},
			wants: wants{
				projects: []string{"first-project"},
				orgs:     []string{"second-org"},
				wantErr:  false,
			},
		},
		{
			name: "fail-remove-non-existent",
			fields: fields{
				orgs: []OrgInfo{
					{Organization: "first-org", Project: "first-project"},
				},
				projects: []ProjectConfiguration{
					{Name: "first-project"},
				},
			},
			args: args{name: "second-org"},
			wants: wants{
				projects: []string{"first-project"},
				orgs:     []string{"first-org"},
				wantErr:  true,
			},
		},
		{
			name: "ignore-missing-project",
			fields: fields{
				orgs: []OrgInfo{
					{Organization: "first-org", Project: "first-project"},
				},
				projects: []ProjectConfiguration{
					{Name: "second-project"},
				},
			},
			args: args{name: "first-org"},
			wants: wants{
				projects: []string{"second-project"},
				wantErr:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStoreForTest(t)

			for _, project := range tt.fields.projects {
				err := s.WriteEnvConfiguration(&project)
				if err != nil {
					t.Errorf("cannot write project env %v: %s", project, err)
				}
				err = s.WriteProjectEnvironment(project.Name, "prod")
				if err != nil {
					t.Errorf("cannot write project info %v: %s", project, err)
				}
			}

			for _, org := range tt.fields.orgs {
				if err := s.WriteOrgInfo(&org); err != nil {
					t.Errorf("cannot write organization %v: %s", org, err)
				}
			}

			if err := s.RemoveOrganization(tt.args.name); (err != nil) != tt.wants.wantErr {
				t.Errorf("RemoveOrganization() error = %v, wantErr %v", err, tt.wants.wantErr)
			}

			projects, err := s.ListConfigurations()
			if err != nil {
				t.Errorf("unexpected error listing projects: %s", err)
			}
			slices.Sort(projects)
			slices.Sort(tt.wants.projects)
			orgs, err := s.ListOrgs()
			if err != nil {
				t.Errorf("unexpected error listing orgs: %s", err)
			}
			slices.Sort(orgs)
			slices.Sort(tt.wants.orgs)
			if diff := cmp.Diff(projects, tt.wants.projects); diff != "" {
				t.Errorf("unexpected projects: %q", diff)
			}
			if diff := cmp.Diff(orgs, tt.wants.orgs); diff != "" {
				t.Errorf("unexpected organizations: %q", diff)
			}
		})
	}
}

func TestStore_RemoveAllKnownCredentials(t *testing.T) {
	type fields struct {
		projects []ProjectConfiguration
		orgs     []OrgInfo
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "single-organization",
			fields: fields{
				orgs: []OrgInfo{
					{Organization: "first-org", Project: "first-project"},
				},
				projects: []ProjectConfiguration{
					{Name: "first-project"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStoreForTest(t)

			for _, project := range tt.fields.projects {
				err := s.WriteEnvConfiguration(&project)
				if err != nil {
					t.Errorf("cannot write project env %v: %s", project, err)
				}
				err = s.WriteProjectEnvironment(project.Name, project.Name)
				if err != nil {
					t.Errorf("cannot write project info %v: %s", project, err)
				}
			}

			for _, org := range tt.fields.orgs {
				if err := s.WriteOrgInfo(&org); err != nil {
					t.Errorf("cannot write organization %v: %s", org, err)
				}
			}

			if err := s.RemoveAllKnownCredentials(); (err != nil) != tt.wantErr {
				t.Errorf("RemoveOrganization() error = %v, wantErr %v", err, tt.wantErr)
			}

			projects, err := s.ListConfigurations()
			if err != nil {
				t.Errorf("unexpected error listing projects: %s", err)
			}
			if len(projects) != 0 {
				t.Errorf("unexpected projects found: %q", strings.Join(projects, ", "))
			}

			orgs, err := s.ListOrgs()
			if err != nil {
				t.Errorf("unexpected error listing orgs: %s", err)
			}
			if len(orgs) != 0 {
				t.Errorf("unexpected organizations found: %q", strings.Join(orgs, ", "))
			}
		})
	}
}

func TestStore_UpsertEnvConfig(t *testing.T) {
	store := newStoreForTest(t)
	envName := "test-env"
	alias := "default"
	apiKey := "test-api-key"

	// 1. Upsert should create a new configuration if it doesn't exist
	if err := store.UpsertEnvConfig(envName, alias, apiKey); err != nil {
		t.Fatalf("UpsertEnvConfig returned unexpected error: %v", err)
	}

	config, err := store.GetEnvConfiguration(envName)
	if err != nil {
		t.Fatalf("GetEnvConfiguration returned unexpected error: %v", err)
	}

	creds, err := config.GetCredentials(alias)
	if err != nil {
		t.Fatalf("GetCredentials returned unexpected error: %v", err)
	}

	if creds.APIKey != apiKey {
		t.Errorf("got APIKey %q, want %q", creds.APIKey, apiKey)
	}

	// 2. Upsert should update existing configuration
	newAPIKey := "new-test-api-key"
	if err := store.UpsertEnvConfig(envName, alias, newAPIKey); err != nil {
		t.Fatalf("UpsertEnvConfig returned unexpected error: %v", err)
	}

	config, err = store.GetEnvConfiguration(envName)
	if err != nil {
		t.Fatalf("GetEnvConfiguration returned unexpected error: %v", err)
	}

	creds, err = config.GetCredentials(alias)
	if err != nil {
		t.Fatalf("GetCredentials returned unexpected error: %v", err)
	}

	if creds.APIKey != newAPIKey {
		t.Errorf("got APIKey %q, want %q", creds.APIKey, newAPIKey)
	}
}
