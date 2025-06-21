# cluster-health-operator
Cluster health first mate

![M8 the First Mate](./assets/cluster_mate.png)

## Description
Cluster-health-operator is a kubernetes native operator that runs inside your clusters and maintains a constant stream of health checks that can be easily accessed via CRDs

## Getting Started

### Prerequisites
- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

## Engineering Goals (TODO):

### 1. Expand `ClusterIssue` CRD
- [ ] Detect and diagnose:
  - [ ] Pod stuck in `Pending` → scheduling failure reason
  - [ ] Pod stuck in `ContainerCreating` → image pull, volume, CNI, runtime
  - [ ] Pod in `CrashLoopBackOff` → exit codes, liveness probe, OOMKilled
  - [ ] Pod with abnormal restarts (e.g., >5 in 10m) → trend-based flag
- [ ] Add `Reason` enum and `diagnostics engine` per pod phase

### 2. Implement `ClusterHealth` CRD
- [ ] Periodic scoring snapshot
- [ ] Metrics across nodes, pods, etcd, API latency
- [ ] Aggregation from live state + daemonset feeds

### 3. Deploy node-level DaemonSet agent
- [ ] Gather:
  - Node status
  - Allocatable + pressure conditions
  - kubelet + runtime metrics
- [ ] Push to central controller or Redis

### 4. Redis-backed Cache API
- [ ] In-memory cache of cluster state
- [ ] Serve CRD data and raw insights via lightweight REST API
- [ ] TTL-based freshness enforcement

### 5. CLI: `choctl`
- [ ] `choctl issues` – List current `ClusterIssue`s
- [ ] `choctl inspect <pod>` – On-demand diagnostics
- [ ] `choctl health` – Pull latest `ClusterHealth`


### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/cluster-health-operator:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/cluster-health-operator:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following are the steps to build the installer and distribute this project to users.

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/cluster-health-operator:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/cluster-health-operator/<tag or branch>/dist/install.yaml
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

