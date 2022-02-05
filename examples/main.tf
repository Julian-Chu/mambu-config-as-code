terraform {
  required_providers {
    mambu = {
      version = "0.3"
      source  = "hashicorp.com/edu/mambu"
    }
  }
}

provider "mambu"{
  mambu_base_url = var.mambu_base_url
  mambu_apikey = var.mambu_apikey
}

data "mambu_custom_fields" "cfs"{
}

output "cfs" {
  value = data.mambu_custom_fields.cfs
}

variable "mambu_base_url" {
  type    = string
  default = ""
}

variable "mambu_apikey" {
  type    = string
  default = ""
}
