variable "miner_architectures" {
  description = "AWS architecture details"
  type        = map(any)

  default = {
    arm64 = {
      excluded_instance_types = [
	"m6g.medium",
	"m6gd.medium",
      ]
    },
    x86_64 = {
      excluded_instance_types = [
	"i3.large",
	"m7a.medium",
      ]
    }
  }
}

data "aws_ami" "amazon_linux" {
  for_each = var.miner_architectures

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

resource "aws_iam_role" "ec2_miner" {
  name = "EC2TorMiner"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
	Service = "ec2.amazonaws.com"
      }
    }
    ]
  })
}

resource "aws_iam_instance_profile" "miner" {
  name = "miner"
  role = aws_iam_role.ec2_miner.name
}

resource "aws_secretsmanager_secret" "tor_miner" {
  name = "tor-miner"
}

resource "aws_iam_policy" "tor_miner_secret_read" {
  name = "TorMinerSecret"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Action = "secretsmanager:GetSecretValue",
      Resource = aws_secretsmanager_secret.tor_miner.arn
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ec2_tor_miner_secret_read" {
  role       = aws_iam_role.ec2_miner.name
  policy_arn = aws_iam_policy.tor_miner_secret_read.arn
}

resource "aws_launch_template" "miner" {
  for_each = var.miner_architectures
  name     = "${each.key}-miner"
  image_id = data.aws_ami.amazon_linux[each.key].id

  update_default_version = true

  instance_requirements {
    vcpu_count {
      min = 1
    }
    memory_mib {
      min = "4096"
    }
    burstable_performance = "excluded"
    excluded_instance_types = each.value.excluded_instance_types
  }

  instance_market_options {
    market_type = "spot"
  }

  iam_instance_profile {
    name = aws_iam_instance_profile.miner.name
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
  associate_public_ip_address = true

  instance_market_options {
    spot_options {
      max_price = 0.7
    }
  }
}
