resource "openstack_networking_secgroup_v2" "xmrig" {
  name        = "xmrig"
  description = "Security group for XMRig nodes"
}

resource "openstack_compute_instance_v2" "miner" {
  count           = var.miner_count
  name            = "miner${count.index + 1}"
  key_pair        = openstack_compute_keypair_v2.keypair.id
  security_groups = ["default", "xmrig"]
  flavor_name     = "gp1.lightspeed"
  image_name      = "Ubuntu-22.04"

  personality {
    file    = "/etc/tor-miner/config_passphrase"
    content = "${var.tor_miner_config_passphrase}\n"
  }

  user_data = file("${path.module}/../setup-miner.sh")
}

resource "dreamhost_dns_record" "miner" {
  count  = var.miner_count
  record = "miner${count.index + 1}.gbenson.net"
  type   = "A"
  value  = openstack_compute_instance_v2.miner[count.index].access_ip_v4
}
