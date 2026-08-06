# CI Workflows

## Overview

Every PR from a branch in this repo automatically keeps a matching branch and PR open in the [spec-tests](https://github.com/ssvlabs/spec-tests) repository, populated with freshly generated JSON test fixtures. When the ssv-spec PR is merged, the spec-tests PR is finalized and merged automatically. Fork PRs are excluded — see [Fork PRs](#fork-prs).

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
| `sync-spec-tests-pr.yaml` | PR opened / updated / reopened, same-repo only | Sync generated files to spec-tests, create/update PR |
| `sync-spec-tests-merge.yaml` | Push to `main` | Push final files to spec-tests branch, merge spec-tests PR |
| `../actions/generate-spec-tests/action.yaml` | (composite, called by all above) | Set up Go, generate JSON fixtures |

---

## Running Locally

The JSON fixtures are no longer committed to this repo — they are generated into a **sibling** `spec-tests` directory, so generation must run before the tests:

```console
foo@bar:~/ssv-spec$ go generate ./...   # or: make generate-jsons
foo@bar:~/ssv-spec$ make test
```

Fixtures resolve to `<parent-of-ssv-spec>/spec-tests/<module>` (`qbft`, `ssv`, `types`), created automatically by `go generate`. A fresh clone that runs `make test` first will fail — `make test` does not depend on `generate-jsons`.

> **Multiple checkouts:** the path is derived from the repo root's parent and is not namespaced per checkout, so two clones or worktrees sharing a parent directory also share one `spec-tests` and overwrite each other's output. Give each checkout its own parent directory.

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

### Fork PRs

Fork PRs get no mirror branch or mirror PR while open, because `pull_request` from a fork receives no secrets. Their fixtures are still generated and tested by `test.yaml`, which needs no secrets, so correctness is verified either way.

After merge, their fixtures reach spec-tests with the next in-repo merge **that has a mirror branch** — each sync replaces the whole generated tree, so it carries the fork's changes along. The gap: an in-repo PR whose own fixtures match spec-tests never pushes a mirror branch, and if such a PR merges next, its merge run finds no mirror branch, sees the fork's fixture diff, and fails loudly rather than pushing. The sync then completes on the following in-repo merge that does have one.

Contributing from a branch in `ssvlabs/ssv-spec` is the recommended path: it gets the mirror PR preview while the PR is open.

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

To create the App: [Registering a GitHub App](https://docs.github.com/en/apps/creating-github-apps/registering-a-github-app/registering-a-github-app) with the two permissions above, installed on `spec-tests` only. Its App ID goes in `SPEC_TESTS_APP_ID` and its generated `.pem` private key in `SPEC_TESTS_APP_PRIVATE_KEY`.
