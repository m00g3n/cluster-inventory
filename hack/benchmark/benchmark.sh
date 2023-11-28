#!/bin/bash

_count=$1
_secret_template_path=$2

function generate_data {
  jq -nR --arg _count $_count '[range(0;($_count|tonumber)) | input]' < <(while true; do uuidgen | awk '{print tolower($0)}' ; done)
}

function create_secrets {
  cat /dev/stdin | jq -r --argjson t "$(<$_secret_template_path)" '.[] as $id | $t | .metadata.name="kubeconfig-"+$id | .metadata.labels["kyma-project.io/runtime-id"]=$id | .metadata.labels["kyma-project.io/shoot-name"]=$id'
}

datetime_postfix=$(date -u +%Y-%m-%dT%H:%M:%S) 
input_filename="/tmp/input_"$datetime_postfix".json"
generate_data $1 $2 | tee $input_filename | create_secrets
