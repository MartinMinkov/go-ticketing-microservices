#!/bin/bash

set -x
set -eo pipefail

kubectl delete deployments --all

kubectl delete all --all -n ingress-nginx

kubectl delete clusterrolebinding ingress-nginx
kubectl delete clusterrolebinding ingress-nginx-admission
kubectl delete clusterrole ingress-nginx-admission
kubectl delete clusterrole ingress-nginx

kubectl delete ingressclass nginx --all-namespaces
kubectl delete ns ingress-nginx
kubectl delete -A ValidatingWebhookConfiguration ingress-nginx-admission

minikube stop

docker system prune -a --volumes -f
docker volume rm $(docker volume ls -q --filter dangling=true)