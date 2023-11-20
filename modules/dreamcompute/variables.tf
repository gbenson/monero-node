variable "admin_ip_prefix" {
  type        = string
  description = "CIDR from which administration may be performed"
}

variable "admin_ssh_key" {
  type        = string
  description = "SSH public key for server admin"
}

variable "pool_count" {
  description = "Number of P2Pool nodes to provision"
  type        = number
}

variable "miner_count" {
  description = "Number of miner instances to provision"
  type        = number
}

variable "tor_miner_config_passphrase" {
  type        = string
  description = "Config passphrase for tor-miner containers"
}
