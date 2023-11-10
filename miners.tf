resource "openstack_networking_secgroup_v2" "xmrig" {
  name        = "xmrig"
  description = "Security group for XMRig nodes"
}

resource "openstack_compute_instance_v2" "miner1" {
  name            = "miner1"
  key_pair        = openstack_compute_keypair_v2.keypair.id
  security_groups = ["default", "xmrig"]
  flavor_name     = "gp1.lightspeed"
  image_name      = "Ubuntu-22.04"
  user_data       = file("setup-miner.sh")
}

resource "dreamhost_dns_record" "miner1" {
  record = "miner1.gbenson.net"
  type   = "A"
  value  = openstack_compute_instance_v2.miner1.access_ip_v4
}
