// Copyright 2023 Intrinsic Innovation LLC

// Package deviceservice implements a DeviceService rpc service that resolves logical EtherCAT
// variable references into concrete hardware addresses and EBI instructions.
package deviceservice

import (
	"context"
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"

	"intrinsic/assets/dependencies/utils"

	log "github.com/golang/glog"

	dagrpcpb "intrinsic/assets/data/proto/v1/data_assets_go_proto"
	dscpb "intrinsic/icon/fieldbus/ethercat/device_service/v1/device_service_config_go_proto"
	dspb "intrinsic/icon/fieldbus/ethercat/device_service/v1/device_service_go_proto"
	esipb "intrinsic/icon/fieldbus/ethercat/device_service/v1/esi_go_proto"
)

var (
	// ErrConfigNil indicates that the user-provided configuration is nil.
	ErrConfigNil = errors.New("config must not be nil")
	// ErrNoDeviceIdentifier indicates that the configuration is missing the mandatory vendor/product ID.
	ErrNoDeviceIdentifier = errors.New("config must have a device identifier")
	// ErrEsiBundleNotFound indicates the requested ESI bundle data asset could not be retrieved.
	ErrEsiBundleNotFound = errors.New("ESI data asset(s) not found")
	// ErrEsiUnmarshal indicates that the ESI XML data is malformed or incompatible.
	ErrEsiUnmarshal = errors.New("unmarshalling ESI data asset")
	// ErrEsiInvalidPath indicates that an ESI bundle contains invalid paths that might violate security constraints.
	ErrEsiInvalidPath = errors.New("ESI bundle contains invalid path")
	// ErrEsiDeepNesting indicates that an ESI file contains nested InfoReferences beyond the supported depth.
	ErrEsiDeepNesting = errors.New("ESI bundle contains unsupported nested InfoReferences")

	esiBundleMsg                esipb.EsiBundle
	esiBundleDataAssetProtoName = string(esiBundleMsg.ProtoReflect().Descriptor().FullName())

	addressPattern = regexp.MustCompile(`^(?i)(?:0x|#x)?([0-9a-f]+)\.(?:0x|#x)?([0-9a-f]+)$`)
)

const (
	// variableReferenceSeparator is the string used to concatenate PDO and Object
	// references for use as unique keys in mapping tables.
	variableReferenceSeparator = "::"
)

// DeviceService implements the DeviceService rpc service.
// It maintains indices derived from the ESI to perform rapid variable resolution.
type DeviceService struct {
	// config holds the user-provided service configuration.
	config *dscpb.DeviceServiceConfig
	// esiBundle contains the ESI bundle loaded as a data asset.
	esiBundle *esipb.EsiBundle

	// objectIndex maps object addresses (Index.SubIndex) to their metadata.
	objectIndex map[objectAddress]*objectMetadata
	// objectNameIndex maps localized name strings to one or more object addresses.
	objectNameIndex map[string][]objectAddress
	// pdoIndex maps PDO indices to their metadata.
	pdoIndex map[uint32]*pdoMetadata

	// activePdos tracks PDOs used in the current resolution request.
	activePdos map[uint32]bool
	// exclusionsToAdd tracks default-active PDOs that must be disabled due to exclusions.
	exclusionsToAdd []uint32
	// exclusionsToRemove tracks default-inactive PDOs that must be enabled.
	exclusionsToRemove []uint32
	// objectsToAdd tracks variables that must be appended to a dynamic PDO mapping.
	objectsToAdd []*dspb.ResolvedConfiguration_EbiPdoInstructions_ObjectAddition
	// resolvedVars stores the final mapping of "pdo::object" to resolution metadata.
	resolvedVars map[string]*dspb.ResolvedVariable

	// supportedOpModes lists the synchronization modes supported by the device (from ESI).
	supportedOpModes []*dspb.OpModeInfo

	// resolvedConfiguration holds the pre-computed configuration response data.
	resolvedConfiguration *dspb.ResolvedConfiguration
}

// fetchESIBundle retrieves the ESI bundle data asset from the DataAsset service.
//
// Parameters:
//   - ctx: The context for the RPC call.
//   - config: The config containing the EsiBundle dependency.
//   - daClient: The client used to fetch the asset data.
//
// Returns:
//   - The populated EsiBundle proto.
//   - An error if the asset is missing, unmarshalling fails, or paths are invalid.
func fetchESIBundle(ctx context.Context, config *dscpb.DeviceServiceConfig, daClient dagrpcpb.DataAssetsClient) (*esipb.EsiBundle, error) {
	iface := "data://" + esiBundleDataAssetProtoName
	anyProto, err := utils.GetDataPayload(ctx, config.GetEsiBundle(), iface, utils.WithDataAssetsClient(daClient))
	if err != nil {
		return nil, fmt.Errorf("get ESI bundle data asset for interface %q failed: %w: %w", iface, ErrEsiBundleNotFound, err)
	}

	esiBundle := &esipb.EsiBundle{}
	if err := anyProto.UnmarshalTo(esiBundle); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrEsiUnmarshal, err)
	}

	if len(esiBundle.GetFiles()) == 0 {
		return nil, fmt.Errorf("ESI bundle has no files: %w", ErrEsiUnmarshal)
	}
	for p := range esiBundle.GetFiles() {
		if path.IsAbs(p) {
			return nil, fmt.Errorf("ESI bundle contains absolute path %q: %w", p, ErrEsiInvalidPath)
		}
		cleaned := path.Clean(p)
		if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
			return nil, fmt.Errorf("ESI bundle contains path %q resolving to outside bundle: %w", p, ErrEsiInvalidPath)
		}
	}

	return esiBundle, nil
}

// NewDeviceService constructs a new service instance and performs ESI indexing.
//
// Parameters:
//   - ctx: Request context.
//   - config: Service configuration containing device IDs and bundle references.
//   - daClient: Client for the Data Asset service.
//
// Returns:
//   - An initialized DeviceService ready to process GetConfiguration calls.
//   - An error if the config is invalid, indexing fails, or the device isn't found in the ESI.
func NewDeviceService(ctx context.Context, config *dscpb.DeviceServiceConfig, daClient dagrpcpb.DataAssetsClient) (*DeviceService, error) {
	if config == nil {
		return nil, ErrConfigNil
	}
	if config.GetDeviceIdentifier() == nil {
		return nil, ErrNoDeviceIdentifier
	}
	esiBundle, err := fetchESIBundle(ctx, config, daClient)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ESI bundle: %w", err)
	}

	s := &DeviceService{config: config, esiBundle: esiBundle}

	di := config.GetDeviceIdentifier()
	if err := s.loadAndIndexESI(ctx, esiBundle, uint32(di.VendorId), uint32(di.ProductCode), uint32(di.Revision)); err != nil {
		return nil, fmt.Errorf("failed to index ESI: %w", err)
	}

	if err := s.computeResolvedConfiguration(ctx); err != nil {
		return nil, fmt.Errorf("failed to resolve configuration: %w", err)
	}

	log.InfoContextf(ctx, "Device service started with config: %v", config)
	return s, nil
}

// GetConfiguration returns the configuration including the resolved VariableReferences, EBI instructions and metadata.
//
// Parameters:
//   - ctx: Request context.
//   - req: The GetConfigurationRequest (currently empty).
//
// Returns:
//   - A response containing the ESI bundle and the fully resolved ResolvedConfiguration.
//   - An error if resolution fails, a variable cannot be found, or PDO conflicts are detected.
func (s *DeviceService) GetConfiguration(ctx context.Context, req *dspb.GetConfigurationRequest) (*dspb.GetConfigurationResponse, error) {
	if s.esiBundle == nil {
		return nil, fmt.Errorf("no ESI bundle loaded")
	}

	return &dspb.GetConfigurationResponse{
		DeviceServiceConfig:   s.config,
		EsiBundle:             s.esiBundle,
		ResolvedConfiguration: s.resolvedConfiguration,
	}, nil
}

// computeResolvedConfiguration resolves all VariableReferences and detects ESI-related errors at startup.
func (s *DeviceService) computeResolvedConfiguration(ctx context.Context) error {
	// Initialize/Reset internal state.
	s.activePdos = make(map[uint32]bool)
	s.exclusionsToAdd = nil
	s.exclusionsToRemove = nil
	s.objectsToAdd = nil
	s.resolvedVars = make(map[string]*dspb.ResolvedVariable)

	// Resolve all mapped interfaces.
	for _, im := range s.config.GetInterfaceMappings() {
		if ds402 := im.GetDeviceData().GetDs402DeviceData(); ds402 != nil {
			if err := s.resolveDs402Device(ctx, ds402); err != nil {
				return err
			}
		}

		if joint := im.GetDeviceData().GetJointDeviceData(); joint != nil {
			if err := s.resolveJointDevice(ctx, joint); err != nil {
				return err
			}
		}

		if adio := im.GetDeviceData().GetAdioDeviceData(); adio != nil {
			if err := s.resolveAdioDevice(ctx, adio); err != nil {
				return err
			}
		}

		if ft := im.GetDeviceData().GetForceTorqueDeviceData(); ft != nil {
			if err := s.resolveForceTorqueDevice(ctx, ft); err != nil {
				return err
			}
		}

		// Handle other DeviceData types here in the future
	}

	// Validate that no mutually exclusive PDOs were activated across different variables.
	for pdoIdx := range s.activePdos {
		meta, ok := s.pdoIndex[pdoIdx]
		if !ok {
			continue
		}
		for _, excludeIdx := range meta.Excludes {
			if s.activePdos[excludeIdx] {
				otherMeta := s.pdoIndex[excludeIdx]
				return fmt.Errorf("PDO conflict detected: %q (#x%04X) and %q (#x%04X) are mutually exclusive", meta.Name, pdoIdx, otherMeta.Name, excludeIdx)
			}
		}
	}

	var instructions *dspb.ResolvedConfiguration_EbiPdoInstructions
	if len(s.exclusionsToAdd) > 0 || len(s.exclusionsToRemove) > 0 || len(s.objectsToAdd) > 0 {
		instructions = &dspb.ResolvedConfiguration_EbiPdoInstructions{
			PdoExclusionsToAdd:    s.exclusionsToAdd,
			PdoExclusionsToRemove: s.exclusionsToRemove,
			ObjectsToAdd:          s.objectsToAdd,
		}
	}

	// Determine active OpMode.
	activeOpMode, err := s.resolveActiveOpMode(ctx)
	if err != nil {
		return err
	}

	s.resolvedConfiguration = &dspb.ResolvedConfiguration{
		EbiPdoInstructions: instructions,
		VariableMappings:   s.resolvedVars,
		SupportedOpModes:   s.supportedOpModes,
		ActiveOpModeName:   activeOpMode,
	}

	return nil
}

// resolveDs402Device resolves all variables required for a DS402 drive.
func (s *DeviceService) resolveDs402Device(ctx context.Context, ds402 *dscpb.Ds402DeviceData) error {
	if err := s.resolveVariable(ctx, ds402.GetStatusWordReference(), PdoDirectionTx); err != nil {
		return err
	}
	if err := s.resolveVariable(ctx, ds402.GetErrorCodeReference(), PdoDirectionTx); err != nil {
		return err
	}
	if err := s.resolveVariable(ctx, ds402.GetControlWordReference(), PdoDirectionRx); err != nil {
		return err
	}
	if ds402.GetDigitalOutputsReference() != nil {
		if err := s.resolveVariable(ctx, ds402.GetDigitalOutputsReference(), PdoDirectionRx); err != nil {
			return err
		}
	}

	// Handle Homing configuration.
	if homing := ds402.GetHomingReference(); homing != nil {
		switch config := homing.GetConfiguration().(type) {
		case *dscpb.HomingReference_Disabled_:
			// Do nothing.
		case *dscpb.HomingReference_Auto_:
			return fmt.Errorf("homing 'auto' configuration is not implemented yet")
		case *dscpb.HomingReference_Manual_:
			m := config.Manual
			// Homing parameters are typically set via SDO/Rx during init, but some
			// can be mapped to PDOs. We treat them as Rx (Commands) or Tx (Status)
			// based on their DS402 nature.
			homingVars := []struct {
				ref *dscpb.VariableReference
				dir PdoDirection
			}{
				{m.GetHomingMethod(), PdoDirectionRx},
				{m.GetHomeOffset(), PdoDirectionRx},
				{m.GetHomingSpeed(), PdoDirectionRx},
				{m.GetHomingCreepSpeed(), PdoDirectionRx},
				{m.GetHomingAcceleration(), PdoDirectionRx},
				{m.GetModesOfOperation(), PdoDirectionRx},
				{m.GetModesOfOperationDisplay(), PdoDirectionTx},
			}
			for _, v := range homingVars {
				if v.ref != nil {
					if err := s.resolveVariable(ctx, v.ref, v.dir); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// resolveJointDevice resolves all variables required for a standard joint device.
func (s *DeviceService) resolveJointDevice(ctx context.Context, joint *dscpb.JointDeviceData) error {
	if err := s.resolveVariable(ctx, joint.GetJointPositionCommandReference(), PdoDirectionRx); err != nil {
		return err
	}
	if err := s.resolveVariable(ctx, joint.GetJointPositionStateReference(), PdoDirectionTx); err != nil {
		return err
	}
	if joint.GetJointVelocityStateReference() != nil {
		if err := s.resolveVariable(ctx, joint.GetJointVelocityStateReference(), PdoDirectionTx); err != nil {
			return err
		}
	}
	if joint.GetJointAccelerationStateReference() != nil {
		if err := s.resolveVariable(ctx, joint.GetJointAccelerationStateReference(), PdoDirectionTx); err != nil {
			return err
		}
	}
	return nil
}

// resolveAdioDevice resolves all variables required for an ADIO device.
// Silently ignores when AdioDeviceData is empty (i.e. no actual ADIO variable has been specified).
func (s *DeviceService) resolveAdioDevice(ctx context.Context, adio *dscpb.AdioDeviceData) error {
	switch v := adio.GetAdioVariable().(type) {
	case *dscpb.AdioDeviceData_DigitalInput:
		if err := s.resolveVariable(ctx, v.DigitalInput.GetVariableReference(), PdoDirectionTx); err != nil {
			return err
		}
	case *dscpb.AdioDeviceData_DigitalOutput:
		if err := s.resolveVariable(ctx, v.DigitalOutput.GetVariableReference(), PdoDirectionRx); err != nil {
			return err
		}
	case *dscpb.AdioDeviceData_AnalogInputs:
		for _, ref := range v.AnalogInputs.GetNamedVariableReferences() {
			if err := s.resolveVariable(ctx, ref, PdoDirectionTx); err != nil {
				return err
			}
		}
	case *dscpb.AdioDeviceData_AnalogOutputs:
		for _, ref := range v.AnalogOutputs.GetNamedVariableReferences() {
			if err := s.resolveVariable(ctx, ref, PdoDirectionRx); err != nil {
				return err
			}
		}
	}
	return nil
}

// resolveForceTorqueDevice resolves all variables required for a force torque sensor.
func (s *DeviceService) resolveForceTorqueDevice(ctx context.Context, ft *dscpb.ForceTorqueDeviceData) error {
	ftVars := []struct {
		ref *dscpb.VariableReference
		dir PdoDirection
	}{
		{ft.GetForceX(), PdoDirectionTx},
		{ft.GetForceY(), PdoDirectionTx},
		{ft.GetForceZ(), PdoDirectionTx},
		{ft.GetTorqueX(), PdoDirectionTx},
		{ft.GetTorqueY(), PdoDirectionTx},
		{ft.GetTorqueZ(), PdoDirectionTx},
		{ft.GetStatusCode(), PdoDirectionTx},
		{ft.GetSampleCounter(), PdoDirectionTx},
		{ft.GetControlCode(), PdoDirectionRx},
	}
	for _, v := range ftVars {
		if v.ref != nil {
			if err := s.resolveVariable(ctx, v.ref, v.dir); err != nil {
				return err
			}
		}
	}
	return nil
}

// resolveVariable translates a logical VariableReference into a hardware address and EBI instructions.
func (s *DeviceService) resolveVariable(ctx context.Context, ref *dscpb.VariableReference, preferredDir PdoDirection) error {
	if ref == nil {
		return nil
	}
	tracingEnabled := s.config.GetOptions().GetEnableVariableResolutionTracing()
	if tracingEnabled {
		log.InfoContextf(ctx, "--- Tracing resolution for variable: pdo=%q, object=%q (PreferredDir=%v) ---", ref.GetPdo(), ref.GetObject(), preferredDir)
	}

	key := ref.GetPdo() + variableReferenceSeparator + ref.GetObject()
	if _, ok := s.resolvedVars[key]; ok {
		if tracingEnabled {
			log.InfoContextf(ctx, "Variable already resolved (cached).")
		}
		return nil
	}

	// 1. Resolve PDO Index (Implicit or Explicit).
	// We do this first because the Object resolution now depends on the PDO context.
	pdoIdx, pdoMeta, err := s.findPdo(ctx, ref.GetPdo(), objectAddress{}, nil, preferredDir)
	if err != nil && ref.GetPdo() != "" {
		// If explicit PDO resolution failed, we can't proceed.
		return fmt.Errorf("explicit PDO %q requested for variable %q was not found in ESI or is ambiguous: %w", ref.GetPdo(), ref.GetObject(), err)
	}

	if tracingEnabled {
		if pdoMeta != nil {
			log.InfoContextf(ctx, "Resolved PDO context: %q (#x%04X, Dir=%v)", pdoMeta.Name, pdoIdx, pdoMeta.Direction)
		} else {
			log.InfoContextf(ctx, "No explicit PDO context provided (implicit resolution).")
		}
	}

	// 2. Resolve Object Address (Hardware index/subindex).
	addr, meta, err := s.findObject(ctx, ref.GetObject(), pdoMeta, preferredDir)
	if err != nil {
		return fmt.Errorf("variable %q could not be resolved. Please check if the name or address matches the ESI file: %w", ref.GetObject(), err)
	}

	if tracingEnabled {
		log.InfoContextf(ctx, "Resolved Object: %q (#x%04X.%d, Mapping=%q)", meta.Name, addr.Index, addr.SubIndex, meta.PdoMapping)
	}

	// 3. If PDO was implicit (empty), we might need to re-resolve it now that we have the object's mapping flags.
	if ref.GetPdo() == "" {
		if tracingEnabled {
			log.InfoContextf(ctx, "Re-resolving PDO based on object mapping flags and preferred direction...")
		}
		pdoIdx, pdoMeta, err = s.findPdo(ctx, "", addr, meta, preferredDir)
		if err != nil {
			return fmt.Errorf("could not find suitable PDO for variable %q: %w", ref.GetObject(), err)
		}
		if tracingEnabled {
			log.InfoContextf(ctx, "Implicit PDO choice: %q (#x%04X)", pdoMeta.Name, pdoIdx)
		}
	}

	// 4. EBI Instruction Generation.
	if err := s.generateInstructions(ctx, pdoIdx, pdoMeta, addr, meta); err != nil {
		return fmt.Errorf("failed to map variable %q to PDO %q. This usually happens if the PDO is 'Fixed' or has a direction mismatch: %w", ref.GetObject(), pdoMeta.Name, err)
	}

	if tracingEnabled {
		log.InfoContextf(ctx, "Successfully resolved to PDO %q and Object %q.", pdoMeta.Name, meta.Name)
	}

	// Record resolution metadata.
	s.resolvedVars[key] = &dspb.ResolvedVariable{
		PdoIndex:      pdoIdx,
		Index:         addr.Index,
		SubIndex:      addr.SubIndex,
		PdoEniName:    pdoMeta.Name,
		ObjectEniName: meta.Name,
	}

	return nil
}

// resolveActiveOpMode determines the synchronization mode to use based on the user configuration
// and the device's supported modes (from ESI).
func (s *DeviceService) resolveActiveOpMode(ctx context.Context) (string, error) {
	selectedOpMode := s.config.GetSyncConfig().GetSelectedOpModeName()

	// 1. Explicit selection by user.
	if selectedOpMode != "" {
		// Validate that the selected mode is actually supported by the device.
		for _, m := range s.supportedOpModes {
			if m.Name == selectedOpMode {
				return selectedOpMode, nil
			}
		}
		// Not found: this is a configuration error.
		var available []string
		for _, m := range s.supportedOpModes {
			available = append(available, fmt.Sprintf("%q", m.Name))
		}
		return "", fmt.Errorf("selected sync mode '%q' is not supported by this device (supported: %s)", selectedOpMode, strings.Join(available, ", "))
	}

	// 2. Automatic selection (user left it empty).

	// If the device reports NO OpModes in ESI, we cannot select one.
	// This implies the device uses its default implicit synchronization (FreeRun/SM).
	if len(s.supportedOpModes) == 0 {
		return "", nil
	}

	// Check if we have Ds402 mappings, which usually benefit from DC.
	hasDs402 := false
	for _, m := range s.config.GetInterfaceMappings() {
		if m.GetDeviceData().GetDs402DeviceData() != nil {
			hasDs402 = true
			break
		}
	}

	// If Ds402 is present, prefer the first available DC mode.
	if hasDs402 {
		for _, m := range s.supportedOpModes {
			if m.IsDc {
				log.InfoContextf(ctx, "Auto-selected DC mode '%q' for DS402 device.", m.Name)
				return m.Name, nil
			}
		}
	}

	// Fallback: Default to the first available mode.
	// ESI spec typically puts the default/preferred mode first.
	return s.supportedOpModes[0].Name, nil
}
