package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"slices"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Image struct {
	ref     string
	inCache bool
}

var (
	stdin            bool
	cacheOnly        bool
	snapshotsPath    string
	excludedPrefixes []string
	forcePull        bool
	noRestoreScript  bool
)

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().BoolVar(&stdin, "stdin", false, "Fetch image list from stdin instead of the local cluster")
	exportCmd.Flags().BoolVar(&cacheOnly, "cache-only", false, "Only list images present in local cache, do not list pods")
	exportCmd.Flags().StringSliceVar(&excludedPrefixes, "exclude", []string{"docker.io/outscale/k8s-image-preloader"}, "Prefixes to skip from the list")
	exportCmd.Flags().BoolVar(&forcePull, "force-pull", false, "Force an image pull before exporting")
	exportCmd.Flags().StringVar(&snapshotsPath, "to", "/snapshot", "Path to snapshot volume")
	exportCmd.Flags().BoolVar(&noRestoreScript, "no-restore-script", false, "Do not copy restore script to the volume")
}

var exportCmd = &cobra.Command{
	Use:     "export",
	Aliases: []string{"e"},
	Short:   "Export a list of images to a path",
	Long: `Export a list of images to a path

By default, fetches the list of all images from the local cluster,
reading the local containerd cache and the list of pods/cronjobs,
and exports all images found to path.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		var imgs []Image
		if stdin {
			imgs = ListFromStdin()
		} else {
			k8s := getK8SClient()
			imgs = ListFromCluster(ctx, k8s, excludedPrefixes)
		}
		Export(ctx, imgs)
		Sync(ctx)
		if !noRestoreScript {
			CopyRestoreScript(ctx)
		}
	},
}

func ListFromCluster(ctx context.Context, k8s *kubernetes.Clientset, excludedPrefixes []string) []Image {
	imgs := listFromCache(ctx, excludedPrefixes)
	if !cacheOnly {
		add := listFromPods(ctx, k8s, excludedPrefixes)
		for _, img := range add {
			if !slices.ContainsFunc(imgs, func(i Image) bool { return i.ref == img.ref }) {
				fmt.Println("Found " + img.ref + " in pod spec")
				imgs = append(imgs, img)
			}
		}
	}
	return imgs
}

func listFromCache(ctx context.Context, excludedPrefixes []string) []Image {
	local := ctr(ctx, "images", "list")
	imgs := make([]Image, 0, len(local))
	done := map[string]bool{}
	for i, line := range local {
		if i == 0 { // skipping header
			continue
		}
		parts := strings.Fields(line)
		ref, digest := parts[0], parts[2]
		if done[digest] {
			fmt.Println("Skipping " + ref)
			continue
		}
		done[digest] = true
		if slices.ContainsFunc(excludedPrefixes, func(prefix string) bool {
			return strings.HasPrefix(ref, prefix)
		}) {
			fmt.Println("Excluding " + ref)
			continue
		}
		fmt.Println("Found " + ref + " in local cache")
		imgs = append(imgs, Image{ref: ref, inCache: true})
	}
	return imgs
}

func listFromPods(ctx context.Context, k8s *kubernetes.Clientset, excludedPrefixes []string) []Image {
	var lst []Image
	nss, err := k8s.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	exitOnError(err)
	for _, ns := range nss.Items {
		fmt.Println("Listing pods in " + ns.Name + "...")
		pods, err := k8s.CoreV1().Pods(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			exitOnError(err)
		}
		for _, pod := range pods.Items {
			imgs := append(listImages(ctx, pod.Spec.Containers, excludedPrefixes),
				listImages(ctx, pod.Spec.InitContainers, excludedPrefixes)...)
			for _, img := range imgs {
				if !slices.ContainsFunc(lst, func(i Image) bool { return i.ref == img }) {
					lst = append(lst, Image{ref: img})
				}
			}
		}
		fmt.Println("Listing cron jobs in " + ns.Name + "...")
		crons, err := k8s.BatchV1().CronJobs(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			exitOnError(err)
		}
		for _, cron := range crons.Items {
			imgs := listImages(ctx, cron.Spec.JobTemplate.Spec.Template.Spec.Containers, excludedPrefixes)
			for _, img := range imgs {
				if !slices.ContainsFunc(lst, func(i Image) bool { return i.ref == img }) {
					lst = append(lst, Image{ref: img})
				}
			}
		}
	}
	return lst
}

func isDomain(ctx context.Context, s string) bool {
	if !strings.Contains(s, ".") {
		return false
	}
	_, err := net.DefaultResolver.LookupHost(ctx, s)
	return err == nil
}

func listImages(ctx context.Context, cs []corev1.Container, excludedPrefixes []string) []string {
	lst := make([]string, 0, len(cs))
	for _, c := range cs {
		parts := strings.Split(c.Image, "/")
		var ref string
		switch {
		case isDomain(ctx, parts[0]):
			ref = c.Image
		case len(parts) == 1:
			ref = "docker.io/library/" + c.Image
		default:
			ref = "docker.io/" + c.Image
		}
		if slices.ContainsFunc(excludedPrefixes, func(prefix string) bool {
			return strings.HasPrefix(ref, prefix)
		}) {
			continue
		}
		lst = append(lst, ref)
	}
	return lst
}

func ListFromStdin() []Image {
	lst := parseLines(os.Stdin, false)
	imgs := make([]Image, 0, len(lst))
	for _, ref := range lst {
		imgs = append(imgs, Image{ref: ref})
	}
	return imgs
}

func Export(ctx context.Context, imgs []Image) {
	for _, img := range imgs {
		if !img.inCache || forcePull {
			fmt.Println("Pulling " + img.ref + "...")
			ctr(ctx, "images", "pull", img.ref)
		}
		fmt.Println("Exporting " + img.ref + "...")
		fname := path.Join(snapshotsPath, strings.ReplaceAll(img.ref, "/", "_")+".tar")
		ctr(ctx, "images", "export", fname, img.ref, "--platform", "linux/amd64")
	}
}

var restoreTemplate = template.Must(template.New("restore").Parse(`#!/bin/bash
RL=$(readlink -f "${BASH_SOURCE[0]}")
DN=$(dirname "$RL")
DIR=$(cd -P "$DN" && pwd)
for archive in ` + "`" + `ls $DIR/*.tar` + "`" + `; do
    {{ .ctrBin }} -n k8s.io images import $archive
done`))

func CopyRestoreScript(ctx context.Context) {
	fmt.Println("Copying restore script...")
	fd, err := os.OpenFile(path.Join(snapshotsPath, "restore.sh"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755) //nolint: gosec
	exitOnError(err)
	err = restoreTemplate.Execute(fd, map[string]string{"ctrBin": ctrPath})
	exitOnError(err)
	err = fd.Sync()
	exitOnError(err)
	err = fd.Close()
	exitOnError(err)
}

func Sync(ctx context.Context) {
	fmt.Println("Sync...")
	cmd := exec.CommandContext(ctx, "sync")
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		os.Exit(cmd.ProcessState.ExitCode())
	}
}
