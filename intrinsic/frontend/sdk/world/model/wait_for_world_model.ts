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

import {firstValueFrom, merge, of} from 'rxjs';
import {delay, filter} from 'rxjs/operators';

import {WorldModel} from './world_model';

const LOAD_TIMEOUT_MS = 3000;
const MAX_NUM_LOAD_RETRIES = 100;

/** Exported for testing purposes only. */
export const TEST_ONLY = {
  LOAD_TIMEOUT_MS,
  MAX_NUM_LOAD_RETRIES,
};

/**
 * Waits for the given WorldModel to finish loading, re-requesting a sync if it
 * doesn't load within LOAD_TIMEOUT_MS. Gives up (and logs an error) after
 * MAX_NUM_LOAD_RETRIES attempts.
 */
export async function waitForWorldModelToLoad(
  model: WorldModel,
): Promise<void> {
  // Every three seconds for five minutes.
  let retry = 0;
  while (!model.isLoaded() && retry < MAX_NUM_LOAD_RETRIES) {
    const loader = merge(
      model.loaded.pipe(filter((loaded) => loaded)),
      of(false).pipe(delay(LOAD_TIMEOUT_MS)),
    );
    const loaded = await firstValueFrom(loader);
    if (loaded) {
      break;
    }

    retry = retry + 1;
    if (retry === MAX_NUM_LOAD_RETRIES) {
      console.error(
        `Could not load world ${model.id} after five minutes! Aborting.`,
      );
    } else {
      console.warn(
        `World ${model.id} didn't load in 3s. Retrying (${retry}/${MAX_NUM_LOAD_RETRIES})...`,
      );
      void model.sync(); // Don't await for it -- it will notify on model.loaded
    }
  }
}
