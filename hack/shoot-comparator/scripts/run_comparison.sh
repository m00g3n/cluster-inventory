#!/bin/bash

echo "Cleaning up the resources created for shoot comparison"

printf "Removing resources needed for fetching results"
kubectl delete -n kcp-system pvc shoot-comparator-pvc-read-only --ignore-not-found
kubectl delete -n kcp-system volumesnapshot shoot-comparator-pvc --ignore-not-found
kubectl delete -n kcp-system po fetch-test-comparison-results --ignore-not-found

printf "Removing resources needed for performing comparison \n"
kubectl delete -n kcp-system job/compare-shoots --ignore-not-found
kubectl delete -n kcp-system pvc/shoot-comparator-pvc --ignore-not-found

kubectl delete -n kcp-system pvc test-prov-shoot-read-only test-kim-shoot-read-only --ignore-not-found
kubectl delete -n kcp-system volumesnapshot test-kim-shoot-spec-storage test-prov-shoot-spec-storage --ignore-not-found

printf "Preparing data for comparison \n"
kubectl apply -f ./manifests/snapshot-for-comparison.yaml

printf "Running comparison job \n"
kubectl apply -f ./manifests/job.yaml

# wait for completion
echo "Waiting for the job to complete. It may take couple of minutes. Please, be patient!"
kubectl wait --for=condition=complete job/compare-shoots -n kcp-system --timeout=10m

# verify if the job succeeded
if kubectl wait --for=condition=failed --timeout=0 job/compare-shoots -n kcp-system 2>/dev/null; then
    echo "Job failed to complete. Exiting..."
    exit 1
fi

echo "Job is complete."
echo "Fetching results"
