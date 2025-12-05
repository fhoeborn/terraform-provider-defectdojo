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

func (t productTypeMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DefectDojo Product Type Membership",

		Attributes: map[string]schema.Attribute{
			"product_type_id": schema.StringAttribute{
				MarkdownDescription: "The id of the Product Type",
				Required:            true,
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "The id of the member",
				Required:            true,
			},
			"role_id": schema.StringAttribute{
				MarkdownDescription: "The id of the role of the member of the Group",
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

type dojoProductTypeMemberResourceData struct {
	ProductType types.String `tfsdk:"product_type_id" ddField:"ProductType"`
	User        types.String `tfsdk:"user_id" ddField:"User"`
	Role        types.String `tfsdk:"role_id" ddField:"Role"`
	Id          types.String `tfsdk:"id" ddField:"Id"`
}

type dojoProductTypeMemberResource struct {
	dd.ProductTypeMember
}

func (ddr *dojoProductTypeMemberResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	tflog.Info(ctx, "createApiCall")

	reqBody := dd.ProductTypeMembersCreateJSONRequestBody(ddr.ProductTypeMember)
	apiResp, err := client.ProductTypeMembersCreateWithResponse(ctx, reqBody)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON201 != nil {
		ddr.ProductTypeMember = *apiResp.JSON201
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *dojoProductTypeMemberResource) readApiCall(state *tfsdk.State, ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "readApiCall")
	apiResp, err := client.ProductTypeMembersRetrieveWithResponse(ctx, idNumber, &dd.ProductTypeMembersRetrieveParams{})
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.ProductTypeMember = *apiResp.JSON200
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *dojoProductTypeMemberResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "updateApiCall")
	reqBody := dd.ProductTypeMembersUpdateJSONRequestBody(ddr.ProductTypeMember)
	apiResp, err := client.ProductTypeMembersUpdateWithResponse(ctx, idNumber, reqBody)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.ProductTypeMember = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *dojoProductTypeMemberResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "deleteApiCall")
	apiResp, err := client.ProductTypeMembersDestroyWithResponse(ctx, idNumber)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, err
}

type productTypeMemberResource struct {
	terraformResource
}

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &productTypeMemberResource{}
var _ resource.ResourceWithImportState = &productTypeMemberResource{}
var _ resource.ResourceWithConfigure = &productTypeMemberResource{}

func NewProductTypeMemberResource() resource.Resource {
	return &productTypeMemberResource{
		terraformResource: terraformResource{
			dataProvider: productTypeMemberDataProvider{},
		},
	}
}

func (r productTypeMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_type_member"
}

type productTypeMemberDataProvider struct{}

func (r productTypeMemberDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data dojoProductTypeMemberResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *dojoProductTypeMemberResourceData) id() types.String {
	return d.Id
}

func (d *dojoProductTypeMemberResourceData) defectdojoResource() defectdojoResource {
	return &dojoProductTypeMemberResource{
		ProductTypeMember: dd.ProductTypeMember{},
	}
}
