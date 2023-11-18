variable "admin_ip_prefix" {
  type        = string
  description = "CIDR from which administration may be performed"
}

variable "admin_ssh_key" {
  type        = string
  description = "SSH public key for server admin"
}