---
title: "Known Issues"
menu:
  main:
    parent: "user"
    identifier: "known-issues"
    weight: 2
---
# Known Issues

Having problems with kind? This guide is covers some known problems and solutions / workarounds.

It may additionally be helpful to:

- check our [issue tracker]
- [file an issue][file an issue] (if there isn't one already)
- reach out and ask for help in [#kind] on the [kubernetes slack]

## Contents
* [kubectl version skew](#kubectl-version-skew)
* [Older Docker Installations](#older-docker-installations)
* [Docker on Btrfs or ZFS](#docker-on-btrfs-or-zfs)
* [Docker Installed With Snap](#docker-installed-with-snap)
* [Failing to apply overlay network](#failing-to-apply-overlay-network)
* [Failure to build node image](#failure-to-build-node-image)
* [Failing to properly start cluster](#failure-to-properly-start-cluster)
* [Pod errors due to "too many open files"](#pod-errors-due-to-too-many-open-files)
* [Docker permission denied](#docker-permission-denied)
* [Windows Containers](#windows-containers)
* [Non-AMD64 Architectures](#nonamd64-architectures)
* [Unable to pull images](#unable-to-pull-images)
* [Chrome OS](#chrome-os)
* [AppArmor](#apparmor)
* [IPv6 port forwarding](#ipv6-port-forwarding)

## Kubectl Version Skew

You may have problems interacting with your kind cluster if your client(s) are
skewed too far from the kind node version. Kubernetes [only supports limited skew][version skew]
between clients and the API server.

This is a issue that frequently occurs when running `kind` alongside Docker For Mac.

This problem is related to a bug in [docker on macOS][for-mac#3663]

If you see something like the following error message:
```bash
$ kubectl edit deploy -n kube-system kubernetes-dashboard
error: SchemaError(io.k8s.api.autoscaling.v2beta1.ExternalMetricStatus): invalid object doesn't have additional properties
```

You can check your client and server versions by running:
{{< codeFromInline lang="bash" >}}
kubectl version
{{< /codeFromInline >}}

If there is a mismatch between the server and client versions, you should install a newer client version. 

If you are using Mac, you can install kubectl via homebrew by running:
{{< codeFromInline lang="bash" >}}
brew install kubernetes-cli
{{< /codeFromInline >}}

And overwrite the symlinks created by Docker For Mac by running:
{{< codeFromInline lang="bash" >}}
brew link --overwrite kubernetes-cli
{{< /codeFromInline >}}

[for-mac#3663]: https://github.com/docker/for-mac/issues/3663

## Older Docker Installations

`kind` is known to have issues with Kubernetes 1.13 or lower when using Docker versions:

- `1.13.1` (released January 2017)
- `17.05.0-ce` (released May 2017)

And possibly other old versions of Docker.

With these versions you must use Kubernetes >= 1.14, or more ideally upgrade Docker instead.

kind is tested with a recent stable `docker-ce` release.

## Docker on Btrfs or ZFS

`kind` cannot run properly if containers on your machine / host are backed by a
[Btrfs](https://en.wikipedia.org/wiki/Btrfs) or [ZFS](https://en.wikipedia.org/wiki/ZFS)
filesystem.

This should only be relevant on Linux, on which you can check with:

{{< codeFromInline lang="bash" >}}
docker info | grep -i storage
{{< /codeFromInline >}}

As a workaround, you'll need to ensure that the storage driver is one that works.
Docker's default of `overlay2` is a good choice, but `aufs` should also work.

You can set the storage driver with the following configuration in `/etc/docker/daemon.json`:

```
{
  "storage-driver": "overlay2"
}
```

After restarting the Docker daemon you should see that Docker is now using the `overlay2` storage driver instead of `btrfs`.


**NOTE**: You'll need to make sure the backing filesystem is not btrfs / ZFS as well,
which may require creating a partition on your host disk with a suitable filesystem and ensuring Docker's
data root is on this (by default `/var/lib/docker`). Ext4 is a reasonable choice.

## Docker Installed with Snap

If you installed Docker with [snap], it is likely that `docker` commands do not
have access to `$TMPDIR`. This may break some kind commands which depend
on using temp directories (`kind build ...`).

Currently a workaround for this is setting the `TEMPDIR` environment variable to
a directory snap does have access to when working with kind.
This can for example be some directory under `$HOME`.

## Failing to apply overlay network
There are two known causes for problems while applying the overlay network
while building a kind cluster:

* Host machine is behind a proxy
* Usage of Docker version 18.09

If you see something like the following error message:
```
 ✗ [kind-1-control-plane] Starting Kubernetes (this may take a minute) ☸
FATA[07:20:43] Failed to create cluster: failed to apply overlay network: exit status 1
```

or the following, when setting the `loglevel` flag to debug,
```
DEBU[16:26:53] Running: /usr/bin/docker [docker exec --privileged kind-1-control-plane /bin/sh -c kubectl apply --kubeconfig=/etc/kubernetes/admin.conf -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version --kubeconfig=/etc/kubernetes/admin.conf | base64 | tr -d '\n')"]
ERRO[16:28:25] failed to apply overlay network: exit status 1 ) ☸
 ✗ [control-plane] Starting Kubernetes (this may take a minute) ☸
ERRO[16:28:25] failed to apply overlay network: exit status 1
DEBU[16:28:25] Running: /usr/bin/docker [docker ps -q -a --no-trunc --filter label=io.k8s.sigs.kind.cluster --format {{.Names}}\t{{.Label "io.k8s.sigs.kind.cluster"}} --filter label=io.k8s.sigs.kind.cluster=1]
DEBU[16:28:25] Running: /usr/bin/docker [docker rm -f -v kind-1-control-plane]
⠈⠁ [control-plane] Pre-loading images 🐋 Error: failed to create cluster: failed to apply overlay network: exit status 1
```

The issue may be due to your host machine being behind a proxy, such as in
[kind#136][kind#136].
We are currently looking into ways of mitigating this issue by preloading CNI
artifacts, see [kind#200][kind#200].
Another possible solution is to enable kind nodes to use a proxy when
downloading images, see [kind#270][kind#270].

The last known case for this issue comes from the host machine 
[using Docker 18.09 in kind#136][kind#136-docker]. 
In this case, a known solution is to upgrade to any Docker version greater than or
equal to Docker 18.09.1.


## Failure to build node image
Building kind's node image may fail due to running out of memory on Docker for Mac or Docker for Windows.
See [kind#229][kind#229].

If you see something like this:

```
    cmd/kube-scheduler
    cmd/kube-proxy
/usr/local/go/pkg/tool/linux_amd64/link: signal: killed
!!! [0116 08:30:53] Call tree:
!!! [0116 08:30:53]  1: /go/src/k8s.io/kubernetes/hack/lib/golang.sh:614 kube::golang::build_some_binaries(...)
!!! [0116 08:30:53]  2: /go/src/k8s.io/kubernetes/hack/lib/golang.sh:758 kube::golang::build_binaries_for_platform(...)
!!! [0116 08:30:53]  3: hack/make-rules/build.sh:27 kube::golang::build_binaries(...)
!!! [0116 08:30:53] Call tree:
!!! [0116 08:30:53]  1: hack/make-rules/build.sh:27 kube::golang::build_binaries(...)
!!! [0116 08:30:53] Call tree:
!!! [0116 08:30:53]  1: hack/make-rules/build.sh:27 kube::golang::build_binaries(...)
make: *** [all] Error 1
Makefile:92: recipe for target 'all' failed
!!! [0116 08:30:54] Call tree:
!!! [0116 08:30:54]  1: build/../build/common.sh:518 kube::build::run_build_command_ex(...)
!!! [0116 08:30:54]  2: build/release-images.sh:38 kube::build::run_build_command(...)
make: *** [quick-release-images] Error 1
ERRO[08:30:54] Failed to build Kubernetes: failed to build images: exit status 2
Error: error building node image: failed to build kubernetes: failed to build images: exit status 2
Usage:
  kind build node-image [flags]

Flags:
      --base-image string   name:tag of the base image to use for the build (default "kindest/base:v20181203-d055041")
  -h, --help                help for node-image
      --image string        name:tag of the resulting image to be built (default "kindest/node:latest")
      --kube-root string    Path to the Kubernetes source directory (if empty, the path is autodetected)
      --type string         build type, one of [bazel, docker] (default "docker")

Global Flags:
      --loglevel string   logrus log level [panic, fatal, error, warning, info, debug] (default "warning")

error building node image: failed to build kubernetes: failed to build images: exit status 2
```

Then you may try increasing the resource limits for the Docker engine on Mac or Windows.

It is recommended that you allocate at least 8GB of RAM to build Kubernetes.

Open the **Preferences** menu.

<img src="/docs/user/images/docker-pref-1.png"/>

Go to the **Advanced** settings page, and change the settings there, see 
[changing Docker's resource limits][Docker resource lims].

<img width="400px" src="/docs/user/images/docker-pref-build.png" alt="Setting 8Gb of memory in Docker for Mac" />

<img width="400px" src="/docs/user/images/docker-pref-build-win.png" alt="Setting 8Gb of memory in Docker for Windows" />

## Failing to properly start cluster
This issue is similar to a 
[failure while building the node image](#failure-to-build-node-image).
If the cluster creation process was successful but you are unable to see any
Kubernetes resources running, for example:

```
$ docker ps
CONTAINER ID        IMAGE                  COMMAND                  CREATED              STATUS              PORTS                      NAMES
c0261f7512fd        kindest/node:v1.12.2   "/usr/local/bin/entr…"   About a minute ago   Up About a minute   0.0.0.0:64907->64907/tcp   kind-1-control-plane
$ docker exec -it c0261f7512fd /bin/sh
# docker ps -a
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
#
```

or `kubectl` being unable to connect to the cluster,
```
$ kind export kubeconfig
$ kubectl cluster-info

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
Unable to connect to the server: EOF
```

Then as in [kind#156][kind#156], you may solve this issue by claiming back some
space on your machine by removing unused data or images left by the Docker 
engine by running:
{{< codeFromInline lang="bash" >}}
docker system prune
{{< /codeFromInline >}}

And / or:
{{< codeFromInline lang="bash" >}}
docker image prune
{{< /codeFromInline >}}

You can verify the issue by exporting the logs (`kind export logs`) and looking
at the kubelet logs, which may have something like the following:
```
Dec 07 00:37:53 kind-1-control-plane kubelet[688]: I1207 00:37:53.229561     688 eviction_manager.go:340] eviction manager: must evict pod(s) to reclaim ephemeral-storage
Dec 07 00:37:53 kind-1-control-plane kubelet[688]: E1207 00:37:53.229638     688 eviction_manager.go:351] eviction manager: eviction thresholds have been met, but no pods are active to evict
```

If on the other hand you are running kind on a btrfs partition and your logs
show something like the following:
```
F0103 17:42:41.470269 3804 kubelet.go:1359] Failed to start ContainerManager failed to get rootfs info: failed to get device for dir "/ var/lib/kubelet": could not find device with major: 0, minor: 67 in cached partitions map
```

This problem seems to be related to a [bug in Docker][moby#9939].

## Pod errors due to "too many open files"

This may be caused by running out of [inotify](https://linux.die.net/man/7/inotify) resources. Resource limits are defined by `fs.inotify.max_user_watches` and `fs.inotify.max_user_instances` system variables. For example, in Ubuntu these default to 8192 and 128 respectively, which is not enough to create a cluster with many nodes.

To increase these limits temporarily run the following commands on the host:
{{< codeFromInline lang="bash" >}}
sudo sysctl fs.inotify.max_user_watches=524288
sudo sysctl fs.inotify.max_user_instances=512
{{< /codeFromInline >}}

To make the changes persistent, edit the file `/etc/sysctl.conf` and add these lines:
{{< codeFromInline lang="bash" >}}
fs.inotify.max_user_watches = 524288
fs.inotify.max_user_instances = 512
{{< /codeFromInline >}}

## Docker permission denied

When using `kind`, we assume that the user you are executing kind as has permission to use docker.
If you initially ran Docker CLI commands using `sudo`, you may see the following error, which indicates that your `~/.docker/` directory was created with incorrect permissions due to the `sudo` commands.

```
WARNING: Error loading config file: /home/user/.docker/config.json
open /home/user/.docker/config.json: permission denied
```

To fix this problem, either follow the docker's docs [manage docker as a non root user][manage docker as a non root user],
or try to use `sudo` before your commands (if you get `command not found` please check [this comment about sudo with kind][sudo with kind]).

## Windows Containers

[Docker Desktop for Windows][docker desktop for windows] supports running both Linux (the default) and Windows Docker containers.

`kind` for Windows requires Linux containers. To switch between Linux and Windows containers see [this page][switch between windows and linux containers].

Windows containers are not like Linux containers and do not support running docker in docker and therefore cannot support kind.

## Non-AMD64 Architectures

KIND does not currently ship pre-built images for non-amd64 architectures.
In the future we may, but currently demand has been low and the cost to build
has been high.

To use kind on other architectures, you need to first build a base image
and then build a node image.

Run `images/base/build.sh` and then taking note of the built image name use `kind build node-image --base-image=kindest/base:tag-i-built`.

There are more details about how to do this in the [Quick Start] guide.

## Unable to pull images

When using named KIND instances you may sometimes see your images failing to pull correctly on pods. This will usually manifest itself with the following output when doing a `kubectl describe pod my-pod`

```
Failed to pull image "docker.io/my-custom-image:tag": rpc error: code = Unknown desc = failed to resolve image "docker.io/library/my-custom-image:tag": no available registry endpoint: pull access denied, repository does not exist or may require authorization: server message: insufficient_scope: authorization failed
```

If this image has been loaded onto your kind cluster using the command `kind load docker-image my-custom-image` then you have likely not provided the name parameter. 

Re-run the command this time adding the `--name my-cluster-name` param:

`kind load docker-image my-custom-image --name my-cluster-name`

## Chrome OS

Kubernetes does not work in the Chrome OS Linux sandbox.

Please see the upstream issue https://bugs.chromium.org/p/chromium/issues/detail?id=878034

For previous discussion see: https://github.com/kubernetes-sigs/kind/issues/763

## AppArmor

If your host has [AppArmor] enabled you may run into [moby/moby/issues/7512](https://github.com/moby/moby/issues/7512#issuecomment-51845976).

You will likely need to disable apparmor on your host or at least any profile(s)
related to applications you are trying to run in KIND.

See Previous Discussion: [kind#1179]

## IPv6 Port Forwarding

Docker assumes that all the IPv6 addresses should be reachable, hence doesn't implement
port mapping using NAT [moby#17666].

You will likely need to use Kubernetes services like NodePort or LoadBalancer to access 
your workloads inside the cluster via the nodes IPv6 addresses.

See Previous Discussion: [kind#1326]

[issue tracker]: https://github.com/kubernetes-sigs/kind/issues
[file an issue]: https://github.com/kubernetes-sigs/kind/issues/new
[#kind]: https://kubernetes.slack.com/messages/CEKK1KTN2/
[kubernetes slack]: http://slack.k8s.io/
[kind#136]: https://github.com/kubernetes-sigs/kind/issues/136
[kind#136-docker]: https://github.com/kubernetes-sigs/kind/issues/136#issuecomment-457015838
[kind#156]: https://github.com/kubernetes-sigs/kind/issues/156
[kind#182]: https://github.com/kubernetes-sigs/kind/issues/182
[kind#200]: https://github.com/kubernetes-sigs/kind/issues/200
[kind#229]: https://github.com/kubernetes-sigs/kind/issues/229
[kind#270]: https://github.com/kubernetes-sigs/kind/issues/270
[kind#1179]: https://github.com/kubernetes-sigs/kind/issues/1179
[kind#1326]: https://github.com/kubernetes-sigs/kind/issues/1326
[moby#9939]: https://github.com/moby/moby/issues/9939
[moby#17666]: https://github.com/moby/moby/issues/17666
[Docker resource lims]: https://docs.docker.com/docker-for-mac/#advanced
[snap]: https://snapcraft.io/
[manage docker as a non root user]: https://docs.docker.com/install/linux/linux-postinstall/#manage-docker-as-a-non-root-user
[sudo with kind]: https://github.com/kubernetes-sigs/kind/issues/713#issuecomment-512665315
[docker desktop for windows]: https://hub.docker.com/editions/community/docker-ce-desktop-windows
[switch between windows and linux containers]: https://docs.docker.com/docker-for-windows/#switch-between-windows-and-linux-containers
[version skew]: https://kubernetes.io/docs/setup/release/version-skew-policy/#supported-version-skew
[Quick Start]: /docs/user/quick-start
[AppArmor]: https://en.wikipedia.org/wiki/AppArmor
