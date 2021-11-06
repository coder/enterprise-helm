#!/usr/bin/env bash
#
# Run Go-based unit tests

set -o errexit
set -o nounset
set -o pipefail

CI=${CI:-""}
PROJECT_ROOT=$(git rev-parse --show-toplevel)
# shellcheck source=lib.sh
source "$PROJECT_ROOT/scripts/lib.sh"

echo "--- Running go test"
export FORCE_COLOR=true

test_args=(
  -v
  -failfast
  "${TEST_ARGS:-}"
)

REPORTDIR="/tmp/testreports"
mkdir -p "$REPORTDIR"
TESTREPORT_JSON="$REPORTDIR/test_go.json"
TESTREPORT_XML="$REPORTDIR/test_go.xml"
COVERAGE="$REPORTDIR/test_go.coverage"

test_args+=(
  "-covermode=set"
  "-coverprofile=$COVERAGE"
)

declare test_status=0

pushd "$PROJECT_ROOT/tests" >/dev/null 2>&1
  # Allow failures to ensure that we can report on slow tests
  set +o errexit

  run_trace false gotestsum \
    --debug \
    --jsonfile="$TESTREPORT_JSON" \
    --junitfile="$TESTREPORT_XML" \
    --hide-summary=skipped \
    --packages="./..." \
    -- "${test_args[@]}"
  test_status=$?

  # Re-enable failures if steps fail
  set -o errexit
popd >/dev/null 2>&1

# These unit tests should all be fast, so this report should be empty
threshold="5s"
echo "--- ಠ_ಠ The following tests took longer than $threshold to complete:"
run_trace false gotestsum tool slowest \
  --jsonfile="$TESTREPORT_JSON" \
  --threshold="$threshold"

exit "$test_status"
