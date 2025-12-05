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

func (t dojoGroupMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DefectDojo Group Membership",

		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				MarkdownDescription: "The id of the Group",
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

type groupMemberResourceData struct {
	Group types.String `tfsdk:"group_id" ddField:"Group"`
	User  types.String `tfsdk:"user_id" ddField:"User"`
	Role  types.String `tfsdk:"role_id" ddField:"Role"`
	Id    types.String `tfsdk:"id" ddField:"Id"`
}

type groupMemberDefectdojoResource struct {
	dd.DojoGroupMember
}

func (ddr *groupMemberDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	tflog.Info(ctx, "createApiCall")

	reqBody := dd.DojoGroupMembersCreateJSONRequestBody(ddr.DojoGroupMember)
	apiResp, err := client.DojoGroupMembersCreateWithResponse(ctx, reqBody)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON201 != nil {
		ddr.DojoGroupMember = *apiResp.JSON201
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *groupMemberDefectdojoResource) readApiCall(state *tfsdk.State, ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "readApiCall")
	apiResp, err := client.DojoGroupMembersRetrieveWithResponse(ctx, idNumber, &dd.DojoGroupMembersRetrieveParams{})
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.DojoGroupMember = *apiResp.JSON200
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *groupMemberDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "updateApiCall")
	reqBody := dd.DojoGroupMembersUpdateJSONRequestBody(ddr.DojoGroupMember)
	apiResp, err := client.DojoGroupMembersUpdateWithResponse(ctx, idNumber, reqBody)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.DojoGroupMember = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *groupMemberDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "deleteApiCall")
	apiResp, err := client.DojoGroupMembersDestroyWithResponse(ctx, idNumber)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, err
}

type dojoGroupMemberResource struct {
	terraformResource
}

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &dojoGroupResource{}
var _ resource.ResourceWithImportState = &dojoGroupResource{}
var _ resource.ResourceWithConfigure = &dojoGroupResource{}

func NewGroupMemberResource() resource.Resource {
	return &dojoGroupMemberResource{
		terraformResource: terraformResource{
			dataProvider: groupMemberDataProvider{},
		},
	}
}

func (r dojoGroupMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_member"
}

type groupMemberDataProvider struct{}

func (r groupMemberDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data groupMemberResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *groupMemberResourceData) id() types.String {
	return d.Id
}

func (d *groupMemberResourceData) defectdojoResource() defectdojoResource {
	return &groupMemberDefectdojoResource{
		DojoGroupMember: dd.DojoGroupMember{},
	}
}
