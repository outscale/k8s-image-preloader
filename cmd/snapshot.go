package cmd

import (
	"context"
	"fmt"

	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v8/apis/volumesnapshot/v1"
	snapshotclient "github.com/kubernetes-csi/external-snapshotter/client/v8/clientset/versioned"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	pvcName           string
	snapshotClass     string
	snapshotName      string
	snapshotNamespace string
)

func init() {
	rootCmd.AddCommand(snapshotCmd)
	snapshotCmd.Flags().StringVar(&pvcName, "pvc", "", "PVC to snashot (required)")
	snapshotCmd.Flags().StringVar(&snapshotClass, "class", "", "VolumeSnapshotClass to use (required)")
	snapshotCmd.Flags().StringVar(&snapshotName, "name", "", "Name of the VolumeSnapshot to create (required)")
	snapshotCmd.Flags().StringVar(&snapshotNamespace, "namespace", "", "Namespace of the VolumeSnapshot to create (required)")
}

var snapshotCmd = &cobra.Command{
	Use:     "snapshot",
	Aliases: []string{"s"},
	Short:   "Create a snapshot from the PVC storing exports of images",
	Long:    `Create a VolumeSnapshot from the PVC storing exports of images.`,
	Run: func(cmd *cobra.Command, args []string) {
		k8s := getSnapshotClient()
		ctx := context.Background()
		snapContentName, err := Snapshot(ctx, k8s)
		exitOnError(err)
		snapContent, err := k8s.SnapshotV1().VolumeSnapshotContents().Get(ctx, snapContentName, metav1.GetOptions{})
		exitOnError(err)
		fmt.Println("Snapshot created:", snapContent.Status.SnapshotHandle)
	},
}

func Snapshot(ctx context.Context, k8s *snapshotclient.Clientset) (string, error) {
	snap := &snapshotv1.VolumeSnapshot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      snapshotName,
			Namespace: snapshotNamespace,
		},
		Spec: snapshotv1.VolumeSnapshotSpec{
			Source: snapshotv1.VolumeSnapshotSource{
				PersistentVolumeClaimName: &pvcName,
			},
			VolumeSnapshotClassName: &snapshotClass,
		},
	}
	snap, err := k8s.SnapshotV1().VolumeSnapshots(snapshotNamespace).Create(ctx, snap, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	for snap.Status == nil || snap.Status.VolumeGroupSnapshotName == nil {
		snap, err = k8s.SnapshotV1().VolumeSnapshots(snapshotNamespace).Get(ctx, snapshotName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
	}
	return *snap.Status.VolumeGroupSnapshotName, nil
}
