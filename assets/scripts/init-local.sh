#!/bin/bash

curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash

kubectl config use-context docker-desktop

kubectl create namespace dev
kubectl config set-context --current --namespace dev

make create-configs-local

kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.2.0/cert-manager.crds.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.0.1/deploy/static/provider/cloud/deploy.yaml

echo 'waiting for ingress container to set up...'
kubectl wait --for=condition=ready --timeout=300s -n ingress-nginx pod -l app.kubernetes.io/component=controller
kubectl apply -f assets/k8s/ingress.yaml
