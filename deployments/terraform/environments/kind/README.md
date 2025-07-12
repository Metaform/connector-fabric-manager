# Development Deployment to Kind

## Docker Images

Build the `tmanager` and `pmanager` Docker images.

## Kind Setup

Install Kind and create a cluster. Set it as the default K8S context. Load the runtime images locally into Kind:

```
kind load docker-image pmanager:latest
kind load docker-image tmanager:latest
```

## Terraform Deployment

Ensure the `tmanager` and `pmanager` and `pull_policy` is set to `Never` (the default). The Terraform scripts are
configured to use the default K8S context. To deploy:

```
terraform init
terraform apply
```