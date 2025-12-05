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

func (t dojoGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DefectDojo Group",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Group",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the Group",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringDefault(""),
				},
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

type groupResourceData struct {
	Name        types.String `tfsdk:"name" ddField:"Name"`
	Description types.String `tfsdk:"description" ddField:"Description"`
	Id          types.String `tfsdk:"id" ddField:"Id"`
}

type groupDefectdojoResource struct {
	dd.DojoGroup
}

func (ddr *groupDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	tflog.Info(ctx, "createApiCall")

	reqBody := dd.DojoGroupsCreateJSONRequestBody(ddr.DojoGroup)
	apiResp, err := client.DojoGroupsCreateWithResponse(ctx, reqBody)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON201 != nil {
		ddr.DojoGroup = *apiResp.JSON201
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *groupDefectdojoResource) readApiCall(state *tfsdk.State, ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "readApiCall")
	apiResp, err := client.DojoGroupsRetrieveWithResponse(ctx, idNumber, &dd.DojoGroupsRetrieveParams{})
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.DojoGroup = *apiResp.JSON200
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *groupDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "updateApiCall")
	reqBody := dd.DojoGroupsUpdateJSONRequestBody(ddr.DojoGroup)
	apiResp, err := client.DojoGroupsUpdateWithResponse(ctx, idNumber, reqBody)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.DojoGroup = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *groupDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "deleteApiCall")
	apiResp, err := client.DojoGroupsDestroyWithResponse(ctx, idNumber)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, err
}

type dojoGroupResource struct {
	terraformResource
}

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &dojoGroupResource{}
var _ resource.ResourceWithImportState = &dojoGroupResource{}
var _ resource.ResourceWithConfigure = &dojoGroupResource{}

func NewGroupResource() resource.Resource {
	return &dojoGroupResource{
		terraformResource: terraformResource{
			dataProvider: groupDataProvider{},
		},
	}
}

func (r dojoGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

type groupDataProvider struct{}

func (r groupDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data groupResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *groupResourceData) id() types.String {
	return d.Id
}

func (d *groupResourceData) defectdojoResource() defectdojoResource {
	return &groupDefectdojoResource{
		DojoGroup: dd.DojoGroup{},
	}
}
