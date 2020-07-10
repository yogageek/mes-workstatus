# 宣告使用 /bin/bash
#!/bin/bash

version=1.0.3
echo "mes/workstatus version=" ${version}
docker build -t iiicondor/mes_workstatus:${version} . 
docker push iiicondor/mes_workstatus:${version}


# 推上k8s
# 1.修改pod.yaml版本
# image: iiicondor/mes_workstatus:1.0.0

# 2.apply
# kubectl apply -f pod.yaml 

# (首次部屬才需要)
# kubectl apply -f service.yaml 
# kubectl apply -f ingress.yaml 
