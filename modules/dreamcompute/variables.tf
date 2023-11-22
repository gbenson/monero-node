variable "admin_ip_prefix" {
  description = "CIDR from which administration may be performed"
  type        = string
}

variable "admin_ssh_key" {
  description = "SSH public key for server admin"
  type        = string
}

variable "pool_count" {
  type        = number
  description = "Number of P2Pool nodes to provision"
}

variable "miner_count" {
  type        = number
  description = "Number of miner instances to provision"
}

variable "tor_miner_config_passphrase" {
  description = "Config passphrase for tor-miner containers"
  type        = string
}
