# KubeEdge sample rpirelay

## Digital Ocean Install

kubectl apply -f tcp-services-ConfigMap.yaml --kubeconfig 

kubectl patch deployment ingress-nginx-controller --namespace ingress-nginx --patch "$(cat patch-deployment-ingress.yaml)" --kubeconfig k8s-master-sfo3-kubeconfig.yaml

kubectl patch service ingress-nginx-controller --namespace ingress-nginx --patch "$(cat patch-service-ingress.yaml)" --kubeconfig k8s-master-sfo3-kubeconfig.yaml

helm upgrade --install cloudcore ./cloudcore --namespace kubeedge --create-namespace -f ./cloudcore/values.yaml --kubeconfig 

kubectl patch ds cilium --namespace kube-system --patch "$(cat patch-ds-cilium.yaml)" --kubeconfig k8s-master-sfo3-kubeconfig.yaml
kubectl patch ds csi-do-node --namespace kube-system --patch "$(cat patch-ds-cilium.yaml)" --kubeconfig k8s-master-sfo3-kubeconfig.yaml
kubectl patch ds do-node-agent --namespace kube-system --patch "$(cat patch-ds-cilium.yaml)" --kubeconfig k8s-master-sfo3-kubeconfig.yaml
kubectl patch ds kube-proxy --namespace kube-system --patch "$(cat patch-ds-cilium.yaml)" --kubeconfig k8s-master-sfo3-kubeconfig.yaml

## Node Raspberry install 

kubectl get secret -nkubeedge tokensecret -o=jsonpath='{.data.tokendata}' --kubeconfig | base64 -d

2ed7fd8fe90902f7ec5d0a8581a8e04bdf442024d05dd11f2b4f53a9c66ca85e.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDc2NDk2NTN9.DHXhbA3HU-aP_bBUvHLGGfJSb65zQJvbiBV_E0cRQW0


sudo keadm join --kubeedge-version=1.10.0 -l devicetype=raspberrypi,sensor=rpirelay --cloudcore-ipport=143.198.246.186:10000 --token=











kubectl create clusterrolebinding service-reader-pod --clusterrole=service-reader  --serviceaccount=default:default

affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
              - matchExpressions:
                  - key: node-role.kubernetes.io/edge
                    operator: DoesNotExist