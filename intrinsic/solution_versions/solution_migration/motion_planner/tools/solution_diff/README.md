<!--
Copyright 2026 Intrinsic Innovation LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->
# solution_diff

> [!NOTE]
> TODO(b/539662184): Added solely for the purpose of facilitating the `MotionPlannerService` (MPS)
> Asset migration (in go/intrinsic-motion-planner-service-assetization-tdd), especially for
> Solutions which have Solution Building Library (SBL) Python script as the Source-of-Truth (SoT).
> Intended to be used to generate a `diff` file which will be fed into a GenAI tool to
> update/create a Pull-Request (PR) on the original SBL Python script. We will delete this
> functionality by 2027-07-20 or earlier as soon as the MPS Asset migration is complete.

`solution_diff` shows the difference between two snapshots of a Solution, with Skill parameter field names and values.
For example, use it to see what `inctl solution migrate --migration motion-planner` changed, and port the same change into the SBL Python scripts that generate a Solution.

In a raw Solution snapshot proto, Skill parameters are stored inside `google.protobuf.Any` messages as encoded bytes.
`solution_diff` fetches the parameter descriptors of each Catalog Asset from the Asset Catalog and decodes the parameters.
Parameters of sideloaded Assets stay encoded, because the Asset Catalog has no descriptors for them.

## Usage

1. Log in: `inctl auth login --org <ORG>`.
2. Find the Solution branch ID (and optionally snapshot IDs):
   - **Branch ID:** In Portal, open the `⋮` menu of the Solution and select **Copy Solution ID** (`..._BRANCH`).
   - **Default snapshots (migration or latest commit):** When `--after_snapshot_id` is omitted, `solution_diff` uses the tip snapshot of `--solution`. When `--before_snapshot_id` is omitted, it uses the primary parent of the `after` snapshot. For example, right after stopping the Solution and running `inctl solution migrate <BRANCH_ID> --migration motion-planner --org <ORG>`, running `solution_diff` with only `--solution=<BRANCH_ID>` diffs the pre-migration snapshot against the migrated snapshot.
     Do not revert an uncommitted migration snapshot before you run the tool: a reverted uncommitted snapshot is deleted.
   - **Explicit snapshots:** Pass `--before_snapshot_id=<SNAP>` and/or `--after_snapshot_id=<SNAP>` to compare specific snapshots (for example, the snapshot ID printed by `SUCCESS: ... (Snapshot: <id>)`).
3. Run the tool:

   - **From an SDK Bazel workspace (`@ai_intrinsic_sdks`):**

     ```sh
     bazel run @ai_intrinsic_sdks//intrinsic/solution_versions/solution_migration/motion_planner/tools/solution_diff -- \
       --solution=<BRANCH_ID> --org=<ORG> \
       [--before_snapshot_id=<SNAP>] [--after_snapshot_id=<SNAP>] \
       --out=$HOME/solution_diff_out
     ```

   - **From a Python environment with the Intrinsic `solutions` wheel installed:**

     ```sh
     python3 -m intrinsic.solution_versions.solution_migration.motion_planner.tools.solution_diff.solution_diff \
       --solution=<BRANCH_ID> --org=<ORG> \
       [--before_snapshot_id=<SNAP>] [--after_snapshot_id=<SNAP>] \
       --out=$HOME/solution_diff_out
     ```

4. Open `~/solution_diff_out/bt_changes.diff`.
   It holds the hunks that change a Skill call.
   Hunks longer than 200 lines (for example, a Process copied into the Solution) are only in `full.diff`; the tool prints how many.
   `before.textproto` and `after.textproto` hold the full decoded snapshots.

## Example: the MPS migration

For the Trinity CNC Solution, the migration added a Motion Planner Service binding to six Skill calls:

```diff
 skill_id: "ai.intrinsic.collect_calibration_data"
 parameters {
   [type.intrinsic.ai/assets/ai.intrinsic.collect_calibration_data/intrinsic_proto.skills.CollectCalibrationDataParams] {
     calibration_case: "MOVING_CAMERA"
     motion_planner_service {
+      name: "motion_planner_service"
     }
```

It also pointed one call at a local copy of a Catalog Process:

```diff
-skill_id: "ai.intrinsic.run_single_camera_multi_view_inference"
+skill_id: "run_single_camera_multi_view_inference"
```
