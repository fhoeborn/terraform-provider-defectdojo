package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGroupMemberResource(t *testing.T) {
	groupName := fmt.Sprintf("dox-test-group-%s", resource.UniqueId())
	username := fmt.Sprintf("dox-test-user-%s", resource.UniqueId())
	email := fmt.Sprintf("%s@example.com", username)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGroupMemberResourceConfig(groupName, username, email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("defectdojo_group_member.test", "group_id", "defectdojo_group.test", "id"),
					resource.TestCheckResourceAttrPair("defectdojo_group_member.test", "user_id", "defectdojo_user.test", "id"),
					resource.TestCheckResourceAttr("defectdojo_group_member.test", "role_id", "1"),
					resource.TestCheckResourceAttrSet("defectdojo_group_member.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_group_member.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing (change role)
			{
				Config: testAccGroupMemberResourceConfigUpdated(groupName, username, email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("defectdojo_group_member.test", "group_id", "defectdojo_group.test", "id"),
					resource.TestCheckResourceAttrPair("defectdojo_group_member.test", "user_id", "defectdojo_user.test", "id"),
					resource.TestCheckResourceAttr("defectdojo_group_member.test", "role_id", "2"),
				),
			},
		},
	})
}

func TestAccGroupMemberResourceMultipleMembers(t *testing.T) {
	groupName := fmt.Sprintf("dox-test-group-%s", resource.UniqueId())
	username1 := fmt.Sprintf("dox-test-user-%s", resource.UniqueId())
	username2 := fmt.Sprintf("dox-test-user-%s", resource.UniqueId())
	email1 := fmt.Sprintf("%s@example.com", username1)
	email2 := fmt.Sprintf("%s@example.com", username2)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with multiple members
			{
				Config: testAccGroupMemberResourceMultipleMembersConfig(groupName, username1, email1, username2, email2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("defectdojo_group_member.test1", "group_id", "defectdojo_group.test", "id"),
					resource.TestCheckResourceAttrPair("defectdojo_group_member.test1", "user_id", "defectdojo_user.test1", "id"),
					resource.TestCheckResourceAttrPair("defectdojo_group_member.test2", "group_id", "defectdojo_group.test", "id"),
					resource.TestCheckResourceAttrPair("defectdojo_group_member.test2", "user_id", "defectdojo_user.test2", "id"),
				),
			},
		},
	})
}

func testAccGroupMemberResourceConfig(groupName, username, email string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_group" "test" {
  name = %[1]q
  description = "test description"
}

resource "defectdojo_user" "test" {
  username = %[2]q
  email = %[3]q
  password = "TestPassword123!"
}

resource "defectdojo_group_member" "test" {
  group_id = defectdojo_group.test.id
  user_id = defectdojo_user.test.id
  role_id = "1"
}
`, groupName, username, email)
}

func testAccGroupMemberResourceConfigUpdated(groupName, username, email string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_group" "test" {
  name = %[1]q
  description = "test description"
}

resource "defectdojo_user" "test" {
  username = %[2]q
  email = %[3]q
  password = "TestPassword123!"
}

resource "defectdojo_group_member" "test" {
  group_id = defectdojo_group.test.id
  user_id = defectdojo_user.test.id
  role_id = "2"
}
`, groupName, username, email)
}

func testAccGroupMemberResourceMultipleMembersConfig(groupName, username1, email1, username2, email2 string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_group" "test" {
  name = %[1]q
  description = "test description"
}

resource "defectdojo_user" "test1" {
  username = %[2]q
  email = %[3]q
  password = "TestPassword123!"
}

resource "defectdojo_user" "test2" {
  username = %[4]q
  email = %[5]q
  password = "TestPassword123!"
}

resource "defectdojo_group_member" "test1" {
  group_id = defectdojo_group.test.id
  user_id = defectdojo_user.test1.id
  role_id = "1"
}

resource "defectdojo_group_member" "test2" {
  group_id = defectdojo_group.test.id
  user_id = defectdojo_user.test2.id
  role_id = "2"
}
`, groupName, username1, email1, username2, email2)
}
