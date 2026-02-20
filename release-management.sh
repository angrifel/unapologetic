#!/usr/bin/env sh

_program_name="release-management.sh"

last_release_tag () {
  git describe --tags --match 'v*' --abbrev=0 2>/dev/null
}

last_release_rev () {
  last_release_tag || git rev-list --max-parents=0 HEAD | tail -1
}

semver_prerelease_core () {
  { last_release_tag || echo "v0.0.0"; } | sed 's/^\(v[0-9]*\.[0-9]*\.[0-9]*\).*/\1/'
}

semver_prerelease_prerelease() {
  git rev-parse --abbrev-ref HEAD | sed 's/[^a-zA-Z0-9]/--/g'
}

semver_prerelease_build() {
  printf '%s-%s' "$(date -u +%Y%m%d-%H%M%S)" "$(git rev-parse --short HEAD)"
}

semver_prerelease() {
  core="$(semver_prerelease_core | awk -F. '{print $1"."$2"."$3+1}')"

  echo "$core-$(semver_prerelease_prerelease)+$(semver_prerelease_build)"
}

validate_semver_version() {
  _version=$1
  [ -z "${_version}" ] && { echo "VERSION is required" >&2 && exit 1; }

  echo "${_version}" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?([+][0-9A-Za-z.-]+)?$' || \
    { echo "Error: VERSION '${_version}' is not a valid semver string (expected: vMAJOR.MINOR.PATCH[-PRERELEASE][+BUILD])" && exit 1; }
}

validate_version_is_most_recent_in_changelog() {
  _version="$1"
  _most_recent=$(awk '/^## \[/{gsub(/^## \[/, ""); gsub(/\].*/, ""); print; exit}' CHANGELOG.md)

  [ "${_most_recent}" != "${_version}" ] && {
    echo "Error: VERSION '${_version}' is not the most recent version in CHANGELOG.md (most recent is '${_most_recent}')" >&2
    exit 1
  }
}

validate_changelog_contains_no_unreleased_sections() {
  grep -qE '^## \[Unreleased\]|^\[Unreleased\]:' CHANGELOG.md && {
    echo "Error: CHANGELOG.md contains an [Unreleased] section or link." >&2
    exit 1
  }
  return 0
}

prepare_changelog_for_release() {
  _version=$1

  validate_semver_version "${_version}"

  _date=$(date +%Y-%m-%d)

  sed -i \
    -e "s/## \[Unreleased\]/## [${_version}] - ${_date}/" \
    -e "s/^\[Unreleased\]: \(.*\/compare\/.*\)\.\.\..*/[${_version}]: \1...${_version}/" \
    CHANGELOG.md
}

extract_version_notes() {
  _version=$1

  validate_semver_version "${_version}"

  validate_changelog_contains_no_unreleased_sections

  validate_version_is_most_recent_in_changelog "${_version}"

  _notes=$(awk -v ver="${_version}" 'index($0,"## [" ver "]")==1{found=1; print; next} found && (/^## / || /^\[.*\]:/) {exit} found{print}' CHANGELOG.md)
  _link=$(awk -v ver="${_version}" 'index($0,"[" ver "]:")==1{print; exit}' CHANGELOG.md)

  [ -z "${_notes}" ] && {
    echo "Error: No release notes found for ${_version} in CHANGELOG.md" >&2
    exit 1
  }

  echo "${_notes}"
  if [ -n "${_link}" ]; then
    echo ""
    echo "${_link}"
  fi
}

prepare_prerelease_and_create_draft_release_on_side_tag() {
  _version="$(semver_prerelease)"
  prepare_changelog_for_release "${_version}" &&
  validate_changelog_contains_no_unreleased_sections &&
  validate_version_is_most_recent_in_changelog "${_version}"
    git add CHANGELOG.md &&
    git commit -m "release ${_version}" &&
    git tag "${_version}" &&
    git reset --hard HEAD~1 &&
    git push origin "${_version}"
}

case "$1" in
    last-release-tag)
      last_release_tag;;
    last-release-rev)
      last_release_rev;;
    semver-prerelease)
      semver_prerelease;;
    prepare-changelog-for-release)
      prepare_changelog_for_release "$2";;
    extract-version-notes)
      extract_version_notes "$2";;
    prepare-prerelease-and-create-draft-release-on-side-tag)
      prepare_prerelease_and_create_draft_release_on_side_tag;;
    '')
      echo "missing command" >&2
      exit 1;;
    *)
      echo "unrecognized command '$1'" >&2
      exit 1;;
esac