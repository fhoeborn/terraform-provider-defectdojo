package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProductsDataSource(t *testing.T) {
	name := fmt.Sprintf("dox-test-repo-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing with name filter
			{
				Config: testAccProductsDataSourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_products.test", "name", name),
					resource.TestCheckResourceAttrSet("data.defectdojo_products.test", "products.#"),
					resource.TestCheckResourceAttr("data.defectdojo_products.test", "products.0.name", name),
					resource.TestCheckResourceAttr("data.defectdojo_products.test", "products.0.description", "test"),
					resource.TestCheckResourceAttrSet("data.defectdojo_products.test", "products.0.id"),
					resource.TestCheckResourceAttrSet("data.defectdojo_products.test", "products.0.product_type_id"),
				),
			},
		},
	})
}

func TestAccProductsDataSourceWithLimit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing with limit
			{
				Config: testAccProductsDataSourceConfigWithLimit(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_products.test", "limit", "5"),
					resource.TestCheckResourceAttrSet("data.defectdojo_products.test", "products.#"),
				),
			},
		},
	})
}

func TestAccProductsDataSourceWithId(t *testing.T) {
	name := fmt.Sprintf("dox-test-repo-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing with specific ID
			{
				Config: testAccProductsDataSourceConfigWithId(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.defectdojo_products.test", "id"),
					resource.TestCheckResourceAttrSet("data.defectdojo_products.test", "products.#"),
					resource.TestCheckResourceAttr("data.defectdojo_products.test", "products.0.name", name),
					resource.TestCheckResourceAttr("data.defectdojo_products.test", "products.0.description", "test"),
				),
			},
		},
	})
}

func testAccProductsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
  description = "test"
  product_type_id = 1
}
data "defectdojo_products" "test" {
  name = defectdojo_product.test.name
  depends_on = [defectdojo_product.test]
}
`, name)
}

func testAccProductsDataSourceConfigWithLimit() string {
	return `
provider "defectdojo" {}
data "defectdojo_products" "test" {
  limit = 5
}
`
}

func testAccProductsDataSourceConfigWithId(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
  description = "test"
  product_type_id = 1
}
data "defectdojo_products" "test" {
  id = defectdojo_product.test.id
  depends_on = [defectdojo_product.test]
}
`, name)
}