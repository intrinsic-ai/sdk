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

#include "intrinsic/math/gaussian_noise.h"

#include <cstddef>

#include "absl/random/bit_gen_ref.h"
#include "absl/random/distributions.h"
#include "intrinsic/eigenmath/types.h"

namespace intrinsic {

GaussianGenerator::GaussianGenerator(absl::BitGenRef gen, double mean,
                                     double stddev)
    : gen_(gen), mean_(mean), stddev_(stddev) {}

double GaussianGenerator::Generate() {
  return absl::Gaussian(gen_, mean_, stddev_);
}

eigenmath::VectorXd GaussianGenerator::Generate(int n) {
  eigenmath::VectorXd r(n);
  for (size_t i = 0; i < n; i++) r(i) = absl::Gaussian(gen_, mean_, stddev_);
  return r;
}

}  // namespace intrinsic
