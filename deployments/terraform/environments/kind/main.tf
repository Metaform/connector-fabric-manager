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

provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = var.kubeconfig_path
}

module "nats" {
  source = "../../modules/nats"

  namespace              = var.namespace
  nats_image             = "nats:2.10-alpine"
  replicas               = 1
  jetstream_storage_size = "256Mi"
  client_nodeport        = 30422
  monitor_nodeport       = 30822
  enable_nodeport        = true

  labels = {
    environment = "development"
    managed-by  = "terraform"
  }
}

module "pmanager" {
  source = "../../modules/pmanager"

  pmanager_image  = "pmanager:latest"
  pull_policy = "Never"  # pull locally from Docker
  enable_nodeport = true

  depends_on = [module.nats]
}

module "tmanager" {
  source = "../../modules/tmanager"

  tmanager_image = "tmanager:latest"
  pull_policy = "Never"  # pull locally from Docker
  enable_nodeport = true

  depends_on = [module.nats]
}