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
    app = "testagent"
  }
  labels = merge(local.default_labels, var.labels)
}

resource "kubernetes_deployment" "testagent" {
  metadata {
    name      = "testagent-server"
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
          image = var.testagent_image
          name  = "testagent"
          image_pull_policy = var.pull_policy

          env {
            name  = "TESTAGENT_URI"
            value = var.nats_url
          }

          env {
            name  = "TESTAGENT_BUCKET"
            value = "cfm-bucket"
          }

          env {
            name  = "TESTAGENT_STREAM"
            value = "cfm-stream"
          }

          resources {
            limits   = var.resources.limits
            requests = var.resources.requests
          }

          # liveness_probe {
          #   http_get {
          #     path = "/health"
          #     port = var.testagent_port
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
          #     port = var.testagent_port
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
          #     port = var.testagent_port
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
