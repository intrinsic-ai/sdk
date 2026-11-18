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

import {Observable, ObservableInput, OperatorFunction, from} from 'rxjs';
import {concatMap, filter, finalize} from 'rxjs/operators';

/**
 * Operator that usually behaves like `concatMap`. If the `exhaustPredicate`
 * returns `true` for a value emitted by the source observable it will be
 * projected and merged to the output observable using `concatMap`. All
 * subsequent emissions from the source observable will be ignored until that
 * previously merged observable completes. This behavior is similar to
 * `exhaustMap`.
 */
export function concatOrExhaustMap<T, R>(
  exhaustPredicate: (value: T) => boolean,
  project: (value: T) => ObservableInput<R>,
): OperatorFunction<T, R> {
  return (source) =>
    new Observable((destination) => {
      let exhausting = false;

      const subscription = source
        .pipe(
          filter(() => !exhausting),
          concatMap((value) => {
            exhausting = exhaustPredicate(value);
            return from(project(value)).pipe(
              finalize(() => {
                if (exhausting) {
                  exhausting = false;
                }
              }),
            );
          }),
        )
        .subscribe({
          next(value) {
            destination.next(value);
          },
          error(err) {
            destination.error(err);
          },
          complete() {
            destination.complete();
          },
        });

      return () => {
        subscription.unsubscribe();
      };
    });
}
