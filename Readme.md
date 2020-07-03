#!/bin/bash

# docker rmi iiicondor/mes_workstatus:1.0.0

docker build -t iiicondor/mes_workstatus:1.0.1 .

docker push iiicondor/mes_workstatus:1.0.1

#部屬步驟

### 建立image
# docker build -t iiicondor/adapter:1.3.0 .
# docker push iiicondor/adapter:1.3.0

### 推上k8s
# 1.修改pod.yaml版本
# image: iiicondor/mes_workstatus:1.0.0
# 2.apply
# kubectl apply -f pod.yaml -n ifactory

# (首次部屬才需要)
# kubectl apply -f service.yaml -n nms
# kubectl apply -f ingress.yaml -n nms
