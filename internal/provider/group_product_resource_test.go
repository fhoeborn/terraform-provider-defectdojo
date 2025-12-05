package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGroupProductResource(t *testing.T) {
	groupName := fmt.Sprintf("dox-test-group-%s", resource.UniqueId())
	productName := fmt.Sprintf("dox-test-product-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGroupProductResourceConfig(groupName, productName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("defectdojo_group_product.test", "group_id", "defectdojo_group.test", "id"),
					resource.TestCheckResourceAttrPair("defectdojo_group_product.test", "product_id", "defectdojo_product.test", "id"),
					resource.TestCheckResourceAttr("defectdojo_group_product.test", "role_id", "1"),
					resource.TestCheckResourceAttrSet("defectdojo_group_product.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_group_product.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing (change role)
			{
				Config: testAccGroupProductResourceConfigUpdated(groupName, productName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("defectdojo_group_product.test", "group_id", "defectdojo_group.test", "id"),
					resource.TestCheckResourceAttrPair("defectdojo_group_product.test", "product_id", "defectdojo_product.test", "id"),
					resource.TestCheckResourceAttr("defectdojo_group_product.test", "role_id", "2"),
				),
			},
		},
	})
}

func TestAccGroupProductResourceMultipleGroups(t *testing.T) {
	groupName1 := fmt.Sprintf("dox-test-group-%s", resource.UniqueId())
	groupName2 := fmt.Sprintf("dox-test-group-%s", resource.UniqueId())
	productName := fmt.Sprintf("dox-test-product-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with multiple groups
			{
				Config: testAccGroupProductResourceMultipleGroupsConfig(groupName1, groupName2, productName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("defectdojo_group_product.test1", "group_id", "defectdojo_group.test1", "id"),
					resource.TestCheckResourceAttrPair("defectdojo_group_product.test1", "product_id", "defectdojo_product.test", "id"),
					resource.TestCheckResourceAttrPair("defectdojo_group_product.test2", "group_id", "defectdojo_group.test2", "id"),
					resource.TestCheckResourceAttrPair("defectdojo_group_product.test2", "product_id", "defectdojo_product.test", "id"),
				),
			},
		},
	})
}

func testAccGroupProductResourceConfig(groupName, productName string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_group" "test" {
  name = %[1]q
  description = "test description"
}

resource "defectdojo_product" "test" {
  name = %[2]q
  description = "test product"
  product_type_id = 1
}

resource "defectdojo_group_product" "test" {
  group_id = defectdojo_group.test.id
  product_id = defectdojo_product.test.id
  role_id = "1"
}
`, groupName, productName)
}

func testAccGroupProductResourceConfigUpdated(groupName, productName string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_group" "test" {
  name = %[1]q
  description = "test description"
}

resource "defectdojo_product" "test" {
  name = %[2]q
  description = "test product"
  product_type_id = 1
}

resource "defectdojo_group_product" "test" {
  group_id = defectdojo_group.test.id
  product_id = defectdojo_product.test.id
  role_id = "2"
}
`, groupName, productName)
}

func testAccGroupProductResourceMultipleGroupsConfig(groupName1, groupName2, productName string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_group" "test1" {
  name = %[1]q
  description = "test description 1"
}

resource "defectdojo_group" "test2" {
  name = %[2]q
  description = "test description 2"
}

resource "defectdojo_product" "test" {
  name = %[3]q
  description = "test product"
  product_type_id = 1
}

resource "defectdojo_group_product" "test1" {
  group_id = defectdojo_group.test1.id
  product_id = defectdojo_product.test.id
  role_id = "1"
}

resource "defectdojo_group_product" "test2" {
  group_id = defectdojo_group.test2.id
  product_id = defectdojo_product.test.id
  role_id = "2"
}
`, groupName1, groupName2, productName)
}
