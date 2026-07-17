<!--
Copyright 2023 Intrinsic Innovation LLC
-->
# `inctl organization` (alias: `inctl organizations`)

The `inctl organization` command hierarchy provides tools to manage Flowstate
organizations, user memberships, role assignments, and organization lifecycles.
Both singular (`inctl organization`) and plural (`inctl organizations`) forms
are supported.

## Commands Overview

*   `list`: Lists available Flowstate organizations that you have access to.
*   `get <organization>`: Gets details of an organization, printing its protobuf
    definition.
*   `create --display-name=<name> --parent=<parent>`: Creates a new empty
    organization under a specified parent organization (can take up to 15
    minutes).
*   `delete <organization>`: Soft-deletes an organization.
*   `join <token>`: Accepts an invitation to join an organization using an
    invitation token.
*   `members list --org=<organization>`: Lists active members of an organization
    (alias: `member`).
*   `members invite --org=<organization> --email=<email>`: Invites a user by
    email with optional initial roles (short-hand for `invitations create`).
*   `members remove --org=<organization> --email=<email>`: Removes an existing
    user from an organization.
*   `invitations create --org=<organization> --email=<email>`: Creates a user
    invitation by email with optional initial roles (alias: `invitation`).
*   `invitations list --org=<organization>`: Lists pending invitations and their
    tokens.
*   `invitations get <token> --org=<organization>`: Gets details of a pending
    user invitation by token.
*   `invitations resend --org=<organization> --token=<token>`: Resends an
    existing user invitation email.
*   `invitations withdraw --org=<organization> --token=<token>`: Withdraws a
    pending user invitation by token.
*   `roles list --org=<organization>`: Lists all available access control roles
    that can be assigned to users within an organization.
*   `role-bindings list --org=<organization>`: Lists role bindings on an
    organization.
*   `role-bindings grant --org=<organization> --email=<email> --role=<role>`:
    Grants a role to a user on an organization.
*   `role-bindings revoke --org=<organization> --name=<resource-name>`: Revokes
    an existing role binding by its resource name.

--------------------------------------------------------------------------------

## Examples

### 1. `list`

Lists available Flowstate organizations that you have access to.

*   **Canonical Invocation:**

    ```bash
    inctl organization list --parent=my-parent-org
    ```

*   **List all organizations without filtering:**

    ```bash
    inctl organization list
    ```

--------------------------------------------------------------------------------

### 2. `get`

Gets details of an organization, printing its protobuf definition.

*   **Example Invocations:**

    ```bash
    inctl organization get my-org
    inctl organization get --org=my-org
    ```

--------------------------------------------------------------------------------

### 3. `create`

Creates a new empty organization under a specified parent organization. Requires
`--display-name` and `--parent`. Note that setting up isolated infrastructure
can take up to 15 minutes.

*   **Canonical Invocation:**

    ```bash
    inctl organization create --display-name="My Robotics Team" --parent=my-parent-org
    ```

    *(Alternatively, `--org=my-parent-org` can also be used instead of
    `--parent` to specify the parent organization.)*

--------------------------------------------------------------------------------

### 4. `delete`

Soft-deletes an organization. Soft-deleted organizations can be recovered within
30 days by contacting support.

*   **Example Invocations:**

    ```bash
    inctl organization delete my-org
    inctl organization delete --org=my-org
    ```

--------------------------------------------------------------------------------

### 5. `members` (alias: `member`)

Manages organization memberships.

*   **`list`** — View all active members:

    ```bash
    inctl organization members list --org=my-org
    ```

*   **`invite`** — Invite a user by email (`--email`) with optional
    comma-separated roles (`--roles`). This is a short-hand for `inctl
    organization invitations create`:

    ```bash
    inctl organization members invite --org=my-org --email=user@example.com --roles=admin
    ```

*   **`remove`** — Remove an existing user from the organization:

    ```bash
    inctl organization members remove --org=my-org --email=user@example.com
    ```

--------------------------------------------------------------------------------

### 6. `invitations` (alias: `invitation`)

Manages organization invitations.

*   **`create`** — Create an invitation for a user by email (`--email`) with
    optional comma-separated roles (`--roles`):

    ```bash
    inctl organization invitations create --org=my-org --email=user@example.com --roles=admin
    ```

*   **`list`** — View all pending invitations and retrieve invitation tokens:

    ```bash
    inctl organization invitations list --org=my-org
    ```

*   **`resend`** — Resend an existing invitation email by token (`--token`):

    ```bash
    inctl organization invitations resend --org=my-org --token=invitation-token
    ```

*   **`withdraw`** — Cancel/withdraw a pending invitation by token (`--token`):

    ```bash
    inctl organization invitations withdraw --org=my-org --token=invitation-token
    ```

--------------------------------------------------------------------------------

### 7. `roles list`

Lists all available access control roles that can be assigned to users within an
organization.

*   **Example Invocation:**

    ```bash
    inctl organization roles list --org=my-org
    ```

--------------------------------------------------------------------------------

### 8. `role-bindings`

Manages role bindings (granting and revoking roles) for users across an
organization.

*   **`list`** — List role bindings on an organization:

    ```bash
    inctl organization role-bindings list --org=my-org
    ```

*   **`grant`** — Grant a role (`--role`) to a user (`--email`):

    ```bash
    inctl organization role-bindings grant --org=my-org --email=user@example.com --role=admin
    ```

*   **`revoke`** — Revoke an existing role binding by its resource name
    (`--name`):

    ```bash
    inctl organization role-bindings revoke --org=my-org --name=role-binding-identifier
    ```
