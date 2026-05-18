#!/bin/bash
set -ex
export KUBECONFIG=/github/workspace/$1
OSC_AKSK=$2
CSI=$3
PRELOADER_IMAGE=$4

export OSC_ACCESS_KEY=`echo $OSC_AKSK|cut -d% -f 1`
export OSC_SECRET_KEY=`echo $OSC_AKSK|cut -d% -f 2`
export OSC_REGION=`echo $OSC_AKSK|cut -d% -f 3`

if [ "$CSI" = "true" ]; then
  echo "Installing CSI"
  kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-8.3/client/config/crd/snapshot.storage.k8s.io_volumesnapshotclasses.yaml
  kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-8.3/client/config/crd/snapshot.storage.k8s.io_volumesnapshotcontents.yaml
  kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-8.3/client/config/crd/snapshot.storage.k8s.io_volumesnapshots.yaml
  kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-8.3/deploy/kubernetes/snapshot-controller/rbac-snapshot-controller.yaml
  kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-8.3/deploy/kubernetes/snapshot-controller/setup-snapshot-controller.yaml
  kubectl create secret generic osc-csi-bsu --from-literal=access_key=$OSC_ACCESS_KEY --from-literal=secret_key=$OSC_SECRET_KEY -n kube-system
  helm upgrade --install osc-bsu-csi-driver oci://docker.io/outscalehelm/osc-bsu-csi-driver \
      --namespace kube-system \
      --set driver.enableVolumeSnapshot=true \
      --set cloud.region=$OSC_REGION
fi

sed -i "s=PRELOADER_IMAGE=$PRELOADER_IMAGE=g" /snapshot.yaml
kubectl apply -f /snapshot.yaml
set +e
echo "Waiting for snapshot"
for i in {1..10}; do
  kubectl get vs -n image-preloader image-preloader-snap && break
  sleep 30
done
kubectl wait --for=jsonpath='{.status.readyToUse}'=true --timeout 10m -n image-preloader vs/image-preloader-snap
kubectl get vsc
kubectl logs -n image-preloader -l "batch.kubernetes.io/job-name=image-preloader"
set -e
VSC=`kubectl get vsc -n image-preloader -o custom-columns=NAME:.metadata.name,VS:.spec.volumeSnapshotRef.name --no-headers=true|grep image-preloader-snap|awk '{print $1}'`
HANDLE=`kubectl get -o template vsc/$VSC --template='{{.status.snapshotHandle}}'`
echo "SNAPSHOT_ID=$HANDLE" >> $GITHUB_OUTPUT