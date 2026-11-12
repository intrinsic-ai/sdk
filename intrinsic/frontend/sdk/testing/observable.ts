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

import {Observable, TimeoutError, firstValueFrom, timer} from 'rxjs';
import {take, takeUntil, timeout, toArray} from 'rxjs/operators';

const DEFAULT_WAIT_MS = 1000;

/** Options for collecting values from an observable. */
export interface CollectOptions {
  waitMs?: number;
  maxValues?: number;
}

/**
 * Error thrown by {@link collectValue} if no value is emitted during the
 * configured wait time.
 */
export class NoValueError extends Error {}

/**
 * Collect the first value from an observable. Defaults to waiting for 1s, until
 * the observable emits a value or until the observable completes. Throws a
 * special {@link NoValueError} if the observable does not emit a value within
 * the wait time. Throws an error if the observable errors.
 */
export async function collectValue<T>(
  observable: Observable<T>,
  waitMs = DEFAULT_WAIT_MS,
): Promise<T> {
  const values = await collectValues(observable, {waitMs, maxValues: 1});
  if (values.length === 0) {
    throw new NoValueError(
      `Expected value to be emitted within ${waitMs}ms, but no value was emitted.`,
    );
  }
  return values[0];
}

/**
 * Collect values from an observable until the expected number of values or the
 * wait time is reached. Defaults to waiting for 1s and collecting values until
 * the observable completes. Throws an error if the observable errors.
 */
export function collectValues<T>(
  observable: Observable<T>,
  options: CollectOptions = {},
): Promise<T[]> {
  const waitMs = options.waitMs ?? DEFAULT_WAIT_MS;
  const maxValues = options.maxValues ?? Infinity;
  const timeout = timer(waitMs);
  return firstValueFrom(
    observable.pipe(take(maxValues), takeUntil(timeout), toArray()),
  );
}

/**
 * Returns whether the observable emits at least one value within the specified
 * wait time. Rejects with an error if the observable errors.
 */
export async function emits<T>(
  observable: Observable<T>,
  waitMs = DEFAULT_WAIT_MS,
): Promise<boolean> {
  try {
    await collectValue(observable, waitMs);
    return true;
  } catch (e: unknown) {
    if (e instanceof NoValueError) {
      return false;
    }
    throw e;
  }
}

/**
 * Returns whether the observable completes within the specified wait time.
 * Rejects with an error if the observable errors.
 */
export function completes<T>(
  observable: Observable<T>,
  waitMs = DEFAULT_WAIT_MS,
): Promise<boolean> {
  return new Promise<boolean>((resolve, reject) => {
    observable.pipe(timeout(waitMs)).subscribe({
      complete: () => {
        resolve(true);
      },
      error: (error) => {
        if (error instanceof TimeoutError) {
          resolve(false);
        } else {
          reject(error);
        }
      },
    });
  });
}
