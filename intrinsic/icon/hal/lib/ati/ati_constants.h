// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_ATI_ATI_CONSTANTS_H_
#define INTRINSIC_ICON_HAL_LIB_ATI_ATI_CONSTANTS_H_

#include <cstdint>

namespace intrinsic::ati {

// EtherCAT object definitions and bitfield values for ATI EtherCAT OEM
// board-powered Force/Torque sensors.
struct AtiEcatOemObjects {
  // Bits for the status field.
  enum StatusCodeBits : uint32_t {
    kMonitorConditionTripped = (1u << 0),
    kSupplyOutOfRange = (1u << 1),
    kVoltsOutOfRange = (1u << 3),
    kCurrentOutOfRange = (1u << 4),
    kDpotFault = (1u << 5),
    kEepromFault = (1u << 6),
    kDacFault = (1u << 7),
    kSimulatedError = (1u << 28),
    kCalibrationChecksumError = (1u << 29),
    kSaturation = (1u << 30),
    kError = (1u << 31),
    kOk = 0,
  };

  // Bits for the Control1 field.
  enum class Control1Bits {
    kSetBias = (1 << 0),
    kSelectGageOutput = (1 << 1),
    kSelectForceOutput = (0 << 1),
    kSetError = (1 << 2),
    kClearMonitor = (1 << 3),

    kNoFilter = (0 << 4),
    kFilter360Hz = (1 << 4),
    kFilter140Hz = (2 << 4),
    kFilter64Hz = (3 << 4),
    kFilter32Hz = (4 << 4),
    kFilter16Hz = (5 << 4),
    kFilter8Hz = (6 << 4),
    kFilter4Hz = (7 << 4),
    kFilter2Hz = (8 << 4),
    kFilterSelectionScaler = 4,

    kCalibrationSelectionScaler = 8,
    kSetBiasAgainstCurrentLoadScaler = 3,

    kInternalSampleRateSelectionScaler = 12,
  };

  // Bits for the Control2 field.
  enum Control2Bits {
    kMonitorScaler = 0,
    kToolTransformIndexScaler = 16,
    kSimulatedErrorControl = (1 << 31),
  };

  // Indexes for COE objects.
  enum {
    kToolTransformationAddress = 0x2020,
    kCalibrationAddress = 0x2040,
    kMonitorConditionAddress = 0x2060,
    kDiagnosticsAddress = 0x2080,
    kVersionAddress = 0x2090,
    kDataAddress = 0x6000,
    kStatusCodeAddress = 0x6010,
    kSampleCounterAddress = 0x6020,
    kControlCodesAddress = 0x7010,
  };

  // Subindexes for COE objects.
  enum {
    kCalibrationSerialSubIndex = 0x1,
    kCalibrationPartSubIndex = 0x2,
    kCalibrationFamilySubIndex = 0x3,
    kCalibrationDateSubIndex = 0x4,
    kCalibrationMatrixStartSubIndex = 0x5,
    kCalibrationForceUnitsSubIndex = 0x29,
    kCalibrationTorqueUnitsSubIndex = 0x2a,
    kCalibrationCountsPerForceSubIndex = 0x31,
    kCalibrationCountsPerTorqueSubIndex = 0x32,
    kCalibrationMaxForceTorqueStartSubIndex = 0x2b,
    kCalibrationGainsStartSubIndex = 0x33,
    kCalibrationOffsetsStartSubIndex = 0x39,
    kControlCodesControl1SubIndex = 0x1,
  };
};  // struct AtiEcatOemObjects

// EtherCAT object definitions and bitfield values for ATI Axia Force/Torque
// sensors.
struct AtiAxiaObjects {
  // Bits for the status field.
  enum StatusCodeBits : uint32_t {
    kInternalTemperatureOutOfRange = (1u << 0),
    kSupplyVoltageOutOfRange = (1u << 1),
    kBrokenGage = (1u << 2),
    kBusy = (1u << 3),
    kHardwareError = (1u << 5),
    kGageOutOfRange = (1u << 27),
    kSimulatedError = (1u << 28),
    kCalibrationChecksumError = (1u << 29),
    kSaturation = (1u << 30),
    kError = (1u << 31),
    kOk = 0,
  };

  // Bits for the Control1 field.
  enum class Control1Bits {
    kSetBias = (1 << 0),
    kSelectGageOutput = (1 << 1),
    kSelectForceOutput = (0 << 1),
    kSetError = (1 << 2),
    kClearMonitor = (1 << 3),

    kNoFilter = (0 << 4),
    kFilter360Hz = (1 << 4),
    kFilter140Hz = (2 << 4),
    kFilter64Hz = (3 << 4),
    kFilter32Hz = (4 << 4),
    kFilter16Hz = (5 << 4),
    kFilter8Hz = (6 << 4),
    kFilter4Hz = (7 << 4),
    kFilter2Hz = (8 << 4),
    kFilterSelectionScaler = 4,

    kCalibrationSelectionScaler = 8,
    kSetBiasAgainstCurrentLoadScaler = 3,

    kInternalSampleRateSelectionScaler = 12,
  };

  // Bits for the Control2 field.
  enum class Control2Bits {
    kMonitorScaler = 0,
    kToolTransformIndexScaler = 16,
    kSimulatedErrorControl = (1 << 31),
  };

  // Indexes for COE objects.
  enum {
    kToolTransformationAddress = 0x2020,
    kCalibrationAddress = 0x2021,
    kDiagnosticsAddress = 0x2080,
    kVersionAddress = 0x2090,
    kDataAddress = 0x6000,
    kStatusCodeAddress = 0x6010,
    kSampleCounterAddress = 0x6020,
    kControlCodesAddress = 0x7010,
  };

  // Subindexes for COE objects.
  enum {
    kCalibrationSerialSubIndex = 0x1,
    kCalibrationPartSubIndex = 0x2,
    kCalibrationFamilySubIndex = 0x3,
    kCalibrationDateSubIndex = 0x4,
    kCalibrationForceUnitsSubIndex = 0x2f,
    kCalibrationTorqueUnitsSubIndex = 0x30,
    kCalibrationMaxForceTorqueStartSubIndex = 0x31,

    kCalibrationCountsPerForceSubIndex = 0x37,
    kCalibrationCountsPerTorqueSubIndex = 0x38,

    kControlCodesControl1SubIndex = 0x1,
  };
};  // struct AtiAxiaObjects

}  // namespace intrinsic::ati

#endif  // INTRINSIC_ICON_HAL_LIB_ATI_ATI_CONSTANTS_H_
