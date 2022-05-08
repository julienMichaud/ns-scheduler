docker build . -t ns-scheduler:1
kubectl delete deployment ns-scheduler
sleep 5
minikube image rm docker.io/library/ns-scheduler:1
minikube image load ns-scheduler:1
kubectl apply -f minikube/deployment.yaml