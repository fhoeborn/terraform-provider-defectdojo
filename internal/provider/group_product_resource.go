package provider

import (
	"context"
	"fmt"

	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (t dojoGroupProductResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DefectDojo Group Access to Product",

		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				MarkdownDescription: "The id of the Group",
				Required:            true,
			},
			"product_id": schema.StringAttribute{
				MarkdownDescription: "The id of the product",
				Required:            true,
			},
			"role_id": schema.StringAttribute{
				MarkdownDescription: "The id of the role of the Group in the product",
				Required:            true,
			},
			"id": schema.StringAttribute{ // the id (for import purposes) MUST be a string
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

type groupProductResourceData struct {
	Group   types.String `tfsdk:"group_id" ddField:"Group"`
	Product types.String `tfsdk:"product_id" ddField:"Product"`
	Role    types.String `tfsdk:"role_id" ddField:"Role"`
	Id      types.String `tfsdk:"id" ddField:"Id"`
}

type groupProductDefectdojoResource struct {
	dd.ProductGroup
}

func (ddr *groupProductDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	tflog.Info(ctx, "createApiCall")

	reqBody := dd.ProductGroupsCreateJSONRequestBody(ddr.ProductGroup)
	apiResp, err := client.ProductGroupsCreateWithResponse(ctx, reqBody)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON201 != nil {
		ddr.ProductGroup = *apiResp.JSON201
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *groupProductDefectdojoResource) readApiCall(state *tfsdk.State, ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "readApiCall")
	apiResp, err := client.ProductGroupsRetrieveWithResponse(ctx, idNumber, &dd.ProductGroupsRetrieveParams{})
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.ProductGroup = *apiResp.JSON200
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *groupProductDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "updateApiCall")
	reqBody := dd.ProductGroupsUpdateJSONRequestBody(ddr.ProductGroup)
	apiResp, err := client.ProductGroupsUpdateWithResponse(ctx, idNumber, reqBody)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.ProductGroup = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *groupProductDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "deleteApiCall")
	apiResp, err := client.ProductGroupsDestroyWithResponse(ctx, idNumber)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, err
}

type dojoGroupProductResource struct {
	terraformResource
}

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &dojoGroupResource{}
var _ resource.ResourceWithImportState = &dojoGroupResource{}
var _ resource.ResourceWithConfigure = &dojoGroupResource{}

func NewGroupProductResource() resource.Resource {
	return &dojoGroupProductResource{
		terraformResource: terraformResource{
			dataProvider: groupProductDataProvider{},
		},
	}
}

func (r dojoGroupProductResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_product"
}

type groupProductDataProvider struct{}

func (r groupProductDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data groupProductResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *groupProductResourceData) id() types.String {
	return d.Id
}

func (d *groupProductResourceData) defectdojoResource() defectdojoResource {
	return &groupProductDefectdojoResource{
		ProductGroup: dd.ProductGroup{},
	}
}
