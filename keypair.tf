variable "public_key" {
  type        = string
  description = ""
}

resource "openstack_compute_keypair_v2" "keypair" {
  name       = "publickey"
  public_key = var.public_key
}
