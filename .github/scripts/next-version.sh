#!/usr/bin/env bash

set -euo pipefail

branch_name="${1:-}"
if [[ -z "$branch_name" ]]; then
	echo "branch name is required" >&2
	exit 1
fi

latest_tag="v0.0.0"
while IFS= read -r tag; do
	if [[ "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
		latest_tag="$tag"
		break
	fi
done < <(git tag --list "v*" --sort=-version:refname)

if [[ ! "$latest_tag" =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
	echo "invalid semantic version tag: $latest_tag" >&2
	exit 1
fi

major="${BASH_REMATCH[1]}"
minor="${BASH_REMATCH[2]}"
patch="${BASH_REMATCH[3]}"

case "$branch_name" in
patch/*)
	((patch += 1))
	;;
minor/*)
	((minor += 1))
	patch=0
	;;
major/*)
	((major += 1))
	minor=0
	patch=0
	;;
*)
	echo "unsupported branch prefix: $branch_name; expected patch/, minor/, or major/" >&2
	exit 1
	;;
esac

printf 'v%d.%d.%d\n' "$major" "$minor" "$patch"
