# Installing Kuberenetes on Ubuntu

[see official docs](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/)

## Common parameters

Param | Value
------|-------
**VM-Name**      | kubernetes-srv
**disk**         | 400Gb
**RAM**          | 16Gb
**Bridge**       | bridge0

## Linux Installer parameters

Param | Value | Comments
------|-------|---------
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
**Server name** | kubernetes-server
**Username**    | acc-server-admin
**Password**    | shwsa3nn

## After Install

**From development linux workstation**

    teleport aria-server .

**In kubernetes-server**    

    sudo scripts/install-k8s.sh

**after init, copy to hold**

    kubeadm join 10.10.112.21:6443 --token bo4hv2.e4xdq52dkz85s9j --discovery-token-ca-cert-hash sha256:<hash>

**Enable for acc-server-admin user**

    mkdir -p $HOME/.kube
    sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
    sudo chown $(id -u):$(id -g) $HOME/.kube/config

**Installing a pod network add-on**
    
    kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/a70459be0084506e4ec919aa1c114638878db11b/Documentation/kube-flannel.yml

> Reference to official docs in each installation of Kubernetes for changes

**Master Isolation for single node installation**

    kubectl taint nodes --all node-role.kubernetes.io/master-

**Reboot server**

    sudo reboot

## Postinstallations

    sudo scripts/post-install.sh

## Deploying official Dashboard UI

[See official docs](https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/)

**Find Bearer Token for login**

    kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep dashboard-user | awk '{print $1}')

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

## Запускаем docker-registry

    sudo docker login 10.10.112.27:5000 (user: kube-registry-user, pw: yf-ujhirt-cbltk-rjhjkm на-горшке-сидел-король)
