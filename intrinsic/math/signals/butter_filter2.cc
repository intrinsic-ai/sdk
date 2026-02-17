// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/math/signals/butter_filter2.h"

#include <array>
#include <cmath>
#include <limits>

#include "absl/log/check.h"
#include "intrinsic/icon/utils/log.h"

namespace intrinsic {
namespace {

// Compute digital filter coefficients from analog coefficients using bilinear
// transform, sb[0]+sb[1]*s + sb[2]*s^2          zb[0]+zb[1]*z^-1 + zb[2]*z^-2
// ------------- -----------    ==>   -----------------------------
// sa[0]+sa[1]*s + sa[2]*s^2          1 + za[1]*z^-1 +za[2]*z^-3
void Bilinear(const std::array<double, 3>& sb, const std::array<double, 3>& sa,
              std::array<double, 3>* zb, std::array<double, 3>* za) {
  CHECK(nullptr != zb);
  CHECK(nullptr != za);
  CHECK_GT(std::fabs(sa[0] + sa[1] + sa[2]),
           std::numeric_limits<double>::epsilon());
  const double inv_denom = 1.0 / (sa[0] + sa[1] + sa[2]);

  (*zb)[0] = (sb[0] + sb[1] + sb[2]) * inv_denom;
  (*zb)[1] = 2.0 * (sb[0] - sb[2]) * inv_denom;
  (*zb)[2] = (sb[0] - sb[1] + sb[2]) * inv_denom;

  (*za)[0] = 1.0;
  (*za)[1] = 2.0 * (sa[0] - sa[2]) * inv_denom;
  (*za)[2] = (sa[0] - sa[1] + sa[2]) * inv_denom;
}

// analog (cutoff) frequency from digital (cutoff) frequency,
// without 2*fs factor (cancels in transfer function)
// [from s = exp(i*omg_analog)  and s= 2*fs*(z-1)/(z+1)]
// fs = sampling frequency
double prewarp(const double omegaD, const double fs) {
  // actually 2*fs*tan(omegaD*0.5/fs), but 2*fs cancels,
  // so omit to avoid large arithmetic with large numbers ...
  CHECK_GT(std::fabs(fs), std::numeric_limits<double>::epsilon());
  return std::tan(omegaD * 0.5 / fs);
}

// analog prototype butterworth lowpass filter coefficients
constexpr std::array<double, 3> kButterProtoB = {{1.0, 0.0, 0.0}};
constexpr std::array<double, 3> kButterProtoA = {{1.0, M_SQRT2, 1.0}};

// analog lowpass prototype to highpass transformation
void ProtoToLP(const double omega, const std::array<double, 3>& sb_in,
               const std::array<double, 3>& sa_in,
               std::array<double, 3>* sb_out, std::array<double, 3>* sa_out) {
  CHECK(nullptr != sb_out);
  CHECK(nullptr != sb_out);

  const double om2 = omega * omega;
  (*sb_out)[0] = sb_in[0] * om2;
  (*sb_out)[1] = sb_in[1] * omega;
  (*sb_out)[2] = sb_in[2];

  (*sa_out)[0] = sa_in[0] * om2;
  (*sa_out)[1] = sa_in[1] * omega;
  (*sa_out)[2] = sa_in[2];
}

// analog lowpass to highpass transformation
void ProtoToHP(const double omega, const std::array<double, 3>& sb_in,
               const std::array<double, 3>& sa_in,
               std::array<double, 3>* sb_out, std::array<double, 3>* sa_out) {
  CHECK(nullptr != sb_out);
  CHECK(nullptr != sb_out);

  const double om2 = omega * omega;
  (*sb_out)[0] = sb_in[2] * om2;
  (*sb_out)[1] = sb_in[1] * omega;
  (*sb_out)[2] = sb_in[0];

  (*sa_out)[0] = sa_in[2] * om2;
  (*sa_out)[1] = sa_in[1] * omega;
  (*sa_out)[2] = sa_in[0];
}
}  // namespace

bool ButterFilter2Coeffs(const double sampling_frequency,
                         const double cutoff_frequency, const FilterType type,
                         std::array<double, 3>* b, std::array<double, 3>* a) {
  CHECK(nullptr != a);
  CHECK(nullptr != b);

  // logic check
  if (2 * cutoff_frequency >= sampling_frequency) {
    INTRINSIC_RT_LOG(ERROR)
        << "cutoff frequency must be smaller than .5* sampling frequency!\n"
           "input: sampling_frequency= "
        << sampling_frequency << ", cutoff_frequency= " << cutoff_frequency;
    return false;
  }

  if ((type != FilterType::LOW_PASS) && (type != FilterType::HIGH_PASS)) {
    INTRINSIC_RT_LOG(ERROR) << "Invalid filter type " << static_cast<int>(type);
    return false;
  }

  if (sampling_frequency <= 0) {
    INTRINSIC_RT_LOG(ERROR)
        << "sampling_frequency must be positive, got " << sampling_frequency;
    return false;
  }

  if (cutoff_frequency <= 0) {
    INTRINSIC_RT_LOG(ERROR)
        << "cutoff_frequency must be positive, got " << cutoff_frequency;
    return false;
  }

  // prewarp analog frequency
  const double omega_warp =
      prewarp(cutoff_frequency * 2.0 * M_PI, sampling_frequency);

  // get analog transfer function coefficients from prototype low-pass filter
  std::array<double, 3> sa = {};
  std::array<double, 3> sb = {};
  switch (type) {
    case FilterType::LOW_PASS:
      ProtoToLP(omega_warp, kButterProtoB, kButterProtoA, &sb, &sa);
      break;
    case FilterType::HIGH_PASS:
      ProtoToHP(omega_warp, kButterProtoB, kButterProtoA, &sb, &sa);
      break;
  }

  // digital coefficients from bilinear transform
  Bilinear(sb, sa, b, a);

  return true;
}

}  // namespace intrinsic
