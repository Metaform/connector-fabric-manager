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

variable "kubeconfig_path" {
  description = "Path to kubeconfig file or context name"
  type        = string
  default     = null  # Current context if not specified
}
