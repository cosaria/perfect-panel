# Branch Protection Audit

This file is the cutover checklist for moving branch protection from subdirectory-era workflow names to the root control plane.

## Current State Audit

Fill the `Current GitHub Setting` column from the repository settings page before cutover.

| Scope | Current GitHub Setting | Root Target | Required After Cutover | Rollback Mapping |
| --- | --- | --- | --- | --- |
| Required status checks | Inspect current rules in GitHub UI | `Monorepo Summary` | Yes | Restore prior required checks and old workflow files together |
| Code owner reviews | Inspect current rules in GitHub UI | `.github/CODEOWNERS` enforced for `.github/`, `server/`, `web/` | Yes | Disable code-owner enforcement and restore prior workflow files together |
| Pull request review count | Inspect current rules in GitHub UI | Keep existing review count unless maintainer chooses otherwise | Manual | Restore prior count |

## Cutover Window

1. Confirm the root workflows are green on at least one feature-branch PR.
2. Confirm publishing workflows passed no-op or shadow validation.
3. Update required checks to `Monorepo Summary`.
4. Enable code-owner enforcement for `.github/`, `server/`, and `web/`.
5. Verify the new checks appear on a fresh PR.
6. Delete the superseded subdirectory workflow files in the same cutover window.

## Rollback

1. Restore the deleted subdirectory workflow files from the cutover commit.
2. Revert required checks to the prior mapping recorded above.
3. Revert code-owner enforcement if it blocked unrelated work.
4. Re-run CI on a fresh PR to confirm the old control plane is live again.
