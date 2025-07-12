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

variable "namespace" {
  description = "Kubernetes namespace for NATS deployment"
  type        = string
  default     = "default"
}

variable "nats_image" {
  description = "NATS container image"
  type        = string
  default     = "nats:2.10-alpine"
}

variable "replicas" {
  description = "Number of NATS replicas"
  type        = number
  default     = 1
}

variable "resources" {
  description = "Resource limits and requests for NATS container"
  type = object({
    limits = object({
      cpu    = string
      memory = string
    })
    requests = object({
      cpu    = string
      memory = string
    })
  })
  default = {
    limits = {
      cpu    = "500m"
      memory = "512Mi"
    }
    requests = {
      cpu    = "100m"
      memory = "128Mi"
    }
  }
}

variable "jetstream_storage_size" {
  description = "Size limit for JetStream storage"
  type        = string
  default     = "256Mi"
}

variable "client_nodeport" {
  description = "NodePort for NATS client connections"
  type        = number
  default     = 30422
}

variable "monitor_nodeport" {
  description = "NodePort for NATS monitoring"
  type        = number
  default     = 30822
}

variable "enable_nodeport" {
  description = "Whether to create NodePort service for external access"
  type        = bool
  default     = true
}

variable "labels" {
  description = "Additional labels to apply to resources"
  type        = map(string)
  default     = {}
}