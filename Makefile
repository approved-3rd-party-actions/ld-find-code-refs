# Note: These commands pertain to the development of ld-find-code-refs.
#       They are not intended for use by the end-users of this program.
SHELL=/bin/bash
GORELEASER_VERSION=v0.169.0

build:
	go build ./cmd/...

init:
	pre-commit install

test: lint
	go test ./...

lint:
	pre-commit run -a --verbose golangci-lint

# Strip debug informatino from production builds
BUILD_FLAGS = -ldflags="-s -w"

compile-macos-binary:
	GOOS=darwin GOARCH=amd64 go build ${BUILD_FLAGS} -o out/ld-find-code-refs ./cmd/ld-find-code-refs

compile-windows-binary:
	GOOS=windows GOARCH=amd64 go build ${BUILD_FLAGS} -o out/ld-find-code-refs.exe ./cmd/ld-find-code-refs

compile-linux-binary:
	GOOS=linux GOARCH=amd64 go build ${BUILD_FLAGS} -o build/package/cmd/ld-find-code-refs ./cmd/ld-find-code-refs

compile-github-actions-binary:
	GOOS=linux GOARCH=amd64 go build ${BUILD_FLAGS} -o build/package/github-actions/ld-find-code-refs-github-action ./build/package/github-actions

compile-bitbucket-pipelines-binary:
	GOOS=linux GOARCH=amd64 go build ${BUILD_FLAGS} -o build/package/bitbucket-pipelines/ld-find-code-refs-bitbucket-pipeline ./build/package/bitbucket-pipelines

# Get the lines added to the most recent changelog update (minus the first 2 lines)
RELEASE_NOTES=<(GIT_EXTERNAL_DIFF='bash -c "diff --unchanged-line-format=\"\" $$2 $$5" || true' git log --ext-diff -1 --pretty= -p CHANGELOG.md)

echo-release-notes:
	@cat $(RELEASE_NOTES)

validate-circle-orb:
	test $(TAG) || (echo "Please provide tag"; exit 1)
	circleci orb validate build/package/circleci/orb.yml || (echo "Unable to validate orb"; exit 1)

publish-dev-circle-orb: validate-circle-orb
	circleci orb publish build/package/circleci/orb.yml launchdarkly/ld-find-code-refs@dev:$(TAG)

publish-release-circle-orb: validate-circle-orb
	circleci orb publish build/package/circleci/orb.yml launchdarkly/ld-find-code-refs@$(TAG)

clean:
	rm -rf out/
	rm -f build/pacakge/cmd/ld-find-code-refs
	rm -f build/package/github-actions/ld-find-code-refs-github-action
	rm -f build/package/bitbucket-pipelines/ld-find-code-refs-bitbucket-pipeline

# We use goreleaser to package and publish:
# 1. docker image for the cli
# 2. github action docker image
# 3. bitbucket pipelines docker image
# 4. circleci orb
# The first three are published to dockerhub. The last one is published to circleci orb registry.
PUBLISH_CMD=curl -sL https://git.io/goreleaser | GOPATH=$(mktemp -d) VERSION=$(GORELEASER_VERSION) bash -s -- --rm-dist --release-notes $(RELEASE_NOTES)

publish:
	$(PUBLISH_CMD)

products-for-release:
	$(PUBLISH_CMD) --skip-publish --skip-validate

.PHONY: init test lint compile-github-actions-binary compile-macos-binary compile-linux-binary compile-windows-binary compile-bitbucket-pipelines-binary echo-release-notes publish-dev-circle-orb publish-release-circle-orb publish-all clean build
