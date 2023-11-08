#!/bin/sh

get_git_dir() {
    git rev-parse --show-toplevel
}

get_tag() {
    git describe --exact-match 2>/dev/null
}

get_pseudo_version() {
    ref=${1:-"HEAD"}
    prefix=${2:-"v0.0.0"}

    latest_tag=$(git describe --tags --abbrev=0)
    if [ -n "${latest_tag}" ]; then
        prefix=${latest_tag}-$(git rev-list --count ${latest_tag}..${ref})
    fi

    # UTC time the revision was created (yyyymmddhhmmss).
    timestamp=$(TZ=UTC git show --no-patch --format='%cd' --date='format-local:%Y%m%d%H%M%S' ${ref})

    # 12-character prefix of the commit hash
    revision=$(git rev-parse --short=12 --verify ${ref})

    echo "${prefix}-${timestamp}-${revision}"
}

get_version() {
    version=$(get_pseudo_version)
    
    tag=$(git describe --exact-match 2>/dev/null || true)
    if [ -n "${tag}" ]; then
        version=${tag}
    fi

    echo ${version}
}

is_dirty() {
    ! git diff --quiet
}

is_tagged() {
    test -n "$(get_tag)"
}
