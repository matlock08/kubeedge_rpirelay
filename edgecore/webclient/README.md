
# Deploying Web Client to interact with edge device

Now we will deploy a Web Client in the kubernetes cluster to allow us to monitor the state of the relays and change them.
You will need to update the rpi-ingress.yaml to reflect your host domain

```
kubectl apply -f webclient/clusterrole.yaml --kubeconfig k8s-master-sfo3-kubeconfig.yaml

kubectl apply -f webclient/clusterrolebinding.yaml --kubeconfig k8s-master-sfo3-kubeconfig.yaml

kubectl apply -f webclient/rpi-service.yaml --kubeconfig k8s-master-sfo3-kubeconfig.yaml

kubectl apply -f webclient/rpi-ingress.yaml --kubeconfig k8s-master-sfo3-kubeconfig.yaml
```

# Testing WebService
Once deployed you should be able to access the Web Service using below commands:

To list the devices available with their model

```
$> curl your.domain.name/device
{"Items":[{"Device":"relay-instance-01","Model":"relay-model"}]}
```

To get the state of an specific device 

```
$> curl your.domain.name/device/relay-instance-01
{"ch1":"ON","ch2":"OFF","ch3":"OFF","device":"relay-instance-01"}
```

And you can change the state of the device by using below command

```
curl -X POST -H "Content-Type: application/json" -d '{"ch1":"ON","ch2":"ON","ch3":"OFF","device":"relay-instance-01"}' https://your.domain.name/device
{"ch1":"ON","ch2":"ON","ch3":"OFF","device":"relay-instance-01"}
```
