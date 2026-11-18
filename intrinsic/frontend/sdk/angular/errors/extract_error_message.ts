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

import {HttpErrorResponse} from '@angular/common/http';

/**
 * An error that may be contained in an HTTP error response if parsing JSON
 * failed. This represents an unexported type in Angular:
 * https://github.com/angular/angular/blob/6073493c5bbb0fa02fe3bdad11c411373fdafd25/packages/common/http/src/response.ts#L143
 *
 * Angular HTTP only returns this on JSON parse issues when using the XHR
 * backend and the fetch backend will return {@link SyntaxError} instead.
 */
interface JsonParseError {
  error: Error;
  text: string;
}

const NO_MESSAGE_SENTINEL = '(no message)';

export const TEST_ONLY = {
  NO_MESSAGE_SENTINEL,
};

/**
 * Extract error message from different types of errors. This includes specific
 * handling for errors thrown by the Angular HTTP client
 * (https://angular.io/guide/http-handle-request-errors).
 *
 * Use `extractErrorMessage` from `extract_error_message.ts` when
 * error handling for the Angular HTTP client is not required.
 * @param err Any error object.
 * @param fallback Message to return if no message can be extracted from `err`.
 *     Defaults to a generic "(no message)" placeholder.
 */
export function extractErrorMessage(
  err: unknown,
  fallback = NO_MESSAGE_SENTINEL,
): string {
  if (err instanceof HttpErrorResponse) {
    if (err.error && typeof err.error === 'string') {
      return err.error;
    }
    // HTTP errors from Angular may contain JSON parsing errors. In that
    // specific case we return a special message to convey this.
    if (isJsonParseError(err)) {
      return 'Response could not be parsed (try reloading this page)';
    }
    // An error status of "0" means that there was an issue where no response is
    // received. Most of the time these are networks errors like CORS or
    // connectivity issues.
    if (err.status === 0) {
      return `Network error (${err.url ?? 'target unknown'})`;
    }
    // If the error response was not handled, extract the error from the
    // contained error rather than the error response instance.
    err = err.error;
  }
  if (hasMessage(err)) {
    return err.message || fallback;
  }
  return rawToString(err) || fallback;
}

function hasMessage(err: unknown): err is {message: string} {
  return (
    err != null &&
    typeof err === 'object' &&
    typeof (err as {message: string}).message === 'string'
  );
}

function isJsonParseError(err: HttpErrorResponse): boolean {
  if (!err.error) {
    return false;
  }
  return (
    // The fetch backend returns the `SyntaxError` from `JSON.parse` directly.
    err.error instanceof SyntaxError ||
    // The XHR backend returns a special error format.
    (err.error as JsonParseError).error instanceof Error
  );
}

function rawToString(err: unknown): string {
  // Returns strings as they come.
  if (typeof err === 'string') {
    return err;
  }
  // We don't want to arbitrarily show other values that may be complex objects
  // or arrays of objects. This also includes `null` and `undefined`.
  return '';
}
