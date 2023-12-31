variable "admin_ip_prefix" {
  description = "CIDR from which administration may be performed"
  type        = string
}

variable "admin_ssh_key" {
  description = "SSH public key for server admin"
  type        = string
}

variable "miner_count" {
  type        = number
  description = "Number of miner instances to provision"
  default     = 0
}

variable "tor_miner_config_passphrase" {
  description = "Config passphrase for tor-miner containers"
  type        = string
}
