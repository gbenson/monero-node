variable "admin_ip_prefix" {
  description = "CIDR from which administration may be performed"
  type        = string
}

variable "admin_ssh_key" {
  description = "SSH public key for server admin"
  type        = string
}
