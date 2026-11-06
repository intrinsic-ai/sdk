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
 * @file  This file contains math utility functions.
 */

/**
 * Returns a number whose value is limited to the given range.
 *
 * @param val The value to clamp
 * @param min The lower boundary of the output range
 * @param max The upper boundary of the output range
 * @returns A number in the range [min, max]
 */
export function clamp(val: number, min: number, max: number) {
  return Math.min(Math.max(val, min), max);
}

/**
 * Convert from degree to radians.
 * @param degree The input angle in degrees.
 * @returns the output angle in radians.
 */
export function degToRad(degree: number) {
  return (degree * Math.PI) / 180;
}

/**
 * Convert from radians to degree.
 * @param radian The input angle in radians.
 * @returns the output angle in degree.
 */
export function radToDeg(radian: number) {
  return (radian * 180) / Math.PI;
}

/**
 * Return true if two numbers are within some tolerance.
 * @param x The first number
 * @param y The second number
 * @param tolerance The maximum difference between 'x' and 'y' to be considered
 * @returns True if 'x' and 'y' are within the tolerance provicd
 */
export function scalarsAlmostEqual(
  x: number,
  y: number,
  tolerance = 1e-6,
): boolean {
  return Math.abs(x - y) <= tolerance;
}

/**
 * Returns a number whose value is rounded to the number of decimal places
 * denoted by the precision.
 *
 * @param value The number to be rounded.
 * @param precision The precision of the rounding. Positive integer.
 * @returns The number value rounded up to a precision defined by user.
 */
export function toPrecision(value: number, precision: number): number {
  if (isNaN(value)) return value;
  if (isNaN(precision)) return value;
  if (precision < 0) return value;
  const scalar = 10 ** precision;
  return Math.round(value * scalar) / scalar;
}

/**
 * Returns the number of decimal places in the given number.
 *
 * @param value The number to infer the precision of.
 * @returns The number of decimal places in the given number.
 */
export function inferPrecision(value: number): number {
  if (isNaN(value)) return 0;

  const valueStr = value.toString();
  const decimalPointIndex = valueStr.indexOf('.');

  if (decimalPointIndex === -1) {
    return 0;
  }

  return valueStr.length - decimalPointIndex - 1;
}
