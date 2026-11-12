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

/** Test spec function for a parameterized test. */
export type ItSpecFn<T> = (param: T) => void | Promise<void>;

/** Function encapsulating a single test. */
export type ItFn = (
  title: string,
  specFn: () => void | Promise<void>,
  timeout?: number,
) => void;

// Provided as a global by the test framework.
declare const it: ItFn;

/**
 * Parameters for parameterized tests. Array values will be inlined into the
 * test title. When using an object, the keys will be used in test titles.
 */
export type Params<T> = readonly T[] | {[key: string]: T};

/** Test cases for parameterized tests. */
export interface TestCases<T> {
  params: Params<T>;
  prefix?: string;
  suffix?: string;
}

/**
 * Runs a test for each of the given cases.
 *
 * Prefer a plain loop over the cases in new tests.
 *
 * @param cases The test cases to run.
 * @param itSpecFn The test body, called with each case's parameter.
 * @param itFn The test framework's `it` function. Defaults to the global `it`.
 */
export function pitMigration<T>(
  cases: TestCases<T>,
  itSpecFn: ItSpecFn<T>,
  itFn: ItFn = it,
): void {
  for (const [name, param] of Object.entries(normalizeParams(cases.params))) {
    const description = buildDescription(name, cases.prefix, cases.suffix);
    itFn(description, async () => {
      await itSpecFn(param);
    });
  }
}

function normalizeParams<T>(params: Params<T>): {[key: string]: T} {
  if (params instanceof Array) {
    const result: {[key: string]: T} = {};
    for (const param of params) {
      result[String(param)] = param;
    }
    return result;
  }
  return params;
}

function buildDescription(
  description: string,
  prefix?: string,
  suffix?: string,
): string {
  const parts = [description];
  if (prefix) {
    parts.unshift(prefix);
  }
  if (suffix) {
    parts.push(suffix);
  }
  return parts.join(' ');
}
