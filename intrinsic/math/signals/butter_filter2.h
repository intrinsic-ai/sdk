// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_MATH_SIGNALS_BUTTER_FILTER2_H_
#define INTRINSIC_MATH_SIGNALS_BUTTER_FILTER2_H_

#include <array>

#include "intrinsic/icon/utils/log.h"

namespace intrinsic {

enum class FilterType { LOW_PASS, HIGH_PASS };

// Calculate coefficients for digital 2nd order Butterworth filter.
// Transfer function is
// \f$H(z)=\frac{b[0]+b[1]z^{-1}+b[2]z^{-2}}{a[0]+a[1]z^{-1}+a[2]z^{-2}}\f$
bool ButterFilter2Coeffs(double sampling_frequency, double cutoff_frequency,
                         FilterType type, std::array<double, 3>* b,
                         std::array<double, 3>* a);

// 2nd Order Butterworth filter.
template <typename T>
class ButterFilter2 {
 public:
  // Initialize the filter, calculate filter coefficients
  // initial_value: initial steady state value for the filter
  // sampling_frequency: sampling frequency (in Hz)
  // cutoff_frequency: cutoff frequency (in Hz)
  // type: filter type (one of FilterType)
  // returns false on error, true on success
  bool Init(T initial_value, double sampling_frequency, double cutoff_frequency,
            FilterType type = FilterType::LOW_PASS);

  // Reset filter to constant values
  // input the steady-state value to set the filter to
  void Reset(const T& input);

  // Returns current filter output
  const T& GetOutput() const { return output_[0]; }

  // Returns filter output for derivative
  const T& GetDotOutput() const { return dot_output_; }

  // Returns filter output for second derivative
  const T& GetDDotOutput() const { return ddot_output_; }

  // Update the filter (process one timestep)
  // input is the filter input
  void Update(const T& input);

  // Returns filter coefficients "a"
  const std::array<double, 3>& GetA() const { return a_; }

  // Returns filter coefficients "b"
  const std::array<double, 3>& GetB() const { return b_; }

 private:
  // sampling frequency
  double sampling_frequency_;

  // cutoff frequency (magnitude response is sqrt(2))
  double cutoff_frequency_;

  // filter coefficients.
  std::array<double, 3> a_ = {};
  std::array<double, 3> b_ = {};

  // Input values: input_[0] is current value at time t, input_[1] at t-1,
  // input_[2] at t-2
  std::array<T, 3> input_ = {};

  // Output values: output_[0] is current value at time t, output_[1] at t-1,
  // output_[2] at t-2
  std::array<T, 3> output_;
  T dot_output_;
  T ddot_output_;
};

template <typename T>
bool ButterFilter2<T>::Init(const T initial_value,
                            const double sampling_frequency,
                            const double cutoff_frequency,
                            const FilterType type) {
  if (!ButterFilter2Coeffs(sampling_frequency, cutoff_frequency, type, &b_,
                           &a_)) {
    INTRINSIC_RT_LOG(ERROR) << "Error computing filter coefficients.";
    return false;
  }

  sampling_frequency_ = sampling_frequency;
  cutoff_frequency_ = cutoff_frequency;

  Reset(initial_value);

  return true;
}

template <typename T>
void ButterFilter2<T>::Reset(const T& input) {
  input_.fill(input);
  output_.fill(input);
  dot_output_ = input_[1] - input_[0];
  ddot_output_ = input_[1] - input_[0];
}

template <typename T>
void ButterFilter2<T>::Update(const T& input) {
  // Shift data and copy new input value
  input_[2] = input_[1];
  input_[1] = input_[0];
  input_[0] = input;
  // Shift data and calculate new output value
  output_[2] = output_[1];
  output_[1] = output_[0];
  output_[0] = (b_[0] * input_[0] + b_[1] * input_[1] + b_[2] * input_[2] -
                a_[1] * output_[1] - a_[2] * output_[2]);

  dot_output_ = sampling_frequency_ * (output_[0] - output_[1]);
  ddot_output_ = sampling_frequency_ * sampling_frequency_ *
                 (output_[0] - 2 * output_[1] + output_[2]);
}

// Definition for double arguments
using ButterFilter2d = ButterFilter2<double>;

}  // namespace intrinsic
#endif  // INTRINSIC_MATH_SIGNALS_BUTTER_FILTER2_H_
