#!/bin/bash

set -x
set -o pipefail

minikube start --memory='8g'

helm upgrade --install ingress-nginx ingress-nginx \
 --repo https://kubernetes.github.io/ingress-nginx \
 --namespace ingress-nginx --create-namespace

minikube addons enable ingress

kubectl create secret generic jwt-secret --from-literal=JWT_KEY=asdf

# kubectl port-forward --namespace=ingress-nginx service/ingress-nginx-controller 8080:80
