terraform {
  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 1.53.0"
    }
    dreamhost = {
      source  = "adamantal/dreamhost"
      version = "0.3.2"
    }
  }
}
