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

/** Frontend-friendly representation of CollisionRules proto. */
export interface WorldCollisionRule {
  leftId: string;
  rightId?: string;
  excluded: boolean;
  margin?: number;
}

/**
 * The CollisionRules class manages the frontend model of the collision rules
 * between objects and entities, namely which entities are excluded from
 * collisions. It is expected that this model is a reflection of the backend
 * data model, and that this class is not in charge of maintaining coherence.
 */
export class CollisionRules {
  private collisionRules: WorldCollisionRule[] = [];

  /**
   * Remove all collision rules that affect the specified pair.
   * For this operation we only match wildcards with wildcards.
   */
  removeRule(
    leftId: string,
    rightId: string | undefined,
    excluded: boolean,
    margin?: number,
  ) {
    const existingRules = this.collisionRules.filter(
      (rule: WorldCollisionRule) => {
        if (margin) {
          return (
            rule.leftId === leftId &&
            rule.rightId === rightId &&
            rule.excluded === excluded &&
            rule.margin === margin
          );
        }
        return (
          rule.leftId === leftId &&
          rule.rightId === rightId &&
          rule.excluded === excluded
        );
      },
    );
    // Rule doesn't exist. Nothing to do.
    if (!existingRules || existingRules.length === 0) {
      return;
    }
    for (const rule of existingRules) {
      if (this.collisionRules.indexOf(rule) !== -1) {
        this.collisionRules.splice(this.collisionRules.indexOf(rule), 1);
      }
    }
  }

  /**
   * Utility function to see if the definition that is passed in is already
   * covered by one of the existing rules.
   */
  ruleEquivalentExists(
    leftId: string,
    rightId: string | undefined,
    excluded: boolean,
    margin?: number,
  ): WorldCollisionRule[] {
    const matchedRules = [];
    for (const rule of this.collisionRules) {
      if (rule.excluded !== excluded) continue;
      if (margin && rule.margin !== margin) continue;
      // Wildcard rules, only match other wildcard rules.
      if (rightId === undefined && rule.rightId !== undefined) continue;
      // If rule has a wildcard, we match any right value, and a swapped match
      // with the left.
      if (rule.rightId === undefined) {
        if (rule.leftId === leftId) {
          matchedRules.push(rule);
        }
      } else {
        if (leftId === rule.leftId && rightId === rule.rightId) {
          matchedRules.push(rule);
        }
      }
    }
    return matchedRules;
  }

  /**
   * Add a collision rule with the specified parameters. This method will always
   * attempt to add a specific object-pair rule into an existing one.
   */
  addRule(
    leftId: string,
    rightId: string | undefined,
    excluded: boolean,
    margin?: number,
  ) {
    // Excluded false and margin 0 are the default behavior when no rule
    // is present, so in this case we can remove the old rule if any, and we
    // don't need to add.
    if (excluded === false && margin === 0) {
      this.removeRule(leftId, rightId, true);
      return;
    }

    // If this rule is covered by other rules, we don't need to add it.
    const existingRules = this.ruleEquivalentExists(
      leftId,
      rightId,
      excluded,
      margin,
    );
    if (existingRules.length > 0) {
      return;
    }

    // If we are adding X:*, we can remove any X:specific rules
    if (rightId === undefined) {
      const rules = this.collisionRules.filter(
        (rule) => rule.excluded === excluded && rule.margin === margin,
      );
      for (const rule of rules) {
        this.removeRule(leftId, rule.rightId, excluded, margin);
      }
    }

    // If this is a margin change, we can modify existing rules.
    if (margin && margin !== 0) {
      const rules = this.collisionRules.filter(
        (rule) =>
          rule.excluded === false &&
          ((rule.leftId === leftId && rule.rightId === rightId) ||
            (rule.leftId === rightId && rule.rightId === leftId)),
      );
      if (rules.length > 0) {
        for (const rule of rules) {
          rule.margin = margin;
        }
        return;
      }
    }

    const newRule: WorldCollisionRule = {
      leftId,
      rightId,
      excluded,
      margin: margin ?? 0,
    };
    this.collisionRules.push(newRule);
  }

  *getRules(): Generator<WorldCollisionRule, void, undefined> {
    yield* this.collisionRules;
  }

  setRules(rules: WorldCollisionRule[]) {
    this.collisionRules = [...rules];
  }
}
