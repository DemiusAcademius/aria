# Installing Kubernetes on Ubuntu Server
## single node Microk8s snap edition

> See https://github.com/ubuntu/microk8s

## Parameters for VM

Param        | Value
-------------|------------
**Vm-Name**  | kube-server
**Disk**     | 400Gb
**Ram**      | 16Gb
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

## After install

    sudo apt update & sudo apt upgrade -y
    sudo nano /etc/fstab
    sudo swapoff -a
    sudo timedatectl set-timezone Europe/Chisinau
    reboot

## Install kubernetes

    sudo snap install microk8s --classic
    sudo snap alias microk8s.kubectl kubectl

## Enable storage & registry

    microk8s.enable storage registry dashboard
    
> See https://dzone.com/articles/kubernetes-metallb-bare-metal-loadbalancer

> See https://github.com/google/metallb

> See https://metallb.universe.tf/