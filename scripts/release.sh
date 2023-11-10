#!/usr/bin/env bash

GITHUB_OWNER=${GITHUB_OWNER:-"cluttrdev"}
GITHUB_REPO=${GITHUB_REPO:-"showdown"}
GITHUB_TOKEN=${GITHUB_TOKEN}

set -euo pipefail

usage() { echo "usage: $(basename -- $0) [-h] [-v] [-n NAME] [-d DIR]" 1>&2; }

dist_dir="dist"
bin_name=$(basename $(git rev-parse --show-toplevel))
verbose=false
while getopts ":n:d:vh" opt; do
    case "${opt}" in
        n)
            bin_name=${OPTARG}
            ;;
        d)
            dist_dir=${OPTARG}
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

#
# Sanity checks
#

is_dirty && (echo "Working directory is dirty" && exit 1)

tag=$(get_tag || true)
test -n "${tag}" || (echo "No tag exactly matches current commit" && exit 1)

test -n "${GITHUB_TOKEN}" || (echo "Missing required github token" && exit 1)

#
# Build release assets
#

${verbose} && echo "Creating distributions archives..."
sh -c "${scriptsdir}/dist.sh -n '${bin_name}' -d '${dist_dir}'"

assets=$(find ${dist_dir}/ -type f \
    -name "${bin_name}_${tag}_*.tar.gz" -o -name "${bin_name}_${tag}_*.zip"
)

#
# Create release changelog
#

changes=$(get_changes)
changelog=$(cat <<-EOF
## Changelog
$(awk '{ print "  - [`" $1 "`][" $1 "] " substr($0, index($0, $2)) }' <<< $changes)

<!-- Link -->
$(awk '{ print "[" $1 "]: https://github.com/cluttrdev/showdown/commit/" $1 }' <<< $changes)
EOF
)

#
# Create release
#

owner=${GITHUB_OWNER}
repo=${GITHUB_REPO}

${verbose} && echo "Creating release for ${tag}..."
response=$(curl -L -s \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer ${GITHUB_TOKEN}" \
    -H "X-GitHub-Api-Version 2022-11-28" \
    -d "{\"tag_name\": \"${tag}\", \"name\": \"${tag}\", \"body\": \"${changelog}\"}" \
    https://api.github.com/repos/${owner}/${repo}/releases \
    2>/dev/null
)

error=$(jq -r '.errors[0].code // empty' <<< $response)
[ -n "$error" ] && {
    echo ${error}
    exit 1
}

release_id=$(jq -r '.id // empty' <<< $response)
[ -n "$release_id" ] || {
    echo "No release with tag: ${tag}"
    exit 1
}

#
# Upload release assets
#

upload_url="https://uploads.github.com/repos/${owner}/${repo}/releases/${release_id}/assets"
for asset in ${assets}; do
    name=$(basename ${asset})
    ${verbose} && echo "uploading asset: ${name}"
    response=$(curl -L -s \
        -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer ${GITHUB_TOKEN}" \
        -H "X-GitHub-Api-Version 2022-11-28" \
        -H "Content-Type: application/octet-stream" \
        ${upload_url}?name=${name} \
        --data-binary "@${asset}"
    )
done
