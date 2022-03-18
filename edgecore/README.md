
## Node Raspberry install 

kubectl get secret -nkubeedge tokensecret -o=jsonpath='{.data.tokendata}' --kubeconfig | base64 -d



sudo keadm join --kubeedge-version=1.10.0 -l devicetype=raspberrypi,sensor=rpirelay --cloudcore-ipport=:10000 --token=


