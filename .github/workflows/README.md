# CI Workflows

## Overview

Every PR to this repo automatically keeps a matching branch and PR open in the [spec-tests](https://github.com/ssvlabs/spec-tests) repository, populated with freshly generated JSON test fixtures. When the ssv-spec PR is merged, the spec-tests PR is finalized and merged automatically.

```
ssv-spec PR opened/updated
        │
        ▼
  Generate JSONs
        │
        ▼
  Push to matching branch in spec-tests (same branch name)
        │
        ▼
  Create / update PR in spec-tests
        │
ssv-spec PR merged
        │
        ▼
  Push final generated JSONs to spec-tests branch
        │
        ▼
  Merge spec-tests PR + delete branch
```

---

## Files

| File | Trigger | Purpose |
|---|---|---|
| `test.yaml` | PR, push to `main` | Build, generate, run tests |
| `sync-spec-tests-pr.yaml` | PR opened / updated / reopened | Sync generated files to spec-tests, create/update PR |
| `sync-spec-tests-merge.yaml` | Push to `main` | Push final files to spec-tests branch, merge spec-tests PR |
| `../.github/actions/generate-spec-tests/action.yaml` | (composite, called by all above) | Set up Go, generate JSON fixtures |

---

## Workflows in Detail

### `test.yaml`

Runs on every PR and every push to `main`.

1. Checkout
2. Run composite action — sets up Go, installs deps, generates JSON fixtures
3. `make test`

### `sync-spec-tests-pr.yaml`

Runs on `pull_request` events: `opened`, `synchronize`, `reopened`.

1. Checkout + generate JSON fixtures (composite action)
2. Get GitHub App token scoped to `spec-tests` only
3. Resolve the bot identity for clean commit attribution
4. Clone `spec-tests`, checkout or create a branch with **the same name as the ssv-spec PR branch**
5. Remove only the previously generated paths (leaves all other files in spec-tests untouched)
6. Copy the freshly generated files, commit, force-push
7. Create the spec-tests PR if it does not exist; update title/body if it does

Branch names `main` and `master` are rejected to prevent accidental overwrites.

### `sync-spec-tests-merge.yaml`

Runs on every push to `main`.

1. Checkout + generate JSON fixtures (composite action)
2. Get GitHub App token scoped to `spec-tests` only
3. Resolve the bot identity
4. Use the built-in `GITHUB_TOKEN` (no app access needed) to look up which ssv-spec PR introduced the merge commit and get its branch name
5. Verify that matching branch exists in `spec-tests` (created by the PR workflow)
6. Push a final sync commit to that branch if anything changed
7. Look up the open spec-tests PR for that branch, merge it, delete the branch

---

## Authentication

| Operation | Token used |
|---|---|
| Read ssv-spec PR metadata | `secrets.GITHUB_TOKEN` (built-in, scoped to this repo) |
| Clone / push / PR on spec-tests | GitHub App token (scoped to `spec-tests` only) |

The GitHub App requires the following permissions on `spec-tests` only:

| Permission | Level |
|---|---|
| Contents | Read and write |
| Pull requests | Read and write |

---

## Required Configuration

Set these in **ssv-spec → Settings → Secrets and variables → Actions**.

**Variables:**

| Name | Example value |
|---|---|
| `SPEC_TESTS_REPO` | `ssvlabs/spec-tests` |
| `SPEC_TESTS_APP_ID` | `12345` |

**Secrets:**

| Name | Value |
|---|---|
| `SPEC_TESTS_APP_PRIVATE_KEY` | Contents of the `.pem` private key file for the GitHub App |

See the [GitHub App setup guide](https://github.com/ssvlabs/ssv-spec/blob/main/.github/workflows/GITHUB_APP_SETUP.md) for step-by-step instructions on creating the app.
