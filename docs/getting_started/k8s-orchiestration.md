
You can run your workload locally, from remote servers, or directly on a Kubernetes cluster. However, deploying your workload on Kubernetes offers several advantages, including:

- Easy Scaling: Kubernetes enables you to scale your workload effortlessly by adjusting the number of replicas or utilizing features like Horizontal Pod Autoscaler.

- Resource Management: Kubernetes provides robust resource management capabilities, allowing you to efficiently allocate and manage resources such as CPU and memory for your workload. You can overwrite default values using `helm-set` flag.


### Installation
To install your workload using Loadbot on a Kubernetes cluster, you can use the following command:

```bash
loadbot install \
    --context dev \
    --namespace default \
    --helm-set workload.replicas=2 \
    --workload-config config.json \
    myworkload 

```

### Uninstall
To uninstall your workload from the Kubernetes cluster, you can use the following command:

```bash
loadbot uninstall \
    --context dev \
    --namespace default \
    myworkload 

```

### Upgrade
To upgrade your workload on the Kubernetes cluster, you can use the following command:

```bash
loadbot upgrade \
    --context dev \
    --namespace default \
    --helm-set workload.replicas=2 \
    --workload-config config.json \
    myworkload 

```

### List
To list the workloads deployed on the Kubernetes cluster, you can use the following command:

```bash
loadbot list --context dev --namespace default
```
