# Terraform TOTP Provider

The TOTP provider is a utility provider, which allows for generating Time-Based One-Time Passwords (TOTP) following the RFC 6238 standard implemented by mobile apps such as Google Authenticator, Microsoft Authenticator or FreeOTP.

## Example Usage

```hcl
data "totp" "my_totp" {
  secret = var.totp_secret

  # optionally specify the algorithm, period and the number of digits
  algorithm = "SHA1"
  period    = 30
  digits    = 6
}
```

The generated TOTP is exposed as `passcode` attribute (i.e. `data.totp.my_totp.passcode`). Check the `examples` folder for additional usage examples.