// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package amp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/amp"
	awstypes "github.com/aws/aws-sdk-go-v2/service/amp/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/enum"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/fwdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	fwflex "github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	fwtypes "github.com/hashicorp/terraform-provider-aws/internal/framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkResource("aws_prometheus_scraper_logging_configuration", name="ScraperLoggingConfiguration")
// @Testing(existsType="github.com/aws/aws-sdk-go-v2/service/amp;amp.DescribeScraperLoggingConfigurationOutput")
func newScraperLoggingConfigurationResource(_ context.Context) (resource.ResourceWithConfigure, error) {
	r := &scraperLoggingConfigurationResource{}

	r.SetDefaultCreateTimeout(5 * time.Minute)
	r.SetDefaultUpdateTimeout(5 * time.Minute)
	r.SetDefaultDeleteTimeout(5 * time.Minute)

	return r, nil
}

type scraperLoggingConfigurationResource struct {
	framework.ResourceWithModel[scraperLoggingConfigurationResourceModel]
	framework.WithTimeouts
}

func (r *scraperLoggingConfigurationResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"scraper_components": schema.SetAttribute{
				CustomType:  fwtypes.SetOfStringEnumType[awstypes.ScraperComponentType](),
				ElementType: fwtypes.StringEnumType[awstypes.ScraperComponentType](),
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"scraper_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"logging_destination": schema.ListNestedBlock{
				CustomType: fwtypes.NewListNestedObjectTypeOf[scraperLoggingDestinationModel](ctx),
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.SizeAtMost(1),
					listvalidator.IsRequired(),
				},
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						names.AttrCloudWatchLogs: schema.ListNestedBlock{
							CustomType: fwtypes.NewListNestedObjectTypeOf[cloudWatchLogDestinationModel](ctx),
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
								listvalidator.SizeAtMost(1),
								listvalidator.IsRequired(),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"log_group_arn": schema.StringAttribute{
										CustomType: fwtypes.ARNType,
										Required:   true,
										Validators: []validator.String{
											stringvalidator.RegexMatches(regexache.MustCompile(`:\*$`), "ARN must end with `:*`"),
										},
									},
								},
							},
						},
					},
				},
			},
			names.AttrTimeouts: timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

func (r *scraperLoggingConfigurationResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data scraperLoggingConfigurationResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	conn := r.Meta().AMPClient(ctx)

	scraperID := fwflex.StringValueFromFramework(ctx, data.ScraperID)
	var input amp.UpdateScraperLoggingConfigurationInput
	response.Diagnostics.Append(fwflex.Expand(ctx, data, &input)...)
	if response.Diagnostics.HasError() {
		return
	}
	input.ScraperComponents = expandScraperComponents(ctx, data.ScraperComponents)

	_, err := conn.UpdateScraperLoggingConfiguration(ctx, &input)

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("creating Prometheus Scraper Logging Configuration (%s)", scraperID), err.Error())

		return
	}

	if _, err := waitScraperLoggingConfigurationCreated(ctx, conn, scraperID, r.CreateTimeout(ctx, data.Timeouts)); err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("waiting for Prometheus Scraper Logging Configuration (%s) create", scraperID), err.Error())

		return
	}

	// Read back to get computed fields (e.g. scraper_components defaulted by API).
	output, err := findScraperLoggingConfigurationByID(ctx, conn, scraperID)
	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("reading Prometheus Scraper Logging Configuration (%s) after create", scraperID), err.Error())

		return
	}
	response.Diagnostics.Append(r.flattenIntoModel(ctx, output, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, data)...)
}

func (r *scraperLoggingConfigurationResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data scraperLoggingConfigurationResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	conn := r.Meta().AMPClient(ctx)

	scraperID := fwflex.StringValueFromFramework(ctx, data.ScraperID)
	output, err := findScraperLoggingConfigurationByID(ctx, conn, scraperID)

	if retry.NotFound(err) {
		response.Diagnostics.Append(fwdiag.NewResourceNotFoundWarningDiagnostic(err))
		response.State.RemoveResource(ctx)

		return
	}

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("reading Prometheus Scraper Logging Configuration (%s)", scraperID), err.Error())

		return
	}

	response.Diagnostics.Append(r.flattenIntoModel(ctx, output, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *scraperLoggingConfigurationResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var new scraperLoggingConfigurationResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &new)...)
	if response.Diagnostics.HasError() {
		return
	}

	conn := r.Meta().AMPClient(ctx)

	scraperID := fwflex.StringValueFromFramework(ctx, new.ScraperID)
	var input amp.UpdateScraperLoggingConfigurationInput
	response.Diagnostics.Append(fwflex.Expand(ctx, new, &input)...)
	if response.Diagnostics.HasError() {
		return
	}
	input.ScraperComponents = expandScraperComponents(ctx, new.ScraperComponents)

	_, err := conn.UpdateScraperLoggingConfiguration(ctx, &input)

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("updating Prometheus Scraper Logging Configuration (%s)", scraperID), err.Error())

		return
	}

	if _, err := waitScraperLoggingConfigurationUpdated(ctx, conn, scraperID, r.UpdateTimeout(ctx, new.Timeouts)); err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("waiting for Prometheus Scraper Logging Configuration (%s) update", scraperID), err.Error())

		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, new)...)
}

func (r *scraperLoggingConfigurationResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data scraperLoggingConfigurationResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	conn := r.Meta().AMPClient(ctx)

	scraperID := fwflex.StringValueFromFramework(ctx, data.ScraperID)
	input := amp.DeleteScraperLoggingConfigurationInput{
		ScraperId:   aws.String(scraperID),
		ClientToken: aws.String(create.UniqueId(ctx)),
	}
	_, err := conn.DeleteScraperLoggingConfiguration(ctx, &input)

	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return
	}

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("deleting Prometheus Scraper Logging Configuration (%s)", scraperID), err.Error())

		return
	}

	if _, err := waitScraperLoggingConfigurationDeleted(ctx, conn, scraperID, r.DeleteTimeout(ctx, data.Timeouts)); err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("waiting for Prometheus Scraper Logging Configuration (%s) delete", scraperID), err.Error())

		return
	}
}

func (r *scraperLoggingConfigurationResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("scraper_id"), request, response)
}

// flattenIntoModel fills data from the API output, handling the autoflex-incompatible scraper_components field.
func (r *scraperLoggingConfigurationResource) flattenIntoModel(ctx context.Context, output *amp.DescribeScraperLoggingConfigurationOutput, data *scraperLoggingConfigurationResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	diags.Append(fwflex.Flatten(ctx, output, data)...)
	if diags.HasError() {
		return diags
	}

	// Flatten scraper_components from []ScraperComponent → set of type strings.
	componentTypes := make([]awstypes.ScraperComponentType, len(output.ScraperComponents))
	for i, c := range output.ScraperComponents {
		componentTypes[i] = c.Type
	}
	data.ScraperComponents = fwflex.FlattenFrameworkStringyValueSetOfStringEnum(ctx, componentTypes)

	return diags
}

func findScraperLoggingConfigurationByID(ctx context.Context, conn *amp.Client, id string) (*amp.DescribeScraperLoggingConfigurationOutput, error) {
	input := amp.DescribeScraperLoggingConfigurationInput{
		ScraperId: aws.String(id),
	}
	output, err := conn.DescribeScraperLoggingConfiguration(ctx, &input)

	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError: err,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.Status == nil {
		return nil, tfresource.NewEmptyResultError()
	}

	return output, nil
}

func statusScraperLoggingConfiguration(conn *amp.Client, id string) retry.StateRefreshFunc {
	return func(ctx context.Context) (any, string, error) {
		output, err := findScraperLoggingConfigurationByID(ctx, conn, id)

		if retry.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, string(output.Status.StatusCode), nil
	}
}

func waitScraperLoggingConfigurationCreated(ctx context.Context, conn *amp.Client, id string, timeout time.Duration) (*amp.DescribeScraperLoggingConfigurationOutput, error) {
	stateConf := &retry.StateChangeConf{
		Pending: enum.Slice(awstypes.ScraperLoggingConfigurationStatusCodeCreating, awstypes.ScraperLoggingConfigurationStatusCodeUpdating),
		Target:  enum.Slice(awstypes.ScraperLoggingConfigurationStatusCodeActive),
		Refresh: statusScraperLoggingConfiguration(conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*amp.DescribeScraperLoggingConfigurationOutput); ok {
		retry.SetLastError(err, errors.New(aws.ToString(output.Status.StatusReason)))

		return output, err
	}

	return nil, err
}

func waitScraperLoggingConfigurationUpdated(ctx context.Context, conn *amp.Client, id string, timeout time.Duration) (*amp.DescribeScraperLoggingConfigurationOutput, error) {
	stateConf := &retry.StateChangeConf{
		Pending: enum.Slice(awstypes.ScraperLoggingConfigurationStatusCodeUpdating),
		Target:  enum.Slice(awstypes.ScraperLoggingConfigurationStatusCodeActive),
		Refresh: statusScraperLoggingConfiguration(conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*amp.DescribeScraperLoggingConfigurationOutput); ok {
		retry.SetLastError(err, errors.New(aws.ToString(output.Status.StatusReason)))

		return output, err
	}

	return nil, err
}

func waitScraperLoggingConfigurationDeleted(ctx context.Context, conn *amp.Client, id string, timeout time.Duration) (*amp.DescribeScraperLoggingConfigurationOutput, error) {
	stateConf := &retry.StateChangeConf{
		Pending: enum.Slice(awstypes.ScraperLoggingConfigurationStatusCodeDeleting, awstypes.ScraperLoggingConfigurationStatusCodeActive),
		Target:  []string{},
		Refresh: statusScraperLoggingConfiguration(conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*amp.DescribeScraperLoggingConfigurationOutput); ok {
		retry.SetLastError(err, errors.New(aws.ToString(output.Status.StatusReason)))

		return output, err
	}

	return nil, err
}

type scraperLoggingConfigurationResourceModel struct {
	framework.WithRegionModel
	LoggingDestination fwtypes.ListNestedObjectValueOf[scraperLoggingDestinationModel] `tfsdk:"logging_destination"`
	ScraperComponents  fwtypes.SetOfStringEnum[awstypes.ScraperComponentType]          `tfsdk:"scraper_components" autoflex:"-"`
	ScraperID          types.String                                                    `tfsdk:"scraper_id"`
	Timeouts           timeouts.Value                                                  `tfsdk:"timeouts"`
}

type scraperLoggingDestinationModel struct {
	CloudwatchLogs fwtypes.ListNestedObjectValueOf[cloudWatchLogDestinationModel] `tfsdk:"cloudwatch_logs"`
}

var (
	_ fwflex.Expander  = scraperLoggingDestinationModel{}
	_ fwflex.Flattener = &scraperLoggingDestinationModel{}
)

func (m scraperLoggingDestinationModel) Expand(ctx context.Context) (any, diag.Diagnostics) {
	var diags diag.Diagnostics

	cwData, d := m.CloudwatchLogs.ToPtr(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	if cwData == nil {
		return nil, diags
	}

	return &awstypes.ScraperLoggingDestinationMemberCloudWatchLogs{
		Value: awstypes.CloudWatchLogDestination{
			LogGroupArn: cwData.LogGroupARN.ValueStringPointer(),
		},
	}, diags
}

func (m *scraperLoggingDestinationModel) Flatten(ctx context.Context, v any) diag.Diagnostics {
	var diags diag.Diagnostics

	switch t := v.(type) {
	case awstypes.ScraperLoggingDestinationMemberCloudWatchLogs:
		var data cloudWatchLogDestinationModel
		diags.Append(fwflex.Flatten(ctx, t.Value, &data)...)
		if diags.HasError() {
			return diags
		}
		m.CloudwatchLogs = fwtypes.NewListNestedObjectValueOfPtrMust(ctx, &data)
	}

	return diags
}

func expandScraperComponents(ctx context.Context, src fwtypes.SetOfStringEnum[awstypes.ScraperComponentType]) []awstypes.ScraperComponent {
	if src.IsNull() || src.IsUnknown() {
		return nil
	}

	componentTypes := fwflex.ExpandFrameworkStringyValueSet[awstypes.ScraperComponentType](ctx, src)
	result := make([]awstypes.ScraperComponent, len(componentTypes))
	for i, t := range componentTypes {
		result[i] = awstypes.ScraperComponent{Type: t}
	}

	return result
}
