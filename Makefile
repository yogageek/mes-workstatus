# TAG ?= $(shell git describe --tags --always)
TAG ?= 1.3.0


#k8s重拉env
# 刪deployment kubectl delete -n nms deployment.apps/smart-exporter 
# pod要加ifnotPresent   imagePullPolicy: Always
# apply kubectl apply -f pod.yaml -n nms
# 查看pod

#run=smart-exporter
#nms-controller-api
#app=rmq-sb105-pod
POD_LABEL=run=nms-adapter-api
NAMESPACE=-n nms
DEPLOYMENTNAME=deployment.apps/nms-adapter-api

POD_ID 	?= $(shell kubectl $(NAMESPACE) -o name get pods -l $(POD_LABEL))

#先取label填入變數
label: 
	kubectl get pod --show-labels $(NAMESPACE)
a:
	kubectl get all $(NAMESPACE)
#刪除目前部屬(重新部屬前做)
d:
	kubectl delete $(NAMESPACE) $(DEPLOYMENTNAME)

#取當前pod name
pn:
	kubectl get pods -l $(POD_LABEL) $(NAMESPACE) -o name

#取所有pod name
allpn:
	kubectl  $(NAMESPACE) get pods -o name	

#取得最新log
newlog:
	kubectl $(NAMESPACE) logs $(POD_ID) -f --tail=20 

#取得最新log
logs: newpodid newlog

#根據POD_NAME取得最新podid
newpodid:
	kubectl $(NAMESPACE) -o name get pods -l $(POD_LABEL)	
	@echo ----->$(POD_LABEL)
	$(eval POD_ID= $(shell kubectl -n nms -o name get pods -l $(POD_LABEL)))

#印出pod id
pid:
	@echo $(POD_ID)

external:
	kubectl get node -o wide $(NAMESPACE)


deploy:
	kubectl apply -f pod.yaml -n nms



w:
	kubectl get pod -l $(POD_LABEL) -w
	
watch:
	watch $(NAMESPACE) 1 kubectl get pod $(POD_ID)
	
ww:
	kubectl get pod $(POD_ID) -w




# --selector app=eventstore -o jsonpath="{.items[*].metadata.name}"


# 1.(首次部屬才需要)
# kubectl apply -f service.yaml -n nms
# kubectl apply -f ingress.yaml -n nms
# 2.apply
# (如果更新image需要先刪deployment)
# kubectl apply -f pod.yaml -n nms





