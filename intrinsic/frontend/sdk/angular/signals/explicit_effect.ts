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

import {
  CreateEffectOptions,
  EffectCleanupRegisterFn,
  EffectRef,
  Signal,
  effect,
  untracked,
} from '@angular/core';

// Convert a list of signals to a tuple that supports inference of generic types
// of the signal values.
type Signals<T> = [
  ...Extract<{[K in keyof T]: Signal<T[K]>}, readonly unknown[]>,
];

/**
 * Variant of `effect` that only tracks the explicitly provided dependencies.
 * This is in contrast to the default behavior of `effect` that implicitly
 * tracks all signals that were encountered during the previous execution of the
 * effect, including those in deeply nested function calls. Explicitly listing
 * all tracked signals can prevent subtle bugs due to unintended signals causing
 * the effect to re-run.
 * @param dependencies Signals that will be tracked by the effect. The effect
 *   will be re-run when any of the dependency signals changes.
 * @param fn Function that runs whenever any of the dependencies changes.
 *   Receives the latest values of the dependencies. Signals read inside this
 *   function will not be tracked.
 * @param options Options for the effect. Same as the options available for the
 *   regular `effect` function.
 */
export function explicitEffect<Values>(
  dependencies: Signals<Values>,
  fn: (values: Values, onCleanup: EffectCleanupRegisterFn) => void,
  options?: CreateEffectOptions | undefined,
): EffectRef {
  return effect((onCleanup) => {
    const values = dependencies.map((dep) => dep());
    untracked(() => {
      // We mapped the dependencies to their values above, so this type cast is correct and safe.
      fn(values as Values, onCleanup);
    });
  }, options);
}
