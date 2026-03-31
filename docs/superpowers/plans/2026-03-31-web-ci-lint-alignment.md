# Web CI Lint Alignment Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Align the GitHub Actions web validation workflow with the repository's root-level Biome lint entrypoint.

**Architecture:** Add one root-level web lint job in the monorepo workflow, keep admin and user jobs focused on typecheck and build, and include the new job in the final summary so CI reports reflect the actual workspace validation contract.

**Tech Stack:** GitHub Actions YAML, Bun, Turbo, Biome

---

### Task 1: Update workflow structure

**Files:**
- Modify: `.github/workflows/monorepo-check.yml`

- [ ] **Step 1: Add a dedicated web lint job**

Insert a `web-lint` job after `server-validate` with the same `changes` gating used by the existing web jobs, using `./.github/actions/setup-bun`, and run:

```yaml
  web-lint:
    name: Web Lint
    runs-on: ubuntu-latest
    needs: changes
    if: needs.changes.outputs.web == 'true' || needs.changes.outputs.shared == 'true'
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Setup Bun workspace
        uses: ./.github/actions/setup-bun

      - name: Lint web workspace
        working-directory: web
        run: bun run lint
```

- [ ] **Step 2: Remove app-local lint steps**

Delete these workflow steps so app jobs no longer duplicate incomplete lint coverage:

```yaml
      - name: Lint admin app
        working-directory: web/apps/admin
        run: bun run lint
```

```yaml
      - name: Lint user app
        working-directory: web/apps/user
        run: bun run lint
```

- [ ] **Step 3: Extend the workflow summary**

Add `web-lint` to `monorepo-summary.needs`, expose its result in `env`, and include one summary line:

```yaml
          WEB_LINT_RESULT: ${{ needs.web-lint.result }}
```

```yaml
            echo "- web-lint: ${WEB_LINT_RESULT}"
```

- [ ] **Step 4: Verify YAML syntax and contract checks**

Run:

```bash
curl -fsSL -o /tmp/actionlint.tgz https://github.com/rhysd/actionlint/releases/download/v1.7.12/actionlint_1.7.12_linux_amd64.tar.gz
tar -xzf /tmp/actionlint.tgz -C /tmp actionlint
/tmp/actionlint
./script/check-workflow-contract.sh
```

Expected: both commands exit successfully with no workflow errors.