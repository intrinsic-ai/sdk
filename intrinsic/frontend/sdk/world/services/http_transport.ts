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

/** A function with the same signature as `fetch` (restricted to url strings). */
export type FetchFn = (url: string, init?: RequestInit) => Promise<Response>;

/**
 * Middleware that can inspect or modify a request before passing it on to
 * `next`, and inspect or modify the response that `next` resolves to.
 */
export type HttpMiddleware = (
  url: string,
  init: RequestInit | undefined,
  next: FetchFn,
) => Promise<Response>;

/**
 * A thin wrapper around `fetch` that the world services send their requests
 * through. It makes the services easy to test with a fake transport, and
 * applies a list of middleware (e.g. `xsrfMiddleware`) to every request.
 */
export class HttpTransport {
  /**
   * @param middleware Applied to every request in order: the first middleware
   *     sees the request first and the response last.
   */
  constructor(private readonly middleware: readonly HttpMiddleware[] = []) {}

  /** Sends a request with `fetch`, passing it through the middleware. */
  fetch(url: string, init?: RequestInit): Promise<Response> {
    // Invoke all middleware in a chain: each middleware calls the next one,
    // and the last one calls `fetch`.
    const dispatch = (
      index: number,
      url: string,
      init?: RequestInit,
    ): Promise<Response> => {
      const middleware = this.middleware[index];
      if (!middleware) {
        return globalThis.fetch(url, init);
      }
      return middleware(url, init, (nextUrl, nextInit) =>
        dispatch(index + 1, nextUrl, nextInit),
      );
    };
    return dispatch(0, url, init);
  }
}
