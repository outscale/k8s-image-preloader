#!/bin/sh
set -ex
./preloader version
./preloader export --force-pull --to $COPY_TO
./preloader snapshot --name $SNAPSHOT_NAME --class $SNAPSHOT_CLASS --namespace $SNAPSHOT_NAMESPACE --pvc $PVC_NAME