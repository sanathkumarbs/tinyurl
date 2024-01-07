export DOCKER_BUILDKIT=1
COMMIT_SHA=$(git rev-parse --short HEAD)
ROOTDIR=$(git rev-parse --show-toplevel)

docker build \
        --no-cache \
        --pull \
        --progress=plain \
        --tag "tiny:${COMMIT_SHA}" \
		--tag "tiny:latest" \
        -f "packaging/services/tiny/Dockerfile" \
        "${ROOTDIR}"