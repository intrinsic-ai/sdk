// Copyright 2023 Intrinsic Innovation LLC

// Package auth manages API keys in a local config files.
package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"intrinsic/config/environments"

	log "github.com/golang/glog"
	"google.golang.org/grpc/metadata"
)

const (

	// AliasDefaultToken is the alias under which the default token is stored.
	AliasDefaultToken = "default"

	projectStoreDirectory = "intrinsic/projects"
	orgStoreDirectory     = "intrinsic/organizations"
	envStoreDirectory     = "intrinsic/environments"
	authConfigExtension   = ".user-token"
	envJSONExtension      = ".env.json"

	directoryMode  os.FileMode = 0o700
	fileMode       os.FileMode = 0o600
	writeFileFlags             = os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	// tokenExchangeServer is the address of the token exchange server.
	// Points to the GRPC gateway.
	tokenExchangeServer = "flowstate.intrinsic.ai"

	// envInDebugAuthStore is the environment variable used to trigger debug logging when reading from the auth store.
	envInDebugAuthStore = "INDEBUG_AUTHSTORE"
)

var debugAuthStore, _ = strconv.ParseBool(os.Getenv(envInDebugAuthStore))

const (
	// KeyProject is used as central flag name for passing a project name to inctl.
	KeyProject = "project"
	// KeyOrganization is used as central flag name for passing an organization name to inctl.
	KeyOrganization = "org"
	// KeyCluster is used as central flag name for passing a cluster name to inctl.
	KeyCluster = "cluster"
	// KeyEnvironment is used as central flag name for passing an environment name to inctl.
	//
	// The environment can be one of prod, staging or dev.
	KeyEnvironment = "env"
)

// RFC3339Time is type alias to correct (un)marshaling time.Time in RFC3339 format
type RFC3339Time time.Time

func (t *RFC3339Time) String() string {
	return time.Time(*t).Format(time.RFC3339)
}

// UnmarshalText implements encoding.TextUnmarshaler interface to unmarshal
// from RFC3339 formatted timestamp
func (t *RFC3339Time) UnmarshalText(text []byte) error {
	tx, err := time.Parse(time.RFC3339, string(text))
	if err != nil {
		return err
	}
	*t = RFC3339Time(tx)
	return nil
}

// MarshalText implements encoding.TextMarshaler interface to output
// timestamp in RFC3339 format.
func (t *RFC3339Time) MarshalText() (text []byte, err error) {
	return []byte((time.Time(*t)).Format(time.RFC3339)), nil
}

// OrgInfo encapsulates the information needed to use organization names in inctl.
type OrgInfo struct {
	Organization string `json:"org"`
	Project      string `json:"project"`
}

// ProjectToken represents cloud project bound API Token for user authorization
type ProjectToken struct {
	APIKey string `json:"apiKey"`
	// Reserved for future use, as we do not have this information yet
	ValidUntil *RFC3339Time `json:"validUntil,omitempty"`
}

// Validate performs validation of API Token
func (p *ProjectToken) Validate() error {
	if p == nil {
		return fmt.Errorf("nil token")
	}
	if p.ValidUntil != nil {
		if time.Now().After(time.Time(*p.ValidUntil)) {
			return fmt.Errorf("project token expired: %s", p.ValidUntil)
		}
	}
	if p.APIKey == "" {
		return fmt.Errorf("missing api key")
	}

	return nil
}

func (p *ProjectToken) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", p.APIKey),
	}, p.Validate()
}

// RequireTransportSecurity always return true to protect credentials
func (p *ProjectToken) RequireTransportSecurity() bool {
	return true
}

// HTTPAuthorization sets Authorization header to given request using project token.
func (p *ProjectToken) HTTPAuthorization(req *http.Request) (*http.Request, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.APIKey))
	return req, p.Validate()
}

type asIDTokenCredentialsOptions struct {
	apiKeyTokenSourceOptions    []APIKeyTokenSourceOption
	tokenExchangeServiceAddress string
}

// AsIDTokenCredentialsOption is a functional option for [ProjectToken.AsIDTokenCredentials].
type AsIDTokenCredentialsOption = func(o *asIDTokenCredentialsOptions)

// WithAPIKeyTokenSourceOptions adds options for creating an [APIKeyTokenSource].
func WithAPIKeyTokenSourceOptions(options ...APIKeyTokenSourceOption) AsIDTokenCredentialsOption {
	return func(opts *asIDTokenCredentialsOptions) {
		opts.apiKeyTokenSourceOptions = append(opts.apiKeyTokenSourceOptions, options...)
	}
}

// WithTokenExchangeServiceAddress sets the address of the token exchange service.
func WithTokenExchangeServiceAddress(address string) AsIDTokenCredentialsOption {
	return func(opts *asIDTokenCredentialsOptions) {
		opts.tokenExchangeServiceAddress = address
	}
}

// AsIDTokenCredentials allows converting Intrinsic API Tokens to Google ID Tokens
// on the fly as [credentials.PerRPCCredentials] implementation.
// This is useful for contacting services which don't accept Intrinsic API Tokens,
// but we want to use this infrastructure to authorize users to them.
func (p *ProjectToken) AsIDTokenCredentials(options ...AsIDTokenCredentialsOption) (*APIKeyTokenSource, error) {
	opts := &asIDTokenCredentialsOptions{
		tokenExchangeServiceAddress: tokenExchangeServer,
	}
	for _, opt := range options {
		opt(opts)
	}

	tsc, err := NewTokensServiceClient(http.DefaultClient, opts.tokenExchangeServiceAddress)
	if err != nil {
		return nil, fmt.Errorf("cannot create token exchange: %w", err)
	}
	return NewAPIKeyTokenSource(p.APIKey, tsc, opts.apiKeyTokenSourceOptions...), nil
}

// ProjectConfiguration contains list of API tokens related to given project
type ProjectConfiguration struct {
	Name string `json:"name"`
	// Tokens map individual API tokens for given project.
	// It is a map of alias: {api_key...}
	Tokens map[string]*ProjectToken `json:"tokens,omitempty"`

	// LastUpdated tracks when the file was last written by store, may be omitted
	LastUpdated *RFC3339Time `json:"lastUpdated,omitempty"`
}

// SetCredentials sets given apiKey to given alias in project configuration and optionally setting validity period.
func (p *ProjectConfiguration) SetCredentials(alias string, apiKey string, validUntil ...time.Time) (*ProjectConfiguration, error) {
	token := &ProjectToken{
		APIKey: apiKey,
	}
	if len(validUntil) > 0 && !validUntil[0].IsZero() {
		expires := RFC3339Time(validUntil[0])
		token.ValidUntil = &expires
	}
	p.Tokens[alias] = token

	return p, token.Validate()
}

// SetDefaultCredentials sets given apiKey to default alias, optionally setting validity period
func (p *ProjectConfiguration) SetDefaultCredentials(apiKey string, validUntil ...time.Time) (*ProjectConfiguration, error) {
	return p.SetCredentials(AliasDefaultToken, apiKey, validUntil...)
}

// HasCredentials checks if given project configuration has apiKey assigned to given alias.
func (p *ProjectConfiguration) HasCredentials(alias string) bool {
	_, ok := p.Tokens[alias]
	return ok
}

// GetCredentials returns ProjectToken object assigned to given alias or error if
// alias was not found.
func (p *ProjectConfiguration) GetCredentials(alias string) (*ProjectToken, error) {
	token, ok := p.Tokens[alias]
	if !ok {
		return nil, fmt.Errorf("token with alias '%s' not found", alias)
	}
	return token, nil
}

// GetDefaultCredentials returns ProjectToken object assigned to default alias,
// or error if not found.
func (p *ProjectConfiguration) GetDefaultCredentials() (*ProjectToken, error) {
	return p.GetCredentials(AliasDefaultToken)
}

// Store provides access to a collection of environment configurations stored as
// files in the users config directory.
type Store struct {
	EnvironmentResolver EnvironmentResolver

	// GetConfigDirFx is an indirection allowing to use custom config dirs in tests.
	GetConfigDirFx func() (string, error)
}

// DefaultStore is default instance of [Store]
var DefaultStore = NewStore()

// NewStore returns a new Store instance.
func NewStore() *Store {
	return &Store{
		EnvironmentResolver: &defaultEnvironmentResolver{},
	}
}

func (s *Store) getConfigDir() (string, error) {
	if s.GetConfigDirFx == nil {
		return os.UserConfigDir()
	}
	return s.GetConfigDirFx()
}

// NewConfiguration returns a new, empty ProjectConfiguration for the given
// project name.
func NewConfiguration(name string) *ProjectConfiguration {
	return &ProjectConfiguration{
		Name:        name,
		Tokens:      make(map[string]*ProjectToken, 0),
		LastUpdated: nil,
	}
}

// environments returns a list of environments
// for which a user is logged into (i.e. has a configuration in the local auth storage).
func (s *Store) environments() ([]string, error) {
	storeLocation, err := s.envStoreLocation()
	if err != nil {
		return nil, fmt.Errorf("cannot find env store: %w", err)
	}

	globPattern := filepath.Join(storeLocation, "*"+authConfigExtension)
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid glob pattern: %w", err)
	}

	envs := make([]string, len(matches))
	for i, match := range matches {
		filename := filepath.Base(match)
		envs[i] = strings.TrimSuffix(filename, authConfigExtension)
	}

	return envs, nil
}

// WriteProjectEnvironment caches the environment of the project,
// so that a user doesn't have to request it again.
// A project doesn't change its environment.
func (s *Store) WriteProjectEnvironment(project, env string) error {
	filename, err := s.projectToEnvFilename(project)
	if err != nil {
		return fmt.Errorf("project to env filename: %w", err)
	}

	if err = os.MkdirAll(filepath.Dir(filename), directoryMode); err != nil {
		return fmt.Errorf("create target directory: %w", err)
	}

	authDebug("AUTHSTORE_WRITE_ATTEMPT=%s", filename)
	f, err := os.OpenFile(filename, writeFileFlags, fileMode)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	defer f.Close()

	if err := json.NewEncoder(f).Encode(&projectToEnv{
		Environment: env,
	}); err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}

// GetConfiguration checks the local auth storage for the environment matching the project
// and returns the token for this environment. The token works for any project
// in the matching environment a user identity has access to.
// If there's no data about the matching environment,
// it tries to guess it from the environments a user is logged into.
// If none of them match the project, a user is asked to relogin.
func (s *Store) GetConfiguration(projectName string) (*ProjectConfiguration, error) {
	env, err := s.projectEnvironment(projectName)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("project environment: %w", err)
		}

		env, err = s.discoverProjectEnvironment(projectName)
		if err != nil {
			return nil, fmt.Errorf("discover project environment: %w", err)
		}
	}

	cfg, err := s.GetEnvConfiguration(env)
	if err != nil {
		return nil, fmt.Errorf("get env configuration: %w; please run `inctl auth login --org=<your_org>@%s`", err, projectName)
	}

	return cfg, nil
}

// discoverProjectEnvironment checks if a project belongs
// to one of environments a user is logged into and caches this data.
// If the project doesn't belong to any of them, a user is asked to login.
func (s *Store) discoverProjectEnvironment(projectName string) (string, error) {
	envs, err := s.environments()
	if err != nil {
		return "", fmt.Errorf("get user envs: %w", err)
	}

	var checkErrs []error
	for _, assumedEnv := range envs {
		ok, err := s.EnvironmentResolver.HasProject(assumedEnv, projectName)
		if err != nil {
			checkErrs = append(checkErrs, fmt.Errorf("check project env %q: %w", assumedEnv, err))
			continue
		}

		if ok {
			if err := s.WriteProjectEnvironment(projectName, assumedEnv); err != nil {
				return "", fmt.Errorf("write project environment: %w", err)
			}
			return assumedEnv, nil
		}
	}

	errMsg := "project %q is not found in any of your environments, please run `inctl auth login --org=<your_org>@%s`"
	errArgs := []any{projectName, projectName}
	if len(checkErrs) > 0 {
		errMsg += " (also encountered check errors: %v)"
		errArgs = append(errArgs, errors.Join(checkErrs...))
	}
	return "", fmt.Errorf(errMsg, errArgs...)
}

func authDebug(format string, args ...any) {
	if debugAuthStore {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}

func readAndDecodeFile(filename string, v any) error {
	authDebug("AUTHSTORE_READ_ATTEMPT=%s", filename)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(v); err != nil {
		// remove it and treat it as non-existent if it's corrupted.
		_ = file.Close()
		authDebug("FILE IS CORRUPTED, AUTHSTORE_DELETE_ATTEMPT=%s", filename)
		_ = os.Remove(filename)
		return os.ErrNotExist
	}

	return nil
}

func readConfigurationFromFile(filename string) (*ProjectConfiguration, error) {
	result := &ProjectConfiguration{}
	if err := readAndDecodeFile(filename, result); err != nil {
		return nil, fmt.Errorf("cannot open configuration file: %w", err)
	}

	// ensure that tokens are always populated
	if result.Tokens == nil {
		result.Tokens = map[string]*ProjectToken{}
	}

	return result, nil
}

func writeConfigToFile(config *ProjectConfiguration, filename string) error {
	// we make sure we have whole directory structure before we create file.
	// os.MkdirAll() calls os.Stat() on path, so there is no point to do it here.
	if err := os.MkdirAll(filepath.Dir(filename), directoryMode); err != nil {
		return fmt.Errorf("cannot create target directory: %w", err)
	}

	authDebug("AUTHSTORE_WRITE_ATTEMPT=%s", filename)
	file, err := os.OpenFile(filename, writeFileFlags, fileMode)
	if err != nil {
		return fmt.Errorf("cannot open configuration file: %w", err)
	}

	defer file.Close()

	// update last modified in UTC time
	now := RFC3339Time(time.Now().UTC())
	config.LastUpdated = &now

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("cannot serialize configuration: %w", err)
	}

	// if sync fails, we did not write into store.
	return file.Sync()
}

// ListConfigurations gives a list of known configurations. It works on
// filesystem level and does not attempt to read the content of configuration.
// Membership in this list does not guarantee valid configuration for given name
// exists. Results is not sorted and is returned in same order as blobbed on
// filesystem.
func (s *Store) ListConfigurations() ([]string, error) {
	storeLocation, err := s.projectStoreLocation()
	if err != nil {
		return nil, fmt.Errorf("cannot find configuration store: %w", err)
	}

	globPattern := filepath.Join(storeLocation, "*"+envJSONExtension)
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		panic(fmt.Errorf("invalid glob pattern, programmer error: %w", err))
	}
	if len(matches) == 0 {
		// this is valid response, there are no projects found.
		return nil, nil
	}
	result := make([]string, 0, len(matches))
	for _, match := range matches {
		filename := filepath.Base(match)
		result = append(result, strings.TrimSuffix(filename, envJSONExtension))
	}

	return result, nil
}

// RemoveConfiguration removes the stored configuration for the given project
// name. Returns nil if no such configuration exists.
func (s *Store) RemoveConfiguration(name string) error {
	filename, err := s.projectToEnvFilename(name)
	if err != nil {
		return fmt.Errorf("cannot remove configuration: %w", err)
	}
	authDebug("AUTHSTORE_DELETE_ATTEMPT=%s", filename)
	return os.Remove(filename)
}

// AuthorizeContext retrieves the default credentials for the given project and adds authorization
// information directly to a context derived from the given context. If the given context already
// has authorization information this function will *not* modify the context.
//
// Always prefer using [ProjectToken] as per-RPC credentials where possible. This is a fallback that
// enables passing API keys over insecure boundaries and must be used carefully.
//
// Warning: This writes an "authorization" header (outgoing metadata) to the context. Only one such
// header is permitted on requests. A context authorized with this method must never be used with a
// connection that specifies per-RPC credentials that also write to the "authorization" header.
// Prominent examples of per-RPC credentials that must not be used with this method are
// [oauth.TokenSource] (used for default credentials) and [ProjectToken].
func (s *Store) AuthorizeContext(ctx context.Context, projectName string) (context.Context, error) {
	configuration, err := s.GetConfiguration(projectName)
	if err != nil {
		return ctx, fmt.Errorf("cannot get configuration: %w", err)
	}
	pt, err := configuration.GetDefaultCredentials()
	if err != nil {
		return ctx, fmt.Errorf("cannot get default credentials: %w", err)
	}
	if err := pt.Validate(); err != nil {
		return ctx, fmt.Errorf("invalid credentials: %w", err)
	}
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok && len(md.Get("authorization")) > 0 {
		return ctx, nil
	}
	return metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", pt.APIKey)), nil
}

func (s *Store) orgStoreLocation() (string, error) {
	configDir, err := s.getConfigDir()
	if err != nil {
		return "", fmt.Errorf("get config directory: %w", err)
	}

	return filepath.Join(configDir, orgStoreDirectory), nil
}

func (s *Store) orgFilename(name string) (string, error) {
	orgDir, err := s.orgStoreLocation()
	if err != nil {
		return "", fmt.Errorf("get config directory: %w", err)
	}

	return filepath.Join(orgDir, fmt.Sprintf("%s.json", name)), nil
}

// WriteOrgInfo writes the information we have about an org to file
func (s *Store) WriteOrgInfo(o *OrgInfo) error {
	filename, err := s.orgFilename(o.Organization)
	if err != nil {
		return err
	}

	// we make sure we have whole directory structure before we create file.
	// os.MkdirAll() calls os.Stat() on path, so there is no point to do it here.
	if err = os.MkdirAll(filepath.Dir(filename), directoryMode); err != nil {
		return fmt.Errorf("create target directory: %w", err)
	}

	authDebug("AUTHSTORE_WRITE_ATTEMPT=%s", filename)
	file, err := os.OpenFile(filename, writeFileFlags, fileMode)
	if err != nil {
		return fmt.Errorf("open configuration file: %w", err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(o); err != nil {
		return fmt.Errorf("serialize configuration: %w", err)
	}

	// if sync fails, we did not write into store.
	return file.Sync()
}

// ReadOrgInfo reads the information about an organization previously written to the auth store.
func (s *Store) ReadOrgInfo(orgName string) (OrgInfo, error) {
	filename, err := s.orgFilename(orgName)
	if err != nil {
		return OrgInfo{}, err
	}

	authDebug("AUTHSTORE_READ_ATTEMPT=%s", filename)

	file, err := os.Open(filename)
	if err != nil {
		return OrgInfo{}, fmt.Errorf("open configuration: %w", err)
	}
	defer file.Close()

	ret := OrgInfo{}
	if err := json.NewDecoder(file).Decode(&ret); err != nil {
		return OrgInfo{}, fmt.Errorf("deserialize configuration: %w", err)
	}

	return ret, nil
}

// ListOrgs gives a list of known organizations. It works on
// filesystem level and does not attempt to read the content of configuration.
// Results are not sorted and the order may change at any time.
func (s *Store) ListOrgs() ([]string, error) {
	storeLocation, err := s.orgStoreLocation()
	if err != nil {
		return nil, fmt.Errorf("cannot find configuration store: %w", err)
	}

	globPattern := filepath.Join(storeLocation, "*.json")
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		panic(fmt.Errorf("invalid glob pattern, programmer error: %w", err))
	}
	if len(matches) == 0 {
		// this is valid response, there are no organizations found.
		return nil, nil
	}

	result := make([]string, 0, len(matches))
	for _, match := range matches {
		filename := filepath.Base(match)
		result = append(result, strings.TrimSuffix(filename, ".json"))
	}

	return result, nil
}

// RemoveOrganization removes the named organization from the store and its
// associated project if no other organization references it. Reference check
// is done to ensure users belonging to multiple organizations in single
// multi-tenancy project don't lose project credentials when removing
// only one organization.
func (s *Store) RemoveOrganization(name string) error {
	info, err := s.ReadOrgInfo(name)
	if err != nil {
		return fmt.Errorf("organization not found: %w", err)
	}

	associatedProject := info.Project
	deleteProject := true

	orgs, err := s.ListOrgs()
	if err != nil {
		return err
	}

	for _, orgName := range orgs {
		if orgName == name {
			continue
		}
		orgInfo, err := s.ReadOrgInfo(orgName)
		if err != nil {
			// we should not hit this unless some "data race" happened
			continue
		}
		if orgInfo.Project == associatedProject {
			deleteProject = false
		}
	}

	filename, err := s.orgFilename(name)
	if err != nil {
		return err
	}
	err = os.Remove(filename)
	if err != nil {
		return fmt.Errorf("cannot remove organization: %w", err)
	}

	if deleteProject {
		// we are going to delete project only if there is only one organization
		// using it. If there are more than one, we are leaving project intact
		err = s.RemoveConfiguration(associatedProject)
		if err != nil {
			// failure to remove project is not considered failure to remove org.
			log.Warningf("cannot remove project configuration: %s", err)
		}
	}

	return nil
}

// RemoveAllKnownCredentials removes all known organizations, environments, and projects
// from authorization store. It operates on filesystem and does not attempt
// to read credentials. Use for full removal of credentials.
func (s *Store) RemoveAllKnownCredentials() error {
	location, err := s.projectStoreLocation()
	if err != nil {
		return err
	}
	err = filepath.WalkDir(location, s.deleteFiles)
	if err != nil {
		return err
	}
	location, err = s.orgStoreLocation()
	if err != nil {
		return err
	}

	err = filepath.WalkDir(location, s.deleteFiles)
	if err != nil {
		return err
	}

	location, err = s.envStoreLocation()
	if err != nil {
		return err
	}

	return filepath.WalkDir(location, s.deleteFiles)
}

func (s *Store) deleteFiles(path string, de fs.DirEntry, err error) error {
	if err != nil {
		if de == nil || de.IsDir() {
			return fs.SkipDir
		}
		return err
	}

	// we are retaining directory structure.
	if !de.IsDir() {
		authDebug("AUTHSTORE_DELETE_ATTEMPT=%s", path)
		if err = os.Remove(path); err != nil {
			// if we fail to remove a file, we just move on.
			log.Warningf("cannot remove %s: %s", path, err)
		}
	}
	return nil
}

func (s *Store) GetEnvConfiguration(env string) (*ProjectConfiguration, error) {
	filename, err := s.getEnvConfigurationFilename(env)
	if err != nil {
		return nil, fmt.Errorf("cannot open configuration for environment %q: %w", env, err)
	}

	cfg, err := readConfigurationFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read env configuration: %w", err)
	}

	pToken, err := cfg.GetDefaultCredentials()
	if err != nil {
		return nil, fmt.Errorf("get default credentials: %w", err)
	}

	if err := pToken.Validate(); err != nil {
		return nil, fmt.Errorf("validate credentials: %w", err)
	}

	return cfg, nil
}

func (s *Store) getEnvConfigurationFilename(env string) (string, error) {
	if env == "" {
		return "", fmt.Errorf("environment name is required")
	}

	storeDir, err := s.envStoreLocation()
	if err != nil {
		return "", fmt.Errorf("cannot find env configurations: %w", err)
	}

	envFilename := env + authConfigExtension
	return filepath.Join(storeDir, envFilename), nil
}

func (s *Store) envStoreLocation() (string, error) {
	configDir, err := s.getConfigDir()
	return filepath.Join(configDir, envStoreDirectory), err
}

// WriteEnvConfiguration writes the given environment configuration to the store.
func (s *Store) WriteEnvConfiguration(config *ProjectConfiguration) error {
	filename, err := s.getEnvConfigurationFilename(config.Name)
	if err != nil {
		return fmt.Errorf("error getting env config filename: %w", err)
	}
	if err := writeConfigToFile(config, filename); err != nil {
		return fmt.Errorf("failed to write environment configuration: %w", err)
	}
	return nil
}

// UpsertEnvConfig sets the given apiKey to the given alias in the environment
// configuration and writes it to the store.
func (s *Store) UpsertEnvConfig(envName, alias, apiKey string) error {
	filename, err := s.getEnvConfigurationFilename(envName)
	if err != nil {
		return fmt.Errorf("cannot get env filename: %w", err)
	}

	config, err := readConfigurationFromFile(filename)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("cannot load '%s' configuration: %w", envName, err)
		}
		config = NewConfiguration(envName)
	}

	config, err = config.SetCredentials(alias, apiKey)
	if err != nil {
		return fmt.Errorf("aborting, invalid credentials: %w", err)
	}

	if err = s.WriteEnvConfiguration(config); err != nil {
		return fmt.Errorf("error writing env config: %w", err)
	}

	return nil
}

type projectToEnv struct {
	Environment string `json:"environment"`
}

// projectEnvironment returns the cached environment value from
// $XDG_CONFIG_HOME/.intrinsic/projectStoreDirectory/<project>.env.json
func (s *Store) projectEnvironment(project string) (string, error) {
	filename, err := s.projectToEnvFilename(project)
	if err != nil {
		return "", fmt.Errorf("project to env filename: %w", err)
	}

	var envData projectToEnv
	if err := readAndDecodeFile(filename, &envData); err != nil {
		return "", fmt.Errorf("open: %w", err)
	}

	return envData.Environment, nil
}

func (s *Store) projectToEnvFilename(project string) (string, error) {
	if project == "" {
		return "", fmt.Errorf("empty project name")
	}

	storeDir, err := s.projectStoreLocation()
	if err != nil {
		return "", fmt.Errorf("cannot find configurations: %w", err)
	}
	projectToEnvFilename := project + envJSONExtension
	return filepath.Join(storeDir, projectToEnvFilename), nil
}

// EnvironmentResolver checks if a project
// belongs to an environment.
type EnvironmentResolver interface {
	HasProject(env, project string) (bool, error)
}

type defaultEnvironmentResolver struct{}

func (r *defaultEnvironmentResolver) HasProject(env, project string) (bool, error) {
	return environments.FromAnyProject(project) == env, nil
}

func (s *Store) projectStoreLocation() (string, error) {
	configDir, err := s.getConfigDir()
	return filepath.Join(configDir, projectStoreDirectory), err
}
