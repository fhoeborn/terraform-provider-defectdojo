package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGroupDataSourceById(t *testing.T) {
	name := fmt.Sprintf("dox-test-group-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing by ID
			{
				Config: testAccGroupDataSourceByIdConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_group.test", "groups.0.name", name),
					resource.TestCheckResourceAttr("data.defectdojo_group.test", "groups.0.description", "test description"),
					resource.TestCheckResourceAttrSet("data.defectdojo_group.test", "groups.0.id"),
				),
			},
		},
	})
}

func TestAccGroupDataSourceByName(t *testing.T) {
	name := fmt.Sprintf("dox-test-group-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing by name
			{
				Config: testAccGroupDataSourceByNameConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_group.test", "name", name),
					resource.TestCheckResourceAttr("data.defectdojo_group.test", "groups.0.name", name),
					resource.TestCheckResourceAttr("data.defectdojo_group.test", "groups.0.description", "test description"),
					resource.TestCheckResourceAttrSet("data.defectdojo_group.test", "groups.0.id"),
				),
			},
		},
	})
}

func TestAccGroupDataSourceWithLimit(t *testing.T) {
	name1 := fmt.Sprintf("dox-test-group-%s", resource.UniqueId())
	name2 := fmt.Sprintf("dox-test-group-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing with limit
			{
				Config: testAccGroupDataSourceWithLimitConfig(name1, name2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_group.test", "limit", "1"),
					// Should only return 1 group due to limit
					resource.TestCheckResourceAttr("data.defectdojo_group.test", "groups.#", "1"),
				),
			},
		},
	})
}

func TestAccGroupDataSourceMultipleGroups(t *testing.T) {
	name1 := fmt.Sprintf("dox-test-group-%s", resource.UniqueId())
	name2 := fmt.Sprintf("dox-test-group-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing multiple groups
			{
				Config: testAccGroupDataSourceMultipleGroupsConfig(name1, name2),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Should return at least 2 groups
					resource.TestCheckResourceAttr("data.defectdojo_group.test", "groups.#", "2"),
					resource.TestCheckResourceAttrSet("data.defectdojo_group.test", "groups.0.id"),
					resource.TestCheckResourceAttrSet("data.defectdojo_group.test", "groups.0.name"),
					resource.TestCheckResourceAttrSet("data.defectdojo_group.test", "groups.1.id"),
					resource.TestCheckResourceAttrSet("data.defectdojo_group.test", "groups.1.name"),
				),
			},
		},
	})
}

func testAccGroupDataSourceByIdConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_group" "test" {
  name = %[1]q
  description = "test description"
}

data "defectdojo_group" "test" {
  id = defectdojo_group.test.id
  depends_on = [defectdojo_group.test]
}
`, name)
}

func testAccGroupDataSourceByNameConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_group" "test" {
  name = %[1]q
  description = "test description"
}

data "defectdojo_group" "test" {
  name = %[1]q
  depends_on = [defectdojo_group.test]
}
`, name)
}

func testAccGroupDataSourceWithLimitConfig(name1, name2 string) string {
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

data "defectdojo_group" "test" {
  limit = 1
  depends_on = [defectdojo_group.test1, defectdojo_group.test2]
}
`, name1, name2)
}

func testAccGroupDataSourceMultipleGroupsConfig(name1, name2 string) string {
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

data "defectdojo_group" "test" {
  limit = 2
  depends_on = [defectdojo_group.test1, defectdojo_group.test2]
}
`, name1, name2)
}
