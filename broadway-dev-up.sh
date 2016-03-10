#!/bin/bash

docker-compose up -d

# ADD TUNNELING TO K8s
dm=$(docker-machine active)
docker-machine ssh "$dm" -f -N -L "8080:localhost:8080"

echo "waiting"

until /usr/local/bin/kubectl get pods &> /dev/null; do
   printf "."
done

echo

kubectl create -f broadway-namespace.yaml
kubectl create -f broadway-rc.yaml
