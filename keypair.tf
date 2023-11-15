resource "openstack_compute_keypair_v2" "keypair" {
  name       = "publickey"
  public_key = var.admin_public_key
}
