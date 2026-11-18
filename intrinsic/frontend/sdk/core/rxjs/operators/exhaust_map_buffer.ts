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
  Observable,
  ObservableInput,
  OperatorFunction,
  Subscription,
  from,
} from 'rxjs';

/**
 * Projects each source value to an Observable which is merged in the output
 * Observable only if the previous projected Observable has completed. The most
 * recent emission received from the source Observable while a projected
 * Observable is not yet complete is buffered, then immediately projected and
 * merged as soon as the previously projected Observable completes.
 */
export function exhaustMapBuffer<T, R>(
  project: (value: T) => ObservableInput<R>,
): OperatorFunction<T, R> {
  return (source) =>
    new Observable((destination) => {
      let innerSubscription: Subscription | null = null;
      let isComplete = false;
      let hasValue = false;
      let lastValue: T | null = null;

      const projectAndSubscribe = (value: T) => {
        innerSubscription = from(project(value)).subscribe({
          next(value) {
            destination.next(value);
          },
          error(e) {
            destination.error(e);
          },
          complete() {
            innerSubscription?.unsubscribe();
            if (hasValue) {
              hasValue = false;
              const value = lastValue as T;
              lastValue = null;
              projectAndSubscribe(value);
            } else {
              innerSubscription = null;
              if (isComplete) {
                destination.complete();
              }
            }
          },
        });
      };

      const subscription = source.subscribe({
        next(outerValue) {
          if (!innerSubscription) {
            projectAndSubscribe(outerValue);
          } else {
            hasValue = true;
            lastValue = outerValue;
          }
        },
        error(e) {
          destination.error(e);
        },
        complete() {
          isComplete = true;
          if (!innerSubscription) {
            destination.complete();
          }
        },
      });

      // Finalization logic that is called when the source is unsubscribed,
      // completed or throws an error.
      return () => {
        subscription.unsubscribe();
        innerSubscription?.unsubscribe();
      };
    });
}
