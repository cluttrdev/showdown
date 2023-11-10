#!/usr/bin/env bash

set -euo pipefail

usage() { echo "usage: $(basename -- $0) [-h] [-v] [-o NAME]" 1>&2; }

output=""
verbose=false
while getopts ":o:vh" opt; do
    case "${opt}" in
        o)
            output=${OPTARG}
            ;;
        v)
            verbose=true
            ;;
        h)
            usage
            exit 0
            ;;
        *)
            usage
            exit 1
            ;;
    esac
done

# load common helpers
scriptsdir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
source "${scriptsdir}/functions.sh"

version=$(get_version)
commit_sha=$(get_commit_sha)
commit_date=$(get_commit_date)

if [ -z "${output}" ]; then
    git_dir=$(get_git_dir)
    bin_dir="${git_dir}/bin"
    bin_name=$(basename ${git_dir})
    output="${bin_dir}/${bin_name}"
fi

goos=${GOOS:-$(go env GOHOSTOS)}
goarch=${GOARCH:-$(go env GOHOSTARCH)}

${verbose} && echo "Building for ${goos}/${goarch}..."
GOOS=${goos} GOARCH=${goarch} go build \
    -o "${output}" \
    -ldflags "-X 'main.version=${version}' -X 'main.commit=${commit_sha}' -X 'main.date=${commit_date}'"
${verbose} && echo "Building for ${goos}/${goarch}... done"

exit 0
