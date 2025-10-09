package main

import "github.com/outscale/k8s-image-preloader/cmd"

func main() {
	// var (
	// 	stdin bool

	// 	policyPath string
	// 	insecure   bool

	// 	excludedPrefixes []string
	// )
	// fs := pflag.CommandLine
	// fs.BoolVar(&stdin, "stdin", false, "Read image list from standard input")
	// fs.StringVar(&policyPath, "policy", "", "Path to a trust policy file")
	// fs.BoolVar(&insecure, "insecure", false, "run the tool without any policy check")
	// fs.StringSliceVar(&excludedPrefixes, "exclude", nil, "Prefixes to skip from the list")
	// pflag.Parse()

	// ctx := context.Background()

	// kubeconfig := os.Getenv("KUBECONFIG")
	// if kubeconfig == "" {
	// 	kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	// }

	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// exitOnError(err)
	// k8s, err := kubernetes.NewForConfig(config)
	// exitOnError(err)
	// var lst []string
	// if stdin {
	// 	lst, err = preloader.ListFromStdin(ctx)
	// } else {
	// 	lst, err = preloader.ListFromCluster(ctx, k8s, excludedPrefixes)
	// }
	// exitOnError(err)
	// err = preloader.Copy(ctx, lst, policyPath, insecure)
	// exitOnError(err)
	// err = preloader.Snapshot(ctx, k8s)
	// exitOnError(err)
	cmd.Execute()
}
