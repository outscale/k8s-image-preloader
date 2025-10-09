package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	snapshotclient "github.com/kubernetes-csi/external-snapshotter/client/v8/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getK8SClient() *kubernetes.Clientset {
	// kubeconfig := os.Getenv("KUBECONFIG")
	// if kubeconfig == "" {
	// 	kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	// }
	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	config, err := rest.InClusterConfig() // in cluster config
	exitOnError(err)
	k8s, err := kubernetes.NewForConfig(config)
	exitOnError(err)
	return k8s
}

func getSnapshotClient() *snapshotclient.Clientset {
	// kubeconfig := os.Getenv("KUBECONFIG")
	// if kubeconfig == "" {
	// 	kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	// }
	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	config, err := rest.InClusterConfig() // in cluster config
	exitOnError(err)
	k8s, err := snapshotclient.NewForConfig(config)
	exitOnError(err)
	return k8s
}

func ctr(ctx context.Context, args ...string) []string {
	args = append([]string{"-n", "k8s.io"}, append(strings.Split(ctrFlags, " "), args...)...)
	cmd := exec.CommandContext(ctx, ctrPath, args...)
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, ctrPath, strings.Join(args, " "))
		fmt.Fprintln(os.Stderr, stderr.String())
		os.Exit(cmd.ProcessState.ExitCode())
	}
	return parseLines(stdout, debug)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "An error occurred:", err)
		os.Exit(1)
	}
}

func parseLines(r io.Reader, debut bool) []string {
	scan := bufio.NewScanner(r)
	var lst []string
	for scan.Scan() {
		line := scan.Text()
		if debug {
			fmt.Println(line)
		}
		lst = append(lst, line)
	}
	exitOnError(scan.Err())
	return lst
}
