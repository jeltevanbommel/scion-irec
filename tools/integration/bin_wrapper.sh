#!/bin/bash

# Wrapper to run integration binaries.

PROG=$1
shift

log() {
    echo "$(date -u +"%F %T.%6N%z") $@" 1>&2
}

set -o pipefail

[ -n "$IA" ] && echo "Listening ia=$IA"

set -o noglob
log "bin_wrapper: Starting $PROG $@"

"$PROG" "$@" |& while read line; do log $line; done
exit_status=$?
set +o noglob

log "bin_wrapper: Stopped, exit code: $exit_status"

exit $exit_status
