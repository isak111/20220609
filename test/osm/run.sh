#!/usr/bin/env bash

set -o errexit

REPO_ROOT=$(git rev-parse --show-toplevel)
DIR="$(cd "$(dirname "$0")" && pwd)"

"$DIR"/install.sh

${REPO_ROOT}/test/workloads/init.sh
"$DIR"/init-pre-test.sh
"$DIR"/test-canary.sh

${REPO_ROOT}/test/workloads/init.sh
"$DIR"/init-pre-test.sh
"$DIR"/test-steps.sh
