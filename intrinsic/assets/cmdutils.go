// Copyright 2023 Intrinsic Innovation LLC

// Package cmdutils provides utils for asset inctl commands.
package cmdutils

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
	"intrinsic/assets/imagetransfer"
	"intrinsic/assets/imageutils"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
	iapb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
	viewpb "intrinsic/assets/proto/view_go_proto"
	"intrinsic/assets/typeutils"
	"intrinsic/assets/viewutils"
	"intrinsic/skills/tools/resource/cmd/bundleimages"
	"intrinsic/tools/inctl/util/orgutil"
)

const (
	// KeyAddress is the name of the address flag.
	KeyAddress = "address"
	// KeyAPIKey is the name of the arg to specify the api-key to use.
	KeyAPIKey = "api_key"
	// KeyAssetTypes is the name of the asset types flag.
	KeyAssetTypes = "asset_types"
	// KeyAuthUser is the name of the auth user flag.
	KeyAuthUser = "auth_user"
	// KeyAuthPassword is the name of the auth password flag.
	KeyAuthPassword = "auth_password"
	// KeyCatalogAddress is the name of the catalog address flag.
	KeyCatalogAddress = "catalog_address"
	// KeyCluster is the name of the cluster flag.
	KeyCluster = "cluster"
	// KeyContext is the name of the context flag.
	KeyContext = "context"
	// KeyDefault is the name of the default flag.
	KeyDefault = "default"
	// KeyDryRun is the name of the dry run flag.
	KeyDryRun = "dry_run"
	// KeyFilter is the name of the filter flag.
	KeyFilter = "filter"
	// KeyIgnoreExisting is the name of the flag to ignore AlreadyExists errors.
	KeyIgnoreExisting = "ignore_existing"
	// KeyManifestFile is the file path to the manifest binary.
	KeyManifestFile = "manifest_file"
	// KeyManifestTarget is the build target to the skill manifest.
	KeyManifestTarget = "manifest_target"
	// KeyOrgPrivate is the name of the org-private flag.
	KeyOrgPrivate = "org_private"
	// KeyOrganization is used as central flag name for passing an organization name to inctl.
	KeyOrganization = orgutil.KeyOrganization
	// KeyPolicy defines the flag used to specify the policy option when
	// interacting with the installed asset service.
	KeyPolicy = "policy"
	// KeyProject is used as central flag name for passing a project name to inctl.
	KeyProject = orgutil.KeyProject
	// KeyProvides is the name of the provided interfaces flag.
	KeyProvides = "provides"
	// KeyRegistry is the name of the registry flag.
	KeyRegistry = "registry"
	// KeyReleaseNotes is the name of the release notes flag.
	KeyReleaseNotes = "release_notes"
	// KeySkipDirectUpload is boolean flag controlling direct upload behavior
	KeySkipDirectUpload = "skip_direct_upload"
	// KeySkipPrompts is the name of the flag to skip user prompts.
	KeySkipPrompts = "skip_prompts"
	// KeySolution is the name of the solution flag.
	KeySolution = "solution"
	// KeyType is the name of the type flag.
	KeyType = "type"
	// KeyTimeout is the name of the timeout flag.
	KeyTimeout = "timeout"
	// KeyUseBorgCredentials is the name of the flag to use borg credentials.
	KeyUseBorgCredentials = "use_borg_credentials"
	// KeyVendor is the name of the vendor flag.
	KeyVendor = "vendor"
	// KeyVersion is the name of the version flag.
	KeyVersion = "version"
	// KeyView is the name of the view flag.
	KeyView = "view"

	envPrefix = "intrinsic"
)

var (
	policyMap = map[string]iapb.UpdatePolicy{
		"":                  iapb.UpdatePolicy_UPDATE_POLICY_UNSPECIFIED,
		"add_new_only":      iapb.UpdatePolicy_UPDATE_POLICY_ADD_NEW_ONLY,
		"update_unused":     iapb.UpdatePolicy_UPDATE_POLICY_UPDATE_UNUSED,
		"update_compatible": iapb.UpdatePolicy_UPDATE_POLICY_UPDATE_COMPATIBLE,
	}
)

// CmdFlags abstracts interaction with inctl command flags.
type CmdFlags struct {
	cmd        *cobra.Command
	viperLocal *viper.Viper
}

// NewCmdFlags returns a new CmdFlags instance.
func NewCmdFlags() *CmdFlags {
	viperLocal := viper.New()
	viperLocal.SetEnvPrefix(envPrefix)

	return NewCmdFlagsWithViper(viperLocal)
}

// NewCmdFlagsWithViper returns a new CmdFlags instance with a custom Viper.
func NewCmdFlagsWithViper(viperLocal *viper.Viper) *CmdFlags {
	return &CmdFlags{cmd: nil, viperLocal: viperLocal}
}

// SetCommand sets the cobra Command to interact with.
//
// The command must be set before any flags are added.
func (cf *CmdFlags) SetCommand(cmd *cobra.Command) {
	cf.cmd = cmd
}

// AddFlagAssetTypes adds a flag for an optional list of asset types.
func (cf *CmdFlags) AddFlagAssetTypes(defaultTypes string) {
	types := typeutils.AllAssetTypes()
	names := make([]string, len(types))
	for i, t := range types {
		names[i] = typeutils.NameFromAssetType(t)
	}

	cf.OptionalString(KeyAssetTypes, defaultTypes, fmt.Sprintf("A comma-separated list of asset types (choose from: %v).", strings.Join(names, ", ")))
}

// GetFlagAssetTypes gets the (enum) values of the asset types flag added by AddFlagAssetTypes.
func (cf *CmdFlags) GetFlagAssetTypes() ([]atypepb.AssetType, error) {
	assetTypesValue := cf.GetString(KeyAssetTypes)
	if assetTypesValue == "" {
		return nil, nil
	}
	assetTypeNames := strings.Split(assetTypesValue, ",")
	assetTypes := make([]atypepb.AssetType, len(assetTypeNames))
	for i, name := range assetTypeNames {
		var err error
		assetTypes[i], err = typeutils.AssetTypeFromName(name)
		if err != nil {
			return nil, err
		}
	}

	return assetTypes, nil
}

// AddFlagsCredentials adds args for specifying credentials.
func (cf *CmdFlags) AddFlagsCredentials() {
	cf.OptionalBool(KeyUseBorgCredentials, false, "Use credentials associated with the current borg user, rather than application-default credentials.")
	cf.OptionalString(KeyAPIKey, "", "The API key to use for authentication.")
}

// GetFlagsCredentials gets the values of the credential args.
func (cf *CmdFlags) GetFlagsCredentials() (useBorgCredentials bool, apiKey string) {
	useBorgCredentials = cf.GetBool(KeyUseBorgCredentials)
	apiKey = cf.GetString(KeyAPIKey)

	return useBorgCredentials, apiKey
}

// AddFlagDefault adds a flag for marking a released asset as default.
func (cf *CmdFlags) AddFlagDefault(assetType string) {
	cf.OptionalBool(KeyDefault, false, fmt.Sprintf("Whether this %s version should be tagged as the default.", assetType))
}

// GetFlagDefault gets the value of the default flag added by AddFlagDefault.
func (cf *CmdFlags) GetFlagDefault() bool {
	return cf.GetBool(KeyDefault)
}

// GetFlagDefaultIsSet returns whether the default flag was set.
func (cf *CmdFlags) GetFlagDefaultIsSet() bool {
	return cf.viperLocal.IsSet(KeyDefault)
}

// AddFlagDryRun adds a flag for performing a dry run.
func (cf *CmdFlags) AddFlagDryRun() {
	cf.OptionalBool(KeyDryRun, false, "Dry-run by validating but not performing any actions.")
}

// GetFlagDryRun gets the value of the dry run flag added by AddFlagDryRun.
func (cf *CmdFlags) GetFlagDryRun() bool {
	return cf.GetBool(KeyDryRun)
}

// AddFlagIgnoreExisting adds a flag to ignore AlreadyExists errors.
func (cf *CmdFlags) AddFlagIgnoreExisting(assetType string) {
	cf.OptionalBool(KeyIgnoreExisting, false, fmt.Sprintf("Ignore errors if the specified %s version already exists in the catalog.", assetType))
}

// GetFlagIgnoreExisting gets the value of the flag added by AddFlagIgnoreExisting.
func (cf *CmdFlags) GetFlagIgnoreExisting() bool {
	return cf.GetBool(KeyIgnoreExisting)
}

// AddFlagAddress adds a flag for the installer service address.
func (cf *CmdFlags) AddFlagAddress() {
	cf.OptionalEnvString(KeyAddress, "xfa.lan:17080", `The address of the cluster.
When not running the cluster on localhost, this should be the address of the relay
(e.g.: dns:///www.endpoints.<gcp_project_name>.cloud.goog:443).`)
}

// GetFlagAddress gets the value of the cluster address flag added by
// AddFlagAddress.
func (cf *CmdFlags) GetFlagAddress() string {
	return cf.GetString(KeyAddress)
}

// AddFlagsAddressClusterSolution adds flags for the address, cluster, and solution when installing
// or working with installed assets.
func (cf *CmdFlags) AddFlagsAddressClusterSolution() {
	cf.OptionalString(KeyAddress, "", "Internal flag to directly set the API server address. Normally, you should use --org instead, which tells inctl to connect via the cloud.")
	cf.OptionalEnvString(KeyCluster, "", "The target Kubernetes cluster ID. If you set this, you must not set --solution.")
	cf.OptionalEnvString(KeySolution, "", "The target solution. Must be running. If you set this, you must not set --cluster.")

	cf.cmd.MarkFlagsMutuallyExclusive(KeyCluster, KeySolution)
}

// GetFlagsAddressClusterSolution gets the values of the address, cluster, and solution flags added
// by AddFlagsAddressClusterSolution.
func (cf *CmdFlags) GetFlagsAddressClusterSolution() (string, string, string, error) {
	address := cf.GetString(KeyAddress)
	cluster := cf.GetString(KeyCluster)
	solution := cf.GetString(KeySolution)

	if address == "" && cluster == "" && solution == "" {
		return "", "", "", fmt.Errorf("at least one of `--%s`, `--%s` or `--%s` must be set", KeyAddress, KeyCluster, KeySolution)
	}
	// This matches these flags being marked as mutually exclusive above.  That
	// does not prevent two environment variables being provided or a
	// combination of flags and variables.  This is a fairly strict check, as
	// we probably want these flags to functionally behave as one where we
	// autodetect the type.  Which could probably be done with a clear order of
	// precedence (address before cluster before solution) without needing to
	// autodetect the kind.  If this is too strict, then we can override the
	// check in clientutils that triggers a lookup if solution is set.
	if cluster != "" && solution != "" {
		return "", "", "", fmt.Errorf("both `--%s=%q` and `--%s=%q` were provided by a flags and/or environment variables, which could be ambiguous", KeyCluster, cluster, KeySolution, solution)
	}

	return address, cluster, solution, nil
}

// AddFlagsManifest adds flags for specifying a manifest.
func (cf *CmdFlags) AddFlagsManifest() {
	cf.OptionalString(KeyManifestFile, "", "The path to the manifest binary file.")
	cf.OptionalString(KeyManifestTarget, "", "The manifest bazel target.")

	cf.cmd.MarkFlagsMutuallyExclusive(KeyManifestFile, KeyManifestTarget)
}

// GetFlagsManifest gets the values of the manifest flags added by AddFlagsManifest.
func (cf *CmdFlags) GetFlagsManifest() (string, string, error) {
	mf := cf.GetString(KeyManifestFile)
	mt := cf.GetString(KeyManifestTarget)

	if mf == "" && mt == "" {
		return "", "", fmt.Errorf("one of --%s or --%s must be provided", KeyManifestFile, KeyManifestTarget)
	}

	return mf, mt, nil
}

// AddFlagOrgPrivate adds a flag for marking a released asset as private to the organization.
func (cf *CmdFlags) AddFlagOrgPrivate() {
	cf.OptionalBool(KeyOrgPrivate, false, "Whether this asset should be private to the organization that owns it.")
}

// GetFlagOrgPrivate gets the value of the org_private flag added by AddFlagOrgPrivate.
func (cf *CmdFlags) GetFlagOrgPrivate() bool {
	return cf.GetBool(KeyOrgPrivate)
}

// GetFlagOrgPrivateIsSet returns whether the org_private flag was set.
func (cf *CmdFlags) GetFlagOrgPrivateIsSet() bool {
	return cf.viperLocal.IsSet(KeyOrgPrivate)
}

// AddFlagPolicy adds a flag for the update policy.
func (cf *CmdFlags) AddFlagPolicy(assetType string) {
	cf.OptionalString(KeyPolicy, "", fmt.Sprintf("The update policy to be used to install the %s. Can be: %v", assetType, maps.Keys(policyMap)))
}

// GetFlagPolicy gets the value of the policy flag.  This flag must be added manually by the
// utility.
func (cf *CmdFlags) GetFlagPolicy() (iapb.UpdatePolicy, error) {
	policy := cf.GetString(KeyPolicy)
	if value, ok := policyMap[policy]; ok {
		return value, nil
	}

	return iapb.UpdatePolicy_UPDATE_POLICY_UNSPECIFIED, fmt.Errorf("%q provided for --%v is invalid; valid values are: %v", policy, KeyPolicy, maps.Keys(policyMap))
}

// AddFlagOrganizationOptional adds an optional flag for the organization.
func (cf *CmdFlags) AddFlagOrganizationOptional() {
	cf.OptionalEnvString(KeyOrganization, "",
		`The Intrinsic organization to use. You can set the environment variable
		INTRINSIC_ORG=organization to set a default organization.`)
}

// GetFlagOrganization gets the value of the organization flag added by AddFlagOrganizationOptional.
func (cf *CmdFlags) GetFlagOrganization() string {
	return cf.GetString(KeyOrganization)
}

// AddFlagsProjectOrg adds both the project and org flag, including the necessary handling.
func (cf *CmdFlags) AddFlagsProjectOrg(opts ...orgutil.WrapCmdOption) {
	// While WrapCmd returns the pointer to make it inline, it's modifying, so we can use it here.
	orgutil.WrapCmd(cf.cmd, cf.viperLocal, opts...)
}

// AddFlagsProjectOrgOptional adds both the project and org flag as optional, including the necessary handling.
func (cf *CmdFlags) AddFlagsProjectOrgOptional(opts ...orgutil.WrapCmdOption) {
	// While WrapCmd returns the pointer to make it inline, it's modifying, so we can use it here.
	orgutil.WrapCmdOptional(cf.cmd, cf.viperLocal, opts...)
}

// AddFlagProject adds a flag for the GCP project.
func (cf *CmdFlags) AddFlagProject() {
	cf.RequiredEnvString(KeyProject, "",
		`The Google Cloud Project (GCP) project to use. You can set the environment variable
		INTRINSIC_PROJECT=project_name to set a default project name.`)
}

// AddFlagCatalogProjectOptional adds an optional flag for the GCP project to use for the catalog.
func (cf *CmdFlags) AddFlagCatalogProjectOptional() {
	cf.OptionalEnvString(KeyProject, "",
		`The Google Cloud Project (GCP) project to use for the catalog. You can set the environment
		variable INTRINSIC_PROJECT=project_name to set a default project name.`)
}

// AddFlagProjectOptional adds an optional flag for the GCP project.
func (cf *CmdFlags) AddFlagProjectOptional() {
	cf.OptionalEnvString(KeyProject, "",
		`The Google Cloud Project (GCP) project to use. You can set the environment variable
		INTRINSIC_PROJECT=project_name to set a default project name.`)
}

// GetFlagProject gets the value of the project flag added by AddFlagProject.
func (cf *CmdFlags) GetFlagProject() string {
	return cf.GetString(KeyProject)
}

// AddFlagProvides adds a flag for specifying provided interfaces.
func (cf *CmdFlags) AddFlagProvides() {
	cf.OptionalString(KeyProvides, "", fmt.Sprintf("A comma-separated list of interfaces that assets must provide in order to be included in the output."))
}

// GetFlagProvides gets the value of the provides flag added by AddFlagProvides.
func (cf *CmdFlags) GetFlagProvides() ([]string, error) {
	providesValue := cf.GetString(KeyProvides)
	if providesValue == "" {
		return nil, nil
	}

	provides := strings.Split(providesValue, ",")
	for i, provide := range provides {
		provides[i] = strings.TrimSpace(provide)
	}

	return provides, nil
}

// AddFlagRegistry adds a flag for the registry when side-loading an asset.
func (cf *CmdFlags) AddFlagRegistry() {
	cf.OptionalEnvString(KeyRegistry, "", fmt.Sprint("The container registry address."))
}

// GetFlagRegistry gets the value of the registry flag added by AddFlagRegistry.
func (cf *CmdFlags) GetFlagRegistry() string {
	return cf.GetString(KeyRegistry)
}

// AddFlagsRegistryAuthUserPassword adds flags for user/password authentication for a private
// container registry.
func (cf *CmdFlags) AddFlagsRegistryAuthUserPassword() {
	cf.OptionalString(KeyAuthUser, "", "The username used to access the private container registry.")
	cf.OptionalString(KeyAuthPassword, "", "The password used to authenticate private container registry access.")
	cf.cmd.MarkFlagsRequiredTogether(KeyAuthUser, KeyAuthPassword)
}

// GetFlagsRegistryAuthUserPassword gets the values of the user/password flags added by
// AddFlagsRegistryAuthUserPassword
func (cf *CmdFlags) GetFlagsRegistryAuthUserPassword() (string, string) {
	return cf.GetString(KeyAuthUser), cf.GetString(KeyAuthPassword)
}

// AddFlagReleaseNotes adds a flag for release notes.
func (cf *CmdFlags) AddFlagReleaseNotes(assetType string) {
	cf.OptionalString(KeyReleaseNotes, "", fmt.Sprintf("Release notes for this version of the %s.", assetType))
}

// GetFlagReleaseNotes gets the value of the release notes flag added by AddFlagReleaseNotes.
func (cf *CmdFlags) GetFlagReleaseNotes() string {
	return cf.GetString(KeyReleaseNotes)
}

// AddFlagSkillReleaseType adds a flag for the type when releasing a skill.
func (cf *CmdFlags) AddFlagSkillReleaseType() {
	targetTypeDescriptions := []string{}

	targetTypeDescriptions = append(targetTypeDescriptions, `"archive": a file path to a skill bundle file.`)

	cf.OptionalString(
		KeyType,
		string(imageutils.Archive),
		fmt.Sprintf("The type of target that is being specified. One of the following:\n   %s", strings.Join(targetTypeDescriptions, "\n   ")),
	)
}

// GetFlagSkillReleaseType gets the value of the type flag added by AddFlagSkillReleaseType.
func (cf *CmdFlags) GetFlagSkillReleaseType() string {
	return cf.GetString(KeyType)
}

// AddFlagSideloadContext adds a flag for the context when side-loading an asset.
func (cf *CmdFlags) AddFlagSideloadContext() {
	cf.OptionalEnvString(KeyContext, "", fmt.Sprintf("The Kubernetes cluster to use. Required unless using localhost for %s.", KeyAddress))
}

// GetFlagSideloadContext gets the value of the context flag added by AddFlagSideloadContext.
func (cf *CmdFlags) GetFlagSideloadContext() string {
	return cf.GetString(KeyContext)
}

// AddFlagSideloadStartType adds a flag for the type when starting an asset.
func (cf *CmdFlags) AddFlagSideloadStartType(assetType string) {
	cf.OptionalString(KeyType, string(imageutils.Archive), fmt.Sprintf(
		`The target's type:
%-10s file path pointing to a %s bundle file`,
		imageutils.Archive,
		assetType,
	))
}

// GetFlagSideloadStartType gets the value of the type flag added by AddFlagSideloadStartType.
func (cf *CmdFlags) GetFlagSideloadStartType() string {
	return cf.GetString(KeyType)
}

// AddFlagSideloadStopType adds a flag for the type when stopping an asset.
func (cf *CmdFlags) AddFlagSideloadStopType(assetType string) {
	cf.OptionalString(KeyType, string(imageutils.ID), fmt.Sprintf(
		`The target's type:
%-10s build target that creates a %s bundle file
%-10s %s id`,
		imageutils.Build,
		assetType,
		imageutils.ID,
		assetType,
	))
}

// GetFlagSideloadStopType gets the value of the type flag added by AddFlagSideloadStopType.
func (cf *CmdFlags) GetFlagSideloadStopType() string {
	return cf.GetString(KeyType)
}

// AddFlagSideloadStartTimeout adds a flag for the timeout when starting an asset.
func (cf *CmdFlags) AddFlagSideloadStartTimeout(assetType string) {
	cf.OptionalString(KeyTimeout, "180s", fmt.Sprintf(`Maximum time to wait for the %s to
become available in the cluster after starting it. Can be set to any valid duration
(\"60s\", \"5m\", ...) or to \"0\" to disable waiting.`, assetType))
}

// GetFlagSideloadStartTimeout gets the value of the flag added by AddFlagSideloadStartTimeout.
func (cf *CmdFlags) GetFlagSideloadStartTimeout() (time.Duration, string, error) {
	timeoutStr := cf.GetString(KeyTimeout)
	timeout, err := parseNonNegativeDuration(timeoutStr)
	if err != nil {
		return timeout, timeoutStr, errors.Wrapf(err, "invalid value passed for --%s", KeyTimeout)
	}

	return timeout, timeoutStr, nil
}

// AddFlagSkipDirectUpload adds a flag for disabling direct upload to workcells
func (cf *CmdFlags) AddFlagSkipDirectUpload(assetType string) {
	usage := fmt.Sprintf("Skips direct upload of %s to workcell. Requires "+
		"external repository. (default false)\nCan be defined via the %s_%s "+
		"environment variable.", assetType, envPrefix, strings.ToUpper(KeySkipDirectUpload))
	cf.OptionalBool(KeySkipDirectUpload, false, usage)
	cf.cmd.PersistentFlags().Lookup(KeySkipDirectUpload).Hidden = true
	cf.viperLocal.BindEnv(KeySkipDirectUpload)
}

// GetFlagSkipDirectUpload gets the value of the flag added by AddFlagSkipDirectUpload
func (cf *CmdFlags) GetFlagSkipDirectUpload() bool {
	return cf.GetBool(KeySkipDirectUpload)
}

// AddFlagSkipPrompts adds a flag for disabling user prompts.
func (cf *CmdFlags) AddFlagSkipPrompts() {
	cf.OptionalBool(KeySkipPrompts, false, "True to skip user prompts.")
}

// GetFlagSkipPrompts gets the value of the flag added by AddFlagSkipPrompts
func (cf *CmdFlags) GetFlagSkipPrompts() bool {
	return cf.GetBool(KeySkipPrompts)
}

// AddFlagVendor adds a flag for the asset vendor.
func (cf *CmdFlags) AddFlagVendor(assetType string) {
	cf.RequiredString(KeyVendor, fmt.Sprintf("The %s vendor.", assetType))
}

// GetFlagVendor gets the value of the vendor flag added by AddFlagVendor.
func (cf *CmdFlags) GetFlagVendor() string {
	return cf.GetString(KeyVendor)
}

// AddFlagVersion adds a flag for the asset version.
func (cf *CmdFlags) AddFlagVersion(assetType string) {
	cf.RequiredString(KeyVersion, fmt.Sprintf("The %s version, in sem-ver format.", assetType))
}

// GetFlagVersion gets the value of the version flag added by AddFlagVersion.
func (cf *CmdFlags) GetFlagVersion() string {
	return cf.GetString(KeyVersion)
}

// AddFlagView adds a flag for the asset view.
func (cf *CmdFlags) AddFlagView() {
	shortStrings := make([]string, 0, len(viewpb.AssetViewType_value))
	for _, v := range viewpb.AssetViewType_value {
		shortStrings = append(shortStrings, viewutils.ShortStringFromEnum(viewpb.AssetViewType(v)))
	}
	slices.Sort(shortStrings)
	cf.OptionalString(KeyView, "", fmt.Sprintf("The view of the asset to return. Can be: %v", shortStrings))
}

// GetFlagView gets the value of the view flag added by AddFlagView.
func (cf *CmdFlags) GetFlagView() (viewpb.AssetViewType, error) {
	return viewutils.EnumFromShortString(cf.GetString(KeyView))
}

// String adds a new string flag.
func (cf *CmdFlags) String(name string, value string, usage string) {
	cf.cmd.PersistentFlags().String(name, value, usage)
	cf.viperLocal.BindPFlag(name, cf.cmd.PersistentFlags().Lookup(name))
}

// RequiredString adds a new required string flag.
func (cf *CmdFlags) RequiredString(name string, usage string) {
	cf.String(name, "", fmt.Sprintf("(required) %s", usage))
	cf.cmd.MarkFlagRequired(name)
}

// OptionalString adds a new optional string flag.
func (cf *CmdFlags) OptionalString(name string, value string, usage string) {
	cf.String(name, value, fmt.Sprintf("(optional) %s", usage))
}

// RequiredEnvString adds a new required string flag that is bound to the corresponding ENV
// variable.
func (cf *CmdFlags) RequiredEnvString(name string, value string, usage string) {
	envVarName := strings.ToUpper(fmt.Sprintf("%s_%s", envPrefix, name))
	cf.envString(name, value, fmt.Sprintf("%s\nRequired unless %s environment variable is defined.", usage, envVarName))

	if cf.GetString(name) == "" {
		cf.cmd.MarkPersistentFlagRequired(name)
	}
}

// OptionalEnvString adds a new optional string flag that is bound to the corresponding ENV
// variable.
func (cf *CmdFlags) OptionalEnvString(name string, value string, usage string) {
	envVarName := strings.ToUpper(fmt.Sprintf("%s_%s", envPrefix, name))
	cf.envString(name, value, fmt.Sprintf("%s\nCan be defined via the %s environment variable.", usage, envVarName))
}

// GetString gets the value of a string flag.
func (cf *CmdFlags) GetString(name string) string {
	return cf.viperLocal.GetString(name)
}

// Bool adds a new bool flag.
func (cf *CmdFlags) Bool(name string, value bool, usage string) {
	cf.cmd.PersistentFlags().Bool(name, value, usage)
	cf.viperLocal.BindPFlag(name, cf.cmd.PersistentFlags().Lookup(name))
}

// RequiredBool adds a new required bool flag.
func (cf *CmdFlags) RequiredBool(name string, usage string) {
	cf.Bool(name, false, fmt.Sprintf("(required) %s", usage))
	cf.cmd.MarkFlagRequired(name)
}

// OptionalBool adds a new optional bool flag.
func (cf *CmdFlags) OptionalBool(name string, value bool, usage string) {
	cf.Bool(name, value, fmt.Sprintf("(optional) %s", usage))
}

// GetBool gets the value of a bool flag.
func (cf *CmdFlags) GetBool(name string) bool {
	return cf.viperLocal.GetBool(name)
}

// Int adds a new int flag.
func (cf *CmdFlags) Int(name string, value int, usage string) {
	cf.cmd.PersistentFlags().Int(name, value, usage)
	cf.viperLocal.BindPFlag(name, cf.cmd.PersistentFlags().Lookup(name))
}

// RequiredInt adds a new required int flag.
func (cf *CmdFlags) RequiredInt(name string, usage string) {
	cf.Int(name, 0, fmt.Sprintf("(required) %s", usage))
	cf.cmd.MarkFlagRequired(name)
}

// OptionalInt adds a new optional int flag.
func (cf *CmdFlags) OptionalInt(name string, value int, usage string) {
	cf.Int(name, value, fmt.Sprintf("(optional) %s", usage))
}

// GetInt gets the value of an int flag.
func (cf *CmdFlags) GetInt(name string) int {
	return cf.viperLocal.GetInt(name)
}

// IsSet checks if value of given flag was set on command line.
// Allows to check if value is coming from user, or is default value.
func (cf *CmdFlags) IsSet(name string) bool {
	return cf.viperLocal.IsSet(name)
}

func (cf *CmdFlags) envString(name string, value string, usage string) {
	cf.String(name, value, usage)
	cf.viperLocal.BindEnv(name)
}

func parseNonNegativeDuration(durationStr string) (time.Duration, error) {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, fmt.Errorf("parsing duration: %w", err)
	}
	if duration < 0 {
		return 0, fmt.Errorf("duration must not be negative, but got %q", durationStr)
	}
	return duration, nil
}

func (cf *CmdFlags) createBasicAuth() *bundleimages.BasicAuth {
	user, pwd := cf.GetFlagsRegistryAuthUserPassword()
	if len(user) == 0 || len(pwd) == 0 {
		return nil
	}
	return &bundleimages.BasicAuth{
		User: user,
		Pwd:  pwd,
	}
}

func (cf *CmdFlags) authOpt() remote.Option {
	if auth := cf.createBasicAuth(); auth != nil {
		return remote.WithAuth(authn.FromConfig(authn.AuthConfig{
			Username: auth.User,
			Password: auth.Pwd,
		}))
	}
	return remote.WithAuthFromKeychain(google.Keychain)
}

// CreateRegistryOpts creates registry options for processing images.
func (cf *CmdFlags) CreateRegistryOpts(ctx context.Context) bundleimages.RegistryOptions {
	return cf.CreateRegistryOptsWithTransferer(
		ctx,
		imagetransfer.RemoteTransferer(remote.WithContext(ctx), cf.authOpt()),
		cf.GetFlagRegistry(),
	)
}

// CreateRegistryOptsWithTransferer creates registry options for processing images.
func (cf *CmdFlags) CreateRegistryOptsWithTransferer(ctx context.Context, transferer imagetransfer.Transferer, registry string) bundleimages.RegistryOptions {
	opts := bundleimages.RegistryOptions{
		Transferer: transferer,
		URI:        registry,
	}
	if auth := cf.createBasicAuth(); auth != nil {
		opts.BasicAuth = *auth
	}
	return opts
}

// MarkHidden marks all flags in the list as hidden for the given command
func (cf *CmdFlags) MarkHidden(flagsToHide ...string) {
	flags := cf.cmd.PersistentFlags()
	for _, flagName := range flagsToHide {
		flag := flags.Lookup(flagName)
		if flag != nil {
			flag.Hidden = true
		}
	}
}
