variable "admin_ip_prefix" {
  description = "CIDR from which administration may be performed"
  type        = string
  sensitive   = true
}

variable "admin_ssh_key" {
  description = "SSH public key for server admin"
  type        = string
  sensitive   = true
}

variable "tor_miner_config_passphrase" {
  description = "Config passphrase for tor-miner containers"
  type        = string
  sensitive   = true
}
