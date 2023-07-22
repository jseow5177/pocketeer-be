#!/usr/bin/env bash

build_service() {
    echo "Building Pocketeer BE service..."
    go build -v -o ./bin/service ./cmd/service
    echo "Done. Here is target information"
    ls -lah -d ./bin/service
}

build_job() {
    echo "Building Pocketeer job"
    go build -v -o ./bin/job ./cmd/job
    echo "Done. Here is target information"
    ls -lah -d ./bin/job
}

usage() {
    cat <<EOF
Build artifacts for Pocketeer BE. Known recipes:
    service     build ./cmd/service into ./bin/service
    job         build ./cmd/job into ./bin/job
EOF
}

pushd "${PROJECT_DIR}" || exit 1

case "$1" in
service)
    build_service
    exit
    ;;
job)
    build_job
    exit
    ;;
mq)
    build_mq
    exit
    ;;
*)
    usage
    exit
    ;;
esac

popd
