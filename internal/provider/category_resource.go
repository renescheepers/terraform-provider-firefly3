// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/renescheepers/terraform-provider-firefly3/internal/client"
)

// Interface guards
var _ resource.Resource = &CategoryResource{}
var _ resource.ResourceWithImportState = &CategoryResource{}

func NewCategoryResource() resource.Resource {
	return &CategoryResource{}
}

type CategoryResource struct {
	client *client.Client
}

type CategoryResourceModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Notes types.String `tfsdk:"notes"`
}

func (r *CategoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_category"
}

func (r *CategoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Firefly III category. Categories allow you to categorize transactions.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the category.",
			},
			"notes": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A description of the category.",
			},
		},
	}
}

func (r *CategoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CategoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CategoryResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	category, diags := r.modelToAPICategory(&data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdCategory, err := r.client.CreateCategory(ctx, category)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create category, got error: %s", err))
		return
	}

	r.apiCategoryToModel(createdCategory, &data)

	tflog.Trace(ctx, "created a category resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CategoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CategoryResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	category, err := r.client.GetCategory(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.Diagnostics.AddWarning("Category not found", fmt.Sprintf("Category %s not found", data.ID.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read category, got error: %s", err))
		return
	}

	r.apiCategoryToModel(category, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CategoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CategoryResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	category, diags := r.modelToAPICategory(&data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedCategory, err := r.client.UpdateCategory(ctx, data.ID.ValueString(), category)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update category, got error: %s, category %+v", err, category))
		return
	}

	r.apiCategoryToModel(updatedCategory, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CategoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CategoryResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCategory(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete category, got error: %s", err))
		return
	}
}

func (r *CategoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CategoryResource) modelToAPICategory(data *CategoryResourceModel) (*client.Category, diag.Diagnostics) {
	var diags diag.Diagnostics

	category := &client.Category{
		Name:  data.Name.ValueString(),
		Notes: data.Notes.ValueString(),
	}

	return category, diags
}

func (r *CategoryResource) apiCategoryToModel(category *client.Category, data *CategoryResourceModel) {
	data.ID = types.StringValue(category.ID)
	data.Name = types.StringValue(category.Name)
	data.Notes = types.StringValue(category.Notes)
}
