package totp

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func dataSourceTotp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTotpRead,

		Schema: map[string]*schema.Schema{
			"secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "SHA1",
				ValidateFunc: validation.StringInSlice(
					[]string{"SHA1", "SHA256", "SHA512", "MD5"},
					false,
				),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"digits": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  6,
			},
			"passcode": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func getAlgorithm(algorithm string) otp.Algorithm {
	switch algorithm {
	case "SHA1":
		return otp.AlgorithmSHA1
	case "SHA256":
		return otp.AlgorithmSHA256
	case "SHA512":
		return otp.AlgorithmSHA512
	case "MD5":
		return otp.AlgorithmMD5
	}
	panic("this should never happen, invalid algorithm name")
}

func dataSourceTotpRead(
	ctx context.Context,
	d *schema.ResourceData,
	meta interface{},
) (diags diag.Diagnostics) {

	secret := d.Get("secret").(string)
	timestamp := time.Now()

	passcode, err := totp.GenerateCodeCustom(secret, timestamp, totp.ValidateOpts{
		Period:    uint(d.Get("period").(int)),
		Digits:    otp.Digits(d.Get("digits").(int)),
		Algorithm: getAlgorithm(d.Get("algorithm").(string)),
	})

	if err != nil {
		return append(diags, diag.Errorf("Error generating TOTP passcode: %s", err)...)
	}

	d.Set("passcode", passcode)
	d.SetId(strconv.FormatInt(timestamp.Unix(), 10))

	return diags
}
