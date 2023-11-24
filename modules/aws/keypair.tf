resource "aws_key_pair" "keypair" {
  key_name   = "admin_ssh_key"
  public_key = var.admin_ssh_key
}
