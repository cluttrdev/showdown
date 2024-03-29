#!/usr/bin/env bash

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

version=$(get_version)

declare -A OSARCHMAP=(
    [linux]="amd64,arm,arm64"
    [darwin]="amd64,arm64"
    [windows]="amd64,arm,arm64"
)

${verbose} && echo "Building binaries..."
for os in ${!OSARCHMAP[@]}; do
    for arch in ${OSARCHMAP[$os]//,/ }; do
        tmp_dir=${dist_dir}/${bin_name}_${version}_${os}_${arch}

        out="${tmp_dir}/${bin_name}"

        ${verbose} && echo "  for ${os}/${arch}"
        GOOS=${os} GOARCH=${arch} ${scriptsdir}/build.sh -o ${out}
    done
done
${verbose} && echo "Building binaries... done"

${verbose} && echo "Creating package archives..."
for os in ${!OSARCHMAP[@]}; do
    dirs=$(find ${dist_dir}/ -mindepth 1 -maxdepth 1 -type d -name ${bin_name}_${version}_${os}_*)

    case "${os}" in
        linux | darwin)
            for dir in ${dirs}; do
                find $dir -printf "%P\n" \
                | tar -czf ${dir}.tar.gz --no-recursion -C ${dir} -T -

                rm -r ${dir}
            done
            ;;
        windows)
            for dir in ${dirs}; do
                (cd ${dir} && zip -q -r - .) > ${dir}.zip

                rm -r ${dir}
            done
            ;;
    esac
done
${verbose} && echo "Creating package archives... done"

exit 0
