# KubeEdge sample



## Master KubeEdge/cloudcore installation

Installation instruction on ubuntu 18.04 LTS

1.- Download/install minikube

``` [bash]
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
```

2.- Download/install keadm

``` [bash]
curl -LO https://github.com/kubeedge/kubeedge/releases/download/v1.8.2/keadm-v1.8.2-linux-amd64.tar.gz
tar xzf keadm-v1.8.2-linux-amd64.tar.gz
sudo mv keadm-v1.8.2-linux-amd64/keadm/keadm /usr/local/bin/
```

3.- Download/install kubectl

``` [bash]

curl -LO "https://storage.googleapis.com/kubernetes-release/release/v1.20.12/bin/linux/amd64/kubectl"
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl
```

4.- Minikube and kubeedge configuration

``` [bash]

minikube start --kubernetes-version=v1.20.12

sudo keadm init --kubeedge-version=1.8.2  --kube-config=$HOME/.kube/config

keadm gettoken --kube-config=$HOME/.kube/config > token.txt

cd /etc/kubeedge

ln -s cloudcore.service /etc/systemd/system/cloudcore.service

systemctl restart cloudcore.service

wget https://raw.githubusercontent.com/kubeedge/kubeedge/master/build/tools/certgen.sh

sed -i 's+/etc/kubernetes/pki+/home/morpheus/.minikube+g' certgen.sh

export CLOUDCOREIPS="192.168.1.79"

./certgen.sh stream

kubectl get cm tunnelport -nkubeedge -oyaml

iptables -t nat -A OUTPUT -p tcp --dport $YOUR-TUNNEL-PORT -j DNAT --to $CLOUDCOREIPS:10003

``` 

## Node KubeEdge/edgecore installation on raspberry pi 3


``` [bash]
sudo apt-get update && sudo apt-get upgrade

# Install Docker on raspberry
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker pi

# Install keadm for arm
curl -fsSL https://github.com/kubeedge/kubeedge/releases/download/v1.8.2/keadm-v1.8.2-linux-arm.tar.gz -o keadm-v1.8.2-linux-arm.tar.gz

tar xzf keadm-v1.8.2-linux-arm.tar.gz
sudo mv keadm-v1.8.2-linux-arm/keadm/keadm /usr/local/bin/

# add "cgroup_enable=memory cgroup_memory=1" to the end of the line in the file /boot/cmdline.txt
sudo sed 's/$/ cgroup_enable=memory cgroup_memory=1/' /boot/cmdline.txt
 
systemctl restart edgecore.service
```

Add the raspberry pi node to the cluster replace IP and PORT as well as TOKEN with correct values

``` [bash]
sudo keadm join --cloudcore-ipport=IP:PORT --token=TOKEN
```


## Docker login to hub
``` [bash]
docker login -u matlock08
``` 


## Docker build for Raspberry Pi on x86 

First we need to prepare buildx

``` [bash]
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
docker buildx rm
docker buildx create --name mybuilder --driver docker-container --use

docker buildx use mybuilder

docker buildx inspect --bootstrap

# Buid for arm64
docker buildx build --platform linux/arm/v7 -t matlock08/kubeedge-led:0.1.0-arm . --push
```

## Deploy kubernetes
``` [bash]
# Raspberry Pi
kubectl apply -f ./edgecore/node1/charts/pod-led-example.yaml

# Jetson Nano
kubectl apply -f ./edgecore/node3/charts/pod-ds-example.yaml
```

## Labeling node son kubernetes
``` [bash]
# Raspberry Pi
kubectl label nodes node1 devicetype=raspberrypi
kubectl label nodes node2 devicetype=raspberrypi

# Jetson Nano
kubectl label nodes node3 devicetype=jetsonnano
```


## Trobleshooting

``` [bash]
# Edit DaemonSet kube-proxy on kube-system, add this above container

      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
              - matchExpressions:
                  - key: node-role.kubernetes.io/edge
                    operator: DoesNotExist
  


```


kubectl proxy --address='0.0.0.0' --disable-filter=true

http://127.0.0.1:8001/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/