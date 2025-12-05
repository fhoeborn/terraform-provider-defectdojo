package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGroupResource(t *testing.T) {
	name := fmt.Sprintf("dox-test-group-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGroupResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_group.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_group.test", "description", "test description"),
					resource.TestCheckResourceAttrSet("defectdojo_group.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccGroupResourceConfigUpdated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_group.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_group.test", "description", "updated description"),
				),
			},
		},
	})
}

func TestAccGroupResourceEmptyDescription(t *testing.T) {
	name := fmt.Sprintf("dox-test-group-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with empty description
			{
				Config: testAccGroupResourceConfigEmptyDescription(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_group.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_group.test", "description", ""),
					resource.TestCheckResourceAttrSet("defectdojo_group.test", "id"),
				),
			},
		},
	})
}

func testAccGroupResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_group" "test" {
  name = %[1]q
  description = "test description"
}
`, name)
}

func testAccGroupResourceConfigUpdated(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_group" "test" {
  name = %[1]q
  description = "updated description"
}
`, name)
}

func testAccGroupResourceConfigEmptyDescription(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_group" "test" {
  name = %[1]q
  description = ""
}
`, name)
}
