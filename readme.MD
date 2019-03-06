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
**Disk usage**     | Use entire disk, LVM, XFS

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

## Install kubernetes (skip if already installed)

    sudo snap install microk8s --classic

## Simplify use of kubectl and docker

    sudo snap alias microk8s.kubectl kubectl
    sudo snap alias microk8s.docker docker

## Enable registry (with storage) and dashboard (with grafana,influxdb and heapster)

    microk8s.enable registry dashboard
    
> See https://dzone.com/articles/kubernetes-metallb-bare-metal-loadbalancer

> See https://github.com/google/metallb

> See https://metallb.universe.tf/