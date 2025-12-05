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

type groupDataSource struct {
	terraformDatasource
}

func (t groupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DefectDojo Group",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the group",
				Computed:            true,
				Optional:            true,
			},
			"limit": schema.Int64Attribute{
				MarkdownDescription: "The amount of groups to fetch",
				Computed:            true,
				Optional:            true,
			},
			"groups": schema.ListNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the group",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the group",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							MarkdownDescription: "Identifier",
							Optional:            true,
						},
					},
				},	
			},
		},
	}
}

func (d groupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func NewGroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

func (r *groupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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


type groupDataSourceData struct {
	Name     types.String `tfsdk:"name"`
	Id       types.String `tfsdk:"id"`
	Limit    types.Int64  `tfsdk:"limit"`
	Groups 	 types.List   `tfsdk:"groups"`
}

type groupData struct {
	Name 		 types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Id     		 types.String `tfsdk:"id"`
}

func (d groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data groupDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	
	if resp.Diagnostics.HasError() {
		return
	}
	var (
		params dd.DojoGroupsListParams
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

	if !data.Name.IsNull() {
		params.Name = ref.Of(data.Name.ValueString())
	}

	if !data.Limit.IsNull() {
		params.Limit = ref.Of(int(data.Limit.ValueInt64()))
	}

	apiResp, err := d.client.DojoGroupsListWithResponse(ctx, &params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Resource",
			fmt.Sprintf("%s", err))
		return
	}

	if apiResp.StatusCode() == 200 {
		var group dd.DojoGroup
		groups := []groupData{}
		for _, element := range *apiResp.JSON200.Results {
			group = element

			var groupData groupData
			groupData.Id = types.StringValue(fmt.Sprintf("%d", group.Id))
			groupData.Name = types.StringValue(group.Name)
			groupData.Description = types.StringValue(*group.Description)
			
			groups = append(groups, groupData)
		}

		data.Groups, _ = types.ListValueFrom(ctx, data.Groups.ElementType(ctx), groups)
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
var _ datasource.DataSource = &groupDataSource{}
