// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/renescheepers/terraform-provider-firefly3/internal/client"
)

// Interface guards
var _ resource.Resource = &RuleGroupResource{}
var _ resource.ResourceWithImportState = &RuleGroupResource{}

func NewRuleGroupResource() resource.Resource {
	return &RuleGroupResource{}
}

type RuleGroupResource struct {
	client *client.Client
}

type RuleGroupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	Order       types.Int32  `tfsdk:"order"`
	Active      types.Bool   `tfsdk:"active"`
}

func (r *RuleGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rule_group"
}

func (r *RuleGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Firefly III rule group. Rule groups are containers for rules and determine the order in which rules are executed.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the rule group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The title of the rule group.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A description of what the rule group is for.",
			},
			"order": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The order of the rule group. Rule groups with a lower order are executed first.",
			},
			"active": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether or not the rule group is active. Defaults to `true`.",
			},
		},
	}
}

func (r *RuleGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *RuleGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RuleGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleGroup := r.modelToAPIRuleGroup(&data)

	createdRuleGroup, err := r.client.CreateRuleGroup(ctx, ruleGroup)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create rule group, got error: %s", err))
		return
	}

	r.apiRuleGroupToModel(createdRuleGroup, &data)

	tflog.Trace(ctx, "created a rule group resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuleGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RuleGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleGroup, err := r.client.GetRuleGroup(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rule group, got error: %s", err))
		return
	}

	r.apiRuleGroupToModel(ruleGroup, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuleGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RuleGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleGroup := r.modelToAPIRuleGroup(&data)

	updatedRuleGroup, err := r.client.UpdateRuleGroup(ctx, data.ID.ValueString(), ruleGroup)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update rule group, got error: %s", err))
		return
	}

	r.apiRuleGroupToModel(updatedRuleGroup, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuleGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RuleGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRuleGroup(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete rule group, got error: %s", err))
		return
	}
}

func (r *RuleGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *RuleGroupResource) modelToAPIRuleGroup(data *RuleGroupResourceModel) *client.RuleGroup {
	ruleGroup := &client.RuleGroup{
		Title:       data.Title.ValueString(),
		Description: data.Description.ValueString(),
		Active:      data.Active.ValueBool(),
	}

	if !data.Order.IsNull() && !data.Order.IsUnknown() {
		ruleGroup.Order = data.Order.ValueInt32()
	}

	return ruleGroup
}

func (r *RuleGroupResource) apiRuleGroupToModel(ruleGroup *client.RuleGroup, data *RuleGroupResourceModel) {
	data.ID = types.StringValue(ruleGroup.ID)
	data.Title = types.StringValue(ruleGroup.Title)
	data.Description = types.StringValue(ruleGroup.Description)
	data.Order = types.Int32Value(ruleGroup.Order)
	data.Active = types.BoolValue(ruleGroup.Active)
}
