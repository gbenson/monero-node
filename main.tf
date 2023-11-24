terraform {
  required_version = ">= 0.14.0"

  backend "local" {
    path = ".tfstate/terraform.tfstate"
  }
}

module "aws" {
  source = "./modules/aws"

  admin_ssh_key = var.admin_ssh_key
}

module "dreamcompute" {
  source = "./modules/dreamcompute"

  admin_ip_prefix = var.admin_ip_prefix
  admin_ssh_key   = var.admin_ssh_key

  pool_count  = 0
  miner_count = 1

  tor_miner_config_passphrase = var.tor_miner_config_passphrase
}
