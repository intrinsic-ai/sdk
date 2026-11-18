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
 * Helper function that may be used to check that a switch statement is
 * exhaustive at compile time. No checks are performed at runtime.
 *
 * To perform exhaustiveness checking, assign the variable being switched as the
 * argument to this function in the default case. The compiler will verify that
 * all enum values were checked in the cases.
 */
export function assumeExhaustive(v: never) {}

/**
 * Helper function that may be used to check that a switch statement is
 * exhaustive at compile time with exceptions. No checks are performed at
 * runtime.
 *
 * To perform exhaustiveness checking, assign the variable being switched as the
 * argument to this function in the default case. The compiler will verify that
 * all enum values were checked in the cases. You can provide fallback behavior
   for all specifically allowed values along with the check.
 */
export function assumeExhaustiveAllowing<ALLOW>(v: ALLOW) {}
