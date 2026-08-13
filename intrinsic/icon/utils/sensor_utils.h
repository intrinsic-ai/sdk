// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#ifndef INTRINSIC_ICON_UTILS_SENSOR_UTILS_H_
#define INTRINSIC_ICON_UTILS_SENSOR_UTILS_H_

#include <array>
#include <optional>

#include "intrinsic/eigenmath/types.h"

namespace intrinsic::icon {

// A class for computing a scalar-valued sensor bias using recursive averaging
// over a given number of target samples.
class DofSensorBias {
 public:
  DofSensorBias() noexcept : sample_count_(0), bias_(0.0) {}

  void ResetBias() { bias_ = 0.0; }

  void ResetSampleCount() { sample_count_ = 0; }

  // Add sample to the recursive average filter. Updates both sample count and
  // bias.
  void AddSample(double sample) {
    if (sample_count_ == 0) bias_ = 0.0;
    sample_count_++;
    bias_ += (sample - bias_) / sample_count_;
  }

  double Bias() const { return bias_; }

  int SampleCount() const { return sample_count_; }

 private:
  int sample_count_ = 0;
  double bias_ = 0.0;
};

// Contains the data to perform taring. It is only used if taring is handled in
// the part, which in practice is only if this part is NOT configured with a
// ForceTorqueCommand interface.
//
// On initialization, biases are set to zero with no taring in progress. Call
// `StartTare` to begin taring.
//
// The class is entirely realtime safe.
class TaringData {
 public:
  // Starts the process of taring. Tick returns zero-values while taring.
  void StartTare(int taring_cycles);

  // Updates the internal data with a new reading. This should be called every
  // cycle.
  void Update(const ::intrinsic::eigenmath::Vector6d& reading);

  // Returns the tared values, or nullopt if a tare is in progress.
  std::optional<::intrinsic::eigenmath::Vector6d> GetTaredValue();

  // Returns true if taring.
  bool TaringInProgress();

 private:
  // Sensor bias which gets set during taring operation and is subsequently
  // subtracted from wrench measurements to correct for unmodelled payload,
  // etc.
  std::array<DofSensorBias, 6> ft_sensor_bias_;
  bool taring_in_progress_ = false;
  int taring_cycles_ = 0;
  std::optional<::intrinsic::eigenmath::Vector6d> latest_reading_;
};

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_UTILS_SENSOR_UTILS_H_
