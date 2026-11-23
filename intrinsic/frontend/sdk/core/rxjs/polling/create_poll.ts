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

import {Observable, ObservableInput, Subscription, defer} from 'rxjs';

import {ExponentialBackoffConfig, retryBackoffIfConfig} from '../operators';

/**
 * Function to be executed in a poll. When this function throws an error the
 * entire poll will be stopped. To prevent this from happening errors may be
 * caught with `try/catch` or `catchError` (for observables) where necessary.
 */
export type PollingFunction<T> = () => ObservableInput<T>;

/** Options to configure a poll. */
export interface PollingOptions {
  /**
   * If this is `true` the interval will begin only after the previous poll has
   * completed. Regardless of how long the poll takes to complete, the next poll
   * only begins after the interval. Otherwise the execution time of the polling
   * function is considered part of the interval. This means that if the polling
   * function takes longer than the polling interval the next poll is dispatched
   * immediately after the previous one completes.
   */
  intervalAfterComplete?: boolean;
  /**
   * If this config is present, the polls are retried on errors with exponential
   * backoff.
   */
  retryBackoffConfig?: ExponentialBackoffConfig;
}

/**
 * Creates a polling observable with the given options wrapped as an observable.
 * Subscribing to the returned poll observable will call the polling function
 * every `intervalMs` milliseconds. The result from the polling function is
 * emitted to the poll observable every iteration.
 */
export function createPoll<T>(
  pollingFn: PollingFunction<T>,
  intervalMs: number,
  options: PollingOptions = {},
): Observable<T> {
  const {intervalAfterComplete, retryBackoffConfig} = options;
  return new Observable<T>((subscriber) => {
    let stop = false;
    let currentPollTimeout = -1;
    let pollSubscription: Subscription | null = null;
    // Recursive polling function that will invoke the `pollingFn` and pipe all
    // outputs to the returned observable. When the `pollingFn` completes, the
    // function will be called again after the `intervalMs` has passed.
    function poll() {
      const pollStart = Date.now();
      pollSubscription?.unsubscribe();
      pollSubscription = defer(() => pollingFn())
        .pipe(retryBackoffIfConfig(retryBackoffConfig))
        .subscribe({
          next(value: T) {
            subscriber.next(value);
          },
          error(error: unknown) {
            subscriber.error(error);
          },
          complete() {
            pollSubscription?.unsubscribe();
            pollSubscription = null;
            // Stop recursing if the outer observable was unsubscribed.
            if (stop) return;
            // If `intervalAfterComplete` is set to `true` then we always
            // wait `intervalMs` before scheduling the next invocation of
            // the polling function. Otherwise the execution time of the
            // polling function is included in the interval by subtracting
            // it. If execution took longer than the specified interval we
            // schedule immediately (0ms).
            const msToNextPoll = intervalAfterComplete
              ? intervalMs
              : Math.max(pollStart + intervalMs - Date.now(), 0);
            // We store the ID of the timeout used to schedule the next
            // poll so it can be canceled when the returned observable is
            // unsubscribed.
            currentPollTimeout = setTimeout(() => {
              currentPollTimeout = -1;
              poll();
            }, msToNextPoll);
          },
        });
    }
    // Begin polling by calling the recursive polling function.
    poll();
    return () => {
      stop = true;
      clearTimeout(currentPollTimeout);
      pollSubscription?.unsubscribe();
    };
  });
}
