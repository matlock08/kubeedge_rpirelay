
## Node Raspberry install 

kubectl get secret -nkubeedge tokensecret -o=jsonpath='{.data.tokendata}' --kubeconfig | base64 -d

2ed7fd8fe90902f7ec5d0a8581a8e04bdf442024d05dd11f2b4f53a9c66ca85e.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDc2NDk2NTN9.DHXhbA3HU-aP_bBUvHLGGfJSb65zQJvbiBV_E0cRQW0


sudo keadm join --kubeedge-version=1.10.0 -l devicetype=raspberrypi,sensor=rpirelay --cloudcore-ipport=143.198.246.186:10000 --token=


