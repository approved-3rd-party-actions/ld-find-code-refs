# Releasing

## Versioning
This project adheres to [Semantic Versioning](http://semver.org). Release version tags should be in the form `MAJOR.MINOR.PATCH`, with no leading v. When releasing, be sure to update the version number in [`version.go`](https://github.com/launchdarkly/ld-find-code-refs/blob/master/internal/version/version.go), and in the [CircleCI orb](https://github.com/launchdarkly/ld-find-code-refs/blob/master/build/package/circleci/orb.yml).

## Github Releases

This project uses [goreleaser](https://goreleaser.com/) to generate github releases. Releases are automated via CircleCI. To generate a new release, simply tag the commit you want to release and push the tag. If the tag ends in -rc(.+), the github release will be marked as "Pre-release." If you'd like to see how release notes are generated, see the .circleci/config.yml publish job.

Make sure you update the changelog before generating a release.

Once a release has been completed, update the [BitBucket pipelines](https://bitbucket.org/launchdarkly/ld-find-code-refs-pipe) repo with the new version number, and push a tag containing the version number along with your commit. Example release commit: https://bitbucket.org/launchdarkly/ld-find-code-refs-pipe/commits/0b1e920c7322cd495f4fc1a09d339342d32606e4

## Docker Hub

We use goreleaser to package and publish four images:

1. docker image for the cli
2. github action docker image
3. bitbucket pipelines docker image
4. circleci orb

The first three are published to dockerhub. The last one is published to circleci orb registry.

## Beta builds

To push a beta build, set the `PRERELEASE=true` environment variable before running a release task. e.g. `make publish-all TAG=1.0.0-beta1`. Note: to publish a beta circle ci orb, run `make publish-dev-circle-orb TAG=$VERSION`