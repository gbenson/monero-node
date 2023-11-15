resource "openstack_networking_secgroup_v2" "default" {
  name        = "default"
  description = "Default security group"
}

# Inbound ICMP and IPv6-ICMP is allowed
resource "openstack_networking_secgroup_rule_v2" "inbound_ipv4_icmp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  security_group_id = openstack_networking_secgroup_v2.default.id
}
resource "openstack_networking_secgroup_rule_v2" "inbound_ipv6_icmp" {
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "ipv6-icmp"
  security_group_id = openstack_networking_secgroup_v2.default.id
}

# All outgoing traffic is allowed
resource "openstack_networking_secgroup_rule_v2" "outbound_ipv4_all" {
  direction         = "egress"
  ethertype         = "IPv4"
  security_group_id = openstack_networking_secgroup_v2.default.id
}
resource "openstack_networking_secgroup_rule_v2" "outbound_ipv6_all" {
  direction         = "egress"
  ethertype         = "IPv6"
  security_group_id = openstack_networking_secgroup_v2.default.id
}

# Inbound SSH is allowed from the admin CIDR only
resource "openstack_networking_secgroup_rule_v2" "inbound_ipv4_ssh" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = var.admin_ip_prefix
  security_group_id = openstack_networking_secgroup_v2.default.id
}
