package provider

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"

	dd "github.com/doximity/defect-dojo-client-go"
	"github.com/doximity/terraform-provider-defectdojo/internal/ref"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type userDataSource struct {
	terraformDatasource
}

func (t userDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source for Defect Dojo Users",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The name of the user",
				Computed:            true,
				Optional:            true,
			},
			"limit": schema.Int64Attribute{
				MarkdownDescription: "The amount of users to fetch",
				Computed:            true,
				Optional:            true,
			},
			"users": schema.ListNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"username": schema.StringAttribute{
							MarkdownDescription: "The username of the User",
							Computed:            true,
						},
						"first_name": schema.StringAttribute{
							MarkdownDescription: "The first name of the User",
							Optional:            true,
							Computed:            true,
						},
						"last_name": schema.StringAttribute{
							MarkdownDescription: "The last name of the User",
							Optional:            true,
							Computed:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "The email of the User",
							Computed:            true,
						},
						"id": schema.StringAttribute{ // the id (for import purposes) MUST be a string
							Computed:            true,
							MarkdownDescription: "Identifier",
						},
					},
				},
			},
		},
	}
}

func (d userDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

func (r *userDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dd.ClientWithResponses)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected dd.ClientWithResponses, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

type userDataSourceData struct {
	Username types.String `tfsdk:"username"`
	Id       types.String `tfsdk:"id"`
	Limit    types.Int64  `tfsdk:"limit"`
	Users    types.List   `tfsdk:"users"`
}

type userData struct {
	Username  types.String `tfsdk:"username" ddField:"Username"`
	FirstName types.String `tfsdk:"first_name" ddField:"First_name"`
	LastName  types.String `tfsdk:"last_name" ddField:"Last_name"`
	Email     types.String `tfsdk:"email" ddField:"Email"`
	Id        types.String `tfsdk:"id" ddField:"Id"`
}

func (d userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data userDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	var (
		params dd.UsersListParams
	)
	if !data.Id.IsNull() {
		idNumber, err := strconv.Atoi(data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not Retrieve Resource",
				"The id field could not be parsed into an integer")
			return
		} else {
			params.Id = &idNumber
		}
	}

	if !data.Username.IsNull() {
		params.Username = ref.Of(data.Username.ValueString())
	}

	if !data.Limit.IsNull() {
		params.Limit = ref.Of(int(data.Limit.ValueInt64()))
	}

	apiResp, err := d.client.UsersListWithResponse(ctx, &params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Resource",
			fmt.Sprintf("%s", err))
		return
	}

	if apiResp.StatusCode() == 200 {
		var user dd.User
		users := []userData{}
		for _, element := range *apiResp.JSON200.Results {
			user = element

			var userData userData
			userData.Id = types.StringValue(fmt.Sprintf("%d", user.Id))
			userData.Username = types.StringValue(user.Username)
			if user.Email != nil {

			}
			userData.Email = types.StringValue(string(*user.Email))
			if user.FirstName != nil {
				userData.FirstName = types.StringValue(*user.FirstName)
			}
			if user.LastName != nil {
				userData.LastName = types.StringValue(*user.LastName)
			}

			users = append(users, userData)
		}

		data.Users, _ = types.ListValueFrom(ctx, data.Users.ElementType(ctx), users)
		diags = resp.State.Set(ctx, &data)
		resp.Diagnostics.Append(diags...)
	} else {
		body, _ := ioutil.ReadAll(apiResp.HTTPResponse.Body)

		resp.Diagnostics.AddError(
			"API Error Retrieving Data Source",
			fmt.Sprintf("Unexpected response code from API: %d", apiResp.StatusCode())+
				fmt.Sprintf("\n\nbody:\n\n%+v", body),
		)
		return
	}
}

// Ensure the implementation satisfies the desired interfaces.
var _ datasource.DataSource = &userDataSource{}
