# Installing Kubernetes on Ubuntu Server
## single node Microk8s snap edition

> See https://github.com/ubuntu/microk8s

## Parameters for VM

Param        | Value       | Comments
-------------|-------------|---
**Vm-Name**  | kube-server
**Disk**     | 400Gb
**Ram**      | 16Gb        | 16384Mib
**Bridge**   | bridge0

## Installer parameters

Param              | Value
-------------------|------
**Subnet**         | 10.10.0.0/16
**Server Address** | 10.10.112.11
**Gateway**        | 10.10.0.254
**Name servers**   | 10.10.0.250, 10.10.0.252
**Search domains** | apa-canal.md
**Disk usage**     | Use entire disk, XFS

## Profile setup

Param           | Value
----------------|------------
**Your name**   | Demius Academius
**Server name** | kube-server
**Username**    | acc-server-admin
**Password**    | shwsa3nn

## Snap setup

    microk8s

## After install

    sudo apt update & sudo apt upgrade -y
    sudo nano /etc/fstab
    sudo swapoff -a
    sudo timedatectl set-timezone Europe/Chisinau
    reboot

## Install kubernetes

skip if already installed

    sudo snap install microk8s --classic

beta version of k8s:

    sudo snap install microk8s --channel=1.14/beta --classic

**Simplify use of kubectl and docker**

    sudo snap alias microk8s.kubectl kubectl
    sudo snap alias microk8s.docker docker

**Enable registry and dashboard**

    microk8s.enable registry dashboard

registry automatically enables storage

dashboard automatically enables grafana,influxdb and heapster

## Deploying official Dashboard UI

**Create Dashboard Account**

    sudo kubectl create -f manifests/dashboard-user.yaml

**Find Bearer Token for login**

    sudo kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep dashboard-user | awk '{print $1}')

**Access to Dashboard UI with NodePort**

    sudo KUBE_EDITOR="nano" kubectl -n kube-system edit service kubernetes-dashboard

edit manifest file for NodePort = 30931
<pre>
    apiVersion: v1
    ...
    name: kubernetes-dashboard
    ...
    spec:
    ...
    ports:
    - port: 443
        protocol: TCP
    targetPort: 8443
    <b>nodePort: 30931</b>
    selector:
        k8s-app: kubernetes-dashboard
    sessionAffinity: None
    type: ClusterIP <b><<< NodePort</b>
    status:
    ...
</pre>

> Access Dashboard via Firefox (certs) at https://10.10.112.11:30931

## Troubleshooting

> if push to local-registry not working, comment localhost ip6 in /etc/hosts

## Additional material
    
> See https://dzone.com/articles/kubernetes-metallb-bare-metal-loadbalancer

> See https://github.com/google/metallb

> See https://metallb.universe.tf/

> See https://help.ubuntu.com/lts/serverguide/advanced-installation.html.en