resource "openstack_networking_secgroup_v2" "monerod" {
  name        = "monerod"
  description = "Security group for Monero nodes"
}

resource "openstack_networking_secgroup_rule_v2" "inbound_ipv4_monero_p2p" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 18080
  port_range_max    = 18080
  security_group_id = openstack_networking_secgroup_v2.monerod.id
}

resource "openstack_networking_secgroup_rule_v2" "inbound_ipv6_monero_p2p" {
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 18080
  port_range_max    = 18080
  security_group_id = openstack_networking_secgroup_v2.monerod.id
}

resource "openstack_compute_instance_v2" "p2pool_node" {
  name            = "p2pool_node"
  key_pair        = openstack_compute_keypair_v2.keypair.id
  security_groups = ["default", "monerod"]
  flavor_name     = "gp1.supersonic" # monerod needs 2 GiB RAM
  image_name      = "Ubuntu-22.04"
  user_data       = file("setup-p2pool-node.sh")
}

resource "openstack_blockstorage_volume_v3" "monerod" {
  name        = "monerod"
  description = "Monero blockchain"
  size        = 80
}

resource "openstack_compute_volume_attach_v2" "p2pool_node" {
  instance_id = openstack_compute_instance_v2.p2pool_node.id
  volume_id   = openstack_blockstorage_volume_v3.monerod.id
}

resource "dreamhost_dns_record" "p2pool_node" {
  record = "p2pool.gbenson.net"
  type   = "A"
  value  = openstack_compute_instance_v2.p2pool_node.access_ip_v4
}
