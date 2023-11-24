#!/bin/bash

_count=$1
_secret_template_path=$2

function generate_data {
  jq -nR --arg _count $_count '[range(0;($_count|tonumber)) | input]' < <(while true; do uuidgen ; done)
}

function create_secrets {
  cat /dev/stdin | jq -r --argjson t "$(<$_secret_template_path)" '.[] as $id | $t | .metadata.name=$id | .metadata.labels["kyma-project.io/runtime-id"]=$id'
}

datetime_postfix=$(date -u +%Y-%m-%dT%H:%M:%S) 
input_filename="/tmp/input_"$datetime_postfix".json"
echo $input_filename 
generate_data $1 $2 | tee $input_filename | create_secrets
