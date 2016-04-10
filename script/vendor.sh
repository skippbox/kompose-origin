#!/usr/bin/env bash
set -e

cd "$(dirname "$BASH_SOURCE")/.."
rm -rf vendor/
source 'script/.vendor-helpers.sh'

clone git github.com/Sirupsen/logrus v0.9.0
clone git github.com/codegangsta/cli 839f07bfe4819fa1434fa907d0804ce6ec45a5df
clone git github.com/docker/distribution 467fc068d88aa6610691b7f1a677271a3fac4aac
clone git github.com/vbatts/tar-split v0.9.11
clone git github.com/docker/docker 9837ec4da53f15f9120d53a6e1517491ba8b0261
clone git github.com/docker/go-units 651fc226e7441360384da338d0fd37f2440ffbe3
clone git github.com/docker/go-connections v0.2.0
clone git github.com/docker/engine-api fd7f99d354831e7e809386087e7ec3129fdb1520
clone git github.com/vdemeester/docker-events b308d2e8d639d928c882913bcb4f85b3a84c7a07
clone git github.com/flynn/go-shlex 3f9db97f856818214da2e1057f8ad84803971cff
clone git github.com/gorilla/context 14f550f51a
clone git github.com/gorilla/mux e444e69cbd
clone git github.com/opencontainers/runc 2441732d6fcc0fb0a542671a4372e0c7bc99c19e
clone git github.com/fsouza/go-dockerclient 39d9fefa6a7fd4ef5a4a02c5f566cb83b73c7293
clone git github.com/stretchr/testify a1f97990ddc16022ec7610326dd9bce31332c116
clone git github.com/davecgh/go-spew 5215b55f46b2b919f50a1df0eaa5886afe4e3b3d
clone git github.com/pmezard/go-difflib d8ed2627bdf02c080bf22230dbb337003b7aba2d
clone git golang.org/x/crypto 4d48e5fa3d62b5e6e71260571bf76c767198ca02 https://github.com/golang/crypto.git
clone git golang.org/x/net 47990a1ba55743e6ef1affd3a14e5bac8553615d https://github.com/golang/net.git
clone git gopkg.in/check.v1 11d3bc7aa68e238947792f30573146a3231fc0f1

clone git github.com/cloudfoundry-incubator/candiedyaml 5cef21e2e4f0fd147973b558d4db7395176bcd95
clone git github.com/Azure/go-ansiterm 388960b655244e76e24c75f48631564eaefade62
clone git github.com/Microsoft/go-winio v0.3.0
clone git github.com/xeipuuv/gojsonpointer e0fe6f68307607d540ed8eac07a342c33fa1b54a
clone git github.com/xeipuuv/gojsonreference e02fc20de94c78484cd5ffb007f8af96be030a45
clone git github.com/xeipuuv/gojsonschema ac452913faa25c08bb78810d3e6f88b8a39f8f25
clone git github.com/kr/pty 5cf931ef8f

# These are moved from Dockerfile
clone git golang.org/x/tools c887be1b2ebd11663d4bf2fbca508c449172339e https://github.com/golang/tools.git
clone git github.com/golang/lint 8f348af5e29faa4262efdc14302797f23774e477
clone git github.com/aktau/github-release v0.6.2

# These are required for the kubernetes dependencies and will need to be kept in sync with
# github.com/kubernetes/kubernetes/Godeps/Godeps.json
clone git bitbucket.org/bertimus9/systemstat 1468fd0db20598383c9393cccaa547de6ad99e5e
clone hg bitbucket.org/ww/goautoneg 75cd24fc2f2c2a2088577d12123ddee5f54e0675
clone git github.com/beorn7/perks b965b613227fddccbfffe13eae360ed3fa822f8d
clone git github.com/blang/semver 31b736133b98f26d5e078ec9eb591666edfd091f
clone git github.com/google/cadvisor 546a3771589bdb356777c646c6eca24914fdd48b
clone git github.com/emicklei/go-restful 777bb3f19bcafe2575ffb2a3e46af92509ae9594
clone git github.com/ghodss/yaml 73d445a93680fa1a78ae23a5839bad48f32ba1ee
clone git github.com/golang/glog 44145f04b68cf362d9c4df2182967c2275eaefed
clone git github.com/golang/protobuf b982704f8bb716bb608144408cff30e15fbde841
clone git github.com/matttproud/golang_protobuf_extensions fc2b8d3a73c4867e51861bbdd5ae3c1f0869dd6a
clone git github.com/pborman/uuid ca53cad383cad2479bbba7f7a1a05797ec1386e4
clone git github.com/juju/ratelimit 77ed1c8a01217656d2080ad51981f6e99adaa177
clone git github.com/google/gofuzz bbcb9da2d746f8bdbd6a936686a0a6067ada0ec5
clone git github.com/prometheus/client_golang 3b78d7a77f51ccbc364d4bc170920153022cfd08
clone git github.com/prometheus/client_model fa8ad6fec33561be4280a8f0514318c79d7f6cb6
clone git github.com/prometheus/common ef7a9a5fb138aa5d3a19988537606226869a0390
clone git github.com/prometheus/procfs 490cc6eb5fa45bf8a8b7b73c8bc82a8160e8531d
clone git github.com/spf13/pflag 08b1a584251b5b62f458943640fc8ebd4d50aaa5
clone git speter.net/go/exp/math/dec/inf 42ca6cd68aa922bc3f32f1e056e61b65945d9ad7 https://code.google.com/p/go-decimal-inf.exp
clone git github.com/ugorji/go f4485b318aadd133842532f841dc205a8e339d74
clone git gopkg.in/yaml.v2 d466437aa4adc35830964cffc5b5f262c63ddcb4
clone git k8s.io/kubernetes v1.2.0 https://github.com/kubernetes/kubernetes.git

clean && mv vendor/src/* vendor
