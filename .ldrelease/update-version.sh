#!/bin/bash

VERSION_GO=internal/version/version.go
sed -i "" "s/const Version =.*/const Version = \"${LD_RELEASE_VERSION}\"/g" ${VERSION_GO}

VERSION_ORB=build/package/circleci/orb.yml
sed -i "" "s#launchdarkly: launchdarkly/ld-find-code-refs@.*#launchdarkly: launchdarkly/ld-find-code-refs@${LD_RELEASE_VERSION}#g" ${VERSION_ORB}
sed -i "" "s#- image: launchdarkly/ld-find-code-refs:.*#- image: launchdarkly/ld-find-code-refs:${LD_RELEASE_VERSION}#g" ${VERSION_ORB}