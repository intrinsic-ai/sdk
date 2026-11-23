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

import type {FetchFn, HttpMiddleware} from './http_transport';

/** Options for `xsrfMiddleware`. */
export interface XsrfMiddlewareOptions {
  /** The name of the header that carries the token, e.g. `X-XSRF-TOKEN`. */
  readonly headerName: string;
  /** The url that new tokens are fetched from. Defaults to `api/token`. */
  readonly tokenUrl?: string;
}

const DEFAULT_TOKEN_URL = 'api/token';

/**
 * How long XSRF tokens are valid for. Tokens are valid for 24 hours (hardcoded
 * in golang.org/x/net/xsrftoken), this is set slightly lower so that a token
 * can be refreshed before a request fails with it.
 */
const TOKEN_VALIDITY_MS = (24 * 60 - 1) * 60 * 1000;

/**
 * Returns middleware that attaches an XSRF token to every request with a
 * relative url. The token is fetched before the first request, and refreshed
 * whenever it expires or the server rejects it with `401 Unauthorized`.
 */
export function xsrfMiddleware(options: XsrfMiddlewareOptions): HttpMiddleware {
  const tokenUrl = options.tokenUrl ?? DEFAULT_TOKEN_URL;
  const token = new XsrfToken();

  return async (url, init, next) => {
    // The token is only meant for our own origin, so absolute urls are sent
    // without it. Token requests are sent without it, too.
    if (/^https?:\/\//i.test(url) || url === tokenUrl) {
      return next(url, init);
    }
    const fetchToken = () => token.refresh(() => requestToken(tokenUrl, next));
    const nextWithToken = (tokenValue: string) => {
      const headers = new Headers(init?.headers);
      headers.set(options.headerName, tokenValue);
      return next(url, {...init, headers});
    };

    const cachedToken = token.get();
    if (!cachedToken) {
      return nextWithToken(await fetchToken());
    }
    const response = await nextWithToken(cachedToken);
    if (response.status === 401) {
      // The token was rejected (e.g. it expired): refresh it and retry once.
      return nextWithToken(await fetchToken());
    }
    return response;
  };
}

async function requestToken(tokenUrl: string, next: FetchFn): Promise<string> {
  const response = await next(tokenUrl);
  if (!response.ok) {
    throw new Error(
      `Failed to fetch XSRF token: ${response.status} ${response.statusText}`,
    );
  }
  return response.text();
}

/** Caches an XSRF token and deduplicates concurrent refreshes. */
class XsrfToken {
  private token = '';
  private pendingToken?: Promise<string>;

  /** Returns the cached token, or an empty string if there is none. */
  get(): string {
    if (this.isExpired()) {
      this.token = '';
    }
    return this.token;
  }

  /** Fetches a new token. Concurrent calls share a single request. */
  refresh(fetchToken: () => Promise<string>): Promise<string> {
    if (!this.pendingToken) {
      this.pendingToken = fetchToken()
        .then((token) => {
          this.token = token;
          return token;
        })
        .finally(() => {
          this.pendingToken = undefined;
        });
    }
    return this.pendingToken;
  }

  private isExpired(): boolean {
    // Tokens have the form `<opaque string>:<Unix timestamp (ms)>`.
    const timestampMs = Number(this.token.match(/:(\d+)$/)?.[1]);
    if (!timestampMs) return false;
    return Date.now() - timestampMs > TOKEN_VALIDITY_MS;
  }
}
