# Kubeedge-traffic-light

- build at edge node:

```bash
$ make
```

- create crds at cloud node:

```bash
$ cd crd
$ kubectl apply -f model.yaml
# replace "<your edge node name>" with your edge node name
$ sed -i 's#raspberrypi#<your edge node name>#' instance.yaml
$ kubectl apply -f instance.yaml
```

**Note: instance must be created after model and deleted before model.**

- create demo at cloud node:

```bash
$ kubectl apply -f deploy.yaml

kubectl edit device relay-instance-01

kubectl --kubeconfig k8s-1-21-9-do-0-sfo3-1645972985631-kubeconfig.yaml get secret -nkubeedge tokensecret -o=jsonpath='{.data.tokendata}'

keadm join --cloudcore-ipport=143.244.210.104:10000 --token=8ba332a0c81ea0bd79c823b8d122ae49f4acfc0d9b7b468666810807d4fd8363.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDY0MzEyMDZ9.0OduQXhkMZlVFkt47ktTT9XH7m8o9DOEX0SsQCkA7C0
```
