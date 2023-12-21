resource "linode_instance" "miner" {
  count  = var.miner_count
  region = "fr-par"
  image  = "linode/ubuntu22.04"
  type   = "g6-dedicated-4"

  authorized_keys = [var.admin_ssh_key]

  metadata {
    user_data = base64encode(
      templatefile("${path.module}/../setup-miner.tftpl", {
	tor_miner_config_passphrase = var.tor_miner_config_passphrase
      }))
  }

  alerts {
    cpu = 0
  }
}

resource "linode_firewall_device" "miner_firewall" {
  count       = length(linode_instance.miner)
  firewall_id = linode_firewall.ssh_admin_only.id
  entity_id   = linode_instance.miner[count.index].id
}
