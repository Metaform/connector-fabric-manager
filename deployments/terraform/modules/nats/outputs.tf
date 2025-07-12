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

output "nats_internal_url" {
  value       = "nats://${kubernetes_service.nats.metadata[0].name}.${var.namespace}.svc.cluster.local:4222"
  description = "NATS URL for internal cluster access"
}

output "nats_external_url" {
  value       = var.enable_nodeport ? "nats://localhost:${var.client_nodeport}" : null
  description = "NATS URL for external access via Kind"
}

output "monitoring_internal_url" {
  value       = "http://${kubernetes_service.nats.metadata[0].name}.${var.namespace}.svc.cluster.local:8222"
  description = "NATS monitoring URL for internal cluster access"
}

output "monitoring_external_url" {
  value       = var.enable_nodeport ? "http://localhost:${var.monitor_nodeport}" : null
  description = "NATS monitoring URL for external access via Kind"
}

output "service_name" {
  value       = kubernetes_service.nats.metadata[0].name
  description = "Name of the NATS service"
}

output "deployment_name" {
  value       = kubernetes_deployment.nats.metadata[0].name
  description = "Name of the NATS deployment"
}

output "namespace" {
  value       = var.namespace
  description = "Namespace where NATS is deployed"
}