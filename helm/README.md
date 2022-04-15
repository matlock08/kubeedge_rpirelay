# KubeEdge CloudCore 1.10.0 install  on Digital Ocean

The guide is for installing KubeEdge Cloudcore version 1.10 on Digital Ocean 

## Digital Ocean Install

From your Digital Ocean account, you will selected from the create drop down the option Kubertes as shown below:

![Create Kubernetes Cluster](/images/cluster-0.png "Create Kubernetes Cluster")

You need to select the kubernetes version 1.22.7, as Kubeedge v 1.10.0 requieres at least kubernetes 1.22.6 
select the datacenter region more appropiated for you.

![Select Kubernetes Version](/images/cluster-1.png "Select Kubernetes Version")

In the cluster capacity we need to select nodes with a minimum capacity of 4 GB (2.5 GB usable RAM ) to run the pods for kubeedge cloud core,
i'm using a single node (which is not recommended ) for ilustration purpose only.

![Select Cluster Capacity](/images/cluster-2.png "Select Cluster Capacity")

Finally name your cluster and click on Create Cluster

![Create Cluster](/images/cluster-3.png "Create Cluster")

The process to create the cluster with take several minutes, during this time we select one add-on that we will use the NGNIX Ingress Controller

![NGNIX Ingress Controller](/images/cluster-4.png "NGNIX Ingress Controller")

## Kubernetes Configuration

Before adding any KubeEdge Edge nodes we need to perform some configuration in the cluster to prepare it.

First kubeedge communicates node on edge with cloud using different port by default 10000,10001,10002,10003 and 10004 we want to add those
tcp ports to the ingress configuration so it sends all the request to the cloudcore service in the namespace kubeedge.

We need to apply the ConfigMap that is used by the ingress controller to map those ports to the services we want with below command.

```kubectl apply -f tcp-services-ConfigMap.yaml --kubeconfig k8s-master-sfo3-kubeconfig.yaml```

Once we have added the Confimap we need to path the ingress controller to add the same ports and to add the argument that specifices the location of the 
configmap.

```kubectl patch deployment ingress-nginx-controller --namespace ingress-nginx --patch "$(cat patch-deployment-ingress.yaml)" --kubeconfig k8s-master-sfo3-kubeconfig.yaml```

Finally we will path the ingress service to accepts connections on those ports.

```kubectl patch service ingress-nginx-controller --namespace ingress-nginx --patch "$(cat patch-service-ingress.yaml)" --kubeconfig k8s-master-sfo3-kubeconfig.yaml```

Check the EXTERNAL-IP address in the Load Balancer, and take note of it

```kubectl get services -o wide -w ingress-nginx-controller --namespace ingress-nginx --kubeconfig k8s-master-sfo3-kubeconfig.yaml```

Now we need to update that IP address in the file helm/values.yaml in the advertiseAddress

```
    advertiseAddress:
        - "XXX.XXX.XXX.XXX"
```

With that completed we can now execute the helm chart to install the cloudcore in our kubernetes cluster in digital ocean

```helm upgrade --install cloudcore ./cloudcore --namespace kubeedge --create-namespace -f ./cloudcore/values.yaml --kubeconfig k8s-master-sfo3-kubeconfig.yaml```

Once the files are installed we need to modify some of the DaemonSet that are in charge of installed pod in all the nodes of the cluster
some of those applications are not needed or conflict with KubeEdge edgecore nodes.

So will change the affinity to exclude from them all the nodes from edgecore by using the label node-role.kubernetes.io/edge in below daemonset

```
kubectl patch ds cilium --namespace kube-system --patch "$(cat patch-ds.yaml)" --kubeconfig k8s-master-sfo3-kubeconfig.yaml
kubectl patch ds csi-do-node --namespace kube-system --patch "$(cat patch-ds.yaml)" --kubeconfig k8s-master-sfo3-kubeconfig.yaml
kubectl patch ds do-node-agent --namespace kube-system --patch "$(cat patch-ds.yaml)" --kubeconfig k8s-master-sfo3-kubeconfig.yaml
kubectl patch ds kube-proxy --namespace kube-system --patch "$(cat patch-ds.yaml)" --kubeconfig k8s-master-sfo3-kubeconfig.yaml
```

After that we should be able to test that cloudcore is exposed by accessing http://EXTERNAL-IP:10000 and we should get a message like this 
***Client sent an HTTP request to an HTTPS server***.



