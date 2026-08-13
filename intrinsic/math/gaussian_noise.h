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

#ifndef INTRINSIC_MATH_GAUSSIAN_NOISE_H_
#define INTRINSIC_MATH_GAUSSIAN_NOISE_H_

#include "absl/random/bit_gen_ref.h"
#include "intrinsic/eigenmath/types.h"

namespace intrinsic {

// Uses absl::Gaussian() with the provided generator, mean and stddev to create
// random values.
class GaussianGenerator {
 public:
  GaussianGenerator() = delete;

  // Maintains a reference to `gen` until destruction.
  explicit GaussianGenerator(absl::BitGenRef gen, double mean = 0.0,
                             double stddev = 1.0);

  // Generates a random number, generated using absl::Gaussian(gen, mean,
  // stddev)
  double Generate();

  // Generates a random vector of size `n`, where each element is set to
  // absl::Gaussian(gen, mean, stddev);
  eigenmath::VectorXd Generate(int n);

 private:
  absl::BitGenRef gen_;  // owned externally.
  double mean_;
  double stddev_;
};

}  // namespace intrinsic

#endif  // INTRINSIC_MATH_GAUSSIAN_NOISE_H_
