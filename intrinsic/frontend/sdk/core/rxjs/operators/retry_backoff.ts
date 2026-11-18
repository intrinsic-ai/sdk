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

import {MonoTypeOperatorFunction, Observable, throwError, timer} from 'rxjs';
import {retry} from 'rxjs/operators';

/** Configuration for the exponential backoff retry strategy. */
export interface ExponentialBackoffConfig {
  initialDelayMs: number;
  maxDelayMs: number;
  /**
   * Maximum number of attempts made when retrying. The initial request is also
   * counted. If not provided, retries indefinitely.
   */
  maxAttempts?: number;
  /**
   * Function that can be used to decide if the received error should cause a
   * retry or not. If this returns `false` for an error it is thrown immediately
   * without retrying.
   */
  retryError?: (error: unknown, nextDelayMs: number) => boolean;
}

/**
 * Returns the delay before the next retry attempts for exponential backoff in
 * milliseconds.
 */
export function getNextExponentialBackoffDelayMs(
  attempt: number,
  initialDelayMs: number,
  maxDelayMs: number,
): number {
  return Math.min(maxDelayMs, initialDelayMs * Math.pow(2, attempt - 1));
}

/**
 * Operator that will mirror the source observable unless it throws an error. If
 * an error is thrown the source-observable is re-subscribed with increasing
 * delays.
 */
// Returning functions like this is specifically allowed in the return-only
// generics guidance. Otherwise we would completely lose typing on the called
// function.
export function retryBackoff<T>(
  config: ExponentialBackoffConfig,
): MonoTypeOperatorFunction<T> {
  const {
    initialDelayMs,
    maxDelayMs,
    maxAttempts,
    retryError = () => true,
  } = config;
  return retry<T>({
    delay: (error: unknown, retryCount: number) => {
      const nextDelayMs = getNextExponentialBackoffDelayMs(
        retryCount,
        initialDelayMs,
        maxDelayMs,
      );
      const shouldRetry =
        (maxAttempts === undefined || retryCount < maxAttempts) &&
        retryError(error, nextDelayMs);
      return shouldRetry ? timer(nextDelayMs) : throwError(() => error);
    },
  });
}

/**
 * Operator that will mirror the source observable if no `config` is given. If
 * the config is provided this behaves exactly like `retryBackoff`.
 */
// Returning functions like this is specifically allowed in the return-only
// generics guidance. Otherwise we would completely lose typing on the called
// function.
export function retryBackoffIfConfig<T>(
  config: ExponentialBackoffConfig | null | undefined,
): MonoTypeOperatorFunction<T> {
  return (source: Observable<T>): Observable<T> =>
    config != null ? source.pipe(retryBackoff(config)) : source;
}
