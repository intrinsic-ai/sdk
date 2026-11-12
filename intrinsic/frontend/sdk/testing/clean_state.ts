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

/**
 * Signature of a `beforeEach`-style hook, as provided by test frameworks such
 * as Jasmine, Jest, Vitest, and Mocha.
 */
export type BeforeEachFn = (
  action: () => void | Promise<void>,
  timeout?: number,
) => void;

// Provided as a global by the test framework.
declare const beforeEach: BeforeEachFn;

/**
 * Creates a state object that is reset before each test.
 *
 * Runs `factory` before each test and replaces the contents of the returned
 * object with the result, so tests can access fresh fixtures as `state.foo`.
 *
 * @param factory Creates the state for a single test.
 * @param beforeEachFn The test framework's `beforeEach` hook. Defaults to the
 *     global `beforeEach`.
 * @returns The state object, which is updated in place before each test.
 */
export function cleanState<T extends object>(
  factory: () => T | Promise<T>,
  beforeEachFn: BeforeEachFn = beforeEach,
): T {
  const state: T = {} as T;
  beforeEachFn(async () => {
    // Modify the state object in place, so that references to it held by the
    // tests stay valid.
    for (const prop of Object.getOwnPropertyNames(state)) {
      Reflect.deleteProperty(state, prop);
    }
    Object.assign(state, await factory());
  });
  return state;
}
