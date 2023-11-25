resource "openstack_compute_keypair_v2" "admin_ssh_key" {
  name       = "admin-ssh-key"
  public_key = var.admin_ssh_key
}
