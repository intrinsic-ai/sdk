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
 * Simple typed event dispatcher for managing listener callbacks.
 * @internal
 */
export class EventEmitter<T> {
  private readonly listeners = new Set<(arg: T) => void>();

  /** Synchronously invokes all currently registered listeners with the emitted value. */
  emit(value: T): void {
    for (const listener of this.listeners) {
      listener(value);
    }
  }

  /** Registers a listener callback and returns a handle to unsubscribe. */
  subscribe(listener: (arg: T) => void): () => void {
    this.listeners.add(listener);
    return () => this.listeners.delete(listener);
  }
}

/**
 * Central distribution point for routing errors, warning notices, and diagnostic debug logs to registered listeners.
 * @internal
 */
export class DiagnosticRouter {
  private readonly errors = new EventEmitter<unknown>();
  private readonly warnings = new EventEmitter<string>();
  private readonly debugs = new EventEmitter<string>();

  /** Emits an error event to all error listeners. */
  error(err: unknown): void {
    this.errors.emit(err);
  }

  /** Emits a non-fatal warning notice to all warning listeners. */
  warn(message: string): void {
    this.warnings.emit(message);
  }

  /** Emits a diagnostic debug log to all debug listeners. */
  debug(message: string): void {
    this.debugs.emit(message);
  }

  /** Registers a callback for operational errors. Returns an unsubscribe handle. */
  onError(listener: (error: unknown) => void): () => void {
    return this.errors.subscribe(listener);
  }

  /** Registers a callback for warning notices. Returns an unsubscribe handle. */
  onWarning(listener: (message: string) => void): () => void {
    return this.warnings.subscribe(listener);
  }

  /** Registers a callback for diagnostic debug logs. Returns an unsubscribe handle. */
  onDebug(listener: (message: string) => void): () => void {
    return this.debugs.subscribe(listener);
  }
}
