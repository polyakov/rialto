# Project Rialto
![Rialto Market](assets/Rialto-Bridge-Market.jpg)


Current repository content:
===========================
1. Apache Directory Server/LDAP

more to come...


Prerequisites:
--------------

1. Run a kubernetes cluster (http://kubernetes.io/docs/getting-started-guides/)
1. Setup helm (https://github.com/kubernetes/helm)
1. Register helm repo using
    ````
       helm repo add hspc-helm http://hspc-helm.preparedmind.net
    ````
1. Setup security groups to allow access on appropriate ports (389 for LDAP)

Install:
--------
1. List available charts using
    ```
    helm search
    ```
1. Download and edit values.yaml file, e.g http://hspc-helm.preparedmind.net/ldap.ApacheDS/values.yaml
1. Install chart using
    ```
    helm install --values values.yaml hspc-helm\ldap.ApacheDS
    ```



ApacheDS Service:
=================

Service consists of:
- Replication Controller running one ApacheDS server
- Service definition exposing this service in the cluster and externally