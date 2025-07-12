#  Copyright (c) 2025 Metaform Systems, Inc
#
#  This program and the accompanying materials are made available under the
#  terms of the Apache License, Version 2.0 which is available at
#  https://www.apache.org/licenses/LICENSE-2.0
#
#  SPDX-License-Identifier: Apache-2.0
#
#  Contributors:
#       Metaform Systems, Inc. - initial API and implementation
#

locals {
  default_labels = {
    app = "nats"
  }
  labels = merge(local.default_labels, var.labels)
}

# NATS Deployment
resource "kubernetes_deployment" "nats" {
  metadata {
    name      = "nats-server"
    namespace = var.namespace
    labels    = local.labels
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = local.default_labels
    }

    template {
      metadata {
        labels = local.labels
      }

      spec {
        container {
          image = var.nats_image
          name  = "nats"

          args = [
            "--jetstream",
            "--store_dir=/tmp/jetstream",
            "--port=4222",
            "--http_port=8222"
          ]

          port {
            container_port = 4222
            name           = "client"
          }

          port {
            container_port = 8222
            name           = "monitor"
          }

          volume_mount {
            name       = "jetstream-storage"
            mount_path = "/tmp/jetstream"
          }

          resources {
            limits   = var.resources.limits
            requests = var.resources.requests
          }

          liveness_probe {
            http_get {
              path = "/healthz"
              port = 8222
            }
            initial_delay_seconds = 30
            period_seconds        = 10
          }

          readiness_probe {
            http_get {
              path = "/healthz"
              port = 8222
            }
            initial_delay_seconds = 5
            period_seconds        = 5
          }
        }

        volume {
          name = "jetstream-storage"
          empty_dir {
            medium     = "Memory"
            size_limit = var.jetstream_storage_size
          }
        }
      }
    }
  }
}

# NATS Service
resource "kubernetes_service" "nats" {
  metadata {
    name      = "nats-service"
    namespace = var.namespace
    labels    = local.labels
  }

  spec {
    selector = local.default_labels

    port {
      name        = "client"
      port        = 4222
      target_port = 4222
    }

    port {
      name        = "monitor"
      port        = 8222
      target_port = 8222
    }

    type = "ClusterIP"
  }
}

# NodePort service for external access
resource "kubernetes_service" "nats_nodeport" {
  count = var.enable_nodeport ? 1 : 0

  metadata {
    name      = "nats-nodeport"
    namespace = var.namespace
    labels    = local.labels
  }

  spec {
    selector = local.default_labels

    port {
      name        = "client"
      port        = 4222
      target_port = 4222
      node_port   = var.client_nodeport
    }

    port {
      name        = "monitor"
      port        = 8222
      target_port = 8222
      node_port   = var.monitor_nodeport
    }

    type = "NodePort"
  }
}