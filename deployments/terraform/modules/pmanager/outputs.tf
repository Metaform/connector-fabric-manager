output "pmanager_service_name" {
  description = "The name of the provision manager service"
  value       = var.pmanager_service
}

output "pmanager_port" {
  description = "The port that the provision manager HTTP server listens on"
  value       = var.pmanager_port
}
