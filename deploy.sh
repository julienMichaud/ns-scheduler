kubectl delete deployment ns-scheduler
docker build . -t ns-scheduler:2
sleep 61
minikube image rm docker.io/library/ns-scheduler:2
minikube image load ns-scheduler:2
kubectl apply -f minikube/deployment.yaml