data "aws_ami" "amazon_linux" {
  for_each = toset(["x86_64", "arm64"])

  most_recent = true
  owners      = ["amazon"]
  filter {
    name   = "architecture"
    values = [each.key]
  }
  filter {
    name   = "name"
    values = ["al2023-ami-minimal-2023*"]
  }
}

resource "aws_launch_template" "miner" {
  for_each = data.aws_ami.amazon_linux
  name     = "${each.value.architecture}-miner"
  image_id = each.value.id

  instance_requirements {
    vcpu_count {
      min = 1
    }
    memory_mib {
      min = "4096"
    }
    burstable_performance = "excluded"
  }

  instance_market_options {
    market_type = "spot"
  }

  vpc_security_group_ids = [aws_security_group.ssh_admin.id]
  key_name               = aws_key_pair.admin_ssh_key.id

  user_data = filebase64("${path.module}/../../setup-miner.sh")
}

resource "aws_instance" "miner" {
  count = 0
  launch_template  {
    id = aws_launch_template.miner["x86_64"].id
  }
  associate_public_ip_address = "true"

  instance_market_options {
    spot_options {
      max_price = 0.7
    }
  }
}
