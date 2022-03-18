#!/bin/bash

apt-get -y update
apt-get install docker.io -y

usermod -aG docker $USER

KUBERNETES_VERSION=v1.20.12
KUBEEDGE_VERSION=1.10.0


echo "Downloading Minikube"
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
echo "Installing Minikube"
sudo install minikube-linux-amd64 /usr/local/bin/minikube



echo "Downloading keadm"
curl -LO "https://github.com/kubeedge/kubeedge/releases/download/v${KUBEEDGE_VERSION}/keadm-v${KUBEEDGE_VERSION}-linux-amd64.tar.gz"
tar xzf keadm-v${KUBEEDGE_VERSION}-linux-amd64.tar.gz
echo "Installing keadm"
sudo mv keadm-v${KUBEEDGE_VERSION}-linux-amd64/keadm/keadm /usr/local/bin/

echo "Downloading kubectl"
curl -LO "https://storage.googleapis.com/kubernetes-release/release/${KUBERNETES_VERSION}/bin/linux/amd64/kubectl"
chmod +x ./kubectl
echo "Installing kubectl"
sudo mv ./kubectl /usr/local/bin/kubectl

echo "Starting minikube"
minikube start --apiserver-ips=164.92.76.102 --kubernetes-version=${KUBERNETES_VERSION} --extra-config=kubeadm.ignore-preflight-errors=NumCPU --force --cpus 1

echo "Initializing keadm"
sudo keadm init --kubeedge-version=${KUBEEDGE_VERSION}  --kube-config=$HOME/.kube/config

echo "keadm token saving"
keadm gettoken --kube-config=$HOME/.kube/config > token.txt

echo "Cert generation"
wget https://raw.githubusercontent.com/kubeedge/kubeedge/master/build/tools/certgen.sh
sed -i 's+/etc/kubernetes/pki+$HOME/.minikube+g' certgen.sh
chmod +x certgen.sh
sudo cp certgen.sh /etc/kubeedge
export CLOUDCOREIPS=164.92.76.102
/etc/kubeedge/certgen.sh stream

echo "Installing CloudCore"
cd /etc/kubeedge
sudo ln -s /etc/kubeedge/cloudcore.service /etc/systemd/system/cloudcore.service

echo "Starting CloudCore"
systemctl restart cloudcore.service