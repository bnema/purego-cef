#!/bin/sh
set -eu

status=0

if ! git diff --quiet HEAD --; then
  echo "generation changed tracked files:" >&2
  git status --short --untracked-files=no >&2
  status=1
fi

untracked_generated=$(git ls-files --others --exclude-standard -- \
  ':(glob)**/*_gen.go' \
  ':(glob)**/mocks/mock_*.go')
if [ -n "$untracked_generated" ]; then
  echo "generation created untracked generated files:" >&2
  printf '%s\n' "$untracked_generated" >&2
  status=1
fi

exit "$status"
