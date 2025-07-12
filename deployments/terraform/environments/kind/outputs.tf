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
  value       = module.nats.nats_internal_url
  description = "NATS URL for internal cluster access"
}

output "nats_external_url" {
  value       = module.nats.nats_external_url
  description = "NATS URL for external access"
}

output "monitoring_urls" {
  value = {
    internal = module.nats.monitoring_internal_url
    external = module.nats.monitoring_external_url
  }
  description = "NATS monitoring URLs"
}
