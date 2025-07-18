# Development Deployment to Kind

## Docker Images

Build the `tmanager` and `pmanager` Docker images.

## Kind Setup

Install Kind and create a cluster. Set it as the default K8S context. 

Services can be exposed on the host by configuring the cluster with port mappings:

```yaml
# kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    extraPortMappings:
      - containerPort: 30081    # PManager
        hostPort: 8181
        protocol: TCP
      - containerPort: 30082    # TManager
        hostPort: 8282
        protocol: TCP
```

Load the runtime images locally into Kind:

```
kind load docker-image pmanager:latest
kind load docker-image tmanager:latest
```

## Terraform Deployment

Ensure the `pull_policy` for `tmanager` and `pmanager` is set to `Never` (the default). The Terraform scripts are
configured to use the default K8S context. To deploy:

```
terraform init
terraform apply
```