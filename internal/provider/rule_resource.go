// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/renescheepers/terraform-provider-firefly3/internal/client"
)

// Interface guards
var _ resource.Resource = &RuleResource{}
var _ resource.ResourceWithImportState = &RuleResource{}

func NewRuleResource() resource.Resource {
	return &RuleResource{}
}

type RuleResource struct {
	client *client.Client
}

type RuleResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Title          types.String `tfsdk:"title"`
	Description    types.String `tfsdk:"description"`
	RuleGroupID    types.String `tfsdk:"rule_group_id"`
	Trigger        types.String `tfsdk:"trigger"`
	Active         types.Bool   `tfsdk:"active"`
	Strict         types.Bool   `tfsdk:"strict"`
	StopProcessing types.Bool   `tfsdk:"stop_processing"`
	Triggers       types.List   `tfsdk:"triggers"`
	Actions        types.List   `tfsdk:"actions"`
}

type RuleTriggerModel struct {
	Type           types.String `tfsdk:"type"`
	Value          types.String `tfsdk:"value"`
	Active         types.Bool   `tfsdk:"active"`
	Prohibited     types.Bool   `tfsdk:"prohibited"`
	StopProcessing types.Bool   `tfsdk:"stop_processing"`
}

type RuleActionModel struct {
	Type           types.String `tfsdk:"type"`
	Value          types.String `tfsdk:"value"`
	Active         types.Bool   `tfsdk:"active"`
	StopProcessing types.Bool   `tfsdk:"stop_processing"`
}

func (r *RuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rule"
}

func (r *RuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Firefly III rule. Rules allow you to automatically categorize, tag, or modify transactions based on triggers.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The title of the rule.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A description of what the rule does.",
			},
			"rule_group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the rule group under which the rule is stored.",
			},
			"trigger": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "When the rule should fire. Must be one of: `store-journal`, `update-journal`, or `manual-activation`.",
				Validators: []validator.String{
					stringvalidator.OneOf("store-journal"),
					// stringvalidator.OneOf("store-journal", "update-journal", "manual-activation"), only store-journal due to a bug in the API.
				},
			},
			"active": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether or not the rule is active. Defaults to `true`.",
			},
			"strict": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "If strict, ALL triggers must match for the rule to fire. Otherwise, just one is enough. Defaults to `true`.",
			},
			"stop_processing": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If true and the rule is triggered, other rules after this one in the group will be skipped. Defaults to `false`.",
			},
			"triggers": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "List of triggers that determine when the rule fires.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The type of trigger (e.g., `description_contains`, `amount_more`, `from_account_is`).",
						},
						"value": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The value to match against. Required for most trigger types, but some (like `has_any_tag`) only need 'true'.",
						},
						"active": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(true),
							MarkdownDescription: "Whether this trigger is active. Defaults to `true`.",
						},
						"prohibited": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
							MarkdownDescription: "If true, the trigger is negated (e.g., 'description is NOT'). Defaults to `false`.",
						},
						"stop_processing": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
							MarkdownDescription: "If true, other triggers will not be checked after this one fires. Defaults to `false`.",
						},
					},
				},
			},
			"actions": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "List of actions to perform when the rule fires.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The type of action (e.g., `set_category`, `add_tag`, `set_description`).",
						},
						"value": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The value for the action. Required for most action types.",
						},
						"active": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(true),
							MarkdownDescription: "Whether this action is active. Defaults to `true`.",
						},
						"stop_processing": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
							MarkdownDescription: "If true, other actions will not fire after this one. Defaults to `false`.",
						},
					},
				},
			},
		},
	}
}

func (r *RuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, diags := r.modelToAPIRule(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdRule, err := r.client.CreateRule(ctx, rule)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create rule, got error: %s", err))
		return
	}

	r.apiRuleToModel(createdRule, &data)

	tflog.Trace(ctx, "created a rule resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RuleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.GetRule(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rule, got error: %s", err))
		return
	}

	r.apiRuleToModel(rule, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, diags := r.modelToAPIRule(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedRule, err := r.client.UpdateRule(ctx, data.ID.ValueString(), rule)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update rule, got error: %s, rule %+v", err, rule))
		return
	}

	r.apiRuleToModel(updatedRule, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RuleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRule(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete rule, got error: %s", err))
		return
	}
}

func (r *RuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *RuleResource) modelToAPIRule(ctx context.Context, data *RuleResourceModel) (*client.Rule, diag.Diagnostics) {
	var diags diag.Diagnostics

	rule := &client.Rule{
		Title:          data.Title.ValueString(),
		Description:    data.Description.ValueString(),
		RuleGroupID:    data.RuleGroupID.ValueString(),
		Trigger:        data.Trigger.ValueString(),
		Active:         data.Active.ValueBool(),
		Strict:         data.Strict.ValueBool(),
		StopProcessing: data.StopProcessing.ValueBool(),
	}

	var triggerModels []RuleTriggerModel
	diags.Append(data.Triggers.ElementsAs(ctx, &triggerModels, false)...)
	if diags.HasError() {
		return nil, diags
	}

	rule.Triggers = make([]client.RuleTrigger, len(triggerModels))
	for i, t := range triggerModels {
		rule.Triggers[i] = client.RuleTrigger{
			Type:           t.Type.ValueString(),
			Value:          t.Value.ValueString(),
			Active:         t.Active.ValueBool(),
			Prohibited:     t.Prohibited.ValueBool(),
			StopProcessing: t.StopProcessing.ValueBool(),
		}
	}

	var actionModels []RuleActionModel
	diags.Append(data.Actions.ElementsAs(ctx, &actionModels, false)...)
	if diags.HasError() {
		return nil, diags
	}

	rule.Actions = make([]client.RuleAction, len(actionModels))
	for i, a := range actionModels {
		rule.Actions[i] = client.RuleAction{
			Type:           a.Type.ValueString(),
			Value:          a.Value.ValueString(),
			Active:         a.Active.ValueBool(),
			StopProcessing: a.StopProcessing.ValueBool(),
		}
	}

	return rule, diags
}

func (r *RuleResource) apiRuleToModel(rule *client.Rule, data *RuleResourceModel) {
	data.ID = types.StringValue(rule.ID)
	data.Title = types.StringValue(rule.Title)
	data.Description = types.StringValue(rule.Description)
	data.RuleGroupID = types.StringValue(rule.RuleGroupID)
	data.Trigger = types.StringValue(rule.Trigger)
	data.Active = types.BoolValue(rule.Active)
	data.Strict = types.BoolValue(rule.Strict)
	data.StopProcessing = types.BoolValue(rule.StopProcessing)

	triggerAttrTypes := map[string]attr.Type{
		"type":            types.StringType,
		"value":           types.StringType,
		"active":          types.BoolType,
		"prohibited":      types.BoolType,
		"stop_processing": types.BoolType,
	}

	triggerValues := make([]attr.Value, len(rule.Triggers))
	for i, t := range rule.Triggers {
		triggerValues[i], _ = types.ObjectValue(triggerAttrTypes, map[string]attr.Value{
			"type":            types.StringValue(t.Type),
			"value":           types.StringValue(t.Value),
			"active":          types.BoolValue(t.Active),
			"prohibited":      types.BoolValue(t.Prohibited),
			"stop_processing": types.BoolValue(t.StopProcessing),
		})
	}
	data.Triggers, _ = types.ListValue(types.ObjectType{AttrTypes: triggerAttrTypes}, triggerValues)

	actionAttrTypes := map[string]attr.Type{
		"type":            types.StringType,
		"value":           types.StringType,
		"active":          types.BoolType,
		"stop_processing": types.BoolType,
	}

	actionValues := make([]attr.Value, len(rule.Actions))
	for i, a := range rule.Actions {
		actionValues[i], _ = types.ObjectValue(actionAttrTypes, map[string]attr.Value{
			"type":            types.StringValue(a.Type),
			"value":           types.StringValue(a.Value),
			"active":          types.BoolValue(a.Active),
			"stop_processing": types.BoolValue(a.StopProcessing),
		})
	}
	data.Actions, _ = types.ListValue(types.ObjectType{AttrTypes: actionAttrTypes}, actionValues)
}
