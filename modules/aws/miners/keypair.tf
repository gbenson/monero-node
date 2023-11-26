resource "aws_key_pair" "admin_ssh_key" {
  key_name   = "admin-ssh-key"
  public_key = var.admin_ssh_key
}
