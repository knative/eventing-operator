#!/usr/bin/env bash

# Copyright 2019 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script runs the end-to-end tests against Knative Eventing
# Operator built from source.  It is started by prow for each PR. For
# convenience, it can also be executed manually.

# If you already have a Knative cluster setup and kubectl pointing
# to it, call this script with the --run-tests arguments and it will use
# the cluster and run the tests.

# Calling this script without arguments will create a new cluster in
# project $PROJECT_ID, start knative in it, run the tests and delete the
# cluster.

source $(dirname $0)/../vendor/knative.dev/test-infra/scripts/e2e-tests.sh

OPERATOR_DIR=$(dirname $0)/..
KNATIVE_EVENTING_DIR=${OPERATOR_DIR}/..

# Namespace used for tests
readonly TEST_NAMESPACE="operator-tests"

declare -A COMPONENTS
COMPONENTS=(
  ["eventing.yaml"]="config"
  ["in-memory-channel.yaml"]="config/channels/in-memory-channel"
)

function install_eventing_operator() {
  header "Installing Knative Eventing operator"

  # Deploy the operator
  ko apply -f config/
  wait_until_pods_running default || fail_test "Eventing Operator did not come up"
}

function knative_setup() {
  echo ">> Creating test namespaces"
  kubectl create namespace $TEST_NAMESPACE
  generate_latest_eventing_manifest
  install_eventing_operator
}

function generate_latest_eventing_manifest() {
  # Go the directory to download the source code of knative eventing
  cd ${KNATIVE_EVENTING_DIR}

  # Download the source code of knative eventing
  git clone https://github.com/knative/eventing.git
  cd eventing
  COMMIT_ID=$(git rev-parse --verify HEAD)
  echo ">> The latest commit ID of Knative Eventing is ${COMMIT_ID}."
  mkdir -p output

  # Generate the manifest
  LABEL_YAML_CMD=(cat)
  local all_yamls=()
  for yaml in "${!COMPONENTS[@]}"; do
    local config="${COMPONENTS[${yaml}]}"
    echo "Building Knative Eventing - ${config}"
    ko resolve -P -f ${config}/ | "${LABEL_YAML_CMD[@]}" > ${yaml}
  done

  EVENTING_YAML=${KNATIVE_EVENTING_DIR}/eventing/eventing.yaml
  IMC_YAML=${KNATIVE_EVENTING_DIR}/eventing/in-memory-channel.yaml
  if [[ -f "${EVENTING_YAML}" && -f "${IMC_YAML}" ]]; then
    echo ">> Replacing the current manifest in operator with the generated manifest"
    rm -rf ${OPERATOR_DIR}/cmd/manager/kodata/knative-eventing/*
    cp ${EVENTING_YAML} ${OPERATOR_DIR}/cmd/manager/kodata/knative-eventing/eventing.yaml
    cp ${IMC_YAML} ${OPERATOR_DIR}/cmd/manager/kodata/knative-eventing/eventing-imc.yaml
  else
    echo ">> The eventing.yaml and in-memory-channel.yaml were not generated, so keep the current manifest"
  fi

  # Go back to the directory of operator
  cd ${OPERATOR_DIR}
}

# Script entry point.

# Skip installing istio as an add-on
initialize $@ --skip-istio-addon

# If we got this far, the operator installed Knative Eventing
header "Running tests for Knative Eventing Operator"
failed=0

# Run the integration tests
go_test_e2e -timeout=20m ./test/e2e || failed=1

# Require that tests succeeded.
(( failed )) && fail_test

success
