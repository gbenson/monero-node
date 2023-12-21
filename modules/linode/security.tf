# Allow SSH from my CIDR only
resource "linode_firewall" "ssh_admin_only" {
  label = "ssh-admin-only"

  inbound {
    label    = "admin-ssh"
    action = "ACCEPT"
    protocol  = "TCP"
    ports     = "22"
    ipv4 = [var.admin_ip_prefix]
  }

  inbound_policy = "DROP"
  outbound_policy = "ACCEPT"
}
