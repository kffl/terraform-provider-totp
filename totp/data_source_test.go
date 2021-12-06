package totp

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testDataSourceConfig_default = `
data "totp" "totp_test" {
  secret = "%s"
}

output "passcode" {
  value = data.totp.totp_test.passcode
}
`

func TestDataSource_default(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testDataSourceConfig_default, "SOMESECRET"),
				Check: func(s *terraform.State) error {
					_, ok := s.RootModule().Resources["data.totp.totp_test"]
					if !ok {
						return fmt.Errorf("data source not found")
					}

					outputs := s.RootModule().Outputs
					passcode := outputs["passcode"].Value.(string)

					passcodeLen := len(passcode)

					if passcodeLen != 6 {
						return fmt.Errorf(`passcode length is %d; want 6`, passcodeLen)
					}

					if !regexp.MustCompile(`^[0-9]+$`).MatchString(passcode) {
						return fmt.Errorf(`passcode should only contain digits; got: %s`, passcode)
					}

					return nil
				},
			},
		},
	})
}

const testDataSourceConfig_custom = `
data "totp" "totp_test" {
  secret    = "%s"
  algorithm = "%s"
  period    = "%s"
  digits    = "%s"
}

output "passcode" {
  value = data.totp.totp_test.passcode
}
`

func TestDataSource_custom(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testDataSourceConfig_custom,
					"ANOTHERSECRET777",
					"SHA512",
					"90",
					"8",
				),
				Check: func(s *terraform.State) error {
					_, ok := s.RootModule().Resources["data.totp.totp_test"]
					if !ok {
						return fmt.Errorf("data source not found")
					}

					outputs := s.RootModule().Outputs
					passcode := outputs["passcode"].Value.(string)

					passcodeLen := len(passcode)

					if passcodeLen != 8 {
						return fmt.Errorf(`passcode length is %d; want 8`, passcodeLen)
					}

					if !regexp.MustCompile(`^[0-9]+$`).MatchString(passcode) {
						return fmt.Errorf(`passcode should only contain digits; got: %s`, passcode)
					}

					return nil
				},
			},
		},
	})
}

func TestDataSource_invalidSecret(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testDataSourceConfig_custom,
					"9NOTINBASE32",
					"SHA512",
					"90",
					"8",
				),
				ExpectError: regexp.MustCompile("Decoding of secret as base32 failed"),
			},
		},
	})
}

func TestDataSource_invalidAlgorithm(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testDataSourceConfig_custom,
					"LEGITSECRET12",
					"SHA98",
					"30",
					"6",
				),
				ExpectError: regexp.MustCompile(
					"SHA98",
				),
			},
		},
	})
}
