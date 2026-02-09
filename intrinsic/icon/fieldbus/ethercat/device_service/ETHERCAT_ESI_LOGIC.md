<!--
Copyright 2023 Intrinsic Innovation LLC
-->
# EtherCAT ESI Parsing & Variable Mapping Logic

This document describes how the `DeviceService` interprets EtherCAT Slave Information (ESI) files and maps logical user references to concrete hardware configurations. For detailed technical specs on the underlying XML structures, refer to the [ETG Specifications](#6-references).

## 1. Architectural Strategy

The service uses a "Fetch-and-Index" strategy to handle the complexity of ESI XML structures:
1.  **Fetching:** ESI files are retrieved as a bundle from the Data Asset service.
2.  **Indexing:** The service performs a deep scan of the XML to build high-performance lookup tables for Objects, Names, and PDOs.
3.  **Resolution:** User requests are resolved against these indices to generate EBI (ENI Builder Input) instructions.

## 2. Parsing Logic (Phase A)

### Localized Name Handling
Modern ESIs (e.g., VIPA 053-1EC00) use the `LcId` attribute to support multiple languages.
```xml
<Name LcId="1033">Inputs</Name>
<Name LcId="1031">Eingänge</Name>
```
The service implements a **Priority Picker**:
*   It searches for `LcId="1033"` (English).
*   If no English name is found, it defaults to the first `<Name>` element encountered without an LcId field (hoping this might be English).
*   If no `<Name>` element without an `LcId` is found, it choses the 1st `<Name>` element in the XML document for that entry, regardless of its `LcId`.

### Modular Device Profiles (MDP)
Devices like a modular EtherCAT coupler (e.g., from VIPA) use a modular structure where variables are often defined inside individual modules rather than a global dictionary.
*   **InfoReferences:** The service follows `<InfoReference>` tags to load external module XMLs. Currently, it supports exactly one level of referencing; nested references within secondary files are unsupported and will trigger an error.
*   **PDO Entry Indexing:** Objects defined directly within a `<TxPdo>` or `<RxPdo>` tag (common in MDP) are automatically indexed and assigned the correct direction (Transmit/Receive) based on their parent PDO.

## 3. Variable Resolution (Phase B)

Users provide a `VariableReference { pdo, object }`. The resolution follows an **Intent-Aware, Context-First** flow to ensure accuracy, especially for bi-directional (RT) objects.

The results of this resolution (mapping of names to hardware addresses and PDOs) are returned in the `ResolvedConfiguration` field of the `GetConfigurationResponse`. For debugging, detailed step-by-step resolution logs can be enabled using the `enable_variable_resolution_tracing` flag in `DeviceServiceOptions`.

You can also download the generated ENI file using `inctl ethercat download_eni...`.

### Step 1: PDO Resolution
The service first identifies the target PDO:
1.  **Explicit:** If `pdo` is provided (Index or Name), use that specific PDO.
2.  **Implicit (Engine Choice):** If `pdo` is empty, the service identifies potential active PDOs. Final implicit selection is deferred until the Object's requirements (direction/mapping) are known in Step 2.

### Step 2: Object Resolution
Resolution of the `object` is scoped by the PDO context and the **Preferred Direction** (Tx/Rx):
1.  **Strict Address:** If `object` looks like `#x6041.0`, it is resolved directly via the global index.
2.  **Contextual Search:** If a PDO context exists (Step 1), the service searches **only** the variables defined within that specific PDO. This eliminates global ambiguity.
3.  **Global Search:** If not found in the specific PDO, the service searches the global **ObjectNameIndex**.
4.  **Ambiguity Filtering (Directional Preference):** If a name search returns multiple candidates (common for "Physical outputs" or "Status word" in complex ESIs), the service filters candidates by:
    *   The direction of the current PDO context.
    *   The **Preferred Direction** associated with the logical interface (e.g., Joint State = Tx, Joint Command = Rx).

### Step 3: Directional Re-Resolution (Implicit only)
If the PDO was implicit (empty) in Step 1, the service selects the best mappable PDO:
1.  **Pass 1 (Reuse):** Prefer PDOs already activated for other variables. This minimizes EBI changes.
2.  **Pass 2 (Default Fixed):** Use default-active fixed PDOs defined in the ESI.
3.  **Pass 3 (Matching Search):** Search for any mappable PDO that matches the **Preferred Direction**.
4.  **Pass 4 (Fallback):** Fallback to any valid mappable PDO.

## 4. EBI Instruction Generation (Phase C)

Once resolved, the service generates a set of EBI instructions that inform the ENI Builder component on how to align the hardware configuration with the user's request:
*   **PDO Activation:** If a variable is mapped to a PDO that is disabled by default, it is added to `pdo_exclusions_to_remove`.
*   **Conflict Resolution:** If the selected PDO excludes a default-active PDO, the default one is added to `pdo_exclusions_to_add`.
*   **Dynamic Mapping:** If the object is not in the PDO's default entry list, an `ObjectAddition` instruction is generated (fails if the PDO is `Fixed="1"`).
*   **Direction Validation:** The service strictly validates that the selected PDO's direction is compatible with the object's `PdoMapping` flags. (e.g. an "RT" object is valid in both Rx and Tx PDOs).

## 5. Examples

### Example 1: Bi-directional (RT) Object
This example uses a standard CiA 402 drive (e.g., Beckhoff EL7041).
*   **User Input:** Interface="Joint Command", Object="Target position"
*   **Logic:** 
    *   "Target position" (#x607A.0) is marked "RT" (Mappable to both Rx and Tx PDOs).
    *   The `DeviceService` logic associates internal interface types (like command inputs) with a **Preferred Direction** (in this case, `Rx`).
*   **Result:** The system picks an RxPDO (e.g., #x1600) to host the variable, even if a TxPDO (e.g., #x1A00) also contains it.

### Example 2: Reusing Active PDOs
*   **User Input:** Variable A in #x1A01. Variable B (pdo="") also fits in #x1A01.
*   **Logic:** Pass 1 (Reuse) identifies #x1A01 as already active.
*   **Result:** Reuses #x1A01, generating no additional activation/exclusion instructions.

## 6. References

For more detailed information on EtherCAT configuration concepts, refer to the following EtherCAT Technology Group (ETG) specifications:
*   **ETG.2000:** EtherCAT Slave Information (ESI) Specification. Defines the structure and content of ESI XML files.
*   **ETG.2100:** EtherCAT Network Information (ENI) Specification. Defines the structure of the network configuration file used by EtherCAT masters.