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

type productsDataSource struct {
	terraformDatasource
}

func (t productsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source for Defect Dojo Products",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Product",
				Computed:            true,
				Optional:            true,
			},
			"limit": schema.Int64Attribute{
				MarkdownDescription: "The amount of Products to fetch",
				Computed:            true,
				Optional:            true,
			},
			"products": schema.ListNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the Product",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the Product",
							Computed:            true,
						},
						"product_type_id": schema.Int64Attribute{
							MarkdownDescription: "The ID of the Product Type",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							MarkdownDescription: "Identifier",
							Optional:            true,
						},
						"prod_numeric_grade": schema.Int64Attribute{
							MarkdownDescription: "The Numeric Grade of the Product",
							Optional:            true,
						},
						"business_criticality": schema.StringAttribute{
							MarkdownDescription: "The Business Criticality of the Product. Valid values are: 'very high', 'high', 'medium', 'low', 'very low', 'none'",
							Optional:            true,
						},
						"platform": schema.StringAttribute{
							MarkdownDescription: "The Platform of the Product. Valid values are: 'web service', 'desktop', 'iot', 'mobile', 'web'",
							Computed:            true,
							Optional:            true,
						},
						"life_cycle": schema.StringAttribute{
							MarkdownDescription: "The Lifecycle state of the Product. Valid values are: 'construction', 'production', 'retirement'",
							Computed:            true,
							Optional:            true,
						},
						"origin": schema.StringAttribute{
							MarkdownDescription: "The Origin of the Product. Valid values are: 'third party library', 'purchased', 'contractor', 'internal', 'open source', 'outsourced'",
							Computed:            true,
							Optional:            true,
						},
						"user_records": schema.Int64Attribute{
							MarkdownDescription: "Estimate the number of user records within the application.",
							Computed:            true,
							Optional:            true,
						},
						"revenue": schema.StringAttribute{
							MarkdownDescription: "Estimate the application's revenue.",
							Computed:            true,
							Optional:            true,
						},
						"external_audience": schema.BoolAttribute{
							MarkdownDescription: "Specify if the application is used by people outside the organization.",
							Computed:            true,
							Optional:            true,
						},
						"internet_accessible": schema.BoolAttribute{
							MarkdownDescription: "Specify if the application is accessible from the public internet.",
							Computed:            true,
							Optional:            true,
						},
						"enable_skip_risk_acceptance": schema.BoolAttribute{
							MarkdownDescription: "Allows simple risk acceptance by checking/unchecking a checkbox.",
							Computed:            true,
							Optional:            true,
						},
						"enable_full_risk_acceptance": schema.BoolAttribute{
							MarkdownDescription: "Allows full risk acceptance using a risk acceptance form, expiration date, uploaded proof, etc.",
							Computed:            true,
							Optional:            true,
						},
						"product_manager_id": schema.Int64Attribute{
							MarkdownDescription: "The ID of the user who is the PM for this product.",
							Computed:            true,
							Optional:            true,
						},
						"technical_contact_id": schema.Int64Attribute{
							MarkdownDescription: "The ID of the user who is the technical contact for this product.",
							Computed:            true,
							Optional:            true,
						},
						"team_manager_id": schema.Int64Attribute{
							MarkdownDescription: "The ID of the user who is the manager for this product.",
							Computed:            true,
							Optional:            true,
						},
						"regulation_ids": schema.SetAttribute{
							MarkdownDescription: "The IDs of the Regulations which apply to this product.",
							Computed:            true,
							ElementType:         types.Int64Type,
							Optional:            true,
						},
						"tags": schema.SetAttribute{
							MarkdownDescription: "Tags to apply to the product",
							Computed:            true,
							ElementType:         types.StringType,
							Optional:            true,
						},
					},
				},
			},
		},
	}
}

func (d productsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_products"
}

func NewProductsDataSource() datasource.DataSource {
	return &productsDataSource{}
}

func (r *productsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type productDataSourceData struct {
	Name     types.String `tfsdk:"name"`
	Id       types.String `tfsdk:"id"`
	Limit    types.Int64  `tfsdk:"limit"`
	Products types.List   `tfsdk:"products"`
}

type productData struct {
	Name                       types.String `tfsdk:"name"`
	Description                types.String `tfsdk:"description"`
	ProductTypeId              types.Int64  `tfsdk:"product_type_id"`
	Id                         types.String `tfsdk:"id"`
	BusinessCriticality        types.String `tfsdk:"business_criticality"`
	EnableFullRiskAcceptance   types.Bool   `tfsdk:"enable_full_risk_acceptance"`
	EnableSimpleRiskAcceptance types.Bool   `tfsdk:"enable_skip_risk_acceptance"`
	ExternalAudience           types.Bool   `tfsdk:"external_audience"`
	InternetAccessible         types.Bool   `tfsdk:"internet_accessible"`
	Lifecycle                  types.String `tfsdk:"life_cycle"`
	Origin                     types.String `tfsdk:"origin"`
	Platform                   types.String `tfsdk:"platform"`
	ProdNumericGrade           types.Int64  `tfsdk:"prod_numeric_grade"`
	ProductManagerId           types.Int64  `tfsdk:"product_manager_id"`
	RegulationIds              types.Set    `tfsdk:"regulation_ids"`
	Revenue                    types.String `tfsdk:"revenue"`
	Tags                       types.Set    `tfsdk:"tags"`
	TeamManagerId              types.Int64  `tfsdk:"team_manager_id"`
	TechnicalContactId         types.Int64  `tfsdk:"technical_contact_id"`
	UserRecords                types.Int64  `tfsdk:"user_records"`
}

func (d productsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data productDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	var (
		params dd.ProductsListParams
	)
	if !data.Id.IsNull() {
		idNumber, err := strconv.Atoi(data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not Retrieve Resource",
				"The id field could not be parsed into an integer")
			return
		} else {
			params.Id = &[]int{idNumber}
		}
	}

	if !data.Name.IsNull() {
		params.Name = ref.Of(data.Name.ValueString())
	}

	if !data.Limit.IsNull() {
		params.Limit = ref.Of(int(data.Limit.ValueInt64()))
	}

	apiResp, err := d.client.ProductsListWithResponse(ctx, &params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Resource",
			fmt.Sprintf("%s", err))
		return
	}

	if apiResp.StatusCode() == 200 {
		var product dd.Product
		products := []productData{}
		for _, element := range *apiResp.JSON200.Results {
			product = element

			var productData productData
			productData.Id = types.StringValue(fmt.Sprintf("%d", product.Id))
			productData.Name = types.StringValue(product.Name)
			productData.Description = types.StringValue(product.Description)
			if product.BusinessCriticality != nil {
				productData.BusinessCriticality = types.StringValue(string(*product.BusinessCriticality))
			}
			if product.EnableFullRiskAcceptance != nil {
				productData.EnableFullRiskAcceptance = types.BoolValue(*product.EnableFullRiskAcceptance)
			}
			if product.EnableSimpleRiskAcceptance != nil {
				productData.EnableSimpleRiskAcceptance = types.BoolValue(*product.EnableSimpleRiskAcceptance)
			}
			if product.ExternalAudience != nil {
				productData.ExternalAudience = types.BoolValue(*product.ExternalAudience)
			}
			productData.InternetAccessible = types.BoolValue(*product.InternetAccessible)
			productData.ProductTypeId = types.Int64Value(int64(product.ProdType))
			if product.Lifecycle != nil {
				productData.Lifecycle = types.StringValue(string(*product.Lifecycle))
			}
			if product.Origin != nil {
				productData.Origin = types.StringValue(string(*product.Origin))
			}
			if product.Platform != nil {
				productData.Platform = types.StringValue(string(*product.Platform))
			}
			if product.ProdNumericGrade != nil {
				productData.ProdNumericGrade = types.Int64Value(int64(*product.ProdNumericGrade))
			}
			if product.ProductManager != nil {
				productData.ProductManagerId = types.Int64Value(int64(*product.ProductManager))
			}
			if product.Regulations != nil {
				productData.RegulationIds, _ = types.SetValueFrom(ctx, types.Int64Type, *product.Regulations)
			}
			if product.Revenue != nil {
				productData.Revenue = types.StringValue(*product.Revenue)
			}
			if product.Tags != nil {
				productData.Tags, _ = types.SetValueFrom(ctx, types.StringType, *product.Tags)
			}
			if product.TeamManager != nil {
				productData.TeamManagerId = types.Int64Value(int64(*product.TeamManager))
			}
			if product.TechnicalContact != nil {
				productData.TechnicalContactId = types.Int64Value(int64(*product.TechnicalContact))
			}
			if product.UserRecords != nil {
				productData.UserRecords = types.Int64Value(int64(*product.UserRecords))
			}

			products = append(products, productData)
		}

		data.Products, _ = types.ListValueFrom(ctx, data.Products.ElementType(ctx), products)
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
var _ datasource.DataSource = &productsDataSource{}
