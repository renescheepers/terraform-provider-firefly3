// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/renescheepers/terraform-provider-firefly3/internal/client"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &Firefly3Provider{}
var _ provider.ProviderWithFunctions = &Firefly3Provider{}
var _ provider.ProviderWithEphemeralResources = &Firefly3Provider{}
var _ provider.ProviderWithActions = &Firefly3Provider{}

// Firefly3Provider defines the provider implementation.
type Firefly3Provider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Firefly3ProviderModel describes the provider data model.
type Firefly3ProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIKey   types.String `tfsdk:"api_key"`
}

func (p *Firefly3Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "firefly3"
	resp.Version = p.version
}

func (p *Firefly3Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Endpoint for the Firefly 3 API",
				Required:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key for the Firefly 3 API",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *Firefly3Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data Firefly3ProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Endpoint.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Endpoint",
			"The provider cannot create the Firefly III API client because the endpoint is not configured.",
		)
		return
	}

	if data.APIKey.IsNull() {
		resp.Diagnostics.AddError(
			"Missing API Key",
			"The provider cannot create the Firefly III API client because the API key is not configured.",
		)
		return
	}

	// Create the Firefly3 API client
	apiClient := client.NewClient(data.Endpoint.ValueString(), data.APIKey.ValueString())
	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *Firefly3Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRuleResource,
	}
}

func (p *Firefly3Provider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *Firefly3Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *Firefly3Provider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *Firefly3Provider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Firefly3Provider{
			version: version,
		}
	}
}
