terraform {
  required_providers {
    totp = {
      version = "0.1.1"
      source  = "kffl/totp"
    }
    http = {
      version = "2.1.0"
      source  = "hashicorp/http"
    }
  }
}

variable "totp_secret" {
  type        = string
  description = "TOTP secret (base32)"
}

data "totp" "my_totp" {
  secret = var.totp_secret

  # optionally specify the algorithm, period and the number of digits
  algorithm = "SHA1"
  period    = 30
  digits    = 6
}

data "http" "api_call" {
  url = "https://example.com/path?token=some-token&totp=${data.totp.my_totp.passcode}"
}

output "response_body" {
  value = data.http.api_call.body
}
