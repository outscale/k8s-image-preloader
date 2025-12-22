# Kubernetes image preloader

[![Project Sandbox](https://docs.outscale.com/fr/userguide/_images/Project-Sandbox-yellow.svg)](https://docs.outscale.com/en/userguide/Open-Source-Projects.html) [![](https://dcbadge.limes.pink/api/server/HUVtY5gT6s?style=flat&theme=default-inverted)](https://discord.gg/HUVtY5gT6s)

<p align="center">
  <img alt="Kubernetes Logo" src="https://upload.wikimedia.org/wikipedia/commons/3/39/Kubernetes_logo_without_workmark.svg" width="120px">
</p>

---

## 🌐 Links

- Documentation: <https://docs.outscale.com/en/>
- Project website: <https://github.com/outscale/<project-name>>
- Join our community on [Discord](https://discord.gg/HUVtY5gT6s)

---

## 📄 Table of Contents

- [Overview](#-overview)
- [Requirements](#-requirements)
- [Usage](#-usage)
- [License](#-license)
- [Contributing](#-contributing)

---

## 🧭 Overview

**Kubernetes image preloader** is a tool that allows you to create a snapshot of images that can be preloaded onto new Kubernetes nodes.

It enables faster startup times and allows nodes to run without internet access.

---

## ✅ Requirements

- A Kubernetes cluster running containerd,
- A CSI driver with snapshotting enabled.

---

## 🚀 Usage

### Creating a snapshot

```bash
kubectl apply -f example/example.yaml
```

This will:
- list all images present in the local containerd cache,
- list all images used by pods/cronjobs on the local cluster and not present in the cache,
- refetch all images to the local containerd cache,
- export all images to a PVC,
- copy a restore script to the PVC,
- take a CSI snapshot of the volume.

> By default, containerd purges some layers, pulling all images again is required before exporting images.

### Preloading images

To preload the images stored on the snapshot, you will need to attach a volume based on the snapshot to the node VM.

Using cluster api, in a OscMachineTemplate resource:
```yaml
        vm:
          [...]
        volumes:
        - device: /dev/xvdb
          size: 2
          type: gp2
          fromSnapshot: snap-xxx
```

> This requires CAPOSC v1.2.0 or later.

And preload images in the KubeadmConfigTemplate/KubeadmControlPlane resources:
```yaml
  spec:
    joinConfiguration:
      [...]
    mounts:
      - - xvdb
        - /preload
        - ext4
        - auto,exec,ro
    preKubeadmCommands:
      - /preload/restore.sh
```

### Manual use

You will need to run:
* `preloader export` to export images to the volume,
* `preloader snapshot` to create a snapshot of the volume,
* `ctr images import` to import each image.

### preloader export flags

```bash
Export a list of images to a path

By default, fetches the list of all images from the local cluster,
reading the local containerd cache and the list of pods/cronjobs,
and exports all images found to path.

Usage:
  preloader export [flags]

Aliases:
  export, e

Flags:
      --cache-only          Only list images present in local cache, do not list pods
      --exclude strings     Prefixes to skip from the list
      --force-pull          Force an image pull before exporting
  -h, --help                help for export
      --no-restore-script   Do not copy restore script to the volume
      --stdin               Fetch image list from stdin instead of the local cluster
      --to string           Path to snapshot volume (default "/snapshot")

Global Flags:
      --ctr-flags string   ctr flags (default "-a /var/run/containerd/containerd.sock")
      --ctr-path string    ctr binary path (default "/usr/local/bin/ctr")
      --debug              log ctr command output
```

### preloader snapshot flags

```bash
Create a VolumeSnapshot from the PVC storing exports of images.

Usage:
  preloader snapshot [flags]

Aliases:
  snapshot, s

Flags:
      --class string       VolumeSnapshotClass to use (required)
  -h, --help               help for snapshot
      --name string        Name of the VolumeSnapshot to create (required)
      --namespace string   Namespace of the VolumeSnapshot to create (required)
      --pvc string         PVC to snashot (required)

Global Flags:
      --ctr-flags string   ctr flags (default "-a /var/run/containerd/containerd.sock")
      --ctr-path string    ctr binary path (default "/usr/local/bin/ctr")
      --debug              log ctr command output
```

---

## 📜 License

**Kubernetes image preloader** is released under the BSD 3-Clause license.

© 2025 Outscale SAS

See [LICENSE](./LICENSE) for full details.

---

## 🤝 Contributing

We welcome contributions!

Please read our [Contributing Guidelines](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md) before submitting a pull request.
