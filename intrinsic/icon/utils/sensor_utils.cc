// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/utils/sensor_utils.h"

#include <optional>

#include "intrinsic/eigenmath/types.h"

namespace intrinsic::icon {

void TaringData::StartTare(int taring_cycles) {
  taring_in_progress_ = true;
  this->taring_cycles_ = taring_cycles;
  for (auto& bias : ft_sensor_bias_) {
    bias.ResetBias();
    bias.ResetSampleCount();
  }
}

void TaringData::Update(const ::intrinsic::eigenmath::Vector6d& reading) {
  latest_reading_ = reading;
  if (taring_in_progress_) {
    // Add samples first.
    for (int i = 0; i < reading.size(); i++) {
      ft_sensor_bias_.at(i).AddSample(reading(i));
    }
    if (ft_sensor_bias_.at(0).SampleCount() >= taring_cycles_) {
      taring_in_progress_ = false;
    }
  }
}

std::optional<::intrinsic::eigenmath::Vector6d> TaringData::GetTaredValue() {
  if (!latest_reading_) {
    return std::nullopt;
  }
  // We are not taring. Apply the bias and return.
  ::intrinsic::eigenmath::Vector6d out;
  for (int i = 0; i < out.size(); i++) {
    out(i) = (*latest_reading_)(i)-ft_sensor_bias_.at(i).Bias();
  }
  return out;
}

bool TaringData::TaringInProgress() { return taring_in_progress_; }

}  // namespace intrinsic::icon
