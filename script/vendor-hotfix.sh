#!/usr/bin/env bash
# hotfix speter.net/go/exp/math/dec/inf
mkdir -p vendor/speter.net/go/exp/math/dec/inf
cp -r vendor/github.com/go-inf/inf/* vendor/speter.net/go/exp/math/dec/inf/
