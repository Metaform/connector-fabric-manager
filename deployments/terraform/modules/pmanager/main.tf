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
    app = "pmanager"
  }
  labels = merge(local.default_labels, var.labels)
}

resource "kubernetes_deployment" "pmanager" {
  metadata {
    name      = "pmanager-server"
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
          image = var.pmanager_image
          name  = "pmanager"
          image_pull_policy = var.pull_policy

          # port {
          #   container_port = var.pmanager_port
          #   name           = "http"
          # }
          #
          # port {
          #   container_port = var.metrics_port
          #   name           = "metrics"
          # }

          env {
            name  = "PM_URI"
            value = var.nats_url
          }

          env {
            name  = "PM_BUCKET"
            value = "cfm-bucket"
          }

          env {
            name  = "PM_STREAM"
            value = "cfm-stream"
          }

          # env {
          #   name  = "PORT"
          #   value = tostring(var.pmanager_port)
          # }

          # env {
          #   name  = "LOG_LEVEL"
          #   value = var.log_level
          # }

          # env {
          #   name  = "METRICS_PORT"
          #   value = tostring(var.metrics_port)
          # }

          resources {
            limits   = var.resources.limits
            requests = var.resources.requests
          }

          # liveness_probe {
          #   http_get {
          #     path = "/health"
          #     port = var.pmanager_port
          #   }
          #   initial_delay_seconds = 30
          #   period_seconds        = 10
          #   timeout_seconds       = 5
          #   failure_threshold     = 3
          # }
          #
          # readiness_probe {
          #   http_get {
          #     path = "/ready"
          #     port = var.pmanager_port
          #   }
          #   initial_delay_seconds = 5
          #   period_seconds        = 5
          #   timeout_seconds       = 3
          #   failure_threshold     = 3
          # }
          #
          # startup_probe {
          #   http_get {
          #     path = "/health"
          #     port = var.pmanager_port
          #   }
          #   initial_delay_seconds = 10
          #   period_seconds        = 10
          #   timeout_seconds       = 3
          #   failure_threshold     = 10
          # }
        }
      }
    }
  }
}

# Pmanager Service
resource "kubernetes_service" "pmanager" {
  metadata {
    name      = "pmanager-service"
    namespace = var.namespace
    labels    = local.labels
  }

  spec {
    selector = local.default_labels

    port {
      name        = "http"
      port        = var.pmanager_port
      target_port = var.pmanager_port
    }

    port {
      name        = "metrics"
      port        = var.metrics_port
      target_port = var.metrics_port
    }

    type = "ClusterIP"
  }
}

# NodePort service for external access
resource "kubernetes_service" "pmanager_nodeport" {
  count = var.enable_nodeport ? 1 : 0

  metadata {
    name      = "pmanager-nodeport"
    namespace = var.namespace
    labels    = local.labels
  }

  spec {
    selector = local.default_labels

    port {
      name        = "http"
      port        = var.pmanager_port
      target_port = var.pmanager_port
      node_port   = var.pmanager_nodeport
    }

    port {
      name        = "metrics"
      port        = var.metrics_port
      target_port = var.metrics_port
      node_port   = var.metrics_nodeport
    }

    type = "NodePort"
  }
}