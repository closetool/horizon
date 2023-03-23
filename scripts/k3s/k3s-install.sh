#!/bin/bash

WEBSIDE="google.com"
INSTALL_K3S_MIRROR=cn
INSTALL_K3S_EXEC="--docker"

if ! which docker >/dev/null 2>&1; then
    echo "==================> Docker is not installed"
	curl https://releases.rancher.com/install-docker/20.10.sh | sh
	sudo systemctl start docker
elif ! docker info &> /dev/null
then
    echo "==================> Docker is installed but the daemon is not running"
    sudo systemctl start docker
else
    echo "==================> Docker is already installed and running"
fi

ping -c 1 $WEBSIDE > /dev/null
if [ $? -eq 0 ]
then
    pwd
    echo "==================> install k3s k3s/get-k3s" ;pwd
	chmod +x ./k3s/get-k3s.sh
    INSTALL_K3S_MIRROR=cn INSTALL_K3S_EXEC="--docker"  ./k3s/get-k3s.sh
else
    echo "==================> install k3s kes/get-k3s-zh"  ;pwd 
	chmod +x ./k3s/get-k3s-zh.sh
    INSTALL_K3S_MIRROR=cn INSTALL_K3S_EXEC="--docker"  ./k3s/get-k3s-zh.sh
fi

## install helm
# TODO: Install using the helm controller
if ! command -v helm >/dev/null 2>&1; then
then
    echo "==================> Helm is not installed, downloading and installing..."
    curl -fsSL -o /tmp/get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
    chmod 700 /tmp/get_helm.sh
    /tmp/get_helm.sh
fi

## Minimum installation
echo "==================> Install horizon, minimum installation"
helm install horizon horizoncd/horizon -n horizoncd --version 2.1.7 --create-namespace -f https://raw.githubusercontent.com/horizoncd/helm-charts/main/horizon-cn-values.yaml

kubectl get pod -n horizoncd

# TODO: http://horizon.h8r.site/