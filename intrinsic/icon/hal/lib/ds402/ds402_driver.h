// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_DS402_DS402_DRIVER_H_
#define INTRINSIC_ICON_HAL_LIB_DS402_DS402_DRIVER_H_

#include <cstdint>

#include "absl/strings/string_view.h"
#include "intrinsic/icon/utils/bitset.h"
#include "intrinsic/icon/utils/realtime_status_or.h"

/*
# The DS402 state diagram
                            ┌──────────────────┐
                            │                  │
                            │  Not ready to    │
                            │                  ◄────── Start
                            │  switch on       │
                            │                  │
                            └────────┬─────────┘
                                     │
                                 automatic
                                     │
                            ┌────────▼─────────┐           ┌──────────────────┐
         ┌──────────────────►                  ◄─────15────┤                  │
         │                  │  Switched on     │           │                  │
         │                  │                  │           │  Fault           │
         │                  │  disabled        │           │                  │
         │            ┌─────►                  ◄────────┐  │                  │
         │            │     └─┬──────────────┬─┘        │  └────────▲─────────┘
         │            │       │              │          │           │
         │            │       2              7          │           │
         │            │       │              │          │           │
         │            │     ┌─▼──────────────┴─┐        │           │
         │            │     │                  │        │           │
         │            │     │  Ready to        │        │           │
        12           10     │                  ◄───┐    │           │
         │            │     │  switch on       │   │    │           │
         │            │     │                  │   │    │           │
         │            │     └─┬──────────────▲─┘   │    │           │
         │            │       │              │     │    │           │
         │            │       3              6     │    │           │
         │            │       │              │     │    │           │
         │            │     ┌─▼──────────────┴─┐   │    │           │
         │            │     │                  │   8    9          14
         │            │     │                  │   │    │           │
         │            └─────┤  Switched on     │   │    │           │
         │                  │                  │   │    │           │
         │                  │                  │   │    │           │
         │                  └─┬──────────────▲─┘   │    │           │
         │                    │              │     │    │           │
         │                    4              5     │    │           │
         │                    │              │     │    │           │
┌────────┴─────────┐        ┌─▼──────────────┴─┐   │    │  ┌────────┴─────────┐
│                  ◄───11───┤                  ├───┘    │  │                  │
│  Quick stop      │        │  Operation       │        │  │  Fault reaction  │
│                  │        │                  │        │  │                  │
│  active          │        │  enabled         │        │  │  active          │
│                  ├───16───►                  ├────────┘  │                  │
└──────────────────┘        └──────────────────┘           └────────▲─────────┘
                                                                    │
                                                                    │
                                                                    │
                                                                 On fault
source:
https://asciiflow.com/#/share/eJzdV0tOwzAQvcrIaxZNSFvgACxBwDYbkxoR0SZV60KrgkCIJQsWUeAcLFFP05Pgps3HieMkdiohRlY1cZ03bz72xEvk4RFBJ2gwtTomTCmmBEbYuXU9gg7QEC%2FIhP27tNHcRifHVufARgummX2TaZTMKXuw0SPsxHiONTNW4DDRrETrJlov0fqJdmTbHhsgkXWwUh7V0MLJOq%2Bd%2BRQmBA8WQP0GrxUmP9%2BKtOGK4gmtAzd9cKlzC76XTu7N531loQGLSPCM%2BiNMXac9XBn%2FsMI7DuZ9HbwojY8sQ%2FVghz8C58Kf3DKjW%2B5DQuE1x0gYNYCrqALJIK7BfEReAU7xbEhVsOWTEbYy74E7xddDxrstbOVs8IlQs1zGoPbgsNZBUFqp4be0jr9q8BSkhZvIp2JVC9METvqiNU3w2udY80Cpzk9zu1CQNjAvM%2F2vAaaRTZXRkXDk94oWV0GnTJe0HFU91GgfrRo0j92W1LQsqXgt1EPgpKePuA%2Beu9YdNon7LkkfupbFB35xEgxLD0%2BNpegwKnZ9RXyJu38SM1DdnV86dtut%2BVgs%2FrGri6fCsWb7KynDBGmt0Vez7HTobL0RhCVzxTMMkQNcUPILYl7CxbHNi5nr3LFLvT8upCJSz8dkwq5Mue6XqtHP9judXWqd7UqZT2WdXoYv4b8xeU%2FK3yce97EuxheASPinOemlT8Iv9ooCkvumeiHkq1sTo8KLzRFVUdlFD5vLP0Y69%2BBms3ts9ISefgEUPgIM)
*/

namespace intrinsic::ds402 {

// Default control word commands, required to initiate state changes in the
// DS402 state machine.
// Note, that transition 14 (fault reaction active to fault) is automatic and
// cannot be commanded.
// clang-format off
enum class ControlWord : uint16_t {
  kShutdown =                      0b00000110,  // =   6, transitions: 2,6,8
  kSwitchOn =                      0b00000111,  // =   7, transitions: 3
  kDisableOperation =              0b00000111,  // =   7, transitions: 5
  kEnableOperation =               0b00001111,  // =  15, transitions: 4
  kDisableVoltage =                0b00000000,  // =   0, transitions: 7,9,10,12
  kQuickStop =                     0b00000010,  // =   2, transitions: 11
  kFaultReset =                    0b10000000,  // = 128, transitions: 15
  kEnableOperationAfterQuickStop = 0b00001111,  // =  15, transitions: 16
};
// clang-format on

// Valid DS402 States.
enum class Ds402State : uint8_t {
  kNotReadyToSwitchOn = 0,
  kSwitchOnDisabled = 1,
  kReadyToSwitchOn = 2,
  kSwitchedOn = 3,
  kOperationEnabled = 4,
  kFault = 5,
  kFaultReactionActive = 6,
  kQuickStopActive = 7,
};

// Complete list of status words indicating a particular state.
// The 4th bit and for some states even the 5th bit can freely be chosen (by the
// device), while still encoding the same state. No additional enumerators will
// be added, since the enumerators below already encode all possible variants
// for each state. Please keep in mind though, that the status word from a ds402
// device might encode more information (in addition to the state) beyond the
// 7th bit.
// clang-format off
enum class StateEncodingStatusWord : uint16_t {
  kNotReadyToSwitchOn0 =  0b0000000,
  kNotReadyToSwitchOn1 =  0b0100000,
  kNotReadyToSwitchOn2 =  0b0010000,
  kNotReadyToSwitchOn3 =  0b0110000,
  kSwitchOnDisabled0 =    0b1000000,
  kSwitchOnDisabled1 =    0b1100000,
  kSwitchOnDisabled2 =    0b1010000,
  kSwitchOnDisabled3 =    0b1110000,
  kReadyToSwitchOn0 =     0b0100001,
  kReadyToSwitchOn1 =     0b0110001,
  kSwitchedOn0 =          0b100011,
  kSwitchedOn1 =          0b110011,
  kOperationEnabled0 =    0b100111,
  kOperationEnabled1 =    0b110111,
  kQuickStopActive0 =     0b00111,
  kQuickStopActive1 =     0b10111,
  kFaultReactionActive0 = 0b001111,
  kFaultReactionActive1 = 0b101111,
  kFaultReactionActive2 = 0b011111,
  kFaultReactionActive3 = 0b111111,
  kFault0 =               0b001000,
  kFault1 =               0b101000,
  kFault2 =               0b011000,
  kFault3 =               0b111000
};
// clang-format on

// Returns the DS402 state from a given status word.
// Returns an error if the provided status word, doesn't encode a DS402 state.
intrinsic::icon::RealtimeStatusOr<Ds402State> ToDs402State(
    uint16_t status_word);

bool IsOperationEnabled(intrinsic::bitset<uint16_t> status_word);

// Defines the fault handling strategy.
enum class FaultHandling {
  kPreserve,  // Do not attempt to clear any faults if present.
  kClear      // Try to clear any fault if present.
};

// Returns the required control word to move the DS402 state machine towards
// `goal_state`. `prev_state` is expected to hold the DS402 state of the drive
// from the previous cycle (i.e. based on which `prev_control_word` was
// computed). `curr_state` is the device state derived from the current
// `status_word` provided by the drive. `prev_control_word` is expected to hold
// the previously sent control word.
//
// If `curr_state` indicates a fault state, then based on the value of
// `fault_handling` the returned control world will have the fault reset bit
// set or not. Some DS402 drives expect a rising edge of the fault reset bit.
// Thus, when the `curr_state` is kFault (and `fault_handling` is kClear), and
// `prev_state` is anything other than kFault the function will return
// `kDisableVoltage` instead of `kFaultReset`.
//
// If the device is busy (i.e. not yet ready, or still transitioning) or already
// at the `goal_state` the function will return `prev_control_word`. An error is
// returned if the `goal_state` is unknown or invalid or if there's no valid
// transition to the `goal_state`.
intrinsic::icon::RealtimeStatusOr<ControlWord> GetNextControlWord(
    Ds402State prev_state, Ds402State curr_state, Ds402State goal_state,
    ControlWord prev_control_word, FaultHandling fault_handling);

// Returns the name of `state` as a string_view.
absl::string_view ToString(Ds402State state);
absl::string_view ToString(ControlWord control_word);

}  // namespace intrinsic::ds402

#endif  // INTRINSIC_ICON_HAL_LIB_DS402_DS402_DRIVER_H_
