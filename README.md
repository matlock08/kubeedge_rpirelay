# KubeEdge sample rpirelay



kubectl create clusterrolebinding service-reader-pod --clusterrole=service-reader  --serviceaccount=default:default

affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
              - matchExpressions:
                  - key: node-role.kubernetes.io/edge
                    operator: DoesNotExist