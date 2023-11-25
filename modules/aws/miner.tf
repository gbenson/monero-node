data "aws_ami" "amazon_linux_x86_64" {
  most_recent = true
  owners      = ["amazon"]
  filter {
    name   = "architecture"
    values = ["x86_64"]
  }
  filter {
    name   = "name"
    values = ["al2023-ami-minimal-2023*"]
  }
}

resource "aws_instance" "miner" {
  count = 0
  ami           = data.aws_ami.amazon_linux_x86_64.id
  instance_type = "c6i.large"
  instance_market_options {
    market_type = "spot"
    spot_options {
      max_price = 0.7
    }
  }

  vpc_security_group_ids      = [aws_security_group.ssh_admin.id]
  key_name                    = aws_key_pair.admin_ssh_key.id
  associate_public_ip_address = "true"

  user_data = file("${path.module}/../setup-miner.sh")
}
