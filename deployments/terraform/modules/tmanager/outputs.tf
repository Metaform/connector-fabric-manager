output "tmanager_service_name" {
  description = "The name of the tenant manager service"
  value       = var.tmanager_service
}

output "tmanager_port" {
  description = "The port that the tenant manager HTTP server listens on"
  value       = var.tmanager_port
}
