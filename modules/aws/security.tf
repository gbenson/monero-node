# Default AWS Virtual Private Cloud
# (the default virtual network segment created by AWS)
resource "aws_default_vpc" "default" {
}

# Security group to allow SSH from my CIDR only
# All outgoing traffic is allowed
resource "aws_security_group" "ssh_admin" {
  name        = "ssh-admin"
  description = "Allows inbound SSH from my CIDR"
  vpc_id      = aws_default_vpc.default.id
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.admin_ip_prefix]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
