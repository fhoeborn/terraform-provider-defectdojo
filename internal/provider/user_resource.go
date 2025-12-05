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

func (t dojoUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DefectDojo User",

		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				MarkdownDescription: "The username of the User",
				Required:            true,
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "The first name of the User",
				Optional:            true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "The last name of the User",
				Optional:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email of the User",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password for the user",
				Required:            true,
				Sensitive:           true,
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

type userResourceData struct {
	Username  types.String `tfsdk:"username" ddField:"Username"`
	FirstName types.String `tfsdk:"first_name" ddField:"First_name"`
	LastName  types.String `tfsdk:"last_name" ddField:"Last_name"`
	Email     types.String `tfsdk:"email" ddField:"Email"`
	Password  types.String `tfsdk:"password" ddField:"Password"`
	Id        types.String `tfsdk:"id" ddField:"Id"`
}

type userDefectdojoResource struct {
	dd.User
}

func (ddr *userDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	tflog.Info(ctx, "createApiCall")

	password := ddr.User.Password
	reqBody := dd.UsersCreateJSONRequestBody(ddr.User)
	apiResp, err := client.UsersCreateWithResponse(ctx, reqBody)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON201 != nil {
		ddr.User = *apiResp.JSON201
		ddr.User.Password = password
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *userDefectdojoResource) readApiCall(state *tfsdk.State, ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "readApiCall")

	apiResp, err := client.UsersRetrieveWithResponse(ctx, idNumber)
	var userDataFromState userResourceData

	state.Get(ctx, &userDataFromState)
	if apiResp.JSON200 != nil {
		ddr.User = *apiResp.JSON200
		password := userDataFromState.Password.ValueString()
		ddr.User.Password = &password
	}

	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *userDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "updateApiCall")
	reqBody := dd.UsersUpdateJSONRequestBody(ddr.User)
	apiResp, err := client.UsersUpdateWithResponse(ctx, idNumber, reqBody)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.User = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *userDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "deleteApiCall")
	apiResp, err := client.UsersDestroyWithResponse(ctx, idNumber)
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, err
}

type dojoUserResource struct {
	terraformResource
}

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &dojoUserResource{}
var _ resource.ResourceWithImportState = &dojoUserResource{}
var _ resource.ResourceWithConfigure = &dojoUserResource{}

func NewUserResource() resource.Resource {
	return &dojoUserResource{
		terraformResource: terraformResource{
			dataProvider: userDataProvider{},
		},
	}
}

func (r dojoUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

type userDataProvider struct{}

func (r userDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data userResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *userResourceData) id() types.String {
	return d.Id
}

func (d *userResourceData) defectdojoResource() defectdojoResource {
	return &userDefectdojoResource{
		User: dd.User{},
	}
}
