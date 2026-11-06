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

import * as three from 'three';

/** Number of decimal places used when writing a pose to a string. */
export const PRECISION = 6;

// Number of decimal places when writing rotations to a string. Higher accuracy
// is required for representing normalized quaternions.
const ROTATION_PRECISION = 17;

/**
 * DUMMY_PRECISION is the default number of decimal places used to compare poses
 * using the `isApprox` class method.
 */
export const DUMMY_PRECISION = 12;

/**
 * DUMMY_TOLERANCE is the default numerical difference used to compare translations
 * using the `isApprox` class method.
 */
export const DUMMY_TOLERANCE = 10 ** -DUMMY_PRECISION;

/**
 * DUMMY_ANGULAR_PRECISION is the default number of decimal places used to
 * compare rotations using the `isApprox` class method.
 */
export const DUMMY_ANGULAR_PRECISION = 7;

/**
 * DUMMY_ANGULAR_TOLERANCE is the default numerical difference used to compare
 * rotations using the `isApprox` class method.
 */
export const DUMMY_ANGULAR_TOLERANCE = 10 ** -DUMMY_ANGULAR_PRECISION;

/**
 * A read-only view of a Pose. Exposes only the methods that do not mutate the
 * pose, so APIs can accept a pose without being able to modify it.
 */
export interface ReadonlyPose {
  /** Returns a copy of the translation component of this pose. */
  getTranslation(): three.Vector3;
  /** Returns a copy of the rotation component of this pose. */
  getRotation(): three.Quaternion;
  /** Returns a mutable copy of this pose. */
  clone(): Pose;
  /** Returns a new pose representing the inverse of this pose. */
  inverse(): Pose;
  /** Returns a new pose representing the composition of this pose with `p`. */
  composedWith(p: ReadonlyPose): Pose;
  /** Returns a new rotation matrix representing this pose's rotation. */
  rotationMatrix(): three.Matrix3;
  /** Returns a new affine matrix representing this pose. */
  matrix(): three.Matrix4;
  /** Returns true if `p` is numerically within some tolerance of this pose. */
  isApprox(
    p: ReadonlyPose,
    tolerance?: {translationTolerance: number; rotationTolerance: number},
  ): boolean;
  /** Returns true if `p` is numerically equal to this pose. */
  equals(p: ReadonlyPose): boolean;
  /** Transforms the point `v` in place and returns it. */
  applyTo(v: three.Vector3): three.Vector3;
  /** Writes the values of the pose as a 7-element comma-delimited string. */
  toString(): string;
}

/**
 * A pose is a rigid transformation between one coordinate space and another.
 * In this class, we use the threejs definitions for vec3 and quaternions to
 * represent translations and rotations, respectively. A common convention is
 * to name pose variables using the `oldToNew` convention, such that the 'old'
 * and 'new' names refer to coordinate frames. For example
 * 'robotBaseToToolFrame'.
 */
export class Pose implements ReadonlyPose {
  constructor(
    private readonly translation = new three.Vector3(),
    private readonly rotation = new three.Quaternion(),
  ) {}

  /**
   * Returns a new pose that has its position set via 't' and the rotation set
   * to be identity.
   *
   * @param t The translation that this pose will represent.
   * @returns New pose representing the given translation.
   */
  static fromTranslation(t: three.Vector3): Pose {
    return new Pose(t);
  }

  /**
   * Returns a new pose that has its rotation set via 'r' and the translation
   * set to zero.
   *
   * @param r The quaternion that this pose's rotation will represent.
   * @returns New pose representing the given rotation.
   */
  static fromRotation(r: three.Quaternion): Pose {
    return new Pose(new three.Vector3(), r);
  }

  /**
   * Returns a new pose by extracting the rotation and translation components
   * from the passed matrix as if it were an affine transformation. Returns
   * undefined if the transform contains non-unit scaling.
   *
   * @param m The matrix that defines a rigid transform.
   * @returns New pose representing the rigid transform.
   */
  static fromMatrix(m: three.Matrix4): Pose | undefined {
    const p = new Pose();
    const scale = new three.Vector3();
    m.decompose(p.translation, p.rotation, scale);
    // If we get nonunit scale, then the matrix does not define a rigid transform.
    if (!scale.toArray().every((x) => Math.abs(x - 1) < 1e-6)) {
      return undefined;
    }
    return p;
  }

  /**
   * Returns a new pose that has its rotation set via 'm' and the translation
   * set to zero. Returns undefined if the matrix contains non-unit scaling.
   *
   * @param m The matrix that this pose's rotation will represent.
   * @returns New pose representing the given rotation.
   */
  static fromRotationMatrix(m: three.Matrix3): Pose | undefined {
    return Pose.fromMatrix(new three.Matrix4().setFromMatrix3(m));
  }

  /**
   * Returns a new pose copied from the original pose.
   *
   * @param p The pose to copy.
   * @returns New copied pose.
   */
  static copyFrom(p: ReadonlyPose): Pose {
    return new Pose().copy(p);
  }

  /**
   * Returns a new pose representing the identity transform..
   *
   * @returns New identity pose.
   */
  static identity(): Pose {
    return new Pose();
  }

  /**
   * Returns a copy of the 3D vector representing the translation component of
   * this pose.
   *
   * @returns A vector denoting the translation of this pose.
   */
  getTranslation(): three.Vector3 {
    return this.translation.clone();
  }

  /**
   * Returns a copy of the quaternion representing the rotation component of
   * this pose.
   *
   * @returns A quaternion denoting the rotation of this pose.
   */
  getRotation(): three.Quaternion {
    return this.rotation.clone();
  }

  /**
   * Updates the translation of this pose. Leaves the input unchanged.
   *
   * @param t A vector denoting the new translation of this pose.
   */
  setTranslation(t: three.Vector3): this {
    this.translation.copy(t);
    return this;
  }

  /**
   * Updates the rotation of this pose. Leaves the input unchanged.
   *
   * @param r A quaternion denoting the new rotation of this pose.
   */
  setRotation(r: three.Quaternion): this {
    this.rotation.copy(r);
    return this;
  }

  /**
   * Updates this pose to represent the provided pose. Leaves the input
   * unchanged.
   *
   * @param p A pose to copy into this pose.
   */
  copy(p: ReadonlyPose): this {
    return this.setTranslation(p.getTranslation()).setRotation(p.getRotation());
  }

  /**
   * Returns a copy of this pose.
   *
   * @returns A copy of this pose.
   */
  clone(): Pose {
    return Pose.copyFrom(this);
  }

  /**
   * Inverts this pose.
   */
  invert(): this {
    this.rotation.invert();
    this.translation.applyQuaternion(this.rotation).negate();
    return this;
  }

  /**
   * Returns a new pose representing the inverse of this pose.
   *
   * @returns The inverse of this pose.
   */
  inverse(): Pose {
    return this.clone().invert();
  }

  /**
   * Updates this pose to represent the composition of this pose with the
   * provided pose. Leaves the input unchanged.
   *
   * @param p A pose to compose this pose with.
   */
  composeWith(p: ReadonlyPose): this {
    this.translation.add(p.getTranslation().applyQuaternion(this.rotation));

    // Remember: quaternions are multiplied "inside out": If we rotate a vector
    // v via quaternions p and then q, then we need to multiply the vector v by
    // p first: qpvp'q'. So we need to premultiply the quaternion here and store
    // it in this rotation.
    this.rotation.copy(p.getRotation().multiply(this.rotation));
    return this;
  }

  /**
   * Returns a new pose representing the composition of this pose with the
   * provided pose. Leaves the input unchanged.
   *
   * @param p A pose to compose this pose with.
   * @returns A new pose representing the composition between this pose and 'p'.
   */
  composedWith(p: ReadonlyPose): Pose {
    return this.clone().composeWith(p);
  }

  /**
   * Returns a new rotation matrix representing the rotation component of this
   * pose.
   *
   * @returns A new rotation matrix representing this pose's rotation.
   */
  rotationMatrix(): three.Matrix3 {
    return new three.Matrix3().setFromMatrix4(this.matrix());
  }

  /**
   * Returns a new 4x4 matrix representing an affine transformation defined by
   * this pose.
   *
   * @returns A new affine matrix representing this pose.
   */
  matrix(): three.Matrix4 {
    const t = this.translation;
    return new three.Matrix4()
      .makeRotationFromQuaternion(this.rotation)
      .setPosition(t.x, t.y, t.z);
  }

  /**
   * Returns true if the provided pose is numerically within some tolerance of
   * this pose. Both the translation and rotation components are checked
   * numerically.
   *
   * @param p A pose to compare this pose with.
   * @param tolerance An optional tolerance for componentwide comparison.
   * @param tolerance.translationTolerance Maximum allowed distance between
   *     the translations.
   * @param tolerance.rotationTolerance Maximum allowed angle, in radians,
   *     between the rotations.
   * @returns true if all values are within the provided tolerance.
   */
  isApprox(
    p: ReadonlyPose,
    tolerance: {translationTolerance: number; rotationTolerance: number} = {
      translationTolerance: DUMMY_TOLERANCE,
      rotationTolerance: DUMMY_ANGULAR_TOLERANCE,
    },
  ): boolean {
    return (
      this.getTranslation().distanceTo(p.getTranslation()) <=
        tolerance.translationTolerance &&
      Math.abs(this.getRotation().angleTo(p.getRotation())) <=
        tolerance.rotationTolerance
    );
  }

  /**
   * Returns true if the provided pose is numerically equal to this pose.
   * Both the translation and rotation components are checked componentwise.
   * Consider using `isApprox` instead of this to avoid numerical instability.
   *
   * @param p A pose to compare this pose with.
   * @returns true if 'p' is numerically equal to this pose.
   */
  equals(p: ReadonlyPose): boolean {
    return (
      this.translation.equals(p.getTranslation()) &&
      this.rotation.equals(p.getRotation())
    );
  }

  /**
   * Transforms a point via the rigid transformation defined by this pose. This
   * mutates the vector passed, and returns the same reference.
   *
   * @param v The point being transformed.
   * @returns The same point 'v' after having applied the transformation.
   */
  applyTo(v: three.Vector3): three.Vector3 {
    return v.applyQuaternion(this.rotation).add(this.translation);
  }

  /**
   * Writes the values of the pose as a 7-element comma-delimited string.
   *
   * @returns The string representing this pose.
   */
  toString(): string {
    // Dollar store curryable function of (prec, x)
    const toFixedStr = (prec: number) => (x: number) => `${x.toFixed(prec)}`;
    const posStr = this.translation
      .toArray()
      .map(toFixedStr(PRECISION))
      .join(',');
    const rotStr = this.rotation
      .toArray()
      .map(toFixedStr(ROTATION_PRECISION))
      .join(',');
    return posStr + ',' + rotStr;
  }

  /**
   * Parses this pose from a string represented as if it were printed via
   * `toString`. Namely, the passed string is expected to be a comma-delimited
   * list of 7 numerical values, which represent the (x, y, z) values of the
   * translation, and the (x, y, z, w) values of the quaternion representing the
   * rotation, respectively.
   *
   * @param s The string representing a pose.
   */
  setFromString(s: string): this {
    const vals = s.split(',').map(Number);
    if (!vals.every((x) => !isNaN(x))) {
      throw new Error(`Failed to set pose, invalid input value: ${s}!`);
    }

    if (vals.length !== 7) {
      throw new Error(
        `Failed to set pose, Input format: 7 (Position + Quaternion) numbers ` +
          `separated by commas. Got ${vals.length}: ${s}`,
      );
    }

    return this.setTranslation(
      new three.Vector3(vals[0], vals[1], vals[2]),
    ).setRotation(new three.Quaternion(vals[3], vals[4], vals[5], vals[6]));
  }
}
