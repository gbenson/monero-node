terraform {
  required_version = ">= 0.14.0"

  backend "local" {
    path = ".tfstate/terraform.tfstate"
  }
}

module "aws_us_east_1" {
  source = "./modules/aws/listener"
}

module "aws_us_east_2" {
  source = "./modules/aws/miners"

  admin_ip_prefix = var.admin_ip_prefix
  admin_ssh_key   = var.admin_ssh_key
}

module "linode" {
  source = "./modules/linode"

  admin_ip_prefix = var.admin_ip_prefix
  admin_ssh_key   = var.admin_ssh_key

  miner_count = 1

  tor_miner_config_passphrase = var.tor_miner_config_passphrase
}
