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

/** A cache entry: the pending or completed result of a load. */
interface CacheEntry<V> {
  readonly promise: Promise<V>;
  /** Aborts the load, if it is still pending. */
  readonly abortController: AbortController;
  /** Set once the load has completed successfully. */
  result?: {readonly value: V};
}

/**
 * A map from keys to the cached results of async loads.
 *
 * The first request for a key starts its load. Requests for the key while the
 * load is pending share it, and requests after it completed get the cached
 * value. If a load fails, its entry is removed, so the next request for the key
 * starts a new load. Entries can also be removed explicitly with `delete` or
 * `clear`, or replaced by a fresh load with `reload`.
 *
 * Each load gets an `AbortSignal`, which is aborted if its entry is deleted or
 * reloaded while the load is still pending. Loads should pass it on to any
 * cancellable work (e.g. `fetch`), so that obsolete loads stop early.
 *
 * ```ts
 * const models = new AsyncCache<string, Model>();
 * // Only the first request for 'a' fetches it.
 * const model = await models.getOrLoad('a', (signal) =>
 *   fetchModel('a', {signal}),
 * );
 * ```
 */
export class AsyncCache<K, V> {
  private readonly entries = new Map<K, CacheEntry<V>>();

  /**
   * Returns the (pending or cached) result for `key`, starting `load` if there
   * is none. `load` is passed a signal that is aborted if the entry is deleted
   * or reloaded before the load completes.
   */
  getOrLoad(key: K, load: (signal: AbortSignal) => Promise<V>): Promise<V> {
    const existing = this.entries.get(key);
    if (existing) return existing.promise;
    const abortController = new AbortController();
    const entry: CacheEntry<V> = {
      abortController,
      // A load that completes after its entry was deleted or replaced only sets
      // the result of its own (detached) entry, so it can't overwrite a newer
      // entry.
      promise: load(abortController.signal).then(
        (value) => {
          entry.result = {value};
          return value;
        },
        (error: unknown) => {
          if (this.entries.get(key) === entry) {
            this.entries.delete(key);
          }
          throw error;
        },
      ),
    };
    this.entries.set(key, entry);
    return entry.promise;
  }

  /**
   * Starts a fresh load for `key`, replacing its current entry (if any), and
   * returns its result. A previous load that is still pending is aborted (see
   * `delete`).
   */
  reload(key: K, load: (signal: AbortSignal) => Promise<V>): Promise<V> {
    this.delete(key);
    return this.getOrLoad(key, load);
  }

  /**
   * Removes the entry for `key`, so that the next request starts a new load.
   * Returns true if there was an entry for `key`.
   *
   * If the entry's load is still pending, its signal is aborted. Requests that
   * are already waiting for it get whatever the load then returns, e.g. an
   * `AbortError` if it honors the signal. Its result is never cached.
   */
  delete(key: K): boolean {
    const entry = this.entries.get(key);
    if (!entry) return false;
    this.entries.delete(key);
    if (!entry.result) {
      entry.abortController.abort();
    }
    return true;
  }

  /**
   * Removes all entries, aborting all pending loads (see `delete`).
   */
  clear(): void {
    for (const entry of this.entries.values()) {
      if (!entry.result) {
        entry.abortController.abort();
      }
    }
    this.entries.clear();
  }

  /** Returns true if `key` has a pending or completed load. */
  has(key: K): boolean {
    return this.entries.has(key);
  }

  /** Returns true if the load for `key` has completed successfully. */
  isLoaded(key: K): boolean {
    return this.entries.get(key)?.result !== undefined;
  }

  /** Returns the loaded value for `key`, or undefined if it isn't loaded. */
  getLoaded(key: K): V | undefined {
    return this.entries.get(key)?.result?.value;
  }

  /** Returns the number of keys with a pending or completed load. */
  size(): number {
    return this.entries.size;
  }

  /**
   * Iterates over the values of all completed loads, e.g. to dispose of them
   * before calling `clear`.
   */
  *loadedValues(): IterableIterator<V> {
    for (const entry of this.entries.values()) {
      if (entry.result) yield entry.result.value;
    }
  }

  /** Returns the number of keys whose load has completed successfully. */
  getNumLoaded(): number {
    let numLoaded = 0;
    for (const entry of this.entries.values()) {
      if (entry.result) numLoaded++;
    }
    return numLoaded;
  }
}
