package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	username := fmt.Sprintf("dox-test-user-%s", resource.UniqueId())
	email := fmt.Sprintf("%s@example.com", username)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserResourceConfig(username, email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_user.test", "username", username),
					resource.TestCheckResourceAttr("defectdojo_user.test", "email", email),
					resource.TestCheckResourceAttrSet("defectdojo_user.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "defectdojo_user.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Update and Read testing
			{
				Config: testAccUserResourceConfigUpdated(username, email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_user.test", "username", username),
					resource.TestCheckResourceAttr("defectdojo_user.test", "email", email),
					resource.TestCheckResourceAttr("defectdojo_user.test", "first_name", "Updated"),
					resource.TestCheckResourceAttr("defectdojo_user.test", "last_name", "Name"),
				),
			},
		},
	})
}

func TestAccUserResourceWithNames(t *testing.T) {
	username := fmt.Sprintf("dox-test-user-%s", resource.UniqueId())
	email := fmt.Sprintf("%s@example.com", username)
	firstName := "Test"
	lastName := "User"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with names
			{
				Config: testAccUserResourceWithNamesConfig(username, email, firstName, lastName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_user.test", "username", username),
					resource.TestCheckResourceAttr("defectdojo_user.test", "email", email),
					resource.TestCheckResourceAttr("defectdojo_user.test", "first_name", firstName),
					resource.TestCheckResourceAttr("defectdojo_user.test", "last_name", lastName),
					resource.TestCheckResourceAttrSet("defectdojo_user.test", "id"),
				),
			},
		},
	})
}

func testAccUserResourceConfig(username, email string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_user" "test" {
  username = %[1]q
  email = %[2]q
  password = "TestPassword123!"
}
`, username, email)
}

func testAccUserResourceConfigUpdated(username, email string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_user" "test" {
  username = %[1]q
  email = %[2]q
  password = "TestPassword123!"
  first_name = "Updated"
  last_name = "Name"
}
`, username, email)
}

func testAccUserResourceWithNamesConfig(username, email, firstName, lastName string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_user" "test" {
  username = %[1]q
  email = %[2]q
  password = "TestPassword123!"
  first_name = %[3]q
  last_name = %[4]q
}
`, username, email, firstName, lastName)
}
