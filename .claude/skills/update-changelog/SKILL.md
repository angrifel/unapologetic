---
name: update-changelog
description: Update CHANGELOG.md following the Keep a Changelog 1.1.0 format
disable-model-invocation: true
---

Update the `CHANGELOG.md` file at the root of the repository following the [Keep a Changelog 1.1.0](https://keepachangelog.com/en/1.1.0/) specification.

## Prerequisites

- `git` must be available in `PATH`
- The current directory must be a git repository

## Instructions

**Important:** Only run the commands explicitly listed in these instructions, exactly as written — do not add, remove, or modify any flags, arguments, or syntax. Do not run additional git commands (e.g., `git tag`, `git log`) beyond what is specified.

**ALWAYS** compare between `BASE_COMMIT` and the current commit. `BASE_COMMIT` will be explained in detail in the next steps.

### Ensure changelog exists
1. If `CHANGELOG.md` does not exist, create it with this structure:

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
```

### Update CHANGELOG.md

1. Read `CHANGELOG.md`.

2. Determine the base commit against which to compare changes. Run the following command to get `BASE_COMMIT`:
   ```
   git describe --tags --match 'v*' --abbrev=0 2>/dev/null || git rev-list --max-parents=0 HEAD | tail -1
   ```
   This always produces a valid commit hash: the most recent `v*` tag, or the first commit if no tags exist. Remember this for the duration of the skill as it will be used in many places.

3. Determine what has changed.
   - Use `git diff --name-only ${BASE_COMMIT}...HEAD` to get the list of files that changed between `BASE_COMMIT` and the current commit.
     - If the list is empty, leave `CHANGELOG.md` unchanged and inform the user that no changes were found since the last release. Stop here.
     - For each file, use `git diff ${BASE_COMMIT}...HEAD -- <file>` to examine its full diff.
     - Thoroughly examine all diffs to understand the changes.
     - ONLY use changes visible in these diffs as the basis for changelog entries.

4. Add entries under the `## [Unreleased]` section. If the section does not exist, create it immediately after the header block.

5. Group changes under the appropriate subsections. Only include subsections that have entries:
   - `### Added` - new features
   - `### Changed` - changes in existing functionality
   - `### Deprecated` - features that have been marked as deprecated and may be removed in future releases
   - `### Removed` - now removed features
   - `### Fixed` - bug fixes
   - `### Security` - vulnerability fixes

6. Each entry should be a concise, human-readable bullet point starting with `- `.

7. Update the comparison link at the bottom of the file. Run `git remote get-url origin` to derive the GitHub `OWNER/REPO`. Run `git rev-parse --abbrev-ref HEAD` to get `CURRENT_BRANCH`. If a `[Unreleased]:` link already exists, replace it. Otherwise, append it.

## Scope

Only include changes that are observable to consumers of this project, or that affect what the project produces. This includes changes to public APIs, user-facing tooling, and CI/CD workflows. Exclude files whose sole purpose is configuring AI assistant behavior (e.g., `CLAUDE.md`, AI prompt files). Include changes to Claude Code skills and other project tooling even when they live under `.claude/` — these affect project workflows and are notable. For language- or project-specific exclusions, consult project conventions.

## Format

- Newest entries go at the top (reverse chronological order).
- Use ISO 8601 dates (`YYYY-MM-DD`) for release entries.
- The `[Unreleased]` comparison link uses the three-dot syntax:
  `[Unreleased]: https://github.com/OWNER/REPO/compare/BASE_COMMIT...CURRENT_BRANCH`
- Do NOT remove or modify existing released entries unless explicitly asked.
- Do NOT duplicate entries that already exist in the changelog.

## Produces

- `CHANGELOG.md` — updated in place with new entries under `## [Unreleased]` and a refreshed comparison link.
