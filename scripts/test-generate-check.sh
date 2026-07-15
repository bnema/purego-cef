#!/bin/sh
set -eu

project_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT HUP INT TERM

cd "$tmp"
git init -q
git config user.email test@example.com
git config user.name test
mkdir -p cef/mocks ignored
printf 'package cef\n' > cef/example_gen.go
printf 'ignored/\n' > .gitignore
git add .gitignore cef/example_gen.go
git commit -qm baseline

# Ignored files present before the check must not look like generation drift.
printf 'ambient cache\n' > ignored/cache_gen.go
"$project_root/scripts/check-generation-drift.sh"

printf '// drift\n' >> cef/example_gen.go
if "$project_root/scripts/check-generation-drift.sh" >/dev/null 2>&1; then
  echo 'expected tracked generation drift to fail the check' >&2
  exit 1
fi
git checkout -q -- cef/example_gen.go

printf 'package mocks\n' > cef/mocks/mock_fixture.go
if "$project_root/scripts/check-generation-drift.sh" >/dev/null 2>&1; then
  echo 'expected untracked generated file to fail the check' >&2
  exit 1
fi

echo 'generate-check drift tests passed'
