# Installing Kuberenetes on Ubuntu

[see official docs](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/)

## demo of real configuration ##

    Content folder contains configs and scripts


    Sources folder contains utilities and scripts written in GO

    

## Common parameters

Param | Value  | Comments
------|--------|----------
**VM-Name**    | aria-server
**disk**       | 400Gb
**RAM**        | 16Gb    | 16384Mib
**CPUs**       | 4
**Network**    | bridge0 | Shared device name
**OS Version** | Ubuntu server 18.10

## Linux Installer parameters

Param | Value
------|-------    
**Subnet**         | 10.10.0.0/16
**Server Address** | 10.10.112.27
**Gateway**        | 10.10.0.254
**Name servers**   | 10.10.0.250, 10.10.0.252
**Search domains** | apa-canal.md
**Disk usage**     | Use entire disk and XFS filesystem

## Profile Setup

Param | Value
------|------
**Your name**   | Demius Academius
**Server name** | aria-server
**Username**    | acc-server-admin
**Password**    | shwsa3nn

## Install K8S

**From development linux workstation**

    teleport aria-server .

**In kubernetes-server**    

    sudo scripts/install.sh

**after init, copy to hold**

    kubeadm join 10.10.112.21:6443 --token bo4hv2.e4xdq52dkz85s9j --discovery-token-ca-cert-hash sha256:<hash>

**Reboot server**

    server will reboot automatically

## Post-install K8S

    sudo scripts/install-after-reboot.sh

## Deploying official Dashboard UI

[See official docs](https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/)

**Find Bearer Token for login and copy it!**

## OK

> SEE: https://aria.acc.md:30931/

## Install aria services

    sudo scripts/install-aria.sh

## Generate UUID on linux: uuidgen!
