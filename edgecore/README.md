# KubeEdge EdgeCore 1.10.0 install on Raspberry PI 3

The guide is for installing KubeEdge Edgecore version 1.10 on a Raspberry Pi 3

## Get Secret Token from CloudCore (DigitalOcean)

On you computer with access to the Kubernetes cluster you need to get the access token for the devices 

```
kubectl get secret -n kubeedge tokensecret -o=jsonpath='{.data.tokendata}' --kubeconfig k8s-master-sfo3-kubeconfig.yaml | base64 -d
```

## Installing RPI-RELAY and Raspberry Pi OS 

For this example we are using the [RPi Relay Board](http://www.ingcool.com/wiki/RPi_Relay_Board "RPi Relay Board wiki") from ingcoll 
this board has 3 relays attached to wiringPi pin 25, 28, 29. This is important as we will reference to them using this wiringPi numbers.

After installing the Hat on the raspberry pi we need to install Raspberry Pi OS from 2021-11-08 as the time of writing this guide.

![RPi Relay Board](/images/raspberrypi-hat.jpg "RPi Relay Board")

You need to configure WIFI and enable ssh access for your convinience.

## Adding the device (Raspberry Pi 3) to the Kubernetes cluster

Replace TOKEN by using the token you obtained in the previous step

Replace EXTERNAL-IP by using the LoadBalancer IP of the ingress controller 

You con add custom label like in this example devicetype and sensor when you add the device node to the cluster

```
sudo keadm join --kubeedge-version=1.10.0 -l devicetype=raspberrypi,sensor=rpirelay --cloudcore-ipport=EXTERNAL-IP:10000 --token=TOKEN
```

This command will download the kubeedge v1.10 as well as installing services and mosquitto if not present on the raspberry pi

once the command finish success fully you should be able to see the new nodes listed on the cluster with the edge role

```
$> kubectl get nodes --kubeconfig k8s-master-sfo3-kubeconfig.yaml                                     
NAME                   STATUS   ROLES        AGE     VERSION
node1                  Ready    agent,edge   4h16m   v1.22.6-kubeedge-v1.10.0
pool-ca0gtt0iw-c3nmz   Ready    <none>       24h     v1.22.7
```

## Deploying sample app from your computer with access to Kubernetes Cluster

Deploy edge application on Kubernetes cluster assigned to node1, the application consist on a model , an instance and the deployment.

```
kubectl apply -f rpirelay/crd/model.yaml --kubeconfig k8s-master-sfo3-kubeconfig.yaml

kubectl apply -f rpirelay/crd/instance.yaml --kubeconfig k8s-master-sfo3-kubeconfig.yaml

kubectl apply -f rpirelay/crd/deploy.yaml --kubeconfig k8s-master-sfo3-kubeconfig.yaml
```

## Deploying Web Client to interact with edge device

Now we will deploy a Web Client in the kubernetes cluster to allow us to monitor the state of the relays and change them.
You will need to update the rpi-ingress.yaml to reflect your host domain

```
kubectl apply -f webclient/clusterrole.yaml --kubeconfig k8s-master-sfo3-kubeconfig.yaml

kubectl apply -f webclient/clusterrolebinding.yaml --kubeconfig k8s-master-sfo3-kubeconfig.yaml

kubectl apply -f webclient/rpi-service.yaml --kubeconfig k8s-master-sfo3-kubeconfig.yaml

kubectl apply -f webclient/rpi-ingress.yaml --kubeconfig k8s-master-sfo3-kubeconfig.yaml
```

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