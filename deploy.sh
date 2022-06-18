kubectl delete deployment ns-scheduler
docker build . -t ns-scheduler:latest -f manifests/deploy/Dockerfile
sleep 10
minikube image rm docker.io/library/ns-scheduler:latest
minikube image load ns-scheduler:latest
kubectl apply -f manifests/deploy/deployment.yaml