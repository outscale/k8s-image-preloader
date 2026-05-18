#!/bin/bash
set -ex
export KUBECONFIG=/github/workspace/$1

kubectl delete -f /snapshot.yaml || /bin/true