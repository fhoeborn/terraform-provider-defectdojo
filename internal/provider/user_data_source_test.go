package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserDataSourceById(t *testing.T) {
	username := fmt.Sprintf("dox-test-user-%s", resource.UniqueId())
	email := fmt.Sprintf("%s@example.com", username)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing by ID
			{
				Config: testAccUserDataSourceByIdConfig(username, email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_user.test", "users.0.username", username),
					resource.TestCheckResourceAttr("data.defectdojo_user.test", "users.0.email", email),
					resource.TestCheckResourceAttrSet("data.defectdojo_user.test", "users.0.id"),
				),
			},
		},
	})
}

func TestAccUserDataSourceByUsername(t *testing.T) {
	username := fmt.Sprintf("dox-test-user-%s", resource.UniqueId())
	email := fmt.Sprintf("%s@example.com", username)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing by username
			{
				Config: testAccUserDataSourceByUsernameConfig(username, email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_user.test", "username", username),
					resource.TestCheckResourceAttr("data.defectdojo_user.test", "users.0.username", username),
					resource.TestCheckResourceAttr("data.defectdojo_user.test", "users.0.email", email),
					resource.TestCheckResourceAttrSet("data.defectdojo_user.test", "users.0.id"),
				),
			},
		},
	})
}

func TestAccUserDataSourceWithLimit(t *testing.T) {
	username1 := fmt.Sprintf("dox-test-user-%s", resource.UniqueId())
	username2 := fmt.Sprintf("dox-test-user-%s", resource.UniqueId())
	email1 := fmt.Sprintf("%s@example.com", username1)
	email2 := fmt.Sprintf("%s@example.com", username2)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing with limit
			{
				Config: testAccUserDataSourceWithLimitConfig(username1, email1, username2, email2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_user.test", "limit", "1"),
					// Should only return 1 user due to limit
					resource.TestCheckResourceAttr("data.defectdojo_user.test", "users.#", "1"),
				),
			},
		},
	})
}

func TestAccUserDataSourceWithNames(t *testing.T) {
	username := fmt.Sprintf("dox-test-user-%s", resource.UniqueId())
	email := fmt.Sprintf("%s@example.com", username)
	firstName := "Test"
	lastName := "User"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing with first and last names
			{
				Config: testAccUserDataSourceWithNamesConfig(username, email, firstName, lastName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_user.test", "users.0.username", username),
					resource.TestCheckResourceAttr("data.defectdojo_user.test", "users.0.email", email),
					resource.TestCheckResourceAttr("data.defectdojo_user.test", "users.0.first_name", firstName),
					resource.TestCheckResourceAttr("data.defectdojo_user.test", "users.0.last_name", lastName),
				),
			},
		},
	})
}

func testAccUserDataSourceByIdConfig(username, email string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_user" "test" {
  username = %[1]q
  email = %[2]q
  password = "TestPassword123!"
}

data "defectdojo_user" "test" {
  id = defectdojo_user.test.id
  depends_on = [defectdojo_user.test]
}
`, username, email)
}

func testAccUserDataSourceByUsernameConfig(username, email string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_user" "test" {
  username = %[1]q
  email = %[2]q
  password = "TestPassword123!"
}

data "defectdojo_user" "test" {
  username = %[1]q
  depends_on = [defectdojo_user.test]
}
`, username, email)
}

func testAccUserDataSourceWithLimitConfig(username1, email1, username2, email2 string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_user" "test1" {
  username = %[1]q
  email = %[2]q
  password = "TestPassword123!"
}

resource "defectdojo_user" "test2" {
  username = %[3]q
  email = %[4]q
  password = "TestPassword123!"
}

data "defectdojo_user" "test" {
  limit = 1
  depends_on = [defectdojo_user.test1, defectdojo_user.test2]
}
`, username1, email1, username2, email2)
}

func testAccUserDataSourceWithNamesConfig(username, email, firstName, lastName string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_user" "test" {
  username = %[1]q
  email = %[2]q
  password = "TestPassword123!"
  first_name = %[3]q
  last_name = %[4]q
}

data "defectdojo_user" "test" {
  username = %[1]q
  depends_on = [defectdojo_user.test]
}
`, username, email, firstName, lastName)
}
