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
  description = "Kubernetes namespace for deployment"
  type        = string
  default     = "default"
}

variable "replicas" {
  description = "Number of testagent replicas"
  type        = number
  default     = 1
}

variable "testagent_image" {
  description = "Docker image"
  type        = string
  default     = "testagent:latest"
}

variable "pull_policy" {
  description = "Docker image pull policy"
  type        = string
}

variable "nats_url" {
  description = "NATS URL"
  type        = string
  default     = "nats://nats-service:4222"
}

variable "log_level" {
  description = "Log level"
  type        = string
  default     = "info"
}

variable "resources" {
  description = "Resource limits and requests"
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
      cpu    = "250m"
      memory = "256Mi"
    }
  }
}

variable "labels" {
  description = "Additional labels to apply to all resources"
  type = map(string)
  default = {}
}

