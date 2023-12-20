resource "linode_instance" "miner" {
  count  = 2
  region = "fr-par"
  image = "linode/ubuntu22.04"
  type = "g6-dedicated-4"
  authorized_keys = [var.admin_ssh_key]
}
