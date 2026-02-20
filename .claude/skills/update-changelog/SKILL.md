---
name: update-changelog
description: Update CHANGELOG.md following the Keep a Changelog 1.1.0 format
disable-model-invocation: true
---

Update the `CHANGELOG.md` file at the root of the repository following the [Keep a Changelog 1.1.0](https://keepachangelog.com/en/1.1.0/) specification.

## Instructions

**Important:** Only run the commands explicitly listed in these instructions. Do not run additional git commands (e.g., `git tag`, `git describe`) beyond what is specified. You may run `git diff` to obtains the differences of files obtained through `make git-diff-stat-since-last-release`

**ALWAYS** compare between the `BASE COMMIT` and the current commit. `BASE COMMIT` will be explained in detail in the next steps.

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

2. Determine the base commit against which to compare changes. Use `make git-rev-list-last-release-commit` to get the `BASE COMMIT`. This target always produces a valid commit hash, even when no prior releases exist. Remember this for the duration of the skill as this will be used in many places.

3. Determine what has changed.
   - Use `make git-diff-stat-since-last-release` to get the diff stat between the `BASE COMMIT` and the current commit.
     - Thoroughly examine the entire diff contents of the files obtained here to understand all changes.
     - ONLY use changes visible in this diff as the basis for changelog entries.

4. Add entries under the `## [Unreleased]` section. If the section does not exist, create it immediately after the header block.

5. Group changes under the appropriate subsections. Only include subsections that have entries:
   - `### Added` - new features
   - `### Changed` - changes in existing functionality
   - `### Deprecated` - features that have been marked as deprecated and may be removed in future releases
   - `### Removed` - now removed features
   - `### Fixed` - bug fixes
   - `### Security` - vulnerability fixes

6. Each entry should be a concise, human-readable bullet point starting with `- `.

## Format rules

- Newest entries go at the top (reverse chronological order).
- Use ISO 8601 dates (`YYYY-MM-DD`) for release entries.
- Include a comparison link at the bottom of the file for the unreleased changes. Derive the GitHub repository URL from `git remote get-url origin`.
- The `[Unreleased]` comparison link should compare `BASE COMMIT` to `HEAD` using the three-dot syntax: `[Unreleased]: https://github.com/OWNER/REPO/compare/BASE_COMMIT...{CURRENT_BRANCH}`. All information needed for this link is already available from the `BASE COMMIT` obtained in step 2, section `Update CHANGELOG.md` — do not query for tags or other version information.
- Do NOT remove or modify existing released entries unless explicitly asked.
- Do NOT duplicate entries that already exist in the changelog.
